package eval

import (
	"github.com/hangingman/gosk/ast"
	"github.com/hangingman/gosk/object"
	"github.com/hangingman/gosk/token"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestEvalStatements(t *testing.T) {
	d1 := &ast.MnemonicStatement{
		Token: token.Token{Type: token.OPCODE, Literal: "DB"},
		Name: &ast.IdentifierArray{
			Tokens: []token.Token{
				{Type: token.OPCODE, Literal: "DB"},
				{Type: token.HEX_LIT, Literal: "0x01"},
				{Type: token.HEX_LIT, Literal: "0x02"},
				{Type: token.HEX_LIT, Literal: "0x03"},
			},
			Values: []string{
				"DB", "0x01", "0x02", "0x03",
			},
		},
	}
	dArr := []ast.Statement{d1, d1, d1}

	evaluated := evalStatements(dArr)
	assert.Equal(t, "*object.ObjectArray", reflect.TypeOf(evaluated).String())
	objArr, ok := evaluated.(*object.ObjectArray)

	assert.True(t, ok)
	assert.Equal(t, 3, len(*objArr))
	for _, obj := range *objArr {
		assert.Equal(t, "*object.Binary", reflect.TypeOf(obj).String())
		bin, _ := obj.(*object.Binary)
		assert.Equal(t, "010203", bin.Inspect())
	}
}
