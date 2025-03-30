package test

import (
	"strings"
	"testing"

	"github.com/HobbyOSs/gosk/internal/ast" // Import ast package
	"github.com/HobbyOSs/gosk/internal/codegen"
	"github.com/HobbyOSs/gosk/internal/gen"
	ocode_client "github.com/HobbyOSs/gosk/internal/ocode_client"
	"github.com/HobbyOSs/gosk/internal/pass1"
	"github.com/HobbyOSs/gosk/pkg/asmdb" // Import asmdb
	"github.com/HobbyOSs/gosk/pkg/cpu"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Pass1Suite struct {
	suite.Suite
}

// go-cmpで比較できない要素をここに列挙する
var IgnoreFields = []string{"Client", "AsmDB"} // Removed Ctx, EquMap, added AsmDB

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

func (s *Pass1Suite) TestStatementToMachineCodeSize() {
	tests := []struct {
		name string
		text string
		want *pass1.Pass1
	}{
		{
			"config",
			"[BITS 32]",
			&pass1.Pass1{
				LOC:              0,
				BitMode:          cpu.MODE_32BIT, // Change cpu.MODE_32BIT to cpu.MODE_32BIT
				SymTable:         make(map[string]int32, 0),
				MacroMap:         make(map[string]ast.Exp), // Expect empty MacroMap
				GlobalSymbolList: []string{},
				ExternSymbolList: []string{},
			},
		},
		{
			"equ",
			"CYLS EQU 10",
			&pass1.Pass1{
				LOC:              0,
				BitMode:          cpu.MODE_16BIT, // Change cpu.MODE_16BIT to cpu.MODE_16BIT
				SymTable:         make(map[string]int32, 0),
				MacroMap:         map[string]ast.Exp{"CYLS": ast.NewNumberExp(ast.ImmExp{}, 10)}, // Expect MacroMap with evaluated value
				GlobalSymbolList: []string{},
				ExternSymbolList: []string{},
			},
		},
		{
			"DB_1",
			"DB 0x90",
			&pass1.Pass1{
				LOC:              1,
				BitMode:          cpu.MODE_16BIT, // Change cpu.MODE_16BIT to cpu.MODE_16BIT
				SymTable:         make(map[string]int32, 0),
				MacroMap:         make(map[string]ast.Exp),
				GlobalSymbolList: []string{},
				ExternSymbolList: []string{},
			},
		},
		{
			"DB_2",
			`DB "HELLO-OS   "`,
			&pass1.Pass1{
				LOC:              11,
				BitMode:          cpu.MODE_16BIT, // Change cpu.MODE_16BIT to cpu.MODE_16BIT
				SymTable:         make(map[string]int32, 0),
				MacroMap:         make(map[string]ast.Exp),
				GlobalSymbolList: []string{},
				ExternSymbolList: []string{},
			},
		},
		{
			"ORG",
			`ORG 0x7c00`,
			&pass1.Pass1{
				LOC:              0x7c00,
				DollarPosition:   0x7c00,
				BitMode:          cpu.MODE_16BIT, // Change cpu.MODE_16BIT to cpu.MODE_16BIT
				SymTable:         make(map[string]int32, 0),
				MacroMap:         make(map[string]ast.Exp),
				GlobalSymbolList: []string{},
				ExternSymbolList: []string{},
			},
		},
		{
			"RESB_1",
			"RESB 18",
			&pass1.Pass1{
				LOC:              18,
				DollarPosition:   0,
				BitMode:          cpu.MODE_16BIT, // Change cpu.MODE_16BIT to cpu.MODE_16BIT
				SymTable:         make(map[string]int32, 0),
				MacroMap:         make(map[string]ast.Exp),
				GlobalSymbolList: []string{},
				ExternSymbolList: []string{},
			},
		},
		{
			"RESB_2",
			"RESB 0x7dfe-$", // This requires '$' evaluation, which might not work correctly yet
			&pass1.Pass1{
				LOC:              0x7dfe,         // Assuming '$' evaluates correctly based on previous LOC
				DollarPosition:   0,              // Assuming ORG wasn't called before
				BitMode:          cpu.MODE_16BIT, // Change cpu.MODE_16BIT to cpu.MODE_16BIT
				SymTable:         make(map[string]int32, 0),
				MacroMap:         make(map[string]ast.Exp),
				GlobalSymbolList: []string{},
				ExternSymbolList: []string{},
			},
		},
		{
			"Label",
			`ORG 0x7c00
		                         label:
		                         # dummy`,
			&pass1.Pass1{
				LOC:              0x7c00,
				DollarPosition:   0x7c00,
				BitMode:          cpu.MODE_16BIT, // Change cpu.MODE_16BIT to cpu.MODE_16BIT
				SymTable:         map[string]int32{"label": 0x7c00},
				MacroMap:         make(map[string]ast.Exp),
				GlobalSymbolList: []string{},
				ExternSymbolList: []string{},
			},
		},
		{
			"integration test for pass1",
			`	ORG		0x7c00

				JMP		entry
				DB		0x90
				DB		"HELLOIPL"
				DW		512
				DB		1
				DW		1
				DB		2
				DW		224
				DW		2880
				DB		0xf0
				DW		9
				DW		18
				DW		2
				DD		0
				DD		2880
				DB		0,0,0x29
				DD		0xffffffff
				DB		"HELLO-OS   "
				DB		"FAT12   "
				RESB	18

		; プログラム本体

		entry:
				MOV		AX,0
				MOV		SS,AX
				MOV		SP,0x7c00
				MOV		DS,AX
				MOV		ES,AX
				MOV		SI,msg
		putloop:
				MOV		AL,[SI]
				ADD		SI,1
				CMP		AL,0
				JE		fin
				MOV		AH,0x0e
				MOV		BX,15
				INT		0x10
				JMP		putloop
		fin:
				HLT
				JMP		fin
		msg:
		`,
			&pass1.Pass1{
				LOC:            31860, // 0x7c00 + size of instructions/data
				DollarPosition: 0x7c00,
				BitMode:        cpu.MODE_16BIT, // Change cpu.MODE_16BIT to cpu.MODE_16BIT
				SymTable: map[string]int32{
					"entry":   31824, // 0x7c00 + size before entry
					"putloop": 31839, // entry + size before putloop
					"fin":     31857, // putloop + size before fin
					"msg":     31860, // fin + size before msg
				},
				MacroMap:         make(map[string]ast.Exp),
				GlobalSymbolList: []string{},
				ExternSymbolList: []string{},
			},
		},
	}

	t := s.T()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsedTree, err := gen.Parse("", []byte(tt.text), gen.Entrypoint("Program"))
			s.Require().NoError(err)
			prog, ok := (parsedTree).(ast.Prog)
			assert.True(t, ok)

			// pass1のEvalを実行
			ctx := &codegen.CodeGenContext{BitMode: cpu.MODE_16BIT} // Change cpu.MODE_16BIT to cpu.MODE_16BIT
			client, err := ocode_client.NewCodegenClient(ctx)
			s.Require().NoError(err)

			pass1 := &pass1.Pass1{
				LOC:              0,
				BitMode:          cpu.MODE_16BIT, // Change cpu.MODE_16BIT to cpu.MODE_16BIT
				SymTable:         make(map[string]int32, 0),
				MacroMap:         make(map[string]ast.Exp), // Initialize MacroMap
				NextImmJumpID:    0,
				GlobalSymbolList: []string{},
				ExternSymbolList: []string{},
				Client:           client,
				AsmDB:            asmdb.NewInstructionDB(), // Initialize AsmDB
			}
			pass1.Eval(prog)

			// Use IgnoreFields defined at the top
			if diff := cmp.Diff(*tt.want, *pass1, cmpopts.IgnoreFields(*pass1, IgnoreFields...)); diff != "" {
				t.Errorf(`pass1.Eval result mismatch (-want +got):\n%s`, diff)
			}

		})
	}

}

func TestPass1Suite(t *testing.T) {
	suite.Run(t, new(Pass1Suite))
}

func (s *Pass1Suite) SetupSuite() {
	setUpColog(true)
}

func (s *Pass1Suite) TearDownSuite() {
}

func (s *Pass1Suite) SetupTest() {
}

func (s *Pass1Suite) TearDownTest() {
}
