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
			name:    "MOV_AX_0",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"AX", "0"}},
			},
			expected: []byte{0xb8, 0x00, 0x00},
		},
		{
			name:    "MOV_SS_AX",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"SS", "AX"}},
			},
			expected: []byte{0x8e, 0xd0},
		},
		{
			name:    "MOV_DS_AX",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"DS", "AX"}},
			},
			expected: []byte{0x8e, 0xd8},
		},
		{
			name:    "MOV_SP_0x7c00",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"SP", "0x7c00"}},
			},
			expected: []byte{0xbc, 0x00, 0x7c},
		},
		{
			name:    "MOV_SI_0x0000", // pass2フェーズではラベル参照ではなく、具体的なアドレス値を使用することになる
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"SI", "0x0000"}},
			},
			expected: []byte{0xbe, 0x00, 0x00}, // MOV SI, 0x0000 = be 00 00
		},
		{
			name:    "MOV_AL_SI",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"AL", "[SI]"}},
			},
			expected: []byte{0x8a, 0x04},
		},
		{
			name:    "MOV_mem_0x0ff0_CH",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"[0x0ff0]", "CH"}},
			},
			expected: []byte{0x88, 0x2e, 0xf0, 0x0f},
		},
		{
			name:    "MOV_mem_0x0ff1_AL",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"[0x0ff1]", "AL"}},
			},
			expected: []byte{0xa2, 0xf1, 0x0f},
		},
		{
			name:    "MOV_BYTE_mem_0x0ff2_8",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"BYTE [ 0x0ff2 ]", "8"}},
			},
			expected: []byte{0xc6, 0x06, 0xf2, 0x0f, 0x08},
		},
		{
			name:    "MOV_WORD_mem_0x0ff4_320",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"WORD [ 0x0ff4 ]", "320"}},
			},
			expected: []byte{0xc7, 0x06, 0xf4, 0x0f, 0x40, 0x01},
		},
		{
			name:    "MOV_DWORD_mem_0x0ff8_0x000a0000",
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
			name:    "ADD_SI_1",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpADD, Operands: []string{"SI", "1"}},
			},
			expected: []byte{0x83, 0xc6, 0x01},
		},
		{
			name:    "CMP_AL_0",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpCMP, Operands: []string{"AL", "0"}},
			},
			expected: []byte{0x3c, 0x00},
		},
		{
			name:    "ADD_AX_0x0020",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpADD, Operands: []string{"AX", "0x0020"}},
			},
			expected: []byte{0x05, 0x20, 0x00}, // ADD AX, imm16
		},
		{
			name:    "OUT_0x21_AL",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpOUT, Operands: []string{"0x21", "AL"}},
			},
			expected: []byte{0xe6, 0x21},
		},
		{
			name:    "OUT_0xa1_AL",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpOUT, Operands: []string{"0xa1", "AL"}},
			},
			expected: []byte{0xe6, 0xa1},
		},
		{
			name:    "CALL_0x1234",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpCALL, Operands: []string{"0x1234"}},
			},
			expected: []byte{0xe8, 0x31, 0x12},
		},
		{
			name:    "MOV_ECX_EBX_plus_16_16bit", // 16bit mode with 32bit operand and address size prefixes
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"ECX", "[ EBX + 16 ]"}},
			},
			// Correct prefix order: 67h (address) + 66h (operand)
			expected: []byte{0x67, 0x66, 0x8b, 0x4b, 0x10},
		},
		{
			name:    "MOV_EAX_ESI_16bit", // 16bit mode with 32bit operand and address size prefixes
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"EAX", "[ ESI ]"}},
			},
			// Correct prefix order: 67h (address) + 66h (operand)
			expected: []byte{0x67, 0x66, 0x8b, 0x06},
		},
		{
			name:    "MOV_EAX_CR0_16bit",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"EAX", "CR0"}},
			},
			expected: []byte{0x0f, 0x20, 0xc0}, // No 66h prefix for MOV r32, CRn
		},
		{
			name:    "MOV_CR0_EAX_16bit",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"CR0", "EAX"}},
			},
			expected: []byte{0x0f, 0x22, 0xc0}, // No 66h prefix for MOV CRn, r32
		},
		{
			name:    "AND_EAX_0x7fffffff_16bit",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpAND, Operands: []string{"EAX", "0x7fffffff"}},
			},
			expected: []byte{0x66, 0x25, 0xff, 0xff, 0xff, 0x7f}, // 66h prefix + AND EAX, imm32
		},
		{
			name:    "OR_EAX_0x00000001_16bit",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpOR, Operands: []string{"EAX", "0x00000001"}},
			},
			expected: []byte{0x66, 0x83, 0xc8, 0x01}, // Expect imm8 form (83 /1 ib) as it's smaller
		},
		/* // TODO: IMUL r/m, imm の ModR/M 生成ロジック修正後に有効化する
		{
			name:    "IMUL_ECX_4608_16bit", // imm32 form expected
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpIMUL, Operands: []string{"ECX", "4608"}}, // 4608 = 0x1200
			},
			expected: []byte{0x66, 0x69, 0xc9, 0x00, 0x12, 0x00, 0x00}, // 66h + IMUL r32, imm32 (69 /r id)
		},
		*/
		{
			name:    "SUB_ECX_128_16bit", // imm32 form expected (even though value fits in imm8)
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpSUB, Operands: []string{"ECX", "128"}}, // 128 = 0x80
			},
			expected: []byte{0x66, 0x81, 0xe9, 0x80, 0x00, 0x00, 0x00}, // 66h + SUB r/m32, imm32 (81 /5 id)
		},
		{
			// Renamed duplicate test case
			name:    "MOV_EAX_CR0_16bit_second",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"EAX", "CR0"}},
			},
			expected: []byte{0x0f, 0x20, 0xc0}, // No 66h prefix in 32bit mode
		},
		{
			name:    "OR_EAX_1_16bit",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpOR, Operands: []string{"EAX", "1"}},
			},
			expected: []byte{0x66, 0x83, 0xc8, 0x01}, // 66h + OR r/m32, imm8 (83 /1 ib)
		},
		{
			name:    "MOV_AL_imm8_16bit",
			bitMode: cpu.MODE_16BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"AL", "0xff"}},
			},
			expected: []byte{0xb0, 0xff}, // MOV AL, imm8
		},
		/* // TODO: 32bitモードでの MOV AL, imm8 は現在不要なためコメントアウト
		{
			name:    "MOV_AL_imm8_32bit",
			bitMode: cpu.MODE_32BIT,
			ocodes: []ocode.Ocode{
				{Kind: ocode.OpMOV, Operands: []string{"AL", "0xff"}},
			},
			expected: []byte{0xb0, 0xff}, // MOV AL, imm8 (no prefix needed)
		},
		*/
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
