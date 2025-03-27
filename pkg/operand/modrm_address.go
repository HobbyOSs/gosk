package operand

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/HobbyOSs/gosk/pkg/cpu" // Added import
)

// ModRMは ModR/Mバイトおよびディスプレースメントを保持するサンプル構造体
type ModRM struct {
	Mod  byte   // 上位2ビット
	Reg  byte   // 中位3ビット
	RM   byte   // 下位3ビット
	Disp []byte // ディスプレースメント(必要な場合のみ)
}

// 16ビットアドレッシング用の簡易マップ
var rm16Table = map[string]byte{
	"bx+si": 0b000,
	"bx+di": 0b001,
	"bp+si": 0b010,
	"bp+di": 0b011,
	"si":    0b100,
	"di":    0b101,
	"bp":    0b110,
	"bx":    0b111,
}

// 32ビットアドレッシング用の簡易マップ（SIB非対応）
var rm32Table = map[string]byte{
	"eax": 0b000,
	"ecx": 0b001,
	"edx": 0b010,
	"ebx": 0b011,
	"esp": 0b100,
	"ebp": 0b101,
	"esi": 0b110,
	"edi": 0b111,
}

// ParseMemoryOperand はメモリオペランド文字列(例: "[bx+si+0x10]" など)を
// パースして (mod, rm, disp) を求める。bitMode には cpu.MODE_16BIT, cpu.MODE_32BIT を使う。
func ParseMemoryOperand(memStr string, bitMode cpu.BitMode) (mod byte, rm byte, disp []byte, err error) { // Changed to cpu.BitMode
	if !strings.Contains(memStr, "[") && !strings.HasSuffix(memStr, "]") {
		return 0, 0, nil, fmt.Errorf("invalid mem operand format: %s", memStr)
	}
	_memStr := removePrefix(memStr)                         // BYTE,WORD,DWORDを除去
	inner := strings.TrimSpace(_memStr[1 : len(_memStr)-1]) // "bx+si+0x10"のように[]を除去

	parts := strings.Split(inner, "+")
	var regs []string
	var dispValue int64
	var hasDisp bool

	for _, p := range parts {
		p = strings.TrimSpace(p)
		// 16進(0x除去), 10進の両方をまとめて試す
		// 例: 0x10 => 16進として解釈 / 10 => 10進 として解釈
		if valHex, errHex := strconv.ParseInt(strings.ReplaceAll(p, "0x", ""), 16, 64); errHex == nil {
			dispValue = valHex
			hasDisp = true
		} else if valDec, errDec := strconv.ParseInt(p, 10, 64); errDec == nil {
			dispValue = valDec
			hasDisp = true
		} else {
			regs = append(regs, strings.ToLower(p))
		}
	}

	switch bitMode {
	case cpu.MODE_16BIT: // Changed to cpu.MODE_16BIT
		baseKey := strings.Join(regs, "+") // 例: "bx+si"
		if baseKey == "" && hasDisp {
			// [0x0ff0] のような直接アドレッシングを処理する
			mod = 0b00 // mod=00 は disp16 参照
			rm = 0b110 // 16ビットモードにおける [disp16] は rm=110
			lo := byte(dispValue & 0xFF)
			hi := byte((dispValue >> 8) & 0xFF)
			disp = []byte{lo, hi}
			return mod, rm, disp, nil
		}

		rmVal, ok := rm16Table[baseKey]
		if !ok {
			// --- START MODIFICATION ---
			// 16bitモードだが、32bitレジスタが使われているか試す
			// (例: [ESI], [EBX+disp])
			// この場合、アドレスサイズプレフィックス(0x67)が必要になる (呼び出し元で判断)
			if len(regs) == 1 { // SIBは考慮しない
				rmVal32, ok32 := rm32Table[regs[0]]
				if ok32 {
					// 32bitレジスタが見つかった
					rmVal = rmVal32 // R/Mビットは32bitテーブルのものを使う
					ok = true       // "見つかった"ことにする
					// Mod/Dispの決定ロジックは後続のswitch-caseに任せる
				}
			}
			if !ok {
				// 16bitでも32bitでも見つからない場合
				return 0, 0, nil, fmt.Errorf("unsupported 16bit mem operand: %q", baseKey)
			}
			// --- END MODIFICATION ---
		}

		// Mod/Dispの決定 (既存ロジック)
		switch {
		case !hasDisp:
			// dispなし => mod=00。ただし [bp]単独はmod=00,r/m=110が disp16=0 扱い
			if baseKey == "bp" {
				mod = 0b00
				rm = rmVal
				disp = []byte{0x00, 0x00} // disp16=0
			} else {
				mod = 0b00
				rm = rmVal
				disp = nil
			}
		case dispValue >= -128 && dispValue <= 127:
			mod = 0b01
			rm = rmVal
			disp = []byte{byte(dispValue & 0xFF)}
		default:
			mod = 0b10
			rm = rmVal
			lo := byte(dispValue & 0xFF)
			hi := byte((dispValue >> 8) & 0xFF)
			disp = []byte{lo, hi}
		}

	case cpu.MODE_32BIT: // Changed to cpu.MODE_32BIT
		if len(regs) == 0 {
			// [disp32] のケース ([0x1000] など)
			mod = 0b00
			rm = 0b101
			disp = toLe32(dispValue)
			break
		}
		if len(regs) > 1 {
			// SIB必須ケースは未対応
			return 0, 0, nil, fmt.Errorf("SIB not supported in this sample: %v", regs)
		}
		baseKey := regs[0]
		rmVal, ok := rm32Table[baseKey]
		if !ok {
			return 0, 0, nil, fmt.Errorf("unsupported 32bit register: %q", baseKey)
		}

		switch {
		case !hasDisp && rmVal == 0b101:
			// [ebp] の mod=00,r/m=101 は disp32=0 扱い
			mod = 0b00
			rm = 0b101
			disp = []byte{0, 0, 0, 0}
		case !hasDisp:
			mod = 0b00
			rm = rmVal
		case dispValue >= -128 && dispValue <= 127:
			mod = 0b01
			rm = rmVal
			disp = []byte{byte(dispValue & 0xFF)}
		default:
			mod = 0b10
			rm = rmVal
			disp = toLe32(dispValue)
		}

	default:
		return 0, 0, nil, fmt.Errorf("unsupported bitMode %d", bitMode)
	}

	return mod, rm, disp, nil
}

