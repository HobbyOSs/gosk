package codegen

import (
	"fmt"
	"strconv"
	"strings" // stringsパッケージをインポート

	"github.com/HobbyOSs/gosk/pkg/cpu" // cpuパッケージをインポート
	"github.com/HobbyOSs/gosk/pkg/ocode"
)

func handleJcc(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	// ジャンプ先アドレスを取得 (JMP_FAR以外の場合)
	var destAddr int64
	var err error
	if params.OCode.Kind != ocode.OpJMP_FAR && len(params.OCode.Operands) >= 1 {
		// JMP/Jcc命令の場合、オペランドは通常1つ (ラベルまたは即値)
		// ラベルはpass2で解決されるため、ここでは数値として扱えるはず
		destAddr, err = strconv.ParseInt(params.OCode.Operands[0], 0, 64)
		if err != nil {
			// JMP_FAR のオペランド形式 (segment:offset) が誤って渡された可能性も考慮
			if strings.Contains(params.OCode.Operands[0], ":") {
				return nil, fmt.Errorf("unexpected segment:offset format for non-JMP_FAR instruction: %v", params.OCode.Operands)
			}
			return nil, fmt.Errorf("invalid jump destination address for %s: %v", params.OCode.Kind, err)
		}
	} else if params.OCode.Kind != ocode.OpJMP_FAR && len(params.OCode.Operands) == 0 {
		// オペランドがないJcc命令はないはずだが念のため
		return nil, fmt.Errorf("%s instruction requires destination address", params.OCode.Kind)
	}

	var machineCode []byte
	var opcode byte

	// 現在のアドレス (ジャンプ命令の次のアドレス) を計算
	// ORG命令で設定されたDollarPositionを考慮する
	currentAddr := int64(ctx.DollarPosition) + int64(params.MachineCodeLen)

	switch params.OCode.Kind {
	// JMP_FAR の処理を追加
	case ocode.OpJMP_FAR:
		// オペランドは "セグメント:オフセット" 形式で1つだけ渡されるはず
		if len(params.OCode.Operands) != 1 {
			return nil, fmt.Errorf("JMP_FAR requires 1 operand (segment:offset), got %d (%v)", len(params.OCode.Operands), params.OCode.Operands)
		}
		operand := params.OCode.Operands[0]
		parts := strings.Split(operand, ":") // コロンで分割
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid operand format for JMP_FAR, expected 'segment:offset', got %s", operand)
		}
		segmentStr := parts[0]
		offsetStr := parts[1]

		segment, err := strconv.ParseInt(segmentStr, 10, 16) // セグメントは16ビット
		if err != nil {
			return nil, fmt.Errorf("invalid segment value for JMP_FAR: %v", err)
		}
		offset, err := strconv.ParseInt(offsetStr, 10, 32) // オフセットは32ビット
		if err != nil {
			return nil, fmt.Errorf("invalid offset value for JMP_FAR: %v", err)
		}

		machineCode = []byte{
			0xea, // JMP ptr16:32 オペコード
			byte(offset),
			byte(offset >> 8),
			byte(offset >> 16),
			byte(offset >> 24),
			byte(segment),
			byte(segment >> 8),
		}

		// オペランドサイズプレフィックス (66h) が必要かチェック
		// JMP ptr16:32 はデフォルトのオペランドサイズが32ビットの場合プレフィックス不要
		// 16ビットモードの場合は必要
		if ctx.BitMode == cpu.MODE_16BIT {
			machineCode = append([]byte{0x66}, machineCode...)
		}

		return machineCode, nil

	case ocode.OpJMP:
		// JMP rel8 (オペコード: eb, オフセット: 1 byte)
		// JMP rel16 (オペコード: e9, オフセット: 2 bytes)
		// JMP rel32 (オペコード: e9, オフセット: 4 bytes)
		relativeOffset := destAddr - currentAddr // ジャンプ先までの相対距離
		offsetSize := getOffsetSize(relativeOffset)

		switch offsetSize {
		case 1:
			// rel8: Opcode(1) + Offset(1) = 2 bytes
			machineCode = []byte{0xeb, byte(relativeOffset - 2)}
		case 2:
			// rel16: Opcode(1) + Offset(2) = 3 bytes
			// 16bitモードでの JMP rel16 は E9 cw (3 bytes)。66h プレフィックスは不要。
			// 32bitモードでの JMP rel16 も E9 cw (3 bytes)。
			// 32bitモードでは rel32 (E9 cd, 5バイト) が一般的だが、ここでは rel16 を生成する。
			machineCode = []byte{0xe9, byte(relativeOffset - 3), byte((relativeOffset - 3) >> 8)}
			// 66h プレフィックスは不要なため、ifブロックを削除
		default: // rel32
			// rel32: Opcode(1) + Offset(4) = 5 bytes
			machineCode = []byte{
				0xe9,
				byte(relativeOffset - 5),
				byte((relativeOffset - 5) >> 8),
				byte((relativeOffset - 5) >> 16),
				byte((relativeOffset - 5) >> 24),
			}
			// TODO: 32bitモードの場合、オペランドサイズプレフィックス(66h)が不要か確認
			// JMP rel32 は 32bitモードでは E9 cd (5バイト)
			// 16bitモードでは 66 E9 cd (6バイト)
			if ctx.BitMode == cpu.MODE_16BIT {
				machineCode = append([]byte{0x66}, machineCode...)
			}
		}
		return machineCode, nil // JMPの場合はここでreturn
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

	switch getOffsetSize(destAddr - currentAddr) {
	case 1:
		offset := destAddr - currentAddr - 2 // rel8: Opcode (1) + Offset (1) = 2 bytes
		machineCode = []byte{opcode, byte(offset)}
	case 2: // rel16
		offset := destAddr - currentAddr - 4 // rel16: Opcode (2) + Offset (2) = 4 bytes
		machineCode = []byte{0x0f, opcode + 0x10, byte(offset), byte(offset >> 8)}
	default: // rel32
		offset := destAddr - currentAddr - 6 // rel32: Opcode (2) + Offset (4) = 6 bytes
		machineCode = []byte{
			0x0f,
			opcode + 0x10, // Jcc rel32 opcode (e.g., 0x87 for JA)
			byte(offset),
			byte(offset >> 8),
			byte(offset >> 16),
			byte(offset >> 24),
		}
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
