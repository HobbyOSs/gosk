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
			types = append(types, getRegisterType(parsed.Reg))

		case parsed.Mem != "":
			if size := getMemorySizeFromPrefix(parsed.Mem); size != "" {
				types = append(types, size)
			} else {
				types = append(types, CodeM)
			}

		case parsed.Imm != "":
			if size := getImmediateSizeFromValue(parsed.Imm); size != "" {
				types = append(types, size)
			} else {
				types = append(types, CodeIMM)
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

	// サイズ未確定のimm/memを他のオペランドから決定
	types = resolveOperandSizes(types)

	return types
}

type Instruction struct {
	Operands []*ParsedOperand `parser:"@@ (',' @@)*"`
}

// レジスタ名からタイプを取得
func getRegisterType(reg string) OperandType {
	switch {
	case len(reg) >= 3 && reg[0] == 'E': // EAX, EBX etc
		return CodeR32
	case len(reg) == 2 && reg[1] == 'L': // AL, BL etc
		return CodeR8
	case len(reg) == 2 && reg[1] == 'X': // AX, BX etc
		return CodeR16
	case len(reg) >= 3 && reg[:3] == "XMM":
		return CodeXMM
	case len(reg) >= 3 && reg[:3] == "YMM":
		return CodeYMM
	case len(reg) >= 3 && reg[:3] == "ZMM":
		return CodeZMM
	case len(reg) >= 2 && reg[:2] == "MM":
		return CodeMM
	case len(reg) >= 2 && reg[:2] == "CR":
		return CodeK
	case len(reg) >= 2 && reg[:2] == "DR":
		return CodeK
	default:
		return CodeR32
	}
}

// メモリプレフィックスからサイズを取得
func getMemorySizeFromPrefix(mem string) OperandType {
	switch {
	case len(mem) >= 5 && mem[:5] == "BYTE ":
		return CodeM8
	case len(mem) >= 6 && mem[:6] == "WORD ":
		return CodeM16
	case len(mem) >= 7 && mem[:7] == "DWORD ":
		return CodeM32
	case len(mem) >= 7 && mem[:7] == "QWORD ":
		return CodeM64
	case len(mem) >= 8 && mem[:8] == "XMMWORD":
		return CodeM128
	case len(mem) >= 8 && mem[:8] == "YMMWORD":
		return CodeM256
	case len(mem) >= 8 && mem[:8] == "ZMMWORD":
		return CodeM512
	default:
		return ""
	}
}

// 即値からサイズを推定
func getImmediateSizeFromValue(imm string) OperandType {
	if len(imm) > 1 && imm[:2] == "0x" {
		hexLen := len(imm) - 2
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
		num := imm
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
}

// オペランドサイズを解決する
func resolveOperandSizes(types []OperandType) []OperandType {
	regSize := getRegisterSizeFromTypes(types)

	for i, t := range types {
		switch t {
		case CodeM:
			types[i] = getMemoryTypeFromRegisterSize(regSize)
		case CodeIMM, CodeIMM4, CodeIMM8, CodeIMM16, CodeIMM32:
			types[i] = getImmediateTypeFromRegisterSize(regSize)
		}
	}
	return types
}

// タイプリストからレジスタサイズを取得
func getRegisterSizeFromTypes(types []OperandType) OperandType {
	for _, t := range types {
		switch t {
		case CodeR8:
			return CodeR8
		case CodeR16:
			return CodeR16
		case CodeR32:
			return CodeR32
		default:
			return CodeR32
		}
	}
	return CodeR32
}

// レジスタサイズからメモリタイプを取得
func getMemoryTypeFromRegisterSize(regSize OperandType) OperandType {
	switch regSize {
	case CodeR8:
		return CodeM8
	case CodeR16:
		return CodeM16
	case CodeR32:
		return CodeM32
	default:
		return CodeM32
	}
}

// レジスタサイズから即値タイプを取得
func getImmediateTypeFromRegisterSize(regSize OperandType) OperandType {
	switch regSize {
	case CodeR8:
		return CodeIMM8
	case CodeR16:
		return CodeIMM16
	case CodeR32:
		return CodeIMM32
	default:
		return CodeIMM32
	}
}

var operandLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Colon", Pattern: `:`},
	{Name: "Seg", Pattern: `\b(CS|DS|ES|FS|GS|SS)\b`},
	{Name: "Reg", Pattern: `\b([ABCD]X|E?[ABCD]X|[ABCD]L|[ABCD]H|SI|DI|MM[0-7]|XMM[0-9]|YMM[0-9]|TR[0-7]|CR[0-7]|DR[0-7])\b`},
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
