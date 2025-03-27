package codegen

import (
	"encoding/binary"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/HobbyOSs/gosk/pkg/asmdb"
	"github.com/HobbyOSs/gosk/pkg/operand"
)

// GenerateModRM generates ModR/M byte based on encoding information and bit mode.
func GenerateModRM(operands []string, modRM *asmdb.Encoding, bitMode ast.BitMode) ([]byte, error) {
	if modRM == nil || modRM.ModRM == nil {
		return nil, nil
	}

	modRMDef := modRM.ModRM
	if modRMDef == nil {
		return nil, nil
	}

	if strings.HasPrefix(modRMDef.Reg, "#") {
		// ModR/M の reg フィールドがオペランドの場合
		regIndex, err := parseIndex(modRMDef.Reg)
		if err != nil {
			return nil, fmt.Errorf("invalid ModRM.Reg format")
		}
		rmIndex, err := parseIndex(modRMDef.Rm)
		if err != nil {
			return nil, fmt.Errorf("invalid ModRM.RM format")
		}

		if regIndex < 0 || regIndex >= len(operands) {
			return nil, fmt.Errorf("ModRM.Reg index out of range")
		}
		if rmIndex < 0 || rmIndex >= len(operands) {
			return nil, fmt.Errorf("ModRM.RM index out of range")
		}

		regOperand := operands[regIndex]
		rmOperand := operands[rmIndex]
		return ModRMByOperand(modRMDef.Mode, regOperand, rmOperand, bitMode)
	} else {
		// ModR/M の reg フィールドが固定値の場合
		regValue, err := strconv.Atoi(modRMDef.Reg)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ModRM.Reg: %v", err)
		}
		rmIndex, err := parseIndex(modRMDef.Rm)
		if err != nil {
			return nil, fmt.Errorf("invalid ModRM.RM format")
		}
		if rmIndex < 0 || rmIndex >= len(operands) {
			return nil, fmt.Errorf("ModRM.RM index out of range")
		}
		rmOperand := operands[rmIndex]
		return ModRMByValue(modRMDef.Mode, regValue, rmOperand, bitMode), nil
	}
}

