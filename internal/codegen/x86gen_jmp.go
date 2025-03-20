package codegen

import (
	"fmt"
	"strconv"

	"github.com/HobbyOSs/gosk/pkg/ocode"
)

// handleJMP handles the JMP instruction in code generation.
func handleJMP(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	return generateJMPCode(ocode.OpJMP, params.OCode, ctx, params.MachineCodeLen)
}

// handleJE handles the JE instruction in code generation.
func handleJE(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	return generateJMPCode(ocode.OpJE, params.OCode, ctx, params.MachineCodeLen)
}

// handleJNC handles the JNC instruction in code generation.
func handleJNC(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	return generateJMPCode(ocode.OpJNC, params.OCode, ctx, params.MachineCodeLen)
}

// handleJAE handles the JAE instruction in code generation.
func handleJAE(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	return generateJMPCode(ocode.OpJAE, params.OCode, ctx, params.MachineCodeLen)
}

// generateJMPCode generates the machine code for JMP and JE instructions.
func generateJMPCode(opKind ocode.OcodeKind, oc ocode.Ocode, ctx *CodeGenContext, currentMachineCodeLen int) ([]byte, error) {
	params := x86genParams{
		OCode:          oc,
		MachineCodeLen: currentMachineCodeLen,
	}
	// ジャンプ先アドレスを取得
	if len(params.OCode.Operands) < 1 {
		return nil, fmt.Errorf("jump instruction requires destination address")
	}
	destAddr, err := strconv.ParseInt(params.OCode.Operands[0], 0, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid jump destination address: %v", err)
	}

	// 現在のアドレス (ジャンプ命令の次のアドレス) を計算
	// ORG命令で設定されたDollarPositionを考慮する
	currentAddr := int64(ctx.DollarPosition) + int64(params.MachineCodeLen)

	var machineCode []byte

	switch opKind {
	case ocode.OpJMP:
		// JMP rel8 (オペコード: eb, オフセット: 1 byte)
		// JMP rel16 (オペコード: e9, オフセット: 2 bytes)
		// JMP rel32 (オペコード: e9, オフセット: 4 bytes)
		switch getOffsetSize(destAddr - currentAddr) {
		case 1:
			offset := destAddr - currentAddr - 2
			machineCode = []byte{0xeb, byte(offset)}
		case 2:
			offset := destAddr - currentAddr - 3
			machineCode = []byte{0xe9, byte(offset), byte(offset >> 8)}
		default:
			// NOP?
		}

	default:
		return handleJcc(params, ctx)
	}

	return machineCode, nil
}

func handleJcc(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	// ジャンプ先アドレスを取得
	if len(params.OCode.Operands) < 1 {
		return nil, fmt.Errorf("jump instruction requires destination address")
	}
	destAddr, err := strconv.ParseInt(params.OCode.Operands[0], 0, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid jump destination address: %v", err)
	}

	var machineCode []byte
	var opcode byte

	switch params.OCode.Kind {
	case ocode.OpJA:
		opcode = 0x77
	case ocode.OpJAE:
		opcode = 0x73
	case ocode.OpJB:
		opcode = 0x72
	case ocode.OpJBE:
		opcode = 0x76
	case ocode.OpJC:
		opcode = 0x72
	case ocode.OpJE:
		opcode = 0x74
	case ocode.OpJG:
		opcode = 0x7F
	case ocode.OpJGE:
		opcode = 0x7D
	case ocode.OpJL:
		opcode = 0x7C
	case ocode.OpJLE:
		opcode = 0x7E
	case ocode.OpJNA:
		opcode = 0x76
	case ocode.OpJNAE:
		opcode = 0x72
	case ocode.OpJNB:
		opcode = 0x73
	case ocode.OpJNBE:
		opcode = 0x77
	case ocode.OpJNC:
		opcode = 0x73
	case ocode.OpJNE:
		opcode = 0x75
	case ocode.OpJNG:
		opcode = 0x7E
	case ocode.OpJNGE:
		opcode = 0x7C
	case ocode.OpJNL:
		opcode = 0x7D
	case ocode.OpJNLE:
		opcode = 0x7F
	case ocode.OpJNO:
		opcode = 0x71
	case ocode.OpJNP:
		opcode = 0x7B
	case ocode.OpJNS:
		opcode = 0x79
	case ocode.OpJNZ:
		opcode = 0x75
	case ocode.OpJO:
		opcode = 0x70
	case ocode.OpJP:
		opcode = 0x7A
	case ocode.OpJPE:
		opcode = 0x7A
	case ocode.OpJPO:
		opcode = 0x7B
	case ocode.OpJS:
		opcode = 0x78
	case ocode.OpJZ:
		opcode = 0x74
	default:
		return nil, fmt.Errorf("invalid opcode kind for generateJMPCode: %v", params.OCode.Kind)
	}

	// 現在のアドレス (ジャンプ命令の次のアドレス) を計算
	// ORG命令で設定されたDollarPositionを考慮する
	currentAddr := int64(ctx.DollarPosition) + int64(params.MachineCodeLen)

	switch getOffsetSize(destAddr - currentAddr) {
	case 1:
		offset := destAddr - currentAddr - 2
		machineCode = []byte{opcode, byte(offset)}
	case 2:
		offset := destAddr - currentAddr - 3
		machineCode = []byte{opcode, byte(offset), byte(offset >> 8)}
	default:
		// NOP?
	}

	return machineCode, nil
}

// -128～127, -32768～32767 などの判定に使う
func getOffsetSize(imm int64) int {
	if imm >= -0x80 && imm <= 0x7f {
		return 1
	}
	if imm >= -0x8000 && imm <= 0x7fff {
		return 2
	}
	return 4
}
