package ast

//go:generate newc
type PlusExp struct {
	BaseExp
	Exp1 Exp
	Exp2 Exp
}

func (plus PlusExp) String() string {
	return plus.Exp1.String() + " + " + plus.Exp2.String()
}

//go:generate newc
type ImmExp struct {
	BaseExp
	Factor Factor
}

func (imm ImmExp) String() string {
	return imm.Factor.String()
}
