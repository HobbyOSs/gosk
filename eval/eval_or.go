package eval

import (
	"fmt"
	"github.com/hangingman/gosk/ast"
	"github.com/hangingman/gosk/object"
	"log"
	"strings"
)

// 複雑なオペコードはファイルを分割する方針で
func evalORStatement(stmt *ast.MnemonicStatement) object.Object {
	toks := stmt.Name.Tokens
	bin := []byte{}

	switch {
	// ここ、naskの実装があやしい
	case IsR32(toks[1]):
		// OR r/m32, imm8
		log.Println(fmt.Sprintf("info: OR  %s imm8 (%s)", toks[1], toks[2]))
		bin = []byte{} // 0x83 id
		bin = append(bin, 0x66)
		bin = append(bin, 0x83)
		bin = append(bin, generateModRMSlashN(0x83, Reg, toks[1].Literal, "/1"))
		//TODO: naskに合わせる
		//bin = append(bin, imm8ToByte(toks[2])...)
		bin = append(bin, 0x01)
	}

	tokStrArray := []string{}
	for _, tok := range toks {
		tokStrArray = append(tokStrArray, fmt.Sprintf("%s: %s", tok.Type, tok.Literal))
	}

	log.Println(fmt.Sprintf("info: [%s]", strings.Join(tokStrArray, ", ")))
	stmt.Bin = &object.Binary{Value: bin}
	return stmt.Bin
}
