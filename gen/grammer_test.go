package gen

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hangingman/gosk/ast"
	"github.com/stretchr/testify/assert"
)

func buildImmExpFromValue(value any) *ast.ImmExp {

	var factor ast.Factor
	switch v := value.(type) {
	case int:
		factor = &ast.NumberFactor{ast.BaseFactor{}, v}
	case string:
		factor = &ast.IdentFactor{ast.BaseFactor{}, v}
	}

	return &ast.ImmExp{Factor: factor}
}

func buildMultExpFromValue(value any) *ast.MultExp {
	return &ast.MultExp{
		HeadExp:   buildImmExpFromValue(value),
		Operators: []string{},
		TailExps:  []*ast.ImmExp{},
	}
}

func buildAddExpFromValue(value any) *ast.AddExp {
	return &ast.AddExp{
		HeadExp:   buildMultExpFromValue(value),
		Operators: []string{},
		TailExps:  []*ast.MultExp{},
	}
}

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
		{"label with space", "Label", "_test:   \n", "_test:"},
		{"line comment1", "Comment", "# sample \n", ""},
		{"line comment2", "Comment", "; sample \n", ""},
		{"line comment1", "Comment", "# sample", ""},
		{"line comment2", "Comment", "; sample", ""},
		// exp
		// stmt
		{"equ macro", "DeclareStmt", "CYLS EQU 10",
			ast.NewDeclareStmt(
				ast.BaseStatement{},
				ast.NewIdentFactor(ast.BaseFactor{}, "CYLS"),
				buildAddExpFromValue(10),
			),
		},
		{"label", "LabelStmt", "_test:\n",
			ast.NewLabelStmt(
				ast.BaseStatement{},
				ast.NewIdentFactor(ast.BaseFactor{}, "_test:"),
			),
		},
		{"single symtable", "ExportSymStmt", "GLOBAL _io_hlt",
			ast.NewExportSymStmt(
				ast.BaseStatement{},
				[]*ast.IdentFactor{
					ast.NewIdentFactor(ast.BaseFactor{}, "_io_hlt"),
				},
			),
		},
		{"single export", "ExternSymStmt", "EXTERN _inthandler21",
			ast.NewExternSymStmt(
				ast.BaseStatement{},
				[]*ast.IdentFactor{
					ast.NewIdentFactor(ast.BaseFactor{}, "_inthandler21"),
				},
			),
		},
		{"multiple export", "ExternSymStmt", "EXTERN _inthandler21, _inthandler27, _inthandler2c",
			ast.NewExternSymStmt(
				ast.BaseStatement{},
				[]*ast.IdentFactor{
					ast.NewIdentFactor(ast.BaseFactor{}, "_inthandler21"),
					ast.NewIdentFactor(ast.BaseFactor{}, "_inthandler27"),
					ast.NewIdentFactor(ast.BaseFactor{}, "_inthandler2c"),
				},
			),
		},
		{"config1", "ConfigStmt", "[BITS 32]",
			ast.NewConfigStmt(
				ast.BaseStatement{},
				ast.Bits,
				&ast.NumberFactor{ast.BaseFactor{}, 32},
			),
		},

		{"opcode simple mnemonic", "MnemonicStmt", "DB 10,20,30",
			ast.NewMnemonicStmt(
				ast.BaseStatement{},
				ast.NewIdentFactor(ast.BaseFactor{}, "DB"),
				[]ast.Exp{
					buildAddExpFromValue(10),
					buildAddExpFromValue(20),
					buildAddExpFromValue(30),
				},
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
