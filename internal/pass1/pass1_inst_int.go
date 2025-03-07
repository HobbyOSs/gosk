package pass1

import (
	"fmt"
	"strings"

	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/HobbyOSs/gosk/pkg/operand"
	"github.com/samber/lo"
)

func processINT(env *Pass1, tokens []*token.ParseToken) {
	// Ocodeの生成
	if len(tokens) != 1 {
		panic("INT instruction requires one operand")
	}

	args := lo.Map(tokens, func(token *token.ParseToken, _ int) string {
		return token.AsString()
	})

	// オペランドを作成してサイズを計算
	operands := operand.
		NewOperandFromString(strings.Join(args, ",")).
		WithBitMode(env.BitMode)
	size, _ := env.AsmDB.FindMinOutputSize("INT", operands)
	env.LOC += int32(size)

	env.Client.Emit(fmt.Sprintf("INT %s", strings.Join(args, ",")))
}
