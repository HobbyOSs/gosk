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
			// Determine register size based on name
			switch {
			case len(parsed.Reg) >= 3 && parsed.Reg[0] == 'R': // RAX, RBX etc
				return CodeR32
			case len(parsed.Reg) >= 3 && parsed.Reg[0] == 'E': // EAX, EBX etc
				return CodeR32
			case len(parsed.Reg) == 2 && parsed.Reg[1] == 'L': // AL, BL etc
				return CodeR8
			case len(parsed.Reg) == 2 && parsed.Reg[1] == 'X': // AX, BX etc
				return CodeR16
			case len(parsed.Reg) >= 3 && parsed.Reg[:3] == "XMM":
				return CodeXMM
			case len(parsed.Reg) >= 3 && parsed.Reg[:3] == "YMM":
				return CodeYMM
			case len(parsed.Reg) >= 3 && parsed.Reg[:3] == "ZMM":
				return CodeZMM
			case len(parsed.Reg) >= 2 && parsed.Reg[:2] == "MM":
				return CodeMM
			case len(parsed.Reg) >= 2 && parsed.Reg[:2] == "CR":
				return CodeK
			case len(parsed.Reg) >= 2 && parsed.Reg[:2] == "DR":
				return CodeK
			default:
				return CodeR32
			}
		case parsed.Mem != "":
			// Determine memory operand size based on prefix
			switch {
			case len(parsed.Mem) >= 5 && parsed.Mem[:5] == "BYTE ":
				return CodeM8
			case len(parsed.Mem) >= 6 && parsed.Mem[:6] == "WORD ":
				return CodeM16
			case len(parsed.Mem) >= 7 && parsed.Mem[:7] == "DWORD ":
				return CodeM32
			case len(parsed.Mem) >= 7 && parsed.Mem[:7] == "QWORD ":
				return CodeM64
			case len(parsed.Mem) >= 8 && parsed.Mem[:8] == "XMMWORD":
				return CodeM128
			case len(parsed.Mem) >= 8 && parsed.Mem[:8] == "YMMWORD":
				return CodeM256
			case len(parsed.Mem) >= 8 && parsed.Mem[:8] == "ZMMWORD":
				return CodeM512
			default:
				return CodeM32
			}
		case parsed.Imm != "":
			// Determine immediate size based on value
			val := parsed.Imm
			if len(val) > 1 && val[:2] == "0x" {
				hexLen := len(val) - 2
				switch {
				case hexLen <= 2:
					return CodeIMM8
				case hexLen <= 4:
					return CodeIMM16
				case hexLen <= 8:
					return CodeIMM32
				default:
					return CodeIMM32
				}
			} else {
				// Decimal immediate
				num := val
				if num[0] == '-' {
					num = num[1:]
				}
				switch {
				case len(num) <= 3 && num <= "127":
					return CodeIMM8
				case len(num) <= 5 && num <= "32767":
					return CodeIMM16
				default:
					return CodeIMM32
				}
			}
		case parsed.Seg != "":
			return CodeR16
		case parsed.Rel != "":
			return CodeREL32
		case parsed.Addr != "":
			return CodeM32
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
