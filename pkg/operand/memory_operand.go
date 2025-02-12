package operand

import "fmt"

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
