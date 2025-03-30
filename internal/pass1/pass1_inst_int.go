package pass1

import (
	"fmt"
	"log" // Add log import

	"github.com/HobbyOSs/gosk/internal/ast" // Add ast import
)

// processINT handles the INT instruction.
func processINT(env *Pass1, operands []ast.Exp) {
	if len(operands) != 1 {
		log.Printf("Error: INT instruction requires exactly one numeric operand.")
		return
	}

	// Evaluate the operand expression
	exp := operands[0]
	evaluatedExp, _ := exp.Eval(env) // Ignore the 'evaluated' flag for now

	// Check if the evaluated result is a NumberExp
	numExp, ok := evaluatedExp.(*ast.NumberExp)
	if !ok {
		// Log the type we actually got after evaluation
		log.Printf("Error: INT directive requires a numeric operand after evaluation, got %T.", evaluatedExp)
		return
	}

	// Check if the interrupt number is within the valid range (0-255)
	interruptNum := numExp.Value // Value is int64
	if interruptNum < 0 || interruptNum > 255 {
		log.Printf("Error: INT operand %d out of range (0-255).", interruptNum)
		return
	}

	// Calculate size: INT 3 is 1 byte (0xCC), others are 2 bytes (0xCD imm8)
	var size int32 = 2
	if interruptNum == 3 {
		size = 1
	}
	env.LOC += size

	// Emit the INT command with the interrupt number.
	env.Client.Emit(fmt.Sprintf("INT %d", interruptNum)) // Use the numeric value
}
