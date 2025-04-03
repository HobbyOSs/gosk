package test

import (
	"log"
	"os"

	"github.com/HobbyOSs/gosk/internal/frontend"
	"github.com/HobbyOSs/gosk/internal/gen"
	"github.com/google/go-cmp/cmp"
)

func (s *Day03Suite) TestHarib00i() {
	//s.T().Skip("未実装命令があるのでskip")
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
		"DATA 0xb0 0x13",
		"DATA 0xb4 0x00",
		"DATA 0xcd 0x10",
		"DATA 0xc6 0x06 0xf2 0x0f 0x08",                     // BYTE [VMODE],8
		"DATA 0xc7 0x06 0xf4 0x0f 0x40 0x01",                // WORD [SCRNX],320
		"DATA 0xc7 0x06 0xf6 0x0f 0xc8 0x00",                // WORD [SCRNY],200
		"DATA 0x66 0xc7 0x06 0xf8 0x0f 0x00 0x00 0x0a 0x00", // DWORD [VRAM],0x000a0000

		"DATA 0xb4 0x02",
		"DATA 0xcd 0x16",
		"DATA 0xa2 0xf1 0x0f",

		"DATA 0xb0 0xff",
		"DATA 0xe6 0x21",
		"DATA 0x90",
		"DATA 0xe6 0xa1",
		"DATA 0xfa",
		"DATA 0xe8 0xb5 0x00", // CALL waitkbdout
		"DATA 0xb0 0xd1",
		"DATA 0xe6 0x64",
		"DATA 0xe8 0xae 0x00", // CALL waitkbdout
		"DATA 0xb0 0xdf",
		"DATA 0xe6 0x60",
		"DATA 0xe8 0xa7 0x00", // CALL waitkbdout

		"DATA 0x0f 0x01 0x16 0x2a 0xc3", // LGDT[GDTR0]
		"DATA 0x0f 0x20 0xc0",
		"DATA 0x66 0x25 0xff 0xff 0xff 0x7f",
		"DATA 0x66 0x83 0xc8 0x01",
		"DATA 0x0f 0x22 0xc0",
		"DATA 0xeb 0x00",

		// pipelineflush:
		"DATA 0xb8 0x08 0x00",
		"DATA 0x8e 0xd8",
		"DATA 0x8e 0xc0",
		"DATA 0x8e 0xe0",
		"DATA 0x8e 0xe8",
		"DATA 0x8e 0xd0",

		// bootpackの転送
		"DATA 0x66 0xbe 0x30 0xc3 0x00 0x00",
		"DATA 0x66 0xbf 0x00 0x00 0x28 0x00",
		"DATA 0x66 0xb9 0x00 0x00 0x02 0x00",
		"DATA 0xe8 0x75 0x00",

		// まずはブートセクタから
		"DATA 0x66 0xbe 0x00 0x7c 0x00 0x00",
		"DATA 0x66 0xbf 0x00 0x00 0x10 0x00",
		"DATA 0x66 0xb9 0x80 0x00 0x00 0x00",
		"DATA 0xe8 0x60 0x00",

		// 残り全部
		"DATA 0x66 0xbe 0x00 0x82 0x00 0x00",
		"DATA 0x66 0xbf 0x00 0x02 0x10 0x00",
		"DATA 0x66 0xb9 0x00 0x00 0x00 0x00",
		"DATA 0x8a 0x0e 0xf0 0x0f",
		"DATA 0x66 0x69 0xc9 0x00 0x12 0x00 0x00",
		"DATA 0x66 0x81 0xe9 0x80 0x00 0x00 0x00",
		"DATA 0xe8 0x39 0x00",

		// bootpackの起動
		"DATA 0x66 0xbb 0x00 0x00 0x28 0x00",
		"DATA 0x67 0x66 0x8b 0x4b 0x10",
		"DATA 0x66 0x83 0xc1 0x03",
		"DATA 0x66 0xc1 0xe9 0x02",
		"DATA 0x74 0x10",
		"DATA 0x67 0x66 0x8b 0x73 0x14",
		"DATA 0x66 0x01 0xde",
		"DATA 0x67 0x66 0x8b 0x7b 0x0c",
		"DATA 0xe8 0x14 0x00",
		// skip:
		"DATA 0x67 0x66 0x8b 0x63 0x0c",
		"DATA 0x66 0xea 0x1b 0x00 0x00 0x00 0x10 0x00",

		// waitkbdout:
		"DATA 0xe4 0x64",
		"DATA 0x24 0x02",
		"DATA 0x75 0xfa",
		"DATA 0xc3",

		// memcpy:
		"DATA 0x67 0x66 0x8b 0x06",
		"DATA 0x66 0x83 0xc6 0x04",
		"DATA 0x67 0x66 0x89 0x07",
		"DATA 0x66 0x83 0xc7 0x04",
		"DATA 0x66 0x83 0xe9 0x01",
		"DATA 0x75 0xea",
		"DATA 0xc3",
		"FILL 8", // resb8
		"DATA 0xff 0xff 0x00 0x00 0x00 0x92 0xcf 0x00",
		"DATA 0xff 0xff 0x00 0x00 0x28 0x9a 0x47 0x00",
		"DATA 0x00 0x00",
		"DATA 0x17 0x00",
		"DATA 0x10 0xc3",
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
