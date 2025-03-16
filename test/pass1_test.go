package test

import (
	"strings"
	"testing"

	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/HobbyOSs/gosk/internal/codegen"
	"github.com/HobbyOSs/gosk/internal/gen"
	client "github.com/HobbyOSs/gosk/internal/ocode_client"
	"github.com/HobbyOSs/gosk/internal/pass1"
	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zeroflucs-given/generics/collections/stack"
)

type Pass1Suite struct {
	suite.Suite
}

// go-cmpで比較できない要素をここに列挙する
var IgnoreFields = []string{"Ctx", "EquMap", "Client"}

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
		ctx  *stack.Stack[*token.ParseToken]
		equ  map[string]*token.ParseToken
		want *pass1.Pass1
	}{
		{
			"config",
			"[BITS 32]",
			stack.NewStack[*token.ParseToken](100),
			nil,
			&pass1.Pass1{
				LOC:              0,
				BitMode:          ast.MODE_32BIT,
				SymTable:         make(map[string]int32, 0),
				GlobalSymbolList: []string{},
				ExternSymbolList: []string{},
				Client: func() client.CodegenClient {
					cli, _ := client.NewCodegenClient(&codegen.CodeGenContext{DollarPosition: 0, BitMode: ast.MODE_32BIT}, nil)
					return cli
				}(),
			},
		},
		{
			"equ",
			"CYLS EQU 10",
			stack.NewStack[*token.ParseToken](100),
			map[string]*token.ParseToken{"CYLS": token.NewParseToken(token.TTNumber, buildImmExpFromValue(10))},
			&pass1.Pass1{
				LOC:              0,
				BitMode:          ast.MODE_16BIT,
				SymTable:         make(map[string]int32, 0),
				GlobalSymbolList: []string{},
				ExternSymbolList: []string{},
				Client: func() client.CodegenClient {
					cli, _ := client.NewCodegenClient(&codegen.CodeGenContext{DollarPosition: 0, BitMode: ast.MODE_16BIT}, nil)
					return cli
				}(),
			},
		},
		{
			"DB_1",
			"DB 0x90",
			stack.NewStack[*token.ParseToken](100),
			nil,
			&pass1.Pass1{
				LOC:              1,
				BitMode:          ast.MODE_16BIT,
				SymTable:         make(map[string]int32, 0),
				GlobalSymbolList: []string{},
				ExternSymbolList: []string{},
				Client: func() client.CodegenClient {
					cli, _ := client.NewCodegenClient(&codegen.CodeGenContext{DollarPosition: 0, BitMode: ast.MODE_16BIT}, nil)
					return cli
				}(),
			},
		},
		{
			"DB_2",
			`DB "HELLO-OS   "`,
			stack.NewStack[*token.ParseToken](100),
			nil,
			&pass1.Pass1{
				LOC:              11,
				BitMode:          ast.MODE_16BIT,
				SymTable:         make(map[string]int32, 0),
				GlobalSymbolList: []string{},
				ExternSymbolList: []string{},
				Client: func() client.CodegenClient {
					cli, _ := client.NewCodegenClient(&codegen.CodeGenContext{DollarPosition: 0, BitMode: ast.MODE_16BIT}, nil)
					return cli
				}(),
			},
		},
		{
			"ORG",
			`ORG 0x7c00`,
			stack.NewStack[*token.ParseToken](100),
			nil,
			&pass1.Pass1{
				LOC:              0x7c00,
				BitMode:          ast.MODE_16BIT,
				SymTable:         make(map[string]int32, 0),
				GlobalSymbolList: []string{},
				ExternSymbolList: []string{},
				Client: func() client.CodegenClient {
					cli, _ := client.NewCodegenClient(&codegen.CodeGenContext{DollarPosition: 0x7c00, BitMode: ast.MODE_16BIT}, nil)
					return cli
				}(),
			},
		},
		{
			"RESB_1",
			"RESB 18",
			stack.NewStack[*token.ParseToken](100),
			nil,
			&pass1.Pass1{
				LOC:              18,
				BitMode:          ast.MODE_16BIT,
				SymTable:         make(map[string]int32, 0),
				GlobalSymbolList: []string{},
				ExternSymbolList: []string{},
				Client: func() client.CodegenClient {
					cli, _ := client.NewCodegenClient(&codegen.CodeGenContext{DollarPosition: 0, BitMode: ast.MODE_16BIT}, nil)
					return cli
				}(),
			},
		},
		{
			"RESB_2",
			"RESB 0x7dfe-$",
			stack.NewStack[*token.ParseToken](100),
			nil,
			&pass1.Pass1{
				LOC:              0x7dfe,
				BitMode:          ast.MODE_16BIT,
				SymTable:         make(map[string]int32, 0),
				GlobalSymbolList: []string{},
				ExternSymbolList: []string{},
				Client: func() client.CodegenClient {
					cli, _ := client.NewCodegenClient(&codegen.CodeGenContext{DollarPosition: 0, BitMode: ast.MODE_16BIT}, nil)
					return cli
				}(),
			},
		},
		{
			"Label",
			`ORG 0x7c00
		                         label:
		                         # dummy`,
			stack.NewStack[*token.ParseToken](100),
			nil,
			&pass1.Pass1{
				LOC:              0x7c00,
				BitMode:          ast.MODE_16BIT,
				SymTable:         make(map[string]int32{"label": 0x7c00}),
				GlobalSymbolList: []string{},
				ExternSymbolList: []string{},
				Client: func() client.CodegenClient {
					cli, err := client.NewCodegenClient(&codegen.CodeGenContext{DollarPosition: 0x7c00, BitMode: ast.MODE_16BIT}, pass1)
					if err != nil {
						panic(err)
					}
					return cli
				}(),
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
			pass1 := &pass1.Pass1{
				LOC:              0,
				BitMode:          ast.MODE_16BIT,
				EquMap:           make(map[string]*token.ParseToken, 0),
				SymTable:         make(map[string]int32, 0),
				GlobalSymbolList: []string{},
				ExternSymbolList: []string{},
				Ctx:              stack.NewStack[*token.ParseToken](100),
				Client:           client.NewCodegenClient(&codegen.CodeGenContext{DollarPosition: 0x7c00, BitMode: ast.MODE_16BIT}, pass1),
			}
			pass1.Eval(prog)
			if diff := cmp.Diff(*tt.want, *pass1, cmpopts.IgnoreFields(pass1.Pass1{}, "Ctx", "EquMap", "Client")); diff != "" {
				t.Errorf(`pass1.Eval("%v") result mismatch:\n%s`, prog, diff)
			}

			// Ctx: stack
			s.Require().Equal(tt.ctx.Capacity(), pass1.Ctx.Capacity(), "Should have same capacity")
			s.Require().Equal(tt.ctx.Count(), pass1.Ctx.Count(), "Should have same count")

			for i := tt.ctx.Count(); i >= 0; i-- {
				ex, _ := tt.ctx.Pop()
				ac, _ := pass1.Ctx.Pop()
				s.Require().Equal(ex, ac)
			}

			// Equ: map
			s.Require().Equal(len(tt.equ), len(pass1.EquMap), "Should have same count")
			for exK, exV := range tt.equ {
				s.Require().Equal(pass1.EquMap[exK], exV)
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
