package ast

//go:generate newc
type DeclareStmt struct {
	BaseStatement
	IdentFactor *IdentFactor
	Exp         Exp // interfaceはポインタにしない
}

func (d DeclareStmt) String() string {
	return d.IdentFactor.String() + " EQU " + d.Exp.String()
}

//go:generate newc
type LabelStmt struct {
	BaseStatement
	IdentFactor *IdentFactor
}

func (l LabelStmt) String() string {
	return l.IdentFactor.String()
}
