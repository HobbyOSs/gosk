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

	// オペランド式を評価します
	exp := operands[0]
	evaluatedExp, _ := exp.Eval(env) // 現時点では 'evaluated' フラグは無視します

	// 評価結果が NumberExp かどうかを確認します
	numExp, ok := evaluatedExp.(*ast.NumberExp)
	if !ok {
		// 評価後に実際に取得した型をログに出力します
		log.Printf("Error: INT directive requires a numeric operand after evaluation, got %T.", evaluatedExp)
		return
	}

	// 割り込み番号が有効な範囲 (0-255) 内にあるか確認します
	interruptNum := numExp.Value // Value は int64 です
	if interruptNum < 0 || interruptNum > 255 {
		log.Printf("Error: INT operand %d out of range (0-255).", interruptNum)
		return
	}

	// サイズを計算します: INT 3 は 1 バイト (0xCC)、その他は 2 バイト (0xCD imm8)
	var size int32 = 2
	if interruptNum == 3 {
		size = 1
	}
	env.LOC += size

	// 割り込み番号を指定して INT コマンドを発行します。
	env.Client.Emit(fmt.Sprintf("INT %d", interruptNum)) // 数値を使用します
}
