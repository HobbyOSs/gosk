package pass1

import (
	"github.com/HobbyOSs/gosk/internal/ast" // Change import to ast
)

// processNoParam now matches the opcodeEvalFn signature.
// The 'operands' slice will be empty for no-parameter instructions.
func processNoParam(env *Pass1, operands []ast.Exp) {
	// パラメータを取らない命令（HLT等）は通常1バイト。
	// Emitはcodegenで行われるため、ここではLOCのインクリメントのみを行う。
	// TODO: 命令によっては1バイトでない場合もあるため、将来的には命令名をenv.AsmDBで調べるべき。
	env.LOC += 1
}
