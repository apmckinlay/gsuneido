// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package metadata

import (
	"math/rand"
	"testing"

	"github.com/apmckinlay/gsuneido/database/stor"
	//	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestOverlay(*testing.T) {
	tbl := NewTableInfoHtbl(0)
	const n = 1000
	data := make([]string, n)
	randStr := str.UniqueRandom(4, 4)
	for i := 0; i < n; i++ {
		data[i] = randStr()
		tbl.Put(&TableInfo{
			table: data[i],
			schema: &TableSchema{
				table: data[i],
				columns: []ColumnSchema{
					{name: "one", field: i},
					{name: "two", field: i * 2},
				},
				indexes: []IndexSchema{
					{fields: []int{i}},
				},
			},
		})
	}
	st := stor.HeapStor(2 * blockSize)
	offInfo := tbl.WriteInfo(st)
	infoPacked := NewInfoPacked(st, offInfo)
	offSchema := tbl.WriteSchema(st)
	schemaPacked := NewSchemaPacked(st, offSchema)
	global := NewTableInfoHtbl(0)
	// startup - nothing in memory

	for i := 0; i < 4; i++ {
		ov := &Overlay{
			rwInfo:       NewTableInfoHtbl(0),
			roInfo:       global,
			infoPacked:   infoPacked,
			schemaPacked: schemaPacked,
		}
		for i := 0; i < 5; i++ {
			ov.GetReadonly(data[rand.Intn(100)])
			ov.GetMutable(data[rand.Intn(100)])
		}
		// end of transaction, merge back to global
		global = global.Merge(ov.rwInfo)
	}

	// persist global info
	global.WriteInfo(st)
}
