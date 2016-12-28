package eval

import (
	"fmt"
	"github.com/hangingman/gosk/ast"
	"github.com/hangingman/gosk/object"
	"log"
	"strings"
)

// 複雑なオペコードはファイルを分割する方針で
func evalANDStatement(stmt *ast.MnemonicStatement) object.Object {
	toks := stmt.Name.Tokens
	bin := []byte{}

	switch {

	case toks[1].Literal == "AL" && IsImm8(toks[2]):
		// AND AL, imm8
		log.Println(fmt.Sprintf("info: AND %s, imm8 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x24 ib
		bin = append(bin, 0x24)
		bin = append(bin, imm8ToByte(toks[2])...)
	case toks[1].Literal == "AX" && IsImm16(toks[2]):
		// AND AX, imm16
		log.Println(fmt.Sprintf("info: AND %s imm16 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x25 iw
		bin = append(bin, 0x25)
		bin = append(bin, imm16ToWord(toks[2])...)
	case toks[1].Literal == "EAX" && IsImm32(toks[2]):
		// AND EAX, imm32
		log.Println(fmt.Sprintf("info: AND  %s imm32 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x25 id
		bin = append(bin, 0x66)
		bin = append(bin, 0x25)
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
