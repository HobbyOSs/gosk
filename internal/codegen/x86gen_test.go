package codegen

import (
	"testing"

	"github.com/HobbyOSs/gosk/pkg/cpu"
	"github.com/HobbyOSs/gosk/pkg/ocode"
	"github.com/stretchr/testify/assert"
)

func TestGenerateX86(t *testing.T) {
	tests := []struct {
		name     string
		ocodes   []ocode.Ocode
		bitMode  cpu.BitMode // Add BitMode field
		expected []byte
	}{
		{
			name:    "DB",
			bitMode: cpu.MODE_16BIT, // Default to 16bit
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpDB, Operands: []string{"2", "224"}},
			},
			expected: []byte{0x02, 0xe0},
		},
		{
			name:    "DW",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpDW, Operands: []string{"4660"}},
			},
			expected: []byte{0x34, 0x12},
		},
		{
			name:    "DD",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpDD, Operands: []string{"305419896"}},
			},
			expected: []byte{0x78, 0x56, 0x34, 0x12},
		},
		{
			name:    "RESB",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpRESB, Operands: []string{"3"}},
			},
			expected: []byte{0x00, 0x00, 0x00},
		},
		{
			name:    "INT",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpINT, Operands: []string{"0x10"}},
			},
			expected: []byte{0xCD, 0x10}, // INT 0x10 = CD 10
		},
		{
			name:    "MOV AX, 0",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"AX", "0"}},
			},
			expected: []byte{0xb8, 0x00, 0x00},
		},
		{
			name:    "MOV SS, AX",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"SS", "AX"}},
			},
			expected: []byte{0x8e, 0xd0},
		},
		{
			name:    "MOV DS, AX",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"DS", "AX"}},
			},
			expected: []byte{0x8e, 0xd8},
		},
		{
			name:    "MOV SP, 0x7c00",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"SP", "0x7c00"}},
			},
			expected: []byte{0xbc, 0x00, 0x7c},
		},
		{
			name:    "MOV SI, 0x0000", // pass2フェーズではラベル参照ではなく、具体的なアドレス値を使用することになる
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"SI", "0x0000"}},
			},
			expected: []byte{0xbe, 0x00, 0x00}, // MOV SI, 0x0000 = be 00 00
		},
		{
			name:    "MOV AL, [ SI ]",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"AL", "[SI]"}},
			},
			expected: []byte{0x8a, 0x04},
		},
		{
			name:    "MOV [ 0x0ff0 ],CH",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"[0x0ff0]", "CH"}},
			},
			expected: []byte{0x88, 0x2e, 0xf0, 0x0f},
		},
		{
			name:    "MOV [ 0x0ff1 ],AL",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"[0x0ff1]", "AL"}},
			},
			expected: []byte{0xa2, 0xf1, 0x0f},
		},
		{
			name:    "MOV BYTE [ 0x0ff2 ], 8",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"BYTE [ 0x0ff2 ]", "8"}},
			},
			expected: []byte{0xc6, 0x06, 0xf2, 0x0f, 0x08},
		},
		{
			name:    "MOV WORD [ 0x0ff4 ], 320",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"WORD [ 0x0ff4 ]", "320"}},
			},
			expected: []byte{0xc7, 0x06, 0xf4, 0x0f, 0x40, 0x01},
		},
		{
			name:    "MOV DWORD [ 0x0ff8 ],0x000a0000",
			bitMode: cpu.MODE_32BIT, // Needs 32bit mode for DWORD
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"DWORD [ 0x0ff8 ]", "0x000a0000"}},
			},
			// Correct expected value for 32-bit mode: C7 /0 id -> MOV r/m32, imm32
			// ModR/M for [disp32] is mod=00, r/m=101 -> 05h
			// No 66h prefix needed as default operand size is 32-bit.
			expected: []byte{0xc7, 0x05, 0xf8, 0x0f, 0x00, 0x00, 0x00, 0x00, 0x0a, 0x00},
		},
		{
			name:    "ADD SI,1",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpADD, Operands: []string{"SI", "1"}},
			},
			expected: []byte{0x83, 0xc6, 0x01},
		},
		{
			name:    "CMP AL,0",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpCMP, Operands: []string{"AL", "0"}},
			},
			expected: []byte{0x3c, 0x00},
		},
		{
			name:    "ADD AX, 0x0020",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpADD, Operands: []string{"AX", "0x0020"}},
			},
			expected: []byte{0x05, 0x20, 0x00}, // ADD AX, imm16
		},
		{
			name:    "OUT 0x21, AL",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpOUT, Operands: []string{"0x21", "AL"}},
			},
			expected: []byte{0xe6, 0x21},
		},
		{
			name:    "OUT 0xa1, AL",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpOUT, Operands: []string{"0xa1", "AL"}},
			},
			expected: []byte{0xe6, 0xa1},
		},
		{
			name:    "CALL 0x1234",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpCALL, Operands: []string{"0x1234"}},
			},
			expected: []byte{0xe8, 0x31, 0x12},
		},
		{
			name:    "MOV ECX, [ EBX + 16 ] (16bit)", // 16bit mode with 32bit operand and address size prefixes
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"ECX", "[ EBX + 16 ]"}},
			},
			// Correct prefix order: 67h (address) + 66h (operand)
			expected: []byte{0x67, 0x66, 0x8b, 0x4b, 0x10},
		},
		{
			name:    "MOV EAX, [ ESI ] (16bit)", // 16bit mode with 32bit operand and address size prefixes
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"EAX", "[ ESI ]"}},
			},
			// Correct prefix order: 67h (address) + 66h (operand)
			expected: []byte{0x67, 0x66, 0x8b, 0x06},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &CodeGenContext{
				MachineCode:    make([]byte, 0),
				VS:             nil,
				BitMode:        tt.bitMode, // Use BitMode from test case
				DollarPosition: 0,
				SymTable:       map[string]int32{},
			}
			result := GenerateX86(tt.ocodes, ctx)
			assert.Equal(t, tt.expected, result, "Test %s failed", tt.name)
		})
	}
}
