package codegen

import (
	"log"
	"strconv"
	"strings"
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
