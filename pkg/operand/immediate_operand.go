package operand

import "fmt"

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
