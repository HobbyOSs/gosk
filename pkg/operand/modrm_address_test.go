package operand

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMemoryOperand(t *testing.T) {
	type testCase struct {
		name     string
		memStr   string
		bitMode  BitMode // Change ast.BitMode to BitMode
		wantMod  byte
		wantRM   byte
		wantDisp []byte
		wantErr  bool
	}
	tests := []testCase{
		{
			name:     "16bit no disp [bx+si]",
			memStr:   "[bx+si]",
			bitMode:  MODE_16BIT, // Remove ast. prefix
			wantMod:  0b00,
			wantRM:   0b000,
			wantDisp: nil,
		},
		{
			name:     "16bit disp8 [bx+si+0x7f]",
			memStr:   "[bx+si+0x7f]",
			bitMode:  MODE_16BIT, // Remove ast. prefix
			wantMod:  0b01,
			wantRM:   0b000,
			wantDisp: []byte{0x7f},
		},
		{
			name:    "16bit disp16 [bx+di+0x1234]",
			memStr:  "[bx+di+0x1234]",
			bitMode: MODE_16BIT, // Remove ast. prefix
			wantMod: 0b10,
			wantRM:  0b001,
			// 0x1234 => LE: 34 12
			wantDisp: []byte{0x34, 0x12},
		},
		{
			name:     "32bit base [eax]",
			memStr:   "[eax]",
			bitMode:  MODE_32BIT, // Remove ast. prefix
			wantMod:  0b00,
			wantRM:   0b000,
			wantDisp: nil,
		},
		{
			name:     "32bit disp8 [ebp+8]",
			memStr:   "[ebp+8]",
			bitMode:  MODE_32BIT, // Remove ast. prefix
			wantMod:  0b01,
			wantRM:   0b101,
			wantDisp: []byte{0x08},
		},
		{
			name:    "32bit SIB not supported => error",
			memStr:  "[eax+ecx]",
			bitMode: MODE_32BIT, // Remove ast. prefix
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mod, rm, disp, err := ParseMemoryOperand(tc.memStr, tc.bitMode)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.wantMod, mod, "mod mismatch")
			assert.Equal(t, tc.wantRM, rm, "r/m mismatch")
			assert.Equal(t, tc.wantDisp, disp, "disp mismatch")
		})
	}
}

func TestCalcModRM(t *testing.T) {
	type testCase struct {
		name      string
		rmOperand string
		regBits   byte
		bitMode   BitMode // Change ast.BitMode to BitMode
		want      []byte
		wantErr   bool
	}
	tests := []testCase{
		{
			name:      "16bit mem disp8",
			rmOperand: "[bx+si+0x10]",
			regBits:   0b010, // 例: reg=010 => DX
			bitMode:   MODE_16BIT, // Remove ast. prefix
			// mod=01, rm=000 => 01 010 000 => 01010000(0x50)？
			// ただしここはバイト合成順に注意
			// mod=01(01000000=0x40) + reg=010<<3(=0x10) => 0x50 + rm=000(=0x0)
			// => 0x50, disp8=0x10 => [0x50, 0x10]
			want: []byte{0x50, 0x10},
		},
		{
			name:      "16bit register operand bx => mod=11, rm=011",
			rmOperand: "bx",
			regBits:   0b011, // (BXとかBLあたり)
			bitMode:   MODE_16BIT, // Remove ast. prefix
			// mod=11(0xC0), reg=011<<3=0x18, rm=011 => combine => 1100 0000 + 0001 1000 + 0000 1011...
			// (ここは計算順に注意) => 0xDB
			want: []byte{0xDB},
		},
		{
			name:      "32bit [eax], reg=001(ECX)",
			rmOperand: "[eax]",
			regBits:   0b001,
			bitMode:   MODE_32BIT, // Remove ast. prefix
			// mod=00, rm=000 => 000 001 000(b)=0x08
			want: []byte{0x08},
		},
		{
			name:      "32bit [ebp+8], reg=010(EDX)",
			rmOperand: "[ebp+8]",
			regBits:   0b010,
			bitMode:   MODE_32BIT, // Remove ast. prefix
			// parse => mod=01, rm=101 => 01 010 101 => 0x55
			// + disp8=0x08 => [0x55, 0x08]
			want: []byte{0x55, 0x08},
		},
		{
			name:      "32bit register ecx => mod=11, rm=001, reg=111 => 110 111 001 => 0xF9",
			rmOperand: "ecx",
			regBits:   0b111,
			bitMode:   MODE_32BIT, // Remove ast. prefix
			want:      []byte{0xF9},
		},
		{
			name:      "SIB required => error",
			rmOperand: "[eax+ecx]",
			regBits:   0,
			bitMode:   MODE_32BIT, // Remove ast. prefix
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := CalcModRM(tc.rmOperand, tc.regBits, tc.bitMode)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
