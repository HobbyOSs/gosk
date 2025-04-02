package pass1

import (
	"fmt"
	"log" // Add log import

	"github.com/HobbyOSs/gosk/internal/ast" // ast インポートを追加
)

// processINT は INT 命令を処理します。
func processINT(env *Pass1, operands []ast.Exp) {
	if len(operands) != 1 {
		log.Printf("Error: INT instruction requires exactly one numeric operand.")
		return
	}

	operandExp := operands[0]
	var size int32 = 2 // デフォルトは INT imm8 の 2 バイト

	// env (Pass1) が持つ GetConstValue を使って定数値を取得し、INT 3 のサイズを判定
	if val, ok := env.GetConstValue(operandExp); ok && val == 3 {
		size = 1 // INT 3 は 1 バイト
	}
	env.LOC += size

	// 修正された ast.ExpToString を使ってオペランド文字列を生成し、codegen に渡す
	operandStr := ast.ExpToString(operandExp)
	if operandStr == "" {
		// ExpToString が空文字列を返した場合のエラーハンドリング (念のため)
		log.Printf("Error: Failed to convert INT operand expression to string: %T", operandExp)
		// ここで処理を中断するか、デフォルトの動作を続けるか検討
		// 今回はログのみ出力し、Emit は試みる
	}
	env.Client.Emit(fmt.Sprintf("INT %s", operandStr))
}
