package pass1

import (
	"fmt"
	"log"

	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/HobbyOSs/gosk/pkg/ng_operand"
)

// processPUSH handles the PUSH instruction.
func processPUSH(env *Pass1, operands []ast.Exp) {
	processPushPopCommon(env, operands, "PUSH")
}

// processPOP handles the POP instruction.
func processPOP(env *Pass1, operands []ast.Exp) {
	processPushPopCommon(env, operands, "POP")
}

// processPushPopCommon is a common handler for PUSH and POP instructions.
func processPushPopCommon(env *Pass1, operands []ast.Exp, instName string) {
	if len(operands) != 1 {
		log.Printf("Error: %s instruction requires 1 operand, got %d", instName, len(operands))
		return
	}

	// Get string representation of the operand
	operandString := operands[0].TokenLiteral()

	// Create ng_operand.Operands from the string
	ngOperands, err := ng_operand.FromString(operandString)
	if err != nil {
		log.Printf("Error creating operands from string '%s' in %s: %v", operandString, instName, err)
		return
	}

	// Set BitMode
	ngOperands = ngOperands.WithBitMode(env.BitMode)

	// Calculate instruction size
	size, err := env.AsmDB.FindMinOutputSize(instName, ngOperands)
	if err != nil {
		log.Printf("Error finding min output size for %s %s: %v", instName, operandString, err)
		// Assume default size or handle error appropriately
		size = 1 // Default size assumption, might need refinement
	}
	env.LOC += int32(size)

	// Emit the command
	env.Client.Emit(fmt.Sprintf("%s %s", instName, operandString))
}
