package codegen

import (
	"fmt"
	"strconv"

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
	// ジャンプ先アドレスを取得
	if len(oc.Operands) < 1 {
		return nil, fmt.Errorf("jump instruction requires destination address")
	}
	destAddr, err := strconv.ParseInt(oc.Operands[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid jump destination address: %v", err)
	}

	// 現在のアドレス (ジャンプ命令の次のアドレス) を計算
	// ORG命令で設定されたDollarPositionを考慮する
	currentAddr := int64(ctx.DollarPosition) + int64(len(ctx.MachineCode))
	fmt.Printf("Debug: DollarPosition=0x%x, len(ctx.MachineCode)=%d, currentAddr=0x%x\n",
		ctx.DollarPosition, len(ctx.MachineCode), currentAddr)
	var offset int64
	var machineCode []byte

	switch opKind {
	case ocode.OpJMP:
		// JMP rel8 (オペコード: eb, オフセット: 1 byte)
		// JMP rel16 (オペコード: e9, オフセット: 2 bytes)
		offset = destAddr - (currentAddr + 2) // 2バイト = オペコード(1) + オフセット(1)
		if offset >= -128 && offset <= 127 {
			// 8ビットオフセットで表現可能な場合
			machineCode = []byte{0xeb, byte(offset)}
			fmt.Printf("JMP rel8: destAddr=0x%x, currentAddr=0x%x, offset=0x%x\n", 
				destAddr, currentAddr, offset)
		} else {
			// 16ビットオフセットが必要な場合
			offset = destAddr - (currentAddr + 3) // 3バイト = オペコード(1) + オフセット(2)
			machineCode = []byte{0xe9, byte(offset), byte(offset >> 8)}
			fmt.Printf("JMP rel16: destAddr=%d, currentAddr=%d, offset=%d\n", destAddr, currentAddr, offset) // デバッグログ
		}
	case ocode.OpJE:
		// JE rel8 (オペコード: 74, オフセット: 1 byte)
		offset = destAddr - (currentAddr + 2)
		if offset >= -128 && offset <= 127 {
			machineCode = []byte{0x74, byte(offset)}
			fmt.Printf("JE rel8: destAddr=%d, currentAddr=%d, offset=%d\n", destAddr, currentAddr, offset) // デバッグログ
		} else {
			return nil, fmt.Errorf("JE instruction with offset larger than 8 bits is not supported")
		}
	default:
		return nil, fmt.Errorf("invalid opcode kind for generateJMPCode: %v", opKind)
	}

	return machineCode, nil
}
