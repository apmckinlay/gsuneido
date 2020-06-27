// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"math/rand"
	"testing"

	"github.com/apmckinlay/gsuneido/database/db19/stor"
	//	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestOverlay(*testing.T) {
	tbl := NewInfoHtbl(0)
	const n = 1000
	data := make([]string, n)
	randStr := str.UniqueRandom(4, 4)
	for i := 0; i < n; i++ {
		data[i] = randStr()
		tbl.Put(&Info{
			Table: data[i],
			// Schema: &TableSchema{
			// 	Table: data[i],
			// 	Columns: []ColumnSchema{
			// 		{Name: "one", Field: i},
			// 		{Name: "two", Field: i * 2},
			// 	},
			// 	Indexes: []IndexSchema{
			// 		{Fields: []int{i}},
			// 	},
			// },
		})
	}
	st := stor.HeapStor(8192)
	offInfo := tbl.Write(st)
	offSchema := tbl.Write(st)
	state := &Overlay{
		baseInfo: NewInfoPacked(st, offInfo),
		baseSchema: NewSchemaPacked(st, offSchema),
		roInfo: NewInfoHtbl(0),
	}
	// startup - nothing in memory

	for i := 0; i < 4; i++ {
		ov := state.NewOverlay()
		for i := 0; i < 5; i++ {
			ov.GetRoInfo(data[rand.Intn(100)])
			ov.GetRwInfo(data[rand.Intn(100)], 0)
		}
		// end of transaction, merge back to global
		state = ov.LayeredOnto(state)
	}

	// persist state
	// state.Write(st)
}
