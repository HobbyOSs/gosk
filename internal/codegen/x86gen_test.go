package codegen

import (
	"testing"

	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/HobbyOSs/gosk/pkg/ocode"
	"github.com/stretchr/testify/assert"
)

func TestGenerateX86(t *testing.T) {
	tests := []struct {
		name     string
		ocodes   []ocode.Ocode
		expected []byte
	}{
		{
			name: "DB",
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpDB, Operands: []string{"2", "224"}},
			},
			expected: []byte{0x02, 0xe0},
		},
		{
			name: "DW",
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpDW, Operands: []string{"4660"}},
			},
			expected: []byte{0x34, 0x12},
		},
		{
			name: "DD",
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpDD, Operands: []string{"305419896"}},
			},
			expected: []byte{0x78, 0x56, 0x34, 0x12},
		},
		{
			name: "RESB",
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpRESB, Operands: []string{"3"}},
			},
			expected: []byte{0x00, 0x00, 0x00},
		},
		{
			name: "INT",
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpINT, Operands: []string{"0x10"}},
			},
			expected: []byte{0xCD, 0x10}, // INT 0x10 = CD 10
		},
		{
			name: "MOV AX, 0",
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"AX", "0"}},
			},
			expected: []byte{0xb8, 0x00, 0x00},
		},
		{
			name: "MOV SS, AX",
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"SS", "AX"}},
			},
			expected: []byte{0x8e, 0xd0},
		},
		{
			name: "MOV DS, AX",
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"DS", "AX"}},
			},
			expected: []byte{0x8e, 0xd8},
		},
		{
			name: "MOV SP, 0x7c00",
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"SP", "0x7c00"}},
			},
			expected: []byte{0xbc, 0x00, 0x7c},
		},
		{
			name: "MOV SI, a_label",
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"SI", "a_label"}},
			},
			expected: []byte{0xbe, 0x00, 0x00},
		},
		{
			name: "MOV AL, [ SI ]",
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"AL", "[SI]"}},
			},
			expected: []byte{0x8a, 0x04},
		},
		{
			name: "MOV [ 0x0ff0 ],CH",
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"[0x0ff0]", "CH"}},
			},
			expected: []byte{0x88, 0x2e, 0xf0, 0x0f},
		},
		{
			name: "MOV [ 0x0ff1 ],AL",
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"[0x0ff1]", "AL"}},
			},
			expected: []byte{0xa2, 0xf1, 0x0f},
		},
		{
			name: "MOV BYTE [ 0x0ff2 ], 8",
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"BYTE [ 0x0ff2 ]", "8"}},
			},
			expected: []byte{0xc6, 0x06, 0xf2, 0x0f, 0x08},
		},
		{
			name: "MOV WORD [ 0x0ff4 ], 320",
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"WORD [ 0x0ff4 ]", "320"}},
			},
			expected: []byte{0xc7, 0x06, 0xf4, 0x0f, 0x40, 0x01},
		},
		{
			name: "MOV DWORD [ 0x0ff8 ],0x000a0000",
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"DWORD [ 0x0ff8 ]", "0x000a0000"}},
			},
			expected: []byte{0x66, 0xc7, 0x06, 0xf8, 0x0f, 0x00, 0x00, 0x0a, 0x00},
		},
		{
			name: "ADD SI,1",
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpADD, Operands: []string{"SI", "1"}},
			},
			expected: []byte{0x83, 0xc6, 0x01},
		},
		{
			name: "CMP AL,0",
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpCMP, Operands: []string{"AL", "0"}},
			},
			expected: []byte{0x3c, 0x00},
		},
		{
			name: "OUT 0x21, AL",
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpOUT, Operands: []string{"0x21", "AL"}},
			},
			expected: []byte{0xe6, 0x21},
		},
		{
			name: "OUT 0xa1, AL",
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpOUT, Operands: []string{"0xa1", "AL"}},
			},
			expected: []byte{0xe6, 0xa1},
		},
		{
			name: "CALL 0x1234",
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpCALL, Operands: []string{"0x1234"}},
			},
			expected: []byte{0xe8, 0x2f, 0x12, 0x00, 0x00}, // CALL 0x1234 = e8 2f 12 00 00 (little endian)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &CodeGenContext{
				MachineCode:    make([]byte, 0),
				VS:             nil,
				BitMode:        ast.MODE_16BIT,
				DollarPosition: 0, // Assume DollarPosition is 0 for simplicity in this test
				SymTable:       map[string]int32{},
			}
			result := GenerateX86(tt.ocodes, ctx)
			assert.Equal(t, tt.expected, result, "Test %s failed", tt.name)
		})
	}
}
