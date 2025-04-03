package codegen

import (
	"fmt"
	"strconv"
	"strings"

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

	// Add prefixes if required
	if opsInterface.Require66h() { // Check for operand size prefix
		code = append(code, 0x66)
	}
	// Add address size prefix if needed (less common for IN)
	// if opsInterface.Require67h() {
	// 	code = append(code, 0x67)
	// }

	code = append(code, opcodeByte) // Append opcode
	if hasImm {
		code = append(code, byte(immValue)) // Append immediate if present
	}

	return code, nil
}
