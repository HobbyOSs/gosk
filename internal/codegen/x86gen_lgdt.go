package codegen

import (
	"log"
	"strings"

	"github.com/HobbyOSs/gosk/pkg/asmdb"
	"github.com/HobbyOSs/gosk/pkg/operand"
)

// handleLGDT handles the LGDT instruction and generates the appropriate machine code.
func handleLGDT(operands []string, ctx *CodeGenContext) ([]byte, error) {
	if len(operands) != 1 {
		log.Printf("error: LGDT requires 1 operand, got %d", len(operands))
		return nil, nil
	}

	// オペランドの解析
	ops := operand.NewOperandFromString(strings.Join(operands, ",")).
		WithBitMode(ctx.BitMode)

	// AsmDBからエンコーディングを取得
	db := asmdb.NewInstructionDB()
	encoding, err := db.FindEncoding("LGDT", ops)
	if err != nil {
		log.Printf("error: Failed to find encoding: %v", err)
		return nil, err
	}

	// エンコーディング情報を使用して機械語を生成
	machineCode := make([]byte, 0)

	// オペコードの追加
	opcode, err := ResolveOpcode(encoding.Opcode, 0) // LGDTではレジスタ番号は使用しない
	if err != nil {
		log.Printf("error: Failed to resolve opcode: %v", err)
		return nil, err
	}
	machineCode = append(machineCode, opcode)

	// ModR/Mの追加
	if encoding.ModRM != nil {
		log.Printf("debug: ModRM: %+v", encoding.ModRM)
		modrm, err := getModRMFromOperands(operands, encoding, ctx.BitMode)
		if err != nil {
			log.Printf("error: Failed to generate ModR/M: %v", err)
			return nil, err
		}
		machineCode = append(machineCode, modrm...)
	}

	// メモリオペランドの解決 (6バイトのデータを追加)
	// TODO: オフセットの計算と追加
	for _, opStr := range operands {
		if _, _, disp, err := operand.ParseMemoryOperand(opStr, ctx.BitMode); err == nil {
			// TODO: 6バイトのデータ (セグメントセレクタ、ベースアドレス、リミット) を追加
			//       今は仮に0で埋める
			for i := 0; i < 6; i++ {
				machineCode = append(machineCode, disp...)
			}
		}
	}

	log.Printf("debug: Generated machine code: % x", machineCode)

	return machineCode, nil
}
