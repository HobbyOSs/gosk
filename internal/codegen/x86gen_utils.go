package codegen

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/HobbyOSs/gosk/pkg/asmdb"
	"github.com/HobbyOSs/gosk/pkg/operand"
)

// GenerateModRM はエンコーディング情報とビットモードに基づいてModR/Mバイトを生成する
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

// parseMode は modeStr から mode バイトを解析する
func parseMode(modeStr string) byte {
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
	return mode
}

// ModRMByOperand はモード、regオペランド、rmオペランド、ビットモードに基づいてModR/Mバイトを生成する
func ModRMByOperand(modeStr string, regOperand string, rmOperand string, bitMode ast.BitMode) ([]byte, error) {
	// ModR/M バイトの生成
	// |  mod  |  reg  |  r/m  |
	// | 7 6 | 5 4 3 | 2 1 0 |

	mode := parseMode(modeStr) // modeの解析

	// regの解析（3ビット）
	reg, err := GetRegisterNumber(regOperand)
	if err != nil {
		return nil, fmt.Errorf("failed to get register number for %s: %w", regOperand, err)
	}
	regBits := byte(reg) << 3

	// r/mの解析
	if strings.Contains(rmOperand, "[") && strings.HasSuffix(rmOperand, "]") {
		// メモリオペランドの場合: pkg/operand.ParseMemoryOperand を使用
		modBits, rmBits, disp, err := operand.ParseMemoryOperand(rmOperand, bitMode)
		if err != nil {
			return nil, fmt.Errorf("failed to parse memory operand '%s' (mode:%d): %w", rmOperand, bitMode, err)
		}

		// 注釈: operand.ParseMemoryOperand は mod ビット (00, 01, 10) を直接返す。
		// parseMode(modeStr) から得られる 'mode' 変数は、命令定義からの意図されたアドレッシングモード (#0, #1, #2, 11) を示す可能性がある。
		// 実際に必要なエンコーディングを反映するため、ParseMemoryOperand によって返された mod ビットを使用すべきである。
		modrmByte := modBits | regBits | rmBits
		log.Printf("debug: ModRMByOperand (mem): modeStr=%s, reg=%s(%b), rm=%s(%b), parsedMod=%b, parsedRm=%b, result=%#x, disp=% x",
			modeStr, regOperand, regBits, rmOperand, rmBits, modBits, rmBits, modrmByte, disp)

		out := []byte{modrmByte}
		if disp != nil {
			out = append(out, disp...)
		}
		return out, nil
	}

	// r/m がレジスタの場合 (mod=11)
	rm, err := GetRegisterNumber(rmOperand)
	if err != nil {
		return nil, fmt.Errorf("failed to get register number for %s: %w", rmOperand, err)
	}
	rmBits := byte(rm)

	// parseMode によって決定された mode を使用
	out := mode | regBits | rmBits
	log.Printf("debug: GenerateModRM: mode=%s(%b), reg=%s(%b), rm=%s(%b), result=%#x", modeStr, mode, regOperand, regBits, rmOperand, rmBits, out)
	return []byte{out}, nil
}

// ModRMByValue はモード、固定reg値、rmオペランドに基づいてModR/Mバイトを生成する
func ModRMByValue(modeStr string, regValue int, rmOperand string, bitMode ast.BitMode) []byte {
	// ModR/M バイトの生成
	// |  mod  |  reg  |  r/m  |
	// | 7 6 | 5 4 3 | 2 1 0 |

	mode := parseMode(modeStr) // modeの解析

	// regの解析（3ビット）
	regBits := byte(regValue) << 3

	// r/mの解析（3ビット）
	if strings.Contains(rmOperand, "[") && strings.HasSuffix(rmOperand, "]") {
		// メモリオペランドの場合: pkg/operand.ParseMemoryOperand を使用
		modBits, rmBits, disp, err := operand.ParseMemoryOperand(rmOperand, bitMode)
		if err != nil {
			log.Printf("error: Failed to parse memory operand for rm in ModRMByValue: %v", err)
			// []byte{0} の代わりにエラーを返すことを検討
			return []byte{0}
		}

		// パースされた mod/rm ビットを使用
		modrmByte := modBits | regBits | rmBits
		log.Printf("debug: ModRMByValue (mem): modeStr=%s, reg=%d(%b), rm=%s(%b), parsedMod=%b, parsedRm=%b, result=%#x, disp=% x",
			modeStr, regValue, regBits, rmOperand, rmBits, modBits, rmBits, modrmByte, disp)

		out := []byte{modrmByte}
		if disp != nil {
			out = append(out, disp...)
		}
		return out
	}

	// r/m がレジスタの場合 (mod=11)
	rm, err := GetRegisterNumber(rmOperand)
	if err != nil {
		log.Printf("error: Failed to get register number for rm: %v", err)
		return []byte{0}
	}
	rmBits := byte(rm)

	// parseMode によって決定された mode を使用
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
