package pass1

import (
	"fmt" // Keep only one fmt import
	"strings"

	"log" // Add log import

	"github.com/HobbyOSs/gosk/internal/ast"   // Add ast import
	"github.com/HobbyOSs/gosk/pkg/ng_operand" // Use ng_operand
	"github.com/samber/lo"
)

// processIN は IN 命令を処理します。
func processIN(env *Pass1, operands []ast.Exp) {
	instName := "IN"
	// オペランドの文字列表現を取得します
	operandStrings := lo.Map(operands, func(exp ast.Exp, _ int) string {
		return exp.TokenLiteral()
	})
	operandString := strings.Join(operandStrings, ",")

	// 結合された文字列から ng_operand.Operands を作成します
	ngOperands, err := ng_operand.FromString(operandString)
	if err != nil {
		log.Printf("Error creating operands from string '%s' in %s: %v", operandString, instName, err)
		return
	}

	// BitMode を設定します
	ngOperands = ngOperands.WithBitMode(env.BitMode)

	// 命令サイズを計算します
	size, err := env.AsmDB.FindMinOutputSize(instName, ngOperands)
	if err != nil {
		log.Printf("Error finding min output size for %s %s: %v", instName, operandString, err)
		// デフォルトサイズを決定するか、処理を停止しますか？ 現時点では、エラー時にサイズ 0 を想定します。
		size = 0 // 現時点では元の動作を維持します。
	}
	env.LOC += int32(size)

	// コマンドを発行します (元のオペランド文字列をカンマ区切りで使用)
	env.Client.Emit(fmt.Sprintf("%s %s", instName, strings.Join(operandStrings, ",")))
}
