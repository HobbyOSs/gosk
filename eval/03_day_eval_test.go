package eval

import (
	"strings"
	"testing"
	"time"
)

func TestDumpHarib(t *testing.T) {
	// input := getAsmSource("03_day_harib00i_asmhead.nas")
	// testAsmSourceOnlyDump(t, input, []string{""})

	// 実際のテスト
	t.Run("harib00a", testHarib00a)
	t.Run("harib00b", testHarib00b)
	t.Run("harib00c", testHarib00c)
	t.Run("harib00d", testHarib00d)
	t.Run("harib00e", testHarib00e)
	t.Run("harib00f", testHarib00f)
	t.Run("harib00g", testHarib00g)
	t.Run("harib00h", testHarib00h)
	t.Run("harib00i", testHarib00i)
}

// TestHelloOS3 naskソース３日目(harib00a)のテスト
func testHarib00a(t *testing.T) {
	time.Sleep(1 * time.Second)
	input := getAsmSource("03_day_harib00a_ipl.nas")

	// wine nask.exe ipl.nas ipl.obj
	// hexdump -C ipl.obj > ipl.hex
	// ... generate .hex file little endian
	answer := `00000000  eb 4e 90 48 41 52 49 42  4f 54 45 00 02 01 01 00  |.N.HARIBOTE.....|
00000010  02 e0 00 40 0b f0 09 00  12 00 02 00 00 00 00 00  |...@............|
00000020  40 0b 00 00 00 00 29 ff  ff ff ff 48 41 52 49 42  |@.....)....HARIB|
00000030  4f 54 45 4f 53 20 46 41  54 31 32 20 20 20 00 00  |OTEOS FAT12   ..|
00000040  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
00000050  b8 00 00 8e d0 bc 00 7c  8e d8 b8 20 08 8e c0 b5  |.......|... ....|
00000060  00 b6 00 b1 02 b4 02 b0  01 bb 00 00 b2 00 cd 13  |................|
00000070  72 03 f4 eb fd be 8a 7c  8a 04 83 c6 01 3c 00 74  |r......|.....<.t|
00000080  f1 b4 0e bb 0f 00 cd 10  eb ee 0a 0a 6c 6f 61 64  |............load|
00000090  20 65 72 72 6f 72 0a 00  00 00 00 00 00 00 00 00  | error..........|
000000a0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
*
000001f0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 55 aa  |..............U.|
` // 00000200 ← hexdumpだと左の行が出るが、Goのdumpでは出ない。動作に支障はないので無視する。

	testAsmSource(t, input, strings.Split(answer, "\n"))
}

// TestHelloOS3 naskソース３日目(harib00b)のテスト
func testHarib00b(t *testing.T) {
	time.Sleep(1 * time.Second)
	input := getAsmSource("03_day_harib00b_ipl.nas")

	// wine nask.exe ipl.nas ipl.obj
	// hexdump -C ipl.obj > ipl.hex
	// ... generate .hex file little endian
	answer := `00000000  eb 4e 90 48 41 52 49 42  4f 54 45 00 02 01 01 00  |.N.HARIBOTE.....|
00000010  02 e0 00 40 0b f0 09 00  12 00 02 00 00 00 00 00  |...@............|
00000020  40 0b 00 00 00 00 29 ff  ff ff ff 48 41 52 49 42  |@.....)....HARIB|
00000030  4f 54 45 4f 53 20 46 41  54 31 32 20 20 20 00 00  |OTEOS FAT12   ..|
00000040  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
00000050  b8 00 00 8e d0 bc 00 7c  8e d8 b8 20 08 8e c0 b5  |.......|... ....|
00000060  00 b6 00 b1 02 be 00 00  b4 02 b0 01 bb 00 00 b2  |................|
00000070  00 cd 13 73 10 83 c6 01  83 fe 05 73 0b b4 00 b2  |...s.......s....|
00000080  00 cd 13 eb e3 f4 eb fd  be 9d 7c 8a 04 83 c6 01  |..........|.....|
00000090  3c 00 74 f1 b4 0e bb 0f  00 cd 10 eb ee 0a 0a 6c  |<.t............l|
000000a0  6f 61 64 20 65 72 72 6f  72 0a 00 00 00 00 00 00  |oad error.......|
000000b0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
*
000001f0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 55 aa  |..............U.|
`
	testAsmSource(t, input, strings.Split(answer, "\n"))
}