// ModRMByOperand generates ModR/M byte based on mode, reg operand, rm operand and bit mode.
func ModRMByOperand(modeStr string, regOperand string, rmOperand string, bitMode ast.BitMode) ([]byte, error) {
	// ModR/M バイトの生成
	// |  mod  |  reg  |  r/m  |
	// | 7 6 | 5 4 3 | 2 1 0 |

	// modeの解析（2ビット）
	var mode byte
	switch modeStr {
	case "#0": // レジスタ間接参照
		mode = 0b00000000
	case "#1": // 8ビット変位レジスタ間接参照
		mode = 0b01000000
	case "#2": // 32ビット変位レジスタ間接参照
		mode = 0b10000000
	case "11": // レジスタ
		mode = 0b11000000
	default:
		mode = 0 // デフォルト値
	}

	// regの解析（3ビット）
	reg, err := GetRegisterNumber(regOperand)
	if err != nil {
		return nil, fmt.Errorf("failed to get register number for %s: %w", regOperand, err)
	}
	regBits := byte(reg) << 3

	// r/mの解析
	if strings.Contains(rmOperand, "[") && strings.HasSuffix(rmOperand, "]") {
		// bitModeに応じて処理を分岐
		if bitMode == ast.MODE_32BIT { // 正しい定数を使用
			// 32bitモードの場合: ParseMemoryOperandを使用
			_, rmNum, disp, err := operand.ParseMemoryOperand(rmOperand, bitMode)
			if err != nil {
				// エラーメッセージを修正: 32bit operand parsing error
				return nil, fmt.Errorf("failed to parse 32bit memory operand '%s': %w", rmOperand, err)
			}
			rmBits := byte(rmNum)

			modrmByte := mode | regBits | rmBits
			log.Printf("debug: ModRMByOperand (mem32): mode=%s(%b), reg=%s(%b), rm=%s(%b), result=%#x", modeStr, mode, regOperand, regBits, rmOperand, rmBits, modrmByte)

			out := []byte{modrmByte}
			if disp != nil {
				out = append(out, disp...)
			}
			return out, nil
		} else if bitMode == ast.MODE_16BIT { // 正しい定数を使用
			// 16bitモードの場合: Manually parse known 16-bit modes
			memContentWithPrefix := strings.TrimSpace(rmOperand[1 : len(rmOperand)-1]) // Content inside []
			// Remove size prefix (BYTE, WORD, etc.) if present
			memContent := memContentWithPrefix
			if strings.HasPrefix(memContent, "BYTE ") {
				memContent = strings.TrimSpace(memContent[5:])
			}
			if strings.HasPrefix(memContent, "WORD ") {
				memContent = strings.TrimSpace(memContent[5:])
			}
			// DWORD should not appear in 16-bit mode, but handle defensively
			if strings.HasPrefix(memContent, "DWORD ") {
				memContent = strings.TrimSpace(memContent[6:])
			}

			var modeBits byte
			var rmBits byte
			var disp []byte
			var dispVal int64
			var dispErr error
			basePart := memContent
			dispSize := 0

			// Check for displacement
			if strings.Contains(memContent, "+") {
				parts := strings.SplitN(memContent, "+", 2)
				basePart = strings.TrimSpace(parts[0])
				dispStr := strings.TrimSpace(parts[1])
				dispVal, dispErr = parseNumeric(dispStr) // Use local parseNumeric
				if dispErr != nil {
					return nil, fmt.Errorf("invalid 16bit displacement '%s': %w", dispStr, dispErr)
				}
			} else if strings.Contains(memContent, "-") {
				// Note: Intel syntax usually doesn't use base-disp, but handle it defensively
				parts := strings.SplitN(memContent, "-", 2)
				basePart = strings.TrimSpace(parts[0])
				dispStr := strings.TrimSpace(parts[1])
				dispVal, dispErr = parseNumeric(dispStr) // Use local parseNumeric
				if dispErr != nil {
					return nil, fmt.Errorf("invalid 16bit displacement '%s': %w", dispStr, dispErr)
				}
				dispVal = -dispVal // Handle subtraction
			}

			// Determine rm bits and preliminary mode/disp based on basePart
			isDirectAddress := false
			switch basePart {
			case "BX+SI":
				rmBits = 0b000
			case "BX+DI":
				rmBits = 0b001
			case "BP+SI":
				rmBits = 0b010
			case "BP+DI":
				rmBits = 0b011
			case "SI":
				rmBits = 0b100
			case "DI":
				rmBits = 0b101
			case "BP":
				rmBits = 0b110 // Special case: [BP] alone implies Mode 01 with disp8=0
			case "BX":
				rmBits = 0b111
			default:
				// Check for direct address [imm16]
				directAddrVal, directAddrErr := parseNumeric(basePart) // Use local parseNumeric
				if directAddrErr == nil {
					rmBits = 0b110        // R/M = 110 for direct address
					modeBits = 0b00000000 // Mode = 00
					dispVal = directAddrVal
					dispSize = 2 // Direct address always uses 16-bit displacement
					isDirectAddress = true
					dispErr = nil // Clear potential error from basePart parsing
				} else {
					// Use original rmOperand in error message for clarity
					return nil, fmt.Errorf("unsupported 16bit base/combination in ModRMByOperand: '%s' from '%s'", basePart, rmOperand)
				}
			}

			// Determine final mode and displacement bytes based on dispVal and basePart
			if !isDirectAddress {
				if dispErr == nil && dispVal != 0 {
					if dispVal >= -128 && dispVal <= 127 {
						modeBits = 0b01000000 // Mode 01 (8-bit disp)
						dispSize = 1
					} else if dispVal >= -32768 && dispVal <= 32767 { // Check 16-bit range
						modeBits = 0b10000000 // Mode 10 (16-bit disp)
						dispSize = 2
					} else {
						return nil, fmt.Errorf("16bit displacement out of range: %d", dispVal)
					}
				} else if basePart == "BP" { // Special case for [BP] or [BP+0] -> Mode 01, disp8=0
					modeBits = 0b01000000
					dispSize = 1
					dispVal = 0 // Ensure dispVal is 0
				} else { // No displacement (and not [BP] alone)
					modeBits = 0b00000000 // Mode 00
					dispSize = 0
					dispVal = 0
				}
			}

			// Prepare displacement bytes
			if dispSize > 0 {
				disp = make([]byte, dispSize)
				if dispSize == 1 {
					disp[0] = byte(dispVal)
				} else { // dispSize == 2
					binary.LittleEndian.PutUint16(disp, uint16(dispVal))
				}
			} else {
				disp = nil
			}

			modrmByte := modeBits | regBits | rmBits
			log.Printf("debug: ModRMByOperand (mem16): calcMode=%b, reg=%s(%b), rm=%s(%b), result=%#x, disp=% x", modeBits, regOperand, regBits, rmOperand, rmBits, modrmByte, disp)

			out := []byte{modrmByte}
			if disp != nil {
				out = append(out, disp...)
			}
			return out, nil

		} else {
			return nil, fmt.Errorf("unknown bitMode: %d", bitMode)
		}
	}

	// r/m がレジスタの場合
	rm, err := GetRegisterNumber(rmOperand)
	if err != nil {
		return nil, fmt.Errorf("failed to get register number for %s: %w", rmOperand, err)
	}
	rmBits := byte(rm)

	out := mode | regBits | rmBits
	log.Printf("debug: GenerateModRM: mode=%s(%b), reg=%s(%b), rm=%s(%b), result=%#x", modeStr, mode, regOperand, regBits, rmOperand, rmBits, out)
	return []byte{out}, nil
}

