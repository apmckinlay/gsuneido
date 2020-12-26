// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// +build !race

package runtime

import (
	"testing"
)

func TestRace(*testing.T) {
	println("*** CONCURRENCY TESTS SHOULD BE RUN WITH -race ***")
}
