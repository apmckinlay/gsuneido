// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"

	"github.com/apmckinlay/gsuneido/database/db19/btree"
	"github.com/apmckinlay/gsuneido/database/db19/meta"
	"github.com/apmckinlay/gsuneido/database/db19/stor"
)

func TestTran(t *testing.T) {
	btree.GetLeafKey = func(st *stor.Stor, ixspec interface{}, off uint64) string {
		return string(st.DataSized(off))
	}
	store := stor.HeapStor(16 * 1024)
	createDb(store)

	ut := NewUpdateTran()
	// write some data
	off, buf := store.AllocSized(10)
	copy(buf, "helloworld")
	// add it to the indexes
	ti := ut.meta.GetRwInfo("mytable", ut.num)
	Assert(t).That(ti.Nrows, Equals(0).Comment("nrows"))
	Assert(t).That(ti.Size, Equals(0).Comment("size"))
	Assert(t).That(len(ti.Indexes), Equals(1).Comment("n indexes"))
	ti.Indexes[0].Insert("helloworld", off)
	ut.Commit()

	off = Persist()

	ReadState(store, off)
}

func createDb(store *stor.Stor) {
	schema := meta.NewSchemaHtbl(0)
	schema.Put(&meta.Schema{
		Table: "mytable",
		Columns: []meta.ColumnSchema{
			{Name: "one", Field: 0},
			{Name: "two", Field: 1},
		},
		Indexes: []meta.IndexSchema{
			{Fields: []int{0}},
		},
	})
	baseSchema := meta.NewSchemaPacked(store, schema.Write(store))

	info := meta.NewInfoHtbl(0)
	info.Put(&meta.Info{
		Table:   "mytable",
		Indexes: []*btree.Overlay{btree.NewOverlay(store).Save()},
	})
	baseInfo := meta.NewInfoPacked(store, info.Write(store))

	roSchema := meta.NewSchemaHtbl(0)
	roSchemaOff := roSchema.Write(store)
	roInfo := meta.NewInfoHtbl(0)
	roInfoOff := roInfo.Write(store)

	UpdateState(func(state *DbState) {
		state.store = store
		state.meta = meta.NewOverlay(baseSchema, baseInfo,
			roSchema, roSchemaOff, roInfo, roInfoOff, nil)
	})
}
