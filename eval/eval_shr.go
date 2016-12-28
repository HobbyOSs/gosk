package eval

import (
	"fmt"
	"github.com/hangingman/gosk/ast"
	"github.com/hangingman/gosk/object"
	"log"
	"strings"
)

// 複雑なオペコードはファイルを分割する方針で
func evalSHRStatement(stmt *ast.MnemonicStatement) object.Object {
	toks := stmt.Name.Tokens
	bin := []byte{}

	switch {

	case IsR8(toks[1]) && IsImm8(toks[2]):
		// SHR r/m8, imm8
		log.Println(fmt.Sprintf("info: SHR %s, imm8 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0xc0 /5 ib
		bin = append(bin, 0xc0)
		bin = append(bin, generateModRMSlashN(0xc0, Reg, toks[1].Literal, "/5"))
		bin = append(bin, imm8ToByte(toks[2])...)
	case IsR16(toks[1]) && IsImm8(toks[2]):
		// SHR r/m16, imm8
		log.Println(fmt.Sprintf("info: SHR %s imm8 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0xc1 /5 ib
		bin = append(bin, 0xc1)
		bin = append(bin, generateModRMSlashN(0xc1, Reg, toks[1].Literal, "/5"))
		bin = append(bin, imm8ToByte(toks[2])...)
	case IsR32(toks[1]) && IsImm8(toks[2]):
		// SHR r/m32, imm8
		log.Println(fmt.Sprintf("info: SHR  %s imm8 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0xc1 /5 ib
		bin = append(bin, 0x66)
		bin = append(bin, 0xc1)
		bin = append(bin, generateModRMSlashN(0xc1, Reg, toks[1].Literal, "/5"))
		bin = append(bin, imm8ToByte(toks[2])...)
	}

	tokStrArray := []string{}
	for _, tok := range toks {
		tokStrArray = append(tokStrArray, fmt.Sprintf("%s: %s", tok.Type, tok.Literal))
	}

	log.Println(fmt.Sprintf("info: [%s]", strings.Join(tokStrArray, ", ")))
	stmt.Bin = &object.Binary{Value: bin}
	return stmt.Bin
}
