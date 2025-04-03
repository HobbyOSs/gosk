package codegen

import (
	"fmt"
	"strconv"
	"strings"
)

// handleIN generates x86 machine code for the IN instruction.
func handleIN(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	if len(params.Operands) != 2 {
		return nil, fmt.Errorf("IN instruction requires 2 operands, but got %d", len(params.Operands))
	}

	dst := strings.ToUpper(params.Operands[0]) // Destination: AL or AX
	src := params.Operands[1]                  // Source: imm8 or DX

	var opcodeByte byte
	var immValue uint64
	var hasImm bool
	var err error

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
		immValue, err = strconv.ParseUint(src, 0, 8) // Parse imm8
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
	code := []byte{opcodeByte}
	if hasImm {
		code = append(code, byte(immValue))
	}

	return code, nil
}
