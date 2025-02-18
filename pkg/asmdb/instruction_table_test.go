package asmdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestX86Instructions(t *testing.T) {
	instrs := X86Instructions()
	assert.NotNil(t, instrs)
	assert.NotEmpty(t, instrs)
	//assert.NotEmpty(t, instrs[0].Mnemonic)
}
