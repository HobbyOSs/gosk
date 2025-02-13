package operand

import "fmt"

type ImmediateOperand struct {
	BaseOperand
	value    int
	internal string
}

func (i ImmediateOperand) AddressingType() AddressingType {
	return i.AddressingType()
}

func (i ImmediateOperand) OperandType() OperandType {
	return i.OperandType()
}

func (i ImmediateOperand) InternalString() string {
	return i.internal
}

func (i ImmediateOperand) Serialize() string {
	return fmt.Sprintf("#%d", i.value)
}

func (i ImmediateOperand) FromString(text string) Operand {
	var value int
	fmt.Sscanf(text, "#%d", &value)
	return ImmediateOperand{BaseOperand: BaseOperand{internal: text}, value: value}
}
