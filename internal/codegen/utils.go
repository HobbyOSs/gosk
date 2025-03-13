package codegen

import (
	"strconv"
	"strings"
)

// getImmediateValue extracts immediate value from operand
func getImmediateValue(operand string, size int) ([]byte, error) {
	// 0xで始まる16進数の場合
	if strings.HasPrefix(operand, "0x") {
		value, err := strconv.ParseUint(operand[2:], 16, size*8)
		if err != nil {
			return make([]byte, size), nil
		}
		return intToBytes(value, size), nil
	}

	// 10進数の場合
	value, err := strconv.ParseInt(operand, 10, size*8)
	if err != nil {
		return make([]byte, size), nil
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
