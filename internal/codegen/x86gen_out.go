package codegen

import (
	"fmt"
	"strconv"
	"strings" // Add strings import
)

func handleOUT(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	if len(params.Operands) != 2 {
		return nil, fmt.Errorf("OUT instruction requires 2 operands, but got %d", len(params.Operands))
	}

	dst := strings.ToUpper(params.Operands[0]) // Destination: imm8 or DX
	src := strings.ToUpper(params.Operands[1]) // Source: AL, AX, or EAX

	var opcodeByte byte
	var immValue uint64
	var hasImm bool
	var err error

	if dst == "DX" {
		// Destination is DX register
		switch src {
		case "AL":
			opcodeByte = 0xEE
		case "AX", "EAX": // Add EAX case
			opcodeByte = 0xEF // Same opcode for AX and EAX
		default:
			return nil, fmt.Errorf("invalid source register '%s' for OUT DX", src)
		}
	} else {
		// Destination is immediate (imm8)
		hasImm = true
		immValue, err = strconv.ParseUint(dst, 0, 8) // Parse imm8 from dst
		if err != nil {
			return nil, fmt.Errorf("invalid immediate value '%s' for OUT instruction: %v", dst, err)
		}

		switch src {
		case "AL":
			opcodeByte = 0xE6
		case "AX", "EAX": // Add EAX case
			opcodeByte = 0xE7 // Same opcode for AX and EAX
		default:
			return nil, fmt.Errorf("invalid source register '%s' for OUT imm8", src)
		}
	}

	// Assemble the code
	code := []byte{opcodeByte}
	if hasImm {
		code = append(code, byte(immValue))
	}

	return code, nil
}
