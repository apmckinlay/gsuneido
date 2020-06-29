// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"os"
	"strconv"
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"

	"github.com/apmckinlay/gsuneido/database/db19/btree"
	"github.com/apmckinlay/gsuneido/database/db19/meta"
	"github.com/apmckinlay/gsuneido/database/db19/stor"
)

func TestTran(t *testing.T) {
	btree.GetLeafKey = func(st *stor.Stor, _ interface{}, off uint64) string {
		return string(st.DataSized(off))
	}

	// store := stor.HeapStor(16 * 1024)
	store, err := stor.MmapStor("test.tmp", stor.CREATE)
	if err != nil {
		panic("can't create test.tmp")
	}

	createDb(store)

	const nout = 100
	for i := 0; i < nout; i++ {
		output1()
	}

	Persist()
	store.Close()

	f, _ := os.OpenFile("test.tmp", os.O_WRONLY|os.O_APPEND, 0644)
	f.Write([]byte("some garbage to add on the end of the file"))
	f.Close()

	store, err = stor.MmapStor("test.tmp", stor.UPDATE)
	if err != nil {
		panic("can't open test.tmp")
	}
	off := store.LastOffset([]byte(magic1))
	UpdateState(func(state *DbState) {
		*state = *ReadState(store, off)
	})

	rt := NewReadTran()
	ti := rt.meta.GetRoInfo("mytable")
	Assert(t).That(ti.Nrows, Equals(nout).Comment("nrows"))
	Assert(t).That(ti.Size, Equals(nout*12).Comment("size"))

	store.Close()
	os.Remove("test.tmp")
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

func output1() {
	ut := NewUpdateTran()
	// write some data
	data := (strconv.Itoa(ut.num) + "transaction")[:12]
	off, buf := ut.store.AllocSized(len(data))
	copy(buf, data)
	// add it to the indexes
	ti := ut.meta.GetRwInfo("mytable", ut.num)
	ti.Nrows++
	ti.Size += uint64(len(data))
	ti.Indexes[0].Insert(data, off)
	ut.Commit()
	Merge(ut.num)
}
