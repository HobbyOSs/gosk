package codegen

import (
	"log"
	"strconv"
	"strings"

	"github.com/HobbyOSs/gosk/pkg/asmdb"
	"github.com/HobbyOSs/gosk/pkg/operand"
)

// handleMOV handles the MOV instruction and generates the appropriate machine code.
func handleMOV(operands []string) []byte {
	if len(operands) != 2 {
		log.Printf("error: MOV requires 2 operands, got %d", len(operands))
		return nil
	}

	// オペランドの解析
	ops := operand.NewOperandFromString(strings.Join(operands, ","))

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
		modrm, err := getModRMFromOperands(operands, encoding)
		if err != nil {
			log.Printf("error: Failed to generate ModR/M: %v", err)
			return nil
		}
		machineCode = append(machineCode, modrm)
	}

	// 即値の追加(必要な場合)
	if encoding.Immediate != nil {
		if imm, err := getImmediateValue(operands[1], encoding.Immediate.Size); err == nil {
			machineCode = append(machineCode, imm...)
		}
	}

	log.Printf("debug: Generated machine code: % x", machineCode)

	return machineCode
}

// getImmediateValue extracts immediate value from operand
func getImmediateValue(operand string, size int) ([]byte, error) {
	// 0xで始まる16進数の場合
	if strings.HasPrefix(operand, "0x") {
		value, err := strconv.ParseUint(operand[2:], 16, size*8)
		if err != nil {
			return nil, err
		}
		return intToBytes(value, size), nil
	}

	// 10進数の場合
	value, err := strconv.ParseInt(operand, 10, size*8)
	if err != nil {
		return nil, err
	}
	return intToBytes(uint64(value), size), nil
}

// intToBytes converts an integer to a byte slice of specified size
func intToBytes(value uint64, size int) []byte {
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		bytes[i] = byte(value >> (i * 8))
	}
	return bytes
}
