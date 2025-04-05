package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type Day04Suite struct {
	suite.Suite
}

func TestDay04Suite(t *testing.T) {
	suite.Run(t, new(Day04Suite))
}

func (s *Day04Suite) SetupSuite() {
	setUpColog(true)
	UseAnsiColorForDiff = false // Set the package-level variable
}

func (s *Day04Suite) TearDownSuite() {}

func (s *Day04Suite) SetupTest() {}

func (s *Day04Suite) TearDownTest() {}
