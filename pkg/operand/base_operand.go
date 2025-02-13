package operand

type BaseOperand struct {
	Internal string
}

func (b *BaseOperand) AddressingType() AddressingType {
	parser := getParser()
	parsed, err := parser.ParseString("", b.Internal)
	if err == nil {
		switch {
		case parsed.Reg != "":
			return CodeGeneralReg
		case parsed.Mem != "":
			return CodeModRMAddress
		case parsed.Imm != "":
			return CodeImmediate
		case parsed.Seg != "":
			return CodeSregField
		case parsed.Rel != "":
			return CodeRelativeOffset
		case parsed.Addr != "":
			return CodeDirectAddress
		}
	}
	return AddressingType("unknown")
}

func (b *BaseOperand) OperandType() OperandType {
	parser := getParser()
	parsed, err := parser.ParseString("", b.Internal)
	if err == nil {
		switch {
		case parsed.Reg != "":
			return CodeDoubleword
		case parsed.Mem != "":
			return CodeDoubleword
		case parsed.Imm != "":
			return CodeDoublewordInteger
		case parsed.Seg != "":
			return CodeWord
		case parsed.Rel != "":
			return CodeWord
		case parsed.Addr != "":
			return CodeDoubleword
		}
	}
	return OperandType("unknown")
}
