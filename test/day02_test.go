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

type Day02Suite struct {
	suite.Suite
}

func (s *Day02Suite) TestHelloos3() {
	s.T().Skip("day02の実装が完了するまでスキップ")

	code := `; hello-os
; TAB=4

		ORG		0x7c00			; このプログラムがどこに読み込まれるのか

; 以下は標準的なFAT12フォーマットフロッピーディスクのための記述

		JMP		entry
		DB		0x90
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

entry:
		MOV		AX,0			; レジスタ初期化
		MOV		SS,AX
		MOV		SP,0x7c00
		MOV		DS,AX
		MOV		ES,AX

		MOV		SI,msg
putloop:
		MOV		AL,[SI]
		ADD		SI,1			; SIに1を足す
		CMP		AL,0
		JE		fin
		MOV		AH,0x0e			; 一文字表示ファンクション
		MOV		BX,15			; カラーコード
		INT		0x10			; ビデオBIOS呼び出し
		JMP		putloop
fin:
		HLT						; 何かあるまでCPUを停止させる
		JMP		fin				; 無限ループ

msg:
		DB		0x0a, 0x0a		; 改行を2つ
		DB		"hello, world"
		DB		0x0a			; 改行
		DB		0

		RESB	0x7dfe-$		; 0x7dfeまでを0x00で埋める命令

		DB		0x55, 0xaa

; 以下はブートセクタ以外の部分の記述

		DB		0xf0, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00
		RESB	4600
		DB		0xf0, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00
		RESB	1469432
`

	temp, err := os.CreateTemp("", "helloos3.img")
	s.Require().NoError(err)
	defer os.Remove(temp.Name()) // clean up

	pt, err := gen.Parse("", []byte(code), gen.Entrypoint("Program"))
	s.Require().NoError(err)
	pass1, _ := frontend.Exec(pt, temp.Name())

	actual, err := ReadFileAsBytes(temp.Name())
	s.Require().NoError(err)
	s.Assert().Equal(int32(1474560), pass1.LOC)

	expected := defineHEX([]string{
		"DATA 0xeb 0x4e",
		"DATA 0x90",
		"DATA 0x48 0x45 0x4c 0x4c 0x4f 0x49 0x50 0x4c",
		"DATA 0x00 0x02",
		"DATA 0x01",
		"DATA 0x01 0x00",
		"DATA 0x02",
		"DATA 0xe0 0x00",
		"DATA 0x40 0x0b",
		"DATA 0xf0",
		"DATA 0x09 0x00",
		"DATA 0x12 0x00",
		"DATA 0x02 0x00",
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x40 0x0b 0x00 0x00",
		"DATA 0x00 0x00 0x29",
		"DATA 0xff 0xff 0xff 0xff",
		"DATA 0x48 0x45 0x4c 0x4c 0x4f 0x2d 0x4f 0x53 0x20 0x20 0x20",
		"DATA 0x46 0x41 0x54 0x31 0x32 0x20 0x20 0x20",
		"FILL 18",

		"DATA 0xb8 0x00 0x00",
		"DATA 0x8e 0xd0",
		"DATA 0xbc 0x00 0x7c",
		"DATA 0x8e 0xd8",
		"DATA 0x8e 0xc0",
		"DATA 0xbe 0x74 0x7c",
		"DATA 0x8a 0x04",
		"DATA 0x83 0xc6 0x01",
		"DATA 0x3c 0x00",
		"DATA 0x74 0x09",
		"DATA 0xb4 0x0e",
		"DATA 0xbb 0x0f 0x00",
		"DATA 0xcd 0x10",
		"DATA 0xeb 0xee",
		"DATA 0xf4",
		"DATA 0xeb 0xfd",

		"DATA 0x0a 0x0a",
		"DATA 0x68 0x65 0x6c 0x6c 0x6f 0x2c 0x20 0x77 0x6f 0x72 0x6c 0x64",
		"DATA 0x0a",
		"DATA 0x00",
		"FILL 378",

		"DATA 0x55 0xaa",
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

func TestDay02Suite(t *testing.T) {
	suite.Run(t, new(Day02Suite))
}

func (s *Day02Suite) SetupSuite() {
	setUpColog(true)
}

func (s *Day02Suite) TearDownSuite() {
}

func (s *Day02Suite) SetupTest() {
}

func (s *Day02Suite) TearDownTest() {
}
