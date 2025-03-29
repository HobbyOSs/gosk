package pass1

import (
	"fmt" // Keep only one fmt import
	"strings"

	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/HobbyOSs/gosk/pkg/ng_operand" // Use ng_operand
	"github.com/samber/lo"
)

// processIN processes the IN instruction in pass1.
// It calculates the instruction size using asmdb and updates the LOC.
func processIN(env *Pass1, tokens []*token.ParseToken) {
	args := lo.Map(tokens, func(token *token.ParseToken, _ int) string {
		return token.AsString()
	})

	// Use ng_operand.FromString factory function
	operands, err := ng_operand.FromString(strings.Join(args, ","))
	if err != nil {
		// TODO: より適切なエラーハンドリングを行う
		fmt.Printf("Error creating operands from string in IN: %v\n", err)
		return // エラーが発生したら処理を中断
	}

	// Set BitMode
	operands = operands.WithBitMode(env.BitMode) // Assuming BitMode is relevant

	// Restore LOC calculation
	// FindMinOutputSize will use the fallback table if the instruction is not in the JSON DB
	size, err := env.AsmDB.FindMinOutputSize("IN", operands)
	if err != nil {
		// Handle error appropriately - maybe log or add to an error list in env
		// For now, just print an error message similar to how OUT might implicitly handle errors
		fmt.Printf("Error finding size for IN %s: %v\n", strings.Join(args, ","), err)
		// Decide on a default size or stop processing? For now, assume 0 size on error.
		size = 0 // Continue with size 0 on error? Let's keep original behavior for now.
	}
	env.LOC += int32(size)

	// Emit debug string
	deb := fmt.Sprintf("IN %s\n", strings.Join(args, ","))
	env.Client.Emit(deb)
}
