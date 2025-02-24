package asmdb

import (
	"encoding/json"
	"errors"

	"github.com/tidwall/gjson"
)

type InstructionDB struct {
	jsonData []byte
}

func NewInstructionDB(data []byte) *InstructionDB {
	return &InstructionDB{jsonData: data}
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

func (db *InstructionDB) FindForms(opcode string, operands []string) ([]InstructionForm, error) {
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

func matchOperands(formOperands *[]Operand, queryOperands []string) bool {
	if formOperands == nil || len(*formOperands) != len(queryOperands) {
		return false
	}
	for i, operand := range *formOperands {
		if operand.Type != queryOperands[i] {
			return false
		}
	}
	return true
}
