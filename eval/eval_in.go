package eval

import (
	"fmt"
	"github.com/hangingman/gosk/ast"
	"github.com/hangingman/gosk/object"
	"log"
	"strings"
)

// 複雑なオペコードはファイルを分割する方針で
func evalINStatement(stmt *ast.MnemonicStatement) object.Object {
	toks := stmt.Name.Tokens
	bin := []byte{}

	switch {

	case toks[1].Literal == "AL" && IsImm8(toks[2]):
		// IN AL , imm8
		log.Println(fmt.Sprintf("info: IN %s, imm8 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0xe4 ib
		bin = append(bin, 0xe4)
		bin = append(bin, imm8ToByte(toks[2])...)
	case toks[1].Literal == "AX" && IsImm8(toks[2]):
		// IN AX , imm8
		log.Println(fmt.Sprintf("info: IN %s, imm8 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0xe5 ib
		bin = append(bin, 0xe5)
		bin = append(bin, imm8ToByte(toks[2])...)
	case toks[1].Literal == "EAX" && IsImm8(toks[2]):
		// IN EAX , imm8
		log.Println(fmt.Sprintf("info: IN %s, imm8 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0xe5 id
		bin = append(bin, 0xe5)
		bin = append(bin, imm8ToByte(toks[2])...)
	case toks[1].Literal == "AL" && toks[2].Literal == "DX":
		// IN AL , DX
		log.Println(fmt.Sprintf("info: IN %s, %s", toks[1], toks[2]))
		bin = []byte{} // 0xec
		bin = append(bin, 0xec)
	case toks[1].Literal == "AX" && toks[2].Literal == "DX":
		// IN AX , DX
		log.Println(fmt.Sprintf("info: IN %s, %s", toks[1], toks[2]))
		bin = []byte{} // 0xed
		bin = append(bin, 0xed)
	case toks[1].Literal == "EAX" && toks[2].Literal == "DX":
		// IN EAX , DX
		log.Println(fmt.Sprintf("info: IN %s, %s", toks[1], toks[2]))
		bin = []byte{} // 0xed
		bin = append(bin, 0xed)
	}

	tokStrArray := []string{}
	for _, tok := range toks {
		tokStrArray = append(tokStrArray, fmt.Sprintf("%s: %s", tok.Type, tok.Literal))
	}

	log.Println(fmt.Sprintf("info: [%s]", strings.Join(tokStrArray, ", ")))
	stmt.Bin = &object.Binary{Value: bin}
	return stmt.Bin
}
