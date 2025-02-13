package operand

import (
	"fmt"

	participle "github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type ParsedOperand struct {
	Reg  string `@Reg`
	Mem  string `| @Mem`
	Imm  string `| @Imm`
	Seg  string `| @Seg`
	Rel  string `| @Rel`
	Addr string `| @Addr`
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
	{Name: "Reg", Pattern: `\b(R[ABCD]X|E[ABCD]X|AL|CL|DL|BL|AH|CH|DH|BH|MM[0-7]|XMM[0-9]|YMM[0-9])\b`},
	{Name: "Mem", Pattern: `\[(?:\w+|\d+|0x[a-fA-F0-9]+)\]`},
	{Name: "Imm", Pattern: `(?:-?\d+|0x[a-fA-F0-9]+)`},
	{Name: "Seg", Pattern: `\b(CS|DS|ES|FS|GS|SS)\b`},
	{Name: "Rel", Pattern: `\b(JMP|CALL|JNZ|JE|JNE|LABEL)\s+\w+\b`},
	{Name: "Addr", Pattern: `\[(?:0x[a-fA-F0-9]+|\d+)\]`},
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
