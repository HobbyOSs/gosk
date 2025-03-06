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

func TestFindEncoding(t *testing.T) {
	db := NewInstructionDB()

	encoding, err := db.FindEncoding("MOV", operand.NewOperandFromString("AL, [SI]")) // MOV AL, [SI]
	assert.NoError(t, err)
	assert.NotNil(t, encoding)
	assert.Equal(t, 2, encoding.GetOutputSize(&OutputSizeOptions{}))

	encoding, err = db.FindEncoding("MOV", operand.NewOperandFromString("NONEXISTENT"))
	assert.Error(t, err)
	assert.Nil(t, encoding)
	assert.Contains(t, err.Error(), "no matching encoding found")

	encoding, err = db.FindEncoding("NONEXISTENT", operand.NewOperandFromString("EAX, 0"))
	assert.Error(t, err)
	assert.Nil(t, encoding)
	assert.Contains(t, err.Error(), "instruction not found")
}

func TestFindMinOutputSize(t *testing.T) {
	db := NewInstructionDB()

	t.Run("MOV AX, 0", func(t *testing.T) {
		size, err := db.FindMinOutputSize("MOV", operand.NewOperandFromString("AX, 0"))
		assert.NoError(t, err)
		assert.Equal(t, 3, size) // prefix(0x66) + opcode(0xB8) + imm16
	})

	t.Run("MOV r16, imm16 should use B8+rw form", func(t *testing.T) {
		operands := operand.NewOperandFromString("AX, 0x1234")
		t.Logf("Operand types: %v", operands.OperandTypes())

		encoding, err := db.FindEncoding("MOV", operands)
		if err != nil {
			t.Logf("Error finding encoding: %v", err)
		}
		if encoding != nil {
			t.Logf("Found encoding: %v", encoding)
		}

		assert.NoError(t, err)
		assert.NotNil(t, encoding)
		// 0xB8+rw encodingが選択されるべき（より短いため）
		assert.Equal(t, "B8", encoding.Opcode.Byte, "B8+rw encoding should be selected")
	})

	// TODO: Fix memory operand with WORD prefix test
	t.Run("MOV with prefix and offset", func(t *testing.T) {
		t.Skip("Memory operand with WORD prefix needs to be fixed")
	})
}
