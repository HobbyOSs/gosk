package pass1

import (
	"testing"

	"github.com/HobbyOSs/gosk/internal/ast"
	"github.com/HobbyOSs/gosk/internal/gen"
	client "github.com/HobbyOSs/gosk/internal/ocode_client"
	"github.com/HobbyOSs/gosk/internal/token"
	"github.com/HobbyOSs/gosk/pkg/asmdb"
	"github.com/comail/colog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zeroflucs-given/generics/collections/stack"
)

type Pass1EvalSuite struct {
	suite.Suite
}

func TestPass1EvalSuite(t *testing.T) {
	//t.Skip()
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
		// {
		// 	bitMode:     ast.MODE_16BIT,
		// 	text:        "ADD [BX], AX",
		// 	expectedLOC: 2,
		// },
		// {
		// 	bitMode:     ast.MODE_16BIT,
		// 	text:        "INT 0x10",
		// 	expectedLOC: 2,
		// },
		// {
		// 	bitMode:     ast.MODE_16BIT,
		// 	text:        "CALL waitkbdout",
		// 	expectedLOC: 5,
		// },
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
			// 0xc6, 0x06, 0xf2, 0x0f, 0x08
			bitMode:     ast.MODE_16BIT,
			text:        "MOV BYTE [ 0x0ff2 ], 8",
			expectedLOC: 5,
		},
		{
			// 0xc7, 0x06, 0xf4, 0x0f, 0x40, 0x01
			bitMode:     ast.MODE_16BIT,
			text:        "MOV [ 0x0ff4 ], 320",
			expectedLOC: 6,
		},
		// {
		// 	// MOV DWORD [VRAM],0x000a0000  ; VRAM=0x0ff8
		// 	// 0x66, 0xc7, 0x06, 0xf8, 0x0f, 0x00, 0x00, 0x0a, 0x00
		// 	bitMode:     ast.MODE_16BIT,
		// 	text:        "MOV DWORD [ 0x0ff8 ], 0x000a0000",
		// 	expectedLOC: 9,
		// },
		// {
		// 	bitMode:     ast.MODE_16BIT,
		// 	text:        "MOV CL, 0x0ff0",
		// 	expectedLOC: 4,
		// },
		// {
		// 	bitMode:     ast.MODE_16BIT,
		// 	text:        "OR EAX, 0x00000001",
		// 	expectedLOC: 4,
		// },
		// {
		// 	bitMode:     ast.MODE_16BIT,
		// 	text:        "IMUL ECX, 4608",
		// 	expectedLOC: 7,
		// },
		// {
		// 	bitMode:     ast.MODE_16BIT,
		// 	text:        "MOV BYTE [ 0x0ff0 ], CH",
		// 	expectedLOC: 4,
		// },
		// {
		// 	bitMode:     ast.MODE_16BIT,
		// 	text:        "SUB ECX, 128",
		// 	expectedLOC: 7,
		// },
		// {
		// 	bitMode:     ast.MODE_16BIT,
		// 	text:        "MOV ECX, [EBX+16]",
		// 	expectedLOC: 5,
		// },
		// {
		// 	bitMode:     ast.MODE_32BIT,
		// 	text:        "MOV AX, SS",
		// 	expectedLOC: 3,
		// },
	}

	for _, tt := range tests {
		s.T().Run(tt.text, func(t *testing.T) {
			pass1 := &Pass1{
				LOC:     0,
				BitMode: tt.bitMode,
				Ctx:     stack.NewStack[*token.ParseToken](10),
				Client:  client.NewCodegenClient(),
				AsmDB:   asmdb.NewInstructionDB(),
			}
			parseTree, err := gen.Parse("", []byte(tt.text), gen.Entrypoint("Program"))
			if !assert.NoError(t, err) {
				t.FailNow()
			}
			program, ok := (parseTree).(ast.Prog)
			if !ok {
				t.FailNow()
			}
			pass1.Eval(program)
			assert.Equal(t, tt.expectedLOC, pass1.LOC, "LOC should match expected value")
		})
	}
}
