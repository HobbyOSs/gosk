package pass1

import (
	"fmt" // Keep only one fmt import
	"log" // Add log import
	"strings"

	"github.com/HobbyOSs/gosk/internal/ast" // Add ast import
	"github.com/HobbyOSs/gosk/pkg/cpu"
	"github.com/HobbyOSs/gosk/pkg/ng_operand" // Use ng_operand
	"github.com/samber/lo"
)

// processArithmeticInst は算術命令の共通ロジックを処理します。
func processArithmeticInst(env *Pass1, operands []ast.Exp, instName string) {
	// オペランドの文字列表現を取得します
	operandStrings := lo.Map(operands, func(exp ast.Exp, _ int) string {
		return exp.TokenLiteral() // TokenLiteral が適切であると仮定します
	})
	operandString := strings.Join(operandStrings, ",")

	// 結合された文字列から ng_operand.Operands を作成します
	ngOperands, err := ng_operand.FromString(operandString)
	if err != nil {
		log.Printf("Error creating operands from string '%s' in %s: %v", operandString, instName, err)
		return
	}

	// BitMode を設定します
	ngOperands = ngOperands.WithBitMode(env.BitMode)

	// 命令サイズを計算します
	size, err := env.AsmDB.FindMinOutputSize(instName, ngOperands)
	if err != nil {
		log.Printf("Error finding min output size for %s %s: %v", instName, operandString, err)
		return
	}
	env.LOC += int32(size)

	// コマンドを発行します (コメントなし)
	env.Client.Emit(fmt.Sprintf("%s %s", instName, ngOperands.Serialize()))
}

// --- 呼び出し元を新しいシグネチャを使用するように更新 ---

// 加算命令
func processADD(env *Pass1, operands []ast.Exp) {
	processArithmeticInst(env, operands, "ADD")
}

// キャリー付き加算命令
func processADC(env *Pass1, operands []ast.Exp) {
	processArithmeticInst(env, operands, "ADC")
}

// 減算命令
func processSUB(env *Pass1, operands []ast.Exp) {
	processArithmeticInst(env, operands, "SUB")
}

// ボロー付き減算命令
func processSBB(env *Pass1, operands []ast.Exp) {
	processArithmeticInst(env, operands, "SBB")
}

// 比較命令
func processCMP(env *Pass1, operands []ast.Exp) {
	processArithmeticInst(env, operands, "CMP")
}

// インクリメント命令
func processINC(env *Pass1, operands []ast.Exp) {
	processArithmeticInst(env, operands, "INC")
}

// デクリメント命令
func processDEC(env *Pass1, operands []ast.Exp) {
	processArithmeticInst(env, operands, "DEC")
}

// 2の補数命令
func processNEG(env *Pass1, operands []ast.Exp) {
	processArithmeticInst(env, operands, "NEG")
}

// 符号なし乗算命令
func processMUL(env *Pass1, operands []ast.Exp) {
	processArithmeticInst(env, operands, "MUL")
}

// 符号付き乗算命令
func processIMUL(env *Pass1, operands []ast.Exp) {
	instName := "IMUL"
	// オペランドの文字列表現を取得します
	operandStrings := lo.Map(operands, func(exp ast.Exp, _ int) string {
		return exp.TokenLiteral()
	})
	operandString := strings.Join(operandStrings, ",")

	// 結合された文字列から ng_operand.Operands を作成します
	ngOperands, err := ng_operand.FromString(operandString)
	if err != nil {
		log.Printf("Error creating operands from string '%s' in %s: %v", operandString, instName, err)
		return
	}

	// BitMode を設定します
	ngOperands = ngOperands.WithBitMode(env.BitMode)

	// FindMinOutputSize を使用して命令サイズを計算します
	calculatedSize, err := env.AsmDB.FindMinOutputSize(instName, ngOperands)
	if err != nil {
		log.Printf("Error finding min output size for %s %s: %v", instName, operandString, err)
		return
	}

	var size int = calculatedSize // デフォルトで計算されたサイズを使用します

	// ★★★ IMUL r32, imm16/32 (16bit mode) のサイズ上書き処理 ★★★
	// FindMinOutputSize が特定のケースで不正なサイズを返す問題への暫定対応。
	// TODO: FindMinOutputSize または依存関係 (asmdb, ng_operand) を修正し、この暫定対応を削除する。
	if env.BitMode == cpu.MODE_16BIT && len(ngOperands.OperandTypes()) == 2 {
		// オペランドが r32 と imm16/imm32 の組み合わせかチェック
		isR32 := ngOperands.IsType(0, ng_operand.CodeR32)
		isImm16 := ngOperands.IsType(1, ng_operand.CodeIMM16)
		isImm32 := ngOperands.IsType(1, ng_operand.CodeIMM32)

		if isR32 && (isImm16 || isImm32) {
			expectedSize := 7 // 16bitモードでの IMUL r32, imm16/32 の期待サイズ
			if calculatedSize < expectedSize {
				// FindMinOutputSize が期待より小さい不正な値を返した場合のみ上書き
				log.Printf("debug: [pass1] Overriding calculated size %d with %d for IMUL r32, imm16/32 (%s) in 16-bit mode.\n",
					calculatedSize, expectedSize, operandString)
				size = expectedSize
			}
			// calculatedSize >= expectedSize の場合は、FindMinOutputSize が正しいか、
			// より大きいサイズを返した可能性があるので、そのまま使う。
		}
	}

	// LOC を加算
	env.LOC += int32(size)

	// コマンドを発行します (コメントなし)
	env.Client.Emit(fmt.Sprintf("%s %s", instName, ngOperands.Serialize()))
}

// 符号なし除算命令
func processDIV(env *Pass1, operands []ast.Exp) {
	processArithmeticInst(env, operands, "DIV")
}

// 符号付き除算命令
func processIDIV(env *Pass1, operands []ast.Exp) {
	processArithmeticInst(env, operands, "IDIV")
}
