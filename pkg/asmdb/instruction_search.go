package asmdb

import (
	"errors"
	"log" // Added for debugging
	"regexp"
	"strings"

	"github.com/HobbyOSs/gosk/pkg/ng_operand"
	"github.com/samber/lo"
)

var jsonData []byte

func init() {
	var err error
	jsonData, err = decompressGzip(compressedJSON)
	if err != nil {
		log.Fatalf("JSONの解凍に失敗しました: %v", err)
	}
}

type InstructionDB struct {
}

func NewInstructionDB() *InstructionDB {
	return &InstructionDB{}
}

// FindEncoding は指定された命令とオペランドに対応するエンコーディングを検索します。
// 最適なエンコーディングを選択するための複雑なロジックが含まれています。
// matchAnyImm が true の場合、imm* タイプの比較を緩和します。これは pass1 でのサイズ計算と
// codegen でのエンコーディング選択の挙動を一致させるために導入されました (activeContext.md 参照)。
func (db *InstructionDB) FindEncoding(opcode string, operands ng_operand.Operands, matchAnyImm bool) (*Encoding, error) { // ng_operand.Operands を使用, matchAnyImm パラメータ追加
	instr, err := GetInstructionByOpcode(opcode)
	if err != nil {
		return nil, errors.New("命令が見つかりません")
	}

	// オペランドにマッチするフォームをフィルタリング
	filteredForms := filterForms(instr.Forms, operands, matchAnyImm) // matchAnyImm を渡す

	if len(filteredForms) == 0 {
		return nil, errors.New("一致するエンコーディングが見つかりません")
	}

	// フィルタリングされたフォームから、さらにオペランドに基づいてエンコーディングを絞り込む
	// (例: 間接メモリ参照の場合は ModRM が必須)
	allEncodings := lo.FlatMap(filteredForms, func(form InstructionForm, _ int) []*Encoding {
		return filterEncodings(form, operands)
	})

	if len(allEncodings) == 0 {
		log.Printf("error: フィルタリング後、適切なエンコーディングが見つかりませんでした。opcode=%s, operands=%s", opcode, operands.InternalString())
		return nil, errors.New("フィルタリング後、適切なエンコーディングが見つかりませんでした")
	}

	// 複数のエンコーディング候補から最適なものを選択する
	// 優先順位: 有効性 > サイズ > アキュムレータ形式 > imm8形式 (符号拡張可能時のみ)
	var minEncoding *Encoding
	if isSignExtendable(opcode) {
		// 符号拡張可能な命令 (ADD, SUB など Opcode 83 系) の場合
		minEncoding = lo.MinBy(allEncodings, func(a, b *Encoding) bool {
			return findBestEncodingForSignExtendable(a, b, operands)
		})
	} else {
		// 符号拡張不可能な命令の場合
		minEncoding = lo.MinBy(allEncodings, func(a, b *Encoding) bool {
			return findBestEncodingForNonSignExtendable(a, b, operands)
		})
	}

	// lo.MinBy が候補なしの場合に nil を返す可能性があるためチェック
	if minEncoding == nil {
		log.Printf("error: lo.MinBy が nil エンコーディングを返しました")
		return nil, errors.New("最小エンコーディングの検索に失敗しました")
	}

	return minEncoding, nil
}

