package ast

import "reflect"

type Exp interface {
	Node
	expressionNode()
	Type() string
}

type BaseExp struct{}

func (b BaseExp) Type() string {
	return reflect.TypeOf(b).Elem().Name()
}
