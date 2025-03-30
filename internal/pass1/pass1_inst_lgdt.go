package pass1

import (
	"fmt" // Keep only one fmt import
	"log" // Add log import

	"github.com/HobbyOSs/gosk/internal/ast" // Add ast import
	"github.com/HobbyOSs/gosk/pkg/cpu"
	"github.com/HobbyOSs/gosk/pkg/ng_operand" // Use ng_operand
)

// processLGDT handles the LGDT instruction.
func processLGDT(env *Pass1, operands []ast.Exp) {
	instName := "LGDT"
	if len(operands) != 1 {
		log.Printf("Error: %s instruction requires exactly one operand.", instName)
		return
	}

	// Get string representation of the operand
	operandString := operands[0].TokenLiteral()

	// Create ng_operand.Operands from the string
	ngOperands, err := ng_operand.FromString(operandString)
	if err != nil {
		log.Printf("Error creating operand from string '%s' in %s: %v", operandString, instName, err)
		return
	}

	// Set BitMode
	ngOperands = ngOperands.WithBitMode(env.BitMode)

	// Check if the operand is a memory type
	opTypes := ngOperands.OperandTypes()
	// Manually check if the type is one of the memory types
	isMem := len(opTypes) == 1 && (opTypes[0] == ng_operand.CodeM ||
		opTypes[0] == ng_operand.CodeM8 ||
		opTypes[0] == ng_operand.CodeM16 ||
		opTypes[0] == ng_operand.CodeM32 ||
		opTypes[0] == ng_operand.CodeM64 ||
		opTypes[0] == ng_operand.CodeMEM)
	if !isMem {
		log.Printf("Error: %s instruction expects a memory operand, got %v (raw: %s)", instName, opTypes, operandString)
		return
	}

	// Calculate LGDT size based on BitMode (using original logic for now)
	// TODO: Consider using FindMinOutputSize for consistency?
	var lgdtSize int32
	switch env.BitMode {
	case cpu.MODE_16BIT:
		// LGDT m16:32 (0F 01 /1) - Opcode(2) + ModRM(1) + disp/addr(2/4 depending on addressing)
		// Size calculation needs refinement based on addressing mode.
		// Assuming ModRM + disp16/32 for simplicity here.
		// FindMinOutputSize might be more accurate.
		lgdtSize = int32(3 + ngOperands.CalcOffsetByteSize()) // Opcode(2) + ModRM(1) + Offset
	case cpu.MODE_32BIT:
		lgdtSize = int32(3 + ngOperands.CalcOffsetByteSize()) // Opcode(2) + ModRM(1) + Offset
	default:
		log.Printf("Error: Unsupported bit mode %v for %s size calculation", env.BitMode, instName)
		return
	}
	env.LOC += lgdtSize

	// Emit the command
	env.Client.Emit(fmt.Sprintf("%s %s ; (size: %d)", instName, ngOperands.Serialize(), lgdtSize))
}
