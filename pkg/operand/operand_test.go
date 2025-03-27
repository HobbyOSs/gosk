package operand

import (
	"testing"
)

func TestRequire66h(t *testing.T) {
	tests := []struct {
		name     string
		operand  string
		bitMode  BitMode
		expected bool
	}{
		{
			name:     "16bit mode with 32bit register",
			operand:  "EAX",
			bitMode:  MODE_16BIT,
			expected: true,
		},
		{
			name:     "16bit mode with 16bit register",
			operand:  "AX",
			bitMode:  MODE_16BIT,
			expected: false,
		},
		{
			name:     "32bit mode with 16bit register",
			operand:  "AX",
			bitMode:  MODE_32BIT,
			expected: true,
		},
		{
			name:     "32bit mode with 32bit register",
			operand:  "EAX",
			bitMode:  MODE_32BIT,
			expected: false,
		},
		{
			name:     "16bit mode with 32bit immediate",
			operand:  "0x12345678",
			bitMode:  MODE_16BIT,
			expected: true,
		},
		{
			name:     "16bit mode with 16bit immediate",
			operand:  "0x1234",
			bitMode:  MODE_16BIT,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := NewOperandFromString(tt.operand).WithBitMode(tt.bitMode)
			if got := op.Require66h(); got != tt.expected {
				t.Errorf("Require66h() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRequire67h(t *testing.T) {
	tests := []struct {
		name     string
		operand  string
		bitMode  BitMode
		expected bool
	}{
		{
			name:     "16bit mode with 32bit memory access",
			operand:  "DWORD [EBX]",
			bitMode:  MODE_16BIT,
			expected: true,
		},
		{
			name:     "16bit mode with 16bit memory access",
			operand:  "WORD [BX]",
			bitMode:  MODE_16BIT,
			expected: false,
		},
		{
			name:     "32bit mode with 16bit memory access",
			operand:  "WORD [BX]",
			bitMode:  MODE_32BIT,
			expected: true,
		},
		{
			name:     "32bit mode with 32bit memory access",
			operand:  "DWORD [EBX]",
			bitMode:  MODE_32BIT,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := NewOperandFromString(tt.operand).WithBitMode(tt.bitMode)
			if got := op.Require67h(); got != tt.expected {
				t.Errorf("Require67h() = %v, want %v", got, tt.expected)
			}
		})
	}
}
