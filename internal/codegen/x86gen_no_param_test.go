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
