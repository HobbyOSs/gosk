package ng_operand

import (
	"strings"

	"github.com/HobbyOSs/gosk/internal/ast" // Import ast
	"github.com/HobbyOSs/gosk/pkg/cpu"
	"github.com/samber/lo" // lo をインポート
)

// Require66h はオペランドサイズプレフィックス (66h) が必要か判定します。
// ビットモードとオペランドの *本来の* サイズ (依存関係解決前) の不一致をチェックします。
func (o *OperandPegImpl) Require66h() bool {
	if len(o.parsedOperands) == 0 {
		return false
	}

	is16bitMode := o.bitMode == cpu.MODE_16BIT
	is32bitMode := o.bitMode == cpu.MODE_32BIT // 16ビットでなければ32ビットと仮定

	for _, parsed := range o.parsedOperands {
		if parsed == nil {
			continue
		}

		// 依存関係解決前のオペランドの本来のサイズを決定
		var inherentSize int // 0: 不明, 8, 16, 32, 64
		baseType := parsed.Type

		switch {
		case isR8Type(baseType) || (baseType == CodeM && parsed.DataType == ast.Byte):
			inherentSize = 8
		case isR16Type(baseType) || (baseType == CodeM && parsed.DataType == ast.Word):
			inherentSize = 16
		case isR32Type(baseType) || (baseType == CodeM && parsed.DataType == ast.Dword):
			inherentSize = 32
		case isR64Type(baseType): // TODO: QWORD がサポートされたら M64 チェックを追加
			inherentSize = 64
		case baseType == CodeIMM || baseType == CodeIMM8 || baseType == CodeIMM16 || baseType == CodeIMM32 || baseType == CodeIMM64:
			// 即値自体から推定されるサイズを使用
			immSize := getImmediateSizeType(parsed.Immediate)
			switch immSize {
			case CodeIMM8:
				inherentSize = 8
			case CodeIMM16:
				inherentSize = 16
			case CodeIMM32:
				inherentSize = 32
			case CodeIMM64:
				inherentSize = 64
			}
		case baseType == CodeM: // 明示的な DataType なしのメモリ
			// メモリアドレスで使用されるレジスタに基づいて推定し、モードサイズにデフォルト設定
			// この部分は resolveMemorySize ロジックと重複するため注意が必要
			// ここでは単純化のため、他の情報がなければモードサイズにデフォルト設定すると仮定
			// より堅牢な解決策には組み合わせたロジックが必要になる可能性がある
			// まずレジスタに基づいて推定を試みる
			mem := parsed.Memory
			if mem != nil {
				if strings.HasPrefix(mem.BaseReg, "E") || strings.HasPrefix(mem.IndexReg, "E") || mem.BaseReg == "ESP" || mem.IndexReg == "ESP" || mem.BaseReg == "EBP" || mem.IndexReg == "EBP" {
					inherentSize = 32
				} else if lo.Contains([]string{"BX", "SI", "DI", "SP", "BP"}, mem.BaseReg) || lo.Contains([]string{"SI", "DI"}, mem.IndexReg) {
					inherentSize = 16
				}
			}
			// それでも不明な場合は、モードに基づいてデフォルト設定
			if inherentSize == 0 {
				if is16bitMode {
					inherentSize = 16
				} else {
					inherentSize = 32
				}
			}
		}

		// 不一致をチェック
		// 16ビットモードでオペランドが本来32ビットの場合、66h が必要
		if is16bitMode && inherentSize == 32 {
			return true
		}
		// 32ビットモードでオペランドが本来16ビットの場合、66h が必要
		// 例外: IMM8 は通常、32ビットモードでは 66h を必要としない (オペコードのバリアントで処理される)
		if is32bitMode && inherentSize == 16 {
			// 後で命令コンテキストに基づいてこれを調整する必要があるかもしれない
			// 現時点では、32ビットモードの16ビットオペランドは 66h が必要と仮定
			return true
		}
	}

	return false
}

// Require67h はアドレスサイズプレフィックス (67h) が必要か判定します。
// ビットモードとメモリオペランドのアドレス指定の不一致をチェックします。
func (o *OperandPegImpl) Require67h() bool { // レシーバーを追加
	is16bitMode := o.bitMode == cpu.MODE_16BIT
	is32bitMode := o.bitMode == cpu.MODE_32BIT // 16ビットでなければ32ビットと仮定

	for _, parsed := range o.parsedOperands { // パース済みオペランドを反復処理
		if parsed == nil || parsed.Memory == nil {
			continue // メモリオペランドでなければスキップ
		}
		mem := parsed.Memory

		// アドレス指定に使われているレジスタを確認
		hasEprefix := strings.HasPrefix(mem.BaseReg, "E") || strings.HasPrefix(mem.IndexReg, "E")
		// 特定の16ビットアドレッシングレジスタをチェック
		has16bitAddrReg := mem.BaseReg == "BX" || mem.BaseReg == "BP" || mem.BaseReg == "SI" || mem.BaseReg == "DI" ||
			mem.IndexReg == "SI" || mem.IndexReg == "DI"
		// ディスプレースメントのみが使用されているかチェック (モードのデフォルトアドレスサイズを意味する)
		// Displacement を "" ではなく 0 と比較
		onlyDisplacement := mem.BaseReg == "" && mem.IndexReg == "" && mem.Displacement != 0

		// 16ビットモードで32ビットアドレッシング (EAXなど) を使用する場合
		if is16bitMode && hasEprefix {
			return true
		}
		// 32ビットモードで16ビットアドレッシング (BX, SI, DI, BP) を使用する場合
		// (ただし、ディスプレースメントのみが使用されている場合は除く。その場合はデフォルトの32ビットサイズを使用するため)
		if is32bitMode && has16bitAddrReg && !onlyDisplacement {
			// 完全な精度のためには、より複雑なSIBバイトの考慮が必要になる可能性がある
			return true
		}
	}

	return false
}
