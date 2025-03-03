package operand

import (
	"regexp"
	"strconv"
	"strings"

	participle "github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type OperandImpl struct {
	Internal string
}

func NewOperandFromString(text string) Operands {
	return &OperandImpl{Internal: text}
}

func (b *OperandImpl) InternalString() string {
	return b.Internal
}

func (b *OperandImpl) Serialize() string {
	return b.Internal
}

func (b *OperandImpl) FromString(text string) Operands {
	return &OperandImpl{Internal: text}
}

func (b *OperandImpl) DetectImmediateSize() int {
	parser := getParser()
	inst, err := parser.ParseString("", b.Internal)
	if err != nil || len(inst.Operands) == 0 {
		return 0
	}

	size := 0
	for _, parsed := range inst.Operands {
		if parsed.Imm != "" {
			s := getImmediateSizeFromValue(parsed.Imm)
			switch s {
			case CodeIMM8:
				size = 1
			case CodeIMM16:
				size = 2
			case CodeIMM32:
				size = 4
			}
		}
	}
	return size
}

type Instruction struct {
	Operands []*ParsedOperand `parser:"@@(',' @@)*"`
}

var operandLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "SegMem", Pattern: `\b(CS|DS|ES|FS|GS|SS):([ABCD]X|SI|DI)\b`}, // このパターンは特別にアドレスとして扱う
	{Name: "Colon", Pattern: `:`},
	{Name: "Seg", Pattern: `\b(CS|DS|ES|FS|GS|SS)\b`},
	{Name: "Reg", Pattern: `\b([ABCD]X|E?[ABCD]X|[ABCD]L|[ABCD]H|SI|DI|MM[0-7]|XMM[0-9]|YMM[0-9]|TR[0-7]|CR[0-7]|DR[0-7])\b`},
	{Name: "MemPrefix", Pattern: `\b(BYTE|WORD|DWORD|QWORD|XMMWORD|YMMWORD|ZMMWORD)\b`},
	{Name: "Addr", Pattern: `(?:FAR\s+PTR|NEAR\s+PTR|PTR)?\s*(?:BYTE|WORD|DWORD|QWORD|XMMWORD|YMMWORD|ZMMWORD)?\s*\[\s*0x[a-fA-F0-9]+\s*\]`},
	{Name: "Mem", Pattern: `(?:BYTE|WORD|DWORD|QWORD|XMMWORD|YMMWORD|ZMMWORD)?\s*\[\s*(?:[A-Za-z_][A-Za-z0-9_]*|\w+\+\w+|\w+-\w+|0x[a-fA-F0-9]+|\d+)\s*\]`},
	{Name: "Imm", Pattern: `(?:0x[a-fA-F0-9]+|-?\d+)`},
	{Name: "Rel", Pattern: `\b(?:SHORT|FAR PTR)?\s*\w+\b`},
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

