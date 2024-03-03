package ast

import "reflect"

type Factor[T any] interface {
    Type() string
    Value() T
}

type BaseFactor struct{}

func (b BaseFactor) Type() string {
    return reflect.TypeOf(b).Elem().Name()
}
