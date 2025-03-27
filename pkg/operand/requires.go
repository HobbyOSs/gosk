package operand

import (
	"github.com/HobbyOSs/gosk/pkg/cpu"
)

// Require66h はオペランドサイズプレフィックスが必要かどうかを判定する
func (b *OperandImpl) Require66h() bool {
	types := b.OperandTypes()
	if len(types) == 0 {
		return false
	}

	switch b.BitMode {
	case cpu.MODE_16BIT: // Changed to cpu.MODE_16BIT
		// 16bitモードで32bitレジスタを使用する場合
		for _, t := range types {
			if t == CodeR32 || t == CodeM32 {
				return true
			}
		}
		// 16bitモードで32bit即値を使用する場合
		if len(types) == 1 {
			parser := getParser()                                   // ここで getParser() を呼び出すのは問題ないか？
			inst, err := parser.ParseString("", b.InternalString()) // b.Internal -> b.InternalString()
			if err == nil && len(inst.Operands) == 1 && inst.Operands[0].Imm != "" {
				imm := getImmediateSizeFromValue(inst.Operands[0].Imm)
				if imm == CodeIMM32 {
					return true
				}
			}
		}
	case cpu.MODE_32BIT: // Changed to cpu.MODE_32BIT
		// 32bitモードで16bitレジスタを使用する場合
		for _, t := range types {
			if t == CodeR16 || t == CodeM16 {
				return true
			}
		}
	}
	return false
}

// Require67h はアドレスサイズプレフィックスが必要かどうかを判定する
func (b *OperandImpl) Require67h() bool {
	types := b.OperandTypes()
	if len(types) == 0 {
		return false
	}

	parsedOperands := b.ParsedOperands()
	if len(parsedOperands) == 0 {
		return false
	}

	switch b.GetBitMode() {
	case cpu.MODE_16BIT:
		// 16bitモードの場合
		for _, parsed := range parsedOperands {
			if parsed.DirectMem != nil || parsed.IndirectMem != nil {
				if parsed.DirectMem != nil && parsed.DirectMem.Prefix != nil {
					t := getMemorySizeFromPrefix(*parsed.DirectMem.Prefix + " " + parsed.DirectMem.Addr)
					return t != CodeM8 && t != CodeM16 && t != CodeM32
				}
				if parsed.IndirectMem != nil && parsed.IndirectMem.Prefix != nil {
					t := getMemorySizeFromPrefix(*parsed.IndirectMem.Prefix + " " + parsed.IndirectMem.Mem)
					return t != CodeM8 && t != CodeM16 && t != CodeM32
				}
				// データサイズプレフィックスなし (例: [0x0ff8], [SI]) の場合は
				// オフセットサイズが2より大きいならば0x67 が必要 (32bitアドレス)
				size := b.CalcOffsetByteSize()
				return size > 2
			}
		}
		return false

	case cpu.MODE_32BIT:
		// 32bitモードの場合
		for _, parsed := range parsedOperands {
			if parsed.DirectMem != nil && parsed.DirectMem.Prefix != nil {
				t := getMemorySizeFromPrefix(*parsed.DirectMem.Prefix + " " + parsed.DirectMem.Addr)
				return t != CodeM32
			}
			if parsed.IndirectMem != nil && parsed.IndirectMem.Prefix != nil {
				t := getMemorySizeFromPrefix(*parsed.IndirectMem.Prefix + " " + parsed.IndirectMem.Mem)
				return t != CodeM32
			}
			size := b.CalcOffsetByteSize()
			return size <= 2
		}
	}
	return false
}
