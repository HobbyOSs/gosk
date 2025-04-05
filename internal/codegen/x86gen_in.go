package codegen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/HobbyOSs/gosk/pkg/cpu"        // Import cpu package
	"github.com/HobbyOSs/gosk/pkg/ng_operand" // Add ng_operand import
)

// handleIN generates x86 machine code for the IN instruction.
func handleIN(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	if len(params.Operands) != 2 {
		return nil, fmt.Errorf("IN instruction requires 2 operands, but got %d", len(params.Operands))
	}

	// Parse operands to check for prefixes
	var opsInterface ng_operand.Operands                                          // Declare opsInterface
	var err error                                                                 // Declare err once
	opsInterface, err = ng_operand.FromString(strings.Join(params.Operands, ",")) // Use = for assignment
	if err != nil {
		return nil, fmt.Errorf("failed to parse IN operands '%s': %v", strings.Join(params.Operands, ","), err)
	}
	opsInterface = opsInterface.WithBitMode(ctx.BitMode) // Set bit mode

	dst := strings.ToUpper(params.Operands[0]) // Destination: AL, AX, EAX
	src := params.Operands[1]                  // Source: imm8 or DX

	var opcodeByte byte
	var immValue uint64
	var hasImm bool
	// var err error // Already declared above

	// Check source type
	if strings.ToUpper(src) == "DX" {
		// Source is DX register
		switch dst {
		case "AL":
			opcodeByte = 0xEC
		case "AX", "EAX": // Add EAX case
			opcodeByte = 0xED // Same opcode for AX and EAX
		default:
			return nil, fmt.Errorf("invalid destination register '%s' for IN DX", dst)
		}
	} else {
		// Source is immediate (imm8)
		hasImm = true
		immValue, err = strconv.ParseUint(src, 0, 8) // Use = for assignment
		if err != nil {
			return nil, fmt.Errorf("invalid immediate value '%s' for IN instruction: %v", src, err)
		}

		switch dst {
		case "AL":
			opcodeByte = 0xE4
		case "AX", "EAX": // Add EAX case
			opcodeByte = 0xE5 // Same opcode for AX and EAX
		default:
			return nil, fmt.Errorf("invalid destination register '%s' for IN imm8", dst)
		}
	}

	// Assemble the code
	code := []byte{}

	// --- Prefix Calculation (similar to asmdb.GetPrefixSize for IN) ---
	prefix66Size := 0
	bitMode := ctx.BitMode
	opTypes := opsInterface.OperandTypes() // Get operand types from parsed interface

	if len(opTypes) > 0 {
		sizeDeterminingOpType := opTypes[0]                                                                     // For IN, the destination (AL/AX/EAX) determines size
		is16bitOp := sizeDeterminingOpType == ng_operand.CodeAX || sizeDeterminingOpType == ng_operand.CodeR16  // Consider r16 as well
		is32bitOp := sizeDeterminingOpType == ng_operand.CodeEAX || sizeDeterminingOpType == ng_operand.CodeR32 // Consider r32 as well

		if bitMode == cpu.MODE_16BIT && is32bitOp { // 16bit mode with 32bit operand (EAX)
			prefix66Size = 1
		} else if bitMode == cpu.MODE_32BIT && is16bitOp { // 32bit mode with 16bit operand (AX)
			prefix66Size = 1
		}
	}

	if prefix66Size > 0 {
		code = append(code, 0x66)
	}
	// Address size prefix (0x67) is generally not needed for IN with DX or imm8
	// if opsInterface.Require67h() {
	// 	code = append(code, 0x67)
	// }

	code = append(code, opcodeByte) // Append opcode
	if hasImm {
		code = append(code, byte(immValue)) // Append immediate if present
	}

	return code, nil
}
