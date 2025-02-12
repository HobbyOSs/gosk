package operand

type Operand interface {
	AddressingType() AddressingType
	OperandType() OperandType
}
