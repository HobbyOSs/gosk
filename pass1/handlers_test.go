package pass1

import (
	"fmt"
	"testing"

	"github.com/HobbyOSs/gosk/ast"
	"github.com/HobbyOSs/gosk/gen"
	"github.com/HobbyOSs/gosk/token"
	"github.com/stretchr/testify/assert"
	"github.com/zeroflucs-given/generics/collections/stack"
)

func buildStack(tokens []*token.ParseToken) *stack.Stack[*token.ParseToken] {
	stack := stack.NewStack[*token.ParseToken](10)
	for _, t := range tokens {
		stack.Push(t)
	}
	return stack
}

func TestAddExp(t *testing.T) {
	tests := []struct {
		name string
		text string
		want *stack.Stack[*token.ParseToken]
	}{
		{
			"+int",
			"30",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTNumber, 30),
			}),
		},
		{
			"-int",
			"-30",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTNumber, -30),
			}),
		},
		{
			"hex",
			"0x0ff0",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTHex, "0x0ff0"),
			}),
		},
		{
			"char",
			"'0x0ff0'",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTIdentifier, "0x0ff0"),
			}),
		},
		{
			"string",
			`"0x0ff0"`,
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTIdentifier, "0x0ff0"),
			}),
		},
		{
			"ident",
			`_testZ009$`,
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTIdentifier, `_testZ009$`),
			}),
		},
		{
			"simple math 1",
			"1 + 1",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTNumber, 2),
			}),
		},
		{
			"simple math 2",
			"4 - 2",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTNumber, 2),
			}),
		},
		{
			"simple math 3",
			"1 + 3 - 2 + 4",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTNumber, 6),
			}),
		},
		{
			"complex math 1",
			"8 * 3 - 1",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTNumber, 23),
			}),
		},
		{
			"displacement 1",
			"ESP+4",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTIdentifier, "ESP"),
				token.NewParseToken(token.TTNumber, 4),
			}),
		},
		{
			"displacement 2",
			"ESP+12+8",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTIdentifier, "ESP"),
				token.NewParseToken(token.TTNumber, 20),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			v := NewVisitor(p)
			handler := NewAddExpHandlerImpl(v, p)
			handler.AddExp(node)

			for i := p.Ctx.Count(); i >= 0; i-- {
				_, expected := tt.want.Pop()
				_, actual := p.Ctx.Pop()
				assert.Equal(t, expected, actual,
					fmt.Sprintf("expected: %+v, actual: %+v\n", expected, actual))
			}
		})
	}
}

func TestMultExp(t *testing.T) {
	tests := []struct {
		name string
		text string
		want *stack.Stack[*token.ParseToken]
	}{
		{
			"simple math 1",
			"1005*8",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTNumber, 8040),
			}),
		},
		{
			"simple math 2",
			"512/4",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTNumber, 128),
			}),
		},
		{
			"simple math 3",
			"512*1024/4",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTNumber, 131072),
			}),
		},
		{
			"scale 1",
			"EDX*4",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTIdentifier, "EDX"),
				token.NewParseToken(token.TTNumber, 4),
			}),
		},
		{
			"scale 2",
			"ESI*8",
			buildStack([]*token.ParseToken{
				token.NewParseToken(token.TTIdentifier, "ESI"),
				token.NewParseToken(token.TTNumber, 8),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
			v := NewVisitor(p)
			handler := NewMultExpHandlerImpl(v, p)
			handler.MultExp(node)

			for i := p.Ctx.Count(); i >= 0; i-- {
				_, expected := tt.want.Pop()
				_, actual := p.Ctx.Pop()
				assert.Equal(t, expected, actual,
					fmt.Sprintf("expected: %+v, actual: %+v\n", expected, actual))
			}
		})
	}
}
