// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

// confirm the behavior of recover
// i.e. Go call stack is as of panic
// but defer's have been done

func TestRecover(*testing.T) {
	a()
}

func a() {
	defer func() {
		if e := recover(); e != nil {
			assert.That(unwound)
			// dbg.PrintStack()
		}
	}()
	b()
}

var unwound = false

func b() {
	defer func() {
		unwound = true
	}()
	c()
}

func c() {
	panic("foo")
}
