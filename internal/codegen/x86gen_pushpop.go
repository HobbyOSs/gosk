package codegen

import (
	"fmt"
	"strconv" // Add strconv import
	"strings"

	"github.com/HobbyOSs/gosk/pkg/cpu" // Add cpu import for BitMode constants
	"github.com/HobbyOSs/gosk/pkg/ng_operand"
)

// handlePUSH generates machine code for the PUSH instruction.
func handlePUSH(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	if len(params.Operands) != 1 { // Check length of original string operands
		return nil, fmt.Errorf("PUSH instruction requires 1 operand, got %d", len(params.Operands))
	}

	opStr := params.Operands[0]
	opsInterface, err := ng_operand.FromString(opStr) // Returns Operands interface
	if err != nil {
		return nil, fmt.Errorf("failed to parse PUSH operand '%s': %v", opStr, err)
	}
	opsInterface = opsInterface.WithBitMode(ctx.BitMode) // Set bit mode

	// Get operand types using the interface method
	opTypes := opsInterface.OperandTypes()
	if len(opTypes) != 1 {
		return nil, fmt.Errorf("expected one operand type for PUSH, got %d from '%s'", len(opTypes), opStr)
	}
	opType := opTypes[0]

	// Assemble the code buffer
	code := []byte{}

	// Add prefixes if required
	if opsInterface.Require66h() { // Check for operand size prefix
		code = append(code, 0x66)
	}
	if opsInterface.Require67h() { // Check for address size prefix
		code = append(code, 0x67)
	}

	// Determine opcode based on operand type using interface methods and correct constants
	switch opType {
	// --- Register Operands ---
	// case ng_operand.CodeREG: // Generic type might not be returned by OperandTypes()
	case ng_operand.CodeR8, ng_operand.CodeR16, ng_operand.CodeR32, ng_operand.CodeR64,
		ng_operand.CodeAL, ng_operand.CodeCL, ng_operand.CodeDL, ng_operand.CodeBL,
		ng_operand.CodeAH, ng_operand.CodeCH, ng_operand.CodeDH, ng_operand.CodeBH,
		ng_operand.CodeAX, ng_operand.CodeCX, ng_operand.CodeDX, ng_operand.CodeBX,
		ng_operand.CodeSP, ng_operand.CodeBP, ng_operand.CodeSI, ng_operand.CodeDI,
		ng_operand.CodeEAX, ng_operand.CodeECX, ng_operand.CodeEDX, ng_operand.CodeEBX,
		ng_operand.CodeESP, ng_operand.CodeEBP, ng_operand.CodeESI, ng_operand.CodeEDI,
		ng_operand.CodeRAX, // ng_operand.CodeRCX, ng_operand.CodeRDX, ng_operand.CodeRBX, // 64bit regs might not have specific types yet
		ng_operand.CodeSREG, ng_operand.CodeES, ng_operand.CodeCS, ng_operand.CodeSS,
		ng_operand.CodeDS, ng_operand.CodeFS, ng_operand.CodeGS:
		// Use the original string for register lookup
		regCode, ok := registerToPushPopCode(opStr)
		if !ok {
			// Handle segment registers separately if needed (e.g., PUSH CS is 0x0E)
			switch strings.ToUpper(opStr) {
			case "ES":
				code = append(code, 0x06)
				return code, nil
			case "CS":
				code = append(code, 0x0E)
				return code, nil
			case "SS":
				code = append(code, 0x16)
				return code, nil
			case "DS":
				code = append(code, 0x1E)
				return code, nil
			case "FS":
				code = append(code, 0x0F, 0xA0)
				return code, nil // 0F A0
			case "GS":
				code = append(code, 0x0F, 0xA8)
				return code, nil // 0F A8
			default:
				return nil, fmt.Errorf("unsupported register for PUSH: %s", opStr)
			}
		}
		// PUSH r16/r32/r64: 50+rd
		code = append(code, 0x50+regCode)
		return code, nil

	// --- Memory Operands ---
	case ng_operand.CodeMEM, // Generic memory types
		ng_operand.CodeM, ng_operand.CodeM8, ng_operand.CodeM16, ng_operand.CodeM32, ng_operand.CodeM64:
		// PUSH m16/m32/m64: FF /6
		memInfo, found := opsInterface.GetMemoryInfo() // Use interface method
		if !found || memInfo == nil {
			return nil, fmt.Errorf("could not get memory info for PUSH operand: %s", opStr)
		}
		// Use calculateModRM from x86gen_utils.go, passing regBits for the /6 extension
		modrmByte, sibByte, dispBytes, err := calculateModRM(memInfo, ctx.BitMode, 6<<3) // Pass regBits = 6 << 3
		if err != nil {
			return nil, fmt.Errorf("failed to calculate ModRM/SIB/Disp for PUSH %s: %w", opStr, err)
		}
		code = append(code, 0xFF) // Append opcode after prefixes
		code = append(code, modrmByte)
		if sibByte != 0 {
			code = append(code, sibByte)
		}
		code = append(code, dispBytes...)
		return code, nil

	// --- Immediate Operands ---
	case ng_operand.CodeIMM, ng_operand.CodeIMM8, ng_operand.CodeIMM16, ng_operand.CodeIMM32, ng_operand.CodeIMM64: // Use correct constants
		// PUSH imm8: 6A ib
		// PUSH imm16/32: 68 iw/id
		// Parse immediate value from the original string
		immVal, err := strconv.ParseInt(opStr, 0, 64) // Parse from opStr
		if err != nil {
			// Could it be a label treated as immediate? Check SymTable.
			if addr, ok := ctx.SymTable[opStr]; ok {
				immVal = int64(addr) // Use address from symbol table
			} else {
				return nil, fmt.Errorf("failed to parse immediate or find symbol '%s' for PUSH: %v", opStr, err)
			}
		}

		if immVal >= -128 && immVal <= 127 { // Check if imm8 fits
			code = append(code, 0x6A, byte(immVal)) // Append opcode and immediate
			return code, nil
		} else if ctx.BitMode == cpu.MODE_16BIT { // Use cpu constant
			code = append(code, 0x68) // Append opcode
			code = append(code, byte(immVal&0xFF), byte((immVal>>8)&0xFF))
			return code, nil
		} else { // 32 or 64 bit mode (PUSH imm64 is not directly supported, uses PUSH imm32)
			code = append(code, 0x68) // Append opcode
			code = append(code, byte(immVal&0xFF), byte((immVal>>8)&0xFF), byte((immVal>>16)&0xFF), byte((immVal>>24)&0xFF))
			return code, nil
		}
	// Add other operand types like segment registers if needed
	default:
		return nil, fmt.Errorf("unsupported operand type %v for PUSH: %s", opType, opStr) // Log opType and opStr
	}
}

