package eval

import (
	"fmt"
	"github.com/hangingman/gosk/ast"
	"github.com/hangingman/gosk/object"
	"log"
	"strings"
)

// 複雑なオペコードはファイルを分割する方針で
func evalOUTStatement(stmt *ast.MnemonicStatement) object.Object {
	toks := stmt.Name.Tokens
	bin := []byte{}

	switch {

	case toks[2].Literal == "AL" && IsImm8(toks[1]):
		// OUT imm8, AL
		log.Println(fmt.Sprintf("info: OUT imm8 (%s), %s", toks[2], toks[1]))
		bin = []byte{} // 0xe6 ib
		bin = append(bin, 0xe6)
		bin = append(bin, imm8ToByte(toks[1])...)
	case toks[2].Literal == "AX" && IsImm16(toks[1]):
		// OUT imm16, AX
		log.Println(fmt.Sprintf("info: OUT imm16 (%s), %s", toks[2], toks[1]))
		bin = []byte{} // 0xe7 iw
		bin = append(bin, 0xe7)
		bin = append(bin, imm16ToWord(toks[1])...)
	case toks[2].Literal == "EAX" && IsImm32(toks[1]):
		// OUT imm32, EAX
		log.Println(fmt.Sprintf("info: OUT imm32 (%s), %s", toks[2], toks[1]))
		bin = []byte{} // 0xe7 id
		bin = append(bin, 0xe7)
		bin = append(bin, imm32ToDword(toks[1])...)
	case toks[1].Literal == "DX" && toks[2].Literal == "AL":
		// OUT DX, AL
		log.Println(fmt.Sprintf("info: OUT %s, %s", toks[1], toks[2]))
		bin = []byte{} // 0xee
		bin = append(bin, 0xee)
	case toks[1].Literal == "DX" && toks[2].Literal == "AX":
		// OUT DX, AX
		log.Println(fmt.Sprintf("info: OUT %s, %s", toks[1], toks[2]))
		bin = []byte{} // 0xef
		bin = append(bin, 0xef)
	case toks[1].Literal == "DX" && toks[2].Literal == "EAX":
		// OUT DX, EAX
		log.Println(fmt.Sprintf("info: OUT %s, %s", toks[1], toks[2]))
		bin = []byte{} // 0xef
		bin = append(bin, 0xef)

	}

	tokStrArray := []string{}
	for _, tok := range toks {
		tokStrArray = append(tokStrArray, fmt.Sprintf("%s: %s", tok.Type, tok.Literal))
	}

	log.Println(fmt.Sprintf("info: [%s]", strings.Join(tokStrArray, ", ")))
	stmt.Bin = &object.Binary{Value: bin}
	return stmt.Bin
}