func (b *OperandImpl) OperandTypes() []OperandType {
	parser := getParser()
	inst, err := parser.ParseString("", b.Internal)
	if err != nil || len(inst.Operands) == 0 {
		return []OperandType{OperandType("unknown")}
	}

	var types []OperandType
	for _, parsed := range inst.Operands {
		switch {
		case parsed.SegMem != "":
			types = append(types, CodeM16)
		case parsed.Reg != "":
			types = append(types, getRegisterType(parsed.Reg))
		case parsed.MemPrefix != "" && parsed.Addr != "":
			types = append(types, getMemorySizeFromPrefix(parsed.MemPrefix+" ["+parsed.Addr+"]"))
		case parsed.MemPrefix != "" && parsed.Mem != "":
			types = append(types, getMemorySizeFromPrefix(parsed.MemPrefix+" "+parsed.Mem))
		case parsed.Imm != "":
			if size := getImmediateSizeFromValue(parsed.Imm); size != "" {
				types = append(types, size)
			} else {
				types = append(types, CodeIMM)
			}
		case parsed.MemPrefix != "" && parsed.Imm != "":
			// BYTE 8 のようなケースに対応
			types = append(types, getMemorySizeFromPrefix(parsed.MemPrefix+" ["+parsed.Imm+"]"))
		case parsed.Seg != "":
			types = append(types, CodeR16)
		case parsed.Rel != "":
			if len(parsed.Rel) >= 5 && parsed.Rel[:5] == "SHORT" {
				types = append(types, CodeREL8)
			} else {
				types = append(types, CodeREL32)
			}
		case parsed.Addr != "":
			types = append(types, CodeM32)
		case parsed.MemPrefix != "" && parsed.Mem != "":
			types = append(types, getMemorySizeFromPrefix(parsed.MemPrefix+" "+parsed.Mem))
		default:
			types = append(types, OperandType("unknown"))
		}
	}

	// サイズ未確定のimm/memを他のオペランドから決定
	types = resolveOperandSizes(types)

	return types
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
		return CodeCR
	case len(reg) >= 2 && reg[:2] == "DR":
		return CodeDR
	case len(reg) >= 2 && reg[:2] == "TR":
		return CodeTR
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

var reBaseOffset = regexp.MustCompile(`^\[\s*([A-Za-z0-9]+)\s*(?:\+|\-)\s*([0-9A-Fa-fx]+)\s*\]$`)
var reDirect = regexp.MustCompile(`^\[\s*([0-9A-Fa-fx]+)\s*\]$`)

// CalcOffsetByteSize
// メモリーアドレス表現にあるoffset値について機械語サイズの計算をする
// * ベースを持たない直接のアドレス表現 e.g. MOV CL,[0x0ff0]; の場合2byteを返す
// * ベースがある場合のアドレス表現     e.g. MOV ECX,[EBX+16]; の場合1byteを返す
func (b *OperandImpl) CalcOffsetByteSize() int {
	parser := getParser()
	inst, err := parser.ParseString("", b.Internal)
	if err != nil {
		return 0
	}

	var total int
	for _, op := range inst.Operands {
		// 例: op.Mem == "[EBX+16]" とか op.Mem == "[0x0ff0]" とかが入る
		if op.Mem != "" {
			size := calcMemOffsetSize(op.Mem)
			total += size
		}
	}
	return total
}

func calcMemOffsetSize(mem string) int {
	// まずベースレジスタがないパターン（[0x0ff0]など）を判定
	if m := reDirect.FindStringSubmatch(mem); m != nil {
		offsetVal, err := parseNumeric(m[1])
		if err != nil {
			return 0
		}
		// ベースなし ⇒ GetOffsetSize相当
		return getOffsetSize(offsetVal)
	}

	// ベースレジスタがあるパターン([EBX+16], [ECX-0x80]など)を判定
	if m := reBaseOffset.FindStringSubmatch(mem); m != nil {
		// m[1] がベース(EBX等), m[2] がオフセット値(16等)
		offsetVal, err := parseNumeric(m[2])
		if err != nil {
			return 0
		}
		// ベース有り ⇒ 0の場合はサイズ0, そうでなければ-128~127 ⇒ 1バイト, …というロジック
		if offsetVal == 0 {
			return 0
		}
		if offsetVal >= -0x80 && offsetVal <= 0x7f {
			return 1
		}
		if offsetVal >= -0x8000 && offsetVal <= 0x7fff {
			return 2
		}
		return 4
	}
	// 上記のどれにも当てはまらない=パターン外。必要に応じて厳密に扱う
	return 0
}

// -128～127, -32768～32767 などの判定に使う
func getOffsetSize(imm int64) int {
	if imm >= -0x80 && imm <= 0x7f {
		return 1
	}
	if imm >= -0x8000 && imm <= 0x7fff {
		return 2
	}
	return 4
}

// 文字列の数値(例: "16", "0x0ff0", "-123")をint64に変換
func parseNumeric(s string) (int64, error) {
	s = strings.ToLower(strings.TrimSpace(s))
	negative := false
	if strings.HasPrefix(s, "-") {
		negative = true
		s = s[1:]
	}

	base := 10
	if strings.HasPrefix(s, "0x") {
		base = 16
		s = s[2:]
	}
	val, err := strconv.ParseInt(s, base, 64)
	if err != nil {
		return 0, err
	}
	if negative {
		val = -val
	}
	return val, nil
}
