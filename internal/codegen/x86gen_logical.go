package codegen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/HobbyOSs/gosk/pkg/asmdb"
	// "github.com/HobbyOSs/gosk/pkg/cpu" // Removed unused import
	"github.com/HobbyOSs/gosk/pkg/operand" // Added import
)

// generateLogicalCode は論理命令の機械語生成の共通処理を行う関数です。
func generateLogicalCode(operands []string, ctx *CodeGenContext, instName string) ([]byte, error) {
	// オペランド数のチェック (ANDは2オペランド)
	if len(operands) != 2 {
		return nil, fmt.Errorf("%s requires 2 operands", instName)
	}

	// オペランドの解析
	ops := operand.NewOperandFromString(strings.Join(operands, ",")).
		WithBitMode(ctx.BitMode). // Added WithBitMode
		WithForceImm8(true)       // 算術命令に合わせて一旦 true に設定

	// AsmDBからエンコーディングを取得
	db := asmdb.NewInstructionDB()
	encoding, err := db.FindEncoding(instName, ops)
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
	var regNum int = -1 // Addendがない場合のデフォルト値
	if encoding.Opcode.Addend != nil {
		operandIndex, err := strconv.Atoi(strings.TrimPrefix(*encoding.Opcode.Addend, "#"))
		if err != nil {
			return nil, fmt.Errorf("failed to parse addend")
		}
		if operandIndex < len(operands) {
			regNum, err = GetRegisterNumber(operands[operandIndex])
			if err != nil {
				// レジスタでない場合もあるのでエラーにしない
				regNum = -1 // レジスタ番号が取得できない場合は -1
			}
		}
	}

	opcode, err := ResolveOpcode(encoding.Opcode, regNum)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve opcode")
	}
	machineCode = append(machineCode, opcode...) // Use spread operator

	// ModR/Mの追加（必要な場合）
	if encoding.ModRM != nil {
		modrm, err := getModRMFromOperands(operands, encoding, ctx.BitMode)
		if err != nil {
			return nil, fmt.Errorf("failed to generate ModR/M")
		}
		machineCode = append(machineCode, modrm...)
	}

	// 即値の追加(必要な場合)
	if encoding.Immediate != nil {
		immIndex, err := parseIndex(encoding.Immediate.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid Immediate.Value format")
		}
		if immIndex < len(operands) {
			if imm, err := getImmediateValue(operands[immIndex], encoding.Immediate.Size); err == nil {
				machineCode = append(machineCode, imm...)
			} else {
				return nil, fmt.Errorf("failed to get immediate value")
			}
		} else {
			return nil, fmt.Errorf("immediate index out of range")
		}
	}

	return machineCode, nil
}

func handleAND(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	return generateLogicalCode(params.Operands, ctx, "AND")
}

func handleOR(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	return generateLogicalCode(params.Operands, ctx, "OR")
}

func handleXOR(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	return generateLogicalCode(params.Operands, ctx, "XOR")
}

func handleSHR(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	return generateLogicalCode(params.Operands, ctx, "SHR")
}

func handleSHL(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	return generateLogicalCode(params.Operands, ctx, "SHL")
}

func handleSAR(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	return generateLogicalCode(params.Operands, ctx, "SAR")
}

func handleNOT(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	// オペランド数のチェック (NOTは1オペランド)
	if len(params.Operands) != 1 {
		return nil, fmt.Errorf("NOT requires 1 operand")
	}

	// オペランドの解析
	ops := operand.NewOperandFromString(params.Operands[0]).
		WithBitMode(ctx.BitMode) // Added WithBitMode

	// AsmDBからエンコーディングを取得
	db := asmdb.NewInstructionDB()
	encoding, err := db.FindEncoding("NOT", ops)
	if err != nil {
		return nil, fmt.Errorf("failed to find encoding for NOT")
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
	// REX prefix handling removed based on user feedback

	// オペコードの追加 (NOTはAddendを使用しない)
	opcode, err := ResolveOpcode(encoding.Opcode, -1) // regNumは不要なので-1
	if err != nil {
		return nil, fmt.Errorf("failed to resolve opcode")
	}
	machineCode = append(machineCode, opcode...) // Use spread operator

	// ModR/Mの追加（必要な場合）
	if encoding.ModRM != nil {
		// NOTは1オペランドなので、getModRMFromOperandsは使えない
		// 必要な情報を直接渡してModR/Mを生成するヘルパーが必要かもしれない
		// ここでは仮実装として、getModRMFromOperandsを流用してみる（要修正）
		modrm, err := getModRMFromOperands(params.Operands, encoding, ctx.BitMode)
		if err != nil {
			return nil, fmt.Errorf("failed to generate ModR/M for NOT")
		}
		machineCode = append(machineCode, modrm...)
	}

	// NOT命令は即値を取らない

	return machineCode, nil
}
