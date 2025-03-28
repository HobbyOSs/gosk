package test

import (
	"log"
	"os"

	"github.com/HobbyOSs/gosk/internal/frontend"
	"github.com/HobbyOSs/gosk/internal/gen"
	"github.com/google/go-cmp/cmp"
)

func (s *Day03Suite) TestHarib00j() {
	s.T().Skip("未実装の命令があるためスキップ")
	code := `; naskfunc
; TAB=4

[FORMAT "WCOFF"]				; オブジェクトファイルを作るモード
[BITS 32]						; 32ビットモード用の機械語を作らせる


; オブジェクトファイルのための情報

[FILE "naskfunc.nas"]			; ソースファイル名情報

		GLOBAL	_io_hlt			; このプログラムに含まれる関数名


; 以下は実際の関数

[SECTION .text]		; オブジェクトファイルではこれを書いてからプログラムを書く

_io_hlt:	; void io_hlt(void);
    	HLT
    	RET
`

	temp, err := os.CreateTemp("", "harib00j.img")
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
		"DATA 0x4c 0x01",           // machine
		"DATA 0x03 0x00",           // numberOfSections
		"DATA 0x00 0x00 0x00 0x00", // timeDateStamp
		"DATA 0x8e 0x00 0x00 0x00", // pointerToSymbolTable (シンボルテーブルへのオフセット; 後で計算される)
		"DATA 0x09 0x00 0x00 0x00", // numberOfSymbols (シンボルの数)
		"DATA 0x00 0x00",           // sizeOfOptionalHeader
		"DATA 0x00 0x00",           // flags

		"DATA 0x2e 0x74 0x65 0x78 0x74 0x00 0x00 0x00", // .text
		"DATA 0x00 0x00 0x00 0x00",                     // virtualSize
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x02 0x00 0x00 0x00",
		"DATA 0x8c 0x00 0x00 0x00",
		"DATA 0x8e 0x00 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x20 0x00 0x10 0x60",

		"DATA 0x2e 0x64 0x61 0x74 0x61 0x00 0x00 0x00", // .data
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x40 0x00 0x10 0xc0",

		"DATA 0x2e 0x62 0x73 0x73 0x00 0x00 0x00 0x00", // .bss
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x80 0x00 0x10 0xc0",

		"DATA 0xf4",
		"DATA 0xc3",

		"DATA 0x2e 0x66 0x69 0x6c 0x65 0x00 0x00 0x00", // .file
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0xfe 0xff 0x00 0x00",
		"DATA 0x67 0x01",
		"DATA 0x6e 0x61 0x73 0x6b 0x66 0x75 0x6e 0x63 0x2e 0x6e 0x61 0x73 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",

		"DATA 0x2e 0x74 0x65 0x78 0x74 0x00 0x00 0x00", // .text
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x01 0x00 0x00 0x00",
		"DATA 0x03 0x01",
		"DATA 0x02 0x00",
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",

		"DATA 0x2e 0x64 0x61 0x74 0x61 0x00 0x00 0x00", // .data
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x02 0x00 0x00 0x00",
		"DATA 0x03 0x01",
		"DATA 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",

		"DATA 0x2e 0x62 0x73 0x73 0x00 0x00 0x00 0x00", // .bss
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x03 0x00 0x00 0x00",
		"DATA 0x03 0x01",
		"DATA 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x00 0x00 0x00 0x00",

		"DATA 0x5f 0x69 0x6f 0x5f 0x68 0x6c 0x74 0x00", // シンボル情報
		"DATA 0x00 0x00 0x00 0x00",
		"DATA 0x01 0x00 0x00 0x00",
		"DATA 0x02 0x00",
		"DATA 0x04 0x00",
		"DATA 0x00 0x00",
	})

	if diff := cmp.Diff(expected, actual); diff != "" {
		log.Printf("error: result mismatch:\n%s", DumpDiff(expected, actual, false))
		s.T().Fail()
	}

	if len(expected) != len(actual) {
		s.T().Fatalf("expected length %d, actual length %d", len(expected), len(actual))
	}
}
