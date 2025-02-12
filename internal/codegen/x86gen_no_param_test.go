package codegen

import (
	"testing"

	"github.com/HobbyOSs/gosk/pkg/ocode"
)

func TestGenerateX86NoParam(t *testing.T) {
	tests := []struct {
		name     string
		ocode    ocode.Ocode
		expected []byte
	}{
		{
			name:     "Single Opcode NOP",
			ocode:    ocode.Ocode{Kind: ocode.OpNOP, Operands: nil},
			expected: []byte{0x90},
		},
		{
			name:     "Single Opcode HLT",
			ocode:    ocode.Ocode{Kind: ocode.OpHLT, Operands: nil},
			expected: []byte{0xF4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateX86NoParam(tt.ocode)
			if !equal(result, tt.expected) {
				t.Errorf("got %v, expected %v", result, tt.expected)
			}
		})
	}
}

func equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
