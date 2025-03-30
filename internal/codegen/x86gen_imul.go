package codegen

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/HobbyOSs/gosk/pkg/asmdb"
	"github.com/HobbyOSs/gosk/pkg/ng_operand"
)

func handleIMUL(params x86genParams, ctx *CodeGenContext) ([]byte, error) {
	instName := "IMUL"
	// IMUL はオペランド数が 1, 2, 3 の場合がある
	if len(params.Operands) < 1 || len(params.Operands) > 3 {
		log.Printf("error: %s requires 1, 2 or 3 operands, got %d", instName, len(params.Operands))
		return nil, fmt.Errorf("%s requires 1, 2 or 3 operands", instName)
	}

	ops, err := ng_operand.FromString(strings.Join(params.Operands, ","))
	if err != nil {
		log.Printf("Error creating operands from string in %s: %v", instName, err)
		return nil, fmt.Errorf("failed to create operands from string")
	}
	ops = ops.WithBitMode(ctx.BitMode)

	// FindEncoding を matchAnyImm = false で呼び出し、正確なエンコーディングを取得
	db := asmdb.NewInstructionDB()
	encoding, err := db.FindEncoding(instName, ops, false) // Use false for codegen
	if err != nil {
		// false で見つからない場合、念のため true で再試行 (通常は不要なはず)
		log.Printf("warn: FindEncoding(matchAnyImm=false) failed for %s %v, retrying with true: %v", instName, params.Operands, err) // Use %v instead of %w
		encoding, err = db.FindEncoding(instName, ops, true)
		if err != nil {
			return nil, fmt.Errorf("failed to find encoding for %s %v even with matchAnyImm=true: %w", instName, params.Operands, err) // fmt.Errorf can use %w
		}
	}

	log.Printf("debug: handleIMUL: Selected encoding Opcode=%s, ModRM=%+v, Immediate=%+v",
		encoding.Opcode.Byte, encoding.ModRM, encoding.Immediate)

	// --- オペランド拡張は行わない ---
	finalOperands := params.Operands // 元のオペランドを使用

	// --- 機械語生成 ---
	machineCode := make([]byte, 0)

	// プレフィックス
	if ops.Require66h() {
		machineCode = append(machineCode, 0x66)
	}
	if ops.Require67h() {
		machineCode = append(machineCode, 0x67)
	}

	// オペコード
	var regNum int = -1
	if encoding.Opcode.Addend != nil {
		operandIndex, err := strconv.Atoi(strings.TrimPrefix(*encoding.Opcode.Addend, "#"))
		if err != nil {
			return nil, fmt.Errorf("failed to parse addend: %w", err)
		}
		if operandIndex < len(params.Operands) { // 元のオペランドで参照
			regNum, err = GetRegisterNumber(params.Operands[operandIndex])
			if err != nil {
				return nil, fmt.Errorf("failed to get register number: %w", err)
			}
		} else {
			log.Printf("warn: Addend index %d out of range for original operands %v", operandIndex, params.Operands)
		}
	}
	opcodeBytes, err := ResolveOpcode(encoding.Opcode, regNum)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve opcode: %w", err)
	}
	machineCode = append(machineCode, opcodeBytes...)

	// ModR/M
	if encoding.ModRM != nil {
		// --- IMUL Opcode 69/6B の ModRM 特殊処理 ---
		tempEncoding := *encoding // encoding のコピーを作成
		// encoding.ModRM が nil でないことを確認してからコピーを作成
		if tempEncoding.ModRM != nil && (tempEncoding.Opcode.Byte == "69" || tempEncoding.Opcode.Byte == "6B") {
			// JSON定義が Reg:"#1", Rm:"#0" となっていると仮定し、
			// 正しい動作 (Reg:"#0", Rm:"#0") になるように強制的に書き換える
			log.Printf("debug: handleIMUL: Applying ModRM workaround for Opcode %s. Original ModRM: %+v", tempEncoding.Opcode.Byte, tempEncoding.ModRM)
			modRMDefCopy := *tempEncoding.ModRM // ModRM 定義のコピーを作成
			modRMDefCopy.Reg = "#0"             // reg フィールドをオペランド 0 参照に強制
			modRMDefCopy.Rm = "#0"              // rm フィールドをオペランド 0 参照に強制 (念のため)
			tempEncoding.ModRM = &modRMDefCopy  // コピーした定義で上書き
			log.Printf("debug: handleIMUL: Modified ModRM for generation: %+v", tempEncoding.ModRM)
		}
		// --- 特殊処理ここまで ---

		// ModR/M の生成には元のオペランド `finalOperands` (=`params.Operands`) と、
		// 必要に応じて変更された `tempEncoding` を使用
		modrm, err := getModRMFromOperands(finalOperands, &tempEncoding, ctx.BitMode)
		if err != nil {
			return nil, fmt.Errorf("failed to generate ModR/M for %v: %w", finalOperands, err)
		}
		machineCode = append(machineCode, modrm...)
	}

	// 即値
	if encoding.Immediate != nil {
		// 即値の取得には元のオペランド `finalOperands` (=`params.Operands`) を使用
		// Immediate.Value が "#1" であることを期待 (2オペランド構文の2番目を参照)
		immIndex, err := parseIndex(encoding.Immediate.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid Immediate.Value format: %s", encoding.Immediate.Value)
		}
		if immIndex < 0 || immIndex >= len(finalOperands) {
			return nil, fmt.Errorf("immediate index %d out of range for final operands %v", immIndex, finalOperands)
		}
		imm, err := getImmediateValue(finalOperands[immIndex], encoding.Immediate.Size)
		if err != nil {
			return nil, fmt.Errorf("failed to get immediate value from '%s': %w", finalOperands[immIndex], err)
		}
		machineCode = append(machineCode, imm...)
	}

	return machineCode, nil
}
