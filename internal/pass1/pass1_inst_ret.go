package pass1

import (
	"github.com/HobbyOSs/gosk/internal/token"
)

// processRET はRET命令のPass1処理を行います。
// RET命令（オペランドなし）は1バイトの機械語を生成します。
func processRET(env *Pass1, tokens []*token.ParseToken) {
	// Emitは呼び出し元（handlers.goのTraverseAST内のOpcodeStmtケース）で実行されるため、
	// ここではLOCのインクリメントのみを行う
	env.LOC += 1
}
