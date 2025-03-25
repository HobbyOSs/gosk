package operand

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

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
				if operands[i].DirectMem != nil {
					size := calcMemOffsetSize(operands[i].DirectMem.Addr)
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

var reBaseOffset = regexp.MustCompile(`^\[\s*([A-Za-z0-9]+)\s*(?:\+|\-)\s*([0-9A-Fa-fx]+)\s*\]$`)
var reDirect = regexp.MustCompile(`^\[\s*([0-9A-Fa-fx]+)\s*\]$`)

// CalcOffsetByteSize
// メモリーアドレス表現にあるoffset値について機械語サイズの計算をする
// * ベースを持たない直接のアドレス表現 e.g. MOV CL,[0x0ff0]; の場合2byteを返す
// * ベースがある場合のアドレス表現     e.g. MOV ECX,[EBX+16]; の場合1byteを返す
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
