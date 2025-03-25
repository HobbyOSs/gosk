package operand

import "github.com/HobbyOSs/gosk/internal/ast"

type Instruction struct {
	Operands []*ParsedOperand `parser:"@@(',' @@)*"`
}

type ParsedOperand struct {
	SegMem      string       `parser:"@SegMem"`
	Reg         string       `parser:"| @Reg"`
	DirectMem   *DirectMem   `parser:"| @@"`
	IndirectMem *IndirectMem `parser:"| @@"`
	Imm         string       `parser:"| @Imm"`
	Seg         string       `parser:"| @Seg"`
	Rel         string       `parser:"| @Rel"`
}

type IndirectMem struct {
	Prefix *string `parser:"@MemSizePrefix?"`
	Mem    string  `parser:"@IndirectMem"`
}

type DirectMem struct {
	Prefix *string `parser:"@MemSizePrefix?"`
	Addr   string  `parser:"@DirectMem"`
}

type Operands interface {
	InternalString() string
	InternalStrings() []string
	OperandTypes() []OperandType
	Serialize() string
	FromString(text string) Operands
	CalcOffsetByteSize() int
	DetectImmediateSize() int
	WithBitMode(mode ast.BitMode) Operands
	WithForceImm8(force bool) Operands
	WithForceRelAsImm(force bool) Operands
	GetBitMode() ast.BitMode
	Require66h() bool // オペランドサイズプレフィックスが必要かどうか
	Require67h() bool // アドレスサイズプレフィックスが必要かどうか
	ParsedOperands() []*ParsedOperand
}
