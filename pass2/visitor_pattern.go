package pass2

import "github.com/HobbyOSs/gosk/ast"

func NewVisitor(p *Pass2) *ast.Visitor {
	v := &ast.Visitor{}
	// program
	v.Handler(NewProgramHandlerImpl(v, p))
	// stmt
	v.Handler(NewDeclareStmtHandlerImpl(v, p))
	v.Handler(NewLabelStmtHandlerImpl(v, p))
	v.Handler(NewExportSymStmtHandlerImpl(v, p))
	v.Handler(NewExternSymStmtHandlerImpl(v, p))
	v.Handler(NewConfigStmtHandlerImpl(v, p))
	v.Handler(NewMnemonicStmtHandlerImpl(v, p))
	// exp
	v.Handler(NewMemoryAddrExpHandlerImpl(v, p))
	v.Handler(NewSegmentExpHandlerImpl(v, p))
	v.Handler(NewAddExpHandlerImpl(v, p))
	v.Handler(NewMultExpHandlerImpl(v, p))
	v.Handler(NewImmExpHandlerImpl(v, p))
	// factor
	v.Handler(NewNumberFactorHandlerImpl(v, p))
	v.Handler(NewStringFactorHandlerImpl(v, p))
	v.Handler(NewHexFactorHandlerImpl(v, p))
	v.Handler(NewIdentFactorHandlerImpl(v, p))
	v.Handler(NewCharFactorHandlerImpl(v, p))

	return v
}
