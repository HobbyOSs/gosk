package pass1

import (
	"fmt" // Keep only one fmt import
	"strings"

	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/HobbyOSs/gosk/pkg/ng_operand" // Use ng_operand
	"github.com/samber/lo"
)

// processArithmeticInst は算術命令の共通処理を行うヘルパー関数です
func processArithmeticInst(env *Pass1, tokens []*token.ParseToken, instName string) {
	args := lo.Map(tokens, func(token *token.ParseToken, _ int) string {
		return token.AsString()
	})

	// isAccumulator 関連コード削除
	// isAccumulator := false
	// if len(args) > 0 {
	// 	matched, _ := regexp.MatchString(`(?i)^(AL|AX|EAX|RAX)$`, args[0])
	// 	isAccumulator = matched
	// }

	// Use ng_operand.FromString factory function
	operands, err := ng_operand.FromString(strings.Join(args, ","))
	if err != nil {
		// TODO: より適切なエラーハンドリングを行う
		fmt.Printf("Error creating operands from string in %s: %v\n", instName, err)
		return // エラーが発生したら処理を中断
	}

	// Set BitMode
	operands = operands.WithBitMode(env.BitMode)
	// WithForceImm8 呼び出し削除
	// if !isAccumulator {
	// 	operands = operands.WithForceImm8(true)
	// }

	// Restore LOC calculation
	size, err := env.AsmDB.FindMinOutputSize(instName, operands)
	if err != nil {
		// TODO: より適切なエラーハンドリングを行う
		fmt.Printf("Error finding min output size for %s %s: %v\n", instName, strings.Join(args, ","), err)
		// エラーが発生しても LOC を加算しないようにする（あるいはデフォルトサイズを加算するか検討）
		return // エラーが発生したら処理を中断
	}
	env.LOC += int32(size)

	// Emit する文字列は元のオペランド文字列を使用
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
