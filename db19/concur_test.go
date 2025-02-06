// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// TestConcur tests that persist doesn't write anything if no activity
func TestConcur(t *testing.T) {
	store := stor.HeapStor(16 * 1024)
	db := CreateDb(store)
	before := store.Size()
	persistInterval := 20 * time.Millisecond
	StartConcur(db, persistInterval)
	time.Sleep(40 * time.Millisecond)
	assert.T(t).This(store.Size()).Is(before)
	db.UpdateState(func(state *DbState) {
		meta := *state.Meta
		state.Meta = &meta
	})
	time.Sleep(40 * time.Millisecond)
	assert.T(t).That(store.Size() > before)
}
