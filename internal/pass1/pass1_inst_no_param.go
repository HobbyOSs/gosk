package pass1

import (
	"fmt" // Import fmt for Sprintf

	"github.com/HobbyOSs/gosk/internal/ast" // Change import to ast
)

// processNoParam now matches the opcodeEvalFn signature and accepts the instruction name.
// The 'operands' slice will be empty for no-parameter instructions.
func processNoParam(env *Pass1, operands []ast.Exp, instName string) { // Add instName parameter
	// パラメータを取らない命令（HLT等）は通常1バイト。
	// TODO: 命令によっては1バイトでない場合もあるため、将来的には命令名をenv.AsmDBで調べるべき。
	env.LOC += 1
	// Emit the instruction name as ocode (改行なし).
	env.Client.Emit(fmt.Sprintf("%s", instName)) // Remove newline
}
