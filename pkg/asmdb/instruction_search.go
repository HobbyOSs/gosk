package asmdb

import (
	"errors" // Keep only one errors import
	"log"
	"regexp"
	"strings" // Add strings import

	// "fmt" // Remove unused fmt import

	"github.com/HobbyOSs/gosk/pkg/ng_operand" // Use ng_operand
	"github.com/samber/lo"
)

var jsonData []byte

func init() {
	var err error
	jsonData, err = decompressGzip(compressedJSON)
	if err != nil {
		log.Fatalf("Failed to decompress JSON: %v", err)
	}
}

type InstructionDB struct {
}

func NewInstructionDB() *InstructionDB {
	return &InstructionDB{}
}

// FindEncoding は指定された命令とオペランドに対応するエンコーディングを検索します。
func (db *InstructionDB) FindEncoding(opcode string, operands ng_operand.Operands) (*Encoding, error) { // Use ng_operand.Operands
	instr, err := GetInstructionByOpcode(opcode)
	if err != nil {
		return nil, errors.New("instruction not found")
	}

	filteredForms := filterForms(instr.Forms, operands)
	log.Printf("debug: filteredForms length after filterForms: %d", len(filteredForms))

	if len(filteredForms) == 0 {
		return nil, errors.New("no matching encoding found")
	}

	// Prioritize accumulator forms if found
	var allEncodings []*Encoding
	isAccFormFound := lo.SomeBy(filteredForms, func(form InstructionForm) bool {
		// Check if the form itself is an accumulator-specific form (e.g., operands like "ax, imm16")
		// This check might need refinement based on how forms are defined.
		// A simpler check might be if filterEncodings for this form returns accEncodings.
		encs := filterEncodings(form, operands) // Call filterEncodings to check its result
		// Check if the returned encodings are accumulator-specific (ModRM==nil, Immediate!=nil)
		return len(encs) > 0 && encs[0].ModRM == nil && encs[0].Immediate != nil
	})

	if isAccFormFound {
		log.Printf("debug: Accumulator form found, filtering encodings for accumulator.")
		// Only consider encodings from accumulator forms
		allEncodings = lo.FlatMap(filteredForms, func(form InstructionForm, _ int) []*Encoding {
			encs := filterEncodings(form, operands)
			// Only return accumulator-specific encodings
			if len(encs) > 0 && encs[0].ModRM == nil && encs[0].Immediate != nil {
				return encs
			}
			return []*Encoding{} // Return empty if not an accumulator encoding
		})
	} else {
		log.Printf("debug: No accumulator form found, considering all encodings.")
		// Flatten the encodings from all filtered forms if no accumulator form was prioritized
		allEncodings = lo.FlatMap(filteredForms, func(form InstructionForm, _ int) []*Encoding {
			return filterEncodings(form, operands)
		})
	}


	if len(allEncodings) == 0 {
		// This might happen if accumulator forms were found but their encodings were filtered out unexpectedly,
		// or if non-accumulator forms had no suitable encodings.
		log.Printf("error: No suitable encodings found after potential accumulator prioritization.")
		return nil, errors.New("no suitable encoding found after filtering")
	}
	// Log details of each candidate encoding
	log.Printf("debug: Found %d candidate encodings before MinBy:", len(allEncodings)) // Restore this log line
	for i, enc := range allEncodings {
		if enc != nil {
			immSize := 0
			if enc.Immediate != nil {
				immSize = enc.Immediate.Size
			}
			log.Printf("  [%d] Opcode: %s, ModRM: %v, ImmSize: %d", i, enc.Opcode.Byte, enc.ModRM != nil, immSize)
		} else {
			log.Printf("  [%d] nil encoding", i)
		}
	}
	// log.Printf("debug: All candidate encodings before MinBy: %+v", allEncodings) // Keep original log commented out for now

	// Find the smallest encoding size
	// Note: This simple MinBy might select the wrong encoding if sizes are equal (e.g., ADD AX, imm vs ADD r/m, imm).
	// The filterForms logic now prioritizes accumulator forms, which should mitigate this.
	minEncoding := lo.MinBy(allEncodings, func(a, b *Encoding) bool {
		// Add nil checks for safety inside MinBy comparison
		if a == nil || b == nil {
			log.Printf("error: nil encoding passed to MinBy comparison (a=%v, b=%v)", a == nil, b == nil)
			return b == nil // Prefer non-nil
		}
		// Compare based on the encoding's defined size (pass nil options)
		sizeA := a.GetOutputSize(nil)
		sizeB := b.GetOutputSize(nil)

		if sizeA != sizeB {
			return sizeA < sizeB
		}

		// If sizes are equal, prioritize imm8 encoding if applicable
		// Check if Immediate fields are not nil before accessing Size
		immSizeA := 0
		if a.Immediate != nil {
			immSizeA = a.Immediate.Size
		}
		immSizeB := 0
		if b.Immediate != nil {
			immSizeB = b.Immediate.Size
		}

		// Prefer the encoding with imm8 (size 1) if the other is larger
		if immSizeA == 1 && immSizeB > 1 {
			return true // Prefer a (imm8)
		}
		if immSizeB == 1 && immSizeA > 1 {
			return false // Prefer b (imm8)
		}

		// If both are imm8 or neither is imm8 (or sizes differ), maintain original order (or based on sizeA < sizeB already handled)
		return false // Default case if no imm8 preference applies
	})

	// Add nil check for minEncoding before returning
	if minEncoding == nil {
		log.Printf("error: lo.MinBy returned nil encoding")
		return nil, errors.New("failed to find minimum encoding")
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
		e := &form.Encodings[i] // ポインタを取得
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

func filterForms(forms []InstructionForm, operands ng_operand.Operands) []InstructionForm { // Use ng_operand.Operands
	// アキュムレータレジスタを優先的に検索
	accForms := lo.Filter(forms, func(form InstructionForm, _ int) bool {
		// Ensure form.Operands is not nil before dereferencing
		if form.Operands == nil {
			return false
		}
		return matchOperandsWithAccumulator(*form.Operands, operands)
	})
	log.Printf("debug: filteredForms length after matchOperandsWithAccumulator: %d", len(accForms))
	// アキュムレータ形式が見つかった場合は、それを優先して返す
	if len(accForms) > 0 {
		return accForms
	}

	// 通常の検索
	strictForms := lo.Filter(forms, func(form InstructionForm, _ int) bool {
		// Ensure form.Operands is not nil before dereferencing
		if form.Operands == nil {
			return false
		}
		return matchOperandsStrict(*form.Operands, operands)
	})
	log.Printf("debug: filteredForms length after matchOperandsStrict: %d", len(strictForms))
	if len(strictForms) > 0 {
		return strictForms
	}

	// 条件緩和検索（sregをr16として扱う）
	relaxedForms := lo.Filter(forms, func(form InstructionForm, _ int) bool {
		// Ensure form.Operands is not nil before dereferencing
		if form.Operands == nil {
			return false
		}
		return matchOperandsRelaxed(*form.Operands, operands)
	})
	log.Printf("debug: filteredForms length after matchOperandsRelaxed: %d", len(relaxedForms)) // Use relaxedForms
	return relaxedForms                                                                         // Return relaxedForms
}

// GetPrefixSize はプレフィックスバイトのサイズを計算します
func (db *InstructionDB) GetPrefixSize(operands ng_operand.Operands) int { // Use ng_operand.Operands
	size := 0
	// operand size prefix (0x66)のみ必要
	if operands.Require66h() {
		size += 1
	}

	return size
}

// Restore FindMinOutputSize method definition
func (db *InstructionDB) FindMinOutputSize(opcode string, operands ng_operand.Operands) (int, error) { // Use ng_operand.Operands
	encoding, err := db.FindEncoding(opcode, operands)
	if err != nil {
		return 0, err
	}

	options := &OutputSizeOptions{
		ImmSize: operands.DetectImmediateSize(),
	}
	size := encoding.GetOutputSize(options) // Pass options here

	// プレフィックスとオフセットのサイズを加算
	minOutputSize := size + db.GetPrefixSize(operands) + operands.CalcOffsetByteSize()
	log.Printf("debug: [pass1] %s %s = %d\n", opcode, operands.InternalString(), minOutputSize)
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

func hasAccumulator(queryOperands ng_operand.Operands) bool { // Use ng_operand.Operands
	hasAccumulator := lo.SomeBy(queryOperands.InternalStrings(), func(op string) bool {
		matched, _ := regexp.MatchString(`(?i)^(AL|AX|EAX|RAX)$`, op)
		return matched
	})
	return hasAccumulator
}

func matchOperandsStrict(formOperands []Operand, queryOperands ng_operand.Operands) bool { // Use ng_operand.Operands
	if formOperands == nil || len(formOperands) != len(queryOperands.OperandTypes()) {
		return false
	}

	for i, operand := range formOperands {
		queryType := string(queryOperands.OperandTypes()[i]) // Convert OperandType to string
		if operand.Type != queryType {
			return false
		}
	}
	return true
}

func matchOperandsRelaxed(formOperands []Operand, queryOperands ng_operand.Operands) bool { // Use ng_operand.Operands
	if formOperands == nil || len(formOperands) != len(queryOperands.OperandTypes()) {
		return false
	}

	for i, operand := range formOperands {
		queryType := string(queryOperands.OperandTypes()[i]) // Convert OperandType to string
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