// TestHelloOS3 naskソース３日目(harib00c)のテスト
func testHarib00c(t *testing.T) {
	time.Sleep(1 * time.Second)
	input := getAsmSource("03_day_harib00c_ipl.nas")

	// wine nask.exe ipl.nas ipl.obj
	// hexdump -C ipl.obj > ipl.hex
	// ... generate .hex file little endian
	answer := `00000000  eb 4e 90 48 41 52 49 42  4f 54 45 00 02 01 01 00  |.N.HARIBOTE.....|
00000010  02 e0 00 40 0b f0 09 00  12 00 02 00 00 00 00 00  |...@............|
00000020  40 0b 00 00 00 00 29 ff  ff ff ff 48 41 52 49 42  |@.....)....HARIB|
00000030  4f 54 45 4f 53 20 46 41  54 31 32 20 20 20 00 00  |OTEOS FAT12   ..|
00000040  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
00000050  b8 00 00 8e d0 bc 00 7c  8e d8 b8 20 08 8e c0 b5  |.......|... ....|
00000060  00 b6 00 b1 02 be 00 00  b4 02 b0 01 bb 00 00 b2  |................|
00000070  00 cd 13 73 10 83 c6 01  83 fe 05 73 1a b4 00 b2  |...s.......s....|
00000080  00 cd 13 eb e3 8c c0 05  20 00 8e c0 80 c1 01 80  |........ .......|
00000090  f9 12 76 d1 f4 eb fd be  ac 7c 8a 04 83 c6 01 3c  |..v......|.....<|
000000a0  00 74 f1 b4 0e bb 0f 00  cd 10 eb ee 0a 0a 6c 6f  |.t............lo|
000000b0  61 64 20 65 72 72 6f 72  0a 00 00 00 00 00 00 00  |ad error........|
000000c0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
*
000001f0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 55 aa  |..............U.|
`

	testAsmSource(t, input, strings.Split(answer, "\n"))
}

// TestHelloOS3 naskソース３日目(harib00d)のテスト
func testHarib00d(t *testing.T) {
	time.Sleep(1 * time.Second)
	input := getAsmSource("03_day_harib00d_ipl.nas")

	// wine nask.exe ipl.nas ipl.obj
	// hexdump -C ipl.obj > ipl.hex
	// ... generate .hex file little endian
	answer := `00000000  eb 4e 90 48 41 52 49 42  4f 54 45 00 02 01 01 00  |.N.HARIBOTE.....|
00000010  02 e0 00 40 0b f0 09 00  12 00 02 00 00 00 00 00  |...@............|
00000020  40 0b 00 00 00 00 29 ff  ff ff ff 48 41 52 49 42  |@.....)....HARIB|
00000030  4f 54 45 4f 53 20 46 41  54 31 32 20 20 20 00 00  |OTEOS FAT12   ..|
00000040  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
00000050  b8 00 00 8e d0 bc 00 7c  8e d8 b8 20 08 8e c0 b5  |.......|... ....|
00000060  00 b6 00 b1 02 be 00 00  b4 02 b0 01 bb 00 00 b2  |................|
00000070  00 cd 13 73 10 83 c6 01  83 fe 05 73 2e b4 00 b2  |...s.......s....|
00000080  00 cd 13 eb e3 8c c0 05  20 00 8e c0 80 c1 01 80  |........ .......|
00000090  f9 12 76 d1 b1 01 80 c6  01 80 fe 02 72 c7 b6 00  |..v.........r...|
000000a0  80 c5 01 80 fd 0a 72 bd  f4 eb fd be c0 7c 8a 04  |......r......|..|
000000b0  83 c6 01 3c 00 74 f1 b4  0e bb 0f 00 cd 10 eb ee  |...<.t..........|
000000c0  0a 0a 6c 6f 61 64 20 65  72 72 6f 72 0a 00 00 00  |..load error....|
000000d0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
*
000001f0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 55 aa  |..............U.|
`

	testAsmSource(t, input, strings.Split(answer, "\n"))
}

