package codegen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/HobbyOSs/gosk/pkg/asmdb"
	"github.com/HobbyOSs/gosk/pkg/ocode"
	"github.com/HobbyOSs/gosk/pkg/operand"
)

type OpcodeHandler = func([]string) ([]byte, error)

// handleADD はADD命令の機械語を生成します
func handleADD(operands []string) ([]byte, error) {
	// オペランドを解析
	ops := operand.NewOperandFromString(strings.Join(operands, ","))

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
	opcodeByte, err := strconv.ParseUint(encoding.Opcode.Byte, 16, 8)
	if err != nil {
		return nil, fmt.Errorf("failed to parse opcode byte: %v", err)
	}
	machineCode = append(machineCode, byte(opcodeByte))

	// ModR/Mの追加（必要な場合）
	if encoding.ModRM != nil {
		modrm := generateModRMByte(encoding.ModRM)
		machineCode = append(machineCode, modrm)
	}

	// 即値の追加（必要な場合）
	if encoding.Immediate != nil {
		if imm, err := getImmediateValue(operands[1], encoding.Immediate.Size); err == nil {
			machineCode = append(machineCode, imm...)
		}
	}

	return machineCode, nil
}

// generateModRMByte はModR/Mバイトを生成します
func generateModRMByte(modrm *asmdb.Modrm) byte {
	// ModR/Mバイトの構造:
	// |  Mod (2)  |  Reg (3)  |  R/M (3)  |

	// デフォルトはレジスタ直接アドレッシング（Mod = 11b）
	mod := byte(0b11000000)

	// レジスタコードを取得
	regCode, _ := strconv.ParseInt(modrm.Reg, 0, 8)
	reg := byte(regCode) << 3

	rmCode, _ := strconv.ParseInt(modrm.Rm, 0, 8)
	rm := byte(rmCode)

	return mod | reg | rm
}

var opcodeHandlers = make(map[ocode.OcodeKind]OpcodeHandler)

func init() {
	// ADD命令のハンドラを登録
	opcodeHandlers[ocode.OpADD] = handleADD
}
