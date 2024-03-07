package ast

import "reflect"

type Statement interface {
	Type() string
	String() string
}

type BaseStatement struct{}

func (b BaseStatement) Type() string {
	return reflect.TypeOf(b).Elem().Name()
}
