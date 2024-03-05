package ast

type Exp[T any] interface {
	Type() string
	Value() T
}