// CalcModRM は最終的な ModR/Mバイト + disp を生成する。
// rmOperand が "[...]" の場合はメモリ、レジスタ名の場合は mod=11 とみなす。
func CalcModRM(rmOperand string, regBits byte, bitMode cpu.BitMode) ([]byte, error) { // Changed to cpu.BitMode
	rmOperand = strings.TrimSpace(rmOperand)

	// メモリオペランドかどうか
	if strings.HasPrefix(rmOperand, "[") && strings.HasSuffix(rmOperand, "]") {
		mod, rm, disp, err := ParseMemoryOperand(rmOperand, bitMode)
		if err != nil {
			return nil, err
		}
		modrmByte := combineModRM(mod, regBits, rm)
		out := []byte{modrmByte}
		if disp != nil {
			out = append(out, disp...)
		}
		return out, nil
	}

	// レジスタオペランド (mod=11)
	regMap16 := map[string]byte{
		"ax": 0b000, "cx": 0b001, "dx": 0b010, "bx": 0b011,
		"sp": 0b100, "bp": 0b101, "si": 0b110, "di": 0b111,
	}
	regMap32 := map[string]byte{
		"eax": 0b000, "ecx": 0b001, "edx": 0b010, "ebx": 0b011,
		"esp": 0b100, "ebp": 0b101, "esi": 0b110, "edi": 0b111,
	}

	rLower := strings.ToLower(rmOperand)
	var rmVal byte
	var ok bool
	switch bitMode {
	case cpu.MODE_16BIT: // Changed to cpu.MODE_16BIT
		rmVal, ok = regMap16[rLower]
	case cpu.MODE_32BIT: // Changed to cpu.MODE_32BIT
		rmVal, ok = regMap32[rLower]
	default:
		return nil, fmt.Errorf("unsupported bitMode %d", bitMode)
	}
	if !ok {
		return nil, fmt.Errorf("unknown register: %s", rmOperand)
	}
	mod := byte(0b11)
	modrmByte := combineModRM(mod, regBits, rmVal)
	return []byte{modrmByte}, nil
}

func combineModRM(mod, reg, rm byte) byte {
	return (mod << 6) | ((reg & 0b111) << 3) | (rm & 0b111)
}

func toLe32(v int64) []byte {
	b := make([]byte, 4)
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	return b
}

// removePrefix は BYTE, WORD, DWORD のようなプレフィックスを取り除く関数
func removePrefix(input string) string {
	// 正規表現パターンを定義
	pattern := regexp.MustCompile(`^(BYTE|WORD|DWORD)\s+(.*)$`)

	// パターンにマッチするか確認
	matches := pattern.FindStringSubmatch(input)
	if len(matches) > 2 {
		// プレフィックスを取り除いた部分を返す
		return strings.TrimSpace(matches[2])
	}

	// マッチしない場合は元の文字列を返す
	return strings.TrimSpace(input)
}
