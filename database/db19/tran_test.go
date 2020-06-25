// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"

	"github.com/apmckinlay/gsuneido/database/db19/btree"
	md "github.com/apmckinlay/gsuneido/database/db19/metadata"
	"github.com/apmckinlay/gsuneido/database/db19/stor"
)

func TestTran(t *testing.T) {
	btree.GetLeafKey = func(st *stor.Stor, ixspec interface{}, off uint64) string {
		return string(st.DataSized(off))
	}
	store := stor.HeapStor(8192)
	createDb(store)

	ut := NewUpdateTran()
	// write some data
	off, buf := store.AllocSized(10)
	copy(buf, "helloworld")
	// add it to the indexes
	ti := ut.meta.GetMutable("mytable", ut.num)
	Assert(t).That(ti.Nrows, Equals(0).Comment("nrows"))
	Assert(t).That(ti.Size, Equals(0).Comment("size"))
	Assert(t).That(len(ti.Indexes), Equals(1).Comment("n indexes"))
	ti.Indexes[0].Insert("helloworld", off)
	ut.Commit()

	off = Persist()

	// ReadState(store, off)
}

func createDb(store *stor.Stor) {
	tbl := md.NewTableInfoHtbl(0)
	tbl.Put(&md.TableInfo{
		Table: "mytable",
		Schema: &md.TableSchema{
			Table: "mytable",
			Columns: []md.ColumnSchema{
				{Name: "one", Field: 0},
				{Name: "two", Field: 1},
			},
			Indexes: []md.IndexSchema{
				{Fields: []int{0}},
			},
		},
		Indexes: []*btree.Overlay{btree.NewOverlay(store).Save()},
	})
	schemaOff := tbl.WriteSchema(store)
	schemaPacked := md.NewSchemaPacked(store, schemaOff)
	infoOff := tbl.WriteInfo(store)
	infoPacked := md.NewInfoPacked(store, infoOff)
	UpdateState(func(state *DbState) {
		state.store = store
		state.baseSchema = schemaPacked
		state.baseInfo = infoPacked
		state.memMeta = md.NewTableInfoHtbl(0)
	})
}
