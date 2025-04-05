package test

import (
	"log"
	"os"

	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/HobbyOSs/gosk/internal/frontend"
	"github.com/HobbyOSs/gosk/internal/gen"
	"github.com/google/go-cmp/cmp"
)

func (s *Day04Suite) TestHarib01f() {
	const naskStatements = `
; naskfunc
; TAB=4

[FORMAT "WCOFF"]				; オブジェクトファイルを作るモード
[INSTRSET "i486p"]				; 486の命令まで使いたいという記述
[BITS 32]						; 32ビットモード用の機械語を作らせる
[FILE "naskfunc.nas"]			; ソースファイル名情報

		GLOBAL	_io_hlt, _io_cli, _io_sti, _io_stihlt
		GLOBAL	_io_in8,  _io_in16,  _io_in32
		GLOBAL	_io_out8, _io_out16, _io_out32
		GLOBAL	_io_load_eflags, _io_store_eflags

[SECTION .text]

_io_hlt:	; void io_hlt(void);
		HLT
		RET

_io_cli:	; void io_cli(void);
		CLI
		RET

_io_sti:	; void io_sti(void);
		STI
		RET

_io_stihlt:	; void io_stihlt(void);
		STI
		HLT
		RET

_io_in8:	; int io_in8(int port);
		MOV		EDX,[ESP+4]		; port
		MOV		EAX,0
		IN		AL,DX
		RET

_io_in16:	; int io_in16(int port);
		MOV		EDX,[ESP+4]		; port
		MOV		EAX,0
		IN		AX,DX
		RET

_io_in32:	; int io_in32(int port);
		MOV		EDX,[ESP+4]		; port
		IN		EAX,DX
		RET

_io_out8:	; void io_out8(int port, int data);
		MOV		EDX,[ESP+4]		; port
		MOV		AL,[ESP+8]		; data
		OUT		DX,AL
		RET

_io_out16:	; void io_out16(int port, int data);
		MOV		EDX,[ESP+4]		; port
		MOV		EAX,[ESP+8]		; data
		OUT		DX,AX
		RET

_io_out32:	; void io_out32(int port, int data);
		MOV		EDX,[ESP+4]		; port
		MOV		EAX,[ESP+8]		; data
		OUT		DX,EAX
		RET

_io_load_eflags:	; int io_load_eflags(void);
		PUSHFD		; PUSH EFLAGS という意味
		POP		EAX
		RET

_io_store_eflags:	; void io_store_eflags(int eflags);
		MOV		EAX,[ESP+4]
		PUSH	EAX
		POPFD		; POP EFLAGS という意味
		RET
`
	// Parse
	pt, err := gen.Parse("naskfunc.nas", []byte(naskStatements))
	s.Require().NoError(err)
	prog, ok := pt.(ast.Prog)
	s.Require().True(ok, "Parsed result is not ast.Prog")

	// Create temp file
	temp, err := os.CreateTemp("", "harib01f_*.img")
	s.Require().NoError(err)
	defer os.Remove(temp.Name()) // clean up

	// Execute frontend
	_, _ = frontend.Exec(prog, temp.Name()) // Ignore both return values as error is handled by os.Exit within Exec

	// Read actual result from temp file
	actual, err := ReadFileAsBytes(temp.Name()) // Use ReadFileAsBytes from test_helper.go
	s.Require().NoError(err)

	expected := []byte{
		0x4C, 0x01, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0xDB, 0x00, 0x00, 0x00, 0x14, 0x00, 0x00, 0x00, // NumberOfSymbols = 0x14 (20)
		0x00, 0x00, 0x00, 0x00, 0x2E, 0x74, 0x65, 0x78, 0x74, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x4F, 0x00, 0x00, 0x00, 0x8C, 0x00, 0x00, 0x00, 0xDB, 0x00, 0x00, 0x00, // PointerToRelocations = 0xDB (symbol table offset)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x10, 0x60, 0x2E, 0x64, 0x61, 0x74,
		0x61, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // PointerToRawData = 0
		0x40, 0x00, 0x10, 0xC0, 0x2E, 0x62, 0x73, 0x73, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x10, 0xC0, 0xF4, 0xC3, 0xFA, 0xC3, // .text section data starts here (offset 0x8c)
		0xFB, 0xC3, 0xFB, 0xF4, 0xC3, 0x8B, 0x54, 0x24, 0x04, 0xB8, 0x00, 0x00, 0x00, 0x00, 0xEC, 0xC3,
		0x8B, 0x54, 0x24, 0x04, 0xB8, 0x00, 0x00, 0x00, 0x00, 0x66, 0xED, 0xC3, 0x8B, 0x54, 0x24, 0x04,
		0xED, 0xC3, 0x8B, 0x54, 0x24, 0x04, 0x8A, 0x44, 0x24, 0x08, 0xEE, 0xC3, 0x8B, 0x54, 0x24, 0x04,
		0x8B, 0x44, 0x24, 0x08, 0x66, 0xEF, 0xC3, 0x8B, 0x54, 0x24, 0x04, 0x8B, 0x44, 0x24, 0x08, 0xEF,
		0xC3, 0x9C, 0x58, 0xC3, 0x8B, 0x44, 0x24, 0x04, 0x50, 0x9D, 0xC3, 0x2E, 0x66, 0x69, 0x6C, 0x65, // Symbol table starts here (offset 0xDB)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFE, 0xFF, 0x00, 0x00, 0x67, 0x01, 0x6E, 0x61, 0x73, // .file aux symbol (naskfunc.nas)
		0x6B, 0x66, 0x75, 0x6E, 0x63, 0x2E, 0x6E, 0x61, 0x73, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2E, // .text symbol
		0x74, 0x65, 0x78, 0x74, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x03, // .text aux symbol (Length=0x4f)
		0x01, 0x4F, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x2E, 0x64, 0x61, 0x74, 0x61, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, // .data symbol
		0x00, 0x00, 0x03, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // .data aux symbol (Length=0)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x2E, 0x62, 0x73, 0x73, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // .bss symbol
		0x00, 0x03, 0x00, 0x00, 0x00, 0x03, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // .bss aux symbol (Length=0)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x5F, 0x69, 0x6F, 0x5F, 0x68, 0x6C, 0x74, 0x00, 0x00, // _io_hlt (Value=0, Section=1)
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x5F, 0x69, 0x6F, 0x5F, 0x63, 0x6C, 0x69, // _io_cli (Value=2, Section=1)
		0x00, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x5F, 0x69, 0x6F, 0x5F, 0x73, // _io_sti (Value=4, Section=1)
		0x74, 0x69, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, // _io_stihlt (Value=6, Section=1) - String table offset
		0x00, 0x04, 0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x5F, // _io_in8 (Value=9, Section=1)
		0x69, 0x6F, 0x5F, 0x69, 0x6E, 0x38, 0x00, 0x09, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, // _io_in16 (Value=0x14, Section=1) - String table offset
		0x00, 0x00, 0x00, 0x00, 0x00, 0x0F, 0x00, 0x00, 0x00, 0x14, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, // _io_in32 (Value=0x20, Section=1) - String table offset
		0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x19, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x01, // _io_out8 (Value=0x26, Section=1) - String table offset
		0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x23, 0x00, 0x00, 0x00, 0x26, 0x00, 0x00, // _io_out16 (Value=0x30, Section=1) - String table offset
		0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2D, 0x00, 0x00, 0x00, 0x30, // _io_out32 (Value=0x3B, Section=1) - String table offset
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x37, 0x00, 0x00, // _io_load_eflags (Value=0x45, Section=1) - String table offset
		0x00, 0x3B, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x41, // _io_store_eflags (Value=0x48, Section=1) - String table offset
		0x00, 0x00, 0x00, 0x45, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x51, 0x00, 0x00, 0x00, 0x48, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x5F, // String table size (0x5F = 95 bytes)
		0x00, 0x00, 0x00, 0x5F, 0x69, 0x6F, 0x5F, 0x73, 0x74, 0x69, 0x68, 0x6C, 0x74, 0x00, 0x5F, 0x69, // String table content
		0x6F, 0x5F, 0x69, 0x6E, 0x31, 0x36, 0x00, 0x5F, 0x69, 0x6F, 0x5F, 0x69, 0x6E, 0x33, 0x32, 0x00,
		0x5F, 0x69, 0x6F, 0x5F, 0x6F, 0x75, 0x74, 0x38, 0x00, 0x5F, 0x69, 0x6F, 0x5F, 0x6F, 0x75, 0x74,
		0x31, 0x36, 0x00, 0x5F, 0x69, 0x6F, 0x5F, 0x6F, 0x75, 0x74, 0x33, 0x32, 0x00, 0x5F, 0x69, 0x6F,
		0x5F, 0x6C, 0x6F, 0x61, 0x64, 0x5F, 0x65, 0x66, 0x6C, 0x61, 0x67, 0x73, 0x00, 0x5F, 0x69, 0x6F,
		0x5F, 0x73, 0x74, 0x6F, 0x72, 0x65, 0x5F, 0x65, 0x66, 0x6C, 0x61, 0x67, 0x73, 0x00,
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		log.Printf("error: result mismatch:\n%s", DumpDiff(expected, actual, false))
		s.T().Fail()
	}
	if len(expected) != len(actual) {
		s.T().Fatalf("expected length %d, actual length %d", len(expected), len(actual))
	}
}
