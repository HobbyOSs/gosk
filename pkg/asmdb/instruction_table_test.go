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
