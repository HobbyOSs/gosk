package asmdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestX86Reference(t *testing.T) {
	ref := X86Reference()
	assert.NotNil(t, ref)
	t.Skip()
}
