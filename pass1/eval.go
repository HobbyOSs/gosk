package pass1

import (
	"log"

	"github.com/comail/colog"
	"github.com/HobbyOSs/gosk/ast"
)

//go:generate newc
type Pass1 struct {
	// LOC(location of counter)
	LOC int32
	// Pass1のシンボルテーブル
	SymTable         map[string]uint32
	GlobalSymbolList []string
	ExternSymbolList []string
}

func init() {
	colog.Register()
	colog.SetDefaultLevel(colog.LInfo)
	colog.SetMinLevel(colog.LInfo)
	colog.SetFlags(log.Lshortfile)
	colog.SetFormatter(&colog.StdFormatter{Colors: false})
}

func (p *Pass1) Eval(program ast.Prog) {
	log.Println("pass1")
}
