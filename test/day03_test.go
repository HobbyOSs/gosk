package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type Day03Suite struct {
	suite.Suite
}

func TestDay03Suite(t *testing.T) {
	suite.Run(t, new(Day03Suite))
}

func (s *Day03Suite) SetupSuite() {}

func (s *Day03Suite) TearDownSuite() {}

func (s *Day03Suite) SetupTest() {}

func (s *Day03Suite) TearDownTest() {}
