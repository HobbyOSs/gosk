package ast

import "fmt"

//go:generate newc
type NumberFactor struct {
	BaseFactor
	Value int
}

func (n NumberFactor) factorNode() {}
func (n NumberFactor) TokenLiteral() string {
	return fmt.Sprintf("%d", n.Value)
}

//go:generate newc
type StringFactor struct {
	BaseFactor
	Value string
}

func (s StringFactor) factorNode() {}
func (s StringFactor) TokenLiteral() string {
	return s.Value
}

//go:generate newc
type HexFactor struct {
	BaseFactor
	Value string
}

func (h HexFactor) factorNode() {}
func (h HexFactor) TokenLiteral() string {
	return h.Value
}

//go:generate newc
type IdentFactor struct {
	BaseFactor
	Value string
}

func (i IdentFactor) factorNode() {}
func (i IdentFactor) TokenLiteral() string {
	return fmt.Sprintf("%s", i.Value)
}

//go:generate newc
type CharFactor struct {
	BaseFactor
	Value string
}

func (c CharFactor) factorNode() {}
func (c CharFactor) TokenLiteral() string {
	return fmt.Sprintf("%s", c.Value)
}
