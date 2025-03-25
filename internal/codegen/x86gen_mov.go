package codegen

import (
	"log"
	"strconv"
	"strings"

	"github.com/HobbyOSs/gosk/pkg/asmdb"
	"github.com/HobbyOSs/gosk/pkg/operand"
)

// handleMOV handles the MOV instruction and generates the appropriate machine code.
func handleMOV(operands []string, ctx *CodeGenContext) []byte {
	if len(operands) != 2 {
		log.Printf("error: MOV requires 2 operands, got %d", len(operands))
		return nil
	}

	// オペランドの解析
	ops := operand.NewOperandFromString(strings.Join(operands, ",")).
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

	// プレフィックスの追加
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
	machineCode = append(machineCode, opcode)

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
		// ModRMがない場合は
		// メモリオペランドが有る場合のoffset取得して設定する必要が有る場合が有る
		for _, opStr := range operands {
			if _, _, disp, err := operand.ParseMemoryOperand(opStr, ctx.BitMode); err == nil {
				machineCode = append(machineCode, disp...)
			}
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
			machineCode = append(machineCode, byte(addr&0xFF), byte((addr>>8)&0xFF))
		} else if imm, err := getImmediateValue(opStr, encoding.Immediate.Size); err == nil {
			machineCode = append(machineCode, imm...)
		}
	}

	log.Printf("debug: Generated machine code: % x", machineCode)

	return machineCode
}
