package operand

import (
	"fmt"

	participle "github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

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
	AddressingType() AddressingType
	OperandType() OperandType
	Serialize() string
	FromString(text string) Operand
}

type OperandBuilder struct{}

var operandLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Colon", Pattern: `:`},
	{Name: "Seg", Pattern: `\b(CS|DS|ES|FS|GS|SS)\b`},
	{Name: "Reg", Pattern: `\b(R?[ABCD]X|E?[ABCD]X|[ABCD]L|[ABCD]H|SI|DI|MM[0-7]|XMM[0-9]|YMM[0-9]|TR[0-7]|CR[0-7]|DR[0-7])\b`},
	{Name: "Addr", Pattern: `\[(?: (?: (?:FAR\s+PTR|NEAR\s+PTR|PTR)\s+ )? \[0x[a-fA-F0-9]+\] )\]`},
	{Name: "Mem", Pattern: `\[(?:[A-Za-z_][A-Za-z0-9_]*|\w+\+\w+|\w+-\w+|0x[a-fA-F0-9]+|\d+)\]`},
	{Name: "Imm", Pattern: `(?:0x[a-fA-F0-9]+|-?\d+)`},
	{Name: "Rel", Pattern: `\b(JMP|CALL|JNZ|JE|JNE|LABEL)\s+\w+\b`},
	{Name: "String", Pattern: `"(?:\\.|[^"\\])*"`},
	{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
})

func getParser() *participle.Parser[ParsedOperand] {
	return participle.MustBuild[ParsedOperand](
		participle.Lexer(operandLexer),
		participle.Unquote(),
	)
}

func (OperandBuilder) Reg(name string) RegisterOperand {
	return RegisterOperand{reg: name}
}

func (OperandBuilder) Imm(value int) ImmediateOperand {
	return ImmediateOperand{value: value, internal: fmt.Sprintf("%d", value)}
}

func (OperandBuilder) Mem(base string, index string, scale int, displacement int) MemoryOperand {
	return MemoryOperand{base: base, index: index, scale: scale, displacement: displacement, internal: fmt.Sprintf("[%s+%s*%d+%d]", base, index, scale, displacement)}
}
