package test

import (
	"log"
	"os"
	"testing"

	"github.com/HobbyOSs/gosk/internal/frontend"
	"github.com/HobbyOSs/gosk/internal/gen"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/suite"
)

type Day01Suite struct {
	suite.Suite
}

func (s *Day01Suite) TestHelloos1() {

	code := `	DB	0xeb, 0x4e, 0x90, 0x48, 0x45, 0x4c, 0x4c, 0x4f
	DB	0x49, 0x50, 0x4c, 0x00, 0x02, 0x01, 0x01, 0x00
	DB	0x02, 0xe0, 0x00, 0x40, 0x0b, 0xf0, 0x09, 0x00
	DB	0x12, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00
	DB	0x40, 0x0b, 0x00, 0x00, 0x00, 0x00, 0x29, 0xff
	DB	0xff, 0xff, 0xff, 0x48, 0x45, 0x4c, 0x4c, 0x4f
	DB	0x2d, 0x4f, 0x53, 0x20, 0x20, 0x20, 0x46, 0x41
	DB	0x54, 0x31, 0x32, 0x20, 0x20, 0x20, 0x00, 0x00
	RESB	16
	DB	0xb8, 0x00, 0x00, 0x8e, 0xd0, 0xbc, 0x00, 0x7c
	DB	0x8e, 0xd8, 0x8e, 0xc0, 0xbe, 0x74, 0x7c, 0x8a
	DB	0x04, 0x83, 0xc6, 0x01, 0x3c, 0x00, 0x74, 0x09
	DB	0xb4, 0x0e, 0xbb, 0x0f, 0x00, 0xcd, 0x10, 0xeb
	DB	0xee, 0xf4, 0xeb, 0xfd, 0x0a, 0x0a, 0x68, 0x65
	DB	0x6c, 0x6c, 0x6f, 0x2c, 0x20, 0x77, 0x6f, 0x72
	DB	0x6c, 0x64, 0x0a, 0x00, 0x00, 0x00, 0x00, 0x00
	RESB	368
	DB	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x55, 0xaa
	DB	0xf0, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00
	RESB	4600
	DB	0xf0, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00
	RESB	1469432
`

	temp, err := os.CreateTemp("", "helloos1.img")
	s.Require().NoError(err)
	defer os.Remove(temp.Name()) // clean up

	pt, err := gen.Parse("", []byte(code), gen.Entrypoint("Program"))
	s.Require().NoError(err)
	pass1, _ := frontend.Exec(pt, temp.Name())

	actual, err := ReadFileAsBytes(temp.Name())
	s.Require().NoError(err)
	s.Assert().Equal(int32(1474560), pass1.LOC)

	expected := defineHEX([]string{
		"DATA 0xeb 0x4e 0x90 0x48 0x45 0x4c 0x4c 0x4f",
		"DATA 0x49 0x50 0x4c 0x00 0x02 0x01 0x01 0x00",
		"DATA 0x02 0xe0 0x00 0x40 0x0b 0xf0 0x09 0x00",
		"DATA 0x12 0x00 0x02 0x00 0x00 0x00 0x00 0x00",
		"DATA 0x40 0x0b 0x00 0x00 0x00 0x00 0x29 0xff",
		"DATA 0xff 0xff 0xff 0x48 0x45 0x4c 0x4c 0x4f",
		"DATA 0x2d 0x4f 0x53 0x20 0x20 0x20 0x46 0x41",
		"DATA 0x54 0x31 0x32 0x20 0x20 0x20 0x00 0x00",
		"FILL 16",
		"DATA 0xb8 0x00 0x00 0x8e 0xd0 0xbc 0x00 0x7c",
		"DATA 0x8e 0xd8 0x8e 0xc0 0xbe 0x74 0x7c 0x8a",
		"DATA 0x04 0x83 0xc6 0x01 0x3c 0x00 0x74 0x09",
		"DATA 0xb4 0x0e 0xbb 0x0f 0x00 0xcd 0x10 0xeb",
		"DATA 0xee 0xf4 0xeb 0xfd 0x0a 0x0a 0x68 0x65",
		"DATA 0x6c 0x6c 0x6f 0x2c 0x20 0x77 0x6f 0x72",
		"DATA 0x6c 0x64 0x0a 0x00 0x00 0x00 0x00 0x00",
		"FILL 368",
		"DATA 0x00 0x00 0x00 0x00 0x00 0x00 0x55 0xaa",
		"DATA 0xf0 0xff 0xff 0x00 0x00 0x00 0x00 0x00",
		"FILL 4600",
		"DATA 0xf0 0xff 0xff 0x00 0x00 0x00 0x00 0x00",
		"FILL 1469432",
	})

	s.Assert().Equal(len(expected), len(actual))
	if diff := cmp.Diff(expected, actual); diff != "" {
		log.Printf("error: result mismatch:\n%s", DumpDiff(expected, actual, false))
	}
}

