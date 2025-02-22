package operand

import (
	participle "github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type OperandImpl struct {
	Internal string
}

func (b *OperandImpl) OperandTypes() []OperandType {
	parser := getParser()
	inst, err := parser.ParseString("", b.Internal)
	if err != nil || len(inst.Operands) == 0 {
		return []OperandType{OperandType("unknown")}
	}

	var types []OperandType
	for _, parsed := range inst.Operands {
		switch {
		case parsed.Reg != "":
			// Determine register size based on name
			switch {
			case len(parsed.Reg) >= 3 && parsed.Reg[0] == 'R': // RAX, RBX etc
				types = append(types, CodeR32)
			case len(parsed.Reg) >= 3 && parsed.Reg[0] == 'E': // EAX, EBX etc
				types = append(types, CodeR32)
			case len(parsed.Reg) == 2 && parsed.Reg[1] == 'L': // AL, BL etc
				types = append(types, CodeR8)
			case len(parsed.Reg) == 2 && parsed.Reg[1] == 'X': // AX, BX etc
				types = append(types, CodeR16)
			case len(parsed.Reg) >= 3 && parsed.Reg[:3] == "XMM":
				types = append(types, CodeXMM)
			case len(parsed.Reg) >= 3 && parsed.Reg[:3] == "YMM":
				types = append(types, CodeYMM)
			case len(parsed.Reg) >= 3 && parsed.Reg[:3] == "ZMM":
				types = append(types, CodeZMM)
			case len(parsed.Reg) >= 2 && parsed.Reg[:2] == "MM":
				types = append(types, CodeMM)
			case len(parsed.Reg) >= 2 && parsed.Reg[:2] == "CR":
				types = append(types, CodeK)
			case len(parsed.Reg) >= 2 && parsed.Reg[:2] == "DR":
				types = append(types, CodeK)
			default:
				types = append(types, CodeR32)
			}
		case parsed.Mem != "":
			// Determine memory operand size based on prefix
			switch {
			case len(parsed.Mem) >= 5 && parsed.Mem[:5] == "BYTE ":
				types = append(types, CodeM8)
			case len(parsed.Mem) >= 6 && parsed.Mem[:6] == "WORD ":
				types = append(types, CodeM16)
			case len(parsed.Mem) >= 7 && parsed.Mem[:7] == "DWORD ":
				types = append(types, CodeM32)
			case len(parsed.Mem) >= 7 && parsed.Mem[:7] == "QWORD ":
				types = append(types, CodeM64)
			case len(parsed.Mem) >= 8 && parsed.Mem[:8] == "XMMWORD":
				types = append(types, CodeM128)
			case len(parsed.Mem) >= 8 && parsed.Mem[:8] == "YMMWORD":
				types = append(types, CodeM256)
			case len(parsed.Mem) >= 8 && parsed.Mem[:8] == "ZMMWORD":
				types = append(types, CodeM512)
			default:
				types = append(types, CodeM32)
			}
		case parsed.Imm != "":
			// Determine immediate size based on value
			val := parsed.Imm
			if len(val) > 1 && val[:2] == "0x" {
				hexLen := len(val) - 2
				switch {
				case hexLen <= 2:
					types = append(types, CodeIMM8)
				case hexLen <= 4:
					types = append(types, CodeIMM16)
				case hexLen <= 8:
					types = append(types, CodeIMM32)
				default:
					types = append(types, CodeIMM32)
				}
			} else {
				// Decimal immediate
				num := val
				if num[0] == '-' {
					num = num[1:]
				}
				switch {
				case len(num) <= 3 && num <= "127":
					types = append(types, CodeIMM8)
				case len(num) <= 5 && num <= "32767":
					types = append(types, CodeIMM16)
				default:
					types = append(types, CodeIMM32)
				}
			}
		case parsed.Seg != "":
			types = append(types, CodeR16)
		case parsed.Rel != "":
			types = append(types, CodeREL32)
		case parsed.Addr != "":
			types = append(types, CodeM32)
		default:
			types = append(types, OperandType("unknown"))
		}
	}
	return types
}

type Instruction struct {
	Operands []*ParsedOperand `parser:"@@ (',' @@)*"`
}

var operandLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Colon", Pattern: `:`},
	{Name: "Seg", Pattern: `\b(CS|DS|ES|FS|GS|SS)\b`},
	{Name: "Reg", Pattern: `\b(R?[ABCD]X|E?[ABCD]X|[ABCD]L|[ABCD]H|SI|DI|MM[0-7]|XMM[0-9]|YMM[0-9]|TR[0-7]|CR[0-7]|DR[0-7])\b`},
	{Name: "MemPrefix", Pattern: `\b(BYTE|WORD|DWORD|QWORD|XMMWORD|YMMWORD|ZMMWORD)\b`},
	{Name: "Ptr", Pattern: `\bPTR\b`},
	{Name: "Addr", Pattern: `\[(?: (?: (?:FAR\s+PTR|NEAR\s+PTR|PTR)\s+ )? \[0x[a-fA-F0-9]+\] )\]`},
	{Name: "Mem", Pattern: `\[(?:[A-Za-z_][A-Za-z0-9_]*|\w+\+\w+|\w+-\w+|0x[a-fA-F0-9]+|\d+)\]`},
	{Name: "Imm", Pattern: `(?:0x[a-fA-F0-9]+|-?\d+)`},
	{Name: "Rel", Pattern: `\b(JMP|CALL|JNZ|JE|JNE|LABEL)\s+\w+\b`},
	{Name: "String", Pattern: `"(?:\\.|[^"\\])*"`},
	{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
	{Name: "Comma", Pattern: `,`},
})

func getParser() *participle.Parser[Instruction] {
	return participle.MustBuild[Instruction](
		participle.Lexer(operandLexer),
		participle.Unquote(),
		participle.Elide("Whitespace"),
	)
}
