package asmdb

import (
	"errors"
	"log"
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
// matchAnyImm が true の場合、imm* タイプの比較を緩和します。
func (db *InstructionDB) FindEncoding(opcode string, operands ng_operand.Operands, matchAnyImm bool) (*Encoding, error) { // ng_operand.Operands を使用, matchAnyImm パラメータ追加
	instr, err := GetInstructionByOpcode(opcode)
	if err != nil {
		return nil, errors.New("命令が見つかりません")
	}

	filteredForms := filterForms(instr.Forms, operands, matchAnyImm) // matchAnyImm を渡す

	if len(filteredForms) == 0 {
		return nil, errors.New("一致するエンコーディングが見つかりません")
	}

	// フィルタリングされたフォームからすべての可能なエンコーディングを取得
	// filterEncodings はアキュムレータ形式とその他の形式を区別せずに返すようになった
	allEncodings := lo.FlatMap(filteredForms, func(form InstructionForm, _ int) []*Encoding {
		return filterEncodings(form, operands)
	})

	if len(allEncodings) == 0 {
		log.Printf("error: フィルタリング後、適切なエンコーディングが見つかりませんでした。opcode=%s, operands=%s", opcode, operands.InternalString())
		return nil, errors.New("フィルタリング後、適切なエンコーディングが見つかりませんでした")
	}

	// 最小のエンコーディングサイズを見つける (優先順位付けロジックを再修正)
	minEncoding := lo.MinBy(allEncodings, func(a, b *Encoding) bool {
		// 1. Nil checks
		if a == nil {
			return b != nil
		}
		if b == nil {
			return false
		}

		// 2. Validity check (imm8 vs signed 8-bit value)
		signExtendable := isSignExtendable(opcode)
		fitsInSignedImm8 := operands.ImmediateValueFitsInSigned8Bits()
		isAImm8 := a.Immediate != nil && a.Immediate.Size == 1
		isBImm8 := b.Immediate != nil && b.Immediate.Size == 1
		isValidA := !(isAImm8 && !fitsInSignedImm8) // Invalid if imm8 form but value doesn't fit signed 8-bit
		isValidB := !(isBImm8 && !fitsInSignedImm8) // Invalid if imm8 form but value doesn't fit signed 8-bit

		if isValidA && !isValidB {
			return true
		} // A valid, B invalid -> A wins
		if !isValidA && isValidB {
			return false
		} // B valid, A invalid -> B wins
		if !isValidA && !isValidB {
			return false
		} // Both invalid -> B wins (arbitrary)

		// --- Both A and B are valid ---

		// 3. Size comparison
		sizeA := a.GetOutputSize(nil)
		sizeB := b.GetOutputSize(nil)
		if sizeA != sizeB {
			return sizeA < sizeB // Smaller size wins
		}

		// --- Sizes are equal ---

		// 4. Accumulator preference (if sizes are equal)
		isAccA := a.ModRM == nil && a.Immediate != nil
		isAccB := b.ModRM == nil && b.Immediate != nil
		if isAccA && !isAccB {
			return true
		} // A is Acc, B is not -> A wins
		if !isAccA && isAccB {
			return false
		} // B is Acc, A is not -> B wins

		// --- Sizes are equal, and either both are Acc or neither are Acc ---

		// 5. Imm8 preference (only if sign-extendable and fits signed 8-bit)
		if signExtendable && fitsInSignedImm8 {
			if isAImm8 && !isBImm8 {
				return true
			} // A is imm8, B is not -> A wins
			if !isAImm8 && isBImm8 {
				return false
			} // B is imm8, A is not -> B wins
		}

		// 6. Default (sizes equal, acc status same, imm8 pref doesn't apply)
		return false // Arbitrarily choose B
	})

	// 返す前に minEncoding の nil チェックを追加
	if minEncoding == nil {
		log.Printf("error: lo.MinBy が nil エンコーディングを返しました")
		return nil, errors.New("最小エンコーディングの検索に失敗しました")
	}

	return minEncoding, nil
}

// filterEncodings は、オペランドに基づいてエンコーディングをフィルタリングします。
// アキュムレータを使用するエンコーディングを優先します。
func filterEncodings(form InstructionForm, operands ng_operand.Operands) []*Encoding {
	accEncodings := make([]*Encoding, 0)
	otherEncodings := make([]*Encoding, 0)
	filteredOtherEncodings := make([]*Encoding, 0)

	isAcc := hasAccumulator(operands)

	// エンコーディングをアキュムレータ用とその他に分類
	for i := range form.Encodings {
		e := &form.Encodings[i]
		// アキュムレータを使用し、ModRMが不要なエンコーディングを優先候補とする
		// (例: ADD AX, imm16 (opcode 0x05) はModRM不要)
		if isAcc && e.ModRM == nil && e.Immediate != nil {
			accEncodings = append(accEncodings, e)
		} else {
			otherEncodings = append(otherEncodings, e)
		}
	}

	// アキュムレータを使用しない場合、またはアキュムレータ用以外のエンコーディングに対するフィルタリング
	if !isAcc {
		// アキュムレータを使用しない場合は、すべてのエンコーディングをそのまま返す
		// (エンコーディングの選択は lo.MinBy に任せる)
		// TODO: アキュムレータ以外の場合もModRMフィルタリングが必要か再検討
		filteredOtherEncodings = otherEncodings
	} else {
		// アキュムレータを使用する場合、その他エンコーディングにModRMフィルタリングを適用
		// hasDirectMem := operands.IsDirectMemory() // Unused variable
		hasIndirectMem := operands.IsIndirectMemory()

		for _, e := range otherEncodings {
			// 直接アドレッシングではModRMが不要なので、ModRMを持つエンコーディングは除外
			// -> このロジックは誤り。C6/C7命令は直接アドレスでもModRMを使う。
			// if hasDirectMem && e.ModRM != nil {
			// 	continue
			// }
			// 間接アドレッシングではModRMが必要なので、ModRMを持たないエンコーディングは除外
			// (アキュムレータ専用形式は除く)
			isAccSpecific := isAcc && e.ModRM == nil && e.Immediate != nil
			if hasIndirectMem && e.ModRM == nil && !isAccSpecific {
				continue
			}
			filteredOtherEncodings = append(filteredOtherEncodings, e)
		}
	}

	// アキュムレータ用エンコーディングが見つかった場合は、それを最優先で返す
	if len(accEncodings) > 0 {
		return accEncodings
	}

	// アキュムレータ用エンコーディングが見つからなかった場合は、フィルタリングされたその他エンコーディングを返す
	return filteredOtherEncodings
}

func filterForms(forms []InstructionForm, operands ng_operand.Operands, matchAnyImm bool) []InstructionForm { // ng_operand.Operands を使用, matchAnyImm パラメータ追加
	// アキュムレータレジスタを優先的に検索
	accForms := lo.Filter(forms, func(form InstructionForm, _ int) bool {
		// 逆参照する前に form.Operands が nil でないことを確認
		if form.Operands == nil {
			return false
		}
		return matchOperandsWithAccumulator(*form.Operands, operands, matchAnyImm) // matchAnyImm を渡す
	})
	// アキュムレータ形式が見つかった場合は、それを優先して返す
	if len(accForms) > 0 {
		return accForms
	}

	// 通常の検索
	strictForms := lo.Filter(forms, func(form InstructionForm, _ int) bool {
		// 逆参照する前に form.Operands が nil でないことを確認
		if form.Operands == nil {
			return false
		}
		return matchOperandsStrict(*form.Operands, operands, matchAnyImm) // matchAnyImm を渡す
	})
	if len(strictForms) > 0 {
		return strictForms
	}

	// 条件緩和検索（sregをr16として扱う）
	relaxedForms := lo.Filter(forms, func(form InstructionForm, _ int) bool {
		// 逆参照する前に form.Operands が nil でないことを確認
		if form.Operands == nil {
			return false
		}
		return matchOperandsRelaxed(*form.Operands, operands)
	})
	return relaxedForms // relaxedForms を返す
}

// isSignExtendable は、指定された命令が imm8 からの符号拡張をサポートするかどうかを返します。
// (例: ADD, SUB, CMP など Opcode 83 系)
func isSignExtendable(opcode string) bool {
	// TODO: より正確なリストに更新する必要があるかもしれません
	switch strings.ToUpper(opcode) {
	case "ADD", "ADC", "SUB", "SBB", "CMP", "AND", "OR", "XOR":
		return true
	default:
		return false
	}
}

// GetPrefixSize はプレフィックスバイトのサイズを計算します
func (db *InstructionDB) GetPrefixSize(operands ng_operand.Operands) int { // ng_operand.Operands を使用
	size := 0
	// オペランドサイズプレフィックス (0x66) のみ必要
	if operands.Require66h() {
		size += 1
	}
	return size
}

// FindMinOutputSize は、指定された命令とオペランドに対して可能な最小の出力サイズを計算します。
// codegen が選択するであろう最適なエンコーディングを考慮します。
func (db *InstructionDB) FindMinOutputSize(opcode string, operands ng_operand.Operands) (int, error) { // ng_operand.Operands を使用
	// codegen と同様に、最適なエンコーディングを見つけるために matchAnyImm = true で検索
	// これにより、即値が小さい場合に imm8 形式が考慮される
	encoding, err := db.FindEncoding(opcode, operands, true)
	if err != nil {
		// フォールバックとして、より厳密なマッチングを試みる (以前の動作に近い)
		// これが必要になるケースは稀だが、念のため
		log.Printf("warn: FindEncoding(matchAnyImm=true) failed for %s %s, retrying with false: %v", opcode, operands.InternalString(), err)
		encoding, err = db.FindEncoding(opcode, operands, false)
		if err != nil {
			log.Printf("error: FindEncoding failed even with matchAnyImm=false for %s %s: %v", opcode, operands.InternalString(), err)
			return 0, err
		}
	}

	// encoding.GetOutputSize は、エンコーディング自体の定義に基づいてサイズを計算する
	// (例: imm8 エンコーディングなら即値は1バイト、imm16 なら2バイト)
	// DetectImmediateSize() の結果はここでは不要 (FindEncoding が最適なものを選択済みのため)
	size := encoding.GetOutputSize(nil) // options は不要

	// プレフィックスとオフセットのサイズを加算
	minOutputSize := size + db.GetPrefixSize(operands) + operands.CalcOffsetByteSize()
	// ログに選択されたエンコーディング情報を追加してデバッグしやすくする (Stringメソッドがないため一旦削除)
	log.Printf("debug: [pass1] %s %s = %d\n", opcode, operands.InternalString(), minOutputSize)
	return minOutputSize, nil
}

// matchOperandsWithAccumulator は、queryOperands にアキュムレータが含まれており、
// formOperands がそれにマッチするかどうかを判定します。
// アキュムレータ専用形式 (例: ADD AX, imm16) を優先的にマッチさせます。
func matchOperandsWithAccumulator(formOperands []Operand, queryOperands ng_operand.Operands, matchAnyImm bool) bool { // matchAnyImm パラメータ追加
	// queryOperands にアキュムレータが含まれていない場合は false
	if !hasAccumulator(queryOperands) {
		return false
	}
	// formOperands と queryOperands の数が一致しない場合は false
	if len(formOperands) != len(queryOperands.OperandTypes()) {
		return false
	}
	// 各オペランドを比較
	for i, formOp := range formOperands {
		queryType := string(queryOperands.OperandTypes()[i])
		formType := formOp.Type
		// タイプが完全に一致する場合はOK
		if formType == queryType {
			continue
		}
		// form がアキュムレータで、query が対応する汎用レジスタの場合もOK
		// (例: form="ax", query="r16" はOK)
		if (formType == "al" && queryType == "r8") ||
			(formType == "ax" && queryType == "r16") ||
			(formType == "eax" && queryType == "r32") {
			continue
		}
		// アキュムレータ以外のオペランドの比較
		if formType != queryType {
			// matchAnyImm が true の場合、imm* タイプ同士は常にマッチ
			isFormImm := strings.HasPrefix(formType, "imm")
			isQueryImm := strings.HasPrefix(queryType, "imm")
			if matchAnyImm && isFormImm && isQueryImm {
				continue
			}
			// それ以外のタイプが不一致なら false
			return false
		}
	}
	// すべてのオペランドがマッチした場合
	return true
}

func hasAccumulator(queryOperands ng_operand.Operands) bool { // ng_operand.Operands を使用
	hasAccumulator := lo.SomeBy(queryOperands.InternalStrings(), func(op string) bool {
		matched, _ := regexp.MatchString(`(?i)^(AL|AX|EAX|RAX)$`, op)
		return matched
	})
	return hasAccumulator
}

func matchOperandsStrict(formOperands []Operand, queryOperands ng_operand.Operands, matchAnyImm bool) bool { // ng_operand.Operands を使用, matchAnyImm パラメータ追加
	queryTypes := queryOperands.OperandTypes() // Get types once
	if formOperands == nil || len(formOperands) != len(queryTypes) {
		return false
	}
	for i, formOp := range formOperands {
		queryType := string(queryTypes[i]) // OperandType を string に変換
		formType := formOp.Type
		if formType != queryType {
			// matchAnyImm が true の場合、imm* タイプ同士は常にマッチ
			isFormImm := strings.HasPrefix(formType, "imm")
			isQueryImm := strings.HasPrefix(queryType, "imm")
			if matchAnyImm && isFormImm && isQueryImm {
				continue
			}
			return false
		}
	}
	return true
}

func matchOperandsRelaxed(formOperands []Operand, queryOperands ng_operand.Operands) bool { // ng_operand.Operands を使用
	if formOperands == nil || len(formOperands) != len(queryOperands.OperandTypes()) {
		return false
	}
	for i, operand := range formOperands {
		queryType := string(queryOperands.OperandTypes()[i]) // OperandType を string に変換
		if operand.Type != queryType {
			// 条件が緩和された場合; sregはr16としても一致を試みる
			if queryType == "sreg" && operand.Type == "r16" {
				continue // sregはr16として扱う
			}
			return false
		}
	}
	return true
}
