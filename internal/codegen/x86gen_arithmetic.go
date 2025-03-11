package codegen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/HobbyOSs/gosk/pkg/asmdb"
	"github.com/HobbyOSs/gosk/pkg/operand"
)

// handleADD はADD命令の機械語を生成します
func handleADD(operands []string) ([]byte, error) {
	// オペランドを解析
	ops := operand.
		NewOperandFromString(strings.Join(operands, ",")).
		WithForceImm8(true)

	// asmdbからエンコーディング情報を取得
	db := asmdb.NewInstructionDB()
	encoding, err := db.FindEncoding("ADD", ops)
	if err != nil {
		return nil, fmt.Errorf("failed to find encoding for ADD: %v", err)
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
			return nil, fmt.Errorf("failed to parse addend: %v", err)
		}

		// operandsからレジスタ名を取得し、番号に変換
		regNum, err = GetRegisterNumber(operands[operandIndex])
		if err != nil {
			return nil, fmt.Errorf("failed to get register number: %v", err)
		}
	}

	opcode, err := ResolveOpcode(encoding.Opcode, regNum)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve opcode: %v", err)
	}
	machineCode = append(machineCode, opcode)

	// ModR/Mの追加（必要な場合）
	if encoding.ModRM != nil {
		modrm, err := getModRMFromOperands(operands, encoding)
		if err != nil {
			return nil, fmt.Errorf("failed to generate ModR/M: %v", err)
		}
		machineCode = append(machineCode, modrm)
	}

	// 即値の追加(必要な場合)
	if encoding.Immediate != nil {
		if imm, err := getImmediateValue(operands[1], encoding.Immediate.Size); err == nil {
			machineCode = append(machineCode, imm...)
		}
	}

	return machineCode, nil
}
