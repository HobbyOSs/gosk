package pass1

import (
	"testing"

	"github.com/HobbyOSs/gosk/internal/ast" // Import ast package
	"github.com/HobbyOSs/gosk/internal/codegen"
	"github.com/HobbyOSs/gosk/internal/gen"
	ocode_client "github.com/HobbyOSs/gosk/internal/ocode_client"
	"github.com/HobbyOSs/gosk/pkg/asmdb"
	"github.com/HobbyOSs/gosk/pkg/cpu" // Keep ast for program argument
	"github.com/comail/colog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Pass1EvalSuite struct {
	suite.Suite
}

func TestPass1EvalSuite(t *testing.T) {
	suite.Run(t, new(Pass1EvalSuite))
}

func (s *Pass1EvalSuite) SetupSuite() {
	setUpColog(colog.LDebug)
}

type EvalTestParam struct {
	name        string      // 識別しやすいように name フィールドを追加
	bitMode     cpu.BitMode // cpu.BitMode を cpu.BitMode に変更
	text        string
	expectedLOC int32
}

func (s *Pass1EvalSuite) TestEvalProgramLOC() {
	tests := []EvalTestParam{
		// パラメータ化テスト；
		// 引数:
		// * 16bit/32bit モード
		// * アセンブラ文の一部
		// * 期待される機械語サイズ
		{
			name:        "ADD_mem_reg",
			bitMode:     cpu.MODE_16BIT,
			text:        "ADD [BX], AX",
			expectedLOC: 2,
		},
		{
			name:        "INT_imm",
			bitMode:     cpu.MODE_16BIT,
			text:        "INT 0x10",
			expectedLOC: 2,
		},
		{
			name:        "MOV_reg_mem",
			bitMode:     cpu.MODE_16BIT,
			text:        "MOV AL, [SI]",
			expectedLOC: 2,
		},
		{
			name:        "MOV_reg_imm",
			bitMode:     cpu.MODE_16BIT,
			text:        "MOV AX, 0",
			expectedLOC: 3,
		},
		{
			name:        "MOV_mem_imm8",
			bitMode:     cpu.MODE_16BIT,
			text:        "MOV BYTE [ 0x0ff2 ], 8",
			expectedLOC: 5,
		},
		{
			name:        "MOV_mem_imm16",
			bitMode:     cpu.MODE_16BIT,
			text:        "MOV WORD [ 0x0ff4 ], 320",
			expectedLOC: 6,
		},
		{
			name:        "MOV_mem_imm32",
			bitMode:     cpu.MODE_16BIT,
			text:        "MOV DWORD [ 0x0ff8 ], 0x000a0000",
			expectedLOC: 9,
		},
		{
			name:        "MOV_mem_reg8",
			bitMode:     cpu.MODE_16BIT,
			text:        "MOV [0x0ff0],CH",
			expectedLOC: 4,
		},
		{
			name:        "IMUL_reg_imm16",
			bitMode:     cpu.MODE_16BIT,
			text:        "IMUL ECX, 4608",
			expectedLOC: 7, // マスターデータは 7 バイト (66 69 c9 00 12 00 00) を示します
		},
		{
			name:        "MOV_mem_reg8_byte",
			bitMode:     cpu.MODE_16BIT,
			text:        "MOV BYTE [ 0x0ff0 ], CH",
			expectedLOC: 4,
		},
		{
			name:        "MOV_reg_seg",
			bitMode:     cpu.MODE_32BIT,
			text:        "MOV AX, SS",
			expectedLOC: 3,
		},
		{
			name:        "MOV_mem_reg8_addr",
			bitMode:     cpu.MODE_16BIT,
			text:        "MOV [ 0x0ff1 ], AL",
			expectedLOC: 3,
		},
		{
			name:        "MOV_reg_mem_disp_16bit",
			bitMode:     cpu.MODE_16BIT,
			text:        "MOV ECX, [EBX+16]",
			expectedLOC: 5, // マスターデータは 5 バイト (67 66 8b 4b 10) を示します
		},
		// --- EQU テストケース ---
		{
			name:    "EQU_simple_mov",
			bitMode: cpu.MODE_16BIT,
			text: `
				MY_CONST EQU 1234
				MOV AX, MY_CONST
			`,
			expectedLOC: 3, // MOV AX, imm16
		},
		{
			name:    "EQU_addr_mov",
			bitMode: cpu.MODE_16BIT,
			text: `
				ADDR EQU 0x100
				MOV BX, [ADDR]
			`,
			expectedLOC: 4, // MOV BX, [imm16]
		},
		{
			name:    "EQU_offset_mov",
			bitMode: cpu.MODE_16BIT,
			text: `
				OFFSET EQU 8
				MOV AL, [BP+OFFSET]
			`,
			expectedLOC: 3, // MOV AL, [BP+imm8]
		},
		{
			name:    "EQU_calc_add",
			bitMode: cpu.MODE_16BIT,
			text: `
				VAL1 EQU 10
				VAL2 EQU VAL1 * 2
				ADD CX, VAL2
			`,
			expectedLOC: 3, // 0x83 /0 ib | ADD r/m16, imm8 (VAL2 = 20)
		},
	}

	for _, tt := range tests {
		// 提供されている場合は tt.name を使用し、それ以外の場合は tt.text にフォールバックします
		testName := tt.name
		if testName == "" {
			testName = tt.text
		}
		s.T().Run(testName, func(t *testing.T) {
			ctx := &codegen.CodeGenContext{BitMode: tt.bitMode}
			client, err := ocode_client.NewCodegenClient(ctx)
			if !assert.NoError(t, err) {
				t.FailNow()
			}

			pass1 := &Pass1{
				LOC:      0,
				BitMode:  tt.bitMode,
				SymTable: make(map[string]int32), // SymTable の初期化を追加
				Client:   client,
				AsmDB:    asmdb.NewInstructionDB(),
				MacroMap: make(map[string]ast.Exp), // MacroMap の初期化を追加
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