// ModRMByValue generates ModR/M byte based on mode, fixed reg value, and rm operand.
func ModRMByValue(modeStr string, regValue int, rmOperand string, bitMode ast.BitMode) []byte {
	// ModR/M バイトの生成
	// |  mod  |  reg  |  r/m  |
	// | 7 6 | 5 4 3 | 2 1 0 |

	// modeの解析（2ビット）
	var mode byte
	switch modeStr {
	case "#0": // レジスタ間接参照
		mode = 0b00000000
	case "#1": // 8ビット変位レジスタ間接参照
		mode = 0b01000000
	case "#2": // 32ビット変位レジスタ間接参照
		mode = 0b10000000
	case "11": // レジスタ
		mode = 0b11000000
	default:
		mode = 0 // デフォルト値
	}

	// regの解析（3ビット）
	regBits := byte(regValue) << 3

	// r/mの解析（3ビット）
	// メモリの場合は0として扱う
	if strings.Contains(rmOperand, "[") && strings.HasSuffix(rmOperand, "]") {
		_, rm, disp, err := operand.ParseMemoryOperand(rmOperand, bitMode)
		if err != nil {
			log.Printf("error: Failed to parse memory operand for rm: %v", err)
			return []byte{0}
		}
		rmBits := byte(rm)

		modrmByte := mode | regBits | rmBits
		log.Printf("debug: ModRMByValue: mode=%s(%b), reg=%d(%b), rm=%s(%b), result=%#x", modeStr, mode, regValue, regBits, rmOperand, rmBits, modrmByte)

		out := []byte{modrmByte}
		if disp != nil {
			out = append(out, disp...)
		}
		return out
	}

	rm, err := GetRegisterNumber(rmOperand)
	if err != nil {
		log.Printf("error: Failed to get register number for rm: %v", err)
		return []byte{0}
	}
	rmBits := byte(rm)

	out := mode | regBits | rmBits
	log.Printf("debug: ModRMByValue: mode=%s(%b), reg=%d(%b), rm=%s(%b), result=%#x", modeStr, mode, regValue, regBits, rmOperand, rmBits, out)
	return []byte{out}
}

