// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

//go:build dbg

package dbg

import "log"

func Assert(f func() bool) {
	if !f() {
		log.Println("ERROR: ASSERT FAILED")
		PrintStack()
		panic("ASSERT FAILED")
	}
}
