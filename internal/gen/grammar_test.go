package gen

import (
	"strings"
	"testing"

	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func buildImmExpFromValue(value any) *ast.ImmExp {
	var factor ast.Factor
	switch v := value.(type) {
	case int:
		factor = &ast.NumberFactor{BaseFactor: ast.BaseFactor{}, Value: v}
	case string:
		if strings.HasPrefix(v, "0x") {
			factor = &ast.HexFactor{BaseFactor: ast.BaseFactor{}, Value: v}
		} else {
			factor = &ast.IdentFactor{BaseFactor: ast.BaseFactor{}, Value: v}
		}
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

func buildSegmentExpFromValue(value any) *ast.SegmentExp {
	return &ast.SegmentExp{
		DataType: "",
		Left:     buildAddExpFromValue(value),
		Right:    nil,
	}
}

func buildMemoryAddrExpFromValue(left any, right any) *ast.MemoryAddrExp {
	return &ast.MemoryAddrExp{
		DataType: "",
		Left:     buildAddExpFromValue(left),
		Right:    buildAddExpFromValue(right),
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
		{"line comment1", "Comment", "# sample \n", "# sample \n"},
		{"line comment2", "Comment", "; sample \n", "; sample \n"},
		{"line comment1", "Comment", "# sample", "# sample"},
		{"line comment2", "Comment", "; sample", "; sample"},
		// exp
		{"simple exp1", "Exp", "10",
			buildSegmentExpFromValue(10),
		},
		{"simple exp2", "Exp", "CYLS",
			buildSegmentExpFromValue("CYLS"),
		},
		{"complex exp1", "Exp", "DWORD 2*8:0x0000001b",
			&ast.SegmentExp{
				DataType: ast.Dword,
				Left: &ast.AddExp{
					HeadExp: &ast.MultExp{
						HeadExp:   buildImmExpFromValue(2),
						Operators: []string{"*"},
						TailExps:  []*ast.ImmExp{buildImmExpFromValue(8)},
					},
					Operators: []string{},
					TailExps:  []*ast.MultExp{},
				},
				Right: buildAddExpFromValue("0x0000001b"),
			},
		},
		{"complex exp2", "Exp", "512*18*2/4",
			&ast.SegmentExp{
				BaseExp:  ast.BaseExp{},
				DataType: ast.None,
				Left: &ast.AddExp{
					HeadExp: &ast.MultExp{
						HeadExp:   buildImmExpFromValue(512),
						Operators: []string{"*", "*", "/"},
						TailExps: []*ast.ImmExp{
							buildImmExpFromValue(18),
							buildImmExpFromValue(2),
							buildImmExpFromValue(4),
						},
					},
					Operators: []string{},
					TailExps:  []*ast.MultExp{},
				},
				Right: nil,
			},
		},
		{"memory address direct", "Exp", "[100]",
			&ast.MemoryAddrExp{
				DataType: ast.None,
				Left:     buildAddExpFromValue(100),
				Right:    nil,
			},
		},
		{"memory address direct (complex)", "Exp", "[CS:0x0020]",
			&ast.MemoryAddrExp{
				DataType: ast.None,
				Left:     buildAddExpFromValue("CS"),
				Right:    buildAddExpFromValue("0x0020"),
			},
		},
		{"memory address register indirect", "Exp", "[BX]",
			&ast.MemoryAddrExp{
				DataType: ast.None,
				Left:     buildAddExpFromValue("BX"),
				Right:    nil,
			},
		},
		{"memory address register indirect (complex)", "Exp", "[CS:ECX]",
			&ast.MemoryAddrExp{
				DataType: ast.None,
				Left:     buildAddExpFromValue("CS"),
				Right:    buildAddExpFromValue("ECX"),
			},
		},
		{"memory address based", "Exp", "[ESP+12]",
			&ast.MemoryAddrExp{
				DataType: ast.None,
				Left: &ast.AddExp{
					HeadExp:   buildMultExpFromValue("ESP"),
					Operators: []string{"+"},
					TailExps:  []*ast.MultExp{buildMultExpFromValue(12)},
				},
				Right: nil,
			},
		},

		// stmt
		{"equ macro", "DeclareStmt", "CYLS EQU 10",
			ast.NewDeclareStmt(
				ast.BaseStatement{},
				ast.NewIdentFactor(ast.BaseFactor{}, "CYLS"),
				buildSegmentExpFromValue(10),
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
				&ast.NumberFactor{BaseFactor: ast.BaseFactor{}, Value: 32},
			),
		},
		{"opcode only", "OpcodeStmt", "HLT",
			ast.NewOpcodeStmt(
				ast.BaseStatement{},
				ast.NewIdentFactor(ast.BaseFactor{}, "HLT"),
			),
		},
		{"1 operand_1", "MnemonicStmt", " ORG 0x7c00 ; comment",
			ast.NewMnemonicStmt(
				ast.BaseStatement{},
				ast.NewIdentFactor(ast.BaseFactor{}, "ORG"),
				[]ast.Exp{
					buildSegmentExpFromValue("0x7c00"),
				},
			),
		},
		{"1 operand_2", "MnemonicStmt", " JMP fin ; comment",
			ast.NewMnemonicStmt(
				ast.BaseStatement{},
				ast.NewIdentFactor(ast.BaseFactor{}, "JMP"),
				[]ast.Exp{
					buildSegmentExpFromValue("fin"),
				},
			),
		},
		{"1 operand_3", "MnemonicStmt", "RESB 0x7dfe-$",
			ast.NewMnemonicStmt(
				ast.BaseStatement{},
				ast.NewIdentFactor(ast.BaseFactor{}, "RESB"),
				[]ast.Exp{
					&ast.SegmentExp{
						DataType: "",
						Left: &ast.AddExp{
							HeadExp:   buildMultExpFromValue("0x7dfe"),
							Operators: []string{"-"},
							TailExps: []*ast.MultExp{
								buildMultExpFromValue("$"),
							},
						},
						Right: nil,
					},
				},
			),
		},
		{"opcode simple mnemonic", "MnemonicStmt", "DB 10,20,30",
			ast.NewMnemonicStmt(
				ast.BaseStatement{},
				ast.NewIdentFactor(ast.BaseFactor{}, "DB"),
				[]ast.Exp{
					buildSegmentExpFromValue(10),
					buildSegmentExpFromValue(20),
					buildSegmentExpFromValue(30),
				},
			),
		},
		// program
		{"1 operand program", "Program", "ORG 0x7c00 ; comment",
			&ast.Program{
				Statements: []ast.Statement{
					ast.NewMnemonicStmt(
						ast.BaseStatement{},
						ast.NewIdentFactor(ast.BaseFactor{}, "ORG"),
						[]ast.Exp{
							buildSegmentExpFromValue("0x7c00"),
						},
					),
				},
			},
		},
		{"cfg program1", "Program", "MOV [CS:DS],8 ; comment",
			&ast.Program{
				Statements: []ast.Statement{
					ast.NewMnemonicStmt(
						ast.BaseStatement{},
						ast.NewIdentFactor(ast.BaseFactor{}, "MOV"),
						[]ast.Exp{
							buildMemoryAddrExpFromValue("CS", "DS"),
							buildSegmentExpFromValue(8),
						},
					),
				},
			},
		},
		{"cfg program2", "Program", "MOV DWORD [VRAM],0x000a0000 ; comment",
			&ast.Program{
				Statements: []ast.Statement{
					ast.NewMnemonicStmt(
						ast.BaseStatement{},
						ast.NewIdentFactor(ast.BaseFactor{}, "MOV"),
						[]ast.Exp{
							&ast.MemoryAddrExp{
								DataType: ast.Dword,
								Left:     buildAddExpFromValue("VRAM"),
								Right:    nil,
							},
							buildSegmentExpFromValue("0x000a0000"),
						},
					),
				},
			},
		},
		{"cfg program3", "Program", "HLT ;\n JMP fin",
			&ast.Program{
				Statements: []ast.Statement{
					ast.NewOpcodeStmt(ast.BaseStatement{}, ast.NewIdentFactor(ast.BaseFactor{}, "HLT")),
					ast.NewMnemonicStmt(
						ast.BaseStatement{},
						ast.NewIdentFactor(ast.BaseFactor{}, "JMP"),
						[]ast.Exp{
							buildSegmentExpFromValue("fin"),
						},
					),
				},
			},
		},
		{"cfg program4", "Program", "_io_hlt:	;\n",
			&ast.Program{
				Statements: []ast.Statement{
					ast.NewLabelStmt(
						ast.BaseStatement{},
						ast.NewIdentFactor(ast.BaseFactor{}, "_io_hlt:"),
					),
				},
			},
		},
		{"cfg program5", "Program", `_farjmp: ;
		JMP FAR [ESP+4] ; eip, cs`,
			&ast.Program{
				Statements: []ast.Statement{
					ast.NewLabelStmt(
						ast.BaseStatement{},
						ast.NewIdentFactor(ast.BaseFactor{}, "_farjmp:"),
					),
					ast.NewMnemonicStmt(
						ast.BaseStatement{},
						ast.NewIdentFactor(ast.BaseFactor{}, "JMP"),
						[]ast.Exp{
							&ast.MemoryAddrExp{
								DataType: "",
								JumpType: "FAR",
								Left: &ast.AddExp{
									HeadExp:   buildMultExpFromValue("ESP"),
									Operators: []string{"+"},
									TailExps:  []*ast.MultExp{buildMultExpFromValue(4)},
								},
								Right: nil,
							},
						},
					),
				},
			},
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
