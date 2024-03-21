package ast

import "reflect"

type Statement interface {
	Node
	Type() string
	statementNode()
}

type BaseStatement struct{}

func (b BaseStatement) Type() string {
	return reflect.TypeOf(b).Elem().Name()
}
