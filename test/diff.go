package test

import (
	"github.com/akedrou/textdiff"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// UseAnsiColorForDiff controls whether ANSI color codes are used in diff output.
// This should be set by test suites during setup.
var UseAnsiColorForDiff bool = false // Default to false

func DumpDiff(expected []byte, actual []byte, checklines bool) string { // Removed useANSIColor parameter
	limit := 1000
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(
		Dump(takeN(expected, limit)),
		Dump(takeN(actual, limit)),
		checklines,
	)
	if UseAnsiColorForDiff { // Use the package-level variable
		return dmp.DiffPrettyText(diffs)
	} else {
		expectedHex := Dump(takeN(expected, limit))
		actualHex := Dump(takeN(actual, limit))
		unifiedDiff := textdiff.Unified("expected", "actual", expectedHex, actualHex)
		return unifiedDiff
	}
}

func takeN(slice []byte, limit int) []byte {
	if len(slice) > limit {
		return slice[:limit]
	}
	return slice
}
