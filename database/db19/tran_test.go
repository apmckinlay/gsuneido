// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/database/db19/btree"
	"github.com/apmckinlay/gsuneido/database/db19/ixspec"
	"github.com/apmckinlay/gsuneido/database/db19/meta"
	"github.com/apmckinlay/gsuneido/database/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestConcurrent(t *testing.T) {
	createDb()
	ck = StartConcur(100 * time.Millisecond)
	var nclients = 8
	var ntrans = 4000
	if testing.Short() {
		nclients = 4
		ntrans = 100
	}
	var wg sync.WaitGroup
	for i := 0; i < nclients; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < ntrans; j++ {
				ut := output1()
				ut.Commit()
				// time.Sleep(time.Duration(rand.Intn(900)) * time.Microsecond)
			}
		}()
	}
	wg.Wait()
	ck.Stop()
	ck = nil

	var nout = nclients * ntrans
	rt := NewReadTran()
	ti := rt.meta.GetRoInfo("mytable")
	assert.T(t).Msg("nrows").This(ti.Nrows).Is(nout)
	assert.T(t).Msg("size").This(ti.Size).Is(nout * 12)

	GetState().store.Close()
}

func TestTran(t *testing.T) {
	ck = NewCheck()
	defer func() { ck = nil }()
	store := createDb()

	const nout = 2000
	for i := 0; i < nout; i++ {
		ut := output1()
		ck.Commit(ut)
		tn := ut.commit()
		Merge(tn)
	}

	Persist()
	store.Close()

	f, _ := os.OpenFile("test.tmp", os.O_WRONLY|os.O_APPEND, 0644)
	f.Write([]byte("some garbage to add on the end of the file"))
	f.Close()

	store, _ = stor.MmapStor("test.tmp", stor.UPDATE)
	off := store.LastOffset([]byte(magic1))
	UpdateState(func(state *DbState) {
		*state = *ReadState(store, off)
	})

	rt := NewReadTran()
	ti := rt.meta.GetRoInfo("mytable")
	assert.T(t).Msg("nrows").This(ti.Nrows).Is(nout)
	assert.T(t).Msg("size").This(ti.Size).Is(nout * 12)

	store.Close()
	os.Remove("test.tmp")
}

func createDb() *stor.Stor {
	btree.GetLeafKey = func(st *stor.Stor, _ *ixspec.T, off uint64) string {
		return string(st.DataSized(off))
	}

	store, err := stor.MmapStor("test.tmp", stor.CREATE)
	if err != nil {
		panic("can't create test.tmp")
	}

	schema := meta.SchemaHamt{}.Mutable()
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

	info := meta.InfoHamt{}.Mutable()
	info.Put(&meta.Info{
		Table:   "mytable",
		Indexes: []*btree.Overlay{btree.NewOverlay(store, nil).Save()},
	})
	baseInfo := meta.NewInfoPacked(store, info.Write(store))

	roSchema := meta.SchemaHamt{}
	roSchemaOff := roSchema.Write(store)
	roInfo := meta.InfoHamt{}
	roInfoOff := roInfo.Write(store)

	UpdateState(func(state *DbState) {
		state.store = store
		state.meta = meta.NewOverlay(baseSchema, baseInfo,
			roSchema, roSchemaOff, roInfo, roInfoOff)
	})

	return store
}

func output1() *UpdateTran {
	ut := NewUpdateTran()
	// write some data
	data := (strconv.Itoa(ut.num()) + "transaction")[:12]
	off, buf := ut.store.AllocSized(len(data))
	copy(buf, data)
	// add it to the indexes
	ti := ut.meta.GetRwInfo("mytable", ut.num())
	ti.Nrows++
	ti.Size += uint64(len(data))
	ti.Indexes[0].Insert(data, off)
	return ut
	// NOTE: does not commit
}
