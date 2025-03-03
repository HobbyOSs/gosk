package asmdb

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/HobbyOSs/gosk/pkg/operand"
	"github.com/samber/lo"
	"github.com/tidwall/gjson"
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
	jsonData []byte
}

func NewInstructionDB() *InstructionDB {
	return &InstructionDB{jsonData: jsonData}
}

func (db *InstructionDB) FindInstruction(opcode string) (*Instruction, bool) {
	path := `instructions.` + opcode
	result := gjson.GetBytes(db.jsonData, path)
	if !result.Exists() {
		return nil, false
	}

	var instr Instruction
	if err := json.Unmarshal([]byte(result.Raw), &instr); err != nil {
		return nil, false
	}
	return &instr, true
}

func (db *InstructionDB) FindForms(opcode string, operands operand.Operands) ([]InstructionForm, error) {
	instr, found := db.FindInstruction(opcode)
	if !found {
		return nil, errors.New("instruction not found")
	}

	var matchedForms []InstructionForm
	for _, form := range instr.Forms {
		if matchOperands(form.Operands, operands) {
			matchedForms = append(matchedForms, form)
		}
	}
	return matchedForms, nil
}

func (db *InstructionDB) FindMinOutputSize(opcode string, operands operand.Operands) (int, error) {
	forms, err := db.FindForms(opcode, operands)
	if err != nil {
		return 0, err
	}

	allOutputSize := lo.FlatMap(forms, func(f InstructionForm, _ int) []int {
		return lo.Map(f.Encodings, func(e Encoding, _ int) int {
			options := &OutputSizeOptions{
				ImmSize: operands.DetectImmediateSize(),
			}
			return e.GetOutputSize(options)
		})
	})
	// メモリーアドレス表現にあるoffset値について
	// 機械語サイズの計算をして足し込む
	offsetByteSize := operands.CalcOffsetByteSize()

	// 最小値を取得
	return lo.Min(allOutputSize) + offsetByteSize, nil
}

func matchOperands(formOperands *[]Operand, queryOperands operand.Operands) bool {
	if formOperands == nil || len(*formOperands) != len(queryOperands.OperandTypes()) {
		return false
	}
	for i, operand := range *formOperands {
		if operand.Type != queryOperands.OperandTypes()[i].String() {
			return false
		}
	}
	return true
}
