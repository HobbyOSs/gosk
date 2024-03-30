package test

import "github.com/sergi/go-diff/diffmatchpatch"

func DumpDiff(expected []byte, actual []byte, checklines bool) string {
	limit := 1000
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(
		Dump(takeN(expected, limit)),
		Dump(takeN(actual, limit)),
		checklines,
	)
	return dmp.DiffPrettyText(diffs)
}

func takeN(slice []byte, limit int) []byte {
	if len(slice) > limit {
		return slice[:limit]
	}
	return slice
}
