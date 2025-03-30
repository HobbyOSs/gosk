package pass1

import (
	"fmt" // Keep only one fmt import
	"strings"

	"log" // Add log import

	"github.com/HobbyOSs/gosk/internal/ast"   // Add ast import
	"github.com/HobbyOSs/gosk/pkg/ng_operand" // Use ng_operand
	"github.com/samber/lo"
)

// processIN handles the IN instruction.
func processIN(env *Pass1, operands []ast.Exp) {
	instName := "IN"
	// Get string representation of operands
	operandStrings := lo.Map(operands, func(exp ast.Exp, _ int) string {
		return exp.TokenLiteral()
	})
	operandString := strings.Join(operandStrings, ",")

	// Create ng_operand.Operands from the combined string
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
		// Decide on a default size or stop processing? For now, assume 0 size on error.
		size = 0 // Keep original behavior for now.
	}
	env.LOC += int32(size)

	// Emit the command
	env.Client.Emit(fmt.Sprintf("%s %s", instName, ngOperands.Serialize()))
}
