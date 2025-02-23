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

func (s *Pass1EvalSuite) TestEvalProgramLOC() {
	tests := []struct {
		name       string
		text       string
		entryPoint string
		expected   int32
	}{
		{
			name:       "simple program",
			text:       "ORG 0x7c00 ; comment",
			entryPoint: "Program",
			expected:   10, // 期待されるLOCの値
		},
		{
			name:       "complex program",
			text:       "MOV [CS:DS],8 ; comment",
			entryPoint: "Program",
			expected:   20, // 期待されるLOCの値
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			pass1 := &Pass1{
				LOC: 0,
			}
			got, err := gen.Parse("", []byte(tt.text), gen.Entrypoint(tt.entryPoint))
			if !assert.NoError(t, err) {
				t.FailNow()
			}
			program, ok := got.(*ast.Program)
			if !ok {
				t.FailNow()
			}
			pass1.Eval(*program)
			assert.Equal(t, tt.expected, pass1.LOC, "LOC should match expected value")
		})
	}
}
