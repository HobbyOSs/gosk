package gen

import (
	"strings"
	"testing"

	"github.com/HobbyOSs/gosk/internal/ast" // Import ast package
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
		TailExps:  []ast.Exp{},
	}
}

func buildAddExpFromValue(value any) *ast.AddExp {
	return &ast.AddExp{
		HeadExp:   buildMultExpFromValue(value),
		Operators: []string{},
		TailExps:  []*ast.MultExp{},
	}
}

// buildAddExpFromValue は単純な値から AddExp を構築するヘルパー (既存)
// buildMultExpFromValue は単純な値から MultExp を構築するヘルパー (既存)
// buildImmExpFromValue は単純な値から ImmExp を構築するヘルパー (既存)

// buildMemoryAddrExpFromValue はメモリアドレス式を構築するヘルパー (既存)
func buildMemoryAddrExpFromValue(left any, right any) *ast.MemoryAddrExp {
	// Note: This helper might need adjustment if MemoryAddrExp structure changes significantly
	// For now, assume it correctly builds based on AddExp for left/right parts.
	var leftExp *ast.AddExp
	if left != nil {
		leftExp = buildAddExpFromValue(left)
	}
	var rightExp *ast.AddExp
	if right != nil {
		rightExp = buildAddExpFromValue(right)
	}

	return &ast.MemoryAddrExp{
		DataType: ast.None, // Default to None, specific tests can override
		Left:     leftExp,
		Right:    rightExp,
	}
}