// findBestEncodingForSignExtendable は、符号拡張可能な命令 (ADD, SUB など) に最適なエンコーディングを比較します。
// 優先順位: 1. 有効性, 2. サイズ, 3. アキュムレータ形式, 4. imm8形式
func findBestEncodingForSignExtendable(a, b *Encoding, operands ng_operand.Operands) bool {
	// 1. nil チェック: 安全な比較のため
	if a == nil {
		return b != nil
	}
	if b == nil {
		return false
	}

	// 2. 有効性チェック: imm8 形式 (Opcode 83 系) が選択可能か
	//    即値が符号付き8ビットに収まらない場合、imm8 形式は無効とする
	fitsInSignedImm8 := operands.ImmediateValueFitsInSigned8Bits()
	isAImm8 := a.Immediate != nil && a.Immediate.Size == 1
	isBImm8 := b.Immediate != nil && b.Immediate.Size == 1
	isValidA := !(isAImm8 && !fitsInSignedImm8) // imm8形式 かつ 8bitに収まらない場合は無効
	isValidB := !(isBImm8 && !fitsInSignedImm8) // imm8形式 かつ 8bitに収まらない場合は無効

	// lo.Switch を使用して有効性を比較
	// 戻り値が nil でない場合 (どちらか一方または両方が無効な場合) は、その結果を返す
	// 戻り値が nil の場合 (両方有効な場合) は、後続の比較に進む
	validityResult := lo.Switch[struct{ A, B bool }, *bool](struct{ A, B bool }{isValidA, isValidB}).
		Case(struct{ A, B bool }{true, false}, lo.ToPtr(true)).   // Aが有効、Bが無効 -> Aを選択
		Case(struct{ A, B bool }{false, true}, lo.ToPtr(false)).  // Bが有効、Aが無効 -> Bを選択
		Case(struct{ A, B bool }{false, false}, lo.ToPtr(false)). // 両方無効 -> Bを選択 (任意)
		Default(nil)                                              // 両方有効 -> 後続の比較へ

	if validityResult != nil {
		return *validityResult // 無効なエンコーディングがある場合はここで決定
	}

	// --- 両方のエンコーディングが有効な場合 ---

	// 3. サイズ比較: より小さいサイズのエンコーディングを優先
	sizeA := a.GetOutputSize(nil)
	sizeB := b.GetOutputSize(nil)
	if sizeA != sizeB {
		return sizeA < sizeB // サイズが小さい方を優先
	}

	// --- サイズが等しい場合 ---

	// 4. アキュムレータ形式の優先: サイズが同じならアキュムレータ専用形式を優先
	//    (例: ADD AX, imm16 (Opcode 05) vs ADD r/m16, imm16 (Opcode 81 /0))
	isAccA := a.ModRM == nil && a.Immediate != nil // ModRMなし、即値ありはアキュムレータ形式の特徴
	isAccB := b.ModRM == nil && b.Immediate != nil
	accPreferenceResult := lo.Switch[struct{ AccA, AccB bool }, *bool](struct{ AccA, AccB bool }{isAccA, isAccB}).
		Case(struct{ AccA, AccB bool }{true, false}, lo.ToPtr(true)).  // Aがアキュムレータ形式、Bが違う -> Aを選択
		Case(struct{ AccA, AccB bool }{false, true}, lo.ToPtr(false)). // Bがアキュムレータ形式、Aが違う -> Bを選択
		Default(nil)                                                   // 両方アキュムレータ形式、または両方違う -> 後続の比較へ

	if accPreferenceResult != nil {
		return *accPreferenceResult
	}

	// --- サイズが等しく、アキュムレータ形式の状況も同じ場合 ---

	// 5. imm8 形式の優先: 符号拡張可能で即値が8ビットに収まる場合、imm8形式 (Opcode 83 系) を優先
	//    これは、より短いエンコーディングを生成するための重要な最適化。
	if fitsInSignedImm8 {
		imm8PreferenceResult := lo.Switch[struct{ Imm8A, Imm8B bool }, *bool](struct{ Imm8A, Imm8B bool }{isAImm8, isBImm8}).
			Case(struct{ Imm8A, Imm8B bool }{true, false}, lo.ToPtr(true)).  // Aがimm8形式、Bが違う -> Aを選択
			Case(struct{ Imm8A, Imm8B bool }{false, true}, lo.ToPtr(false)). // Bがimm8形式、Aが違う -> Bを選択
			Default(nil)                                                     // 両方imm8形式、または両方違う -> デフォルトへ

		if imm8PreferenceResult != nil {
			return *imm8PreferenceResult
		}
	}

	// 6. デフォルト: 上記の条件で決まらない場合は、任意でBを選択 (比較関数は false を返す)
	return false
}

// findBestEncodingForNonSignExtendable は、符号拡張不可能な命令に最適なエンコーディングを比較します。
// 優先順位: 1. サイズ, 2. アキュムレータ形式
// imm8 形式の優先は行わない。
// 有効性チェックは削除 (MOV r8, imm8 などは常に有効なため)
func findBestEncodingForNonSignExtendable(a, b *Encoding, operands ng_operand.Operands) bool {
	// 1. nil チェック: 安全な比較のため
	if a == nil {
		return b != nil
	}
	if b == nil {
		return false
	}

	// 2. サイズ比較: より小さいサイズのエンコーディングを優先
	sizeA := a.GetOutputSize(nil)
	sizeB := b.GetOutputSize(nil)
	if sizeA != sizeB {
		return sizeA < sizeB
	}

	// --- サイズが等しい場合 ---

	// 3. アキュムレータ形式の優先: サイズが同じならアキュムレータ専用形式を優先
	isAccA := a.ModRM == nil && a.Immediate != nil
	isAccB := b.ModRM == nil && b.Immediate != nil
	accPreferenceResult := lo.Switch[struct{ AccA, AccB bool }, *bool](struct{ AccA, AccB bool }{isAccA, isAccB}).
		Case(struct{ AccA, AccB bool }{true, false}, lo.ToPtr(true)).  // Aがアキュムレータ形式、Bが違う -> Aを選択
		Case(struct{ AccA, AccB bool }{false, true}, lo.ToPtr(false)). // Bがアキュムレータ形式、Aが違う -> Bを選択
		Default(nil)                                                   // 両方アキュムレータ形式、または両方違う -> デフォルトへ

	if accPreferenceResult != nil {
		return *accPreferenceResult
	}

	// --- サイズが等しく、アキュムレータ形式の状況も同じ場合 ---

	// 4. デフォルト: 上記の条件で決まらない場合は、任意でBを選択
	return false
}

