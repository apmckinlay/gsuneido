// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbg

import "fmt"

const maxLen = 35

func PrintString(s string) {
	elipsis := ""
	if len(s) > maxLen {
		s = s[:maxLen]
		elipsis = " ..."
	}
	fmt.Printf("%x%s\n", s, elipsis)
	for _, b := range []byte(s) {
		if b < 32 || 126 < b {
			b = ' '
		}
		fmt.Print(" ", string(b))
	}
	fmt.Println(elipsis)
}

func PrintSlice(s []byte) {
	elipsis := ""
	if len(s) > maxLen {
		s = s[:maxLen]
		elipsis = " ..."
	}
	fmt.Printf("%x%s\n", s, elipsis)
	for _, b := range []byte(s) {
		if b < 32 || 126 < b {
			b = ' '
		}
		fmt.Print(" ", string(b))
	}
	fmt.Println(elipsis)
}
