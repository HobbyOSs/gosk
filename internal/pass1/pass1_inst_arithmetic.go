package pass1

import (
	"fmt"
	"strings"

	"regexp"

	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/HobbyOSs/gosk/pkg/operand"
	"github.com/samber/lo"
)

// processArithmeticInst は算術命令の共通処理を行うヘルパー関数です
func processArithmeticInst(env *Pass1, tokens []*token.ParseToken, instName string) {
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

// 加算命令
func processADD(env *Pass1, tokens []*token.ParseToken) {
	processArithmeticInst(env, tokens, "ADD")
}

// キャリー付き加算命令
func processADC(env *Pass1, tokens []*token.ParseToken) {
	processArithmeticInst(env, tokens, "ADC")
}

// 減算命令
func processSUB(env *Pass1, tokens []*token.ParseToken) {
	processArithmeticInst(env, tokens, "SUB")
}

// ボロー付き減算命令
func processSBB(env *Pass1, tokens []*token.ParseToken) {
	processArithmeticInst(env, tokens, "SBB")
}

// 比較命令
func processCMP(env *Pass1, tokens []*token.ParseToken) {
	processArithmeticInst(env, tokens, "CMP")
}

// インクリメント命令
func processINC(env *Pass1, tokens []*token.ParseToken) {
	processArithmeticInst(env, tokens, "INC")
}

// デクリメント命令
func processDEC(env *Pass1, tokens []*token.ParseToken) {
	processArithmeticInst(env, tokens, "DEC")
}

// 2の補数命令
func processNEG(env *Pass1, tokens []*token.ParseToken) {
	processArithmeticInst(env, tokens, "NEG")
}

// 符号なし乗算命令
func processMUL(env *Pass1, tokens []*token.ParseToken) {
	processArithmeticInst(env, tokens, "MUL")
}

// 符号付き乗算命令
func processIMUL(env *Pass1, tokens []*token.ParseToken) {
	processArithmeticInst(env, tokens, "IMUL")
}

// 符号なし除算命令
func processDIV(env *Pass1, tokens []*token.ParseToken) {
	processArithmeticInst(env, tokens, "DIV")
}

// 符号付き除算命令
func processIDIV(env *Pass1, tokens []*token.ParseToken) {
	processArithmeticInst(env, tokens, "IDIV")
}
