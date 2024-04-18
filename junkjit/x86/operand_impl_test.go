package x86

import (
	"testing"

	"github.com/HobbyOSs/gosk/asmdb"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name  string
		text  string
		wantA asmdb.AddressingType
	}{
		{"+int", "30", asmdb.CodeImmediate},
		{"-int", "-30", asmdb.CodeImmediate},
		{"hex", "0x0ff0", asmdb.CodeImmediate},
		{"char", "'0x0ff0'", asmdb.CodeImmediate},
		{"string", "\"0x0ff0\"", asmdb.CodeImmediate},
		{"ident", "_testZ009$", asmdb.CodeImmediate},
		{"CR0", "CR0", asmdb.CodeCRField},
		{"CR8", "CR8", asmdb.CodeCRField},
		{"DR0", "DR0", asmdb.CodeDebugField},
		{"DR7", "DR7", asmdb.CodeDebugField},
		{"Sreg1", "CS", asmdb.CodeSregField},
		{"Sreg2", "DS", asmdb.CodeSregField},
		{"Sreg3", "ES", asmdb.CodeSregField},
		{"Sreg4", "SS", asmdb.CodeSregField},
		{"Sreg5", "FS", asmdb.CodeSregField},
		{"Sreg6", "GS", asmdb.CodeSregField},
		{"GR1", "AX", asmdb.CodeGeneralReg},
		{"GR2", "SI", asmdb.CodeGeneralReg},
		{"GR3", "AH", asmdb.CodeGeneralReg},
		{"GR4", "EAX", asmdb.CodeGeneralReg},
		{"GR5", "EBX", asmdb.CodeGeneralReg},
		{"GR6", "RAX", asmdb.CodeGeneralReg},

		//{"memory address direct", "Exp", "[100]",
		// 	&ast.MemoryAddrExp{
		// 		DataType: ast.None,
		// 		Left:     buildAddExpFromValue(100),
		// 		Right:    nil,
		// 	},
		//},
		//{"memory address direct (complex)", "Exp", "[CS:0x0020]",
		// 	&ast.MemoryAddrExp{
		// 		DataType: ast.None,
		// 		Left:     buildAddExpFromValue("CS"),
		// 		Right:    buildAddExpFromValue("0x0020"),
		// 	},
		//},
		//{"memory address register indirect", "Exp", "[BX]",
		// 	&ast.MemoryAddrExp{
		// 		DataType: ast.None,
		// 		Left:     buildAddExpFromValue("BX"),
		// 		Right:    nil,
		// 	},
		//},
		//{"memory address register indirect (complex)", "Exp", "[CS:ECX]",
		// 	&ast.MemoryAddrExp{
		// 		DataType: ast.None,
		// 		Left:     buildAddExpFromValue("CS"),
		// 		Right:    buildAddExpFromValue("ECX"),
		// 	},
		//},
		//{"memory address based", "Exp", "[ESP+12]",
		// 	&ast.MemoryAddrExp{
		// 		DataType: ast.None,
		// 		Left: &ast.AddExp{
		// 			HeadExp:   buildMultExpFromValue("ESP"),
		// 			Operators: []string{"+"},
		// 			TailExps:  []*ast.MultExp{buildMultExpFromValue(12)},
		// 		},
		// 		Right: nil,
		// 	},
		//},
		//
		//// stmt
		//{"equ macro", "DeclareStmt", "CYLS EQU 10",
		// 	ast.NewDeclareStmt(
		// 		ast.BaseStatement{},
		// 		ast.NewIdentFactor(ast.BaseFactor{}, "CYLS"),
		// 		buildSegmentExpFromValue(10),
		// 	),
		//},
		//{"label", "LabelStmt", "_test:\n",
		// 	ast.NewLabelStmt(
		// 		ast.BaseStatement{},
		// 		ast.NewIdentFactor(ast.BaseFactor{}, "_test:"),
		// 	),
		//},
		//{"single symtable", "ExportSymStmt", "GLOBAL _io_hlt",
		// 	ast.NewExportSymStmt(
		// 		ast.BaseStatement{},
		// 		[]*ast.IdentFactor{
		// 			ast.NewIdentFactor(ast.BaseFactor{}, "_io_hlt"),
		// 		},
		// 	),
		//},
		//{"single export", "ExternSymStmt", "EXTERN _inthandler21",
		// 	ast.NewExternSymStmt(
		// 		ast.BaseStatement{},
		// 		[]*ast.IdentFactor{
		// 			ast.NewIdentFactor(ast.BaseFactor{}, "_inthandler21"),
		// 		},
		// 	),
		//},
		//{"multiple export", "ExternSymStmt", "EXTERN _inthandler21, _inthandler27, _inthandler2c",
		// 	ast.NewExternSymStmt(
		// 		ast.BaseStatement{},
		// 		[]*ast.IdentFactor{
		// 			ast.NewIdentFactor(ast.BaseFactor{}, "_inthandler21"),
		// 			ast.NewIdentFactor(ast.BaseFactor{}, "_inthandler27"),
		// 			ast.NewIdentFactor(ast.BaseFactor{}, "_inthandler2c"),
		// 		},
		// 	),
		//},
		//{"config1", "ConfigStmt", "[BITS 32]",
		// 	ast.NewConfigStmt(
		// 		ast.BaseStatement{},
		// 		ast.Bits,
		// 		&ast.NumberFactor{ast.BaseFactor{}, 32},
		// 	),
		//},
		//{"opcode only", "OpcodeStmt", "HLT",
		// 	ast.NewMnemonicStmt(
		// 		ast.BaseStatement{},
		// 		ast.NewIdentFactor(ast.BaseFactor{}, "HLT"),
		// 		[]ast.Exp{},
		// 	),
		//},
		//{"1 operand_1", "MnemonicStmt", " ORG 0x7c00 ; comment",
		// 	ast.NewMnemonicStmt(
		// 		ast.BaseStatement{},
		// 		ast.NewIdentFactor(ast.BaseFactor{}, "ORG"),
		// 		[]ast.Exp{
		// 			buildSegmentExpFromValue("0x7c00"),
		// 		},
		// 	),
		//},
		//{"1 operand_2", "MnemonicStmt", " JMP fin ; comment",
		// 	ast.NewMnemonicStmt(
		// 		ast.BaseStatement{},
		// 		ast.NewIdentFactor(ast.BaseFactor{}, "JMP"),
		// 		[]ast.Exp{
		// 			buildSegmentExpFromValue("fin"),
		// 		},
		// 	),
		//},
		//{"1 operand_3", "MnemonicStmt", "RESB 0x7dfe-$",
		// 	ast.NewMnemonicStmt(
		// 		ast.BaseStatement{},
		// 		ast.NewIdentFactor(ast.BaseFactor{}, "RESB"),
		// 		[]ast.Exp{
		// 			&ast.SegmentExp{
		// 				DataType: "",
		// 				Left: &ast.AddExp{
		// 					HeadExp:   buildMultExpFromValue("0x7dfe"),
		// 					Operators: []string{"-"},
		// 					TailExps: []*ast.MultExp{
		// 						buildMultExpFromValue("$"),
		// 					},
		// 				},
		// 				Right: nil,
		// 			},
		// 		},
		// 	),
		//},
		//{"opcode simple mnemonic", "MnemonicStmt", "DB 10,20,30",
		// 	ast.NewMnemonicStmt(
		// 		ast.BaseStatement{},
		// 		ast.NewIdentFactor(ast.BaseFactor{}, "DB"),
		// 		[]ast.Exp{
		// 			buildSegmentExpFromValue(10),
		// 			buildSegmentExpFromValue(20),
		// 			buildSegmentExpFromValue(30),
		// 		},
		// 	),
		//},
		//// program
		//{"1 operand program", "Program", "ORG 0x7c00 ; comment",
		// 	&ast.Program{
		// 		Statements: []ast.Statement{
		// 			ast.NewMnemonicStmt(
		// 				ast.BaseStatement{},
		// 				ast.NewIdentFactor(ast.BaseFactor{}, "ORG"),
		// 				[]ast.Exp{
		// 					buildSegmentExpFromValue("0x7c00"),
		// 				},
		// 			),
		// 		},
		// 	},
		//},
		//{"cfg program1", "Program", "MOV [CS:DS],8 ; comment",
		// 	&ast.Program{
		// 		Statements: []ast.Statement{
		// 			ast.NewMnemonicStmt(
		// 				ast.BaseStatement{},
		// 				ast.NewIdentFactor(ast.BaseFactor{}, "MOV"),
		// 				[]ast.Exp{
		// 					buildMemoryAddrExpFromValue("CS", "DS"),
		// 					buildSegmentExpFromValue(8),
		// 				},
		// 			),
		// 		},
		// 	},
		//},
		//{"cfg program2", "Program", "MOV DWORD [VRAM],0x000a0000 ; comment",
		// 	&ast.Program{
		// 		Statements: []ast.Statement{
		// 			ast.NewMnemonicStmt(
		// 				ast.BaseStatement{},
		// 				ast.NewIdentFactor(ast.BaseFactor{}, "MOV"),
		// 				[]ast.Exp{
		// 					&ast.MemoryAddrExp{
		// 						DataType: ast.Dword,
		// 						Left:     buildAddExpFromValue("VRAM"),
		// 						Right:    nil,
		// 					},
		// 					buildSegmentExpFromValue("0x000a0000"),
		// 				},
		// 			),
		// 		},
		// 	},
		//},
		//{"cfg program3", "Program", "HLT ;\n JMP fin",
		// 	&ast.Program{
		// 		Statements: []ast.Statement{
		// 			ast.NewMnemonicStmt(
		// 				ast.BaseStatement{},
		// 				ast.NewIdentFactor(ast.BaseFactor{}, "HLT"),
		// 				[]ast.Exp{},
		// 			),
		// 			ast.NewMnemonicStmt(
		// 				ast.BaseStatement{},
		// 				ast.NewIdentFactor(ast.BaseFactor{}, "JMP"),
		// 				[]ast.Exp{
		// 					buildSegmentExpFromValue("fin"),
		// 				},
		// 			),
		// 		},
		// 	},
		//},
		//{"cfg program4", "Program", "_io_hlt:	;\n",
		// 	&ast.Program{
		// 		Statements: []ast.Statement{
		// 			ast.NewLabelStmt(
		// 				ast.BaseStatement{},
		// 				ast.NewIdentFactor(ast.BaseFactor{}, "_io_hlt:"),
		// 			),
		// 		},
		// 	},
		//},
		//{"cfg program5", "Program", `_farjmp: ;
		//JMP FAR [ESP+4] ; eip, cs`,
		// 	&ast.Program{
		// 		Statements: []ast.Statement{
		// 			ast.NewLabelStmt(
		// 				ast.BaseStatement{},
		// 				ast.NewIdentFactor(ast.BaseFactor{}, "_farjmp:"),
		// 			),
		// 			ast.NewMnemonicStmt(
		// 				ast.BaseStatement{},
		// 				ast.NewIdentFactor(ast.BaseFactor{}, "JMP"),
		// 				[]ast.Exp{
		// 					&ast.MemoryAddrExp{
		// 						DataType: "",
		// 						JumpType: "FAR",
		// 						Left: &ast.AddExp{
		// 							HeadExp:   buildMultExpFromValue("ESP"),
		// 							Operators: []string{"+"},
		// 							TailExps:  []*ast.MultExp{buildMultExpFromValue(4)},
		// 						},
		// 						Right: nil,
		// 					},
		// 				},
		// 			),
		// 		},
		// 	},
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewX86Operand(tt.text)
			if !assert.NoError(t, err) {
				t.FailNow()
			}

			if diff := cmp.Diff(tt.wantA, got.AddressingType()); diff != "" {
				t.Errorf(`AddressingType("%v") result mismatch:\n%s`, tt.text, diff)
			}
		})
	}
}