// GetRegisterNumber はレジスタ名からレジスタ番号（0-7）を取得する
func GetRegisterNumber(regName string) (int, error) {
	switch regName {
	case "AL", "AX", "EAX", "RAX", "ES", "CR0":
		return 0, nil
	case "CL", "CX", "ECX", "RCX", "CS":
		return 1, nil
	case "DL", "DX", "EDX", "RDX", "SS", "CR2":
		return 2, nil
	case "BL", "BX", "EBX", "RBX", "DS", "CR3":
		return 3, nil
	case "AH", "SP", "ESP", "RSP", "FS", "CR4":
		return 4, nil
	case "CH", "BP", "EBP", "RBP", "GS":
		return 5, nil
	case "DH", "SI", "ESI", "RSI":
		return 6, nil
	case "BH", "DI", "EDI", "RDI":
		return 7, nil
	default:
		return 0, fmt.Errorf("unknown register: %s", regName)
	}
}

// ResolveOpcode はOpcodeとレジスタ番号を受け取り、最終的なオペコードバイト列を算出する。
// regNum はレジスタの番号（0-7）を表す。
func ResolveOpcode(op asmdb.Opcode, regNum int) ([]byte, error) {
	opBytes := []byte{}
	opStr := op.Byte

	// オペコード文字列をバイトごとに処理
	if len(opStr)%2 != 0 {
		return nil, fmt.Errorf("invalid opcode string length: %s", opStr)
	}
	for i := 0; i < len(opStr); i += 2 {
		byteStr := opStr[i : i+2]
		byteVal, err := strconv.ParseUint(byteStr, 16, 8)
		if err != nil {
			return nil, fmt.Errorf("invalid opcode byte string: %s in %s", byteStr, opStr)
		}
		opBytes = append(opBytes, byte(byteVal))
	}

	// Addendがある場合、最後のバイトにレジスタ番号を加算
	if op.Addend != nil && len(opBytes) > 0 {
		regBits := byte(regNum & 0x07)
		lastByteIndex := len(opBytes) - 1
		opBytes[lastByteIndex] |= regBits
		log.Printf("debug: ResolveOpcode: base=%s, addend=%v, reg=%d, result=% x", opStr, op.Addend, regNum, opBytes)
	} else {
		log.Printf("debug: ResolveOpcode: base=%s, result=% x", opStr, opBytes)
	}

	return opBytes, nil
}

// getModRMFromOperands はオペランドからModR/Mバイトを生成する
func getModRMFromOperands(operands []string, modRM *asmdb.Encoding, bitMode ast.BitMode) ([]byte, error) {
	modrmByte, err := GenerateModRM(operands, modRM, bitMode)
	if err != nil {
		return nil, err
	}
	return modrmByte, nil
}

func parseIndex(indexStr string) (int, error) {
	if strings.HasPrefix(indexStr, "#") {
		indexStr = indexStr[1:]
	}
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		return -1, fmt.Errorf("invalid index format")
	}
	return index, nil
}

// parseNumeric (copied from pkg/operand/operand_util.go)
// 文字列の数値(例: "16", "0x0ff0", "-123")をint64に変換
func parseNumeric(s string) (int64, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	negative := false
	if strings.HasPrefix(s, "-") {
		negative = true
		s = s[1:]
	}

	base := 10
	if strings.HasPrefix(s, "0x") {
		base = 16
		s = s[2:]
	}
	val, err := strconv.ParseInt(s, base, 64)
	if err != nil {
		return 0, err
	}
	if negative {
		val = -val
	}
	return val, nil
}
