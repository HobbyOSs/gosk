package pass1

import (
	"fmt"
	"log"

	"github.com/HobbyOSs/gosk/internal/ast" // astパッケージをインポート
	"github.com/HobbyOSs/gosk/pkg/cpu"
)

// estimateJumpSize estimates the size of a near jump/call instruction in Pass 1.
// This is an estimation because the final offset size (rel8/rel16/32) might change in Pass 2.
func estimateJumpSize(instName string, bitMode cpu.BitMode) int32 {
	isJcc := instName != "JMP" && instName != "CALL" // Assume others are Jcc

	// Default to near relative jump/call sizes (opcode + rel16/32)
	// JMP rel16/32 (E9 cw/cd): 1 + 2/4 = 3/5 bytes
	// CALL rel16/32 (E8 cw/cd): 1 + 2/4 = 3/5 bytes
	// Jcc rel16/32 (0F 8x cw/cd): 2 + 2/4 = 4/6 bytes
	size := int32(5) // Assume rel32 for JMP/CALL initially
	if bitMode == cpu.MODE_16BIT {
		size = 3 // Assume rel16 for JMP/CALL in 16-bit mode
	}

	if isJcc {
		size = 6 // Assume rel32 for Jcc initially
		if bitMode == cpu.MODE_16BIT {
			size = 4 // Assume rel16 for Jcc in 16-bit mode
		}
	}
	// Note: We don't estimate short jumps (rel8) here, Pass 2 will optimize if possible.
	return size
}

// processCalcJcc handles JMP and conditional jump instructions.
func processCalcJcc(env *Pass1, operands []ast.Exp, instName string) {
	if len(operands) != 1 {
		log.Printf("Error: %s instruction requires exactly one operand, got %d", instName, len(operands))
		return
	}

	operand := operands[0]
	evaluatedOperand, _ := operand.Eval(env) // Evaluate the operand first, explicitly ignore 'evaluated' flag

	// Determine estimated size and emit Ocode based on the *evaluated* operand type
	var estimatedSize int32
	var ocode string

	switch op := evaluatedOperand.(type) {
	case *ast.SegmentExp: // Handle FAR jumps (e.g., JMP FAR label, JMP seg:off)
		log.Printf("[pass1] Processing evaluated SegmentExp for %s: %s", instName, op.TokenLiteral())
		// Evaluate segment and offset parts *again* (Eval on SegmentExp itself might not fully resolve)
		segEval, segOk := op.Left.Eval(env)
		offEval, offOk := op.Right.Eval(env)

		if !segOk || !offOk {
			log.Printf("Error: Could not fully evaluate segment or offset for FAR %s.", instName)
			// Emit placeholder based on original operand string if evaluation failed
			estimatedSize = 7 // Assume ptr16:32
			ocode = fmt.Sprintf("%s {{expr:%s}}", instName, operand.TokenLiteral())
		} else {
			// Try to get values if they are numbers, otherwise use placeholders
			var segStr, offStr string
			if segNum, ok := segEval.(*ast.NumberExp); ok {
				segStr = fmt.Sprintf("%d", segNum.Value)
			} else {
				segStr = fmt.Sprintf("{{expr:%s}}", op.Left.TokenLiteral()) // Placeholder for segment
			}
			if offNum, ok := offEval.(*ast.NumberExp); ok {
				offStr = fmt.Sprintf("%d", offNum.Value)
			} else {
				offStr = fmt.Sprintf("{{expr:%s}}", op.Right.TokenLiteral()) // Placeholder for offset
			}

			estimatedSize = 7 // JMP ptr16:32 (EA + ptr16:32)
			ocode = fmt.Sprintf("%s_FAR %s:%s", instName, segStr, offStr)
		}

	case *ast.ImmExp:
		if factor, ok := op.Factor.(*ast.IdentFactor); ok { // Unresolved label
			label := factor.Value
			log.Printf("[pass1] Processing label '%s' for %s", label, instName)
			// Register label if not exists
			if _, exists := env.SymTable[label]; !exists {
				env.SymTable[label] = 0 // Placeholder address
			}
			estimatedSize = estimateJumpSize(instName, env.BitMode)
			ocode = fmt.Sprintf("%s {{.%s}}", instName, label) // Ocode with label placeholder
		} else {
			// Should not happen if ImmExp.Eval works correctly, but handle defensively
			log.Printf("Error: Unexpected factor type %T within evaluated ImmExp for %s.", op.Factor, instName)
			estimatedSize = estimateJumpSize(instName, env.BitMode)                 // Estimate size anyway
			ocode = fmt.Sprintf("%s {{expr:%s}}", instName, operand.TokenLiteral()) // Placeholder with original expr
		}

	case *ast.NumberExp: // Resolved immediate address
		targetAddr := op.Value
		log.Printf("[pass1] Processing immediate address %d (0x%x) for %s", targetAddr, targetAddr, instName)
		// Pass 1 cannot reliably calculate relative offset. Use a placeholder.
		estimatedSize = estimateJumpSize(instName, env.BitMode)
		// Use a placeholder indicating an immediate address for Pass 2
		ocode = fmt.Sprintf("%s {{addr:%d}}", instName, targetAddr)

	case *ast.AddExp, *ast.MultExp: // Partially evaluated expression (e.g., label + offset)
		log.Printf("[pass1] Processing partially evaluated expression for %s: %s", instName, op.TokenLiteral())
		// Cannot fully resolve in Pass 1. Use a placeholder for Pass 2.
		estimatedSize = estimateJumpSize(instName, env.BitMode)
		ocode = fmt.Sprintf("%s {{expr:%s}}", instName, op.TokenLiteral()) // Placeholder with the expression string

	// TODO: Handle MemoryAddrExp if needed (e.g., JMP DWORD PTR [EAX])

	default:
		log.Printf("Error: Invalid evaluated operand type %T for %s instruction.", evaluatedOperand, instName)
		// Attempt to use original operand string as a fallback placeholder
		estimatedSize = estimateJumpSize(instName, env.BitMode) // Estimate size
		ocode = fmt.Sprintf("%s {{expr:%s}}", instName, operand.TokenLiteral())
	}

	// Update LOC and emit Ocode
	if estimatedSize == 0 {
		log.Printf("WARN: Estimated size is 0 for %s %s. Defaulting to 2.", instName, operand.TokenLiteral())
		estimatedSize = 2 // Avoid LOC not advancing
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
