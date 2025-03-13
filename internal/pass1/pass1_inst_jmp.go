package pass1

import (
	"fmt"
	"log"
	"strings"

	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/samber/lo"
)

func processCalcJcc(env *Pass1, tokens []*token.ParseToken, instName string) {
	// JMP命令のオペランドは1つ (ジャンプ先ラベル)
	if len(tokens) != 1 {
		log.Fatalf("%s instruction requires exactly one operand, got %d", instName, len(tokens))
		return
	}

	arg := tokens[0]

	if arg.TokenType == token.TTIdentifier {
		label := arg.AsString()

		// ラベルをSymTableに登録 (仮アドレスを割り当てる)
		if _, ok := env.SymTable[label]; !ok {
			env.SymTable[label] = 0 // Pass 1では仮アドレス
		}
		// TODO: 相対ジャンプのオフセットサイズを決定 (rel8 or rel16)
		// 現状は rel8 を仮定

		// 機械語サイズを計算 (JMP rel8 は 2 bytes)
		env.LOC += 2

		// Ocodeを生成 (ジャンプ先アドレスはプレースホルダー)
		args := lo.Map(tokens, func(t *token.ParseToken, _ int) string {
			return t.AsString()
		})
		// プレースホルダーとしてラベルを使用
		env.Client.Emit(fmt.Sprintf("%s %s", instName, strings.Join(args, ",")))
		return
	}
	dataSize := checkUintRange(arg.ToUInt())
	env.LOC += int32(dataSize)
}

func checkUintRange(value uint) int {
	switch {
	case value <= uint(^uint8(0)):
		return 2
	case value <= uint(^uint16(0)):
		return 4
	case value <= uint(^uint32(0)):
		return 6
	default:
		log.Fatal("The value is larger than uint32")
	}
	return 0
}
