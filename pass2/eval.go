package pass2

import (
	"log"

	"github.com/comail/colog"
	"github.com/HobbyOSs/gosk/ast"
)

func init() {
	colog.Register()
	colog.SetDefaultLevel(colog.LInfo)
	colog.SetMinLevel(colog.LInfo)
	colog.SetFlags(log.Lshortfile)
	colog.SetFormatter(&colog.StdFormatter{Colors: false})
}

func Eval(program ast.Prog) {
	log.Println("pass2")
}
