package asmdb

import (
	"testing"

	"github.com/HobbyOSs/gosk/pkg/operand"
	"github.com/stretchr/testify/assert"
)

func TestFindInstruction(t *testing.T) {
	db := NewInstructionDB()

	instr, found := db.FindInstruction("MOV")
	assert.True(t, found)
	assert.NotNil(t, instr)
	assert.NotEmpty(t, instr.Summary)
	assert.NotEmpty(t, instr.Forms)

	instr, found = db.FindInstruction("NONEXISTENT")
	assert.False(t, found)
	assert.Nil(t, instr)
}

func TestFindForms(t *testing.T) {
	db := NewInstructionDB()

	forms, err := db.FindForms("MOV", operand.NewOperandFromString("AL, [SI]")) // MOV AL, [SI]
	assert.NoError(t, err)
	assert.NotEmpty(t, forms)
	assert.Equal(t, 2, forms[0].Encodings[0].GetOutputSize(&OutputSizeOptions{}))

	forms, err = db.FindForms("MOV", operand.NewOperandFromString("NONEXISTENT"))
	assert.NoError(t, err)
	assert.Empty(t, forms)

	forms, err = db.FindForms("NONEXISTENT", operand.NewOperandFromString("EAX, 0"))
	assert.Error(t, err)
	assert.Empty(t, forms)
}

func FindMinOutputSize(t *testing.T) {
	db := NewInstructionDB()

	size, err := db.FindMinOutputSize("MOV", operand.NewOperandFromString("AX, 0")) // MOV AX, 0
	assert.NoError(t, err)
	assert.Equal(t, 3, size)
}
