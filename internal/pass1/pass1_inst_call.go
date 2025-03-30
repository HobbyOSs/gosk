package pass1

import (
	"github.com/HobbyOSs/gosk/internal/ast" // Add ast import
)

// processCALL handles the CALL instruction by delegating to processCalcJcc.
func processCALL(env *Pass1, operands []ast.Exp) {
	processCalcJcc(env, operands, "CALL") // Delegate to the common handler
}