// TestHelloOS3 naskソース３日目(harib00e)のテスト
func testHarib00e(t *testing.T) {
	time.Sleep(1 * time.Second)
	input := getAsmSource("03_day_harib00e_ipl.nas")

	// wine nask.exe ipl.nas ipl.obj
	// hexdump -C ipl.obj > ipl.hex
	// ... generate .hex file little endian
	answer := `00000000  eb 4e 90 48 41 52 49 42  4f 54 45 00 02 01 01 00  |.N.HARIBOTE.....|
00000010  02 e0 00 40 0b f0 09 00  12 00 02 00 00 00 00 00  |...@............|
00000020  40 0b 00 00 00 00 29 ff  ff ff ff 48 41 52 49 42  |@.....)....HARIB|
00000030  4f 54 45 4f 53 20 46 41  54 31 32 20 20 20 00 00  |OTEOS FAT12   ..|
00000040  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
00000050  b8 00 00 8e d0 bc 00 7c  8e d8 b8 20 08 8e c0 b5  |.......|... ....|
00000060  00 b6 00 b1 02 be 00 00  b4 02 b0 01 bb 00 00 b2  |................|
00000070  00 cd 13 73 10 83 c6 01  83 fe 05 73 2e b4 00 b2  |...s.......s....|
00000080  00 cd 13 eb e3 8c c0 05  20 00 8e c0 80 c1 01 80  |........ .......|
00000090  f9 12 76 d1 b1 01 80 c6  01 80 fe 02 72 c7 b6 00  |..v.........r...|
000000a0  80 c5 01 80 fd 0a 72 bd  f4 eb fd be c0 7c 8a 04  |......r......|..|
000000b0  83 c6 01 3c 00 74 f1 b4  0e bb 0f 00 cd 10 eb ee  |...<.t..........|
000000c0  0a 0a 6c 6f 61 64 20 65  72 72 6f 72 0a 00 00 00  |..load error....|
000000d0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
*
000001f0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 55 aa  |..............U.|
`

	testAsmSource(t, input, strings.Split(answer, "\n"))
}

// TestHelloOS3 naskソース３日目(harib00f)のテスト
func testHarib00f(t *testing.T) {
	time.Sleep(1 * time.Second)
	input1 := getAsmSource("03_day_harib00f_ipl.nas")

	// wine nask.exe ipl.nas ipl.obj
	// hexdump -C ipl.obj > ipl.hex
	// ... generate .hex file little endian
	answer1 := `00000000  eb 4e 90 48 41 52 49 42  4f 54 45 00 02 01 01 00  |.N.HARIBOTE.....|
00000010  02 e0 00 40 0b f0 09 00  12 00 02 00 00 00 00 00  |...@............|
00000020  40 0b 00 00 00 00 29 ff  ff ff ff 48 41 52 49 42  |@.....)....HARIB|
00000030  4f 54 45 4f 53 20 46 41  54 31 32 20 20 20 00 00  |OTEOS FAT12   ..|
00000040  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
00000050  b8 00 00 8e d0 bc 00 7c  8e d8 b8 20 08 8e c0 b5  |.......|... ....|
00000060  00 b6 00 b1 02 be 00 00  b4 02 b0 01 bb 00 00 b2  |................|
00000070  00 cd 13 73 10 83 c6 01  83 fe 05 73 2e b4 00 b2  |...s.......s....|
00000080  00 cd 13 eb e3 8c c0 05  20 00 8e c0 80 c1 01 80  |........ .......|
00000090  f9 12 76 d1 b1 01 80 c6  01 80 fe 02 72 c7 b6 00  |..v.........r...|
000000a0  80 c5 01 80 fd 0a 72 bd  e9 55 45 be c3 7c 8a 04  |......r..UE..|..|
000000b0  83 c6 01 3c 00 74 09 b4  0e bb 0f 00 cd 10 eb ee  |...<.t..........|
000000c0  f4 eb fd 0a 0a 6c 6f 61  64 20 65 72 72 6f 72 0a  |.....load error.|
000000d0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
*
000001f0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 55 aa  |..............U.|
`

	testAsmSource(t, input1, strings.Split(answer1, "\n"))

	time.Sleep(1 * time.Second)
	input2 := getAsmSource("03_day_harib00f_haribote.nas")
	answer2 := `00000000  f4 eb fd                                          |...|
`
	testAsmSource(t, input2, strings.Split(answer2, "\n"))

}

