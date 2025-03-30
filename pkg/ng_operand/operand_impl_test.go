package ng_operand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImmediateValueFitsIn8Bits(t *testing.T) {
	tests := []struct {
		name     string
		value    int64
		expected bool
	}{
		{"Zero", 0, true},
		{"Positive Max", 127, true},
		{"Positive Out of Range", 128, false},
		{"Negative Min", -128, true},
		{"Negative Out of Range", -129, false},
		{"Large Positive", 4608, false},
		{"Large Negative", -4608, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a dummy ParsedOperandPeg with the immediate value
			parsedOp := &ParsedOperandPeg{
				Type:      CodeIMM, // Type doesn't matter for this specific test
				Immediate: tt.value,
			}
			// Create OperandPegImpl with this single operand
			opImpl := NewOperandPegImpl([]*ParsedOperandPeg{parsedOp})

			assert.Equal(t, tt.expected, opImpl.ImmediateValueFitsIn8Bits())
		})
	}

	// Test case with no immediate operand
	t.Run("NoImmediateOperand", func(t *testing.T) {
		parsedOp := &ParsedOperandPeg{
			Type: CodeAX, // A register type
		}
		opImpl := NewOperandPegImpl([]*ParsedOperandPeg{parsedOp})
		assert.False(t, opImpl.ImmediateValueFitsIn8Bits())
	})

	// Test case with nil parsedOperands (edge case)
	t.Run("NilParsedOperands", func(t *testing.T) {
		opImpl := NewOperandPegImpl(nil)
		assert.False(t, opImpl.ImmediateValueFitsIn8Bits())
	})
}
