package codegen

import (
	// "fmt" // Remove unused fmt import
	"log"     // Keep only one log import
	"strconv" // Keep only one strconv import
	"strings" // Keep only one strings import

	"github.com/HobbyOSs/gosk/pkg/asmdb"
	"github.com/HobbyOSs/gosk/pkg/ng_operand" // Use ng_operand
)

// handleMOV handles the MOV instruction and generates the appropriate machine code.
func handleMOV(operands []string, ctx *CodeGenContext) []byte {
	if len(operands) != 2 {
		log.Printf("error: MOV requires 2 operands, got %d", len(operands))
		return nil
	}

	// オペランドの解析
	ops, err := ng_operand.FromString(strings.Join(operands, ","))
	if err != nil {
		// TODO: より適切なエラーハンドリングを行う
		log.Printf("Error creating operands from string in MOV: %v", err)
		return nil
	}
	ops = ops.WithBitMode(ctx.BitMode).
		WithForceRelAsImm(true)

	// AsmDBからエンコーディングを取得
	db := asmdb.NewInstructionDB()
	encoding, err := db.FindEncoding("MOV", ops)
	if err != nil {
		log.Printf("error: Failed to find encoding: %v", err)
		return nil
	}

	// エンコーディング情報を使用して機械語を生成
	machineCode := make([]byte, 0)

	// プレフィックスの追加 (0x67 アドレスサイズ, 0x66 オペランドサイズ の順)
	if ops.Require67h() {
		machineCode = append(machineCode, 0x67)
	}
	if ops.Require66h() {
		machineCode = append(machineCode, 0x66)
	}

	// オペコードの追加
	// Addendから処理すべきオペランドのインデックスを取得
	var regNum int
	if encoding.Opcode.Addend != nil {
		operandIndex, err := strconv.Atoi(strings.TrimPrefix(*encoding.Opcode.Addend, "#"))
		if err != nil {
			log.Printf("error: Failed to parse addend: %v", err)
			return nil
		}

		// operandsからレジスタ名を取得し、番号に変換
		regNum, err = GetRegisterNumber(operands[operandIndex])
		if err != nil {
			log.Printf("error: %v", err)
			return nil
		}
	}

	opcode, err := ResolveOpcode(encoding.Opcode, regNum)
	if err != nil {
		log.Printf("error: Failed to resolve opcode: %v", err)
		return nil
	}
	machineCode = append(machineCode, opcode...) // Use spread operator

	// ModR/Mの追加（必要な場合）
	if encoding.ModRM != nil {
		log.Printf("debug: ModRM: %+v", encoding.ModRM)
		modrm, err := getModRMFromOperands(operands, encoding, ctx.BitMode)
		if err != nil {
			log.Printf("error: Failed to generate ModR/M: %v", err)
			return nil
		}
		machineCode = append(machineCode, modrm...)
	// } else {
	// TODO: Handle case where ModRM is nil but displacement might be needed.
	// This requires inspecting the parsed operands from 'ops'.
	// The old logic using operand.ParseMemoryOperand needs replacement.
	// This might be handled within getModRMFromOperands or requires new logic here.
	// For now, commenting out the old logic.
	// // ModRMがない場合は
	// // メモリオペランドが有る場合のoffset取得して設定する必要が有る場合が有る
	// for _, opStr := range operands {
	// 	// Replace operand.ParseMemoryOperand logic
	// 	// Need to get displacement from ng_operand.Operands 'ops'
	// }
	}

	// 即値の追加(必要な場合)
	if encoding.Immediate != nil {
		immIndex, err := parseIndex(encoding.Immediate.Value)
		if err != nil {
			log.Printf("error: invalid Immediate.Value format")
			return nil
		}

		opStr := operands[immIndex]
		// アドレスが即値扱いされるパターンと通常の即値を処理する
		if addr, ok := ctx.SymTable[opStr]; ok {
			// Assuming immediate size dictates how many bytes of the address to use
			switch encoding.Immediate.Size {
			case 1:
				machineCode = append(machineCode, byte(addr&0xFF))
			case 2:
				machineCode = append(machineCode, byte(addr&0xFF), byte((addr>>8)&0xFF))
			case 4:
				machineCode = append(machineCode, byte(addr&0xFF), byte((addr>>8)&0xFF), byte((addr>>16)&0xFF), byte((addr>>24)&0xFF))
			default:
				log.Printf("error: unsupported immediate size for symbol address: %d", encoding.Immediate.Size)
				return nil
			}
		} else if imm, err := getImmediateValue(opStr, encoding.Immediate.Size); err == nil {
			machineCode = append(machineCode, imm...)
		} else {
			log.Printf("error: Failed to get immediate value or symbol address for %s: %v", opStr, err) // Added error handling
			return nil                                                                                  // Return nil on error
		}
	}

	log.Printf("debug: Generated machine code: % x", machineCode)

	return machineCode
}
