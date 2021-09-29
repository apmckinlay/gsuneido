// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func init() {
	MakeSuTran = func(ut *UpdateTran) *rt.SuTran {
		return rt.NewSuTran(nil, true)
	}
}

func TestConcurrent(t *testing.T) {
	db := createDb()
	StartConcur(db, 50*time.Millisecond)
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
	db.ck = nil

	ck(db.Check())
	var nout = nclients * ntrans
	rt := db.NewReadTran()
	ti := rt.meta.GetRoInfo("mytable")
	assert.T(t).Msg("nrows").This(ti.Nrows).Is(nout)
	assert.T(t).Msg("size").This(ti.Size).Is(nout * 23)

	db.Close()
	ck(CheckDatabase("tmp.db"))
	os.Remove("tmp.db")
}

func TestTran(t *testing.T) {
	var err error
	db := createDb()
	db.CheckerSync()

	const nout = 4000
	for i := 0; i < nout; i++ {
		ut := output1(db)
		db.CommitMerge(ut) // commit synchronously
		if i%100 == 50 {
			if i%500 != 250 {
				db.Persist(&execPersistSingle{}, false)
			} else {
				db.Close()
				db, err = OpenDatabase("tmp.db")
				ck(err)
				db.CheckerSync()
			}
		}
	}
	db.Persist(&execPersistSingle{}, true)
	ck(db.Check())
	db.Close()

	db, err = OpenDatabaseRead("tmp.db")
	ck(err)
	rt := db.NewReadTran()
	ti := rt.meta.GetRoInfo("mytable")
	assert.T(t).Msg("nrows").This(ti.Nrows).Is(nout)
	assert.T(t).Msg("size").This(ti.Size).Is(nout * 23)
	db.Close()
	ck(CheckDatabase("tmp.db"))
	os.Remove("tmp.db")
}

func createDb() *Database {
	db, err := CreateDatabase("tmp.db")
	ck(err)
	createTbl(db)
	return db
}

func createTbl(db *Database) {
	db.Create(&schema.Schema{
		Table:   "mytable",
		Columns: []string{"one", "two"},
		Indexes: []schema.Index{{Mode: 'k', Columns: []string{"one"}}},
	})
}

var recnum int32

func output1(db *Database) *UpdateTran {
	n := atomic.AddInt32(&recnum, 1)
	ut := db.NewUpdateTran()
	data := (strconv.Itoa(int(n)) + "transaction")[:12]
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

func ck(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func TestSchemaChange(*testing.T) {
	store := stor.HeapStor(8192)
	db, err := CreateDb(store)
	ck(err)
	db.CheckerSync()
	createTbl(db)
	db.AlterCreate(&schema.Schema{
		Table:   "mytable",
		Indexes: []schema.Index{{Mode: 'i', Columns: []string{"two"}}}})

	state0 := db.GetState()
	testWith := func(fn func()) {
		ut := output1(db)
		// commit synchronously
		tables := db.ck.(*Check).commit(ut)
		ut.commit()

		fn()

		merges := &mergeList{}
		merges.add(tables)
		db.Merge(ut.meta, mergeSingle, merges)

		// restore state
		db.UpdateState(func(state *DbState) {
			*state = *state0
		})
	}
	testWith(func() {
		// no changes
	})
	testWith(func() {
		// drop table
		ck(db.Drop("mytable"))
	})
	testWith(func() {
		// modify table
		db.UpdateState(func(state *DbState) {
			state.Meta = state.Meta.TouchTable("mytable")
		})
	})
	testWith(func() {
		// drop index
		assert.That(db.AlterDrop(&schema.Schema{Table: "mytable",
			Indexes: []schema.Index{{Columns: []string{"two"}}}}))
	})
	testWith(func() {
		// modify indexes
		db.UpdateState(func(state *DbState) {
			state.Meta = state.Meta.TouchIndexes("mytable")
		})
	})
}

func TestTooMany(*testing.T) {
	store := stor.HeapStor(8192)
	db, err := CreateDb(store)
	ck(err)
	db.CheckerSync()
	for i := 0; i < maxTrans; i++ {
		assert.That(nil != db.NewUpdateTran())
	}
	assert.That(nil == db.NewUpdateTran())
}

func TestExclusive(*testing.T) {
	store := stor.HeapStor(8192)
	db, err := CreateDb(store)
	ck(err)
	db.CheckerSync()

	createTbl(db)
	assert.That(db.ck.AddExclusive("mytable"))
	ut := db.NewUpdateTran()
	assert.This(db.ck.Write(ut.ct, "mytable", []string{""})).Is(false)
	assert.This(ut.ct.conflict.Load()).Is("conflict with index creation (mytable)")
	db.ck.EndExclusive("mytable")

	ut = db.NewUpdateTran()
	assert.That(db.ck.Write(ut.ct, "mytable", []string{""}))
	assert.This(db.ck.AddExclusive("mytable")).Is(false)
	ut.Abort()

	ut = db.NewUpdateTran()
	assert.That(db.ck.Write(ut.ct, "mytable", []string{""}))
	ut.Commit()
}
