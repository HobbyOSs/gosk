package pass1

import "github.com/HobbyOSs/gosk/ast"

// 各フィールドを適切なハンドラーの新しいインスタンスで初期化
func NewVisitor(p *Pass1) *ast.Visitor {
	v := &ast.Visitor{}
	v.Handler(NewProgramHandlerImpl(v, p))
	v.Handler(NewDeclareStmtHandlerImpl(v, p))
	v.Handler(NewLabelStmtHandlerImpl(v, p))
	v.Handler(NewExportSymStmtHandlerImpl(v, p))
	v.Handler(NewExternSymStmtHandlerImpl(v, p))
	v.Handler(NewConfigStmtHandlerImpl(v, p))
	v.Handler(NewMnemonicStmtHandlerImpl(v, p))
	v.Handler(NewNumberFactorHandlerImpl(v, p))
	v.Handler(NewStringFactorHandlerImpl(v, p))
	v.Handler(NewHexFactorHandlerImpl(v, p))
	v.Handler(NewIdentFactorHandlerImpl(v, p))
	v.Handler(NewCharFactorHandlerImpl(v, p))

	return v
}
