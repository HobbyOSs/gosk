package ast

import (
	"reflect"
)

type Factor interface {
	Node
	Type() string
	factorNode()
}

type BaseFactor struct{}

func (b BaseFactor) Type() string {
	return reflect.TypeOf(b).Elem().Name()
}
