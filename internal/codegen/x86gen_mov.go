package codegen

import (
	"log"
	"strconv"
	"strings"

	"github.com/HobbyOSs/gosk/pkg/asmdb"
	"github.com/HobbyOSs/gosk/pkg/operand"
)

// handleMOV handles the MOV instruction and generates the appropriate machine code.
func handleMOV(operands []string) []byte {
	if len(operands) != 2 {
		log.Printf("error: MOV requires 2 operands, got %d", len(operands))
		return nil
	}

	// オペランドの解析
	ops := operand.NewOperandFromString(strings.Join(operands, ","))

	// AsmDBからエンコーディングを取得
	db := asmdb.NewInstructionDB()
	encoding, err := db.FindEncoding("MOV", ops)
	if err != nil {
		log.Printf("error: Failed to find encoding: %v", err)
		return nil
	}

	// エンコーディング情報を使用して機械語を生成
	machineCode := make([]byte, 0)

	// プレフィックスの追加
	if ops.Require66h() {
		machineCode = append(machineCode, 0x66)
	}

	// オペコードの追加
	opcodeByte, err := strconv.ParseUint(encoding.Opcode.Byte, 16, 8)
	if err != nil {
		log.Printf("error: Failed to parse opcode byte: %v", err)
		return nil
	}
	machineCode = append(machineCode, byte(opcodeByte))

	// ModR/Mの追加（必要な場合）
	if encoding.ModRM != nil {
		modrm := generateModRM(encoding.ModRM)
		machineCode = append(machineCode, modrm)
	}

	// 即値の追加（必要な場合）
	if encoding.Immediate != nil {
		if imm, err := getImmediateValue(operands[1], encoding.Immediate.Size); err == nil {
			machineCode = append(machineCode, imm...)
		}
	}

	log.Printf("debug: Generated machine code: % x", machineCode)

	return machineCode
}

// generateModRM generates ModR/M byte based on ModRM encoding
func generateModRM(modrm *asmdb.Modrm) byte {
	// ModR/Mバイトの生成
	// |  mod  |  reg  |  r/m  |
	// | 7 6 | 5 4 3 | 2 1 0 |

	// modeの解析（2ビット）
	var mode byte
	switch modrm.Mode {
	case "11":
		mode = 0b11000000
	default:
		mode = 0
	}

	// regの解析（3ビット）
	reg, _ := strconv.ParseUint(modrm.Reg, 2, 3)
	regBits := byte(reg) << 3

	// r/mの解析（3ビット）
	rm, _ := strconv.ParseUint(modrm.Rm, 2, 3)
	rmBits := byte(rm)

	return mode | regBits | rmBits
}

// getImmediateValue extracts immediate value from operand
func getImmediateValue(operand string, size int) ([]byte, error) {
	// 0xで始まる16進数の場合
	if strings.HasPrefix(operand, "0x") {
		value, err := strconv.ParseUint(operand[2:], 16, size*8)
		if err != nil {
			return nil, err
		}
		return intToBytes(value, size), nil
	}

	// 10進数の場合
	value, err := strconv.ParseInt(operand, 10, size*8)
	if err != nil {
		return nil, err
	}
	return intToBytes(uint64(value), size), nil
}

// intToBytes converts an integer to a byte slice of specified size
func intToBytes(value uint64, size int) []byte {
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		bytes[i] = byte(value >> (i * 8))
	}
	return bytes
}
