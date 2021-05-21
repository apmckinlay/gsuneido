// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package dbms

import (
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestTimestamp(t *testing.T) {
	StartTimestamps()
	prev := Timestamp()
	for i := 0; i < 1100; i++ {
		ts := Timestamp()
		assert.T(t).That(ts.Compare(prev) > 0)
		prev = ts
	}
	if !testing.Short() {
		prev = runtime.Now()
		time.Sleep(1100 * time.Millisecond)
		assert.T(t).That(Timestamp().Compare(prev) > 0)
	}
}
