package gen

import (
	"reflect"
	"testing"

	"github.com/hangingman/gosk/ast"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name       string
		entryPoint string
		text       string
		want       interface{}
	}{
		{"letter", "Letter", "a", []rune{'a'}},
		{"+int", "NumberFactor", "30", ast.NewNumberFactor(ast.BaseFactor{}, 30)},
		{"-int", "NumberFactor", "-30", ast.NewNumberFactor(ast.BaseFactor{}, -30)},
		{"hex", "HexFactor", "0x0ff0", ast.NewHexFactor(ast.BaseFactor{}, "0x0ff0")},
		{"char", "CharFactor", "'0x0ff0'", ast.NewCharFactor(ast.BaseFactor{}, "0x0ff0")},
		{"string", "StringFactor", "\"0x0ff0\"", ast.NewStringFactor(ast.BaseFactor{}, "0x0ff0")},
		{"ident", "IdentFactor", "_testZ009$", ast.NewIdentFactor(ast.BaseFactor{}, "_testZ009$")},
		{"label", "Label", "_test:\n", "_test:"},
		{"line comment1", "Comment", "# sample \n", ""},
		{"line comment2", "Comment", "; sample \n", ""},
		{"line comment1", "Comment", "# sample", ""},
		{"line comment2", "Comment", "; sample", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse("", []byte(tt.text), Entrypoint(tt.entryPoint))
			if !assert.NoError(t, err) {
				t.FailNow()
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse(%v) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}
