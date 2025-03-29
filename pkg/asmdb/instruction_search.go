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
func (db *InstructionDB) FindEncoding(opcode string, operands ng_operand.Operands) (*Encoding, error) { // ng_operand.Operands を使用
	instr, err := GetInstructionByOpcode(opcode)
	if err != nil {
		return nil, errors.New("命令が見つかりません")
	}

	filteredForms := filterForms(instr.Forms, operands)

	if len(filteredForms) == 0 {
		return nil, errors.New("一致するエンコーディングが見つかりません")
	}

	// アキュムレータ形式が見つかった場合、それを優先する
	var allEncodings []*Encoding
	isAccFormFound := lo.SomeBy(filteredForms, func(form InstructionForm) bool {
		// フォーム自体がアキュムレータ固有の形式（例：オペランドが "ax, imm16"）かどうかを確認
		// このチェックはフォームの定義方法に基づいて改良が必要かもしれない
		// より簡単なチェックは、このフォームに対する filterEncodings が accEncodings を返すかどうか
		encs := filterEncodings(form, operands) // filterEncodings を呼び出して結果を確認
		// 返されたエンコーディングがアキュムレータ固有（ModRM==nil, Immediate!=nil）かどうかを確認
		return len(encs) > 0 && encs[0].ModRM == nil && encs[0].Immediate != nil
	})

	if isAccFormFound {
		// アキュムレータ形式からのみエンコーディングを考慮する
		allEncodings = lo.FlatMap(filteredForms, func(form InstructionForm, _ int) []*Encoding {
			encs := filterEncodings(form, operands)
			// アキュムレータ固有のエンコーディングのみを返す
			if len(encs) > 0 && encs[0].ModRM == nil && encs[0].Immediate != nil {
				return encs
			}
			return []*Encoding{} // アキュムレータエンコーディングでない場合は空を返す
		})
	} else {
		// アキュムレータ形式が優先されなかった場合、フィルタリングされたすべてのフォームからエンコーディングをフラット化する
		allEncodings = lo.FlatMap(filteredForms, func(form InstructionForm, _ int) []*Encoding {
			return filterEncodings(form, operands)
		})
	}

	if len(allEncodings) == 0 {
		// これは、アキュムレータ形式が見つかったが、そのエンコーディングが予期せずフィルタリングされた場合、
		// または非アキュムレータ形式に適したエンコーディングがなかった場合に発生する可能性がある
		log.Printf("error: アキュムレータの優先順位付けの後、適切なエンコーディングが見つかりませんでした。")
		return nil, errors.New("フィルタリング後、適切なエンコーディングが見つかりませんでした")
	}

	// 最小のエンコーディングサイズを見つける
	minEncoding := lo.MinBy(allEncodings, func(a, b *Encoding) bool {
		// 安全のため MinBy 比較内で nil チェックを追加
		if a == nil || b == nil {
			log.Printf("error: nil エンコーディングが MinBy 比較に渡されました (a=%v, b=%v)", a == nil, b == nil)
			return b == nil // nil でない方を優先
		}
		// エンコーディングの定義済みサイズに基づいて比較（nil オプションを渡す）
		sizeA := a.GetOutputSize(nil)
		sizeB := b.GetOutputSize(nil)

		if sizeA != sizeB {
			return sizeA < sizeB
		}

		// サイズが等しい場合、該当すれば imm8 エンコーディングを優先する
		// Size にアクセスする前に Immediate フィールドが nil でないことを確認
		immSizeA := 0
		if a.Immediate != nil {
			immSizeA = a.Immediate.Size
		}
		immSizeB := 0
		if b.Immediate != nil {
			immSizeB = b.Immediate.Size
		}

		// 他方より大きい場合、imm8（サイズ1）を持つエンコーディングを優先する
		if immSizeA == 1 && immSizeB > 1 {
			return true // a (imm8) を優先
		}
		if immSizeB == 1 && immSizeA > 1 {
			return false // b (imm8) を優先
		}

		// 両方が imm8 であるか、どちらも imm8 でない場合（またはサイズが異なる場合）、元の順序を維持する（または既に処理された sizeA < sizeB に基づく）
		return false // imm8 優先が適用されない場合のデフォルトケース
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
		hasDirectMem := operands.IsDirectMemory()
		hasIndirectMem := operands.IsIndirectMemory()

		for _, e := range otherEncodings {
			// 直接アドレッシングではModRMが不要なので、ModRMを持つエンコーディングは除外
			if hasDirectMem && e.ModRM != nil {
				continue
			}
			// 間接アドレッシングではModRMが必要なので、ModRMを持たないエンコーディングは除外
			if hasIndirectMem && e.ModRM == nil {
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

func filterForms(forms []InstructionForm, operands ng_operand.Operands) []InstructionForm { // ng_operand.Operands を使用
	// アキュムレータレジスタを優先的に検索
	accForms := lo.Filter(forms, func(form InstructionForm, _ int) bool {
		// 逆参照する前に form.Operands が nil でないことを確認
		if form.Operands == nil {
			return false
		}
		return matchOperandsWithAccumulator(*form.Operands, operands)
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
		return matchOperandsStrict(*form.Operands, operands)
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

// GetPrefixSize はプレフィックスバイトのサイズを計算します
func (db *InstructionDB) GetPrefixSize(operands ng_operand.Operands) int { // ng_operand.Operands を使用
	size := 0
	// オペランドサイズプレフィックス (0x66) のみ必要
	if operands.Require66h() {
		size += 1
	}
	return size
}

// Restore FindMinOutputSize method definition
func (db *InstructionDB) FindMinOutputSize(opcode string, operands ng_operand.Operands) (int, error) { // ng_operand.Operands を使用
	encoding, err := db.FindEncoding(opcode, operands)
	if err != nil {
		return 0, err
	}

	options := &OutputSizeOptions{
		ImmSize: operands.DetectImmediateSize(),
	}
	size := encoding.GetOutputSize(options) // ここで options を渡す

	// プレフィックスとオフセットのサイズを加算
	minOutputSize := size + db.GetPrefixSize(operands) + operands.CalcOffsetByteSize()
	log.Printf("debug: [pass1] %s %s = %d\n", opcode, operands.InternalString(), minOutputSize) // このログは保持
	return minOutputSize, nil
}

// matchOperandsWithAccumulator は、queryOperands にアキュムレータが含まれており、
// formOperands がそれにマッチするかどうかを判定します。
// アキュムレータ専用形式 (例: ADD AX, imm16) を優先的にマッチさせます。
func matchOperandsWithAccumulator(formOperands []Operand, queryOperands ng_operand.Operands) bool {
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
			// 即値タイプの比較を緩和: imm, imm8, imm16, imm32, imm64 は互換性があるとみなす
			isFormImm := strings.HasPrefix(formType, "imm")
			isQueryImm := strings.HasPrefix(queryType, "imm")
			if isFormImm && isQueryImm {
				continue // 両方とも即値タイプならOKとする
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

func matchOperandsStrict(formOperands []Operand, queryOperands ng_operand.Operands) bool { // ng_operand.Operands を使用
	if formOperands == nil || len(formOperands) != len(queryOperands.OperandTypes()) {
		return false
	}
	for i, operand := range formOperands {
		queryType := string(queryOperands.OperandTypes()[i]) // OperandType を string に変換
		if operand.Type != queryType {
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
