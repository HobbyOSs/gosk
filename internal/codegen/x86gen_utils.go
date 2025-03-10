package codegen

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/HobbyOSs/gosk/pkg/asmdb"
)

// GenerateModRM generates ModR/M byte based on mode, reg, and rm strings.
func GenerateModRM(modeStr string, regStr string, rmStr string) byte {
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
	regStr = strings.TrimPrefix(regStr, "#")
	reg, _ := strconv.ParseUint(regStr, 10, 3) // 10進数として解釈
	regBits := byte(reg) << 3

	// r/mの解析（3ビット）
	rmStr = strings.TrimPrefix(rmStr, "#")
	rm, _ := strconv.ParseUint(rmStr, 10, 3) // 10進数として解釈
	rmBits := byte(rm)

	result := mode | regBits | rmBits
	log.Printf("debug: GenerateModRM: mode=%s(%b), reg=%s(%b), rm=%s(%b), result=%#x", modeStr, mode, regStr, regBits, rmStr, rmBits, result)
	return result
}

// GetRegisterNumber はレジスタ名からレジスタ番号（0-7）を取得する
func GetRegisterNumber(regName string) (int, error) {
	switch regName {
	case "AL", "AX", "EAX":
		return 0, nil
	case "CL", "CX", "ECX":
		return 1, nil
	case "DL", "DX", "EDX":
		return 2, nil
	case "BL", "BX", "EBX":
		return 3, nil
	case "AH", "SP", "ESP":
		return 4, nil
	case "CH", "BP", "EBP":
		return 5, nil
	case "DH", "SI", "ESI":
		return 6, nil
	case "BH", "DI", "EDI":
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

	log.Printf("debug: ResolveOpcode: base=%#x, addend=%v, reg=%d, result=%#x",
		opcodeByte, op.Addend, regNum, result)
	return result, nil
}
