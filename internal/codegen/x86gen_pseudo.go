package codegen

import "strconv"

// 配列型で渡される文字列はすべて数値型を想定する
// 0x00のようなhex notationや "hello" のような文字列はastで事前に数値の配列に変換
func handleDB(args []string) []byte {
	var binary []byte
	for _, arg := range args {
		num, _ := strconv.Atoi(arg)
		binary = append(binary, byte(num))
	}
	return binary
}

func handleDW(args []string) []byte {
	var binary []byte
	for _, arg := range args {
		num, _ := strconv.Atoi(arg)
		binary = append(binary, byte(num&0xFF), byte((num>>8)&0xFF))
	}
	return binary
}

func handleDD(args []string) []byte {
	var binary []byte
	for _, arg := range args {
		num, _ := strconv.Atoi(arg)
		binary = append(binary, byte(num&0xFF), byte((num>>8)&0xFF), byte((num>>16)&0xFF), byte((num>>24)&0xFF))
	}
	return binary
}
