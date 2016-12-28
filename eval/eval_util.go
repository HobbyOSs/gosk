package eval

import (
	"encoding/hex"
	"fmt"
	"github.com/hangingman/gosk/token"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func IsNil(x interface{}) bool {
	return x == nil || reflect.ValueOf(x).IsNil()
}

func IsNotNil(x interface{}) bool {
	return !IsNil(x)
}

func int2Byte(i int) []byte {
	// TODO: ここはもう少しうまくやりたい
	if i < 0 {
		hexStr := fmt.Sprintf("%16x", uint(i))
		bs, _ := hex.DecodeString(hexStr[14:16])
		log.Println(fmt.Sprintf("info: int2Byte %s = %x", hexStr, bs[0:1]))
		return bs[0:1]
	}

	hexStr := fmt.Sprintf("%08x", uint(i))
	bs, _ := hex.DecodeString(hexStr[6:])
	log.Println(fmt.Sprintf("info: int2Byte %s = %x", hexStr, bs[0:1]))
	return bs[0:1]
}

func int2Word(i int) []byte {
	hexStr := fmt.Sprintf("%08x", uint(i))
	bs, _ := hex.DecodeString(hexStr[4:])
	// リトルエンディアンで格納する
	bs[0], bs[1] = bs[1], bs[0]
	log.Println(fmt.Sprintf("info: int2Word %s = %x", hexStr, bs[0:2]))
	return bs[0:2]
}

func int2Dword(i int) []byte {
	hexStr := fmt.Sprintf("%08x", uint(i))
	bs, _ := hex.DecodeString(hexStr)

	// リトルエンディアンで格納する
	bs[0], bs[1], bs[2], bs[3] = bs[3], bs[2], bs[1], bs[0]
	log.Println(fmt.Sprintf("info: int2Dword %s = %x", hexStr, bs[0:4]))
	return bs[0:4]
}

func isByteHex(tok token.Token) bool {
	return tok.Type == token.HEX_LIT && len(tok.Literal) == 4
}

func isWordHex(tok token.Token) bool {
	return tok.Type == token.HEX_LIT && len(tok.Literal) == 6
}

func isDwordHex(tok token.Token) bool {
	return tok.Type == token.HEX_LIT && len(tok.Literal) == 10
}

func imm8ToByte(tok token.Token) []byte {
	if tok.Type == token.HEX_LIT {
		bs, _ := hex.DecodeString(string([]rune(tok.Literal)[2:]))
		return bs
	}
	if tok.Type == token.INT {
		v, _ := strconv.Atoi(tok.Literal)
		return int2Byte(v)
	}
	return nil
}

func imm16ToWord(tok token.Token) []byte {
	if tok.Type == token.HEX_LIT {
		bs, _ := hex.DecodeString(string([]rune(tok.Literal)[2:]))
		return []byte{bs[1], bs[0]}
	}
	if tok.Type == token.INT {
		v, _ := strconv.Atoi(tok.Literal)
		return int2Word(v)
	}
	return nil
}

func imm32ToDword(tok token.Token) []byte {
	if tok.Type == token.HEX_LIT {
		tmp := fmt.Sprintf("%8v", tok.Literal[2:])
		hexStr := strings.Replace(tmp, " ", "0", -1)
		bs, _ := hex.DecodeString(string([]rune(hexStr))[:])
		// リトルエンディアンで格納する
		bs[0], bs[1], bs[2], bs[3] = bs[3], bs[2], bs[1], bs[0]
		return bs[0:4]
	}
	if tok.Type == token.INT {
		v, _ := strconv.Atoi(tok.Literal)
		return int2Dword(v)
	}
	return nil
}

func plusRb(opcode byte, register string) byte {
	// +rb
	return int2Byte(int(opcode) + r8CodeMap[register])[0]
}

func plusRw(opcode byte, register string) byte {
	// +rw
	return int2Byte(int(opcode) + r16CodeMap[register])[0]
}

func plusRd(opcode byte, register string) byte {
	// +rd
	return int2Byte(int(opcode) + r32CodeMap[register])[0]
}
