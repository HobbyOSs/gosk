package codegen

import (
	"log"
	"strconv"
	"strings"

	"github.com/morikuni/failure"
)

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

func handleRESB(args []string, currentBufferSize int) []byte {
	var binary []byte

	if strings.Contains(args[0], `-$`) {
		rangeOfResb := args[0][:len(args[0])-len(`-$`)]
		reserveSize, err := strconv.ParseInt(rangeOfResb, 0, 32)
		if err != nil {
			log.Fatal(failure.Wrap(err))
		}

		needToAppendSize := reserveSize - int64(currentBufferSize)
		binary = append(binary, make([]byte, needToAppendSize)...)
		return binary
	}

	reserveSize, err := strconv.ParseInt(args[0], 0, 32)
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}

	binary = append(binary, make([]byte, reserveSize)...)
	return binary
}
