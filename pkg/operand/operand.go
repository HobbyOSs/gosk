package operand

type SegmentedReg struct {
	Seg   string `@Seg`
	Colon string `@Colon`
	Reg   string `@Reg`
}

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

type Operands interface {
	InternalString() string
	OperandTypes() []OperandType
	Serialize() string
	FromString(text string) Operands
}

type OperandBuilder struct{}
