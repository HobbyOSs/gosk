package pass1

import (
	"fmt"
	"log" // Add log import

	"github.com/HobbyOSs/gosk/internal/ast" // Add ast import
	"github.com/HobbyOSs/gosk/pkg/ng_operand"
)

// processMOV handles the MOV instruction using string representation of operands.
func processMOV(env *Pass1, operands []ast.Exp) {
	if len(operands) != 2 {
		log.Printf("Error: MOV instruction requires exactly two operands.")
		return
	}

	// Get string representation of operands
	// Assuming TokenLiteral() provides the correct representation for ng_operand.FromString
	op1Str := operands[0].TokenLiteral()
	op2Str := operands[1].TokenLiteral()
	operandString := op1Str + "," + op2Str

	// Create ng_operand.Operands from the combined string
	ngOperands, err := ng_operand.FromString(operandString)
	if err != nil {
		log.Printf("Error creating operands from string '%s' in MOV: %v", operandString, err)
		return // エラーが発生したら処理を中断
	}

	// Set BitMode
	ngOperands = ngOperands.WithBitMode(env.BitMode)
	// ngOperands = ngOperands.WithForceRelAsImm(true) // Re-evaluate if needed for MOV

	// Calculate instruction size
	size, err := env.AsmDB.FindMinOutputSize("MOV", ngOperands)
	if err != nil {
		log.Printf("Error finding min output size for MOV %s: %v", operandString, err)
		// Fallback or default size? For now, just log and don't update LOC.
		return
	}
	env.LOC += int32(size)

	// Emit the command
	// Use the original strings or the serialized version from ngOperands
	env.Client.Emit(fmt.Sprintf("MOV %s ; (size: %d)", ngOperands.Serialize(), size))
}
