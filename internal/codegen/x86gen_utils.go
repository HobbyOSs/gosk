package codegen

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/HobbyOSs/gosk/pkg/asmdb"
)

// GenerateModRM generates ModR/M byte based on mode, reg, and rm.
func GenerateModRM(modeStr string, regOperand string, rmOperand string) byte {
	// ModR/Mバイトの生成
	// |  mod  |  reg  |  r/m  |
	// | 7 6 | 5 4 3 | 2 1 0 |

	// modeの解析（2ビット）
	var mode byte
	switch modeStr {
	case "11":
		mode = 0b11000000
	case "00":
		mode = 0b00000000
	case "01":
		mode = 0b01000000
	case "10":
		mode = 0b10000000
	default:
		mode = 0 // デフォルト値
	}

	// regの解析（3ビット）
	reg, err := GetRegisterNumber(regOperand)
	if err != nil {
		log.Printf("error: Failed to get register number for reg: %v", err)
		return 0
	}
	regBits := byte(reg) << 3

	// r/mの解析（3ビット）
	// メモリの場合は0として扱う
	if strings.HasPrefix(rmOperand, "[") && strings.HasSuffix(rmOperand, "]") {
		// TODO: メモリオペランドの解析
		rmBits := byte(0)
		return mode | regBits | rmBits
	}

	rm, err := GetRegisterNumber(rmOperand)
	if err != nil {
		log.Printf("error: Failed to get register number for rm: %v", err)
		return 0
	}
	rmBits := byte(rm)

	result := mode | regBits | rmBits
	log.Printf("debug: GenerateModRM: mode=%s(%b), reg=%s(%b), rm=%s(%b), result=%#x", modeStr, mode, regOperand, regBits, rmOperand, rmBits, result)
	return result
}

// GetRegisterNumber はレジスタ名からレジスタ番号（0-7）を取得する
func GetRegisterNumber(regName string) (int, error) {
	switch regName {
	case "AL", "AX", "EAX", "RAX", "ES":
		return 0, nil
	case "CL", "CX", "ECX", "RCX", "CS":
		return 1, nil
	case "DL", "DX", "EDX", "RDX", "SS":
		return 2, nil
	case "BL", "BX", "EBX", "RBX", "DS":
		return 3, nil
	case "AH", "SP", "ESP", "RSP", "FS":
		return 4, nil
	case "CH", "BP", "EBP", "RBP", "GS":
		return 5, nil
	case "DH", "SI", "ESI", "RSI":
		return 6, nil
	case "BH", "DI", "EDI", "RDI":
		return 7, nil
	default:
		return 0, fmt.Errorf("unknown register: %s", regName)
	}
}

// ResolveOpcode はOpcodeとレジスタ番号を受け取り、最終的なオペコードを算出する。
// regNum はレジスタの番号（0-7）を表す。
func ResolveOpcode(op asmdb.Opcode, regNum int) (byte, error) {
	// Byteを16進数文字列から数値に変換
	opcodeByte, err := strconv.ParseUint(op.Byte, 16, 8)
	if err != nil {
		return 0, fmt.Errorf("invalid opcode byte: %v", err)
	}

	// Addendがnilなら基本オペコードをそのまま返す
	if op.Addend == nil {
		return byte(opcodeByte), nil
	}

	// レジスタ番号の下位3ビットを取得してORする
	regBits := byte(regNum & 0x07)
	result := byte(opcodeByte) | regBits

	log.Printf("debug: ResolveOpcode: base=%#x, addend=%v, reg=%d, result=%#x", opcodeByte, op.Addend, regNum, result)
	return result, nil
}

// getModRMFromOperands はオペランドからModR/Mバイトを生成する
func getModRMFromOperands(operands []string, modRM *asmdb.Encoding) (byte, error) {
	if modRM == nil || modRM.ModRM == nil {
		return 0, nil
	}

	// ModR/M バイトを生成
	modrmByte := GenerateModRM(modRM.ModRM.Mode, operands[0], operands[1])

	return modrmByte, nil
}
