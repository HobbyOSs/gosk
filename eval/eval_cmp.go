package eval

import (
	"fmt"
	"github.com/hangingman/gosk/ast"
	"github.com/hangingman/gosk/object"
	"log"
	"strings"
)

// 複雑なオペコードはファイルを分割する方針で
func evalCMPStatement(stmt *ast.MnemonicStatement) object.Object {
	toks := stmt.Name.Tokens
	bin := []byte{}

	switch {

	case toks[1].Literal == "AL" && IsImm8(toks[2]):
		// CMP r/m8 , imm8
		log.Println(fmt.Sprintf("info: CMP %s, imm8 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x3c /0 ib
		bin = append(bin, 0x3c)
		bin = append(bin, imm8ToByte(toks[2])...)
	case toks[1].Literal == "AX" && IsImm16(toks[2]):
		// CMP r/m16, imm16
		log.Println(fmt.Sprintf("info: CMP %s, imm16 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x3d /0 iw
		bin = append(bin, 0x3d)
		bin = append(bin, imm16ToWord(toks[2])...)
	case toks[1].Literal == "EAX" && IsImm32(toks[2]):
		// CMP r/m32, imm32
		log.Println(fmt.Sprintf("info: CMP %s, imm32 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x3d /0 id
		bin = append(bin, 0x3d)
		bin = append(bin, imm32ToDword(toks[2])...)
	case IsR8(toks[1]) && IsImm8(toks[2]):
		// CMP r/m8, imm8
		log.Println(fmt.Sprintf("info: CMP %s, imm8 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x80 /7 ib
		bin = append(bin, 0x80)
		bin = append(bin, generateModRMSlashN(0x80, Reg, toks[1].Literal, "/7"))
		bin = append(bin, imm8ToByte(toks[2])...)
	case IsR16(toks[1]) && IsImm8(toks[2]):
		// CMP r/m16, imm8
		log.Println(fmt.Sprintf("info: CMP %s, imm8 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x83 /7 ib
		bin = append(bin, 0x83)
		bin = append(bin, generateModRMSlashN(0x83, Reg, toks[1].Literal, "/7"))
		bin = append(bin, imm8ToByte(toks[2])...)
	case IsR16(toks[1]) && IsImm16(toks[2]):
		// CMP r/m16, imm16
		log.Println(fmt.Sprintf("info: CMP %s, imm16 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x83 /7 ib
		bin = append(bin, 0x83)
		bin = append(bin, generateModRMSlashN(0x83, Reg, toks[1].Literal, "/7"))
		bin = append(bin, imm16ToWord(toks[2])...)
	case IsR32(toks[1]) && IsImm8(toks[2]):
		// CMP r/m32, imm8
		log.Println(fmt.Sprintf("info: CMP %s, imm8 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x83 /7 ib
		bin = append(bin, 0x83)
		bin = append(bin, generateModRMSlashN(0x83, Reg, toks[1].Literal, "/7"))
		bin = append(bin, imm8ToByte(toks[2])...)
	case IsR32(toks[1]) && IsImm32(toks[2]):
		// CMP r/m32, imm32
		log.Println(fmt.Sprintf("info: CMP %s, imm32 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x83 /7 ib
		bin = append(bin, 0x83)
		bin = append(bin, generateModRMSlashN(0x83, Reg, toks[1].Literal, "/7"))
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
