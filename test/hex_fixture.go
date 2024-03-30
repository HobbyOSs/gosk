package test

import (
	"fmt"
	"strconv"
	"strings"
)

// 16進数文字列をテスト用に用意するDSL
func defineHEX(dsl []string) []byte {
	var result []byte
	for _, line := range dsl {
		tokens := strings.Fields(line)
		switch tokens[0] {
		case "#", ";":
			continue
		case "DATA":
			// DATA命令の処理
			for _, hexStr := range tokens[1:] {
				if val, err := strconv.ParseUint(hexStr, 0, 8); err == nil {
					result = append(result, byte(val))
				} else {
					fmt.Println("Error parsing DATA value:", err)
				}
			}
		case "FILL":
			// FILL命令の処理
			num, err := strconv.Atoi(tokens[1])
			if err != nil {
				fmt.Println("Error parsing FILL value:", err)
				continue
			}
			var fillByte byte = 0x00 // デフォルト値
			if len(tokens) > 2 {
				if val, err := strconv.ParseUint(tokens[2], 0, 8); err == nil {
					fillByte = byte(val)
				} else {
					fmt.Println("Error parsing FILL byte value:", err)
				}
			}
			for i := 0; i < num; i++ {
				result = append(result, fillByte)
			}
		}
	}
	return result
}