// buildSegmentExp はセグメント式を構築するヘルパー (新規または修正)
// Note: This helper assumes a structure where Left and Right are AddExp.
// It's used for explicit segment expressions like "DWORD 2*8:0x..."
func buildSegmentExp(dataType ast.DataType, leftExp *ast.AddExp, rightExp *ast.AddExp) *ast.SegmentExp {
	return &ast.SegmentExp{
		DataType: dataType,
		Left:     leftExp,
		Right:    rightExp,
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
		{"simple exp1", "Exp", "10", // Expect AddExp for simple numbers
			buildAddExpFromValue(10),
		},
		{"simple exp2", "Exp", "CYLS", // Expect AddExp for simple identifiers
			buildAddExpFromValue("CYLS"),
		},
		{"complex exp1 (SegmentExp)", "Exp", "DWORD 2*8:0x0000001b", // Expect SegmentExp for explicit segment syntax
			buildSegmentExp(
				ast.Dword,
				&ast.AddExp{ // Left part of SegmentExp
					HeadExp: &ast.MultExp{
						HeadExp:   buildImmExpFromValue(2),
						Operators: []string{"*"},
						TailExps:  []ast.Exp{buildImmExpFromValue(8)}, // Correct type: []ast.Exp
					},
					Operators: []string{},
					TailExps:  []*ast.MultExp{},
				},
				buildAddExpFromValue("0x0000001b"), // Right part of SegmentExp
			),
		},
		{"complex exp2 (AddExp)", "Exp", "512*18*2/4", // Expect AddExp for complex arithmetic without segment/memory syntax
			&ast.AddExp{
				HeadExp: &ast.MultExp{
					HeadExp:   buildImmExpFromValue(512),
					Operators: []string{"*", "*", "/"},
					TailExps: []ast.Exp{ // Correct type: []ast.Exp
						buildImmExpFromValue(18),
						buildImmExpFromValue(2),
						buildImmExpFromValue(4),
					},
				},
				Operators: []string{},
				TailExps:  []*ast.MultExp{},
			},
		},
		{"memory address direct", "Exp", "[100]", // Expect MemoryAddrExp
			buildMemoryAddrExpFromValue(100, nil),
		},
		{"memory address direct (complex)", "Exp", "[CS:0x0020]", // Expect MemoryAddrExp
			buildMemoryAddrExpFromValue("CS", "0x0020"),
		},
		{"memory address register indirect", "Exp", "[BX]", // Expect MemoryAddrExp
			buildMemoryAddrExpFromValue("BX", nil),
		},
		{"memory address register indirect (complex)", "Exp", "[CS:ECX]", // Expect MemoryAddrExp
			buildMemoryAddrExpFromValue("CS", "ECX"),
		},
		{"memory address based", "Exp", "[ESP+12]", // Expect MemoryAddrExp
			&ast.MemoryAddrExp{
				DataType: ast.None,
				Left: &ast.AddExp{ // The expression inside [] is an AddExp
					HeadExp:   buildMultExpFromValue("ESP"),
					Operators: []string{"+"},
					TailExps:  []*ast.MultExp{buildMultExpFromValue(12)},
				},
				Right: nil,
			},
		},
		{ // Renamed from "jmp dword far" to be more specific about the AST type
			name:       "segment exp (dword far)",
			entryPoint: "Exp",
			text:       "DWORD 2*8:0x0000001b",
			want: buildSegmentExp( // Use the new helper
				ast.Dword,
				&ast.AddExp{
					HeadExp: &ast.MultExp{
						HeadExp:   buildImmExpFromValue(2),
						Operators: []string{"*"},
						TailExps:  []ast.Exp{buildImmExpFromValue(8)}, // Correct type: []ast.Exp
					},
					Operators: []string{},
					TailExps:  []*ast.MultExp{},
				},
				buildAddExpFromValue("0x0000001b"),
			),
		},

		// stmt
		{"equ macro", "DeclareStmt", "CYLS EQU 10",
			ast.NewDeclareStmt(
				ast.BaseStatement{},
				ast.NewIdentFactor(ast.BaseFactor{}, "CYLS"),
				buildAddExpFromValue(10), // Expect AddExp for the value
			),
		},
		{"label", "LabelStmt", "_test:\n", // LabelStmt remains the same
			ast.NewLabelStmt( // No changes needed here
				ast.BaseStatement{},
				ast.NewIdentFactor(ast.BaseFactor{}, "_test:"),
			),
		},
		{"single symtable", "ExportSymStmt", "GLOBAL _io_hlt", // ExportSymStmt remains the same
			ast.NewExportSymStmt( // No changes needed here
				ast.BaseStatement{},
				[]*ast.IdentFactor{
					ast.NewIdentFactor(ast.BaseFactor{}, "_io_hlt"),
				},
			),
		},
		{"single export", "ExternSymStmt", "EXTERN _inthandler21", // ExternSymStmt remains the same
			ast.NewExternSymStmt( // No changes needed here
				ast.BaseStatement{},
				[]*ast.IdentFactor{
					ast.NewIdentFactor(ast.BaseFactor{}, "_inthandler21"),
				},
			),
		},
		{"multiple export", "ExternSymStmt", "EXTERN _inthandler21, _inthandler27, _inthandler2c", // ExternSymStmt remains the same
			ast.NewExternSymStmt( // No changes needed here
				ast.BaseStatement{},
				[]*ast.IdentFactor{
					ast.NewIdentFactor(ast.BaseFactor{}, "_inthandler21"),
					ast.NewIdentFactor(ast.BaseFactor{}, "_inthandler27"),
					ast.NewIdentFactor(ast.BaseFactor{}, "_inthandler2c"),
				},
			),
		},
		{"config1", "ConfigStmt", "[BITS 32]", // ConfigStmt remains the same
			ast.NewConfigStmt( // No changes needed here
				ast.BaseStatement{},
				ast.Bits,
				&ast.NumberFactor{BaseFactor: ast.BaseFactor{}, Value: 32},
			),
		},
		{"opcode only", "OpcodeStmt", "HLT", // OpcodeStmt remains the same
			ast.NewOpcodeStmt( // No changes needed here
				ast.BaseStatement{},
				ast.NewIdentFactor(ast.BaseFactor{}, "HLT"),
			),
		},
		{"1 operand_1", "MnemonicStmt", " ORG 0x7c00 ; comment",
			ast.NewMnemonicStmt(
				ast.BaseStatement{},
				ast.NewIdentFactor(ast.BaseFactor{}, "ORG"),
				[]ast.Exp{
					buildAddExpFromValue("0x7c00"), // Expect AddExp for operand
				},
			),
		},
		{"1 operand_2", "MnemonicStmt", " JMP fin ; comment",
			ast.NewMnemonicStmt(
				ast.BaseStatement{},
				ast.NewIdentFactor(ast.BaseFactor{}, "JMP"),
				[]ast.Exp{
					buildAddExpFromValue("fin"), // Expect AddExp for operand
				},
			),
		},
		{"1 operand_3", "MnemonicStmt", "RESB 0x7dfe-$",
			ast.NewMnemonicStmt(
				ast.BaseStatement{},
				ast.NewIdentFactor(ast.BaseFactor{}, "RESB"),
				[]ast.Exp{
					// Expect AddExp for the complex expression operand
					&ast.AddExp{
						HeadExp:   buildMultExpFromValue("0x7dfe"),
						Operators: []string{"-"},
						TailExps: []*ast.MultExp{
							buildMultExpFromValue("$"),
						},
					},
				},
			),
		},
		{"opcode simple mnemonic", "MnemonicStmt", "DB 10,20,30",
			ast.NewMnemonicStmt(
				ast.BaseStatement{},
				ast.NewIdentFactor(ast.BaseFactor{}, "DB"),
				[]ast.Exp{
					buildAddExpFromValue(10), // Expect AddExp for operands
					buildAddExpFromValue(20),
					buildAddExpFromValue(30),
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
							buildAddExpFromValue("0x7c00"), // Expect AddExp
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
							buildMemoryAddrExpFromValue("CS", "DS"), // MemoryAddrExp
							buildAddExpFromValue(8),                 // AddExp
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
							&ast.MemoryAddrExp{ // MemoryAddrExp with DataType
								DataType: ast.Dword,
								Left:     buildAddExpFromValue("VRAM"),
								Right:    nil,
							},
							buildAddExpFromValue("0x000a0000"), // AddExp
						},
					),
				},
			},
		},
		{"cfg program3", "Program", "HLT ;\n JMP fin",
			&ast.Program{
				Statements: []ast.Statement{
					ast.NewOpcodeStmt(ast.BaseStatement{}, ast.NewIdentFactor(ast.BaseFactor{}, "HLT")), // OpcodeStmt
					ast.NewMnemonicStmt(
						ast.BaseStatement{},
						ast.NewIdentFactor(ast.BaseFactor{}, "JMP"),
						[]ast.Exp{
							buildAddExpFromValue("fin"), // AddExp
						},
					),
				},
			},
		},
		{"cfg program4", "Program", "_io_hlt:	;\n",
			&ast.Program{
				Statements: []ast.Statement{
					ast.NewLabelStmt( // LabelStmt
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
					ast.NewLabelStmt( // LabelStmt
						ast.BaseStatement{},
						ast.NewIdentFactor(ast.BaseFactor{}, "_farjmp:"),
					),
					ast.NewMnemonicStmt(
						ast.BaseStatement{},
						ast.NewIdentFactor(ast.BaseFactor{}, "JMP"),
						[]ast.Exp{
							&ast.MemoryAddrExp{ // MemoryAddrExp with JumpType
								DataType: ast.None, // Explicitly None if not specified
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
