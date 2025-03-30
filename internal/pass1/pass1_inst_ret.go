package pass1

import (
	"log" // Add log import

	"github.com/HobbyOSs/gosk/internal/ast" // Add ast import
)

// processRET handles the RET instruction.
// RET instruction (no operands) generates 1 byte of machine code (0xC3).
func processRET(env *Pass1, operands []ast.Exp) {
	if len(operands) != 0 {
		log.Printf("Warning: RET instruction should not have operands, but got %d.", len(operands))
	}
	// RET instruction size is 1 byte.
	env.LOC += 1
	// Emit the RET command.
	env.Client.Emit("RET")
}
