package asmdb

import (
	"testing"

	"github.com/HobbyOSs/gosk/pkg/operand"
	"github.com/stretchr/testify/assert"
)

func TestFindInstruction(t *testing.T) {
	instr, err := GetInstructionByOpcode("MOV")
	assert.NoError(t, err)
	assert.NotNil(t, instr)
	assert.NotEmpty(t, instr.Summary)
	assert.NotEmpty(t, instr.Forms)

	instr, err = GetInstructionByOpcode("NONEXISTENT")
	assert.Nil(t, err)
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
}

func TestFindEncoding_ModRM(t *testing.T) {
	db := NewInstructionDB()

	// ModRM が必要なケース
	encoding, err := db.FindEncoding("MOV", operand.NewOperandFromString("AL, [SI]")) // MOV AL, [SI]
	assert.NoError(t, err)
	assert.NotNil(t, encoding)
	assert.NotNil(t, encoding.ModRM, "ModRM should be required for MOV AL, [SI]")

	encoding, err = db.FindEncoding("MOV", operand.NewOperandFromString("[SI], AL")) // MOV [SI], AL
	assert.NoError(t, err)
	assert.NotNil(t, encoding)
	assert.NotNil(t, encoding.ModRM, "ModRM should be required for MOV [SI], AL")

	// ModRM が不要なケース
	encoding, err = db.FindEncoding("MOV", operand.NewOperandFromString("AL, [0x1234]")) // MOV AL, [0x1234]
	assert.NoError(t, err)
	assert.NotNil(t, encoding)
	assert.Nil(t, encoding.ModRM, "ModRM should not be required for MOV AL, [0x1234]")

	encoding, err = db.FindEncoding("MOV", operand.NewOperandFromString("AX, 0x1234")) // MOV AX, 0x1234
	assert.NoError(t, err)
	assert.NotNil(t, encoding)
	assert.Nil(t, encoding.ModRM, "ModRM should not be required for MOV AX, 0x1234")
}

func TestFindMinOutputSize(t *testing.T) {
	db := NewInstructionDB()

	t.Run("MOV AX, 0", func(t *testing.T) {
		size, err := db.FindMinOutputSize(
			"MOV", operand.NewOperandFromString("AX, 0").WithForceRelAsImm(true),
		)
		assert.NoError(t, err)
		assert.Equal(t, 3, size) // prefix(0x66) + opcode(0xB8) + imm16
	})

	t.Run("MOV r16, imm16 should use B8+rw form", func(t *testing.T) {
		operands := operand.
			NewOperandFromString("AX, 0x1234").
			WithForceImm8(true)

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
