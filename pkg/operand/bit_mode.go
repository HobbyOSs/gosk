package operand

type BitMode int

const (
	MODE_16BIT BitMode = 16
	MODE_32BIT BitMode = 32
)

var intToBitMode = map[int]BitMode{
	16: MODE_16BIT,
	32: MODE_32BIT,
}

func NewBitMode(i int) (BitMode, bool) {
	b, ok := intToBitMode[i]
	return b, ok
}
