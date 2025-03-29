package ng_operand

import (
	"strings"

	"github.com/HobbyOSs/gosk/pkg/cpu"
)

// Require66h はオペランドサイズプレフィックス (66h) が必要か判定します。
// ビットモードとオペランドサイズ（解決済み）の不一致をチェックします。
func (o *OperandPegImpl) Require66h() bool { // Add receiver
	opTypes := o.OperandTypes() // Get resolved types
	if len(opTypes) == 0 {
		return false
	}

	is16bitMode := o.bitMode == cpu.MODE_16BIT
	is32bitMode := o.bitMode == cpu.MODE_32BIT // Assume 32-bit if not 16-bit

	for _, opType := range opTypes { // Iterate over all resolved operand types
		is16bitOperand := opType == CodeR16 || opType == CodeM16 || opType == CodeIMM16 ||
			opType == CodeAX || opType == CodeCX || opType == CodeDX || opType == CodeBX ||
			opType == CodeSP || opType == CodeBP || opType == CodeSI || opType == CodeDI
		is32bitOperand := opType == CodeR32 || opType == CodeM32 || opType == CodeIMM32 ||
			opType == CodeEAX || opType == CodeECX || opType == CodeEDX || opType == CodeEBX ||
			opType == CodeESP || opType == CodeEBP || opType == CodeESI || opType == CodeEDI
		// TODO: Add R8/M8/IMM8 checks if needed, though they usually don't trigger 66h.

		// 16bitモードで32bitオペランドを使用する場合
		if is16bitMode && is32bitOperand {
			return true
		}
		// 32bitモードで16bitオペランドを使用する場合
		if is32bitMode && is16bitOperand {
			// IMM8 は 32bit モードでも 66h を必要としないことが多い (e.g., ADD EAX, 1)
			// More accurate check might need instruction context.
			if opType == CodeIMM8 { // Check if it's specifically IMM8
				continue // Skip IMM8 for now (simplification)
			}
			return true
		}
	}

	return false
}

// Require67h はアドレスサイズプレフィックス (67h) が必要か判定します。
// ビットモードとメモリオペランドのアドレス指定の不一致をチェックします。
func (o *OperandPegImpl) Require67h() bool { // Add receiver
	is16bitMode := o.bitMode == cpu.MODE_16BIT
	is32bitMode := o.bitMode == cpu.MODE_32BIT // Assume 32-bit if not 16-bit

	for _, parsed := range o.parsedOperands { // Iterate over parsed operands
		if parsed == nil || parsed.Memory == nil {
			continue // メモリオペランドでなければスキップ
		}
		mem := parsed.Memory

		// アドレス指定に使われているレジスタを確認
		hasEprefix := strings.HasPrefix(mem.BaseReg, "E") || strings.HasPrefix(mem.IndexReg, "E")
		// Check for specific 16-bit addressing registers
		has16bitAddrReg := mem.BaseReg == "BX" || mem.BaseReg == "BP" || mem.BaseReg == "SI" || mem.BaseReg == "DI" ||
			mem.IndexReg == "SI" || mem.IndexReg == "DI"
		// Check if only displacement is used (implies default address size for the mode)
		// Compare Displacement to 0 instead of ""
		onlyDisplacement := mem.BaseReg == "" && mem.IndexReg == "" && mem.Displacement != 0

		// 16bitモードで32bitアドレッシング (EAXなど) を使用する場合
		if is16bitMode && hasEprefix {
			return true
		}
		// 32bitモードで16bitアドレッシング (BX, SI, DI, BP) を使用する場合
		// (But not if only displacement is used, as that uses the default 32-bit size)
		if is32bitMode && has16bitAddrReg && !onlyDisplacement {
			// More complex SIB byte considerations might be needed for full accuracy.
			return true
		}
	}

	return false
}
