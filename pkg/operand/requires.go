package operand

import (
	"regexp" // regexp パッケージをインポート

	"github.com/HobbyOSs/gosk/pkg/cpu"
	"github.com/samber/lo" // samber/lo をインポート
)

// Require66h はオペランドサイズプレフィックスが必要かどうかを判定する
func (b *OperandImpl) Require66h() bool {
	operandTypes := b.OperandTypes()
	parsedOperands := b.ParsedOperands() // ParsedOperands() を一度だけ呼び出す
	return requireOperandSizePrefix(b.BitMode, operandTypes, parsedOperands)
}

// Require67h はアドレスサイズプレフィックス(0x67)が必要かどうかを判定する
// メモリオペランドで使用される実効アドレスサイズに基づいて判定する
func (b *OperandImpl) Require67h() bool {
	parsedOperands := b.ParsedOperands()
	return requireAddressSizePrefix(b.BitMode, parsedOperands)
}

// requireOperandSizePrefix はビットモードとオペランドタイプ/パース結果に基づいてオペランドサイズプレフィックスが必要かどうかを判定する
func requireOperandSizePrefix(bitMode cpu.BitMode, operandTypes []OperandType, parsedOperands []*ParsedOperand) bool {
	if len(operandTypes) == 0 {
		return false
	}

	switch bitMode {
	case cpu.MODE_16BIT:
		return is32BitRegisterMemoryOperandType(operandTypes) || is32BitImmediateOperandValue(parsedOperands)
	case cpu.MODE_32BIT:
		return is16BitRegisterMemoryOperandType(operandTypes)
	default:
		return false
	}
}

// is32BitRegisterMemoryOperandType はオペランドタイプに32bitレジスタ/メモリが含まれているか判定する
func is32BitRegisterMemoryOperandType(operandTypes []OperandType) bool {
	for _, operandType := range operandTypes {
		if operandType == CodeR32 || operandType == CodeM32 {
			return true
		}
	}
	return false
}

// is16BitRegisterMemoryOperandType はオペランドタイプに16bitレジスタ/メモリが含まれているか判定する
func is16BitRegisterMemoryOperandType(operandTypes []OperandType) bool {
	for _, operandType := range operandTypes {
		if operandType == CodeR16 || operandType == CodeM16 {
			return true
		}
	}
	return false
}

// is32BitImmediateOperandValue はパースされたオペランドに32bit即値が含まれているか判定する
func is32BitImmediateOperandValue(parsedOperands []*ParsedOperand) bool {
	if len(parsedOperands) == 1 && parsedOperands[0].Imm != "" {
		immSize := getImmediateSizeFromValue(parsedOperands[0].Imm)
		return immSize == CodeIMM32
	}
	return false
}

// requireAddressSizePrefix はビットモードとパースされたオペランドに基づいてアドレスサイズプレフィックスが必要かどうかを判定する
func requireAddressSizePrefix(bitMode cpu.BitMode, parsedOperands []*ParsedOperand) bool {
	// メモリアクセスが存在するかチェック
	hasMemoryAccess := lo.SomeBy(parsedOperands, func(op *ParsedOperand) bool {
		return op.IndirectMem != nil || op.DirectMem != nil
	})
	if !hasMemoryAccess {
		return false
	}

	// 32bitアドレッシングを要求するオペランドが存在するかチェック
	requires32BitAddressing := lo.SomeBy(parsedOperands, func(op *ParsedOperand) bool {
		switch {
		case op.IndirectMem != nil:
			if is32bitRegInIndirectMem(op.IndirectMem.Mem) {
				return true // 32bitレジスタが含まれていれば32bitアドレッシング
			}
			// TODO: [disp32] の判定 (現状は16bitとみなされる)
			return false
		case op.DirectMem != nil:
			// TODO: [0x12345678] の判定 (現状は16bitとみなす)
			return false
		default:
			// メモリアクセスでない場合は false
			return false
		}
	})

	// 結果を返す
	switch bitMode {
	case cpu.MODE_16BIT:
		// 16bitモード: 32bitアドレッシングが必要ならプレフィックス要
		return requires32BitAddressing
	case cpu.MODE_32BIT:
		// 32bitモード: 16bitアドレッシングが必要ならプレフィックス要
		// (requires32BitAddressing が false なら 16bit アドレッシングとみなす)
		return !requires32BitAddressing
	default:
		return false
	}
}

// is32bitRegInIndirectMem は IndirectMem 文字列に32bitレジスタが含まれているか判定する
func is32bitRegInIndirectMem(memStr string) bool {
	// 正規表現で32bitレジスタ(EAX, EBX等)をチェック (\bは単語境界)
	re := regexp.MustCompile(`\b(EAX|EBX|ECX|EDX|ESI|EDI|ESP|EBP)\b`)
	return re.MatchString(memStr)
}
