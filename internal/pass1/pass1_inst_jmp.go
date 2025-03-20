package pass1

import (
	"fmt"
	"log"

	"github.com/HobbyOSs/gosk/internal/token"
)

func processCalcJcc(env *Pass1, tokens []*token.ParseToken, instName string) {

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
		// Forward reference（前方参照）の問題により、Pass1フェーズではオフセットを正確に計算できない
		// そのため、現状はrel8（2バイト）を仮定し、必要に応じてPass2フェーズで調整する
		//
		// 例：
		//   JMP label   ; ラベルが前方にある場合、この時点でラベルの位置が不明
		//   ...
		//   label:      ; ラベルの実際の位置はPass2まで確定しない

		// 機械語サイズを計算 (JMP rel8 は 2 bytes)
		env.LOC += 2

		// Ocodeを生成 (ジャンプ先アドレスはプレースホルダー)
		// プレースホルダーとしてラベルを使用
		env.Client.Emit(fmt.Sprintf("%s {{.%s}}", instName, label))
	case token.TTNumber, token.TTHex:
		// 機械語サイズを計算
		offset := arg.ToInt32() - int32(env.DollarPosition)
		env.LOC += getOffsetSize(offset)
		env.LOC += 1

		// ダミーのラベルを作る
		fakeLabel := fmt.Sprintf("imm_jmp_%d", env.NextImmJumpID)
		env.NextImmJumpID++
		env.SymTable[fakeLabel] = arg.ToInt32()

		// Ocodeを生成 (ジャンプ先アドレスはダミー)
		env.Client.Emit(fmt.Sprintf("%s {{.%s}}", instName, fakeLabel))
	default:
		log.Fatalf("invalid JMP operand: %v", arg)
	}
}

// -128～127, -32768～32767 などの判定に使う
func getOffsetSize(imm int32) int32 {
	if imm >= -0x80 && imm <= 0x7f {
		return 1
	}
	if imm >= -0x8000 && imm <= 0x7fff {
		return 2
	}
	return 4
}
