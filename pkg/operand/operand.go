 package operand

import "github.com/HobbyOSs/gosk/pkg/cpu"

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
	WithBitMode(mode cpu.BitMode) Operands // Add back
	WithForceImm8(force bool) Operands
	WithForceRelAsImm(force bool) Operands
	GetBitMode() cpu.BitMode // Add back
	Require66h() bool // オペランドサイズプレフィックスが必要かどうか
	Require67h() bool // アドレスサイズプレフィックスが必要かどうか
	ParsedOperands() []*ParsedOperand
}
