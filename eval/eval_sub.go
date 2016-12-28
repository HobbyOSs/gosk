package eval

import (
	"fmt"
	"github.com/hangingman/gosk/ast"
	"github.com/hangingman/gosk/object"
	"log"
	"strings"
)

// 複雑なオペコードはファイルを分割する方針で
func evalSUBStatement(stmt *ast.MnemonicStatement) object.Object {
	toks := stmt.Name.Tokens
	bin := []byte{}

	switch {

	case IsR16(toks[1]) && IsImm8(toks[2]):
		// SUB r/m16, imm8
		log.Println(fmt.Sprintf("info: SUB r/m16 (%s), imm8 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x83 /5 ib
		bin = append(bin, 0x83)
		bin = append(bin, generateModRMSlashN(0x83, Reg, toks[1].Literal, "/5"))
		bin = append(bin, imm8ToByte(toks[2])...)
	case IsR32(toks[1]) && IsImm8(toks[2]):
		// SUB r/m32, imm8
		log.Println(fmt.Sprintf("info: SUB r/m32 (%s), imm8 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x83 /5 ib
		bin = append(bin, 0x66)
		bin = append(bin, 0x83)
		bin = append(bin, generateModRMSlashN(0x83, Reg, toks[1].Literal, "/5"))
		bin = append(bin, imm8ToByte(toks[2])...)
	case IsR8(toks[1]) && IsImm8(toks[2]):
		// SUB r/m8 , imm8
		log.Println(fmt.Sprintf("info: SUB r/m8 (%s), imm8 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x80 /5 ib
		bin = append(bin, 0x80)
		bin = append(bin, generateModRMSlashN(0x80, Reg, toks[1].Literal, "/5"))
		bin = append(bin, imm8ToByte(toks[2])...)
	case IsR16(toks[1]) && IsImm16(toks[2]):
		// SUB r/m16, imm16
		log.Println(fmt.Sprintf("info: SUB r/m16 (%s), imm16 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x81 /5 iw
		bin = append(bin, 0x81)
		bin = append(bin, generateModRMSlashN(0x81, Reg, toks[1].Literal, "/5"))
		bin = append(bin, imm16ToWord(toks[2])...)
	case IsR32(toks[1]) && IsImm32(toks[2]):
		// SUB r/m32, imm32
		log.Println(fmt.Sprintf("info: SUB r/m32 (%s), imm32 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x81 /5 id
		bin = append(bin, 0x66)
		bin = append(bin, 0x81)
		bin = append(bin, generateModRMSlashN(0x81, Reg, toks[1].Literal, "/5"))
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
