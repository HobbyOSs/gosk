package pass1

import (
	"fmt"
	"testing"

	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/HobbyOSs/gosk/internal/gen"
	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/comail/colog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zeroflucs-given/generics/collections/stack"
)

type Pass1HandlersSuite struct {
	suite.Suite
}

func TestPass1HandlersSuite(t *testing.T) {
	suite.Run(t, new(Pass1HandlersSuite))
}

func (s *Pass1HandlersSuite) SetupSuite() {
	setUpColog(colog.LDebug)
}

func (s *Pass1HandlersSuite) TestAddExp() {
	tests := []struct {
		name string
		text string
		want *stack.Stack[*token.ParseToken]
	}{
		{
			"+int",
			"30",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTNumber, buildImmExpFromValue(30)),
			}),
		},
		{
			"-int",
			"-30",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTNumber, buildImmExpFromValue(-30)),
			}),
		},
		{
			"hex",
			"0x0ff0",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTHex, buildImmExpFromValue("0x0ff0")),
			}),
		},
		{
			"char",
			"'0x0ff0'",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTIdentifier, buildImmExpFromValue("'0x0ff0'")),
			}),
		},
		{
			"string",
			`"0x0ff0"`,
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTIdentifier, buildImmExpFromValue(`"0x0ff0"`)),
			}),
		},
		{
			"ident",
			`_testZ009$`,
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTIdentifier, buildImmExpFromValue(`_testZ009$`)),
			}),
		},
		{
			"simple math 1",
			"1 + 1",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTNumber, buildImmExpFromValue(2)),
			}),
		},
		{
			"simple math 2",
			"4 - 2",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTNumber, buildImmExpFromValue(2)),
			}),
		},
		{
			"simple math 3",
			"1 + 3 - 2 + 4",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTNumber, buildImmExpFromValue(6)),
			}),
		},
		{
			"complex math 1",
			"8 * 3 - 1",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTNumber, buildImmExpFromValue(23)),
			}),
		},
		{
			"displacement 1",
			"ESP+4",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTIdentifier, buildImmExpFromValue("ESP")),
				token.NewParseToken(token.TTNumber, buildImmExpFromValue(4)),
			}),
		},
		{
			"displacement 2",
			"ESP+12+8",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTIdentifier, buildImmExpFromValue("ESP")),
				token.NewParseToken(token.TTNumber, buildImmExpFromValue(20)),
			}),
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := gen.Parse("", []byte(tt.text), gen.Entrypoint("AddExp"))
			if !assert.NoError(t, err) {
				t.FailNow()
			}
			node, ok := got.(*ast.AddExp)
			if !ok {
				t.FailNow()
			}

			p := &Pass1{
				Ctx: stack.NewStack[*token.ParseToken](10),
			}
			TraverseAST(node, p)

			for i := p.Ctx.Count(); i >= 0; i-- {
				_, expected := tt.want.Pop()
				_, actual := p.Ctx.Pop()
				assert.Equal(t, expected, actual,
					fmt.Sprintf("expected: %+v, actual: %+v\n", expected, actual))
			}
		})
	}
}

func (s *Pass1HandlersSuite) TestMultExp() {
	tests := []struct {
		name string
		text string
		want *stack.Stack[*token.ParseToken]
	}{
		{
			"simple math 1",
			"1005*8",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTNumber, buildImmExpFromValue(8040)),
			}),
		},
		{
			"simple math 2",
			"512/4",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTNumber, buildImmExpFromValue(128)),
			}),
		},
		{
			"simple math 3",
			"512*1024/4",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTNumber, buildImmExpFromValue(131072)),
			}),
		},
		{
			"scale 1",
			"EDX*4",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTIdentifier, buildImmExpFromValue("EDX")),
				token.NewParseToken(token.TTNumber, buildImmExpFromValue(4)),
			}),
		},
		{
			"scale 2",
			"ESI*8",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTIdentifier, buildImmExpFromValue("ESI")),
				token.NewParseToken(token.TTNumber, buildImmExpFromValue(8)),
			}),
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			got, err := gen.Parse("", []byte(tt.text), gen.Entrypoint("MultExp"))
			if !assert.NoError(t, err) {
				t.FailNow()
			}
			node, ok := got.(*ast.MultExp)
			if !ok {
				t.FailNow()
			}

			p := &Pass1{
				Ctx: stack.NewStack[*token.ParseToken](10),
			}
			TraverseAST(node, p)

			for i := p.Ctx.Count(); i >= 0; i-- {
				_, expected := tt.want.Pop()
				_, actual := p.Ctx.Pop()
				assert.Equal(t, expected, actual,
					fmt.Sprintf("expected: %+v, actual: %+v\n", expected, actual))
			}
		})
	}
}
