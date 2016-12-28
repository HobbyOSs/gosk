package object

import (
	"encoding/hex"
)

type ObjectType string

const (
	OBJ_ARRAY  = "OBJECT_ARRAY"
	BINARY_OBJ = "BINARY"
	RECALC_OBJ = "RECALC"
)

type Object interface {
	Type() ObjectType
	Inspect() string
	// GetID() string
}

type ObjectArray []Object

func (oa *ObjectArray) Type() ObjectType {
	return OBJ_ARRAY
}

func (oa *ObjectArray) Inspect() string {
	result := "[ "
	for _, obj := range *oa {
		result += obj.Inspect()
		result += ", "
	}
	result += " ]"
	return result
}

// Binary
type Binary struct {
	Value []byte
}

func (b *Binary) Type() ObjectType {
	return BINARY_OBJ
}

func (b *Binary) Inspect() string {
	return hex.EncodeToString(b.Value)
}

// Recalc
type Recalc struct {
	Value func() []byte
}

func (r *Recalc) Type() ObjectType {
	return RECALC_OBJ
}

func (r *Recalc) Inspect() string {
	return hex.EncodeToString(r.Value())
}
