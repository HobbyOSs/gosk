package codegen

import (
	"fmt"
	"log"
	"strconv" // Added missing import
	"strings" // Added missing import

	"github.com/HobbyOSs/gosk/pkg/asmdb"
	"github.com/HobbyOSs/gosk/pkg/ng_operand" // Use ng_operand
)

// generateArithmeticCode は算術命令の機械語生成の共通処理を行う関数です。
func generateArithmeticCode(operands []string, ctx *CodeGenContext, instName string) ([]byte, error) {
	if len(operands) != 2 {
		log.Printf("error: %s requires 2 operands, got %d", instName, len(operands))
		return nil, fmt.Errorf("%s requires 2 operands", instName)
	}

	// オペランドの解析
	ops, err := ng_operand.FromString(strings.Join(operands, ","))
	if err != nil {
		// TODO: より適切なエラーハンドリングを行う
		log.Printf("Error creating operands from string in %s: %v", instName, err)
		return nil, fmt.Errorf("failed to create operands from string")
	}
	ops = ops.WithBitMode(ctx.BitMode)

	// AsmDBからエンコーディングを取得 (matchAnyImm = true)
	db := asmdb.NewInstructionDB()
	encoding, err := db.FindEncoding(instName, ops, true)
	if err != nil {
		return nil, fmt.Errorf("failed to find encoding for %s", instName)
	}

	// エンコーディング情報を使用して機械語を生成
	machineCode := make([]byte, 0)

	// プレフィックスの追加
	if ops.Require66h() {
		machineCode = append(machineCode, 0x66)
	}
	if ops.Require67h() {
		machineCode = append(machineCode, 0x67)
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
	machineCode = append(machineCode, opcode...) // Use spread operator

	// ModR/Mの追加（必要な場合）
	if encoding.ModRM != nil {
		modrm, err := getModRMFromOperands(operands, encoding, ctx.BitMode)
		if err != nil {
			log.Printf("error: Failed to generate ModR/M: %v", err)
			return nil, fmt.Errorf("failed to generate ModR/M")
		}
		machineCode = append(machineCode, modrm...)
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
		} else {
			log.Printf("error: Failed to get immediate value: %v", err) // Added error handling
			return nil, fmt.Errorf("failed to get immediate value")
		}
	}

	return machineCode, nil
}

func handleADD(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	return generateArithmeticCode(params.Operands, ctx, "ADD")
}

func handleCMP(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	return generateArithmeticCode(params.Operands, ctx, "CMP")
}

func handleIMUL(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	return generateArithmeticCode(params.Operands, ctx, "IMUL")
}

func handleSUB(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	return generateArithmeticCode(params.Operands, ctx, "SUB")
}
