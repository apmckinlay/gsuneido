// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"math/rand"
	"testing"

	"github.com/apmckinlay/gsuneido/database/db19/stor"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestOverlay(*testing.T) {
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
	offInfo := tbl.Write(st)
	offSchema := tbl.Write(st)
	state := &Meta{
		baseInfo:   NewInfoPacked(st, offInfo),
		baseSchema: NewSchemaPacked(st, offSchema),
	}
	// startup - nothing in memory

	for i := 0; i < 4; i++ {
		m := state.Mutable()
		for i := 0; i < 5; i++ {
			m.GetRoInfo(data[rand.Intn(100)])
			m.GetRwInfo(data[rand.Intn(100)], 0)
		}
		// end of transaction, merge back to global
		state = m.LayeredOnto(state)
	}

	// persist state
	// state.Write(st)
}
