package eval

import (
	"fmt"
	"github.com/hangingman/gosk/ast"
	"github.com/hangingman/gosk/object"
	"log"
	"strings"
)

// 複雑なオペコードはファイルを分割する方針で
func evalADDStatement(stmt *ast.MnemonicStatement) object.Object {
	toks := stmt.Name.Tokens
	bin := []byte{}

	switch {

	case toks[1].Literal == "AL" && IsImm8(toks[2]):
		// ADD AL , imm8
		log.Println(fmt.Sprintf("info: ADD %s, imm8 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x04 ib
		bin = append(bin, 0x04)
		bin = append(bin, imm8ToByte(toks[2])...)
	case toks[1].Literal == "AX" && IsImm16(toks[2]):
		// ADD AX , imm16
		log.Println(fmt.Sprintf("info: ADD %s, imm16 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x05 iw
		bin = append(bin, 0x05)
		bin = append(bin, imm16ToWord(toks[2])...)
	case toks[1].Literal == "EAX" && IsImm32(toks[2]):
		// ADD EAX , imm32
		log.Println(fmt.Sprintf("info: ADD %s, imm32 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x05 id
		bin = append(bin, 0x05)
		bin = append(bin, imm32ToDword(toks[2])...)
	case IsR8(toks[1]) && IsImm8(toks[2]):
		// ADD r/m8 , imm8
		log.Println(fmt.Sprintf("info: ADD r/m8 (%s), imm8 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x80 /0 ib
		bin = append(bin, 0x80)
		bin = append(bin, generateModRMSlashN(0x80, Reg, toks[1].Literal, "/0"))
		bin = append(bin, imm8ToByte(toks[2])...)
	case IsR16(toks[1]) && IsImm8(toks[2]):
		// ADD r/m16, imm8
		log.Println(fmt.Sprintf("info: ADD r/m16 (%s), imm8 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x83 /0 ib
		bin = append(bin, 0x83)
		bin = append(bin, generateModRMSlashN(0x83, Reg, toks[1].Literal, "/0"))
		bin = append(bin, imm8ToByte(toks[2])...)
	case IsR32(toks[1]) && IsImm8(toks[2]):
		// ADD r/m32, imm8
		log.Println(fmt.Sprintf("info: ADD r/m32 (%s), imm8 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x83 /0 ib
		bin = append(bin, 0x66)
		bin = append(bin, 0x83)
		bin = append(bin, generateModRMSlashN(0x83, Reg, toks[1].Literal, "/0"))
		bin = append(bin, imm8ToByte(toks[2])...)
	case IsR16(toks[1]) && IsImm16(toks[2]):
		// ADD r/m16, imm16
		log.Println(fmt.Sprintf("info: ADD r/m16 (%s), imm16 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x81 /0 iw
		bin = append(bin, 0x81)
		bin = append(bin, generateModRMSlashN(0x81, Reg, toks[1].Literal, "/0"))
		bin = append(bin, imm16ToWord(toks[2])...)
	case IsR32(toks[1]) && IsImm32(toks[2]):
		// ADD r/m32, imm32
		log.Println(fmt.Sprintf("info: ADD r/m32 (%s), imm32 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x81 /0 id
		bin = append(bin, 0x66)
		bin = append(bin, 0x81)
		bin = append(bin, generateModRMSlashN(0x81, Reg, toks[1].Literal, "/0"))
		bin = append(bin, imm32ToDword(toks[2])...)
	case IsR8(toks[1]) && IsR8(toks[2]):
		// ADD r/m8 , r8
		log.Println(fmt.Sprintf("info: ADD r/m8 (%s), r8 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x00 /r
		bin = append(bin, 0x00)
		bin = append(bin, generateModRMSlashR(0x80, Reg, toks[1].Literal, toks[2].Literal, false))
	case IsR16(toks[1]) && IsR16(toks[2]):
		// ADD r/m16, r16
		log.Println(fmt.Sprintf("info: ADD r/m16 (%s), r16 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x01 /r
		bin = append(bin, 0x01)
		bin = append(bin, generateModRMSlashR(0x83, Reg, toks[1].Literal, toks[2].Literal, false))
	case IsR32(toks[1]) && IsR32(toks[2]):
		// ADD r/m32, r32
		log.Println(fmt.Sprintf("info: ADD r/m32 (%s), r32 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x01 /r
		bin = append(bin, 0x66)
		bin = append(bin, 0x01)
		bin = append(bin, generateModRMSlashR(0x83, Reg, toks[1].Literal, toks[2].Literal, false))
	}

	tokStrArray := []string{}
	for _, tok := range toks {
		tokStrArray = append(tokStrArray, fmt.Sprintf("%s: %s", tok.Type, tok.Literal))
	}

	log.Println(fmt.Sprintf("info: [%s]", strings.Join(tokStrArray, ", ")))
	stmt.Bin = &object.Binary{Value: bin}
	return stmt.Bin
}
