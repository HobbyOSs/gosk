package pass1

import (
	"fmt"
	"log" // Add log import

	"github.com/HobbyOSs/gosk/internal/ast"   // Add ast import
	"github.com/HobbyOSs/gosk/pkg/ng_operand" // Use ng_operand
)

// processLIDT handles the LIDT instruction.
// LIDT m16&32
// オペランドはメモリアドレスのみ
func processLIDT(env *Pass1, operands []ast.Exp) {
	instName := "LIDT"
	if len(operands) != 1 {
		log.Printf("Error: %s instruction requires exactly one operand.", instName)
		return
	}

	// Get string representation of the operand
	operandString := operands[0].TokenLiteral()

	// Create ng_operand.Operands from the string
	ngOperands, err := ng_operand.FromString(operandString)
	if err != nil {
		log.Printf("Error creating operand from string '%s' in %s: %v", operandString, instName, err)
		return
	}

	// Set BitMode
	ngOperands = ngOperands.WithBitMode(env.BitMode)

	// Check if the operand is a memory type
	opTypes := ngOperands.OperandTypes()
	// Manually check if the type is one of the memory types
	isMem := len(opTypes) == 1 && (opTypes[0] == ng_operand.CodeM ||
		opTypes[0] == ng_operand.CodeM8 ||
		opTypes[0] == ng_operand.CodeM16 ||
		opTypes[0] == ng_operand.CodeM32 ||
		opTypes[0] == ng_operand.CodeM64 ||
		opTypes[0] == ng_operand.CodeMEM)
	if !isMem {
		log.Printf("Error: %s instruction expects a memory operand, got %v (raw: %s)", instName, opTypes, operandString)
		return
	}

	// Calculate size using FindMinOutputSize (consistent with pass1_inst_mov.go)
	size, err := env.AsmDB.FindMinOutputSize(instName, ngOperands) // Use env.AsmDB
	if err != nil {
		log.Printf("Error finding min output size for %s %s: %v", instName, ngOperands.Serialize(), err)
		return
	}
	env.LOC += int32(size) // Update LOC with the calculated size

	// Emit the command
	env.Client.Emit(fmt.Sprintf("%s %s", instName, ngOperands.Serialize()))
}
