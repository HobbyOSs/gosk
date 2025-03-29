package ng_operand // Changed package name

import (
	"testing"

	"github.com/HobbyOSs/gosk/pkg/cpu" // Uncomment cpu import
)

func TestRequire66h(t *testing.T) {
	// t.Skip("Skipping test until NewOperandPegImpl is implemented") // Unskip test
	tests := []struct { // Use cpu.BitMode // Uncomment tests variable definition
		name     string
		operand  string
		bitMode  cpu.BitMode
		expected bool
	}{
		{
			name:     "16bit mode with 32bit register",
			operand:  "EAX",
			bitMode:  cpu.MODE_16BIT,
			expected: true,
		},
		{
			name:     "16bit mode with 16bit register",
			operand:  "AX",
			bitMode:  cpu.MODE_16BIT,
			expected: false,
		},
		{
			name:     "32bit mode with 16bit register",
			operand:  "AX",
			bitMode:  cpu.MODE_32BIT,
			expected: true,
		},
		{
			name:     "32bit mode with 32bit register",
			operand:  "EAX",
			bitMode:  cpu.MODE_32BIT,
			expected: false,
		},
		{
			name:     "16bit mode with 32bit immediate",
			operand:  "0x12345678",
			bitMode:  cpu.MODE_16BIT,
			expected: true,
		},
		{
			name:     "16bit mode with 16bit immediate",
			operand:  "0x1234",
			bitMode:  cpu.MODE_16BIT,
			expected: false,
		},
	}

	// Uncomment test logic
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedOp, err := ParseOperandString(tt.operand) // Use new parser
			if err != nil {
				t.Fatalf("ParseOperandString failed for '%s': %v", tt.operand, err)
			}
			op := NewOperandPegImpl(parsedOp).WithBitMode(tt.bitMode) // Needs implementation of methods in OperandPegImpl
			if got := op.Require66h(); got != tt.expected {
				t.Errorf("Require66h() for operand '%s' = %v, want %v", tt.operand, got, tt.expected)
			}
		})
	}
}

func TestRequire67h(t *testing.T) {
	// t.Skip("Skipping test until NewOperandPegImpl is implemented") // Unskip test
	/* // Comment out tests variable definition */ // Uncomment tests variable definition
	tests := []struct { // Use cpu.BitMode
		name     string
		operand  string
		bitMode  cpu.BitMode
		expected bool
	}{
		{
			name:     "16bit mode with 32bit memory access",
			operand:  "DWORD [EBX]",
			bitMode:  cpu.MODE_16BIT, // Changed to cpu.MODE_16BIT
			expected: true,
		},
		{
			name:     "16bit mode with 16bit memory access",
			operand:  "WORD [BX]",
			bitMode:  cpu.MODE_16BIT, // Changed to cpu.MODE_16BIT
			expected: false,
		},
		{
			name:     "32bit mode with 16bit memory access",
			operand:  "WORD [BX]",
			bitMode:  cpu.MODE_32BIT, // Changed to cpu.MODE_32BIT
			expected: true,
		},
		{
			name:     "32bit mode with 32bit memory access",
			operand:  "DWORD [EBX]",
			bitMode:  cpu.MODE_32BIT, // Changed to cpu.MODE_32BIT
			expected: false,
		},
	}
	// */ // Uncomment tests variable definition

	// Original test logic commented out // Uncomment test logic
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedOp, err := ParseOperandString(tt.operand) // Use new parser
			if err != nil {
				t.Fatalf("ParseOperandString failed for '%s': %v", tt.operand, err)
			}
			op := NewOperandPegImpl(parsedOp).WithBitMode(tt.bitMode) // Needs implementation of methods in OperandPegImpl
			if got := op.Require67h(); got != tt.expected {
				t.Errorf("Require67h() for operand '%s' = %v, want %v", tt.operand, got, tt.expected)
			}
		})
	}
}
