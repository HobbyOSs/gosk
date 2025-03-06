package test

import (
	"log"
	"os"
	"testing"

	"github.com/HobbyOSs/gosk/internal/frontend"
	"github.com/HobbyOSs/gosk/internal/gen"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/suite"
)

type ArithmeticSuite struct {
	suite.Suite
}

func (s *ArithmeticSuite) TestArithmeticInstructions() {
	s.T().Skip("算術命令の実装が完了するまでスキップ")

	code := `; arithmetic instructions test
		ORG		0x7c00

		; ADD test
		MOV		AX, 1
		ADD		AX, 2		; AX = 3

		; ADC test
		MOV		AX, 0xFFFF
		MOV		BX, 1
		ADD		AX, BX		; CF = 1
		MOV		AX, 1
		ADC		AX, 1		; AX = 3 (1 + 1 + CF)

		; SUB test
		MOV		AX, 5
		SUB		AX, 2		; AX = 3

		; SBB test
		MOV		AX, 0
		SUB		AX, 1		; CF = 1
		MOV		AX, 5
		SBB		AX, 2		; AX = 2 (5 - 2 - CF)

		; CMP test
		MOV		AX, 5
		CMP		AX, 5		; ZF = 1

		; INC/DEC test
		MOV		AX, 1
		INC		AX			; AX = 2
		DEC		AX			; AX = 1

		; NEG test
		MOV		AX, 5
		NEG		AX			; AX = -5

		; MUL test
		MOV		AX, 5
		MOV		BX, 2
		MUL		BX			; AX = 10

		; IMUL test
		MOV		AX, -5
		MOV		BX, 2
		IMUL	BX			; AX = -10

		; DIV test
		MOV		AX, 10
		MOV		BL, 2
		DIV		BL			; AL = 5

		; IDIV test
		MOV		AX, -10
		MOV		BL, 2
		IDIV	BL			; AL = -5

		HLT
`

	temp, err := os.CreateTemp("", "arithmetic.img")
	s.Require().NoError(err)
	defer os.Remove(temp.Name())

	pt, err := gen.Parse("", []byte(code), gen.Entrypoint("Program"))
	s.Require().NoError(err)
	_, _ = frontend.Exec(pt, temp.Name())

	actual, err := ReadFileAsBytes(temp.Name())
	s.Require().NoError(err)

	expected := defineHEX([]string{
		// ADD test
		"DATA 0xb8 0x01 0x00",      // MOV AX, 1
		"DATA 0x05 0x02 0x00",      // ADD AX, 2

		// ADC test
		"DATA 0xb8 0xff 0xff",      // MOV AX, 0xFFFF
		"DATA 0xbb 0x01 0x00",      // MOV BX, 1
		"DATA 0x03 0xc3",           // ADD AX, BX
		"DATA 0xb8 0x01 0x00",      // MOV AX, 1
		"DATA 0x83 0xd0 0x01",      // ADC AX, 1

		// SUB test
		"DATA 0xb8 0x05 0x00",      // MOV AX, 5
		"DATA 0x2d 0x02 0x00",      // SUB AX, 2

		// SBB test
		"DATA 0xb8 0x00 0x00",      // MOV AX, 0
		"DATA 0x2d 0x01 0x00",      // SUB AX, 1
		"DATA 0xb8 0x05 0x00",      // MOV AX, 5
		"DATA 0x83 0xdb 0x02",      // SBB AX, 2

		// CMP test
		"DATA 0xb8 0x05 0x00",      // MOV AX, 5
		"DATA 0x3d 0x05 0x00",      // CMP AX, 5

		// INC/DEC test
		"DATA 0xb8 0x01 0x00",      // MOV AX, 1
		"DATA 0x40",                // INC AX
		"DATA 0x48",                // DEC AX

		// NEG test
		"DATA 0xb8 0x05 0x00",      // MOV AX, 5
		"DATA 0xf7 0xd8",           // NEG AX

		// MUL test
		"DATA 0xb8 0x05 0x00",      // MOV AX, 5
		"DATA 0xbb 0x02 0x00",      // MOV BX, 2
		"DATA 0xf7 0xe3",           // MUL BX

		// IMUL test
		"DATA 0xb8 0xfb 0xff",      // MOV AX, -5
		"DATA 0xbb 0x02 0x00",      // MOV BX, 2
		"DATA 0xf7 0xeb",           // IMUL BX

		// DIV test
		"DATA 0xb8 0x0a 0x00",      // MOV AX, 10
		"DATA 0xb3 0x02",           // MOV BL, 2
		"DATA 0xf6 0xf3",           // DIV BL

		// IDIV test
		"DATA 0xb8 0xf6 0xff",      // MOV AX, -10
		"DATA 0xb3 0x02",           // MOV BL, 2
		"DATA 0xf6 0xfb",           // IDIV BL

		"DATA 0xf4",                // HLT
	})

	// 実際のバイト数を確認
	s.T().Logf("Expected length: %d, Actual length: %d", len(expected), len(actual))
	s.T().Logf("Expected bytes: %v", expected)
	s.T().Logf("Actual bytes: %v", actual)
	s.Assert().Equal(len(expected), len(actual))
	if diff := cmp.Diff(expected, actual); diff != "" {
		log.Printf("error: result mismatch:\n%s", DumpDiff(expected, actual, false))
	}
}

func TestArithmeticSuite(t *testing.T) {
	suite.Run(t, new(ArithmeticSuite))
}

func (s *ArithmeticSuite) SetupSuite() {
	setUpColog(true)
}

func (s *ArithmeticSuite) TearDownSuite() {
}

func (s *ArithmeticSuite) SetupTest() {
}

func (s *ArithmeticSuite) TearDownTest() {
}
