package operand

type RegisterOperand struct {
	BaseOperand
	reg string
}

func (r RegisterOperand) InternalString() string {
	return r.Internal
}

func (r RegisterOperand) AddressingType() AddressingType {
	return r.AddressingType()
}

func (r RegisterOperand) OperandType() OperandType {
	return r.OperandType()
}

func (r RegisterOperand) Serialize() string {
	return r.reg
}

func (r RegisterOperand) FromString(text string) Operand {
	return RegisterOperand{BaseOperand: BaseOperand{Internal: text}, reg: text}
}
