// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"math/rand"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestMeta(t *testing.T) {
	tbl := InfoHamt{}.Mutable()
	const n = 1000
	data := make([]string, n)
	randStr := str.UniqueRandom(4, 4)
	for i := 0; i < n; i++ {
		data[i] = randStr()
		tbl.Put(&Info{
			Table: data[i],
		})
	}
	tbl = tbl.Freeze()
	st := stor.HeapStor(32 * 1024)
	st.Alloc(1) // avoid offset 0
	meta := &Meta{info: tbl}

	for i := 0; i < 4; i++ {
		m := meta.Mutable()
		for i := 0; i < 5; i++ {
			m.GetRoInfo(data[rand.Intn(100)])
			m.GetRwInfo(data[rand.Intn(100)], 0)
		}
		// end of transaction, merge back to global
		meta = m.LayeredOnto(meta)
	}

	// persist state
	meta.Write(st, true)

	// test that nothing is written if no changes
	size := st.Size()
	meta.Write(st, false)
	assert.T(t).This(st.Size()).Is(size)
}

func TestMergeSize(t *testing.T) {
	assert := assert.T(t)
	test := func (clock, expected_npersists, expected_timespan int) {
		t.Helper()
		npersists, timespan := mergeSize(clock, false)
		assert.Msg("npersists").This(npersists).Is(expected_npersists)
		assert.Msg("timespan").This(timespan).Is(expected_timespan)
	}
	test(0, 0, 0)
	test(1, 1, 1)
	test(0b100111, 3, 7)
}
