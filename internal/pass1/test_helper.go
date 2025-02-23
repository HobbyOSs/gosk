package pass1

import (
	"log"
	"strings"

	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/comail/colog"
	"github.com/zeroflucs-given/generics/collections/stack"
)

func setUpColog(logLevel colog.Level) {
	colog.Register()
	colog.SetDefaultLevel(colog.LInfo)
	colog.SetMinLevel(logLevel)
	colog.SetFlags(log.Lshortfile)
	colog.SetFormatter(&colog.StdFormatter{Colors: false})
}

func buildImmExpFromValue(value any) *ast.ImmExp {
	var factor ast.Factor
	switch v := value.(type) {
	case int:
		factor = &ast.NumberFactor{BaseFactor: ast.BaseFactor{}, Value: v}
	case string:
		if strings.HasPrefix(v, "0x") {
			factor = &ast.HexFactor{BaseFactor: ast.BaseFactor{}, Value: v}
		} else if strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'") {
			factor = &ast.CharFactor{BaseFactor: ast.BaseFactor{}, Value: v[1 : len(v)-1]}
		} else if strings.HasPrefix(v, `"`) && strings.HasSuffix(v, `"`) {
			factor = &ast.StringFactor{BaseFactor: ast.BaseFactor{}, Value: v[1 : len(v)-1]}
		} else {
			factor = &ast.IdentFactor{BaseFactor: ast.BaseFactor{}, Value: v}
		}
	}

	return &ast.ImmExp{Factor: factor}
}

func buildStack(tokens []*token.ParseToken) *stack.Stack[*token.ParseToken] {
	stack := stack.NewStack[*token.ParseToken](10)
	for _, t := range tokens {
		stack.Push(t)
	}
	return stack
}
