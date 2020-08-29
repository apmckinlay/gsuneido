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
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestConcurrent(t *testing.T) {
	db := createDb()
	StartConcur(db, 100 * time.Millisecond)
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
				ut := output1(db)
				ut.Commit()
				// time.Sleep(time.Duration(rand.Intn(900)) * time.Microsecond)
			}
		}()
	}
	wg.Wait()
	db.ck.Stop()

	var nout = nclients * ntrans
	rt := db.NewReadTran()
	ti := rt.meta.GetRoInfo("mytable")
	assert.T(t).Msg("nrows").This(ti.Nrows).Is(nout)
	assert.T(t).Msg("size").This(ti.Size).Is(nout * 23)

	db.store.Close()
}

func TestTran(t *testing.T) {
	db := createDb()
	db.ck = NewCheck()

	const nout = 2000
	for i := 0; i < nout; i++ {
		ut := output1(db)
		db.ck.Commit(ut)
		tn := ut.commit()
		db.Merge(tn)
	}

	db.Persist()
	db.store.Close()

	f, _ := os.OpenFile("test.tmp", os.O_WRONLY|os.O_APPEND, 0644)
	f.Write([]byte("some garbage to add on the end of the file"))
	f.Close()

	db.store, _ = stor.MmapStor("test.tmp", stor.UPDATE)
	off := db.store.LastOffset([]byte(magic1))
	db.UpdateState(func(state *DbState) {
		*state = *ReadState(db.store, off)
	})

	rt := db.NewReadTran()
	ti := rt.meta.GetRoInfo("mytable")
	assert.T(t).Msg("nrows").This(ti.Nrows).Is(nout)
	assert.T(t).Msg("size").This(ti.Size).Is(nout * 23)

	db.store.Close()
	os.Remove("test.tmp")
}

func createDb() *Database {
	var db Database
	db.state.set(&DbState{})

	store, err := stor.MmapStor("test.tmp", stor.CREATE)
	if err != nil {
		panic("can't create test.tmp")
	}
	db.store = store

	schema := meta.SchemaHamt{}.Mutable()
	is := ixspec.T{Cols: []int{0}}
	schema.Put(&meta.Schema{
		Table: "mytable",
		Columns: []meta.ColumnSchema{
			{Name: "one", Field: 0},
			{Name: "two", Field: 1},
		},
		Indexes: []meta.IndexSchema{{Fields: []int{0}, Ixspec: is}},
	})
	baseSchema := meta.NewSchemaPacked(store, schema.Write(store))

	info := meta.InfoHamt{}.Mutable()
	info.Put(&meta.Info{
		Table:   "mytable",
		Indexes: []*btree.Overlay{btree.NewOverlay(store, &is).Save()},
	})
	baseInfo := meta.NewInfoPacked(store, info.Write(store))

	roSchema := meta.SchemaHamt{}
	roSchemaOff := roSchema.Write(store)
	roInfo := meta.InfoHamt{}
	roInfoOff := roInfo.Write(store)

	db.UpdateState(func(state *DbState) {
		state.store = store
		state.meta = meta.NewOverlay(baseSchema, baseInfo,
			roSchema, roSchemaOff, roInfo, roInfoOff)
	})

	return &db
}

func output1(db *Database) *UpdateTran {
	ut := db.NewUpdateTran()
	data := (strconv.Itoa(ut.num()) + "transaction")[:12]
	ut.Output("mytable", mkrec(data, "data"))
	return ut
	// NOTE: does not commit
}

func mkrec(args ...string) rt.Record {
	var b rt.RecordBuilder
	for _, a := range args {
		b.Add(rt.SuStr(a))
	}
	return b.Build()
}
