package operand

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/HobbyOSs/gosk/internal/ast"

	participle "github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/samber/lo"
)

var operandTypesCache = make(map[string][]OperandType)

type OperandImpl struct {
	Internal      string
	BitMode       ast.BitMode
	ForceImm8     bool
	ForceRelAsImm bool
}

func NewOperandFromString(text string) Operands {
	return &OperandImpl{Internal: text, BitMode: ast.MODE_16BIT, ForceImm8: false, ForceRelAsImm: false}
}

func (b *OperandImpl) WithForceRelAsImm(force bool) Operands {
	b.ForceRelAsImm = force
	return b
}

func (b *OperandImpl) WithForceImm8(force bool) Operands {
	b.ForceImm8 = force
	return b
}

func (b *OperandImpl) InternalString() string {
	return b.Internal
}

var internalStringsCache = make(map[string][]string)

func (b *OperandImpl) InternalStrings() []string {
	if cached, exists := internalStringsCache[b.Internal]; exists {
		return cached
	}

	parser := getParser()
	inst, err := parser.ParseString("", b.Internal)
	if err != nil || len(inst.Operands) == 0 {
		return []string{}
	}

	var results []string
	for _, parsed := range inst.Operands {
		switch {
		case parsed.SegMem != "":
			results = append(results, parsed.SegMem)
		case parsed.Reg != "":
			results = append(results, parsed.Reg)
		case parsed.Addr != nil:
			results = append(results, parsed.Addr.Addr)
		case parsed.Mem != nil:
			results = append(results, parsed.Mem.Mem)
		case parsed.Imm != "":
			results = append(results, parsed.Imm)
		case parsed.Seg != "":
			results = append(results, parsed.Seg)
		case parsed.Rel != "":
			results = append(results, parsed.Rel)
		}
	}

	internalStringsCache[b.Internal] = results
	return results
}

func (b *OperandImpl) Serialize() string {
	return b.Internal
}

func (b *OperandImpl) FromString(text string) Operands {
	return &OperandImpl{Internal: text, BitMode: b.BitMode}
}

func (b *OperandImpl) WithBitMode(mode ast.BitMode) Operands {
	b.BitMode = mode
	return b
}

func (b *OperandImpl) GetBitMode() ast.BitMode {
	return b.BitMode
}

func (b *OperandImpl) DetectImmediateSize() int {
	parser := getParser()
	inst, err := parser.ParseString("", b.Internal)
	if err != nil || len(inst.Operands) == 0 {
		return 0
	}

	if len(inst.Operands) == 1 {
		parsed := inst.Operands[0]
		if parsed.Imm != "" {
			s := getImmediateSizeFromValue(parsed.Imm)
			switch s {
			case CodeIMM8:
				return 1
			case CodeIMM16:
				return 2
			case CodeIMM32:
				return 4
			}
		}
		return 0
	}

	for _, parsed := range inst.Operands {
		if parsed.Addr != nil && parsed.Addr.Prefix != nil {
			t := getMemorySizeFromPrefix(*parsed.Addr.Prefix + " " + parsed.Addr.Addr)
			switch t {
			case CodeM8:
				return 1
			case CodeM16:
				return 2
			case CodeM32:
				return 4
			}
			break
		}
		if parsed.Mem != nil && parsed.Mem.Prefix != nil {
			t := getMemorySizeFromPrefix(*parsed.Mem.Prefix + " " + parsed.Mem.Mem)
			switch t {
			case CodeM8:
				return 1
			case CodeM16:
				return 2
			case CodeM32:
				return 4
			}
			break
		}
		if parsed.Reg != "" {
			if b.ForceImm8 {
				return 1
			}
			t := getRegisterType(parsed.Reg)
			switch t {
			case CodeR8:
				return 1
			case CodeR16:
				return 2
			case CodeR32:
				return 4
			}
			break
		}
	}
	return 0
}

type Instruction struct {
	Operands []*ParsedOperand `parser:"@@(',' @@)*"`
}

type ParsedOperand struct {
	SegMem string `parser:"@SegMem"`
	Reg    string `parser:"| @Reg"`
	Addr   *Addr  `parser:"| @@"`
	Mem    *Mem   `parser:"| @@"`
	Imm    string `parser:"| @Imm"`
	Seg    string `parser:"| @Seg"`
	Rel    string `parser:"| @Rel"`
}

type Mem struct {
	Prefix *string `parser:"@MemSizePrefix?"`
	Mem    string  `parser:"@Mem"`
}

type Addr struct {
	Prefix *string `parser:"@MemSizePrefix?"`
	Addr   string  `parser:"@Addr"`
}

var operandLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Comma", Pattern: `,`},
	{Name: "SegMem", Pattern: `(CS|DS|ES|FS|GS|SS):([ABCD]X|SI|DI)`}, // このパターンは特別にアドレスとして扱う
	{Name: "Colon", Pattern: `:`},
	{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
	{Name: "MemSizePrefix", Pattern: `(BYTE|WORD|DWORD|QWORD|XMMWORD|YMMWORD|ZMMWORD)`},
	{Name: "Seg", Pattern: `(CS|DS|ES|FS|GS|SS)`},
	{Name: "Reg", Pattern: `([ABCD]X|E?[ABCD]X|[ABCD]L|[ABCD]H|SI|DI|SP|BP|MM[0-7]|XMM[0-9]|YMM[0-9]|TR[0-7]|CR[0-7]|DR[0-7])`},
	{Name: "Addr", Pattern: `(?:FAR\s+PTR|NEAR\s+PTR|PTR)?\s*\[\s*0x[a-fA-F0-9]+\s*\]`},
	{Name: "Mem", Pattern: `(?:BYTE|WORD|DWORD|QWORD|XMMWORD|YMMWORD|ZMMWORD)?\s*\[\s*(?:[A-Za-z_][A-Za-z0-9_]*|\w+\+\w+|\w+-\w+|0x[a-fA-F0-9]+|\d+)\s*\]`},
	{Name: "Imm", Pattern: `(0x[a-fA-F0-9]+|-?\d+)`},
	{Name: "Rel", Pattern: `(?:SHORT|FAR PTR)?\s*\w+`},
	{Name: "String", Pattern: `"(?:\\.|[^"\\])*"`},
})

func getParser() *participle.Parser[Instruction] {
	return participle.MustBuild[Instruction](
		participle.Lexer(operandLexer),
		participle.Unquote(),
		participle.Elide("Whitespace"),
	)
}

func (b *OperandImpl) OperandTypes() []OperandType {
	if cached, exists := operandTypesCache[b.Internal]; exists {
		return cached
	}

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
		case parsed.Addr != nil && parsed.Addr.Prefix != nil:
			types = append(types, getMemorySizeFromPrefix(*parsed.Addr.Prefix+" "+parsed.Addr.Addr))
		case parsed.Mem != nil && parsed.Mem.Prefix != nil:
			types = append(types, getMemorySizeFromPrefix(*parsed.Mem.Prefix+" "+parsed.Mem.Mem))
		case parsed.Imm != "":
			if b.ForceImm8 {
				types = append(types, CodeIMM8)
			} else {
				types = append(types, CodeIMM)
			}
		case parsed.Seg != "":
			types = append(types, CodeSREG)
		case parsed.Addr != nil:
			types = append(types, CodeM)
		case parsed.Mem != nil:
			types = append(types, CodeM)
		case parsed.Rel != "":
			// ラベル指定
			if b.ForceRelAsImm {
				types = append(types, CodeIMM) // ForceRelAsImm が true なら Imm として扱う
			} else {
				if len(parsed.Rel) >= 5 && parsed.Rel[:5] == "SHORT" {
					types = append(types, CodeREL8)
				} else {
					types = append(types, CodeREL32)
				}
			}
		default:
			types = append(types, OperandType("unknown"))
		}
	}
	// サイズ未確定のimm/memを他のオペランドから決定
	types = b.resolveOperandSizes(types, inst.Operands)

	operandTypesCache[b.Internal] = types
	return types
}

var (
	regR8Pattern  = regexp.MustCompile(`^[ABCD][HL]$`)
	regR16Pattern = regexp.MustCompile(`^(?:[ABCD]X|SP|BP|SI|DI)$`)
	regR32Pattern = regexp.MustCompile(`^E[ABCD]X|ESI|EDI$`)
)

// レジスタ名からタイプを取得
func getRegisterType(reg string) OperandType {
	// Prefixで判定できるもの
	switch {
	case strings.HasPrefix(reg, "XMM"):
		return CodeXMM
	case strings.HasPrefix(reg, "YMM"):
		return CodeYMM
	case strings.HasPrefix(reg, "ZMM"):
		return CodeZMM
	case strings.HasPrefix(reg, "MM"):
		return CodeMM
	case strings.HasPrefix(reg, "CR"):
		return CodeCR
	case strings.HasPrefix(reg, "DR"):
		return CodeDR
	case strings.HasPrefix(reg, "TR"):
		return CodeTR
	}

	// 正規表現で判定するもの
	switch {
	case regR8Pattern.MatchString(reg):
		return CodeR8
	case regR16Pattern.MatchString(reg):
		return CodeR16
	case regR32Pattern.MatchString(reg):
		return CodeR32
	}

	return CodeR32
}

