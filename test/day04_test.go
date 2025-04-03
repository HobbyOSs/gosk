package test

import (
	"testing"

	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/HobbyOSs/gosk/internal/codegen"
	"github.com/HobbyOSs/gosk/internal/gen"
	ocode_client "github.com/HobbyOSs/gosk/internal/ocode_client"
	"github.com/HobbyOSs/gosk/internal/pass1"
	"github.com/HobbyOSs/gosk/internal/pass2"
	"log" // Add log import

	"github.com/HobbyOSs/gosk/pkg/asmdb"
	"github.com/HobbyOSs/gosk/pkg/cpu"
	"github.com/google/go-cmp/cmp" // Add cmp import
	"github.com/stretchr/testify/suite"
)

type Day04Suite struct {
	suite.Suite
}

func TestDay04Suite(t *testing.T) {
	suite.Run(t, new(Day04Suite))
}

func (s *Day04Suite) SetupSuite() {
	setUpColog(true)
	UseAnsiColorForDiff = false // Set the package-level variable
}

func (s *Day04Suite) TearDownSuite() {}

func (s *Day04Suite) SetupTest() {}

func (s *Day04Suite) TearDownTest() {}

func (s *Day04Suite) TestHarib01a() {
	const naskStatements = `
; naskfunc
; TAB=4

[FORMAT "WCOFF"]				; オブジェクトファイルを作るモード
[INSTRSET "i486p"]				; 486の命令まで使いたいという記述
[BITS 32]						; 32ビットモード用の機械語を作らせる
[FILE "naskfunc.nas"]			; ソースファイル名情報

		GLOBAL	_io_hlt,_write_mem8

[SECTION .text]

_io_hlt:	; void io_hlt(void);
		HLT
		RET

_write_mem8:	; void write_mem8(int addr, int data);
		MOV		ECX,[ESP+4]		; [ESP+4]にaddrが入っているのでそれをECXに読み込む
		MOV		AL,[ESP+8]		; [ESP+8]にdataが入っているのでそれをALに読み込む
		MOV		[ECX],AL
		RET
`
	// Parse
	pt, err := gen.Parse("naskfunc.nas", []byte(naskStatements))
	s.Require().NoError(err)
	prog, ok := pt.(ast.Prog)
	s.Require().True(ok, "Parsed result is not ast.Prog")

	// Eval (Pass1 & Pass2 simulation)
	ctx := &codegen.CodeGenContext{
		BitMode:          cpu.MODE_32BIT, // Day04 is 32bit
		SymTable:         make(map[string]int32),
		GlobalSymbolList: []string{},
		MachineCode:      []byte{},
	}
	client, _ := ocode_client.NewCodegenClient(ctx)

	p1 := &pass1.Pass1{
		LOC:              0,
		BitMode:          ctx.BitMode,
		SymTable:         ctx.SymTable,
		GlobalSymbolList: ctx.GlobalSymbolList,
		ExternSymbolList: []string{},
		Client:           client,
		AsmDB:            asmdb.NewInstructionDB(),
		MacroMap:         make(map[string]ast.Exp),
	}
	p1.Eval(prog)
	ctx.SourceFileName = p1.SourceFileName // Copy SourceFileName

	p2 := &pass2.Pass2{
		BitMode:          p1.BitMode,
		OutputFormat:     p1.OutputFormat,
		SourceFileName:   p1.SourceFileName,
		CurrentSection:   p1.CurrentSection,
		SymTable:         p1.SymTable,
		GlobalSymbolList: p1.GlobalSymbolList,
		ExternSymbolList: p1.ExternSymbolList,
		Client:           p1.Client,
		DollarPos:        p1.DollarPosition,
	}
	err = p2.Eval(prog)
	s.Require().NoError(err)

	binoutContainer := ctx.MachineCode // Get binary from CodeGenContext

	expected := []byte{
		0x4C, 0x01, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x99, 0x00, 0x00, 0x00, 0x0A, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x2E, 0x74, 0x65, 0x78, 0x74, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x0D, 0x00, 0x00, 0x00, 0x8C, 0x00, 0x00, 0x00, 0x99, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x10, 0x60, 0x2E, 0x64, 0x61, 0x74,
		0x61, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x40, 0x00, 0x10, 0xC0, 0x2E, 0x62, 0x73, 0x73, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x10, 0xC0, 0xF4, 0xC3, 0x8B, 0x4C,
		0x24, 0x04, 0x8A, 0x44, 0x24, 0x08, 0x88, 0x01, 0xC3, 0x2E, 0x66, 0x69, 0x6C, 0x65, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0xFE, 0xFF, 0x00, 0x00, 0x67, 0x01, 0x6E, 0x61, 0x73, 0x6B, 0x66,
		0x75, 0x6E, 0x63, 0x2E, 0x6E, 0x61, 0x73, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2E, 0x74, 0x65,
		0x78, 0x74, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x03, 0x01, 0x0D,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x2E, 0x64, 0x61, 0x74, 0x61, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00,
		0x00, 0x03, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x2E, 0x62, 0x73, 0x73, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x03, 0x00, 0x00, 0x00, 0x03, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x5F, 0x69, 0x6F, 0x5F, 0x68, 0x6C, 0x74,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04,
		0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x10, 0x00, 0x00,
		0x00, 0x5F, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5F, 0x6D, 0x65, 0x6D, 0x38, 0x00,
	}

	if diff := cmp.Diff(expected, binoutContainer); diff != "" {
		log.Printf("error: result mismatch:\n%s", DumpDiff(expected, binoutContainer, false))
		s.T().Fail()
	}
	if len(expected) != len(binoutContainer) {
		s.T().Fatalf("expected length %d, actual length %d", len(expected), len(binoutContainer))
	}
}

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

	// Eval (Pass1 & Pass2 simulation)
	ctx := &codegen.CodeGenContext{
		BitMode:          cpu.MODE_32BIT, // Day04 is 32bit
		SymTable:         make(map[string]int32),
		GlobalSymbolList: []string{},
		MachineCode:      []byte{},
	}
	client, _ := ocode_client.NewCodegenClient(ctx)

	p1 := &pass1.Pass1{
		LOC:              0,
		BitMode:          ctx.BitMode,
		SymTable:         ctx.SymTable,
		GlobalSymbolList: ctx.GlobalSymbolList,
		ExternSymbolList: []string{},
		Client:           client,
		AsmDB:            asmdb.NewInstructionDB(),
		MacroMap:         make(map[string]ast.Exp),
	}
	p1.Eval(prog)
	ctx.SourceFileName = p1.SourceFileName // Copy SourceFileName

	p2 := &pass2.Pass2{
		BitMode:          p1.BitMode,
		OutputFormat:     p1.OutputFormat,
		SourceFileName:   p1.SourceFileName,
		CurrentSection:   p1.CurrentSection,
		SymTable:         p1.SymTable,
		GlobalSymbolList: p1.GlobalSymbolList,
		ExternSymbolList: p1.ExternSymbolList,
		Client:           p1.Client,
		DollarPos:        p1.DollarPosition,
	}
	err = p2.Eval(prog)
	s.Require().NoError(err)

	binoutContainer := ctx.MachineCode // Get binary from CodeGenContext

	expected := []byte{
		0x4C, 0x01, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0xDB, 0x00, 0x00, 0x00, 0x14, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x2E, 0x74, 0x65, 0x78, 0x74, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x4F, 0x00, 0x00, 0x00, 0x8C, 0x00, 0x00, 0x00, 0xDB, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x20, 0x00, 0x10, 0x60, 0x2E, 0x64, 0x61, 0x74,
		0x61, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x40, 0x00, 0x10, 0xC0, 0x2E, 0x62, 0x73, 0x73, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0x00, 0x10, 0xC0, 0xF4, 0xC3, 0xFA, 0xC3,
		0xFB, 0xC3, 0xFB, 0xF4, 0xC3, 0x8B, 0x54, 0x24, 0x04, 0xB8, 0x00, 0x00, 0x00, 0x00, 0xEC, 0xC3,
		0x8B, 0x54, 0x24, 0x04, 0xB8, 0x00, 0x00, 0x00, 0x00, 0x66, 0xED, 0xC3, 0x8B, 0x54, 0x24, 0x04,
		0xED, 0xC3, 0x8B, 0x54, 0x24, 0x04, 0x8A, 0x44, 0x24, 0x08, 0xEE, 0xC3, 0x8B, 0x54, 0x24, 0x04,
		0x8B, 0x44, 0x24, 0x08, 0x66, 0xEF, 0xC3, 0x8B, 0x54, 0x24, 0x04, 0x8B, 0x44, 0x24, 0x08, 0xEF,
		0xC3, 0x9C, 0x58, 0xC3, 0x8B, 0x44, 0x24, 0x04, 0x50, 0x9D, 0xC3, 0x2E, 0x66, 0x69, 0x6C, 0x65,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFE, 0xFF, 0x00, 0x00, 0x67, 0x01, 0x6E, 0x61, 0x73,
		0x6B, 0x66, 0x75, 0x6E, 0x63, 0x2E, 0x6E, 0x61, 0x73, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2E,
		0x74, 0x65, 0x78, 0x74, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x03,
		0x01, 0x4F, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x2E, 0x64, 0x61, 0x74, 0x61, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00,
		0x00, 0x00, 0x03, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x2E, 0x62, 0x73, 0x73, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x03, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x5F, 0x69, 0x6F, 0x5F, 0x68, 0x6C,
		0x74, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x5F, 0x69, 0x6F, 0x5F,
		0x63, 0x6C, 0x69, 0x00, 0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x5F, 0x69,
		0x6F, 0x5F, 0x73, 0x74, 0x69, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x5F, 0x69, 0x6F, 0x5F, 0x69, 0x6E, 0x38, 0x00, 0x09, 0x00, 0x00, 0x00, 0x01, 0x00,
		0x00, 0x00, 0x02, 0x00, 0x5F, 0x69, 0x6F, 0x5F, 0x69, 0x6E, 0x31, 0x36, 0x14, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x5F, 0x69, 0x6F, 0x5F, 0x69, 0x6E, 0x33, 0x32, 0x20, 0x00,
		0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x5F, 0x69, 0x6F, 0x5F, 0x6F, 0x75, 0x74, 0x38,
		0x26, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0F, 0x00,
		0x00, 0x00, 0x30, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x19, 0x00, 0x00, 0x00, 0x3B, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x23, 0x00, 0x00, 0x00, 0x45, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x02, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x33, 0x00, 0x00, 0x00, 0x48, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x44, 0x00, 0x00, 0x00, 0x5F, 0x69, 0x6F, 0x5F, 0x73, 0x74, 0x69, 0x68, 0x6C, 0x74,
		0x00, 0x5F, 0x69, 0x6F, 0x5F, 0x6F, 0x75, 0x74, 0x31, 0x36, 0x00, 0x5F, 0x69, 0x6F, 0x5F, 0x6F,
		0x75, 0x74, 0x33, 0x32, 0x00, 0x5F, 0x69, 0x6F, 0x5F, 0x6C, 0x6F, 0x61, 0x64, 0x5F, 0x65, 0x66,
		0x6C, 0x61, 0x67, 0x73, 0x00, 0x5F, 0x69, 0x6F, 0x5F, 0x73, 0x74, 0x6F, 0x72, 0x65, 0x5F, 0x65,
		0x66, 0x6C, 0x61, 0x67, 0x73, 0x00,
	}

	if diff := cmp.Diff(expected, binoutContainer); diff != "" {
		log.Printf("error: result mismatch:\n%s", DumpDiff(expected, binoutContainer, false))
		s.T().Fail()
	}
	if len(expected) != len(binoutContainer) {
		s.T().Fatalf("expected length %d, actual length %d", len(expected), len(binoutContainer))
	}
}
