package pass1

import (
	"fmt" // Keep only one fmt import
	"strings"

	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/HobbyOSs/gosk/pkg/cpu"        // Import cpu package
	"github.com/HobbyOSs/gosk/pkg/ng_operand" // Use ng_operand
	"github.com/samber/lo"
)

func processLGDT(env *Pass1, tokens []*token.ParseToken) {

	// オペランドの解析
	if len(tokens) != 1 {
		fmt.Printf("LGDT instruction expects exactly one operand")
		return
	}

	args := lo.Map(tokens, func(token *token.ParseToken, _ int) string {
		return token.AsString()
	})

	// Use ng_operand.FromString factory function
	operands, err := ng_operand.FromString(strings.Join(args, ","))
	if err != nil {
		// TODO: より適切なエラーハンドリングを行う
		fmt.Printf("Error creating operands from string in LGDT: %v\n", err)
		return // エラーが発生したら処理を中断
	}

	// Set BitMode
	operands = operands.WithBitMode(env.BitMode)

	// Check if the operand is a memory type using OperandTypes()
	opTypes := operands.OperandTypes()
	// Check if the first operand type is one of the memory types
	isMem := len(opTypes) == 1 && (opTypes[0] == ng_operand.CodeM ||
		opTypes[0] == ng_operand.CodeM8 ||
		opTypes[0] == ng_operand.CodeM16 ||
		opTypes[0] == ng_operand.CodeM32 ||
		opTypes[0] == ng_operand.CodeM64 ||
		opTypes[0] == ng_operand.CodeMEM) // Include CodeMEM
	if !isMem {
		fmt.Printf("LGDT instruction expects a memory operand, got %v (raw: %s)\n", opTypes, operands.InternalString())
		return
	}

	// Calculate LGDT size based on BitMode
	var lgdtSize int32
	switch env.BitMode {
	case cpu.MODE_16BIT:
		lgdtSize = 2 + 1 + 2 // Opcode + ModRM + disp16
	case cpu.MODE_32BIT:
		lgdtSize = 2 + 1 + 4 // Opcode + ModRM + disp32
	default:
		fmt.Printf("Error: Unsupported bit mode %v for LGDT size calculation\n", env.BitMode)
		return
	}
	env.LOC += lgdtSize

	env.Client.Emit(fmt.Sprintf("LGDT %s\n", strings.Join(args, ",")))
}
