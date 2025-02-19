package operand

import (
	participle "github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type OperandImpl struct {
	Internal string
}

func (b *OperandImpl) OperandType() OperandType {
	parser := getParser()
	parsed, err := parser.ParseString("", b.Internal)
	if err == nil {
		switch {
		case parsed.Reg != "":
			return CodeDoubleword
		case parsed.Mem != "":
			return CodeDoubleword
		case parsed.Imm != "":
			return CodeDoublewordInteger
		case parsed.Seg != "":
			return CodeWord
		case parsed.Rel != "":
			return CodeWord
		case parsed.Addr != "":
			return CodeDoubleword
		}
	}
	return OperandType("unknown")
}

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
