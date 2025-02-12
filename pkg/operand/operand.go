package operand

import "fmt"

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

type ImmediateOperand struct {
	value int
}

func (i ImmediateOperand) AddressingType() AddressingType {
	return CodeImmediate
}

func (i ImmediateOperand) OperandType() OperandType {
	return CodeDoubleword
}

func (i ImmediateOperand) Serialize() string {
	return fmt.Sprintf("#%d", i.value)
}

func (i ImmediateOperand) FromString(text string) Operand {
	var value int
	fmt.Sscanf(text, "#%d", &value)
	return ImmediateOperand{value: value}
}

type MemoryOperand struct {
	base         string
	index        string
	scale        int
	displacement int
}

func (m MemoryOperand) AddressingType() AddressingType {
	return CodeModRMAddress
}

func (m MemoryOperand) OperandType() OperandType {
	return CodeDoubleword
}

func (m MemoryOperand) Serialize() string {
	return fmt.Sprintf("[%s %s*%d +%d]", m.base, m.index, m.scale, m.displacement)
}

func (m MemoryOperand) FromString(text string) Operand {
	var base, index string
	var scale, displacement int
	fmt.Sscanf(text, "[%s %s*%d +%d]", &base, &index, &scale, &displacement)
	return MemoryOperand{base: base, index: index, scale: scale, displacement: displacement}
}
