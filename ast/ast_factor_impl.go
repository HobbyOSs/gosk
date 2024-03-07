package ast

import "fmt"

//go:generate newc
type NumberFactor struct {
	BaseFactor
	Value int
}

func (n NumberFactor) String() string {
	return fmt.Sprintf("%d", n.Value)
}

//go:generate newc
type StringFactor struct {
	BaseFactor
	Value string
}

func (s StringFactor) String() string {
	return s.Value
}

//go:generate newc
type HexFactor struct {
	BaseFactor
	Value string
}

func (h HexFactor) String() string {
	return h.Value
}

//go:generate newc
type IdentFactor struct {
	BaseFactor
	Value string
}

func (i IdentFactor) String() string {
	return fmt.Sprintf("%s", i.Value)
}

//go:generate newc
type CharFactor struct {
	BaseFactor
	Value string
}

func (c CharFactor) String() string {
	return fmt.Sprintf("%s", c.Value)
}