// メモリプレフィックスからサイズを取得
func getMemorySizeFromPrefix(_mem string) OperandType {

	mem := strings.ToUpper(strings.TrimSpace(_mem))

	switch {
	case strings.HasPrefix(mem, "BYTE"):
		return CodeM8
	case strings.HasPrefix(mem, "WORD"):
		return CodeM16
	case strings.HasPrefix(mem, "DWORD"):
		return CodeM32
	case strings.HasPrefix(mem, "QWORD"):
		return CodeM64
	case strings.HasPrefix(mem, "XMMWORD"):
		return CodeM128
	case strings.HasPrefix(mem, "YMMWORD"):
		return CodeM256
	case strings.HasPrefix(mem, "ZMMWORD"):
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
func (b *OperandImpl) resolveOperandSizes(types []OperandType, operands []*ParsedOperand) []OperandType {
	regSize := getOperandSizeFromTypesLo(types, operands)

	for i, t := range types {
		switch t {
		case CodeM:
			types[i] = getMemoryTypeFromRegisterSize(regSize)
		case CodeIMM, CodeIMM4, CodeIMM8, CodeIMM16, CodeIMM32:
			if b.ForceImm8 {
				types[i] = CodeIMM8
			} else {
				types[i] = getImmediateTypeFromRegisterSize(regSize)
			}
		}
	}
	return types
}

// タイプリストからレジスタサイズを取得 (samber/lo バージョン)
func getOperandSizeFromTypesLo(types []OperandType, operands []*ParsedOperand) OperandType {
	foundType, _ := lo.Find(types, func(t OperandType) bool {
		return lo.Contains([]OperandType{CodeR8, CodeM8, CodeR16, CodeM16, CodeR32, CodeM32}, t)
	})

	return lo.Switch[OperandType, OperandType](foundType).
		Case(CodeR8, CodeR8).
		Case(CodeM8, CodeR8).
		Case(CodeR16, CodeR16).
		Case(CodeM16, CodeR16).
		Case(CodeR32, CodeR32).
		Case(CodeM32, CodeR32).
		DefaultF(func() OperandType {
			if lo.Contains(types, CodeM) {
				i := lo.IndexOf(types, CodeM) // CodeM のインデックスを取得 (最初の出現箇所)
				if operands[i].Addr != nil {
					size := calcMemOffsetSize(operands[i].Addr.Addr)
					return lo.Switch[int, OperandType](size).
						Case(1, CodeR8).
						Case(2, CodeR16).
						Default(CodeR32)
				}
			}
			return CodeR32 // デフォルト値
		})
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

// Require66h はオペランドサイズプレフィックスが必要かどうかを判定する
func (b *OperandImpl) Require66h() bool {
	types := b.OperandTypes()
	if len(types) == 0 {
		return false
	}

	switch b.BitMode {
	case ast.MODE_16BIT:
		// 16bitモードで32bitレジスタを使用する場合
		for _, t := range types {
			if t == CodeR32 || t == CodeM32 {
				return true
			}
		}
		// 16bitモードで32bit即値を使用する場合
		if len(types) == 1 {
			parser := getParser()
			inst, err := parser.ParseString("", b.Internal)
			if err == nil && len(inst.Operands) == 1 && inst.Operands[0].Imm != "" {
				imm := getImmediateSizeFromValue(inst.Operands[0].Imm)
				if imm == CodeIMM32 {
					return true
				}
			}
		}
	case ast.MODE_32BIT:
		// 32bitモードで16bitレジスタを使用する場合
		for _, t := range types {
			if t == CodeR16 || t == CodeM16 {
				return true
			}
		}
	}
	return false
}

// Require67h はアドレスサイズプレフィックスが必要かどうかを判定する
func (b *OperandImpl) Require67h() bool {
	types := b.OperandTypes()
	if len(types) == 0 {
		return false
	}

	switch b.BitMode {
	case ast.MODE_16BIT:
		// 16bitモードで32bitメモリアクセスを行う場合
		for _, t := range types {
			if t == CodeM32 {
				return true
			}
		}
	case ast.MODE_32BIT:
		// 32bitモードで16bitメモリアクセスを行う場合
		for _, t := range types {
			if t == CodeM16 {
				return true
			}
		}
	}
	return false
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
		// 例: op.Mem == "[EBX+16]" とか op.Addr == "[0x0ff0]" とかが入る
		if op.Mem != nil {
			size := calcMemOffsetSize(op.Mem.Mem)
			total += size
		}
		if op.Addr != nil {
			size := calcMemOffsetSize(op.Addr.Addr)
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
