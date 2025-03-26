package pass1

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/HobbyOSs/gosk/pkg/operand"
	"github.com/samber/lo"
)

// processLogicalInst は論理命令の共通処理を行うヘルパー関数です
func processLogicalInst(env *Pass1, tokens []*token.ParseToken, instName string) {
	args := lo.Map(tokens, func(token *token.ParseToken, _ int) string {
		return token.AsString()
	})

	isAccumulator := false
	if len(args) > 0 {
		matched, _ := regexp.MatchString(`(?i)^(AL|AX|EAX|RAX)$`, args[0])
		isAccumulator = matched
	}

	operands := operand.
		NewOperandFromString(strings.Join(args, ",")).
		WithBitMode(env.BitMode)

	// アキュムレータでない場合は `force_imm8` を有効にする
	if !isAccumulator {
		operands = operands.WithForceImm8(true)
	}

	size, _ := env.AsmDB.FindMinOutputSize(instName, operands)
	env.LOC += int32(size)
	env.Client.Emit(fmt.Sprintf("%s %s\n", instName, strings.Join(args, ",")))
}

// AND命令
func processAND(env *Pass1, tokens []*token.ParseToken) {
	processLogicalInst(env, tokens, "AND")
}
