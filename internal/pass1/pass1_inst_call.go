package pass1

import (
	"fmt"
	"log"

	"github.com/HobbyOSs/gosk/internal/ast"   // Add ast import
	"github.com/HobbyOSs/gosk/pkg/ng_operand" // Import ng_operand
)

// processCALL handles the CALL instruction.
func processCALL(env *Pass1, operands []ast.Exp) {
	const instName = "CALL"
	if len(operands) != 1 {
		log.Printf("Error: %s instruction requires exactly one operand, got %d", instName, len(operands))
		return
	}

	operand := operands[0]

	switch op := operand.(type) {
	case *ast.ImmExp: // Handles labels (IdentFactor)
		if factor, ok := op.Factor.(*ast.IdentFactor); ok {
			label := factor.Value

			// ラベルをSymTableに登録 (仮アドレスを割り当てる)
			if _, exists := env.SymTable[label]; !exists {
				env.SymTable[label] = 0 // Pass 1では仮アドレス
			}
			// Pass1では正確なサイズ計算は難しい。FindMinOutputSize を試みる。
			// 文字列表現から ng_operand を作成
			operandString := op.TokenLiteral()
			ngOperands, err := ng_operand.FromString(operandString)
			if err != nil {
				log.Printf("Error creating operand from string '%s' in %s: %v", operandString, instName, err)
				return
			}
			ngOperands = ngOperands.WithBitMode(env.BitMode)

			size, err := env.AsmDB.FindMinOutputSize(instName, ngOperands)
			if err != nil {
				log.Printf("Error finding size for %s %s: %v. Assuming default size 3.", instName, operandString, err)
				size = 3 // Default size on error
			}
			env.LOC += int32(size)

			// Ocodeを生成 (ジャンプ先アドレスはプレースホルダー)
			env.Client.Emit(fmt.Sprintf("%s {{.%s}} ; (size: %d)", instName, label, size))
		} else {
			log.Printf("Error: Invalid factor type %T within ImmExp for %s operand.", op.Factor, instName)
		}
	case *ast.NumberExp: // Handles immediate address
		targetAddr := op.Value // int64

		// 機械語サイズを計算 (Opcode + Offset)
		// TODO: $ を考慮したオフセット計算が必要な場合、Eval で処理されているべき
		//       ここでは単純なターゲットアドレスからのオフセットと仮定
		// offset := targetAddr - int64(env.LOC) // Incorrect: LOC is updated *after* size calculation
		// Pass1 では正確なオフセット計算は難しい。FindMinOutputSize を試みる。
		operandString := op.TokenLiteral()
		ngOperands, err := ng_operand.FromString(operandString)
		if err != nil {
			log.Printf("Error creating operand from string '%s' in %s: %v", operandString, instName, err)
			return
		}
		ngOperands = ngOperands.WithBitMode(env.BitMode)

		size, err := env.AsmDB.FindMinOutputSize(instName, ngOperands)
		if err != nil {
			log.Printf("Error finding size for %s %s: %v. Assuming default size 3.", instName, operandString, err)
			size = 3 // Default size on error
		}
		env.LOC += int32(size)

		// ダミーのラベルを作る
		fakeLabel := fmt.Sprintf("imm_call_%d", env.NextImmJumpID)
		env.NextImmJumpID++
		env.SymTable[fakeLabel] = int32(targetAddr) // Store target address

		// Ocodeを生成 (ジャンプ先アドレスはダミー)
		env.Client.Emit(fmt.Sprintf("%s {{.%s}} ; (size: %d)", instName, fakeLabel, size))
	// TODO: Handle other operand types like MemoryAddrExp (e.g., CALL DWORD PTR [EAX]) if necessary
	default:
		log.Printf("Error: Invalid operand type %T for %s instruction.", operand, instName)
	}
}
