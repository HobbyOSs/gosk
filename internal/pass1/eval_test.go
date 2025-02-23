package pass1

import (
	"testing"

	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/HobbyOSs/gosk/internal/gen"
	"github.com/comail/colog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Pass1EvalSuite struct {
	suite.Suite
}

func TestPass1EvalSuite(t *testing.T) {
	t.Skip()
	suite.Run(t, new(Pass1EvalSuite))
}

func (s *Pass1EvalSuite) SetupSuite() {
	setUpColog(colog.LDebug)
}

type EvalTestParam struct {
	bitMode     ast.BitMode
	text        string
	expectedLOC int32
}

func (s *Pass1EvalSuite) TestEvalProgramLOC() {
	tests := []EvalTestParam{
		// パラメタライズテスト；
		// 引数:
		// * 16bit/32bitモード
		// * アセンブラ文の一部
		// * 期待される機械語サイズ
		{
			bitMode:     ast.MODE_16BIT,
			text:        "ADD [BX], AX",
			expectedLOC: 2,
		},
		{
			bitMode:     ast.MODE_16BIT,
			text:        "INT 0x10",
			expectedLOC: 2,
		},
		{
			bitMode:     ast.MODE_16BIT,
			text:        "CALL waitkbdout",
			expectedLOC: 5,
		},
		{
			bitMode:     ast.MODE_16BIT,
			text:        "MOV AL, [SI]",
			expectedLOC: 2,
		},
		{
			bitMode:     ast.MODE_16BIT,
			text:        "MOV AX, 0",
			expectedLOC: 3,
		},
		{
			bitMode:     ast.MODE_16BIT,
			text:        "MOV 0x0ff2, 8",
			expectedLOC: 5,
		},
		{
			bitMode:     ast.MODE_16BIT,
			text:        "MOV 0x0ff4, 320",
			expectedLOC: 6,
		},
		{
			bitMode:     ast.MODE_16BIT,
			text:        "MOV 0x0ff8, 0x000a0000",
			expectedLOC: 9,
		},
		{
			bitMode:     ast.MODE_16BIT,
			text:        "MOV CL, 0x0ff0",
			expectedLOC: 4,
		},
		{
			bitMode:     ast.MODE_16BIT,
			text:        "OR EAX, 0x00000001",
			expectedLOC: 4,
		},
		{
			bitMode:     ast.MODE_16BIT,
			text:        "IMUL ECX, 4608",
			expectedLOC: 7,
		},
		{
			bitMode:     ast.MODE_16BIT,
			text:        "MOV 0x0ff0, CH",
			expectedLOC: 4,
		},
		{
			bitMode:     ast.MODE_16BIT,
			text:        "SUB ECX, 128",
			expectedLOC: 7,
		},
		{
			bitMode:     ast.MODE_16BIT,
			text:        "MOV ECX, [EBX+16]",
			expectedLOC: 5,
		},
		{
			bitMode:     ast.MODE_32BIT,
			text:        "MOV AX, SS",
			expectedLOC: 3,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.text, func(t *testing.T) {
			pass1 := &Pass1{
				LOC:     0,
				BitMode: tt.bitMode,
			}
			got, err := gen.Parse("", []byte(tt.text), gen.Entrypoint("Program"))
			if !assert.NoError(t, err) {
				t.FailNow()
			}
			program, ok := got.(*ast.Program)
			if !ok {
				t.FailNow()
			}
			pass1.Eval(*program)
			assert.Equal(t, tt.expectedLOC, pass1.LOC, "LOC should match expected value")
		})
	}
}
