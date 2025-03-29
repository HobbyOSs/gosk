package test

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/comail/colog"
	"github.com/samber/lo"
)

func setUpColog(debug bool) {
	colog.Register()
	colog.SetDefaultLevel(colog.LInfo) // Keep default level as Info
	colog.SetMinLevel(colog.LDebug)    // テスト時は常に Debug レベル
	colog.SetFlags(log.Lshortfile)
	colog.SetFormatter(&colog.StdFormatter{Colors: false})
}

// 16進数文字列をテスト用に用意するDSL
func defineHEX(dsl []string) []byte {
	var result []byte
	for _, line := range dsl {
		// CSV Reader を作成
		var inQuotes bool
		// FieldsFunc でカスタム区切りルール
		tokens := strings.FieldsFunc(line, func(r rune) bool {
			if r == '"' {
				inQuotes = !inQuotes // クォートの開始・終了
				return false         // クォート自体は区切りにしない
			}
			// クォート外のスペースだけ区切りとする
			return !inQuotes && (r == ' ' || r == '\t')
		})

		// クォートを除去
		for i, field := range tokens {
			tokens[i] = strings.Trim(field, `"`)
		}

		switch tokens[0] {
		case "#", ";":
			continue
		case "DATA":
			// DATA命令の処理
			for _, str := range tokens[1:] {
				if bin, err := strconv.ParseUint(str, 0, 8); err == nil {
					result = append(result, byte(bin))
				} else {
					bytes := lo.Map([]rune(str), func(r rune, _ int) byte {
						return byte(r)
					})
					result = append(result, bytes...)
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
