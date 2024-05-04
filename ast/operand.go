package ast

type Operand interface {
	AddressingType() AddressingType
	OperandType() OperandType
}
