package pass1

import (
	"github.com/HobbyOSs/gosk/internal/token"
)

func processNoParam(env *Pass1, tokens []*token.ParseToken) {
	// パラメータを取らない命令（HLT等）は1バイトの機械語を生成
	// Emitは呼び出し元（handlers.goのTraverseAST内のOpcodeStmtケース）で実行されるため、
	// ここではLOCのインクリメントのみを行う
	env.LOC += 1
}