// TestHelloOS3 naskソース３日目(harib00g)のテスト
func testHarib00g(t *testing.T) {
	time.Sleep(1 * time.Second)
	input1 := getAsmSource("03_day_harib00g_ipl10.nas")

	// wine nask.exe ipl.nas ipl.obj
	// hexdump -C ipl.obj > ipl.hex
	// ... generate .hex file little endian
	answer1 := `00000000  eb 4e 90 48 41 52 49 42  4f 54 45 00 02 01 01 00  |.N.HARIBOTE.....|
00000010  02 e0 00 40 0b f0 09 00  12 00 02 00 00 00 00 00  |...@............|
00000020  40 0b 00 00 00 00 29 ff  ff ff ff 48 41 52 49 42  |@.....)....HARIB|
00000030  4f 54 45 4f 53 20 46 41  54 31 32 20 20 20 00 00  |OTEOS FAT12   ..|
00000040  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
00000050  b8 00 00 8e d0 bc 00 7c  8e d8 b8 20 08 8e c0 b5  |.......|... ....|
00000060  00 b6 00 b1 02 be 00 00  b4 02 b0 01 bb 00 00 b2  |................|
00000070  00 cd 13 73 10 83 c6 01  83 fe 05 73 32 b4 00 b2  |...s.......s2...|
00000080  00 cd 13 eb e3 8c c0 05  20 00 8e c0 80 c1 01 80  |........ .......|
00000090  f9 12 76 d1 b1 01 80 c6  01 80 fe 02 72 c7 b6 00  |..v.........r...|
000000a0  80 c5 01 80 fd 0a 72 bd  88 2e f0 0f e9 51 45 be  |......r......QE.|
000000b0  c7 7c 8a 04 83 c6 01 3c  00 74 09 b4 0e bb 0f 00  |.|.....<.t......|
000000c0  cd 10 eb ee f4 eb fd 0a  0a 6c 6f 61 64 20 65 72  |.........load er|
000000d0  72 6f 72 0a 00 00 00 00  00 00 00 00 00 00 00 00  |ror.............|
000000e0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
*
000001f0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 55 aa  |..............U.|
`

	testAsmSource(t, input1, strings.Split(answer1, "\n"))

	time.Sleep(1 * time.Second)
	input2 := getAsmSource("03_day_harib00g_haribote.nas")
	answer2 := `00000000  b0 13 b4 00 cd 10 f4 eb  fd                       |.........|
`
	testAsmSource(t, input2, strings.Split(answer2, "\n"))

}

// TestHelloOS3 naskソース３日目(harib00h)のテスト
func testHarib00h(t *testing.T) {
	time.Sleep(1 * time.Second)
	input1 := getAsmSource("03_day_harib00h_ipl10.nas")

	// wine nask.exe ipl.nas ipl.obj
	// hexdump -C ipl.obj > ipl.hex
	// ... generate .hex file little endian
	answer1 := `00000000  eb 4e 90 48 41 52 49 42  4f 54 45 00 02 01 01 00  |.N.HARIBOTE.....|
00000010  02 e0 00 40 0b f0 09 00  12 00 02 00 00 00 00 00  |...@............|
00000020  40 0b 00 00 00 00 29 ff  ff ff ff 48 41 52 49 42  |@.....)....HARIB|
00000030  4f 54 45 4f 53 20 46 41  54 31 32 20 20 20 00 00  |OTEOS FAT12   ..|
00000040  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
00000050  b8 00 00 8e d0 bc 00 7c  8e d8 b8 20 08 8e c0 b5  |.......|... ....|
00000060  00 b6 00 b1 02 be 00 00  b4 02 b0 01 bb 00 00 b2  |................|
00000070  00 cd 13 73 10 83 c6 01  83 fe 05 73 32 b4 00 b2  |...s.......s2...|
00000080  00 cd 13 eb e3 8c c0 05  20 00 8e c0 80 c1 01 80  |........ .......|
00000090  f9 12 76 d1 b1 01 80 c6  01 80 fe 02 72 c7 b6 00  |..v.........r...|
000000a0  80 c5 01 80 fd 0a 72 bd  88 2e f0 0f e9 51 45 be  |......r......QE.|
000000b0  c7 7c 8a 04 83 c6 01 3c  00 74 09 b4 0e bb 0f 00  |.|.....<.t......|
000000c0  cd 10 eb ee f4 eb fd 0a  0a 6c 6f 61 64 20 65 72  |.........load er|
000000d0  72 6f 72 0a 00 00 00 00  00 00 00 00 00 00 00 00  |ror.............|
000000e0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
*
000001f0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 55 aa  |..............U.|
`

	testAsmSource(t, input1, strings.Split(answer1, "\n"))

	time.Sleep(1 * time.Second)
	input2 := getAsmSource("03_day_harib00h_haribote.nas")
	answer2 := `00000000  b0 13 b4 00 cd 10 c6 06  f2 0f 08 c7 06 f4 0f 40  |...............@|
00000010  01 c7 06 f6 0f c8 00 66  c7 06 f8 0f 00 00 0a 00  |.......f........|
00000020  b4 02 cd 16 a2 f1 0f f4  eb fd                    |..........|
`
	testAsmSource(t, input2, strings.Split(answer2, "\n"))

}

