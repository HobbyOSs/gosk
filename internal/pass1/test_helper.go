package pass1

import (
	"log"
	"strings"

	"github.com/HobbyOSs/gosk/internal/ast" // Added import
	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/comail/colog"
	"github.com/zeroflucs-given/generics/collections/stack"
)

func setUpColog(logLevel colog.Level) { // Keep logLevel param for potential future use, but ignore it for now
	colog.Register()
	colog.SetDefaultLevel(colog.LInfo) // Keep default level as Info
	colog.SetMinLevel(colog.LDebug)    // テスト時は常に Debug レベル
	colog.SetFlags(log.Lshortfile)
	colog.SetFormatter(&colog.StdFormatter{Colors: false})
}

func buildImmExpFromValue(value any) *ast.ImmExp { // Restored ast.ImmExp
	var factor ast.Factor // Restored ast.Factor
	switch v := value.(type) {
	case int:
		factor = &ast.NumberFactor{BaseFactor: ast.BaseFactor{}, Value: v} // Restored ast types
	case string:
		if strings.HasPrefix(v, "0x") {
			factor = &ast.HexFactor{BaseFactor: ast.BaseFactor{}, Value: v} // Restored ast types
		} else if strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'") {
			factor = &ast.CharFactor{BaseFactor: ast.BaseFactor{}, Value: v[1 : len(v)-1]} // Restored ast types
		} else if strings.HasPrefix(v, `"`) && strings.HasSuffix(v, `"`) {
			factor = &ast.StringFactor{BaseFactor: ast.BaseFactor{}, Value: v[1 : len(v)-1]} // Restored ast types
		} else {
			factor = &ast.IdentFactor{BaseFactor: ast.BaseFactor{}, Value: v} // Restored ast types
		}
	}

	return &ast.ImmExp{Factor: factor} // Restored ast.ImmExp
}

func buildStack(tokens []*token.ParseToken) *stack.Stack[*token.ParseToken] {
	stack := stack.NewStack[*token.ParseToken](10)
	for _, t := range tokens {
		stack.Push(t)
	}
	return stack
}
