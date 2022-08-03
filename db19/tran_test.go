// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
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
				db.persist(&execPersistSingle{})
			} else {
				db.Close()
				db, err = OpenDatabase("tmp.db")
				ck(err)
				db.CheckerSync()
				ck(db.Check())
			}
		}
	}
	db.persist(&execPersistSingle{})
	ck(db.Check())
	rt := db.NewReadTran()
	ti := rt.meta.GetRoInfo("mytable")
	assert.T(t).Msg("nrows").This(ti.Nrows).Is(nout)
	assert.T(t).Msg("size").This(ti.Size).Is(nout * 23)
	db.Close()

	db, err = OpenDatabaseRead("tmp.db")
	ck(err)
	ck(db.Check())
	rt = db.NewReadTran()
	ti = rt.meta.GetRoInfo("mytable")
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

var recnum atomic.Int32

func output1(db *Database) *UpdateTran {
	n := recnum.Add(1)
	ut := db.NewUpdateTran()
	data := (strconv.Itoa(int(n)) + "transaction")[:12]
	ut.Output(nil, "mytable", mkrec(data, "data"))
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
	ut2 := db.NewUpdateTran()
	ut := db.NewUpdateTran()
	db.RunExclusive("mytable", func() {})
	assert.This(db.ck.Output(ut.ct, "mytable", []string{""})).Is(false)
	assert.This(ut.ct.failure.Load()).Is("conflict with exclusive (mytable)")
	// still fails because ut2 started before EndExclusive
	assert.This(db.ck.Output(ut2.ct, "mytable", []string{""})).Is(false)
	assert.This(ut2.ct.failure.Load()).Is("conflict with exclusive (mytable)")

	ut = db.NewUpdateTran()
	assert.That(db.ck.Output(ut.ct, "mytable", []string{""}))
	ut.Commit()
}

func TestRangeEnd(t *testing.T) {
	end := func(n int, flds ...string) string {
		return rangeEnd(strings.Join(flds, "\x00\x00"), n)
	}
	assert := assert.T(t).This
	assert(end(1)).Is("\x00\x00" + ixkey.Max)
	assert(end(1, "foo")).Is("foo\x00\x00" + ixkey.Max)
	assert(end(2, "foo")).Is("foo\x00\x00\x00\x00" + ixkey.Max)
	assert(end(2, "foo", "bar")).Is("foo\x00\x00bar\x00\x00" + ixkey.Max)
}

func TestOutputDupConflict(*testing.T) {
	checkerAbortT1 = true
	defer func() { checkerAbortT1 = false }()
	store := stor.HeapStor(8192)
	db, err := CreateDb(store)
	ck(err)
	db.CheckerSync()
	createTbl(db)
	t1 := db.NewUpdateTran()
	t2 := db.NewUpdateTran()
	t1.Output(nil, "mytable", mkrec("1"))
	assert.This(func() { t2.Output(nil, "mytable", mkrec("1")) }).
		Panics("conflicted")
}
