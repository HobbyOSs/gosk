package operand

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

type ParsedOperand struct {
	SegMem    string `@SegMem`
	Reg       string `| @Reg`
	Addr      string `| @Addr`
	Mem       string `| @Mem`
	Imm       string `| @Imm`
	Seg       string `| @Seg`
	Rel       string `| @Rel`
	MemPrefix string `| @MemPrefix`
}

type Operands interface {
	InternalString() string
	OperandTypes() []OperandType
	Serialize() string
	FromString(text string) Operands
	CalcOffsetByteSize() int
	DetectImmediateSize() int
}
