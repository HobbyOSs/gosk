package ast

type BitMode int

const (
	ID_16BIT_MODE BitMode = 16
	ID_32BIT_MODE BitMode = 32
	ID_64BIT_MODE BitMode = 64
)

var intToBitMode = map[int]BitMode{
	16: ID_16BIT_MODE,
	32: ID_32BIT_MODE,
	64: ID_64BIT_MODE,
}

func NewBitMode(i int) (BitMode, bool) {
	b, ok := intToBitMode[i]
	return b, ok
}
