package operand_test

import (
	"testing"

	"github.com/HobbyOSs/gosk/pkg/operand"
)

func TestBaseOperand_AddressingType(t *testing.T) {
	tests := []struct {
		name     string
		internal string
		expected operand.AddressingType
	}{
		{"General Register", "EAX", operand.CodeGeneralReg},
		{"Memory Address", "[EBX]", operand.CodeModRMAddress},
		{"Immediate Value", "0x10", operand.CodeImmediate},
		{"Segment Register", "CS", operand.CodeSregField},
		{"Relative Offset", "JMP LABEL", operand.CodeRelativeOffset},
		{"Direct Address", "[0x1234]", operand.CodeDirectAddress},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := operand.BaseOperand{Internal: tt.internal}
			if got := b.AddressingType(); got != tt.expected {
				t.Errorf("AddressingType() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBaseOperand_OperandType(t *testing.T) {
	tests := []struct {
		name     string
		internal string
		expected operand.OperandType
	}{
		{"General Register", "EAX", operand.CodeDoubleword},
		{"Memory Address", "[EBX]", operand.CodeDoubleword},
		{"Immediate Value", "0x10", operand.CodeDoublewordInteger},
		{"Segment Register", "CS", operand.CodeWord},
		{"Relative Offset", "JMP LABEL", operand.CodeWord},
		{"Direct Address", "[0x1234]", operand.CodeDoubleword},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := operand.BaseOperand{Internal: tt.internal}
			if got := b.OperandType(); got != tt.expected {
				t.Errorf("OperandType() = %v, want %v", got, tt.expected)
			}
		})
	}
}
