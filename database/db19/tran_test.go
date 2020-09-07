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
	"github.com/apmckinlay/gsuneido/database/db19/meta/schema"
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestConcurrent(t *testing.T) {
	db := createDb()
	StartConcur(db, 100*time.Millisecond)
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

	db.Close()
	os.Remove("tmp.db")
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
	db.Close()

	db = OpenDatabaseRead("tmp.db")

	rt := db.NewReadTran()
	ti := rt.meta.GetRoInfo("mytable")
	assert.T(t).Msg("nrows").This(ti.Nrows).Is(nout)
	assert.T(t).Msg("size").This(ti.Size).Is(nout * 23)

	db.Close()
	os.Remove("tmp.db")
}

func createDb() *Database {
	db := CreateDatabase("tmp.db")
	is := ixspec.T{Cols: []int{0}}
	ts := &meta.Schema{Schema: schema.Schema{
		Table:   "mytable",
		Columns: []string{"one", "two"},
		Indexes: []schema.Index{{Fields: []int{0}, Ixspec: is}},
	}}
	ti := &meta.Info{
		Table:   "mytable",
		Indexes: []*btree.Overlay{btree.NewOverlay(db.store, &is).Save()},
	}
	db.LoadedTable(ts, ti)
	return db
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
