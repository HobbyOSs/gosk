package asmdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindInstruction(t *testing.T) {
	data, err := decompressGzip(compressedJSON)
	assert.NoError(t, err)

	db := NewInstructionDB(data)

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
	data, err := decompressGzip(compressedJSON)
	assert.NoError(t, err)

	db := NewInstructionDB(data)

	forms, err := db.FindForms("MOV", []string{"r8", "m8"}) // MOV AL, [SI]
	assert.NoError(t, err)
	assert.NotEmpty(t, forms)
	assert.Equal(t, 2, forms[0].Encodings[0].GetOutputSize())

	forms, err = db.FindForms("MOV", []string{"r8", "imm8"}) // MOV AX, 0
	assert.NoError(t, err)
	assert.NotEmpty(t, forms)
	assert.Equal(t, 3, forms[0].Encodings[0].GetOutputSize())

	forms, err = db.FindForms("MOV", []string{"NONEXISTENT"})
	assert.NoError(t, err)
	assert.Empty(t, forms)

	forms, err = db.FindForms("NONEXISTENT", []string{"r32", "r32"})
	assert.Error(t, err)
	assert.Empty(t, forms)
}
