// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"math/rand"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestMeta(*testing.T) {
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
	meta.Write(st)

	//TODO test that nothing written if no changes
	// size := st.Size()
	// meta.Write(st)
	// assert.T(t).This(st.Size()).Is(size)
}

// func TestMetaUnchanged(t *testing.T) {
// 	m := CreateMeta()
// 	offs := m.Write(nil)
// 	assert.T(t).This(offs).Is(offsets{})
// }
