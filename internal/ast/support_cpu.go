package ast

import "strings"

type SupCPU int

// http://ref.x86asm.net/index.html
// --------------------------------
// 00: 8086
// 01: 80186
// 02: 80286
// 03: 80386 (i386)
// 04: 80486 (i486)
// P1 (05): Pentium (1)
// PX (06): Pentium with MMX
// PP (07): Pentium Pro
// P2 (08): Pentium II
// P3 (09): Pentium III
// P4 (10): Pentium 4
// C1 (11): Core (1)
// C2 (12): Core 2
// C7 (13): Core i7
// IT (99): Itanium (only geek editions)
const (
	SUP_8086           SupCPU = 0
	SUP_80186          SupCPU = 1
	SUP_80286          SupCPU = 2
	SUP_80386          SupCPU = 3
	SUP_80486          SupCPU = 4
	SUP_Pentium        SupCPU = 5
	SUP_PentiumWithMMX SupCPU = 6
	SUP_PentiumPro     SupCPU = 7
	SUP_Pentium2       SupCPU = 8
	SUP_Pentium3       SupCPU = 9
	SUP_Pentium4       SupCPU = 10
	SUP_Core           SupCPU = 11
	SUP_Core2          SupCPU = 12
	SUP_Core7          SupCPU = 13
	SUP_Itanium        SupCPU = 99
)

var intToSupCPU = map[int]SupCPU{
	0:  SUP_8086,
	1:  SUP_80186,
	2:  SUP_80286,
	3:  SUP_80386,
	4:  SUP_80486,
	5:  SUP_Pentium,
	6:  SUP_PentiumWithMMX,
	7:  SUP_PentiumPro,
	8:  SUP_Pentium2,
	9:  SUP_Pentium3,
	10: SUP_Pentium4,
	11: SUP_Core,
	12: SUP_Core2,
	13: SUP_Core7,
	99: SUP_Itanium,
}

func NewSupCPU(i int) (SupCPU, bool) {
	b, ok := intToSupCPU[i]
	return b, ok
}

var codeToSupCPU = map[string]SupCPU{
	"00": SUP_8086,
	"01": SUP_80186,
	"02": SUP_80286,
	"03": SUP_80386,
	"04": SUP_80486,
	"P1": SUP_Pentium,
	"PX": SUP_PentiumWithMMX,
	"PP": SUP_PentiumPro,
	"P2": SUP_Pentium2,
	"P3": SUP_Pentium3,
	"P4": SUP_Pentium4,
	"C1": SUP_Core,
	"C2": SUP_Core2,
	"C7": SUP_Core7,
	"IT": SUP_Itanium,
}

type SupCPURange struct {
	Start    SupCPU // 範囲の開始
	End      SupCPU // 範囲の終了
	AnyLater bool   // 任意の後続プロセッサをサポート
	LateStep bool   // 後続プロセッサの特定ステッピングのみサポート
}

func NewSupCPUByCode(code string) (*SupCPURange, bool) {
	if strings.Contains(code, "-") {
		codes := strings.Split(code, "-")
		start, ok := codeToSupCPU[codes[0]]
		if !ok {
			return nil, false
		}
		end, ok := codeToSupCPU[codes[1]]
		if !ok {
			return nil, false
		}
		return &SupCPURange{Start: start, End: end}, true
	}
	if strings.HasSuffix(code, "++") {
		startCode := strings.TrimSuffix(code, "++")
		start, ok := codeToSupCPU[startCode]
		if !ok {
			return nil, false
		}
		return &SupCPURange{Start: start, AnyLater: false, LateStep: true}, true
	}
	if strings.HasSuffix(code, "+") {
		startCode := strings.TrimSuffix(code, "+")
		start, ok := codeToSupCPU[startCode]
		if !ok {
			return nil, false
		}
		return &SupCPURange{Start: start, AnyLater: true, LateStep: false}, true
	}

	start, ok := codeToSupCPU[code]
	if !ok {
		return nil, false
	}
	return &SupCPURange{Start: start, End: start}, true
}

func (r SupCPURange) IsSupported(targetCPU SupCPU) bool {
	if r.AnyLater || r.LateStep {
		return targetCPU >= r.Start
	}
	return r.Start <= targetCPU && targetCPU <= r.End
}
