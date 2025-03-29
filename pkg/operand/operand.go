package operand

import "github.com/HobbyOSs/gosk/pkg/cpu"

// Instruction 構造体を変更: 最初のオペランドと後続を分離
type Instruction struct {
	FirstOperand *ParsedOperand  `parser:"@@"`
	RestOperands []*CommaOperand `parser:"@@*"` // カンマと後続オペランドの繰り返し
}

// CommaOperand 構造体を追加
type CommaOperand struct {
	Comma   string         `parser:"@Comma"` // カンマトークンをキャプチャ
	Operand *ParsedOperand `parser:"@@"`
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
	Require66h() bool        // オペランドサイズプレフィックスが必要かどうか
	Require67h() bool        // アドレスサイズプレフィックスが必要かどうか
	ParsedOperands() []*ParsedOperand
}
