package pass1

import (
	"fmt"
	"log"
	"strings"

	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/HobbyOSs/gosk/pkg/operand"
	"github.com/samber/lo"
)

func processMOV(env *Pass1, tokens []*token.ParseToken) {
	args := lo.Map(tokens, func(token *token.ParseToken, _ int) string {
		return token.AsString()
	})

	operands := operand.NewOperandFromString(strings.Join(args, ","))
	operandTypes := lo.Map(operands.OperandTypes(), func(t operand.OperandType, _ int) string {
		return t.String()
	})
	size, _ := env.AsmDB.FindMinOutputSize("MOV", operandTypes)

	log.Printf("debug: %s\noperands:%s\nsize: %d\n", operands.Serialize(), operandTypes, size)

	env.LOC += int32(size)
	env.Client.Emit(fmt.Sprintf("MOV %s\n", strings.Join(args, ",")))
}
