package ast

//go:generate newc
type NumberFactor struct {
	BaseFactor
	value int
}

func (n *NumberFactor) Value() int {
	return n.value
}

//go:generate newc
type StringFactor struct {
	BaseFactor
	value string
}

func (s *StringFactor) Value() string {
	return s.value
}

//go:generate newc
type HexFactor struct {
	BaseFactor
	value string
}

func (h *HexFactor) Value() string {
	return h.value
}

//go:generate newc
type IdentFactor struct {
	BaseFactor
	value string
}

func (i *IdentFactor) Value() string {
	return i.value
}

//go:generate newc
type CharFactor struct {
	BaseFactor
	value string
}

func (c *CharFactor) Value() string {
	return c.value
}
