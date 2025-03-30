package pass1

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/morikuni/failure"
)

func processALIGNB(env *Pass1, tokens []*token.ParseToken) {
	arg := tokens[0]
	unit := arg.ToInt32()
	var loc int32 = 0
	if env.LOC%unit != 0 { // 現在のLOCが境界に揃っていない場合のみ計算
		nearestSize := env.LOC/unit + 1
		loc = nearestSize*unit - env.LOC
	}
	env.LOC += loc
	// ALIGNB は ocode を発行しない (LOC調整のみ)
}

func processDB(env *Pass1, tokens []*token.ParseToken) {
	var loc int32 = 0
	ocodes := []int32{}

	for _, t := range tokens {
		switch t.TokenType {
		case token.TTNumber, token.TTHex:
			loc += 1
			ocodes = append(ocodes, t.ToInt32())
		case token.TTIdentifier:
			ident := t.AsString()
			if labelLOC, ok := env.SymTable[ident]; ok { // ラベルが存在する場合
				loc += 1                               // DBは1バイト
				ocodes = append(ocodes, labelLOC&0xFF) // ラベルアドレスの下位1バイトを追加
			} else { // ラベルが存在しない場合は文字列として扱う
				loc += int32(len(ident))
				ocodes = append(stringToOcodes(ident), ocodes...)
			}
		default:
			log.Fatalf("unsupported token type in DB: %v", t.TokenType)
		}
	}

	env.LOC += loc
	emitCommand(env, "DB", ocodes)
}

func processDW(env *Pass1, tokens []*token.ParseToken) {
	var loc int32 = 0
	ocodes := []int32{}

	for _, t := range tokens {
		switch t.TokenType {
		case token.TTNumber, token.TTHex:
			loc += 2
			ocodes = append(ocodes, t.ToInt32())
		case token.TTIdentifier:
			ident := t.AsString()
			if labelLOC, ok := env.SymTable[ident]; ok { // ラベルが存在する場合
				loc += 2                          // DWは2バイト
				ocodes = append(ocodes, labelLOC) // ラベルアドレス(int32)をそのまま追加
			} else { // ラベルが存在しない場合は文字列として扱う
				loc += int32(len(ident))
				ocodes = append(stringToOcodes(ident), ocodes...)
			}
		default:
			log.Fatalf("unsupported token type in DW: %v", t.TokenType)
		}
	}

	env.LOC += loc
	emitCommand(env, "DW", ocodes)
}

func processDD(env *Pass1, tokens []*token.ParseToken) {
	var loc int32 = 0
	ocodes := []int32{}

	for _, t := range tokens {
		switch t.TokenType {
		case token.TTNumber, token.TTHex:
			loc += 4
			ocodes = append(ocodes, t.ToInt32())
		case token.TTIdentifier:
			labelName := t.AsString()
			labelLOC, ok := env.SymTable[labelName] // SymTable を使用
			if !ok {
				log.Fatalf("label not found: %s", labelName) // log.Fatalf を使用
			}
			loc += 4
			ocodes = append(ocodes, labelLOC) // ラベルのアドレスを追加
		default:
			log.Fatalf("unsupported token type in DD: %v", t.TokenType) // log.Fatalf を使用
		}
	}

	env.LOC += loc
	emitCommand(env, "DD", ocodes)
}

func processORG(env *Pass1, tokens []*token.ParseToken) {
	arg := tokens[0].AsString()
	size, err := strconv.ParseInt(arg, 0, 32)
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}
	env.LOC = int32(size)
	env.DollarPosition += uint32(size) // エントリーポイントのアドレスを加算
}

func processRESB(env *Pass1, tokens []*token.ParseToken) {
	arg := tokens[0].AsString()

	if strings.Contains(arg, `-$`) {
		rangeOfResb := arg[:len(arg)-len(`-$`)]
		reserveSize, err := strconv.ParseInt(rangeOfResb, 0, 32)
		if err != nil {
			log.Fatal(failure.Wrap(err))
		}
		env.LOC += int32(reserveSize)
		env.Client.Emit(fmt.Sprintf("RESB %d-$", reserveSize))
		return
	}

	size, err := strconv.ParseInt(arg, 10, 32)
	if err != nil {
		log.Fatal(failure.Wrap(err))
	}
	env.LOC += int32(size)
	env.Client.Emit(fmt.Sprintf("RESB %d", size))
}
