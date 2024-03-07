package ast

import "reflect"

type Exp interface {
	Type() string
	String() string
}

type BaseExp struct{}

func (b BaseExp) Type() string {
	return reflect.TypeOf(b).Elem().Name()
}
