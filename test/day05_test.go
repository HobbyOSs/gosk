package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// Day05Suite is already defined in day05_harib02i_test.go
// We just need to run it here.

func TestDay05Suite(t *testing.T) {
	suite.Run(t, new(Day05Suite))
}

// SetupSuite, TearDownSuite, SetupTest, TearDownTest can be defined in
// day05_harib02i_test.go if needed, or here if they apply to all Day05 tests.
// For now, we keep this file minimal.
// If common setup/teardown is needed later, add it to Day05Suite in day05_harib02i_test.go
// or create a base suite if multiple Day05 test files emerge.
