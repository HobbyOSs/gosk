package test

import (
	"log"
	"os"

	"github.com/HobbyOSs/gosk/internal/frontend"
	"github.com/HobbyOSs/gosk/internal/gen"
	"github.com/google/go-cmp/cmp"
)

func (s *Day03Suite) TestHarib00i() {
	//s.T().Skip("未実装の命令があるためスキップ")
	code := `; haribote-os boot asm
; TAB=4

BOTPAK	EQU		0x00280000		; bootpackのロード先
DSKCAC	EQU		0x00100000		; ディスクキャッシュの場所
DSKCAC0	EQU		0x00008000		; ディスクキャッシュの場所（リアルモード）

; BOOT_INFO関係
CYLS	EQU		0x0ff0			; ブートセクタが設定する
LEDS	EQU		0x0ff1
VMODE	EQU		0x0ff2			; 色数に関する情報。何ビットカラーか？
SCRNX	EQU		0x0ff4			; 解像度のX
SCRNY	EQU		0x0ff6			; 解像度のY
VRAM	EQU		0x0ff8			; グラフィックバッファの開始番地

		ORG		0xc200			; このプログラムがどこに読み込まれるのか

; 画面モードを設定

		MOV		AL,0x13			; VGAグラフィックス、320x200x8bitカラー
		MOV		AH,0x00
		INT		0x10
		MOV		BYTE [VMODE],8	; 画面モードをメモする（C言語が参照する）
		MOV		WORD [SCRNX],320
		MOV		WORD [SCRNY],200
		MOV		DWORD [VRAM],0x000a0000

; キーボードのLED状態をBIOSに教えてもらう

		MOV		AH,0x02
		INT		0x16			; keyboard BIOS
		MOV		[LEDS],AL

; PICが一切の割り込みを受け付けないようにする
;	AT互換機の仕様では、PICの初期化をするなら、
;	こいつをCLI前にやっておかないと、たまにハングアップする
;	PICの初期化はあとでやる

		MOV		AL,0xff
		OUT		0x21,AL
		NOP						; OUT命令を連続させるとうまくいかない機種があるらしいので
		OUT		0xa1,AL

		CLI						; さらにCPUレベルでも割り込み禁止

; CPUから1MB以上のメモリにアクセスできるように、A20GATEを設定

		CALL	waitkbdout
		MOV		AL,0xd1
		OUT		0x64,AL
		CALL	waitkbdout
		MOV		AL,0xdf			; enable A20
		OUT		0x60,AL
		CALL	waitkbdout

; プロテクトモード移行

[INSTRSET "i486p"]				; 486の命令まで使いたいという記述

		LGDT	[GDTR0]			; 暫定GDTを設定
		MOV		EAX,CR0
		AND		EAX,0x7fffffff	; bit31を0にする（ページング禁止のため）
		OR		EAX,0x00000001	; bit0を1にする（プロテクトモード移行のため）
		MOV		CR0,EAX
		JMP		pipelineflush
pipelineflush:
		MOV		AX,1*8			;  読み書き可能セグメント32bit
		MOV		DS,AX
		MOV		ES,AX
		MOV		FS,AX
		MOV		GS,AX
		MOV		SS,AX

; bootpackの転送

		MOV		ESI,bootpack	; 転送元
		MOV		EDI,BOTPAK		; 転送先
		MOV		ECX,512*1024/4
		CALL	memcpy

; ついでにディスクデータも本来の位置へ転送

; まずはブートセクタから

		MOV		ESI,0x7c00		; 転送元
		MOV		EDI,DSKCAC		; 転送先
		MOV		ECX,512/4
		CALL	memcpy

; 残り全部

		MOV		ESI,DSKCAC0+512	; 転送元
		MOV		EDI,DSKCAC+512	; 転送先
		MOV		ECX,0
		MOV		CL,BYTE [CYLS]
		IMUL	ECX,512*18*2/4	; シリンダ数からバイト数/4に変換
		SUB		ECX,512/4		; IPLの分だけ差し引く
		CALL	memcpy

; asmheadでしなければいけないことは全部し終わったので、
;	あとはbootpackに任せる

; bootpackの起動

		MOV		EBX,BOTPAK
		MOV		ECX,[EBX+16]
		ADD		ECX,3			; ECX += 3;
		SHR		ECX,2			; ECX /= 4;
		JZ		skip			; 転送するべきものがない
		MOV		ESI,[EBX+20]	; 転送元
		ADD		ESI,EBX
		MOV		EDI,[EBX+12]	; 転送先
		CALL	memcpy
skip:
		MOV		ESP,[EBX+12]	; スタック初期値
		JMP		DWORD 2*8:0x0000001b

waitkbdout:
		IN		 AL,0x64
		AND		 AL,0x02
		JNZ		waitkbdout		; ANDの結果が0でなければwaitkbdoutへ
		RET

memcpy:
		MOV		EAX,[ESI]
		ADD		ESI,4
		MOV		[EDI],EAX
		ADD		EDI,4
		SUB		ECX,1
		JNZ		memcpy			; 引き算した結果が0でなければmemcpyへ
		RET
; memcpyはアドレスサイズプリフィクスを入れ忘れなければ、ストリング命令でも書ける

		ALIGNB	16
GDT0:
		RESB	8				; ヌルセレクタ
		DW		0xffff,0x0000,0x9200,0x00cf	; 読み書き可能セグメント32bit
		DW		0xffff,0x0000,0x9a28,0x0047	; 実行可能セグメント32bit（bootpack用）

		DW		0
GDTR0:
		DW		8*3-1
		DD		GDT0

		ALIGNB	16
bootpack:
`

	temp, err := os.CreateTemp("", "harib00i.img")
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
		"DATA 0xb8 0x00 0x06",
		"DATA 0xb7 0x07",
		"DATA 0xb9 0x00 0x00",
		"DATA 0xba 0x4f 0x18",
		"DATA 0xcd 0x10",
		"DATA 0xbe 0x9e 0x7c",
		"DATA 0xbb 0x00 0x00",
		"DATA 0x8a 0x04",
		"DATA 0x83 0xc6 0x01",
		"DATA 0x3c 0x00",
		"DATA 0x74 0x13",
		"DATA 0xb4 0x0e",
		"DATA 0x8a 0x9b 0xaa 0x7c",
		"DATA 0xcd 0x10",
		"DATA 0x83 0xc3 0x01",
		"DATA 0x83 0xfb 0x10",
		"DATA 0x75 0xee",
		"DATA 0xb8 0x12 0x10",
		"DATA 0xbb 0x00 0x00",
		"DATA 0xb9 0x10 0x00",
		"DATA 0xba 0xba 0x7c",
		"DATA 0xcd 0x10",
		"DATA 0xb8 0x00 0xb8",
		"DATA 0x8e 0xc0",
		"DATA 0xbb 0x00 0x00",
		"DATA 0xc6 0x07 0x00",
		"DATA 0x83 0xc3 0x01",
		"DATA 0x81 0xfb 0x00 0xfa",
		"DATA 0x75 0xf6",
		"DATA \"hello, world\"",
		"DATA 0x0a",
		"DATA 0x00",
		"DATA 0x01 0x02 0x03 0x04 0x05 0x06 0x07 0x08",
		"DATA 0x09 0x0a 0x0b 0x0c 0x0d 0x0e 0x0f",
		"DATA 0x00 0x00 0x00",
		"DATA 0xff 0x00 0x00",
		"DATA 0x00 0xff 0x00",
		"DATA 0xff 0xff 0x00",
		"DATA 0x00 0x00 0xff",
		"DATA 0xff 0x00 0xff",
		"DATA 0x00 0xff 0xff",
		"DATA 0xff 0xff 0xff",
		"DATA 0xc6 0xc6 0xc6",
		"DATA 0x84 0x00 0x00",
		"DATA 0x00 0x84 0x00",
		"DATA 0x84 0x84 0x00",
		"DATA 0x00 0x00 0x84",
		"DATA 0x84 0x00 0x84",
		"DATA 0x00 0x84 0x84",
		"DATA 0x84 0x84 0x84",
		"DATA 0xf4",
		"DATA 0xeb 0xfe",
	})

	if len(expected) != len(actual) {
		s.T().Fatalf("expected length %d, actual length %d", len(expected), len(actual))
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		log.Printf("error: result mismatch:\n%s", DumpDiff(expected, actual, false))
		s.T().Fail()
	}
}
