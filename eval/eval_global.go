package eval

import (
	"fmt"
	"github.com/hangingman/gosk/ast"
	"github.com/hangingman/gosk/object"
	"github.com/hangingman/gosk/token"
	"log"
	"strings"
)

func evalGLOBALStatement(stmt *ast.MnemonicStatement) object.Object {
	toks := []string{}
	bin := &object.Binary{Value: []byte{}}

	for _, tok := range stmt.Name.Tokens {
		if tok.Type == token.IDENT {
			globalSymbolList = append(globalSymbolList, tok.Literal)
		}
		toks = append(toks, fmt.Sprintf("%s: %s", tok.Type, tok.Literal))
	}

	log.Println(fmt.Sprintf("info: [%s]", strings.Join(toks, ", ")))
	stmt.Bin = bin
	return stmt.Bin
}
