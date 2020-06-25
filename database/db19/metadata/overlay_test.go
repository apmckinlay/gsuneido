// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package metadata

import (
	"math/rand"
	"testing"

	"github.com/apmckinlay/gsuneido/database/db19/stor"
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
			Table: data[i],
			Schema: &TableSchema{
				Table: data[i],
				Columns: []ColumnSchema{
					{Name: "one", Field: i},
					{Name: "two", Field: i * 2},
				},
				Indexes: []IndexSchema{
					{Fields: []int{i}},
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
			rwMeta:       NewTableInfoHtbl(0),
			roMeta:       global,
			baseInfo:   infoPacked,
			baseSchema: schemaPacked,
		}
		for i := 0; i < 5; i++ {
			ov.GetReadonly(data[rand.Intn(100)])
			ov.GetMutable(data[rand.Intn(100)], 0)
		}
		// end of transaction, merge back to global
		global = ov.LayeredOnto(global)
	}

	// persist global info
	global.WriteInfo(st)
}
