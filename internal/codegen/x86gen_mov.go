package codegen

import (
	// "fmt" // Remove unused fmt import
	// Add binary package for displacement conversion
	"log"     // Keep only one log import
	"strconv" // Keep only one strconv import
	"strings" // Keep only one strings import

	"github.com/HobbyOSs/gosk/pkg/asmdb"      // Add cpu package for BitMode constants
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
	} else {
		// ModRMがない場合、直接アドレス指定 (moffs) の可能性がある。
		// ng_operand の DisplacementBytes() を使ってディスプレースメントを取得・追加する。
		dispBytes := ops.DisplacementBytes()
		if dispBytes != nil {
			log.Printf("debug: Adding direct memory displacement (moffs) using DisplacementBytes(): %x", dispBytes)
			machineCode = append(machineCode, dispBytes...)
		} else {
			// ModRMも直接メモリアドレスもない場合、エラーか、あるいは即値のみのパターンかもしれない
			// （例：MOV AX, imm16 は ModRM なしだが、即値処理は後段で行われる）
			// ここでは特に何もしない
			log.Printf("debug: No ModRM and DisplacementBytes() returned nil in MOV.")
		}
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
