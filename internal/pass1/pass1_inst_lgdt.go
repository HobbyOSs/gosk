package pass1

import (
	"fmt"
	"strings"

	"github.com/HobbyOSs/gosk/internal/token"
	// "github.com/HobbyOSs/gosk/pkg/cpu" // Removed unused import
	"github.com/HobbyOSs/gosk/pkg/operand" // Added import
	"github.com/samber/lo"
)

func processLGDT(env *Pass1, tokens []*token.ParseToken) {

	// オペランドの解析
	if len(tokens) != 1 {
		fmt.Printf("LGDT instruction expects exactly one operand")
		return
	}

	args := lo.Map(tokens, func(token *token.ParseToken, _ int) string {
		return token.AsString()
	})

	operands := operand.
		NewOperandFromString(strings.Join(args, ",")).
		WithBitMode(env.BitMode)
	if operands.ParsedOperands()[0].DirectMem == nil && operands.ParsedOperands()[0].IndirectMem == nil {
		fmt.Printf("LGDT instruction expects a memory operand")
		return
	}

	size, _ := env.AsmDB.FindMinOutputSize("LGDT", operands)
	env.LOC += int32(size)

	env.Client.Emit(fmt.Sprintf("LGDT %s\n", strings.Join(args, ",")))
}
