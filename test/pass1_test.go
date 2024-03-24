package test

import (
	"testing"

	"github.com/HobbyOSs/gosk/ast"
	"github.com/HobbyOSs/gosk/gen"
	"github.com/HobbyOSs/gosk/pass1"
	"github.com/HobbyOSs/gosk/token"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zeroflucs-given/generics/collections/stack"
)

type Pass1Suite struct {
	suite.Suite
}

// mapやstackの内部はgo-cmpで比較できなかった
var IgnoreFields = []string{"Ctx", "EquMap"}

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
				BitMode:          pass1.ID_32BIT_MODE,
				SymTable:         make(map[string]uint32, 0),
				GlobalSymbolList: []string{},
				ExternSymbolList: []string{},
			},
		},
		{
			"equ",
			"CYLS EQU 10",
			stack.NewStack[*token.ParseToken](100),
			map[string]*token.ParseToken{"CYLS": token.NewParseToken(token.TTNumber, 10)},
			&pass1.Pass1{
				LOC:              0,
				BitMode:          pass1.ID_16BIT_MODE,
				SymTable:         make(map[string]uint32, 0),
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
			pass1 := &pass1.Pass1{
				LOC:              0,
				BitMode:          pass1.ID_16BIT_MODE,
				EquMap:           make(map[string]*token.ParseToken, 0),
				SymTable:         make(map[string]uint32, 0),
				GlobalSymbolList: []string{},
				ExternSymbolList: []string{},
				Ctx:              stack.NewStack[*token.ParseToken](100),
			}
			pass1.Eval(prog)
			if diff := cmp.Diff(*tt.want, *pass1, cmpopts.IgnoreFields(*pass1, IgnoreFields...)); diff != "" {
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
}

func (s *Pass1Suite) TearDownSuite() {
}

func (s *Pass1Suite) SetupTest() {
}

func (s *Pass1Suite) TearDownTest() {
}
