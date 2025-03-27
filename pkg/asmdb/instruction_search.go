package asmdb

import (
	"errors"
	"log"
	"regexp"

	"github.com/HobbyOSs/gosk/pkg/operand" // Restored operand import
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
func (db *InstructionDB) FindEncoding(opcode string, operands operand.Operands) (*Encoding, error) {
	instr, err := GetInstructionByOpcode(opcode)
	if err != nil {
		return nil, errors.New("instruction not found")
	}

	filteredForms := filterForms(instr.Forms, operands)
	log.Printf("debug: filteredForms length after filterForms: %d", len(filteredForms))

	if len(filteredForms) == 0 {
		return nil, errors.New("no matching encoding found")
	}

	// Flatten the encodings from all filtered forms
	allEncodings := lo.FlatMap(filteredForms, func(form InstructionForm, _ int) []*Encoding {
		return filterEncodings(form, operands)
	})

	// Find the smallest encoding size
	minEncoding := lo.MinBy(allEncodings, func(a, b *Encoding) bool {
		optionsA := &OutputSizeOptions{ImmSize: operands.DetectImmediateSize()}
		optionsB := &OutputSizeOptions{ImmSize: operands.DetectImmediateSize()}

		return a.GetOutputSize(optionsA) < b.GetOutputSize(optionsB)
	})

	return minEncoding, nil
}

// ModRM 要否によるフィルタリング(accがあるときのみのルール)
func filterEncodings(form InstructionForm, operands operand.Operands) []*Encoding {

	encodings := make([]*Encoding, 0)

	if !hasAccumulator(operands) {
		for _, e := range form.Encodings {
			encodings = append(encodings, &e)
		}
		return encodings
	}

	hasDirectMem := false
	hasIndirectMem := false

	for _, op := range operands.ParsedOperands() {
		// 直接アドレッシングではModRMが不要
		if op.DirectMem != nil {
			hasDirectMem = true
		}
		// 間接アドレッシングではModRMが必要
		if op.IndirectMem != nil {
			hasIndirectMem = true
		}
	}

	for _, e := range form.Encodings {
		// 直接アドレッシングではModRMが不要なのでエンコーディングとしては除外
		if hasDirectMem && e.ModRM != nil {
			continue
		}
		// 間接アドレッシングではModRMが必要なのでエンコーディングとしては除外
		if hasIndirectMem && e.ModRM == nil {
			continue
		}
		encodings = append(encodings, &e)
	}
	return encodings
}

func filterForms(forms []InstructionForm, operands operand.Operands) []InstructionForm {
	var filteredForms []InstructionForm

	// アキュムレータレジスタを優先的に検索
	filteredForms = append(filteredForms, lo.Filter(forms, func(form InstructionForm, _ int) bool {
		return matchOperandsWithAccumulator(*form.Operands, operands)
	})...)
	log.Printf("debug: filteredForms length after matchOperandsWithAccumulator: %d", len(filteredForms))

	// 通常の検索
	_forms := lo.Filter(forms, func(form InstructionForm, _ int) bool {
		return matchOperandsStrict(*form.Operands, operands)
	})
	filteredForms = append(filteredForms, _forms...)
	if len(filteredForms) > 0 {
		return filteredForms
	}

	// 条件緩和検索（sregをr16として扱う）
	filteredForms = lo.Filter(forms, func(form InstructionForm, _ int) bool {
		return matchOperandsRelaxed(*form.Operands, operands)
	})
	log.Printf("debug: filteredForms length after matchOperandsRelaxed: %d", len(filteredForms))
	return filteredForms
}

// GetPrefixSize はプレフィックスバイトのサイズを計算します
func (db *InstructionDB) GetPrefixSize(operands operand.Operands) int {
	size := 0

	// operand size prefix (0x66)のみ必要
	if operands.Require66h() {
		size += 1
	}

	return size
}

func (db *InstructionDB) FindMinOutputSize(opcode string, operands operand.Operands) (int, error) {
	encoding, err := db.FindEncoding(opcode, operands)
	if err != nil {
		return 0, err
	}

	options := &OutputSizeOptions{
		ImmSize: operands.DetectImmediateSize(),
	}
	size := encoding.GetOutputSize(options)

	// プレフィックスとオフセットのサイズを加算
	minOutputSize := size + db.GetPrefixSize(operands) + operands.CalcOffsetByteSize()
	log.Printf("debug: [pass1] %s %s = %d\n", opcode, operands.InternalString(), minOutputSize)
	return minOutputSize, nil
}

func matchOperandsWithAccumulator(formOperands []Operand, queryOperands operand.Operands) bool {
	// formOperandsにアキュムレータレジスタが含まれているかチェック
	if !hasAccumulator(queryOperands) {
		return false
	}

	// アキュムレータレジスタを優先的にマッチングするロジック
	for i, operand := range formOperands {
		queryType := queryOperands.OperandTypes()[i].String()
		if operand.Type != queryType {
			// アキュムレータレジスタの場合、特定の条件でマッチングを試みる
			if (operand.Type == "al" && queryType == "r8") ||
				(operand.Type == "ax" && queryType == "r16") ||
				(operand.Type == "eax" && queryType == "r32") {
				continue
			}
			return false
		}
	}
	return true
}

func hasAccumulator(queryOperands operand.Operands) bool {
	hasAccumulator := lo.SomeBy(queryOperands.InternalStrings(), func(op string) bool {
		matched, _ := regexp.MatchString(`(?i)^(AL|AX|EAX|RAX)$`, op)
		return matched
	})
	return hasAccumulator
}

func matchOperandsStrict(formOperands []Operand, queryOperands operand.Operands) bool {
	if formOperands == nil || len(formOperands) != len(queryOperands.OperandTypes()) {
		return false
	}

	for i, operand := range formOperands {
		queryType := queryOperands.OperandTypes()[i].String()
		if operand.Type != queryType {
			return false
		}
	}
	return true
}

func matchOperandsRelaxed(formOperands []Operand, queryOperands operand.Operands) bool {
	if formOperands == nil || len(formOperands) != len(queryOperands.OperandTypes()) {
		return false
	}

	for i, operand := range formOperands {
		queryType := queryOperands.OperandTypes()[i].String()
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
