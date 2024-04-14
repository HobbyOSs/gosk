package asmdb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestX86Reference(t *testing.T) {
	ref := X86Reference()
	assert.NotNil(t, ref)
	assert.Len(t, ref.InstructionsBy("ADC"), 10)
	assert.Equal(t, "Add with Carry", ref.InstructionsBy("ADC")[0].Description)

	assert.Equal(t, "FA", ref.InstructionsBy("CLI")[0].Opcode)
	assert.Equal(t, "F4", ref.InstructionsBy("HLT")[0].Opcode)

	for _, i := range ref.InstructionsBy("MOV") {
		fmt.Printf("%+v\n", i)
		fmt.Printf("%+v\n", i.Operand1)
		fmt.Printf("%+v\n", i.Operand2)
	}

	assert.Len(t, ref.InstructionsBy("MOV"), 0)
}
