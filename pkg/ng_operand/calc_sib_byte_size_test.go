package ng_operand

import (
	"testing"

	"github.com/HobbyOSs/gosk/pkg/cpu"
	"github.com/stretchr/testify/assert"
)

func TestCalcSibByteSize(t *testing.T) {
	tests := []struct {
		name     string
		operand  string
		bitMode  cpu.BitMode
		expected int
	}{
		// 32ビットモードでSIBバイトが必要なケース
		{"32bit_BaseIndex", "[EAX+EBX]", cpu.MODE_32BIT, 1},
		{"32bit_BaseIndexScale", "[EAX+EBX*4]", cpu.MODE_32BIT, 1},
		{"32bit_ESPBase", "[ESP]", cpu.MODE_32BIT, 1}, // ESPをベースレジスタとする場合
		{"32bit_ESPBaseIndex", "[ESP+EAX]", cpu.MODE_32BIT, 1},
		{"32bit_ESPBaseIndexScale", "[ESP+EAX*4]", cpu.MODE_32BIT, 1},
		{"32bit_IndexScaleOnly", "[EBX*4]", cpu.MODE_32BIT, 1}, // ベースレジスタなし、インデックスレジスタあり

		// 32ビットモードでSIBバイトが不要なケース
		{"32bit_DirectAddress", "[0x1000]", cpu.MODE_32BIT, 0}, // 直接アドレス指定
		{"32bit_EBPBasedNoIndex", "[EBP]", cpu.MODE_32BIT, 0}, // EBPベースでインデックスなし
		{"32bit_EBPBasedDisp", "[EBP+4]", cpu.MODE_32BIT, 0}, // EBPベースでインデックスなし、ディスプレースメントあり
		{"32bit_BaseOnlyNoESP", "[EAX]", cpu.MODE_32BIT, 0}, // ESP以外のベースレジスタのみ

		// 16ビットモードでは常にSIBバイトは不要
		{"16bit_BaseIndex", "[BX+SI]", cpu.MODE_16BIT, 0},
		{"16bit_DirectAddress", "[0x1000]", cpu.MODE_16BIT, 0},
		{"16bit_BaseOnly", "[BX]", cpu.MODE_16BIT, 0},

		// メモリオペランドではないケース
		{"RegisterOperand", "EAX", cpu.MODE_32BIT, 0},
		{"ImmediateOperand", "123", cpu.MODE_32BIT, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op, err := FromString(tt.operand)
			assert.NoError(t, err)
			opImpl := op.(*OperandPegImpl).WithBitMode(tt.bitMode).(*OperandPegImpl) // 型アサーションとビットモード設定

			assert.Equal(t, tt.expected, opImpl.CalcSibByteSize())
		})
	}

	// Edge case: nil parsedOperands
	t.Run("NilParsedOperands", func(t *testing.T) {
		opImpl := NewOperandPegImpl(nil)
		assert.Equal(t, 0, opImpl.CalcSibByteSize())
	})
}
