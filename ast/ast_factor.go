package ast

import (
	"reflect"
)

type Factor interface {
	Type() string
	String() string
}

type BaseFactor struct{}

func (b BaseFactor) Type() string {
	return reflect.TypeOf(b).Elem().Name()
}
