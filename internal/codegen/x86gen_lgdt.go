package codegen

import (
	"encoding/binary"
	"fmt"
	"log"
	"strings"
)

// handleLGDT handles the LGDT instruction and generates the appropriate machine code.
// LGDT m -> 0F 01 /2
func handleLGDT(operands []string, ctx *CodeGenContext) ([]byte, error) {
	if len(operands) != 1 {
		return nil, fmt.Errorf("LGDT requires 1 operand, got %d", len(operands))
	}

	opStr := operands[0]
	// Expecting memory operand like "[ label ]"
	if !strings.HasPrefix(opStr, "[") || !strings.HasSuffix(opStr, "]") {
		return nil, fmt.Errorf("invalid LGDT operand format: %s", opStr)
	}
	label := strings.TrimSpace(opStr[1 : len(opStr)-1])

	// Lookup label address in symbol table
	addr, ok := ctx.SymTable[label]
	if !ok {
		return nil, fmt.Errorf("label not found: %s", label)
	}

	// Opcode: 0F 01
	machineCode := []byte{0x0F, 0x01}

	// ModR/M: mod=00, reg=2 (010), r/m=5 (101) -> 00010101 -> 0x15
	// This indicates disp32 addressing mode.
	machineCode = append(machineCode, 0x15)

	// Displacement: 32-bit address of the label
	dispBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(dispBytes, uint32(addr))
	machineCode = append(machineCode, dispBytes...)

	log.Printf("debug: Generated LGDT machine code: % x", machineCode)

	return machineCode, nil
}
