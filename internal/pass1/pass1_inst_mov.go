package pass1

import (
	"fmt"
	"strings"

	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/HobbyOSs/gosk/pkg/operand"
	"github.com/samber/lo"
)

func processMOV(env *Pass1, tokens []*token.ParseToken) {
	args := lo.Map(tokens, func(token *token.ParseToken, _ int) string {
		return token.AsString()
	})

	operands := operand.NewOperandFromString(strings.Join(args, ",")).
		WithBitMode(env.BitMode)
	size, _ := env.AsmDB.FindMinOutputSize("MOV", operands)
	env.LOC += int32(size)
	env.Client.Emit(fmt.Sprintf("MOV %s\n", strings.Join(args, ",")))
}
