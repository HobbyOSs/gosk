package codegen

import (
	"github.com/HobbyOSs/gosk/pkg/ocode"
)

// handleJMP handles the JMP instruction in code generation.
func handleJMP(oc ocode.Ocode, ctx *CodeGenContext) ([]byte, error) {
	// TODO: Pass2で解決されたジャンプ先ラベルのアドレスをSymTableから取得
	// TODO: 相対ジャンプのオフセットを計算 (ジャンプ元アドレス - ジャンプ先アドレス)
	// TODO: オフセットサイズに応じて、JMP rel8 または JMP rel16 の機械語コードを生成
	//   - JMP rel8 (オペコード: eb, オフセット: 1 byte)
	//   - JMP rel16 (オペコード: e9, オフセット: 2 bytes)

	// 仮実装: JMP rel8 (eb) を生成
	return []byte{0xeb, 0x00}, nil
}
