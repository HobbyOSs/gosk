package asmdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestX86Instructions(t *testing.T) {
	instrs := X86Instructions()
	assert.NotNil(t, instrs)
	assert.NotEmpty(t, instrs)

	if instr, exists := instrs["MOV"]; exists {
		assert.NotEmpty(t, instr.Summary)
		assert.NotEmpty(t, instr.Forms)
		assert.NotEmpty(t, instr.Forms[0].Encodings)
	}
}

func TestGetInstructionByOpcode(t *testing.T) {
	instr, err := GetInstructionByOpcode("NOP")
	assert.NoError(t, err)
	assert.NotNil(t, instr)
	assert.NotEmpty(t, instr.Summary)
	assert.NotEmpty(t, instr.Forms)

	instr, err = GetInstructionByOpcode("MOV")
	assert.NoError(t, err)
	assert.NotNil(t, instr)
	assert.NotEmpty(t, instr.Summary)
	assert.NotEmpty(t, instr.Forms)

	instr, err = GetInstructionByOpcode("NONEXISTENT")
	assert.NoError(t, err)
	assert.Nil(t, instr)
}

func TestSegmentRegisterLookup(t *testing.T) {
	// MOV 命令の情報を取得
	instruction, err := GetInstructionByOpcode("MOV")
	if err != nil {
		t.Fatalf("Failed to get MOV instruction: %v", err)
	}
	assert.NotNil(t, instruction, "MOV instruction not found")

	// セグメントレジスタを含むオペランドの組み合わせを探す
	found := false
	for _, form := range instruction.Forms {
		if form.Operands == nil {
			continue
		}
		operands := *form.Operands
		if len(operands) == 2 &&
			((operands[0].Type == "r16" && operands[1].Type == "sreg") ||
				(operands[0].Type == "sreg" && operands[1].Type == "r16")) {
			found = true
			break
		}
	}

	// 見つからなかった場合はテストを失敗とする
	assert.True(t, found, "MOV instruction with segment register operands not found")
}
