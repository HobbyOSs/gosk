package operand

type BaseOperand struct {
	internal string
}

func (b *BaseOperand) AddressingType() AddressingType {
	parser := getParser()
	if _, err := parser.ParseString("", b.internal); err == nil {
		// TODO
		return CodeGeneralReg
	}
	return ""
}

func (b *BaseOperand) OperandType() OperandType {
	parser := getParser()
	if _, err := parser.ParseString("", b.internal); err == nil {
		// TODO
		return CodeDoubleword
	}
	return ""
}
