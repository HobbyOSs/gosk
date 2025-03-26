package codegen

import (
	"testing"

	"github.com/HobbyOSs/gosk/pkg/ocode"
	"github.com/stretchr/testify/assert"
)

func TestGenerateX86NoParam(t *testing.T) {
	for kind, expected := range opcodeMap {
		t.Run(kind.String(), func(t *testing.T) {
			ocode := ocode.Ocode{Kind: kind, Operands: nil}
			result := GenerateX86NoParam(ocode)
			assert.Equal(t, []byte{expected}, result, "Opcode %s should generate correct binary", kind.String())
		})
	}
}

func TestHandleRET(t *testing.T) {
	testCases := []struct {
		name     string
		ocode    ocode.Ocode
		expected []byte
		wantErr  bool
	}{
		{
			name:     "RET instruction",
			ocode:    ocode.Ocode{Kind: ocode.OpRET, Operands: nil},
			expected: []byte{0xC3},
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := handleRET(tc.ocode)

			if tc.wantErr {
				assert.Error(t, err, "handleRET() should return an error")
			} else {
				assert.NoError(t, err, "handleRET() should not return an error")
				assert.Equal(t, tc.expected, result, "handleRET() should generate correct binary")
			}
		})
	}
}
