package pass1

import (
	"fmt"
	"log" // Add log import
	"strings"

	"github.com/HobbyOSs/gosk/internal/ast" // Add ast import
	"github.com/HobbyOSs/gosk/pkg/ng_operand"
	"github.com/samber/lo"
)

// processMOV handles the MOV instruction using string representation of operands.
func processMOV(env *Pass1, operands []ast.Exp) {
	if len(operands) != 2 {
		log.Printf("Error: MOV instruction requires exactly two operands.")
		return
	}

	args := lo.Map(operands, func(op ast.Exp, _ int) string {
		return op.TokenLiteral()
	})

	// Create ng_operand.Operands from the combined string
	ngOperands, err := ng_operand.FromString(strings.Join(args, ","))
	if err != nil {
		log.Printf("Error creating operands from string '%s' in MOV: %v", strings.Join(args, ","), err)
		return // エラーが発生したら処理を中断
	}
	log.Printf("debug: processMOV: ngOperands after FromString: %s", ngOperands.Serialize())

	// Set BitMode and ForceRelAsImm
	ngOperands = ngOperands.WithBitMode(env.BitMode).
		WithForceRelAsImm(true) // Force relative symbols (like labels) to be treated as immediates for size calculation

	// Calculate instruction size
	size, err := env.AsmDB.FindMinOutputSize("MOV", ngOperands)
	if err != nil {
		// Log operands separately for clarity
		log.Printf("Error finding min output size for MOV (op1: '%s', op2: '%s'): %v", args[0], args[1], err)
		// Fallback or default size? For now, just log and don't update LOC.
		return
	}
	log.Printf("debug: processMOV: Calculated size=%d", size)
	env.LOC += int32(size)

	// Emit the command
	// Use the original strings or the serialized version from ngOperands
	env.Client.Emit(fmt.Sprintf("MOV %s", strings.Join(args, ",")))
}
