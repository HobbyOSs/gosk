package codegen

import (
	"fmt"

	"github.com/HobbyOSs/gosk/pkg/ocode"
)

// handleJMP handles the JMP instruction in code generation.
func handleJMP(oc ocode.Ocode, ctx *CodeGenContext) ([]byte, error) {
	return generateJMPCode(ocode.OpJMP, oc, ctx)
}

// handleJE handles the JE instruction in code generation.
func handleJE(oc ocode.Ocode, ctx *CodeGenContext) ([]byte, error) {
	return generateJMPCode(ocode.OpJE, oc, ctx)
}

// generateJMPCode generates the machine code for JMP and JE instructions.
func generateJMPCode(opKind ocode.OcodeKind, oc ocode.Ocode, ctx *CodeGenContext) ([]byte, error) {
	// TODO: Pass2で解決されたジャンプ先ラベルのアドレスをSymTableから取得
	// TODO: 相対ジャンプのオフセットを計算 (ジャンプ元アドレス - ジャンプ先アドレス)
	// TODO: オフセットサイズに応じて、JMP rel8 または JMP rel16 の機械語コードを生成
	//   - JMP rel8 (オペコード: eb (JMP), 74 (JE), オフセット: 1 byte)
	//   - JMP rel16 (オペコード: e9 (JMP),  , オフセット: 2 bytes)

	// 仮実装: JMP rel8 を生成
	switch opKind {
	case ocode.OpJMP:
		return []byte{0xeb, 0x00}, nil
	case ocode.OpJE:
		return []byte{0x74, 0x00}, nil
	default:
		return nil, fmt.Errorf("invalid opcode kind for generateJMPCode: %v", opKind)
	}
}