// filterEncodings は、オペランドに基づいてエンコーディングをフィルタリングします。
// 最適なエンコーディングの選択 (アキュムレータ優先など) は lo.MinBy に委譲するため、
// ここでは主に ModRM の要不要に基づいたフィルタリングを行います。
func filterEncodings(form InstructionForm, operands ng_operand.Operands) []*Encoding {
	// isAcc := hasAccumulator(operands) // アキュムレータ判定は MinBy で行うためここでは不要
	hasIndirectMem := operands.IsIndirectMemory() // 間接メモリアクセスか？

	// 間接メモリアクセス時のフィルタリング
	filteredEncodings := lo.Filter(form.Encodings, func(e Encoding, _ int) bool {
		// 間接メモリアクセスでは ModRM バイトが必須。
		// ModRM を持たないエンコーディングは除外する。
		// ただし、アキュムレータ専用形式 (ModRMなし、Immediateあり) は常に有効とみなす。
		// (例: ADD AX, imm16 は ModRM 不要だが有効)
		isAccSpecific := e.ModRM == nil && e.Immediate != nil
		if hasIndirectMem && e.ModRM == nil && !isAccSpecific {
			return false // 間接メモリ参照なのにModRMがない形式は除外 (アキュムレータ専用を除く)
		}
		return true // それ以外は有効な候補
	})

	// Encoding 構造体のスライスから、*Encoding のスライスに変換して返す
	// (lo.MinBy がポインタを扱うため)
	return lo.Map(filteredEncodings, func(e Encoding, _ int) *Encoding { return &e })
}

// filterForms は、オペランドにマッチする命令フォームをフィルタリングします。
// 優先順位: 1. アキュムレータ形式, 2. 厳密マッチ, 3. 緩和マッチ (sreg -> r16)
func filterForms(forms []InstructionForm, operands ng_operand.Operands, matchAnyImm bool) []InstructionForm { // ng_operand.Operands を使用, matchAnyImm パラメータ追加
	// 1. アキュムレータレジスタを含む形式を優先的に検索
	accForms := lo.Filter(forms, func(form InstructionForm, _ int) bool {
		// form.Operands が nil の場合、安全に処理をスキップ
		if form.Operands == nil {
			return false
		}
		match := matchOperandsWithAccumulator(*form.Operands, operands, matchAnyImm) // matchAnyImm を渡す
		// log.Printf("debug: [filterForms] acc check: form=%v, query=%s, match=%t", form.Operands, operands.OperandTypes(), match) // Detailed log if needed
		return match
	})

	// アキュムレータ形式が見つかった場合は、それを最優先で返す
	if len(accForms) > 0 {
		return accForms
	}

	// 2. 通常の厳密なオペランドタイプマッチング
	strictForms := lo.Filter(forms, func(form InstructionForm, _ int) bool {
		if form.Operands == nil {
			return false
		}
		match := matchOperandsStrict(*form.Operands, operands, matchAnyImm) // matchAnyImm を渡す
		// log.Printf("debug: [filterForms] strict check: form=%v, query=%s, match=%t", form.Operands, operands.OperandTypes(), match) // Detailed log if needed
		return match
	})

	if len(strictForms) > 0 {
		return strictForms
	}

	// 3. 条件緩和検索: sreg を r16 として扱う
	//    (例: MOV r/m16, Sreg (Opcode 8C) のような命令に対応するため)
	relaxedForms := lo.Filter(forms, func(form InstructionForm, _ int) bool {
		if form.Operands == nil {
			return false
		}
		match := matchOperandsRelaxed(*form.Operands, operands)
		// log.Printf("debug: [filterForms] relaxed check: form=%v, query=%s, match=%t", form.Operands, operands.OperandTypes(), match) // Detailed log if needed
		return match
	})
	return relaxedForms // 緩和マッチの結果を返す (見つからなければ空スライス)
}

