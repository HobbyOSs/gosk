package pass1

import (
	"fmt"
	"log"

	"github.com/HobbyOSs/gosk/internal/ast" // astパッケージをインポート
	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/HobbyOSs/gosk/pkg/cpu" // cpuパッケージをインポート
)

// evalSimpleExp は式ノードを評価し、結果をint32で返すヘルパー関数 (TraverseASTを呼び出す形式に修正)
func evalSimpleExp(exp ast.Node, env *Pass1) (int32, error) {
	// 式ノードをTraverseASTで評価し、結果をスタックに積む
	TraverseAST(exp, env)
	// スタックから評価結果を取得
	resultToken := pop(env)

	// 結果トークンが数値として解釈できるか確認し、int32で返す
	// TODO: ラベル参照などの解決が必要な場合に対応
	if resultToken.IsNumber() || resultToken.TokenType == token.TTHex {
		return resultToken.ToInt32(), nil
	} else if resultToken.TokenType == token.TTIdentifier {
		// TODO: ラベルやEQUの解決
		log.Printf("WARN: Identifier '%s' evaluation in evalSimpleExp is not implemented yet.", resultToken.AsString())
		return 0, fmt.Errorf("identifier evaluation not implemented: %s", resultToken.AsString())
	}

	return 0, fmt.Errorf("cannot evaluate expression result to int32: %v (Type: %s)", resultToken.Data, resultToken.TokenType)
}

func processCalcJcc(env *Pass1, tokens []*token.ParseToken, instName string) {

	if len(tokens) != 1 {
		log.Fatalf("%s instruction requires exactly one operand, got %d", instName, len(tokens))
		return
	}

	arg := tokens[0]

	// SegmentExpの場合の処理を追加
	// TokenTypeもチェックするように修正
	if arg.TokenType == token.TTIdentifier { // TokenTypeがIdentifierであることを確認 (TraverseASTの暫定対応)
		if segExp, ok := arg.Data.(*ast.SegmentExp); ok {
			// SegmentExpの処理
			log.Printf("[pass1] Processing SegmentExp for %s: %s", instName, segExp.TokenLiteral())

			// Left (セグメント) と Right (オフセット) を評価
			segment, err := evalSimpleExp(segExp.Left, env)
			if err != nil {
				log.Fatalf("Failed to evaluate segment expression for %s: %v", instName, err)
			}

			if segExp.Right == nil {
				// JMP DWORD label のようなケース (Rightがnil) は現状未対応
				log.Fatalf("SegmentExp without Right part is not supported for JMP FAR: %s", segExp.TokenLiteral())
			}
			offset, err := evalSimpleExp(segExp.Right, env)
			if err != nil {
				log.Fatalf("Failed to evaluate offset expression for %s: %v", instName, err)
			}

			// 機械語サイズを計算
			// JMP ptr16:32 は通常 7 バイト (EA + 4バイトオフセット + 2バイトセレクタ)
			// 16ビットモードではオペランドサイズプレフィックス(66h)が付くため 8 バイト
			size := int32(7)
			if env.BitMode == cpu.MODE_16BIT {
				size = 8
			}
			env.LOC += size

			// Ocodeを生成 (セグメントとオフセットをコロン区切りで1つのオペランドとして指定)
			env.Client.Emit(fmt.Sprintf("%s_FAR %d:%d", instName, segment, offset)) // 例: JMP_FAR 16:27
			return                                                                  // SegmentExpの処理はここで終了
		} // 追加した if segExp, ok := ... の閉じ括弧
	}

	// 既存のラベル・数値処理
	switch arg.TokenType {
	case token.TTIdentifier:
		// ラベル参照の場合 (SegmentExpではないIdentFactorなど)
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

		// 機械語サイズを計算 (JMP/Jcc rel8 は 2 bytes)
		env.LOC += 2 // TODO: Jcc命令の種類やオフセットサイズによってサイズが変わる可能性

		// Ocodeを生成 (ジャンプ先アドレスはプレースホルダー)
		// プレースホルダーとしてラベルを使用
		env.Client.Emit(fmt.Sprintf("%s {{.%s}}", instName, label))
	case token.TTNumber, token.TTHex:
		// 数値参照の場合
		// 機械語サイズを計算
		offsetVal := arg.ToInt32()
		relativeOffset := offsetVal - int32(env.DollarPosition) // 相対オフセット計算
		// JMP/Jcc命令のオペコードサイズは1バイトまたは2バイト (0F xx)
		// rel8: opcode(1) + offset(1) = 2 bytes
		// rel16/32: opcode(1 or 2) + offset(2 or 4)
		offsetSize := getOffsetSize(relativeOffset)
		opcodeSize := int32(1)                   // JMP rel8/rel16/rel32 は E9/EB (1 byte)
		if offsetSize > 1 && instName != "JMP" { // Jcc rel16/32 は 0F 8x (2 bytes)
			opcodeSize = 2
		}
		env.LOC += opcodeSize + offsetSize

		// ダミーのラベルを作る
		fakeLabel := fmt.Sprintf("imm_jmp_%d", env.NextImmJumpID)
		env.NextImmJumpID++
		env.SymTable[fakeLabel] = offsetVal // ジャンプ先の絶対アドレスを登録

		// Ocodeを生成 (ジャンプ先アドレスはダミー)
		// 追加: 即値のログ出力
		log.Printf("[pass1] %s immediate value: %d (0x%x)", instName, offsetVal, offsetVal)

		// Ocodeを生成 (ジャンプ先アドレスはダミー)
		env.Client.Emit(fmt.Sprintf("%s {{.%s}}", instName, fakeLabel))
	default:
		// SegmentExpがTTIdentifierとして渡ってきたが、上記のSegmentExp処理でハンドルされなかった場合など
		log.Fatalf("invalid JMP operand type or unhandled case: %T, value: %v", arg.Data, arg)
	}
}

// getOffsetSize は相対オフセットのサイズ (バイト数) を返す
func getOffsetSize(imm int32) int32 {
	if imm >= -128 && imm <= 127 {
		return 1 // rel8
	}
	// rel16の判定を追加
	if imm >= -32768 && imm <= 32767 {
		return 2 // rel16
	}
	return 4 // rel32
}