func (s *Day01Suite) TestHelloos2() {

	code := `; hello-os
; TAB=4

; 以下は標準的なFAT12フォーマットフロッピーディスクのための記述

		DB		0xeb, 0x4e, 0x90
		DB		"HELLOIPL"		; ブートセクタの名前を自由に書いてよい（8バイト）
		DW		512				; 1セクタの大きさ（512にしなければいけない）
		DB		1				; クラスタの大きさ（1セクタにしなければいけない）
		DW		1				; FATがどこから始まるか（普通は1セクタ目からにする）
		DB		2				; FATの個数（2にしなければいけない）
		DW		224				; ルートディレクトリ領域の大きさ（普通は224エントリにする）
		DW		2880			; このドライブの大きさ（2880セクタにしなければいけない）
		DB		0xf0			; メディアのタイプ（0xf0にしなければいけない）
		DW		9				; FAT領域の長さ（9セクタにしなければいけない）
		DW		18				; 1トラックにいくつのセクタがあるか（18にしなければいけない）
		DW		2				; ヘッドの数（2にしなければいけない）
		DD		0				; パーティションを使ってないのでここは必ず0
		DD		2880			; このドライブ大きさをもう一度書く
		DB		0,0,0x29		; よくわからないけどこの値にしておくといいらしい
		DD		0xffffffff		; たぶんボリュームシリアル番号
		DB		"HELLO-OS   "	; ディスクの名前（11バイト）
		DB		"FAT12   "		; フォーマットの名前（8バイト）
		RESB	18				; とりあえず18バイトあけておく

; プログラム本体

		DB		0xb8, 0x00, 0x00, 0x8e, 0xd0, 0xbc, 0x00, 0x7c
		DB		0x8e, 0xd8, 0x8e, 0xc0, 0xbe, 0x74, 0x7c, 0x8a
		DB		0x04, 0x83, 0xc6, 0x01, 0x3c, 0x00, 0x74, 0x09
		DB		0xb4, 0x0e, 0xbb, 0x0f, 0x00, 0xcd, 0x10, 0xeb
		DB		0xee, 0xf4, 0xeb, 0xfd

; メッセージ部分

		DB		0x0a, 0x0a		; 改行を2つ
		DB		"hello, world"
		DB		0x0a			; 改行
		DB		0

		RESB	0x1fe-$			; 0x001feまでを0x00で埋める命令

		DB		0x55, 0xaa

; 以下はブートセクタ以外の部分の記述

		DB		0xf0, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00
		RESB	4600
		DB		0xf0, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00
		RESB	1469432
`

	temp, err := os.CreateTemp("", "helloos2.img")
	s.Require().NoError(err)
	defer os.Remove(temp.Name()) // clean up

	pt, err := gen.Parse("", []byte(code), gen.Entrypoint("Program"))
	s.Require().NoError(err)
	pass1, _ := frontend.Exec(pt, temp.Name())

	actual, err := ReadFileAsBytes(temp.Name())
	s.Require().NoError(err)
	s.Assert().Equal(int32(1474692), pass1.LOC)

	expected := defineHEX([]string{
		"DATA 0xeb  0x4e  0x90",
		"DATA 0x48  0x45  0x4c  0x4c  0x4f  0x49  0x50  0x4c",
		"DATA 0x00  0x02",
		"DATA 0x01",
		"DATA 0x01  0x00",
		"DATA 0x02",
		"DATA 0xe0  0x00",
		"DATA 0x40  0x0b",
		"DATA 0xf0",
		"DATA 0x09  0x00",
		"DATA 0x12  0x00",
		"DATA 0x02  0x00",
		"DATA 0x00  0x00  0x00  0x00",
		"DATA 0x40  0x0b  0x00  0x00",
		"DATA 0x00  0x00  0x29",
		"DATA 0xff  0xff  0xff  0xff",
		"DATA 0x48  0x45  0x4c  0x4c  0x4f  0x2d  0x4f  0x53  0x20  0x20  0x20",
		"DATA 0x46  0x41  0x54  0x31  0x32  0x20  0x20  0x20",
		"FILL 18",

		"DATA 0xb8  0x00  0x00  0x8e  0xd0  0xbc  0x00  0x7c",
		"DATA 0x8e  0xd8  0x8e  0xc0  0xbe  0x74  0x7c  0x8a",
		"DATA 0x04  0x83  0xc6  0x01  0x3c  0x00  0x74  0x09",
		"DATA 0xb4  0x0e  0xbb  0x0f  0x00  0xcd  0x10  0xeb",
		"DATA 0xee  0xf4  0xeb  0xfd",

		"DATA 0x0a  0x0a",
		"DATA 0x68  0x65  0x6c  0x6c  0x6f  0x2c  0x20  0x77  0x6f  0x72  0x6c  0x64",
		"DATA 0x0a",
		"DATA 0x00",

		"FILL 378",
		"DATA 0x55  0xaa",
		"DATA 0xf0  0xff  0xff  0x00  0x00  0x00  0x00  0x00",
		"FILL 4600",
		"DATA 0xf0  0xff  0xff  0x00  0x00  0x00  0x00  0x00",
		"FILL 1469432",
	})

	s.Assert().Equal(len(expected), len(actual))
	if diff := cmp.Diff(expected, actual); diff != "" {
		log.Printf("error: result mismatch:\n%s", DumpDiff(expected, actual, false))
	}
}

func TestDay01Suite(t *testing.T) {
	suite.Run(t, new(Day01Suite))
}

func (s *Day01Suite) SetupSuite() {
	setUpColog(true)
	UseAnsiColorForDiff = false // Set the package-level variable
}

func (s *Day01Suite) TearDownSuite() {
}

func (s *Day01Suite) SetupTest() {
}

func (s *Day01Suite) TearDownTest() {
}
