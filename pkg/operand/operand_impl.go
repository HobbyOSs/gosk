package operand

import (
	"github.com/HobbyOSs/gosk/pkg/cpu" // Added import

	participle "github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var operandTypesCache = make(map[string][]OperandType)
var parsedOperandsCache = make(map[string][]*ParsedOperand)

type OperandImpl struct {
	Internal      string
	BitMode       cpu.BitMode // Changed to cpu.BitMode
	ForceImm8     bool
	ForceRelAsImm bool
}

func NewOperandFromString(text string) Operands {
	// Changed MODE_16BIT to cpu.MODE_16BIT
	return &OperandImpl{Internal: text, BitMode: cpu.MODE_16BIT, ForceImm8: false, ForceRelAsImm: false}
}

func (b *OperandImpl) WithForceRelAsImm(force bool) Operands {
	b.ForceRelAsImm = force
	return b
}

func (b *OperandImpl) WithForceImm8(force bool) Operands {
	b.ForceImm8 = force
	return b
}

func (b *OperandImpl) ParsedOperands() []*ParsedOperand {
	if cached, exists := parsedOperandsCache[b.Internal]; exists {
		return cached
	}

	parser := getParser()
	inst, err := parser.ParseString("", b.Internal)
	if err != nil || len(inst.Operands) == 0 {
		return []*ParsedOperand{} // エラー時は空のスライスを返し、キャッシュしない
	}

	parsedOperandsCache[b.Internal] = inst.Operands
	return inst.Operands
}

func (b *OperandImpl) InternalString() string {
	return b.Internal
}

var internalStringsCache = make(map[string][]string)
var internalParsedCache = make(map[string]*Instruction)

func (b *OperandImpl) getInternalParsed() *Instruction {
	if cached, exists := internalParsedCache[b.Internal]; exists {
		return cached
	}
	parser := getParser()
	inst, err := parser.ParseString("", b.Internal)
	if err != nil {
		return nil
	}
	internalParsedCache[b.Internal] = inst
	return inst
}

func (b *OperandImpl) InternalStrings() []string {
	if cached, exists := internalStringsCache[b.Internal]; exists {
		return cached
	}

	inst := b.getInternalParsed()
	if inst == nil {
		return []string{}
	}

	var results []string
	for _, parsed := range inst.Operands {
		switch {
		case parsed.SegMem != "":
			results = append(results, parsed.SegMem)
		case parsed.Reg != "":
			results = append(results, parsed.Reg)
		case parsed.DirectMem != nil:
			results = append(results, parsed.DirectMem.Addr)
		case parsed.IndirectMem != nil:
			results = append(results, parsed.IndirectMem.Mem)
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

func (b *OperandImpl) WithBitMode(mode cpu.BitMode) Operands { // Changed to cpu.BitMode
	b.BitMode = mode
	return b
}

func (b *OperandImpl) GetBitMode() cpu.BitMode { // Changed to cpu.BitMode
	return b.BitMode
}

func (b *OperandImpl) DetectImmediateSize() int {
	inst := b.getInternalParsed()
	if inst == nil {
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
		if parsed.DirectMem != nil && parsed.DirectMem.Prefix != nil {
			t := getMemorySizeFromPrefix(*parsed.DirectMem.Prefix + " " + parsed.DirectMem.Addr)
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
		if parsed.IndirectMem != nil && parsed.IndirectMem.Prefix != nil {
			t := getMemorySizeFromPrefix(*parsed.IndirectMem.Prefix + " " + parsed.IndirectMem.Mem)
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

var operandLexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Comma", Pattern: `,`},
	{Name: "SegMem", Pattern: `(CS|DS|ES|FS|GS|SS):([ABCD]X|SI|DI)`}, // このパターンは特別にアドレスとして扱う
	{Name: "Colon", Pattern: `:`},
	{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
	{Name: "MemSizePrefix", Pattern: `(BYTE|WORD|DWORD|QWORD|XMMWORD|YMMWORD|ZMMWORD)`},
	{Name: "Seg", Pattern: `(CS|DS|ES|FS|GS|SS)`},
	{Name: "Reg", Pattern: `([ABCD]X|E?[ABCD]X|[ABCD]L|[ABCD]H|SI|DI|SP|BP|MM[0-7]|XMM[0-9]|YMM[0-9]|TR[0-7]|CR[0-7]|DR[0-7])`},
	{Name: "DirectMem", Pattern: `(?:FAR\s+PTR|NEAR\s+PTR|PTR)?\s*\[\s*0x[a-fA-F0-9]+\s*\]`},
	{Name: "IndirectMem", Pattern: `(?:BYTE|WORD|DWORD|QWORD|XMMWORD|YMMWORD|ZMMWORD)?\s*\[\s*(?:[A-Za-z_][A-Za-z0-9_]*|\w+\+\w+|\w+-\w+|0x[a-fA-F0-9]+|\d+)\s*\]`},
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

	inst := b.getInternalParsed()
	if inst == nil {
		return []OperandType{OperandType("unknown")}
	}

	var types []OperandType
	for _, parsed := range inst.Operands {
		switch {
		case parsed.SegMem != "":
			types = append(types, CodeM16)
		case parsed.Reg != "":
			types = append(types, getRegisterType(parsed.Reg))
		case parsed.DirectMem != nil && parsed.DirectMem.Prefix != nil:
			types = append(types, getMemorySizeFromPrefix(*parsed.DirectMem.Prefix+" "+parsed.DirectMem.Addr))
		case parsed.IndirectMem != nil && parsed.IndirectMem.Prefix != nil:
			types = append(types, getMemorySizeFromPrefix(*parsed.IndirectMem.Prefix+" "+parsed.IndirectMem.Mem))
		case parsed.Imm != "":
			if b.ForceImm8 {
				types = append(types, CodeIMM8)
			} else {
				types = append(types, CodeIMM)
			}
		case parsed.Seg != "":
			types = append(types, CodeSREG)
		case parsed.DirectMem != nil:
			types = append(types, CodeM)
		case parsed.IndirectMem != nil:
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

func (b *OperandImpl) CalcOffsetByteSize() int {
	inst := b.getInternalParsed()
	if inst == nil {
		return 0
	}

	var total int
	for _, op := range inst.Operands {
		// 例: op.IndirectMem == "[EBX+16]" とか op.DirectMem == "[0x0ff0]" とかが入る
		if op.IndirectMem != nil {
			size := calcMemOffsetSize(op.IndirectMem.Mem)
			total += size
		}
		if op.DirectMem != nil {
			size := calcMemOffsetSize(op.DirectMem.Addr)
			total += size
		}
	}
	return total
}
