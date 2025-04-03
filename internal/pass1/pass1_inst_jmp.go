package pass1

import (
	"fmt"
	"log"

	"github.com/HobbyOSs/gosk/internal/ast" // astパッケージをインポート
	"github.com/HobbyOSs/gosk/pkg/cpu"
)

// estimateJumpSize は Pass 1 でのジャンプ/コール命令のサイズを推定します。
// Pass 2 で short jump に最適化される可能性を考慮し、16bit モードでは short jump サイズを推定します。
func estimateJumpSize(instName string, bitMode cpu.BitMode) int32 {
	// --- 16bit モード: short jump (2バイト) を推定 ---
	if bitMode == cpu.MODE_16BIT {
		// JMP rel8 (EB rb): 2 bytes
		// Jcc rel8 (7x rb): 2 bytes
		// CALL rel16/32 (E8 cw/cd): 3/5 bytes (CALL は short がないので near を推定)
		if instName == "CALL" {
			log.Printf("debug: [estimateJumpSize] Assuming 3 bytes (near call) for CALL in 16-bit mode.")
			return 3
		}
		log.Printf("debug: [estimateJumpSize] Assuming 2 bytes (short jump) for %s in 16-bit mode.", instName)
		return 2
	}

	// --- 32/64bit モード: near jump を推定 ---
	isJcc := instName != "JMP" && instName != "CALL"
	size := int32(5) // Default: JMP/CALL rel32 (E9/E8 cd) = 1 + 4 = 5 bytes
	if isJcc {
		size = 6 // Default: Jcc rel32 (0F 8x cd) = 2 + 4 = 6 bytes
	}
	log.Printf("debug: [estimateJumpSize] Assuming %d bytes (near jump/call) for %s in %d-bit mode.", size, instName, bitMode)
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
	case *ast.NumberExp: // ケース 1: 解決された即値アドレス
		targetAddr := op.Value
		log.Printf("[pass1] Case 1: Processing immediate address %d (0x%x) for %s", targetAddr, targetAddr, instName)
		// 数値を直接オペランドとして渡します。
		// Pass 2 のテンプレート解決は不要です。

		// 16bit モードで即値アドレスへの JMP/CALL は near (3 バイト) と推定
		if env.BitMode == cpu.MODE_16BIT {
			// JMP/CALL rel16 (E9/E8 cw) は 1 + 2 = 3 バイト
			estimatedSize = 3
			log.Printf("debug: [processCalcJcc] Assuming %d bytes (near jump/call to immediate) for %s in 16-bit mode.", estimatedSize, instName)
		} else {
			// 32/64bit モードでは estimateJumpSize を使用 (near jump/call を推定)
			estimatedSize = estimateJumpSize(instName, env.BitMode)
		}
		ocode = fmt.Sprintf("%s %d", instName, targetAddr) // 数値文字列を直接設定

	case *ast.SegmentExp: // FAR ジャンプ (seg:off)
		log.Printf("[pass1] Processing evaluated SegmentExp for %s: %s", instName, op.TokenLiteral())
		// op.Left と op.Right は SegmentExp.Eval によって既に評価されている可能性があります。
		// それらが定数に評価されるかどうかを確認します。
		segVal, segIsConst := env.GetConstValue(op.Left) // op.Left は評価済みの可能性がある *ast.AddExp
		offVal, offIsConst := 0, false
		if op.Right != nil {
			offVal, offIsConst = env.GetConstValue(op.Right) // op.Right は評価済みの可能性がある *ast.AddExp
		} else {
			// op.Right が nil の場合、オフセットがありません。FAR ジャンプには必須です。
			log.Printf("Error: FAR jump/call operand %s is missing offset.", op.TokenLiteral())
			estimatedSize = 7                                                       // ptr16:32 を仮定
			ocode = fmt.Sprintf("%s {{expr:%s}}", instName, operand.TokenLiteral()) // 元のオペランドを使用
			break                                                                   // switch case から抜ける
		}

		// ビットモードに応じてサイズを推定
		if env.BitMode == cpu.MODE_16BIT {
			estimatedSize = 8 // 66h + EA + ptr16:32
		} else {
			estimatedSize = 7 // EA + ptr16:32
		}

		// セグメントとオフセットの両方が定数に解決された場合のみ _FAR 形式を使用
		if segIsConst && offIsConst {
			segStr := fmt.Sprintf("%d", segVal)
			offStr := fmt.Sprintf("%d", offVal)
			// estimatedSize は上で設定済み
			ocode = fmt.Sprintf("%s_FAR %s:%s", instName, segStr, offStr)
			log.Printf("[pass1] Case 2a: Processing fully resolved FAR address %s:%s for %s (estimated size: %d)", segStr, offStr, instName, estimatedSize)
		} else {
			// セグメントかオフセットのどちらか、または両方が定数でない場合
			log.Printf("[pass1] Case 2b: FAR jump/call operand %s requires Pass 2 resolution for %s (estimated size: %d).", op.TokenLiteral(), instName, estimatedSize)
			// estimatedSize は上で設定済み
			// 定数解決できない場合でも _FAR サフィックスを付与し、オペランドはそのまま渡す
			ocode = fmt.Sprintf("%s_FAR %s", instName, op.TokenLiteral())
		}

	case *ast.ImmExp: // ケース 3: 即値式 (ラベルまたは '$' など)
		// ケース 3a: 単純なラベル (IdentFactor であり '$' でない)
		if factor, ok := op.Factor.(*ast.IdentFactor); ok && factor.Value != "$" {
			label := factor.Value
			log.Printf("[pass1] Case 3a: Processing label '%s' for %s", label, instName)
			// ラベルが存在しない場合はプレースホルダーを追加します
			if _, exists := env.SymTable[label]; !exists {
				log.Printf("debug: [processCalcJcc] Label '%s' not found in SymTable yet. Adding placeholder.", label)
				env.SymTable[label] = 0 // プレースホルダーアドレス
			}
			estimatedSize = estimateJumpSize(instName, env.BitMode)
			ocode = fmt.Sprintf("%s {{.%s}}", instName, label) // ラベルプレースホルダー
		} else {
			// ケース 3b: ラベルでない ImmExp (例: '$' が NumberExp に評価された場合や予期しない Factor)
			// デフォルトの処理にフォールスルーします (ocode は default で設定)
			log.Printf("[pass1] Case 3b: Evaluated ImmExp is not a simple label for %s: %s. Falling through to default.", instName, op.TokenLiteral())
			// default ケースで処理するため、ここでは ocode と estimatedSize を設定しません
		}

	// MemoryAddrExp の処理が必要な場合 (例: JMP DWORD PTR [EAX])
	// TODO: JMP/CALL [memory] の処理を追加

	// デフォルトケース: AddExp, MultExp, ラベル以外の ImmExp, MemoryAddrExp, その他の未解決の式、または予期しない型
	default:
		// ImmExp のケース 3b からフォールスルーしてきた場合もここで処理
		// default ケースに来た場合、ocode と estimatedSize が未設定の可能性があるため、ここで設定する
		if ocode == "" { // ImmExp ケース 3b から来た場合など、まだ設定されていない場合
			if _, ok := op.(*ast.ImmExp); ok {
				// ImmExp のログは上で出力済みなので、ここでは一般的なログを出力
				log.Printf("[pass1] Case 4 (from ImmExp): Processing non-label ImmExp %s for %s", op.TokenLiteral(), instName)
			} else {
				log.Printf("[pass1] Case 4: Processing expression requiring Pass 2 resolution for %s: %s (Type: %T)", instName, evaluatedOperand.TokenLiteral(), evaluatedOperand)
			}
			estimatedSize = estimateJumpSize(instName, env.BitMode)
			// 評価されたオペランドの文字列表現をそのまま使用します (プレースホルダーなし)
			ocode = fmt.Sprintf("%s %s", instName, evaluatedOperand.TokenLiteral())
		}
		// すでに他のケース (例: ImmExp 3a) で ocode が設定されている場合は、ここでは何もしません
	}

	// LOC を更新し、Ocode を発行します
	if estimatedSize == 0 {
		log.Printf("WARN: Estimated size is 0 for %s %s. Defaulting to 2.", instName, operand.TokenLiteral())
		estimatedSize = 2 // LOC が進まないのを避けます
	}
	env.LOC += estimatedSize
	env.Client.Emit(ocode)
}
