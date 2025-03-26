package pass1

import (
	"fmt"
	"log"
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

// OR命令
func processOR(env *Pass1, tokens []*token.ParseToken) {
	processLogicalInst(env, tokens, "OR")
}

// XOR命令
func processXOR(env *Pass1, tokens []*token.ParseToken) {
	processLogicalInst(env, tokens, "XOR")
}

// NOT命令
func processNOT(env *Pass1, tokens []*token.ParseToken) {
	if len(tokens) != 1 {
		log.Fatalf("error: NOT instruction requires 1 operand, but got %d", len(tokens))
	}
	arg := tokens[0].AsString()

	operands := operand.
		NewOperandFromString(arg).
		WithBitMode(env.BitMode)

	size, _ := env.AsmDB.FindMinOutputSize("NOT", operands)
	env.LOC += int32(size)
	env.Client.Emit(fmt.Sprintf("NOT %s\n", arg))
}

// SHR命令
func processSHR(env *Pass1, tokens []*token.ParseToken) {
	processLogicalInst(env, tokens, "SHR")
}

// SHL命令
func processSHL(env *Pass1, tokens []*token.ParseToken) {
	processLogicalInst(env, tokens, "SHL")
}

// SAR命令
func processSAR(env *Pass1, tokens []*token.ParseToken) {
	processLogicalInst(env, tokens, "SAR")
}
