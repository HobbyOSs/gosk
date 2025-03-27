package operand

import (
	"strings" // strings パッケージをインポート

	"github.com/HobbyOSs/gosk/pkg/cpu"
)

// Require66h はオペランドサイズプレフィックスが必要かどうかを判定する
func (b *OperandImpl) Require66h() bool {
	types := b.OperandTypes()
	if len(types) == 0 {
		return false
	}

	switch b.BitMode {
	case cpu.MODE_16BIT:
		// 16bitモードで32bitレジスタ/メモリを使用する場合
		for _, t := range types {
			if t == CodeR32 || t == CodeM32 {
				return true
			}
		}
		// 16bitモードで32bit即値を使用する場合
		// ParsedOperands() を利用して即値のサイズを判定する
		parsedOperands := b.ParsedOperands()
		// Imm フィールドが string 型であるため、空文字列 "" でチェックする
		if len(parsedOperands) == 1 && parsedOperands[0].Imm != "" {
			// Imm フィールドが string 型であるため、間接参照 * を削除する
			immSize := getImmediateSizeFromValue(parsedOperands[0].Imm)
			if immSize == CodeIMM32 {
				return true
			}
		}

	case cpu.MODE_32BIT:
		// 32bitモードで16bitレジスタ/メモリを使用する場合
		for _, t := range types {
			if t == CodeR16 || t == CodeM16 {
				return true
			}
		}
	}
	return false
}

// Require67h はアドレスサイズプレフィックス(0x67)が必要かどうかを判定する
// メモリオペランドで使用される実効アドレスサイズに基づいて判定する
func (b *OperandImpl) Require67h() bool {
	parsedOperands := b.ParsedOperands()
	if len(parsedOperands) == 0 {
		// パース結果がない場合、OperandTypes からメモリオペランドの存在を確認するかもしれないが、
		// アドレスサイズを特定できないため、一旦 false とする。
		// TODO: ラベル参照 ([label]) の場合のアドレスサイズ解決が必要になる可能性がある。
		return false
	}

	requires32BitAddressing := false // オペランドのいずれかが32bitアドレッシングを要求するか
	hasMemoryAccess := false         // メモリアクセスが実際にあるか

	for _, parsed := range parsedOperands {
		is32BitForThisOperand := false // このオペランドが32bitアドレッシングか

		if parsed.IndirectMem != nil {
			hasMemoryAccess = true
			memStr := parsed.IndirectMem.Mem // 例: "[EBX+16]", "[SI]"
			// 32bitレジスタが含まれているかチェック
			if strings.Contains(memStr, "EAX") || strings.Contains(memStr, "EBX") ||
				strings.Contains(memStr, "ECX") || strings.Contains(memStr, "EDX") ||
				strings.Contains(memStr, "ESI") || strings.Contains(memStr, "EDI") ||
				strings.Contains(memStr, "ESP") || strings.Contains(memStr, "EBP") {
				is32BitForThisOperand = true
			} else {
				// 32bitレジスタがなく、16bitアドレッシングレジスタ ([BX], [BP], [SI], [DI]) のみか、
				// またはディスプレースメントのみ ([1234]) の場合。
				// ディスプレースメントのみの場合のサイズを CalcOffsetByteSize で確認する。
				// ただし、CalcOffsetByteSize はオペランド全体のサイズを返すため、
				// ここで個別にサイズを計算するのは難しい。
				// 一旦、32bitレジスタがなければ16bitアドレッシングとみなす。
				// TODO: [disp32] のようなケースを正しく判定するには CalcOffsetByteSize の改善か、
				//       個別のオペランドに対するサイズ計算が必要。
			}
		} else if parsed.DirectMem != nil {
			hasMemoryAccess = true
			// DirectMem の場合、アドレス値のサイズで判定する。
			// CalcOffsetByteSize はオペランド全体のサイズを返すため、
			// ここで個別にサイズを計算するのは難しい。
			// TODO: [0x12345678] のようなケースを正しく判定するには CalcOffsetByteSize の改善か、
			//       個別のオペランドに対するサイズ計算が必要。
			//       現状では、DirectMem は 16bit アドレス ([0x1234]) とみなす。
		}

		if is32BitForThisOperand {
			requires32BitAddressing = true
		}
	}

	if !hasMemoryAccess {
		return false // メモリアクセスがなければ不要
	}

	// 現在のモードと、要求されるアドレッシングモードを比較
	switch b.GetBitMode() {
	case cpu.MODE_16BIT:
		// 16bitモードで、実効アドレッシングが32bitの場合に 0x67 が必要
		return requires32BitAddressing
	case cpu.MODE_32BIT:
		// 32bitモードで、実効アドレッシングが16bitの場合に 0x67 が必要
		// (requires32BitAddressing が false であれば 16bit アドレッシングとみなす)
		return !requires32BitAddressing
	default:
		return false
	}
}
