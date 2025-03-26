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

	// CALL rel32 (オペコード: e8, オフセット: 4 bytes)
	offset := destAddr - currentAddr - 5
	machineCode = []byte{0xe8, byte(offset), byte(offset >> 8), byte(offset >> 16), byte(offset >> 24)}

	return machineCode, nil
}
