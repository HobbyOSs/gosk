package eval

import (
	"fmt"
	"github.com/hangingman/gosk/token"
	"log"
	"strconv"
	"strings"
)

type Mod int

const (
	RegReg    Mod = iota // mod=00: [レジスター+レジスター]
	RegDisp8             // mod=01: [レジスター+disp8]
	RegDisp16            // mod=10: [レジスター+disp16/32]
	Reg                  // mod=11: レジスター
)

func (c Mod) String() string {
	switch c {
	case RegReg:
		return "RegReg"
	case RegDisp8:
		return "RegDisp8"
	case RegDisp16:
		return "RegDisp16"
	case Reg:
		return "Reg"
	default:
		return "Unknown"
	}
}

// ModR/Mの条件と2bitのデータの対応
var mod2byteMap = map[Mod]string{
	RegReg: "00", RegDisp8: "01", RegDisp16: "10", Reg: "11",
}

// /0,/1.../7の条件と3bitのデータの対応
var slash2bitMap = map[string]string{
	"/0": "000",
	"/1": "001",
	"/2": "010",
	"/3": "011",
	"/4": "100",
	"/5": "101",
	"/6": "110",
	"/7": "111",
}

func getRMFromReg(srcReg string) string {
	log.Println(fmt.Sprintf("debug: Get rm reg=%s", srcReg))

	var regBits int

	switch {
	case strings.HasPrefix(srcReg, "[EAX") && strings.HasSuffix(srcReg, "]"):
		regBits = 0 // 0x0000 "000"
	case strings.HasPrefix(srcReg, "[ECX") && strings.HasSuffix(srcReg, "]"):
		regBits = 1 // 0x0000 "001"
	case strings.HasPrefix(srcReg, "[EDX") && strings.HasSuffix(srcReg, "]"):
		regBits = 2 // 0x0000 "010"
	case strings.HasPrefix(srcReg, "[EBX") && strings.HasSuffix(srcReg, "]"):
		regBits = 3 // 0x0000 "011"
	case strings.HasPrefix(srcReg, "[0x") && strings.HasSuffix(srcReg, "]"):
		regBits = 6 // 0x0000 "110"
	case strings.HasPrefix(srcReg, "[ESI") && strings.HasSuffix(srcReg, "]"):
		regBits = 6 // 0x0000 "110"
	case strings.HasPrefix(srcReg, "[SI") && strings.HasSuffix(srcReg, "]"):
		regBits = 4 // 0x00   "100"
	case strings.HasPrefix(srcReg, "[EDI") && strings.HasSuffix(srcReg, "]"):
		regBits = 7 // 0x00   "111"
	case strings.HasPrefix(srcReg, "[BX") && strings.HasSuffix(srcReg, "]"):
		regBits = 7 // 0x00   "111"
	case strings.HasPrefix(srcReg, "[DI") && strings.HasSuffix(srcReg, "]"):
		regBits = 5 // 0x00   "101"
	case IsR8(token.Token{Type: token.REGISTER, Literal: srcReg}):
		regBits = r8CodeMap[srcReg]
	case IsR16(token.Token{Type: token.REGISTER, Literal: srcReg}):
		regBits = r16CodeMap[srcReg]
	case IsR32(token.Token{Type: token.REGISTER, Literal: srcReg}):
		regBits = r32CodeMap[srcReg]
	case IsSreg(token.Token{Type: token.SEG_REGISTER, Literal: srcReg}):
		regBits = sregCodeMap[srcReg]
	case IsCtl(token.Token{Type: token.CTL_REGISTER, Literal: srcReg}):
		regBits = 0 // 資料がない。。。
	default:
		// 当てはまらなければ 110 で
		regBits = 6 // 0x0000 "110"
	}

	ans := fmt.Sprintf("%03b", regBits)
	log.Println(fmt.Sprintf("debug: Get rm = %s", ans))
	return ans
}

func concatNoSwap(m Mod, dstReg string, srcReg string) string {
	modrm := mod2byteMap[m]       // [mod]
	modrm += getRMFromReg(dstReg) // [reg]
	modrm += getRMFromReg(srcReg) // [r/m]
	return modrm
}

func concatSwap(m Mod, dstReg string, srcReg string) string {
	modrm := mod2byteMap[m]       // [mod]
	modrm += getRMFromReg(srcReg) // [r/m]
	modrm += getRMFromReg(dstReg) // [reg]
	return modrm
}

// generateModRM オペコードと２つのレジスタについてModR/Mを作成する
// 仕様書に '/r' の形式でModR/Mを求められる場合に使用する
func generateModRMSlashR(opcode byte, m Mod, dstReg string, srcReg string, forceSwap bool) byte {
	log.Println(fmt.Sprintf("info: ModR/M /r opcode=%x type=%s dst=%s src=%s", opcode, m, dstReg, srcReg))
	//
	// Generate ModR/M byte with arguments
	// [mod] 2bit
	// [reg] 3bit
	// [r/m] 3bit
	//
	// example) MOV DS, AX (mov dst, src)
	//          [mod] = Reg
	//          [reg] = DS(011)
	//          [r/m] = AX(000)
	var srcRM string = getRMFromReg(srcReg)
	var srcHasOperator = strings.Contains(srcReg, "+")
	var modrm string = ""

	if forceSwap {
		// 無理くり逆転させたい場合使う
		modrm = concatSwap(m, dstReg, srcReg)
	} else {
		if srcRM[2:3] == "0" || srcHasOperator {
			modrm = concatNoSwap(m, dstReg, srcReg)
		} else {
			modrm = concatSwap(m, dstReg, srcReg)
		}
	}

	i, _ := strconv.ParseUint(modrm, 2, 0)
	log.Println(fmt.Sprintf("info: ModR/M => %s(%x)", modrm, i))
	return byte(i)
}

// generateModRM オペコードと２つのレジスタについてModR/Mを作成する
// 仕様書に '/r' の形式でModR/Mを求められる場合に使用する
// @param opcode   オペコード
// @param m        どの形式の命令か
// @param dstReg   宛先のレジスタ
// @param regField '/0', '/1' ... '/7' までの文字列
func generateModRMSlashN(opcode byte, m Mod, dstReg string, regField string) byte {
	log.Println(fmt.Sprintf("info: ModR/M %s opcode=%x type=%s dst=%s", regField, opcode, m, dstReg))
	//
	// Generate ModR/M byte with arguments
	// [mod] 2bit
	// [reg] 3bit
	// [r/m] 3bit
	//
	modrm := mod2byteMap[m]         // [mod]
	modrm += slash2bitMap[regField] // [reg]
	modrm += getRMFromReg(dstReg)   // [r/m]

	i, _ := strconv.ParseUint(modrm, 2, 0)
	log.Println(fmt.Sprintf("info: ModR/M => %s(%x)", modrm, i))
	return byte(i)
}
