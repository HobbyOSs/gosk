package pass1

import (
	"fmt" // Keep only one fmt import
	"log" // Keep only one log import
	"strings"

	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/HobbyOSs/gosk/pkg/ng_operand" // Use ng_operand
	"github.com/samber/lo"
)

// processLogicalInst は論理命令の共通処理を行うヘルパー関数です
func processLogicalInst(env *Pass1, tokens []*token.ParseToken, instName string) {
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

	// Set BitMode (WithForceImm8 削除)
	operands = operands.WithBitMode(env.BitMode)
	// if !isAccumulator {
	// }

	// Restore LOC calculation
	size, err := env.AsmDB.FindMinOutputSize(instName, operands)
	if err != nil {
		// TODO: より適切なエラーハンドリングを行う
		fmt.Printf("Error finding min output size for %s %s: %v\n", instName, strings.Join(args, ","), err)
		return // エラーが発生したら処理を中断
	}
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

	// Use ng_operand.FromString factory function
	operands, err := ng_operand.FromString(arg)
	if err != nil {
		// TODO: より適切なエラーハンドリングを行う
		fmt.Printf("Error creating operands from string in NOT: %v\n", err)
		return // エラーが発生したら処理を中断
	}

	// Set BitMode (WithForceImm8 削除)
	operands = operands.WithBitMode(env.BitMode) // 算術命令に合わせて一旦 true に設定 -> このコメントも不要か？

	// Restore LOC calculation
	size, err := env.AsmDB.FindMinOutputSize("NOT", operands)
	if err != nil {
		// TODO: より適切なエラーハンドリングを行う
		fmt.Printf("Error finding min output size for NOT %s: %v\n", arg, err)
		return // エラーが発生したら処理を中断
	}
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
