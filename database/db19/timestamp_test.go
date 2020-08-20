// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestTimestamp(t *testing.T) {
	prev := Timestamp()
	for i := 0; i < 1100; i++ {
		ts := Timestamp()
		assert.T(t).That(ts.After(prev))
		prev = ts
	}
	if !testing.Short() {
		prev = time.Now()
		time.Sleep(1100 * time.Millisecond)
		assert.T(t).That(Timestamp().After(prev))
	}
}
