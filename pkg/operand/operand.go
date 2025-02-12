package operand


type Operand interface {
	AddressingType() AddressingType
	OperandType() OperandType
	Serialize() string
	FromString(text string) Operand
}

type OperandBuilder struct{}

func (OperandBuilder) Reg(name string) RegisterOperand {
	return RegisterOperand{reg: name}
}

func (OperandBuilder) Imm(value int) ImmediateOperand {
	return ImmediateOperand{value: value}
}

func (OperandBuilder) Mem(base string, index string, scale int, displacement int) MemoryOperand {
	return MemoryOperand{base: base, index: index, scale: scale, displacement: displacement}
}

