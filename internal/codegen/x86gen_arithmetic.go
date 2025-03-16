package codegen

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/HobbyOSs/gosk/pkg/asmdb"
	"github.com/HobbyOSs/gosk/pkg/operand"
)

// generateArithmeticCode は算術命令の機械語生成の共通処理を行う関数です。
func generateArithmeticCode(operands []string, ctx *CodeGenContext, instName string) ([]byte, error) {
	if len(operands) != 2 {
		log.Printf("error: %s requires 2 operands, got %d", instName, len(operands))
		return nil, fmt.Errorf("%s requires 2 operands", instName)
	}

	// オペランドの解析
	ops := operand.NewOperandFromString(strings.Join(operands, ",")).
		WithForceImm8(true)

	// AsmDBからエンコーディングを取得
	db := asmdb.NewInstructionDB()
	encoding, err := db.FindEncoding(instName, ops)
	if err != nil {
		log.Printf("error: Failed to find encoding for %s: %v", instName, err)
		return nil, fmt.Errorf("failed to find encoding for %s", instName)
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
			return nil, fmt.Errorf("failed to parse addend")
		}

		// operandsからレジスタ名を取得し、番号に変換
		regNum, err = GetRegisterNumber(operands[operandIndex])
		if err != nil {
			log.Printf("error: %v", err)
			return nil, err
		}
	}

	opcode, err := ResolveOpcode(encoding.Opcode, regNum)
	if err != nil {
		log.Printf("error: Failed to resolve opcode: %v", err)
		return nil, fmt.Errorf("failed to resolve opcode")
	}
	machineCode = append(machineCode, opcode)

	// ModR/Mの追加（必要な場合）
	if encoding.ModRM != nil {
		log.Printf("debug: ModRM: %+v", encoding.ModRM)
		modrm, err := getModRMFromOperands(operands, encoding, ctx.BitMode)
		if err != nil {
			log.Printf("error: Failed to generate ModR/M: %v", err)
			return nil, fmt.Errorf("failed to generate ModR/M")
		}
		machineCode = append(machineCode, byte(modrm))
	}

	// 即値の追加(必要な場合)
	if encoding.Immediate != nil {
		immIndex, err := parseIndex(encoding.Immediate.Value)
		if err != nil {
			log.Printf("error: invalid Immediate.Value format")
			return nil, fmt.Errorf("invalid Immediate.Value format")
		}
		if imm, err := getImmediateValue(operands[immIndex], encoding.Immediate.Size); err == nil {
			machineCode = append(machineCode, imm...)
		}
	}

	log.Printf("debug: Generated machine code: % x", machineCode)

	// TODO: フラグ設定

	return machineCode, nil
}

func handleADD(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	return generateArithmeticCode(params.Operands, ctx, "ADD")
}

func handleCMP(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	return generateArithmeticCode(params.Operands, ctx, "CMP")
}
