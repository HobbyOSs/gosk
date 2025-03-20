package test

import (
	"log"
	"os"

	"github.com/HobbyOSs/gosk/internal/frontend"
	"github.com/HobbyOSs/gosk/internal/gen"
	"github.com/google/go-cmp/cmp"
)

func (s *Day03Suite) TestHarib00b() {
	code := `; haribote-ipl
; TAB=4

		ORG		0x7c00			; このプログラムがどこに読み込まれるのか

; 以下は標準的なFAT12フォーマットフロッピーディスクのための記述

		JMP		entry
		DB		0x90
		DB		"HARIBOTE"		; ブートセクタの名前を自由に書いてよい（8バイト）
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
		DB		"HARIBOTEOS "	; ディスクの名前（11バイト）
		DB		"FAT12   "		; フォーマットの名前（8バイト）
		RESB	18				; とりあえず18バイトあけておく

; プログラム本体

entry:
		MOV		AX,0			; レジスタ初期化
		MOV		SS,AX
		MOV		SP,0x7c00
		MOV		DS,AX

; ディスクを読む

		MOV		AX,0x0820
		MOV		ES,AX
		MOV		CH,0			; シリンダ0
		MOV		DH,0			; ヘッド0
		MOV		CL,2			; セクタ2

		MOV		SI,0			; 失敗回数を数えるレジスタ
retry:
		MOV		AH,0x02			; AH=0x02 : ディスク読み込み
		MOV		AL,1			; 1セクタ
		MOV		BX,0
		MOV		DL,0x00			; Aドライブ
		INT		0x13			; ディスクBIOS呼び出し
		JNC		fin				; エラーがおきなければfinへ
		ADD		SI,1			; SIに1を足す
		CMP		SI,5			; SIと5を比較
		JAE		error			; SI >= 5 だったらerrorへ
		MOV		AH,0x00
		MOV		DL,0x00			; Aドライブ
		INT		0x13			; ドライブのリセット
		JMP		retry

; 読み終わったけどとりあえずやることないので寝る

fin:
		HLT						; 何かあるまでCPUを停止させる
		JMP		fin				; 無限ループ

error:
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
msg:
		DB		0x0a, 0x0a		; 改行を2つ
		DB		"load error"
		DB		0x0a			; 改行
		DB		0

		RESB	0x7dfe-$		; 0x7dfeまでを0x00で埋める命令

		DB		0x55, 0xaa
`

	temp, err := os.CreateTemp("", "harib00b.img")
	if err != nil {
		s.T().Fatal(err)
	}
	defer os.Remove(temp.Name()) // clean up

	pt, err := gen.Parse("", []byte(code), gen.Entrypoint("Program"))
	if err != nil {
		s.T().Fatal(err)
	}
	_, _ = frontend.Exec(pt, temp.Name())

	actual, err := ReadFileAsBytes(temp.Name())
	if err != nil {
		s.T().Fatal(err)
	}

	expected := defineHEX([]string{
		"DATA 0xeb 0x4e",
		"DATA 0x90",
		"DATA 0x48 0x41 0x52 0x49 0x42 0x4f 0x54 0x45",
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
		"DATA 0x48 0x41 0x52 0x49 0x42 0x4f 0x54 0x45 0x4f 0x53 0x20",
		"DATA 0x46 0x41 0x54 0x31 0x32 0x20 0x20 0x20",
		"FILL 18",
		"DATA 0xb8 0x00 0x00",
		"DATA 0x8e 0xd0",
		"DATA 0xbc 0x00 0x7c",
		"DATA 0x8e 0xd8",
		"DATA 0xb8 0x20 0x08",
		"DATA 0x8e 0xc0",
		"DATA 0xb5 0x00",
		"DATA 0xb6 0x00",
		"DATA 0xb1 0x02",
		"DATA 0xbe 0x00 0x00",
		"DATA 0xb4 0x02",
		"DATA 0xb0 0x01",
		"DATA 0xbb 0x00 0x00",
		"DATA 0xb2 0x00",
		"DATA 0xcd 0x13",
		"DATA 0x73 0x10",
		"DATA 0x83 0xc6 0x01",
		"DATA 0x83 0xfe 0x05",
		"DATA 0x73 0x0b",
		"DATA 0xb4 0x00",
		"DATA 0xb2 0x00",
		"DATA 0xcd 0x13",
		"DATA 0xeb 0xe3",
		"DATA 0xf4",
		"DATA 0xeb 0xfd",
		"DATA 0xbe 0x9d 0x7c",
		"DATA 0x8a 0x04",
		"DATA 0x83 0xc6 0x01",
		"DATA 0x3c 0x00",
		"DATA 0x74 0xf1",
		"DATA 0xb4 0x0e",
		"DATA 0xbb 0x0f 0x00",
		"DATA 0xcd 0x10",
		"DATA 0xeb 0xee",
		"DATA 0x0a 0x0a",
		"DATA 0x6c 0x6f 0x61 0x64 0x20 0x65 0x72 0x72 0x6f 0x72",
		"DATA 0x0a",
		"DATA 0x00",
		"FILL 339",
		"DATA 0x55 0xaa",
	})

	if diff := cmp.Diff(expected, actual); diff != "" {
		log.Printf("error: result mismatch:\n%s", DumpDiff(expected, actual, false))
		s.T().Fail()
	}

	if len(expected) != len(actual) {
		s.T().Fatalf("expected length %d, actual length %d", len(expected), len(actual))
	}
}
