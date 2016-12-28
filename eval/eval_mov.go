package eval

import (
	"fmt"
	"github.com/hangingman/gosk/ast"
	"github.com/hangingman/gosk/object"
	"github.com/hangingman/gosk/token"
	"log"
	"strings"
)

// 複雑なオペコードはファイルを分割する方針で
func evalMOVStatement(stmt *ast.MnemonicStatement) object.Object {
	toks := stmt.Name.Tokens
	bin := &object.Binary{Value: []byte{}}

	// (1) MOV AX, 0xff
	// (2) MOV [XX], AX
	// (3) MOV AX, [XX]
	//     MOV AX, labl
	// (4) MOV Sreg, R16
	//     MOV R16, Sreg
	// (5) MOV CR, r32
	//     MOV r32, CR
	switch {
	//
	// (1) MOV r8~r32, immX
	//
	case IsR8(toks[1]) && IsImm8(toks[2]):
		// MOV r8 , imm8
		log.Println(fmt.Sprintf("info: MOV r8 (%s), imm8 (%s)", toks[1], toks[2]))
		bin.Value = append(bin.Value, plusRb(0xb0, toks[1].Literal))
		bin.Value = append(bin.Value, imm8ToByte(toks[2])...)
	case IsR16(toks[1]) && IsImm16(toks[2]):
		// MOV r16, imm16
		log.Println(fmt.Sprintf("info: MOV r16 (%s), imm16 (%s)", toks[1], toks[2]))
		// 0xB8+rw
		bin.Value = append(bin.Value, plusRw(0xb8, toks[1].Literal))
		bin.Value = append(bin.Value, imm16ToWord(toks[2])...)
	case IsR32(toks[1]) && (IsImm16(toks[2]) || IsImm32(toks[2])):
		// MOV r32, imm32
		log.Println(fmt.Sprintf("info: MOV r32 (%s), imm32 (%s)", toks[1], toks[2]))
		// 0xB8+rd
		bin.Value = append(bin.Value, 0x66)
		bin.Value = append(bin.Value, plusRd(0xb8, toks[1].Literal))
		bin.Value = append(bin.Value, imm32ToDword(toks[2])...)

		//
		// (2) MOV moffs8~32, Acc
		//
	case toks[1].Type == token.LBRACKET && toks[2].Type != token.REGISTER && toks[3].Type == token.RBRACKET && toks[4].Literal == "AL":
		// MOV moffs8, AL
		log.Println(fmt.Sprintf("info: MOV moffs8 (%s), AL (%s)", toks[2], toks[4]))
		bin.Value = append(bin.Value, 0xa2)
		bin.Value = append(bin.Value, imm16ToWord(toks[2])...)
	case toks[1].Type == token.LBRACKET && toks[2].Type != token.REGISTER && toks[3].Type == token.RBRACKET && toks[4].Literal == "AX":
		// MOV moffs16, AX
		log.Println(fmt.Sprintf("info: MOV moffs16 (%s), AX (%s)", toks[2], toks[4]))
		bin.Value = append(bin.Value, 0xa3)
		bin.Value = append(bin.Value, imm16ToWord(toks[2])...)
	case toks[1].Type == token.LBRACKET && toks[2].Type != token.REGISTER && toks[3].Type == token.RBRACKET && toks[4].Literal == "EAX":
		// MOV moffs32, EAX
		log.Println(fmt.Sprintf("info: MOV moffs32 (%s), EAX (%s)", toks[2], toks[4]))
		bin.Value = append(bin.Value, 0x67)
		bin.Value = append(bin.Value, 0x66)
		bin.Value = append(bin.Value, 0xa3)
		bin.Value = append(bin.Value, imm16ToWord(toks[2])...)

		//
		// (2) MOV m8~m32, rX (not accumulator)
		//
	case toks[1].Type == token.LBRACKET && toks[3].Type == token.RBRACKET && IsR8(toks[4]):
		// MOV m8 , r8
		log.Println(fmt.Sprintf("info: MOV m8 (%s), r8 (%s)", toks[2], toks[4]))
		disp := "[" + toks[2].Literal + "]"
		bin.Value = append(bin.Value, 0x88)
		bin.Value = append(bin.Value, generateModRMSlashR(0x88, RegReg, disp, toks[4].Literal, false))
		bin.Value = append(bin.Value, imm16ToWord(toks[2])...)
	case toks[1].Type == token.LBRACKET && toks[3].Type == token.RBRACKET && IsR16(toks[4]):
		// MOV m16, r16
		log.Println(fmt.Sprintf("info: MOV m16 (%s), r16 (%s)", toks[2], toks[4]))
		disp := "[" + toks[2].Literal + "]"
		bin.Value = append(bin.Value, 0x89)
		bin.Value = append(bin.Value, generateModRMSlashR(0x88, RegReg, disp, toks[4].Literal, false))
		bin.Value = append(bin.Value, imm16ToWord(toks[2])...)
	case toks[1].Type == token.LBRACKET && toks[3].Type == token.RBRACKET && IsR32(toks[4]):
		// MOV m32, r32
		log.Println(fmt.Sprintf("info: MOV m32 (%s), r32 (%s)", toks[2], toks[4]))
		disp := "[" + toks[2].Literal + "]"
		bin.Value = append(bin.Value, 0x67)
		bin.Value = append(bin.Value, 0x66)
		bin.Value = append(bin.Value, 0x89)
		// FIXME: ここのmod/rmがどうしても逆になる
		bin.Value = append(bin.Value, generateModRMSlashR(0x88, RegReg, disp, toks[4].Literal, true))
		bin.Value = append(bin.Value, imm32ToDword(toks[2])...)

		//
		// (3) MOV rX, [m8~m32]
		//     MOV rX, [m8~m32+disp]
		//     MOV rX, BYTE []
		//     MOV rX, WORD []
		//     MOV rX, DWORD []
	case IsR8(toks[1]) && toks[2].Type == token.LBRACKET && toks[4].Type == token.RBRACKET:
		// MOV r8 , imm8 で immが参照（ex: [SI]）
		log.Println(fmt.Sprintf("info: MOV r8 (%s), disp8 (%s)", toks[1], toks[3]))
		disp := "[" + toks[3].Literal + "]"
		// 0x8a
		bin.Value = append(bin.Value, 0x8a)
		bin.Value = append(bin.Value, generateModRMSlashN(0x8a, RegReg, disp, "/0"))
	case IsR16(toks[1]) && toks[2].Type == token.LBRACKET && toks[4].Type == token.RBRACKET:
		// MOV r16 , imm16 で immが参照
		log.Println(fmt.Sprintf("info: MOV r16 (%s), disp16 (%s)", toks[1], toks[3]))
		disp := "[" + toks[3].Literal + "]"
		// 0x8b
		bin.Value = append(bin.Value, 0x8b)
		bin.Value = append(bin.Value, generateModRMSlashN(0x8b, RegReg, disp, "/0"))
	case IsR32(toks[1]) && toks[2].Type == token.LBRACKET && toks[4].Type == token.RBRACKET:
		// MOV r32 , imm32 で immが参照
		log.Println(fmt.Sprintf("info: MOV r32 (%s), disp32 (%s)", toks[1], toks[3]))
		disp := "[" + toks[3].Literal + "]"
		// 0x8b
		bin.Value = append(bin.Value, 0x67)
		bin.Value = append(bin.Value, 0x66)
		bin.Value = append(bin.Value, 0x8b)
		bin.Value = append(bin.Value, generateModRMSlashN(0x8b, RegReg, disp, "/0"))

	case IsR32(toks[1]) && toks[2].Type == token.LBRACKET && toks[4].IsOperator():
		// MOV r32 , [imm32 + int]
		log.Println(fmt.Sprintf("info: MOV r32 (%s), disp32 (%s) + %s", toks[1], toks[3], toks[5]))
		disp := "[" + toks[3].Literal + "+" + toks[5].Literal + "]"
		// 0x8b
		bin.Value = append(bin.Value, 0x67)
		bin.Value = append(bin.Value, 0x66)
		bin.Value = append(bin.Value, 0x8b)
		bin.Value = append(bin.Value, generateModRMSlashR(0x8b, RegDisp8, toks[1].Literal, disp, false))
		bin.Value = append(bin.Value, imm8ToByte(toks[5])...)

	case IsR8(toks[1]) && toks[2].Literal == "BYTE" && toks[3].Type == token.LBRACKET && toks[5].Type == token.RBRACKET:
		// MOV r8 , imm8 で immが参照（ex: [SI]）、データ型指定
		log.Println(fmt.Sprintf("info: MOV r8 (%s), %s disp8 (%s)", toks[1], toks[2], toks[4]))
		disp := "[" + toks[4].Literal + "]"
		// 0x8a
		bin.Value = append(bin.Value, 0x8a)
		bin.Value = append(bin.Value, generateModRMSlashR(0x8a, RegReg, toks[1].Literal, disp, false))
		bin.Value = append(bin.Value, imm16ToWord(toks[4])...)

	case IsR8(toks[1]) && toks[2].Type == token.IDENT:
		// MOV r8 , imm8 で immがラベル
		// callbackを配置し今のバイト数を設定する
		log.Println(fmt.Sprintf("info: MOV r8 (%s), imm8(label) (%s)", toks[1], toks[2]))
		// 0xB0+rb
		opcode := plusRb(0xb0, toks[1].Literal)
		bin.Value = append(bin.Value, opcode) // ダミーデータ
		bin.Value = append(bin.Value, 0x00)   // ダミーデータ
		labelManage.AddLabelCallback(
			[]byte{opcode},
			toks[2].Literal,
			bin,
			-dollarPosition,
			int2Byte,
		)
	case IsR16(toks[1]) && toks[2].Type == token.IDENT:
		// MOV r16 , imm16 で immがラベルの場合
		// callbackを配置し今のバイト数を設定する
		log.Println(fmt.Sprintf("info: MOV r16 (%s), imm16(label) (%s)", toks[1], toks[2]))
		// 0xB8+rw
		opcode := plusRw(0xb8, toks[1].Literal)
		bin.Value = append(bin.Value, opcode) // ダミーデータ
		bin.Value = append(bin.Value, 0x00)   // ダミーデータ
		bin.Value = append(bin.Value, 0x00)   // ダミーデータ
		labelManage.AddLabelCallback(
			[]byte{opcode},
			toks[2].Literal,
			bin,
			-dollarPosition,
			int2Word,
		)
	case IsR32(toks[1]) && toks[2].Type == token.IDENT:
		// MOV r32 , imm32 で immがラベルの場合
		// callbackを配置し今のバイト数を設定する
		log.Println(fmt.Sprintf("info: MOV r32 (%s), imm32(label) (%s)", toks[1], toks[2]))
		// 0xB8+rd
		opcode := plusRd(0xb8, toks[1].Literal)
		bin.Value = append(bin.Value, 0x66)   // ダミーデータ
		bin.Value = append(bin.Value, opcode) //
		bin.Value = append(bin.Value, 0x00)   //
		bin.Value = append(bin.Value, 0x00)   //
		bin.Value = append(bin.Value, 0x00)   //
		bin.Value = append(bin.Value, 0x00)   //
		labelManage.AddLabelCallback(
			[]byte{0x66, opcode},
			toks[2].Literal,
			bin,
			-dollarPosition,
			int2Dword,
		)
		// (4)
	case IsSreg(toks[1]) && IsR16(toks[2]):
		// MOV Sreg, r/m16
		log.Println(fmt.Sprintf("info: MOV Sreg (%s), r/m16 (%s)", toks[1], toks[2]))
		// 0x8E /r
		bin.Value = append(bin.Value, 0x8e)
		bin.Value = append(bin.Value, generateModRMSlashR(0x8e, Reg, toks[1].Literal, toks[2].Literal, false))
	case IsR16(toks[1]) && IsSreg(toks[2]):
		// MOV r/m16, Sreg
		log.Println(fmt.Sprintf("info: MOV r/m16 (%s), Sreg (%s)", toks[1], toks[2]))
		// 0x8C /r
		bin.Value = append(bin.Value, 0x8c)
		bin.Value = append(bin.Value, generateModRMSlashR(0x8e, Reg, toks[1].Literal, toks[2].Literal, false))
	case toks[1].Literal == "BYTE" && toks[2].Type == token.LBRACKET && toks[4].Type == token.RBRACKET:
		// MOV r/m8, imm8
		log.Println(fmt.Sprintf("info: MOV r/m8 (%s), Imm (%s)", toks[3], toks[5]))
		bin.Value = append(bin.Value, 0xc6)
		bin.Value = append(bin.Value, generateModRMSlashN(0xc6, RegReg, toks[5].Literal, "/0"))
		bin.Value = append(bin.Value, imm16ToWord(toks[3])...)
		bin.Value = append(bin.Value, imm8ToByte(toks[5])...)
	case toks[1].Literal == "WORD" && toks[2].Type == token.LBRACKET && toks[4].Type == token.RBRACKET:
		// MOV r/m16, imm16pp
		log.Println(fmt.Sprintf("info: MOV r/m16 (%s), Imm (%s)", toks[3], toks[5]))
		bin.Value = append(bin.Value, 0xc7)
		bin.Value = append(bin.Value, generateModRMSlashN(0xc7, RegReg, toks[5].Literal, "/0"))
		bin.Value = append(bin.Value, imm16ToWord(toks[3])...)
		bin.Value = append(bin.Value, imm16ToWord(toks[5])...)
	case toks[1].Literal == "DWORD" && toks[2].Type == token.LBRACKET && toks[4].Type == token.RBRACKET:
		// MOV r/m32, imm32
		log.Println(fmt.Sprintf("info: MOV r/m32 (%s), Imm (%s)", toks[3], toks[5]))
		bin.Value = append(bin.Value, 0x66)
		bin.Value = append(bin.Value, 0xc7)
		bin.Value = append(bin.Value, generateModRMSlashN(0xc7, RegReg, toks[5].Literal, "/0"))
		bin.Value = append(bin.Value, imm16ToWord(toks[3])...)
		bin.Value = append(bin.Value, imm32ToDword(toks[5])...)

	case IsCtl(toks[1]) && IsR32(toks[2]):
		// MOV CR0, r32
		log.Println(fmt.Sprintf("info: MOV CR0(%s), R32(%s)", toks[1], toks[2]))
		bin.Value = append(bin.Value, 0x0f)
		bin.Value = append(bin.Value, 0x22)
		bin.Value = append(bin.Value, generateModRMSlashR(0x0f, Reg, toks[1].Literal, toks[2].Literal, false))
	case IsR32(toks[1]) && IsCtl(toks[2]):
		// MOV r32, CR0
		log.Println(fmt.Sprintf("info: MOV R32(%s), CR0(%s)", toks[1], toks[2]))
		bin.Value = append(bin.Value, 0x0f)
		bin.Value = append(bin.Value, 0x20)
		bin.Value = append(bin.Value, generateModRMSlashR(0x0f, Reg, toks[1].Literal, toks[2].Literal, false))

	}

	tokStrArray := []string{}
	for _, tok := range toks {
		tokStrArray = append(tokStrArray, fmt.Sprintf("%s: %s", tok.Type, tok.Literal))
	}

	log.Println(fmt.Sprintf("info: [%s]", strings.Join(tokStrArray, ", ")))

	stmt.Bin = bin
	return stmt.Bin
}