// isSignExtendable は、指定された命令が imm8 からの符号拡張 (Opcode 83 系) をサポートするかどうかを返します。
// この判定は、imm8 形式を優先するかどうかの決定に使われます。
func isSignExtendable(opcode string) bool {
	// TODO: より正確なリストに更新する必要があるかもしれません (Intelマニュアル参照)
	switch strings.ToUpper(opcode) {
	case "ADD", "ADC", "SUB", "SBB", "CMP", "AND", "OR", "XOR": // Opcode 83 を持つ主要な命令
		return true
	default:
		return false
	}
}

// GetPrefixSize は必要なプレフィックスバイト (現在はオペランドサイズプレフィックス 0x66 のみ) のサイズを計算します。
func (db *InstructionDB) GetPrefixSize(operands ng_operand.Operands) int { // ng_operand.Operands を使用
	size := 0
	// 16/32ビットモード間でオペランドサイズが異なる場合に 0x66 が必要
	if operands.Require66h() {
		size += 1 // オペランドサイズプレフィックス
	}
	if operands.Require67h() {
		size += 1 // アドレスサイズプレフィックス
	}
	return size
}

// FindMinOutputSize は、pass1 での LOC (Location Counter) 計算のために、
// 指定された命令とオペランドに対する最小の出力バイトサイズを計算します。
// codegen が選択するであろう最適なエンコーディング (imm8 優先など) を考慮する必要があります。
func (db *InstructionDB) FindMinOutputSize(opcode string, operands ng_operand.Operands) (int, error) { // ng_operand.Operands を使用
	// codegen と同じロジックで最適なエンコーディングを検索する
	// matchAnyImm = true にすることで、即値が小さい場合に imm8 形式が考慮されるようにする
	// これにより、pass1 と codegen のサイズ解釈のずれを防ぐ (activeContext.md 参照)
	encoding, err := db.FindEncoding(opcode, operands, true)
	if err != nil {
		// matchAnyImm=true で見つからない場合、フォールバックとして false で再試行
		// (基本的には true で見つかるはずだが、予期せぬケースへの対応)
		log.Printf("warn: FindEncoding(matchAnyImm=true) failed for %s %s, retrying with false: %v", opcode, operands.InternalString(), err)
		encoding, err = db.FindEncoding(opcode, operands, false)
		if err != nil {
			log.Printf("error: FindEncoding failed even with matchAnyImm=false for %s %s: %v", opcode, operands.InternalString(), err)
			return 0, err // フォールバックでも見つからなければエラー
		}
	}

	// 選択された最適なエンコーディングの基本サイズを取得
	// (例: imm8 なら即値1バイト、imm16/32 なら即値2/4バイトとして計算される)
	size := encoding.GetOutputSize(nil) // オプションは不要

	// 基本サイズにプレフィックス、オフセット、SIB バイトのサイズを加算
	sibSize := operands.CalcSibByteSize() // Use the interface method
	minOutputSize := size + db.GetPrefixSize(operands) + operands.CalcOffsetByteSize() + sibSize
	// デバッグ用に計算結果をログ出力 (SIB サイズも含む)
	log.Printf("debug: [pass1] %s %s = %d (base:%d, prefix:%d, offset:%d, sib:%d)\n",
		opcode, operands.InternalString(), minOutputSize,
		size, db.GetPrefixSize(operands), operands.CalcOffsetByteSize(), sibSize)
	return minOutputSize, nil
}

