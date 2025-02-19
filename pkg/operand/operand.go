package operand

type SegmentedReg struct {
	Seg   string `@Seg`
	Colon string `@Colon`
	Reg   string `@Reg`
}

type SegmentedMem struct {
	Seg   string `@Seg`
	Colon string `@Colon`
	Mem   string `@Mem`
}

type ParsedOperand struct {
	SegReg *SegmentedReg `@@`
	SegMem *SegmentedMem `| @@`
	Reg    string        `| @Reg`
	Addr   string        `| @Addr`
	Mem    string        `| @Mem`
	Imm    string        `| @Imm`
	Seg    string        `| @Seg`
	Rel    string        `| @Rel`
}

type Operand interface {
	InternalString() string
	OperandType() OperandType
	Serialize() string
	FromString(text string) Operand
}

type OperandBuilder struct{}
