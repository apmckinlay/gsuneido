// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

/*
Package verify provides a simple assertion facility

For example:

	verify.That(size >= 0)
*/
package verify

import (
	"log"
	"runtime/debug"
)

// That panics if its argument is false
func That(cond bool) {
	if !cond {
		debug.PrintStack()
		log.Panicln("verify failed")
	}
}
