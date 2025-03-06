/*
// hexdump: a hexdumper utility written in Golang
//
// Copyright 2015 Jason E. Aten <j.e.aten -a-t- g-m-a-i-l dot c-o-m>
// License: MIT
*/
package test

import (
	"fmt"
)

func Dump(by []byte) string {
	hexdump := ""

	n := len(by)
	rowcount := 0
	stop := (n / 16) * 16
	k := 0
	for i := 0; i <= stop; i += 16 {
		k++
		hexdump += fmt.Sprintf("%06x  ", i*16) // offset
		if i+16 < n {
			rowcount = 16
			for j := 0; j < 16; j++ {
				hexdump += fmt.Sprintf("%02x ", by[i+j])
			}
		} else {
			rowcount = min(k*16, n) % 16
			for j := 0; j < rowcount; j++ {
				hexdump += fmt.Sprintf("%02x ", by[i+j])
			}
			for j := rowcount; j < 16; j++ {
				hexdump += fmt.Sprintf("   ")
			}
		}

		hexdump += fmt.Sprintf("'%s'\n", viewString(by[i:(i+rowcount)]))
	}

	return hexdump
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func viewString(b []byte) string {
	r := []rune(string(b))
	for i := range r {
		if r[i] < 32 || r[i] > 126 {
			r[i] = '.'
		}
	}
	return string(r)
}
