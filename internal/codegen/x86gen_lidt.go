package codegen

import (
	"fmt"
	"log"

	"github.com/HobbyOSs/gosk/pkg/asmdb"
	"github.com/HobbyOSs/gosk/pkg/ng_operand"
)

// handleLIDT handles the LIDT instruction and generates the appropriate machine code.
// LIDT m -> 0F 01 /3
func handleLIDT(operands []string, ctx *CodeGenContext) ([]byte, error) {
	instName := "LIDT" // Change instruction name
	if len(operands) != 1 {
		return nil, fmt.Errorf("%s requires 1 operand, got %d", instName, len(operands))
	}

	// Deserialize operand string from pass1 using FromString
	op, err := ng_operand.FromString(operands[0]) // Use FromString
	if err != nil {
		// エラーハンドリング改善: FromString が失敗した場合のエラーメッセージ
		return nil, fmt.Errorf("failed to parse operand string '%s' for %s: %w", operands[0], instName, err)
	}
	op = op.WithBitMode(ctx.BitMode) // Ensure bit mode is set

	// Ensure it's a memory operand by checking OperandTypes
	opTypes := op.OperandTypes()
	isMem := len(opTypes) == 1 && (opTypes[0] == ng_operand.CodeM ||
		opTypes[0] == ng_operand.CodeM8 ||
		opTypes[0] == ng_operand.CodeM16 ||
		opTypes[0] == ng_operand.CodeM32 ||
		opTypes[0] == ng_operand.CodeM64 ||
		opTypes[0] == ng_operand.CodeMEM)
	if !isMem {
		// Use Serialize() for error message
		return nil, fmt.Errorf("%s requires a memory operand, got %s", instName, op.Serialize())
	}

	// Create a local InstructionDB instance and find the encoding, similar to x86gen_mov.go
	db := asmdb.NewInstructionDB()
	// Use the Operands interface directly. matchAnyImm=false for codegen.
	encoding, err := db.FindEncoding(instName, op, false)
	if err != nil {
		return nil, fmt.Errorf("failed to find encoding for %s with operand %s: %w", instName, op.Serialize(), err)
	}
	if encoding == nil { // Check if encoding is nil after FindEncoding
		return nil, fmt.Errorf("could not find valid encoding form for %s with operand %s", instName, op.Serialize())
	}

	// Generate base opcode bytes (0F 01)
	opcodeBytes, err := ResolveOpcode(encoding.Opcode, -1) // -1 indicates no register addend
	if err != nil {
		return nil, fmt.Errorf("failed to resolve opcode for %s: %w", instName, err)
	}

	// Generate ModR/M, SIB, Displacement using the utility function
	// Pass the original serialized operand string from pass1 to GenerateModRM
	// Pass the encoding found via FindEncoding
	modrmSibDispBytes, err := GenerateModRM([]string{operands[0]}, encoding, ctx.BitMode) // Pass original string and found encoding
	if err != nil {
		// エラーハンドリング改善: GenerateModRM が失敗した場合のエラーメッセージ
		return nil, fmt.Errorf("failed to generate ModR/M for %s operand %s: %w", instName, operands[0], err) // Use original string in error
	}

	// Combine opcode and ModR/M bytes
	machineCode := append(opcodeBytes, modrmSibDispBytes...)

	// ログメッセージ改善: 元のオペランド文字列も表示
	log.Printf("debug: Generated %s machine code: %x (operand: %s)", instName, machineCode, operands[0]) // Change log message

	return machineCode, nil
}
