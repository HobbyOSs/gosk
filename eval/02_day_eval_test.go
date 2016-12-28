package eval

import (
	"strings"
	"testing"
)

// TestHelloOS3 naskソース２日目(helloos3)のテスト
func TestAdd(t *testing.T) {
	input := "ADD		SI,1			; SIに1を足す"
	answer := []string{"00000000  83 c6 01                                          |...|", ""}

	testAsmSource(t, input, answer)
}

// TestHelloOS3 naskソース２日目(helloos3)のテスト
func TestMovDisp(t *testing.T) {
	input := "MOV		AL,[SI]"
	answer := []string{"00000000  8a 04                                             |..|", ""}

	testAsmSource(t, input, answer)
}

// TestHelloOS3 naskソース２日目(helloos3)のテスト
func TestCmp(t *testing.T) {
	input := "CMP		AL,0"
	answer := []string{"00000000  3c 00                                             |<.|", ""}

	testAsmSource(t, input, answer)
}

// TestHelloOS3 naskソース２日目(helloos3)のテスト
func TestDumpHelloOS3(t *testing.T) {
	input := getAsmSource("02_day_helloos3_helloos.nas")
	testAsmSourceOnlyDump(t, input, []string{""})
}

// TestHelloOS3 naskソース２日目(helloos3)のテスト
func TestHelloOS3(t *testing.T) {
	input := getAsmSource("02_day_helloos3_helloos.nas")

	// wine nask.exe helloos.nas > helloos.obj
	// hexdump -C helloos.obj > helloos.hex
	// ... generate .hex file little endian
	answer := `00000000  eb 4e 90 48 45 4c 4c 4f  49 50 4c 00 02 01 01 00  |.N.HELLOIPL.....|
00000010  02 e0 00 40 0b f0 09 00  12 00 02 00 00 00 00 00  |...@............|
00000020  40 0b 00 00 00 00 29 ff  ff ff ff 48 45 4c 4c 4f  |@.....)....HELLO|
00000030  2d 4f 53 20 20 20 46 41  54 31 32 20 20 20 00 00  |-OS   FAT12   ..|
00000040  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
00000050  b8 00 00 8e d0 bc 00 7c  8e d8 8e c0 be 74 7c 8a  |.......|.....t|.|
00000060  04 83 c6 01 3c 00 74 09  b4 0e bb 0f 00 cd 10 eb  |....<.t.........|
00000070  ee f4 eb fd 0a 0a 68 65  6c 6c 6f 2c 20 77 6f 72  |......hello, wor|
00000080  6c 64 0a 00 00 00 00 00  00 00 00 00 00 00 00 00  |ld..............|
00000090  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
*
000001f0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 55 aa  |..............U.|
00000200  f0 ff ff 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
00000210  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
*
00001400  f0 ff ff 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
00001410  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
*
00168000`

	testAsmSource(t, input, strings.Split(answer, "\n"))
}