// matchOperandsWithAccumulator は、問い合わせオペランド (queryOperands) にアキュムレータ (AL/AX/EAX) が含まれ、
// かつ命令フォームのオペランド (formOperands) がそれにマッチするかどうかを判定します。
// アキュムレータ専用形式 (例: ADD AX, imm16 (Opcode 05)) を優先的にマッチさせるために使用されます。
func matchOperandsWithAccumulator(formOperands []Operand, queryOperands ng_operand.Operands, matchAnyImm bool) bool { // matchAnyImm パラメータ追加
	queryTypes := queryOperands.OperandTypes()
	// 問い合わせオペランドにアキュムレータが含まれていない、またはオペランド数が不一致なら false
	if !hasAccumulator(queryOperands) || len(formOperands) != len(queryTypes) {
		return false
	}

	// すべてのオペランドが条件を満たすかチェック
	return lo.EveryBy(formOperands, func(formOp Operand) bool {
		i := lo.IndexOf(formOperands, formOp) // 対応する問い合わせオペランドのインデックスを取得
		queryType := string(queryTypes[i])
		formType := formOp.Type

		// 1. タイプが完全に一致する場合 (例: form="r16", query="r16")
		if formType == queryType {
			return true
		}
		// 2. form がアキュムレータで、query が対応する汎用レジスタの場合
		//    (例: form="ax", query="r16" はマッチ)
		if (formType == "al" && queryType == "r8") ||
			(formType == "ax" && queryType == "r16") ||
			(formType == "eax" && queryType == "r32") {
			return true
		}
		// 3. アキュムレータ以外のオペランドの比較 (主に即値)
		//    matchAnyImm が true の場合、imm* タイプ同士はサイズ違いでもマッチとみなす
		//    (例: form="imm16", query="imm8" は matchAnyImm=true ならマッチ)
		isFormImm := strings.HasPrefix(formType, "imm")
		isQueryImm := strings.HasPrefix(queryType, "imm")
		if matchAnyImm && isFormImm && isQueryImm {
			return true
		}
		// 上記以外でタイプが不一致なら false
		return false
	})
}

// hasAccumulator は、オペランドリストにアキュムレータレジスタ (AL/AX/EAX/RAX) が含まれるかどうかを判定します。
func hasAccumulator(queryOperands ng_operand.Operands) bool { // ng_operand.Operands を使用
	// オペランドの文字列表現をチェック
	return lo.SomeBy(queryOperands.InternalStrings(), func(op string) bool {
		// 大文字小文字を区別せずに完全一致で判定
		matched, _ := regexp.MatchString(`(?i)^(AL|AX|EAX|RAX)$`, op)
		return matched
	})
}

// matchOperandsStrict は、命令フォームのオペランド (formOperands) と
// 問い合わせオペランド (queryOperands) のタイプが厳密に一致するかどうかを判定します。
// matchAnyImm が true の場合、imm* タイプ同士はサイズ違いでも一致とみなします。
func matchOperandsStrict(formOperands []Operand, queryOperands ng_operand.Operands, matchAnyImm bool) bool { // ng_operand.Operands を使用, matchAnyImm パラメータ追加
	queryTypes := queryOperands.OperandTypes() // 問い合わせオペランドのタイプを取得
	// オペランド数が異なる場合は false
	if formOperands == nil || len(formOperands) != len(queryTypes) {
		return false
	}

	// すべてのオペランドタイプが一致するかチェック
	return lo.EveryBy(formOperands, func(formOp Operand) bool {
		i := lo.IndexOf(formOperands, formOp)
		queryType := string(queryTypes[i]) // OperandType を string に変換
		formType := formOp.Type

		// 1. タイプが完全に一致する場合
		if formType == queryType {
			return true
		}
		// 2. matchAnyImm が true で、両方が imm* タイプの場合
		isFormImm := strings.HasPrefix(formType, "imm")
		isQueryImm := strings.HasPrefix(queryType, "imm")
		if matchAnyImm && isFormImm && isQueryImm {
			return true
		}
		// 上記以外は不一致
		return false
	})
}

// matchOperandsRelaxed は、オペランドタイプの一致を判定しますが、
// 問い合わせオペランドの "sreg" (セグメントレジスタ) をフォームの "r16" としても
// マッチするように条件を緩和します。これは MOV r/m16, Sreg (Opcode 8C) などに対応するためです。
func matchOperandsRelaxed(formOperands []Operand, queryOperands ng_operand.Operands) bool { // ng_operand.Operands を使用
	queryTypes := queryOperands.OperandTypes()
	// オペランド数が異なる場合は false
	if formOperands == nil || len(formOperands) != len(queryTypes) {
		return false
	}

	// すべてのオペランドが条件を満たすかチェック
	return lo.EveryBy(formOperands, func(formOp Operand) bool {
		i := lo.IndexOf(formOperands, formOp)
		queryType := string(queryTypes[i]) // OperandType を string に変換
		formType := formOp.Type

		// 1. タイプが完全に一致する場合
		if formType == queryType {
			return true
		}
		// 2. 緩和条件: query が "sreg" で form が "r16" の場合
		if queryType == "sreg" && formType == "r16" {
			return true // sreg は r16 として扱う
		}
		// 上記以外は不一致
		return false
	})
}
