package codegen

import (
	"fmt"
	"strconv"
	"strings" // Add strings import

	"github.com/HobbyOSs/gosk/pkg/cpu"        // Import cpu package
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

	// --- Prefix Calculation (similar to asmdb.GetPrefixSize for OUT) ---
	prefix66Size := 0
	bitMode := ctx.BitMode
	opTypes := opsInterface.OperandTypes() // Get operand types from parsed interface

	if len(opTypes) > 1 { // OUT has 2 operands
		// For OUT, the source (AL/AX/EAX) determines the size
		sizeDeterminingOpType := opTypes[1]
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
	// Address size prefix (0x67) is generally not needed for OUT with DX or imm8
	// if opsInterface.Require67h() {
	// 	code = append(code, 0x67)
	// }

	code = append(code, opcodeByte) // Append opcode
	if hasImm {
		code = append(code, byte(immValue)) // Append immediate if present
	}

	return code, nil
}
