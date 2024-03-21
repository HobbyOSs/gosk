package ast

type Prog interface {
	Node
	Type() string
	program()
}
