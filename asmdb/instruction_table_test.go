package asmdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestX86Reference(t *testing.T) {
	ref := X86Reference()
	assert.NotNil(t, ref)
	assert.Len(t, ref.InstructionsBy("ADC"), 10)
	assert.Equal(t, "Add with Carry", ref.InstructionsBy("ADC")[0].Description)
}
