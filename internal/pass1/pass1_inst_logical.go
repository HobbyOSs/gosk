package pass1

import (
	"fmt" // Keep only one fmt import
	"log" // Keep only one log import
	"strings"

	"github.com/HobbyOSs/gosk/internal/ast"   // Add ast import
	"github.com/HobbyOSs/gosk/pkg/ng_operand" // Use ng_operand
	"github.com/samber/lo"
)

// processLogicalInst handles common logic for logical instructions.
func processLogicalInst(env *Pass1, operands []ast.Exp, instName string) {
	// Get string representation of operands
	operandStrings := lo.Map(operands, func(exp ast.Exp, _ int) string {
		return exp.TokenLiteral() // Assuming TokenLiteral is suitable
	})
	operandString := strings.Join(operandStrings, ",")

	// isAccumulator := false
	// if len(args) > 0 {
	// 	matched, _ := regexp.MatchString(`(?i)^(AL|AX|EAX|RAX)$`, args[0])
	// 	isAccumulator = matched
	// }

	// Create ng_operand.Operands from the combined string
	ngOperands, err := ng_operand.FromString(operandString)
	if err != nil {
		log.Printf("Error creating operands from string '%s' in %s: %v", operandString, instName, err)
		return
	}

	// Set BitMode
	ngOperands = ngOperands.WithBitMode(env.BitMode)

	// Calculate instruction size
	size, err := env.AsmDB.FindMinOutputSize(instName, ngOperands)
	if err != nil {
		log.Printf("Error finding min output size for %s %s: %v", instName, operandString, err)
		return
	}
	env.LOC += int32(size)

	// Emit the command
	env.Client.Emit(fmt.Sprintf("%s %s", instName, ngOperands.Serialize()))
}

// --- Update callers to use the new signature ---

// AND命令
func processAND(env *Pass1, operands []ast.Exp) {
	processLogicalInst(env, operands, "AND")
}

// OR命令
func processOR(env *Pass1, operands []ast.Exp) {
	processLogicalInst(env, operands, "OR")
}

// XOR命令
func processXOR(env *Pass1, operands []ast.Exp) {
	processLogicalInst(env, operands, "XOR")
}

// NOT命令
func processNOT(env *Pass1, operands []ast.Exp) {
	if len(operands) != 1 {
		log.Printf("Error: NOT instruction requires exactly one operand.")
		return
	}

	// Get string representation of the operand
	operandString := operands[0].TokenLiteral()

	// Create ng_operand.Operands from the string
	ngOperands, err := ng_operand.FromString(operandString)
	if err != nil {
		log.Printf("Error creating operand from string '%s' in NOT: %v", operandString, err)
		return
	}

	// Set BitMode
	ngOperands = ngOperands.WithBitMode(env.BitMode)

	// Calculate instruction size
	size, err := env.AsmDB.FindMinOutputSize("NOT", ngOperands)
	if err != nil {
		log.Printf("Error finding min output size for NOT %s: %v", operandString, err)
		return
	}
	env.LOC += int32(size)

	// Emit the command
	env.Client.Emit(fmt.Sprintf("NOT %s", ngOperands.Serialize()))
}

// SHR命令
func processSHR(env *Pass1, operands []ast.Exp) {
	processLogicalInst(env, operands, "SHR")
}

// SHL命令
func processSHL(env *Pass1, operands []ast.Exp) {
	processLogicalInst(env, operands, "SHL")
}

// SAR命令
func processSAR(env *Pass1, operands []ast.Exp) {
	processLogicalInst(env, operands, "SAR")
}
