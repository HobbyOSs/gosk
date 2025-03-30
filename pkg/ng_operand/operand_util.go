package ng_operand

// needsResolution はオペランドタイプがサイズ解決を必要とするか判定します。
func needsResolution(opType OperandType) bool {
	return opType == CodeM || opType == CodeIMM || opType == CodeIMM8 || opType == CodeIMM16
}

// isR32Type は指定された型が32ビット汎用レジスタ型かどうかを判定します。
func isR32Type(opType OperandType) bool {
	switch opType {
	case CodeR32, CodeEAX, CodeECX, CodeEDX, CodeEBX, CodeESP, CodeEBP, CodeESI, CodeEDI:
		return true
	default:
		return false
	}
}

// isR8Type は指定された型が8ビット汎用レジスタ型かどうかを判定します。
func isR8Type(opType OperandType) bool {
	switch opType {
	case CodeR8, CodeAL, CodeCL, CodeDL, CodeBL, CodeAH, CodeCH, CodeDH, CodeBH:
		return true
	default:
		return false
	}
}

// isR64Type は指定された型が64ビット汎用レジスタ型かどうかを判定します。
func isR64Type(opType OperandType) bool {
	// R8-R15 が実装されたら追加する
	switch opType {
	// RAX, RBX などの64ビットレジスタ用の定数が存在するか、追加されることを想定
	// 例: case CodeR64, CodeRAX, CodeRBX, CodeRCX, CodeRDX, CodeRSI, CodeRDI, CodeRSP, CodeRBP:
	case CodeR64: // 特定の64ビットコードが定義/使用されるまでのプレースホルダー
		return true
	default:
		return false
	}
}

// isR16Type は指定された型が16ビット汎用レジスタ型かどうかを判定します。
func isR16Type(opType OperandType) bool {
	switch opType {
	case CodeR16, CodeAX, CodeCX, CodeDX, CodeBX, CodeSP, CodeBP, CodeSI, CodeDI:
		return true
	default:
		return false
	}
}

// isCREGType は指定された型が制御レジスタ型かどうかを判定します。
func isCREGType(opType OperandType) bool {
	return opType == CodeCREG
}

// isRegisterType は指定された型がレジスタ型かどうかを判定します。
func isRegisterType(opType OperandType) bool {
	switch opType {
	case CodeR8, CodeR16, CodeR32, CodeR64,
		CodeAL, CodeCL, CodeDL, CodeBL, CodeAH, CodeCH, CodeDH, CodeBH,
		CodeAX, CodeCX, CodeDX, CodeBX, CodeSP, CodeBP, CodeSI, CodeDI,
		CodeEAX, CodeECX, CodeEDX, CodeEBX, CodeESP, CodeEBP, CodeESI, CodeEDI,
		CodeSREG, CodeCREG, CodeDREG, CodeTREG, CodeMM, CodeXMM, CodeYMM:
		return true
	default:
		return false
	}
}
