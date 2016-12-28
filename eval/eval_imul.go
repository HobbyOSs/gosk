package eval

import (
	"fmt"
	"github.com/hangingman/gosk/ast"
	"github.com/hangingman/gosk/object"
	"log"
	"strings"
)

// 複雑なオペコードはファイルを分割する方針で
func evalIMULStatement(stmt *ast.MnemonicStatement) object.Object {
	toks := stmt.Name.Tokens
	bin := []byte{}

	switch {

	case IsR16(toks[1]) && IsImm8(toks[2]):
		// IMUL r16 , imm8
		log.Println(fmt.Sprintf("info: IMUL r16 (%s), imm8 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x6b /r ib
		bin = append(bin, 0x6b)
		bin = append(bin, generateModRMSlashR(0x6b, Reg, toks[1].Literal, toks[1].Literal, false))
		bin = append(bin, imm8ToByte(toks[2])...)
	case IsR32(toks[1]) && IsImm8(toks[2]):
		// IMUL r32, imm8
		log.Println(fmt.Sprintf("info: IMUL r32 (%s), imm8 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x6b /r ib
		bin = append(bin, 0x66)
		bin = append(bin, 0x6b)
		bin = append(bin, generateModRMSlashR(0x6b, Reg, toks[1].Literal, toks[1].Literal, false))
		bin = append(bin, imm8ToByte(toks[2])...)
	case IsR16(toks[1]) && IsImm16(toks[2]):
		// IMUL r16, imm16
		log.Println(fmt.Sprintf("info: IMUL r16 (%s), imm16 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x69 /r iw
		bin = append(bin, 0x69)
		bin = append(bin, generateModRMSlashR(0x69, Reg, toks[1].Literal, toks[1].Literal, false))
		bin = append(bin, imm16ToWord(toks[2])...)
	case IsR32(toks[1]) && IsImm32(toks[2]):
		// IMUL r32, imm32
		log.Println(fmt.Sprintf("info: IMUL r32 (%s), imm32 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x69 /r id
		bin = append(bin, 0x66)
		bin = append(bin, 0x69)
		bin = append(bin, generateModRMSlashR(0x69, Reg, toks[1].Literal, toks[1].Literal, false))
		bin = append(bin, imm32ToDword(toks[2])...)
	}

	tokStrArray := []string{}
	for _, tok := range toks {
		tokStrArray = append(tokStrArray, fmt.Sprintf("%s: %s", tok.Type, tok.Literal))
	}

	log.Println(fmt.Sprintf("info: [%s]", strings.Join(tokStrArray, ", ")))
	stmt.Bin = &object.Binary{Value: bin}
	return stmt.Bin
}
