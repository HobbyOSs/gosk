package codegen

import "github.com/HobbyOSs/gosk/pkg/ocode"

// opcodeMapRET はRET命令のオペコードを定義します。
var opcodeMapRET = map[ocode.OcodeKind]byte{
	ocode.OpRET: 0xC3,
}

// handleRET はRET命令のOcodeを処理し、対応するx86機械語を生成します。
func handleRET(ocode ocode.Ocode) ([]byte, error) {
	var binary []byte
	if code, exists := opcodeMapRET[ocode.Kind]; exists {
		binary = append(binary, code)
	}
	// RET命令はエラーを返さない想定
	return binary, nil
}
