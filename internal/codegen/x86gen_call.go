package codegen

import (
	"fmt"
	"strconv"

	"github.com/HobbyOSs/gosk/pkg/ocode"
)

func handleCALL(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	// ジャンプ先アドレスを取得
	if len(params.OCode.Operands) < 1 {
		return nil, fmt.Errorf("call instruction requires destination address")
	}
	destAddr, err := strconv.ParseInt(params.OCode.Operands[0], 0, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid call destination address: %v", err)
	}

	var machineCode []byte

	// 現在のアドレス (CALL命令の次のアドレス) を計算
	// ORG命令で設定されたDollarPositionを考慮する
	currentAddr := int64(ctx.DollarPosition) + int64(params.MachineCodeLen)

	// kind が OpCALL であることを確認
	if params.OCode.Kind != ocode.OpCALL {
		return nil, fmt.Errorf("invalid opcode kind for handleCALL: %v", params.OCode.Kind)
	}

	// オフセットを計算 (仮に rel32 として計算)
	offset32 := destAddr - currentAddr - 5

	// オフセットが rel16 の範囲内か確認
	if offset32 >= -32768 && offset32 <= 32767 {
		// CALL rel16 (オペコード: e8, オフセット: 2 bytes)
		offset16 := destAddr - currentAddr - 3 // rel16 の命令長は 3 バイト
		machineCode = []byte{0xe8, byte(offset16), byte(offset16 >> 8)}
	} else {
		// CALL rel32 (オペコード: e8, オフセット: 4 bytes)
		machineCode = []byte{0xe8, byte(offset32), byte(offset32 >> 8), byte(offset32 >> 16), byte(offset32 >> 24)}
	}

	return machineCode, nil
}
