package pass1

import (
	"fmt"
	"log"

	"github.com/HobbyOSs/gosk/internal/ast" // astパッケージをインポート
	"github.com/HobbyOSs/gosk/pkg/ng_operand" // Import ng_operand
	// "github.com/HobbyOSs/gosk/internal/token" // Remove unused token import
	// "github.com/HobbyOSs/gosk/pkg/cpu" // Remove duplicate and unused cpu import
)

// evalSimpleExp evaluates an expression node and returns the result as int32.
// It now accepts ast.Exp instead of ast.Node.
func evalSimpleExp(exp ast.Exp, env *Pass1) (int32, error) {
	// Evaluate the expression node using TraverseAST
	evalNode := TraverseAST(exp, env) // TraverseAST now returns ast.Exp

	// Check if the result is a number (NumberExp)
	if numExp, ok := evalNode.(*ast.NumberExp); ok {
		return int32(numExp.Value), nil
	}

	// 評価結果が未解決の識別子 (ImmExp with IdentFactor) かどうかを確認
	if immExp, ok := evalNode.(*ast.ImmExp); ok {
		if identFactor, ok := immExp.Factor.(*ast.IdentFactor); ok {
			// TODO: ラベルやEQUの解決 (Pass2で行うか、ここでシンボルテーブルを参照するか)
			// 現時点では未解決としてエラーを返すか、0を返す
			log.Printf("WARN: Identifier '%s' evaluation in evalSimpleExp is not fully implemented yet (returning 0).", identFactor.Value)
			// return 0, fmt.Errorf("identifier evaluation not implemented: %s", identFactor.Value)
			return 0, nil // 仮に0を返す
		}
	}

	// その他の評価不能なケース
	return 0, fmt.Errorf("cannot evaluate expression result to int32: %v (Type: %T)", evalNode.TokenLiteral(), evalNode)
}

// processCalcJcc handles JMP and conditional jump instructions.
func processCalcJcc(env *Pass1, operands []ast.Exp, instName string) {
	if len(operands) != 1 {
		log.Printf("Error: %s instruction requires exactly one operand, got %d", instName, len(operands))
		return
	}

	operand := operands[0]

	switch op := operand.(type) {
	case *ast.SegmentExp: // Handle FAR jumps (e.g., JMP FAR label, JMP seg:off)
		log.Printf("[pass1] Processing SegmentExp for %s: %s", instName, op.TokenLiteral())

		// Evaluate segment and offset
		segment, errSeg := evalSimpleExp(op.Left, env) // Pass ast.Exp
		if errSeg != nil {
			log.Printf("Error evaluating segment expression for %s: %v", instName, errSeg)
			return
		}
		if op.Right == nil {
			log.Printf("Error: SegmentExp without Right part is not supported for %s FAR.", instName)
			return
		}
		offset, errOff := evalSimpleExp(op.Right, env) // Pass ast.Exp
		if errOff != nil {
			log.Printf("Error evaluating offset expression for %s: %v", instName, errOff)
			return
		}

		// Calculate size (JMP ptr16:16/32)
		size := int32(7) // EA + ptr16:32 (4 byte offset + 2 byte selector)
		// if env.BitMode == cpu.MODE_16BIT { size = ? } // Adjust for 16-bit ptr16:16 if needed
		env.LOC += size

		// Emit Ocode (placeholder format)
		env.Client.Emit(fmt.Sprintf("%s_FAR %d:%d ; (size: %d)", instName, segment, offset, size))

	case *ast.ImmExp: // Handle labels (IdentFactor)
		if factor, ok := op.Factor.(*ast.IdentFactor); ok {
			label := factor.Value
			// Register label in SymTable (placeholder address)
			if _, exists := env.SymTable[label]; !exists {
				env.SymTable[label] = 0 // Placeholder for Pass 1
			}
			// Calculate size using FindMinOutputSize
			operandString := op.TokenLiteral()
			ngOperands, err := ng_operand.FromString(operandString)
			if err != nil {
				log.Printf("Error creating operand from string '%s' in %s: %v", operandString, instName, err)
				return
			}
			ngOperands = ngOperands.WithBitMode(env.BitMode)
			// Jcc might need relative address handling forced differently than CALL?
			// ngOperands = ngOperands.WithForceRelAsImm(false) // Example if needed

			size, err := env.AsmDB.FindMinOutputSize(instName, ngOperands)
			if err != nil {
				log.Printf("Error finding size for %s %s: %v. Assuming default size 2.", instName, operandString, err)
				size = 2 // Default size (rel8) on error
			}
			env.LOC += int32(size)

			// Emit Ocode with label placeholder (no comment)
			env.Client.Emit(fmt.Sprintf("%s {{.%s}}", instName, label))
		} else {
			log.Printf("Error: Invalid factor type %T within ImmExp for %s operand.", op.Factor, instName)
		}

	case *ast.NumberExp: // Handle immediate address (relative jump target)
		targetAddr := op.Value // int64
		// Calculate relative offset (this is tricky in Pass 1 as LOC changes)
		// Assume the offset calculation happens relative to the *end* of the current instruction.
		// We need the size first. Assume rel8/rel16/rel32 based on offset range.
		// This calculation is inherently problematic in a single pass without fixups.
		// Let's estimate size based on typical jump instructions.
		// JMP rel8 (EB cb) = 2 bytes
		// JMP rel16/32 (E9 cw/cd) = 3/5 bytes
		// Jcc rel8 (7x cb) = 2 bytes
		// Jcc rel16/32 (0F 8x cw/cd) = 4/6 bytes

		// Calculate size using FindMinOutputSize
		operandString := op.TokenLiteral()
		ngOperands, err := ng_operand.FromString(operandString)
		if err != nil {
			log.Printf("Error creating operand from string '%s' in %s: %v", operandString, instName, err)
			return
		}
		ngOperands = ngOperands.WithBitMode(env.BitMode)

		size, err := env.AsmDB.FindMinOutputSize(instName, ngOperands)
		if err != nil {
			log.Printf("Error finding size for %s %s: %v. Assuming default size 2.", instName, operandString, err)
			size = 2 // Default size (rel8) on error
		}
		env.LOC += int32(size)

		// Create dummy label for immediate jump
		fakeLabel := fmt.Sprintf("imm_jmp_%d", env.NextImmJumpID)
		env.NextImmJumpID++
		env.SymTable[fakeLabel] = int32(targetAddr) // Store absolute target address

		// Emit Ocode with dummy label placeholder (no comment)
		log.Printf("[pass1] %s immediate value: %d (0x%x)", instName, targetAddr, targetAddr)
		env.Client.Emit(fmt.Sprintf("%s {{.%s}}", instName, fakeLabel))

	// TODO: Handle other operand types like MemoryAddrExp (e.g., JMP DWORD PTR [EAX]) if necessary
	default:
		log.Printf("Error: Invalid operand type %T for %s instruction.", operand, instName)
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
