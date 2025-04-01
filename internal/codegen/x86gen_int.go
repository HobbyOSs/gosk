package codegen

import (
	"strconv"

	"github.com/HobbyOSs/gosk/pkg/ocode"
)

func handleINT(ocode ocode.Ocode) []byte {
	// INT命令は2バイトの命令
	// 1バイト目: 0xCD (INT命令のオペコード)
	// 2バイト目: 割り込み番号
	binary := []byte{0xCD}

	// 割り込み番号を取得
	if len(ocode.Operands) != 1 {
		panic("INT instruction requires one operand")
	}

	// 0xを除去して16進数として解析
	intNum := ocode.Operands[0]
	if len(intNum) > 2 && intNum[:2] == "0x" {
		intNum = intNum[2:]
	}
	// Parse as decimal (base 10)
	num, err := strconv.ParseInt(intNum, 10, 8) // Change base to 10
	if err != nil {
		panic("Failed to parse INT number (decimal): " + err.Error()) // Update panic message
	}

	// 割り込み番号を追加
	binary = append(binary, byte(num))

	return binary
}
