package pass1

import (
	"github.com/HobbyOSs/gosk/internal/ast" // ast インポートを追加
)

// processCALL は CALL 命令を processCalcJcc に委譲して処理します。
func processCALL(env *Pass1, operands []ast.Exp) {
	processCalcJcc(env, operands, "CALL") // 共通ハンドラに委譲
}
