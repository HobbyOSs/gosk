package pass1

import (
	"fmt"
	"strings"

	"github.com/HobbyOSs/gosk/internal/token"
	// "github.com/HobbyOSs/gosk/pkg/cpu" // Removed unused import
	"github.com/HobbyOSs/gosk/pkg/operand" // Added import
	"github.com/samber/lo"
)

func processMOV(env *Pass1, tokens []*token.ParseToken) {
	args := lo.Map(tokens, func(token *token.ParseToken, _ int) string {
		return token.AsString()
	})

	operands := operand.
		NewOperandFromString(strings.Join(args, ",")).
		WithBitMode(env.BitMode).
		WithForceRelAsImm(true)
	size, _ := env.AsmDB.FindMinOutputSize("MOV", operands)
	env.LOC += int32(size)

	deb := fmt.Sprintf("MOV %s\n", strings.Join(args, ","))
	env.Client.Emit(deb)
}
