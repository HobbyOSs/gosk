package gen

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
		// xxx factor
		{"+int", "NumberFactor", "30", ast.NewNumberFactor(ast.BaseFactor{}, 30)},
		{"-int", "NumberFactor", "-30", ast.NewNumberFactor(ast.BaseFactor{}, -30)},
		{"hex", "HexFactor", "0x0ff0", ast.NewHexFactor(ast.BaseFactor{}, "0x0ff0")},
		{"char", "CharFactor", "'0x0ff0'", ast.NewCharFactor(ast.BaseFactor{}, "0x0ff0")},
		{"string", "StringFactor", "\"0x0ff0\"", ast.NewStringFactor(ast.BaseFactor{}, "0x0ff0")},
		{"ident", "IdentFactor", "_testZ009$", ast.NewIdentFactor(ast.BaseFactor{}, "_testZ009$")},
		// factor
		{"+int", "Factor", "30", ast.NewNumberFactor(ast.BaseFactor{}, 30)},
		{"-int", "Factor", "-30", ast.NewNumberFactor(ast.BaseFactor{}, -30)},
		{"hex", "Factor", "0x0ff0", ast.NewHexFactor(ast.BaseFactor{}, "0x0ff0")},
		{"char", "Factor", "'0x0ff0'", ast.NewCharFactor(ast.BaseFactor{}, "0x0ff0")},
		{"string", "Factor", "\"0x0ff0\"", ast.NewStringFactor(ast.BaseFactor{}, "0x0ff0")},
		{"ident", "Factor", "_testZ009$", ast.NewIdentFactor(ast.BaseFactor{}, "_testZ009$")},
		{"label", "Label", "_test:\n", "_test:"},
		{"line comment1", "Comment", "# sample \n", ""},
		{"line comment2", "Comment", "; sample \n", ""},
		{"line comment1", "Comment", "# sample", ""},
		{"line comment2", "Comment", "; sample", ""},
		// stmt
		{"equ macro", "DeclareStmt", "CYLS EQU 10",
			ast.NewDeclareStmt(
				ast.BaseStatement{},
				ast.NewIdentFactor(ast.BaseFactor{}, "CYLS"),
				ast.NewImmExp(ast.BaseExp{}, ast.NewNumberFactor(ast.BaseFactor{}, 10)),
			),
		},
		{"label1", "LabelStmt", "_test:\n",
			ast.NewLabelStmt(
				ast.BaseStatement{},
				ast.NewIdentFactor(ast.BaseFactor{}, "_test:"),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse("", []byte(tt.text), Entrypoint(tt.entryPoint))
			if !assert.NoError(t, err) {
				t.FailNow()
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf(`Parse("%v") result mismatch:\n%s`, tt.text, diff)
			}
		})
	}
}
