package codegen

import (
	"encoding/binary"
	"fmt"
	"log"
	"strings"

	"github.com/HobbyOSs/gosk/pkg/cpu" // Import cpu package for BitMode
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

	var modrmByte byte
	var dispBytes []byte

	// Determine ModR/M and displacement based on BitMode
	switch ctx.BitMode {
	case cpu.MODE_16BIT:
		// ModR/M: mod=00, reg=2 (010), r/m=6 (110) -> 00010110 -> 0x16 for [disp16]
		modrmByte = 0x16
		// Displacement: 16-bit address of the label
		dispBytes = make([]byte, 2)
		binary.LittleEndian.PutUint16(dispBytes, uint16(addr))
	case cpu.MODE_32BIT:
		// ModR/M: mod=00, reg=2 (010), r/m=5 (101) -> 00010101 -> 0x15 for [disp32]
		modrmByte = 0x15
		// Displacement: 32-bit address of the label
		dispBytes = make([]byte, 4)
		binary.LittleEndian.PutUint32(dispBytes, uint32(addr))
	default:
		return nil, fmt.Errorf("unsupported bit mode for LGDT: %v", ctx.BitMode)
	}

	machineCode = append(machineCode, modrmByte)
	machineCode = append(machineCode, dispBytes...)

	log.Printf("debug: Generated LGDT machine code (%s): % x", ctx.BitMode, machineCode)

	return machineCode, nil
}
