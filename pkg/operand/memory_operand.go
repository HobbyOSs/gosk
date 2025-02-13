package operand

import "fmt"

type MemoryOperand struct {
	BaseOperand
	base         string
	index        string
	scale        int
	displacement int
	internal     string
}

func (m MemoryOperand) InternalString() string {
	return m.internal
}

func (m MemoryOperand) AddressingType() AddressingType {
	return m.AddressingType()
}

func (m MemoryOperand) OperandType() OperandType {
	return m.OperandType()
}

func (m MemoryOperand) Serialize() string {
	return fmt.Sprintf("[%s %s*%d +%d]", m.base, m.index, m.scale, m.displacement)
}

func (m MemoryOperand) FromString(text string) Operand {
	var base, index string
	var scale, displacement int
	fmt.Sscanf(text, "[%s %s*%d +%d]", &base, &index, &scale, &displacement)
	return MemoryOperand{BaseOperand: BaseOperand{internal: text}, base: base, index: index, scale: scale, displacement: displacement}
}
