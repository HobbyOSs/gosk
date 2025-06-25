package ast

import "fmt"

//go:generate go tool newc
type NumberFactor struct {
	BaseFactor
	Value int
}

func (n NumberFactor) factorNode() {}
func (n NumberFactor) TokenLiteral() string {
	return fmt.Sprintf("%d", n.Value)
}

//go:generate go tool newc
type StringFactor struct {
	BaseFactor
	Value string
}

func (s StringFactor) factorNode() {}
func (s StringFactor) TokenLiteral() string {
	return s.Value
}

//go:generate go tool newc
type HexFactor struct {
	BaseFactor
	Value string
}

func (h HexFactor) factorNode() {}
func (h HexFactor) TokenLiteral() string {
	return h.Value
}

//go:generate go tool newc
type IdentFactor struct {
	BaseFactor
	Value string
}

func (i IdentFactor) factorNode() {}
func (i IdentFactor) TokenLiteral() string {
	return fmt.Sprintf("%s", i.Value)
}

//go:generate go tool newc
type CharFactor struct {
	BaseFactor
	Value string
}

func (c CharFactor) factorNode() {}
func (c CharFactor) TokenLiteral() string {
	return fmt.Sprintf("%s", c.Value)
}

func FactorToString(f Factor) string {
	switch x := f.(type) {
	case *NumberFactor:
		return x.TokenLiteral()
	case *HexFactor:
		return x.TokenLiteral()
	case *IdentFactor:
		return x.TokenLiteral()
	case *StringFactor:
		return x.TokenLiteral()
	case *CharFactor:
		return x.TokenLiteral()
	default:
		return ""
	}
}