// TestHelloOS3 naskソース３日目(harib00i)のテスト
func testHarib00i(t *testing.T) {
	time.Sleep(1 * time.Second)
	input1 := getAsmSource("03_day_harib00i_asmhead.nas")

	// wine nask.exe ipl.nas ipl.obj
	// hexdump -C ipl.obj > ipl.hex
	// ... generate .hex file little endian
	answer1 := `00000000  b0 13 b4 00 cd 10 c6 06  f2 0f 08 c7 06 f4 0f 40  |...............@|
00000010  01 c7 06 f6 0f c8 00 66  c7 06 f8 0f 00 00 0a 00  |.......f........|
00000020  b4 02 cd 16 a2 f1 0f b0  ff e6 21 90 e6 a1 fa e8  |..........!.....|
00000030  b5 00 b0 d1 e6 64 e8 ae  00 b0 df e6 60 e8 a7 00  |.....d......` + "`" + `...|
00000040  0f 01 16 2a c3 0f 20 c0  66 25 ff ff ff 7f 66 83  |...*.. .f%....f.|
00000050  c8 01 0f 22 c0 eb 00 b8  08 00 8e d8 8e c0 8e e0  |..."............|
00000060  8e e8 8e d0 66 be 30 c3  00 00 66 bf 00 00 28 00  |....f.0...f...(.|
00000070  66 b9 00 00 02 00 e8 75  00 66 be 00 7c 00 00 66  |f......u.f..|..f|
00000080  bf 00 00 10 00 66 b9 80  00 00 00 e8 60 00 66 be  |.....f......` + "`" + `.f.|
00000090  00 82 00 00 66 bf 00 02  10 00 66 b9 00 00 00 00  |....f.....f.....|
000000a0  8a 0e f0 0f 66 69 c9 00  12 00 00 66 81 e9 80 00  |....fi.....f....|
000000b0  00 00 e8 39 00 66 bb 00  00 28 00 67 66 8b 4b 10  |...9.f...(.gf.K.|
000000c0  66 83 c1 03 66 c1 e9 02  74 10 67 66 8b 73 14 66  |f...f...t.gf.s.f|
000000d0  01 de 67 66 8b 7b 0c e8  14 00 67 66 8b 63 0c 66  |..gf.{....gf.c.f|
000000e0  ea 1b 00 00 00 10 00 e4  64 24 02 75 fa c3 67 66  |........d$.u..gf|
000000f0  8b 06 66 83 c6 04 67 66  89 07 66 83 c7 04 66 83  |..f...gf..f...f.|
00000100  e9 01 75 ea c3 00 00 00  00 00 00 00 00 00 00 00  |..u.............|
00000110  00 00 00 00 00 00 00 00  ff ff 00 00 00 92 cf 00  |................|
00000120  ff ff 00 00 28 9a 47 00  00 00 17 00 10 c3 00 00  |....(.G.........|
`

	testAsmSource(t, input1, strings.Split(answer1, "\n"))

	time.Sleep(1 * time.Second)
	input2 := getAsmSource("03_day_harib00i_ipl10.nas")
	answer2 := `00000000  eb 4e 90 48 41 52 49 42  4f 54 45 00 02 01 01 00  |.N.HARIBOTE.....|
00000010  02 e0 00 40 0b f0 09 00  12 00 02 00 00 00 00 00  |...@............|
00000020  40 0b 00 00 00 00 29 ff  ff ff ff 48 41 52 49 42  |@.....)....HARIB|
00000030  4f 54 45 4f 53 20 46 41  54 31 32 20 20 20 00 00  |OTEOS FAT12   ..|
00000040  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
00000050  b8 00 00 8e d0 bc 00 7c  8e d8 b8 20 08 8e c0 b5  |.......|... ....|
00000060  00 b6 00 b1 02 be 00 00  b4 02 b0 01 bb 00 00 b2  |................|
00000070  00 cd 13 73 10 83 c6 01  83 fe 05 73 32 b4 00 b2  |...s.......s2...|
00000080  00 cd 13 eb e3 8c c0 05  20 00 8e c0 80 c1 01 80  |........ .......|
00000090  f9 12 76 d1 b1 01 80 c6  01 80 fe 02 72 c7 b6 00  |..v.........r...|
000000a0  80 c5 01 80 fd 0a 72 bd  88 2e f0 0f e9 51 45 be  |......r......QE.|
000000b0  c7 7c 8a 04 83 c6 01 3c  00 74 09 b4 0e bb 0f 00  |.|.....<.t......|
000000c0  cd 10 eb ee f4 eb fd 0a  0a 6c 6f 61 64 20 65 72  |.........load er|
000000d0  72 6f 72 0a 00 00 00 00  00 00 00 00 00 00 00 00  |ror.............|
000000e0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 00 00  |................|
*
000001f0  00 00 00 00 00 00 00 00  00 00 00 00 00 00 55 aa  |..............U.|
`
	testAsmSource(t, input2, strings.Split(answer2, "\n"))

}
