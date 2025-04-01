package pass1

import (
	"fmt"
	"log"

	"github.com/HobbyOSs/gosk/internal/ast" // astパッケージをインポート
	"github.com/HobbyOSs/gosk/pkg/cpu"
)

// estimateJumpSize は Pass 1 での near ジャンプ/コール命令のサイズを推定します。
// これは推定値です。最終的なオフセットサイズ (rel8/rel16/32) は Pass 2 で変更される可能性があるためです。
func estimateJumpSize(instName string, bitMode cpu.BitMode) int32 {
	isJcc := instName != "JMP" && instName != "CALL" // その他は Jcc と仮定します

	// デフォルトでは near 相対ジャンプ/コールサイズ (オペコード + rel16/32) になります
	// JMP rel16/32 (E9 cw/cd): 1 + 2/4 = 3/5 バイト
	// CALL rel16/32 (E8 cw/cd): 1 + 2/4 = 3/5 バイト
	// Jcc rel16/32 (0F 8x cw/cd): 2 + 2/4 = 4/6 バイト
	size := int32(5) // 最初は JMP/CALL に rel32 を仮定します
	if bitMode == cpu.MODE_16BIT {
		size = 3 // 16 ビットモードでは JMP/CALL に rel16 を仮定します
	}

	if isJcc {
		size = 6 // 最初は Jcc に rel32 を仮定します
		if bitMode == cpu.MODE_16BIT {
			size = 4 // 16 ビットモードでは Jcc に rel16 を仮定します
		}
	}
	// 注意: ここでは short ジャンプ (rel8) は推定しません。Pass 2 で可能であれば最適化されます。
	return size
}

// processCalcJcc は JMP および条件付きジャンプ命令を処理します。
func processCalcJcc(env *Pass1, operands []ast.Exp, instName string) {
	if len(operands) != 1 {
		log.Printf("Error: %s instruction requires exactly one operand, got %d", instName, len(operands))
		return
	}

	operand := operands[0]
	evaluatedOperand, _ := operand.Eval(env) // 最初にオペランドを評価し、'evaluated' フラグを明示的に無視します

	// *評価された* オペランドタイプに基づいて推定サイズを決定し、Ocode を発行します
	var estimatedSize int32
	var ocode string

	switch op := evaluatedOperand.(type) {
	case *ast.SegmentExp: // FAR ジャンプを処理します (例: JMP FAR label, JMP seg:off)
		log.Printf("[pass1] Processing evaluated SegmentExp for %s: %s", instName, op.TokenLiteral())
		// セグメントとオフセット部分を *再度* 評価します (SegmentExp 自体の Eval では完全には解決されない場合があります)
		segEval, segOk := op.Left.Eval(env)
		offEval, offOk := op.Right.Eval(env)

		if !segOk || !offOk {
			log.Printf("Error: Could not fully evaluate segment or offset for FAR %s.", instName)
			// 評価に失敗した場合は、元のオペランド文字列に基づいてプレースホルダーを発行します
			estimatedSize = 7 // ptr16:32 を仮定します
			ocode = fmt.Sprintf("%s {{expr:%s}}", instName, operand.TokenLiteral())
		} else {
			// 数値の場合は値を取得し、それ以外の場合はプレースホルダーを使用します
			var segStr, offStr string
			if segNum, ok := segEval.(*ast.NumberExp); ok {
				segStr = fmt.Sprintf("%d", segNum.Value)
			} else {
				segStr = fmt.Sprintf("{{expr:%s}}", op.Left.TokenLiteral()) // セグメントのプレースホルダー
			}
			if offNum, ok := offEval.(*ast.NumberExp); ok {
				offStr = fmt.Sprintf("%d", offNum.Value)
			} else {
				offStr = fmt.Sprintf("{{expr:%s}}", op.Right.TokenLiteral()) // オフセットのプレースホルダー
			}

			estimatedSize = 7 // JMP ptr16:32 (EA + ptr16:32)
			ocode = fmt.Sprintf("%s_FAR %s:%s", instName, segStr, offStr)
		}

	case *ast.ImmExp:
		if factor, ok := op.Factor.(*ast.IdentFactor); ok { // 未解決のラベル
			label := factor.Value
			log.Printf("[pass1] Processing label '%s' for %s", label, instName)
			// ラベルが存在するか、およびその値を確認します
			if addr, exists := env.SymTable[label]; exists {
				log.Printf("debug: [processCalcJcc] Label '%s' found in SymTable with address 0x%x (%d)", label, addr, addr)
			} else {
				log.Printf("debug: [processCalcJcc] Label '%s' not found in SymTable yet. Adding placeholder.", label)
				env.SymTable[label] = 0 // プレースホルダーアドレス
			}
			estimatedSize = estimateJumpSize(instName, env.BitMode)
			ocode = fmt.Sprintf("%s {{.%s}}", instName, label) // ラベルプレースホルダー付きの Ocode
		} else {
			// ImmExp.Eval が正しく機能すれば発生しないはずですが、防御的に処理します
			log.Printf("Error: Unexpected factor type %T within evaluated ImmExp for %s.", op.Factor, instName)
			estimatedSize = estimateJumpSize(instName, env.BitMode)                 // とにかくサイズを推定します
			ocode = fmt.Sprintf("%s {{expr:%s}}", instName, operand.TokenLiteral()) // 元の式を持つプレースホルダー
		}

	case *ast.NumberExp: // 解決された即値アドレス
		targetAddr := op.Value
		log.Printf("[pass1] Processing immediate address %d (0x%x) for %s", targetAddr, targetAddr, instName)
		// Pass 1 では相対オフセットを確実に計算できません。プレースホルダーを使用します。
		estimatedSize = estimateJumpSize(instName, env.BitMode)
		// Pass 2 用に即値アドレスを示すプレースホルダーを使用します
		ocode = fmt.Sprintf("%s {{addr:%d}}", instName, targetAddr)

	case *ast.AddExp, *ast.MultExp: // 部分的に評価された式 (例: label + offset)
		log.Printf("[pass1] Processing partially evaluated expression for %s: %s", instName, op.TokenLiteral())
		// Pass 1 では完全に解決できません。Pass 2 用にプレースホルダーを使用します。
		estimatedSize = estimateJumpSize(instName, env.BitMode)
		ocode = fmt.Sprintf("%s {{expr:%s}}", instName, op.TokenLiteral()) // 式文字列を持つプレースホルダー

	// MemoryAddrExp の処理が必要な場合 (例: JMP DWORD PTR [EAX])

	default:
		log.Printf("Error: Invalid evaluated operand type %T for %s instruction.", evaluatedOperand, instName)
		// フォールバックプレースホルダーとして元のオペランド文字列を使用しようとします
		estimatedSize = estimateJumpSize(instName, env.BitMode) // サイズを推定します
		ocode = fmt.Sprintf("%s {{expr:%s}}", instName, operand.TokenLiteral())
	}

	// LOC を更新し、Ocode を発行します
	if estimatedSize == 0 {
		log.Printf("WARN: Estimated size is 0 for %s %s. Defaulting to 2.", instName, operand.TokenLiteral())
		estimatedSize = 2 // LOC が進まないのを避けます
	}
	env.LOC += estimatedSize
	env.Client.Emit(ocode)
}

// getOffsetSize は相対オフセットのサイズ (バイト数) を返す (Pass 2 で使用)
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
