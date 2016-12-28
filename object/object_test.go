package object

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBinaryInspect(t *testing.T) {
	bin := Binary{Value: []byte{1, 2, 3}}
	assert.Equal(t, bin.Value, []byte{1, 2, 3})
	assert.Equal(t, bin.Inspect(), "010203")
}
