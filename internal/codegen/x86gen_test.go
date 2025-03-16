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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &CodeGenContext{
				MachineCode:    make([]byte, 0),
				VS:             nil,
				BitMode:        ast.MODE_16BIT,
				DollarPosition: 0,
				SymTable:       map[string]int32{},
			}
			result := GenerateX86(tt.ocodes, ast.MODE_16BIT, ctx)
			assert.Equal(t, tt.expected, result, "Test %s failed", tt.name)
		})
	}
}
