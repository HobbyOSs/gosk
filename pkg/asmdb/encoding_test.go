package asmdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOutputSize(t *testing.T) {
	instrs := X86Instructions()

	// Test for a simple instruction
	if instr, exists := instrs["NOP"]; exists {
		assert.NotEmpty(t, instr.Forms)
		encoding := instr.Forms[0].Encodings[0]
		assert.Equal(t, 1, encoding.GetOutputSize(&OutputSizeOptions{}), "Expected output size for NOP is 1")
	}

	if instr, exists := instrs["CLI"]; exists {
		assert.NotEmpty(t, instr.Forms)
		encoding := instr.Forms[0].Encodings[0]
		assert.Equal(t, 1, encoding.GetOutputSize(&OutputSizeOptions{}), "Expected output size for CLI is 1")
	}
}
