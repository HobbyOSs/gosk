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

func TestFindForm(t *testing.T) {
	db := NewInstructionDB()

	form, err := db.FindForm("MOV", operand.NewOperandFromString("AL, [SI]")) // MOV AL, [SI]
	assert.NoError(t, err)
	assert.NotNil(t, form)
	assert.Equal(t, 2, form.Encodings[0].GetOutputSize(&OutputSizeOptions{}))

	form, err = db.FindForm("MOV", operand.NewOperandFromString("NONEXISTENT"))
	assert.Error(t, err)
	assert.Nil(t, form)
	assert.Contains(t, err.Error(), "no matching form found")

	form, err = db.FindForm("NONEXISTENT", operand.NewOperandFromString("EAX, 0"))
	assert.Error(t, err)
	assert.Nil(t, form)
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

		form, err := db.FindForm("MOV", operands)
		if err != nil {
			t.Logf("Error finding form: %v", err)
		}
		if form != nil {
			t.Logf("Found form operands: %v", form.Operands)
			t.Logf("Found form encodings: %v", form.Encodings)
		}

		assert.NoError(t, err)
		assert.NotNil(t, form)
		// 0xB8+rw formが選択されるべき（より短いため）
		// B8+rwのエンコーディングを探す
		foundB8 := false
		for _, enc := range form.Encodings {
			if enc.Opcode.Byte == "B8" {
				foundB8 = true
				break
			}
		}
		assert.True(t, foundB8, "B8+rw encoding should be found")
	})

	// TODO: Fix memory operand with WORD prefix test
	t.Run("MOV with prefix and offset", func(t *testing.T) {
		t.Skip("Memory operand with WORD prefix needs to be fixed")
	})
}
