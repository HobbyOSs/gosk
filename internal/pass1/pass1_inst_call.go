package pass1

import (
	"fmt"
	"log"

	"github.com/HobbyOSs/gosk/internal/token"
)

func processCALL(env *Pass1, tokens []*token.ParseToken, instName string) {
	if len(tokens) != 1 {
		log.Fatalf("%s instruction requires exactly one operand, got %d", instName, len(tokens))
		return
	}

	arg := tokens[0]

	switch arg.TokenType {
	case token.TTIdentifier:
		label := arg.AsString()

		// ラベルをSymTableに登録 (仮アドレスを割り当てる)
		if _, ok := env.SymTable[label]; !ok {
			env.SymTable[label] = 0 // Pass 1では仮アドレス
		}
		// CALL rel32 は 5 bytes
		env.LOC += 5

		// Ocodeを生成 (ジャンプ先アドレスはプレースホルダー)
		// プレースホルダーとしてラベルを使用
		env.Client.Emit(fmt.Sprintf("%s {{.%s}}", instName, label))
	case token.TTNumber, token.TTHex:
		// 機械語サイズを計算
		offset := arg.ToInt32() - int32(env.DollarPosition)
		env.LOC += getOffsetSize(offset)
		env.LOC += 1

		// ダミーのラベルを作る
		fakeLabel := fmt.Sprintf("imm_call_%d", env.NextImmJumpID)
		env.NextImmJumpID++
		env.SymTable[fakeLabel] = arg.ToInt32()

		// Ocodeを生成 (ジャンプ先アドレスはダミー)
		env.Client.Emit(fmt.Sprintf("%s {{.%s}}", instName, fakeLabel))
	default:
		log.Fatalf("invalid CALL operand: %v", arg)
	}
}