// handlePOP generates machine code for the POP instruction.
func handlePOP(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	if len(params.Operands) != 1 { // Check length of original string operands
		return nil, fmt.Errorf("POP instruction requires 1 operand, got %d", len(params.Operands))
	}

	opStr := params.Operands[0]
	opsInterface, err := ng_operand.FromString(opStr) // Returns Operands interface
	if err != nil {
		return nil, fmt.Errorf("failed to parse POP operand '%s': %v", opStr, err)
	}
	opsInterface = opsInterface.WithBitMode(ctx.BitMode) // Set bit mode

	// Get operand types using the interface method
	opTypes := opsInterface.OperandTypes()
	if len(opTypes) != 1 {
		return nil, fmt.Errorf("expected one operand type for POP, got %d from '%s'", len(opTypes), opStr)
	}
	opType := opTypes[0]

	// Assemble the code buffer
	code := []byte{}

	// Add prefixes if required
	if opsInterface.Require66h() { // Check for operand size prefix
		code = append(code, 0x66)
	}
	if opsInterface.Require67h() { // Check for address size prefix
		code = append(code, 0x67)
	}

	// Determine opcode based on operand type using interface methods and correct constants
	switch opType {
	// --- Register Operands ---
	// case ng_operand.CodeREG: // Generic type might not be returned by OperandTypes()
	case ng_operand.CodeR8, ng_operand.CodeR16, ng_operand.CodeR32, ng_operand.CodeR64,
		ng_operand.CodeAL, ng_operand.CodeCL, ng_operand.CodeDL, ng_operand.CodeBL,
		ng_operand.CodeAH, ng_operand.CodeCH, ng_operand.CodeDH, ng_operand.CodeBH,
		ng_operand.CodeAX, ng_operand.CodeCX, ng_operand.CodeDX, ng_operand.CodeBX,
		ng_operand.CodeSP, ng_operand.CodeBP, ng_operand.CodeSI, ng_operand.CodeDI,
		ng_operand.CodeEAX, ng_operand.CodeECX, ng_operand.CodeEDX, ng_operand.CodeEBX,
		ng_operand.CodeESP, ng_operand.CodeEBP, ng_operand.CodeESI, ng_operand.CodeEDI,
		ng_operand.CodeRAX, // ng_operand.CodeRCX, ng_operand.CodeRDX, ng_operand.CodeRBX, // 64bit regs might not have specific types yet
		ng_operand.CodeSREG, ng_operand.CodeES, ng_operand.CodeCS, ng_operand.CodeSS,
		ng_operand.CodeDS, ng_operand.CodeFS, ng_operand.CodeGS:
		// Use the original string for register lookup
		regCode, ok := registerToPushPopCode(opStr)
		if !ok {
			// Handle segment registers separately if needed (e.g., POP ES is 0x07)
			switch strings.ToUpper(opStr) {
			case "ES":
				code = append(code, 0x07)
				return code, nil
			// case "CS": // POP CS is invalid
			case "SS":
				code = append(code, 0x17)
				return code, nil
			case "DS":
				code = append(code, 0x1F)
				return code, nil
			case "FS":
				code = append(code, 0x0F, 0xA1)
				return code, nil // 0F A1
			case "GS":
				code = append(code, 0x0F, 0xA9)
				return code, nil // 0F A9
			default:
				return nil, fmt.Errorf("unsupported register for POP: %s", opStr)
			}
		}
		// POP r16/r32/r64: 58+rd
		code = append(code, 0x58+regCode)
		return code, nil

	// --- Memory Operands ---
	case ng_operand.CodeMEM, // Generic memory types
		ng_operand.CodeM, ng_operand.CodeM8, ng_operand.CodeM16, ng_operand.CodeM32, ng_operand.CodeM64:
		// POP m16/m32/m64: 8F /0
		memInfo, found := opsInterface.GetMemoryInfo() // Use interface method
		if !found || memInfo == nil {
			return nil, fmt.Errorf("could not get memory info for POP operand: %s", opStr)
		}
		// Use calculateModRM from x86gen_utils.go, passing regBits for the /0 extension
		modrmByte, sibByte, dispBytes, err := calculateModRM(memInfo, ctx.BitMode, 0<<3) // Pass regBits = 0 << 3
		if err != nil {
			return nil, fmt.Errorf("failed to calculate ModRM/SIB/Disp for POP %s: %w", opStr, err)
		}
		code = append(code, 0x8F) // Append opcode after prefixes
		code = append(code, modrmByte)
		if sibByte != 0 {
			code = append(code, sibByte)
		}
		code = append(code, dispBytes...)
		return code, nil
	// POP imm is invalid
	default:
		return nil, fmt.Errorf("unsupported operand type %v for POP: %s", opType, opStr) // Log opType and opStr
	}
}

// registerToPushPopCode maps a register name to its 3-bit code used in PUSH/POP opcodes (50+rd / 58+rd).
func registerToPushPopCode(reg string) (byte, bool) {
	// Normalize to uppercase
	regUpper := strings.ToUpper(reg)

	// Map common registers to their codes
	// Note: This mapping assumes 16/32/64 bit modes share the base code.
	// Size prefixes (like 0x66) are handled elsewhere if needed.
	switch regUpper {
	case "AX", "EAX", "RAX":
		return 0, true
	case "CX", "ECX", "RCX":
		return 1, true
	case "DX", "EDX", "RDX":
		return 2, true
	case "BX", "EBX", "RBX":
		return 3, true
	case "SP", "ESP", "RSP":
		return 4, true
	case "BP", "EBP", "RBP":
		return 5, true
	case "SI", "ESI", "RSI":
		return 6, true
	case "DI", "EDI", "RDI":
		return 7, true
	// Add segment registers if needed (e.g., PUSH CS uses 0x0E)
	default:
		return 0, false // Register not directly supported by 50+rd/58+rd encoding
	}
}
