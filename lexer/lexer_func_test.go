package lexer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsDigit(t *testing.T) {
	assert.Equal(t, isDigit('あ'), false)
	assert.Equal(t, isDigit('0'), true)
}

func TestIsLetter(t *testing.T) {
	assert.Equal(t, isLetter('あ'), false)
	assert.Equal(t, isLetter('0'), true)
	assert.Equal(t, isLetter('a'), true)
	assert.Equal(t, isLetter('_'), true)
}
