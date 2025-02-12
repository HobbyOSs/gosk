package operand_test

import (
	"testing"

	"github.com/HobbyOSs/gosk/pkg/operand"
)

func TestOperandBuilder(t *testing.T) {
	tests := []struct {
		name     string
		operand  operand.Operand
		expected string
	}{
		{
			name:     "Register Operand",
			operand:  operand.OperandBuilder{}.Reg("EAX"),
			expected: "EAX",
		},
		{
			name:     "Immediate Operand",
			operand:  operand.OperandBuilder{}.Imm(42),
			expected: "#42",
		},
		{
			name:     "Memory Operand",
			operand:  operand.OperandBuilder{}.Mem("EBX", "ECX", 2, 8),
			expected: "[EBX ECX*2 +8]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.operand.Serialize(); got != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, got)
			}
		})
	}
}
