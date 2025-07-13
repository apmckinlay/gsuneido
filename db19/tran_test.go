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

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func init() {
	MakeSuTran = func(ut *UpdateTran) *core.SuTran {
		return core.NewSuTran(nil, true)
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
	for range nclients {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range ntrans {
				ut := output1(db)
				ut.Commit()
				// time.Sleep(time.Duration(rand.Intn(900)) * time.Microsecond)
			}
		}()
	}
	wg.Wait()
	db.ck.Stop()
	db.ck = nil

	db.MustCheck()
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
	for i := range nout {
		ut := output1(db)
		db.CommitMerge(ut) // commit synchronously
		if i%100 == 50 {
			if i%500 != 250 {
				db.persist(&execPersistSingle{}, false)
			} else {
				db.Close()
				db, err = OpenDatabase("tmp.db")
				ck(err)
				db.CheckerSync()
				db.MustCheck()
			}
		}
	}
	db.persist(&execPersistSingle{}, false)
	db.MustCheck()
	rt := db.NewReadTran()
	ti := rt.meta.GetRoInfo("mytable")
	assert.T(t).Msg("nrows").This(ti.Nrows).Is(nout)
	assert.T(t).Msg("size").This(ti.Size).Is(nout * 23)
	db.Close()

	db, err = OpenDb("tmp.db", stor.Read, true)
	ck(err)
	db.MustCheck()
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

func mkrec(args ...string) core.Record {
	var b core.RecordBuilder
	for _, a := range args {
		b.Add(core.SuStr(a))
	}
	return b.Build()
}

func ck(err error) {
	if err != nil {
		panic(err.Error())
	}
}
func TestTooMany(*testing.T) {
	db := CreateDb(stor.HeapStor(8192))
	db.CheckerSync()
	for range MaxTrans {
		assert.That(nil != db.NewUpdateTran())
	}
	assert.That(nil == db.NewUpdateTran())
}

func TestExclusive(*testing.T) {
	db := CreateDb(stor.HeapStor(8192))
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
	db := CreateDb(stor.HeapStor(8192))
	db.CheckerSync()
	createTbl(db)
	t1 := db.NewUpdateTran()
	t2 := db.NewUpdateTran()
	t1.Output(nil, "mytable", mkrec("1"))
	assert.This(func() { t2.Output(nil, "mytable", mkrec("1")) }).
		Panics("conflicted")
}

func TestGetIndexI(*testing.T) {
	db := CreateDb(stor.HeapStor(8192))
	StartConcur(db, 50*time.Millisecond)
	createTbl(db)

	ut := db.NewUpdateTran()
	it := index.NewOverIter("mytable", 0)
	it.Next(ut)                           // incorrectly got r/o info
	ut.Output(nil, "mytable", mkrec("1")) // updates r/w info
	it.Rewind()
	it.Next(ut) // output wasn't visible through r/o info
	assert.That(!it.Eof())
	key, _ := it.Cur()
	assert.This(core.Unpack(key)).Is(core.SuStr("1"))
	ut.Commit()

	db.MustCheck()
}

func TestGetIndexI2(t *testing.T) {
	db := CreateDb(stor.HeapStor(8192))
	StartConcur(db, 50*time.Millisecond)
	createTbl(db)

	ut := db.NewUpdateTran()
	ut.GetIndexI("mytable", 0) // creates mut's
	ut.Commit()                // moves mut's to layers but does not merge

	ut = db.NewUpdateTran()
	ut.Output(nil, "mytable", mkrec("1")) // merges wrong layer
	ut.Commit()

	db.MustCheck()
}

func TestUpdateUpdateSameBug(t *testing.T) {
	db := CreateDb(stor.HeapStor(8192))
	StartConcur(db, 50*time.Millisecond)
	createTbl(db)

	// Insert initial record with known size and get its offset
	ut := db.NewUpdateTran()
	initialRec := mkrec("short", "data") // small record
	ut.Output(nil, "mytable", initialRec)
	ut.Commit()

	ut = db.NewUpdateTran()

	ts := ut.getSchema("mytable")
	key := ts.Indexes[0].Ixspec.Key(initialRec)
	dbRec := ut.Lookup("mytable", 0, key)
	if dbRec == nil {
		t.Fatal("Failed to find inserted record")
	}
	recordOffset := dbRec.Off
	
	// first update
	firstUpdateRec := mkrec("short", "much_longer_data_string_here")
	ut.Update(nil, "mytable", recordOffset, firstUpdateRec)

	// second update using the original offset
	secondUpdateRec := mkrec("short", "medium_data")
	assert.This(func() {
		ut.Update(nil, "mytable", recordOffset, secondUpdateRec)
	}).Panics("update & update on same record")
	
	db.MustCheck()
}

func TestUpdateDeleteSameBug(t *testing.T) {
	db := CreateDb(stor.HeapStor(8192))
	StartConcur(db, 50*time.Millisecond)
	createTbl(db)

	ut := db.NewUpdateTran()
	initialRec := mkrec("delete", "data") // small record
	ut.Output(nil, "mytable", initialRec)
	ut.Commit()

	ut = db.NewUpdateTran()
	
	ts := ut.getSchema("mytable")
	key := ts.Indexes[0].Ixspec.Key(initialRec)
	dbRec := ut.Lookup("mytable", 0, key)
	if dbRec == nil {
		t.Fatal("Failed to find inserted record")
	}
	recordOffset := dbRec.Off

	// first update
	updateRec := mkrec("delete", "update")
	ut.Update(nil, "mytable", recordOffset, updateRec)

	// then delete using the original offset
	assert.This(func() {
	ut.Delete(nil, "mytable", recordOffset)
	}).Panics("update & delete on same record")

	db.MustCheck()
}
