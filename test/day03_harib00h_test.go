package test

import (
	"log"
	"os"

	"github.com/HobbyOSs/gosk/internal/frontend"
	"github.com/HobbyOSs/gosk/internal/gen"
	"github.com/google/go-cmp/cmp"
)

func (s *Day03Suite) TestHarib00h() {
	s.T().Skip("未実装の命令があるためスキップ")
	code := `; haribote-os
; TAB=4

; BOOT_INFO関係
CYLS	EQU		0x0ff0			; ブートセクタが設定する
LEDS	EQU		0x0ff1
VMODE	EQU		0x0ff2			; 色数に関する情報。何ビットカラーか？
SCRNX	EQU		0x0ff4			; 解像度のX
SCRNY	EQU		0x0ff6			; 解像度のY
VRAM	EQU		0x0ff8			; グラフィックバッファの開始番地

		ORG		0xc200			; このプログラムがどこに読み込まれるのか

		MOV		AL,0x13			; VGAグラフィックス、320x200x8bitカラー
		MOV		AH,0x00
		INT		0x10
		MOV		BYTE [VMODE],8	; 画面モードをメモする
		MOV		WORD [SCRNX],320
		MOV		WORD [SCRNY],200
		MOV		DWORD [VRAM],0x000a0000

; キーボードのLED状態をBIOSに教えてもらう

		MOV		AH,0x02
		INT		0x16 			; keyboard BIOS
		MOV		[LEDS],AL

fin:
		HLT
		JMP		fin
`

	temp, err := os.CreateTemp("", "harib00h.img")
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
		"DATA \"HARIBOTE\"",
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
		"DATA \"HARIBOTEOS \"",
		"DATA \"FAT12   \"",
		"FILL 18",
		"DATA 0xb8 0x00 0x00",
		"DATA 0x8e 0xd0",
		"DATA 0xbc 0x00 0x7c",
		"DATA 0x8e 0xd8",
		"DATA 0xb0 0x13",
		"DATA 0xb4 0x00",
		"DATA 0xcd 0x10",
		"DATA 0xc6 0x06 0xf2 0x0f 0x08",
		"DATA 0xc7 0x06 0xf4 0x0f 0x40 0x01",
		"DATA 0xc7 0x06 0xf6 0x0f 0xc8 0x00",
		"DATA 0x66 0xc7 0x06 0xf8 0x0f 0x00 0x00 0x0a 0x00",
		"DATA 0xb4 0x02",
		"DATA 0xcd 0x16",
		"DATA 0xa2 0xf1 0x0f",
		"DATA 0xf4",
		"DATA 0xeb 0xfd",
	})

	if diff := cmp.Diff(expected, actual); diff != "" {
		log.Printf("error: result mismatch:\n%s", DumpDiff(expected, actual, false))
		s.T().Fail()
	}

	if len(expected) != len(actual) {
		s.T().Fatalf("expected length %d, actual length %d", len(expected), len(actual))
	}
}
