package operand

type RegisterOperand struct {
	reg string
}

func (r RegisterOperand) AddressingType() AddressingType {
	return CodeGeneralReg
}

func (r RegisterOperand) OperandType() OperandType {
	return CodeWord
}

func (r RegisterOperand) Serialize() string {
	return r.reg
}

func (r RegisterOperand) FromString(text string) Operand {
	return RegisterOperand{reg: text}
}
