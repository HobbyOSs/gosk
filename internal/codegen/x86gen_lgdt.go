package codegen

import (
	"fmt"
	"log"

	"github.com/HobbyOSs/gosk/pkg/asmdb"
	"github.com/HobbyOSs/gosk/pkg/ng_operand"
)

// handleLGDT handles the LGDT instruction and generates the appropriate machine code.
// LGDT m -> 0F 01 /2
func handleLGDT(operands []string, ctx *CodeGenContext) ([]byte, error) {
	instName := "LGDT"
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
		// Use Serialize() for error message
		return nil, fmt.Errorf("%s requires a memory operand, got %s", instName, op.Serialize())
	}

	// Get all instructions using the exported function
	allInstructions := asmdb.X86Instructions()
	lgdtInst, ok := allInstructions[instName]
	if !ok || len(lgdtInst.Forms) == 0 || len(lgdtInst.Forms[0].Encodings) == 0 {
		// Fallback or error if LGDT definition is missing in asmdb (should not happen with fallback)
		return nil, fmt.Errorf("internal error: LGDT instruction definition not found via asmdb.X86Instructions()")
	}
	// Use the first encoding form defined in the fallback table
	encoding := lgdtInst.Forms[0].Encodings[0] // Get the base encoding (Opcode 0F01, ModRM Reg 2)

	// Generate base opcode bytes (0F 01) using the retrieved encoding info
	opcodeBytes, err := ResolveOpcode(encoding.Opcode, -1) // -1 indicates no register addend
	if err != nil {
		return nil, fmt.Errorf("failed to resolve opcode for %s: %w", instName, err)
	}

	// Generate ModR/M, SIB, Displacement using the utility function
	// Pass the original serialized operand string from pass1 to GenerateModRM
	// Pass the encoding retrieved directly from asmdb
	modrmSibDispBytes, err := GenerateModRM([]string{operands[0]}, &encoding, ctx.BitMode) // Pass the retrieved encoding
	if err != nil {
		// エラーハンドリング改善: GenerateModRM が失敗した場合のエラーメッセージ
		return nil, fmt.Errorf("failed to generate ModR/M for %s operand %s: %w", instName, operands[0], err) // Use original string in error
	}

	// Combine opcode and ModR/M bytes
	baseMachineCode := append(opcodeBytes, modrmSibDispBytes...)

	// Add address-size prefix (0x67) if needed
	prefixBytes := []byte{}
	// Check if address size prefix (0x67) is needed using the Require67h method
	if op.Require67h() { // Use the correct method from ng_operand
		prefixBytes = append(prefixBytes, 0x67)
		// Use %v for ctx.BitMode as it might not be a string
		log.Printf("debug: Added address-size prefix (0x67) for %s in %v mode (operand: %s)", instName, ctx.BitMode, op.Serialize())
	}

	machineCode := append(prefixBytes, baseMachineCode...)

	// ログメッセージ改善: 元のオペランド文字列も表示
	log.Printf("debug: Generated %s machine code: %x (operand: %s, prefixes: %x)", instName, machineCode, operands[0], prefixBytes)

	return machineCode, nil
}
