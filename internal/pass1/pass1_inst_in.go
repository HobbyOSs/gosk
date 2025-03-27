package pass1

import (
	"fmt"
	"strings"

	"github.com/HobbyOSs/gosk/internal/token"
	// "github.com/HobbyOSs/gosk/pkg/cpu" // Removed unused import
	"github.com/HobbyOSs/gosk/pkg/operand" // Added import
	"github.com/samber/lo"
)

// processIN processes the IN instruction in pass1.
// It calculates the instruction size using asmdb and updates the LOC.
func processIN(env *Pass1, tokens []*token.ParseToken) {
	args := lo.Map(tokens, func(token *token.ParseToken, _ int) string {
		return token.AsString()
	})

	operands := operand.
		NewOperandFromString(strings.Join(args, ",")).
		WithBitMode(env.BitMode) // Assuming BitMode is relevant, like in processOUT

	// FindMinOutputSize will use the fallback table if the instruction is not in the JSON DB
	size, err := env.AsmDB.FindMinOutputSize("IN", operands)
	if err != nil {
		// Handle error appropriately - maybe log or add to an error list in env
		// For now, just print an error message similar to how OUT might implicitly handle errors
		fmt.Printf("Error finding size for IN %s: %v\n", strings.Join(args, ","), err)
		// Decide on a default size or stop processing? For now, assume 0 size on error.
		size = 0
	}

	env.LOC += int32(size)

	// Emit debug string
	deb := fmt.Sprintf("IN %s\n", strings.Join(args, ","))
	env.Client.Emit(deb)
}
