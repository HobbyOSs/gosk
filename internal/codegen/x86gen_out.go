package codegen

import (
	"fmt"
	"strconv"
	"strings" // Add strings import

	"github.com/HobbyOSs/gosk/pkg/ng_operand" // Add ng_operand import
)

func handleOUT(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	if len(params.Operands) != 2 {
		return nil, fmt.Errorf("OUT instruction requires 2 operands, but got %d", len(params.Operands))
	}

	// Parse operands to check for prefixes
	opsInterface, err := ng_operand.FromString(strings.Join(params.Operands, ","))
	if err != nil {
		return nil, fmt.Errorf("failed to parse OUT operands '%s': %v", strings.Join(params.Operands, ","), err)
	}
	opsInterface = opsInterface.WithBitMode(ctx.BitMode) // Set bit mode

	dst := strings.ToUpper(params.Operands[0]) // Destination: imm8 or DX
	src := strings.ToUpper(params.Operands[1]) // Source: AL, AX, or EAX

	var opcodeByte byte
	var immValue uint64
	var hasImm bool

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
		var parseErr error                                // Use a separate variable for this scope
		immValue, parseErr = strconv.ParseUint(dst, 0, 8) // Parse imm8 from dst
		if parseErr != nil {
			return nil, fmt.Errorf("invalid immediate value '%s' for OUT instruction: %v", dst, parseErr)
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
	code := []byte{}

	// Add prefixes if required
	if opsInterface.Require66h() { // Check for operand size prefix
		code = append(code, 0x66)
	}
	// Add address size prefix if needed (less common for OUT)
	// if opsInterface.Require67h() {
	// 	code = append(code, 0x67)
	// }

	code = append(code, opcodeByte) // Append opcode
	if hasImm {
		code = append(code, byte(immValue)) // Append immediate if present
	}

	return code, nil
}
