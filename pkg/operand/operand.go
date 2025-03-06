package operand

import "github.com/HobbyOSs/gosk/internal/ast"

func equalOperandTypes(a, b []OperandType) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type Operands interface {
	InternalString() string
	OperandTypes() []OperandType
	Serialize() string
	FromString(text string) Operands
	CalcOffsetByteSize() int
	DetectImmediateSize() int
	WithBitMode(mode ast.BitMode) Operands
	GetBitMode() ast.BitMode
}
