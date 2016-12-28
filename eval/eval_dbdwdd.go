package eval

import (
	"fmt"
	"github.com/hangingman/gosk/ast"
	"github.com/hangingman/gosk/object"
	"github.com/hangingman/gosk/token"
	"log"
	"strconv"
	"strings"
)

// 複雑なオペコードはファイルを分割する方針で
func evalDBStatement(stmt *ast.MnemonicStatement) object.Object {
	return evalDStatements(stmt, int2Byte, "BYTE")
}

func evalDWStatement(stmt *ast.MnemonicStatement) object.Object {
	return evalDStatements(stmt, int2Word, "WORD")
}

func evalDDStatement(stmt *ast.MnemonicStatement) object.Object {
	return evalDStatements(stmt, int2Dword, "DWORD")
}

func evalDStatements(stmt *ast.MnemonicStatement, f func(int) []byte, usingType string) object.Object {
	toks := []string{}
	bytes := []byte{}

	for _, tok := range stmt.Name.Tokens {
		if tok.Type == token.HEX_LIT {
			switch {
			case IsImm8(tok):
				bytes = append(bytes, imm8ToByte(tok)...)
			case IsImm16(tok):
				bytes = append(bytes, imm16ToWord(tok)...)
			case IsImm32(tok):
				bytes = append(bytes, imm32ToDword(tok)...)
			default:
				// do nothing
			}
		} else if tok.Type == token.STR_LIT {
			// "を取り除いて処理
			strLength := len(tok.Literal)
			bs := []byte(tok.Literal[1 : strLength-1])
			bytes = append(bytes, bs...)
		} else if tok.Type == token.INT {
			// Go言語のintは常にint64 -> uint8
			int64Val, _ := strconv.Atoi(tok.Literal)
			bs := f(int64Val)
			bytes = append(bytes, bs...)
		} else if tok.Type == token.IDENT {
			if from, ok := labelManage.labelBytesMap[tok.Literal]; ok {
				// ラベルが見つかっていればバイト数を計算して設定する
				absSize := from + dollarPosition
				// FIX: ここもなんかおかしい、naskはDDなのにWordでデータ入れてる
				bytes = append(bytes, int2Word(absSize)...)
			}
		}
		toks = append(toks, fmt.Sprintf("%s: %s", tok.Type, tok.Literal))
	}

	log.Println(fmt.Sprintf("info: [%s]", strings.Join(toks, ", ")))
	stmt.Bin = &object.Binary{Value: bytes}
	return stmt.Bin
}
