// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/hash"
	"github.com/apmckinlay/gsuneido/util/str"
)

const (
	dbfile          = "tmp.db"
	ntables         = 1009
	maxcols         = 211
	maxidxs         = 11
	maxidxcols      = 5
	nrows           = 1_000_000
	nthreads        = 4 // must divide evenly into nrows
	tablesPerTran   = 7
	rowsPerTable    = 7
	persistInterval = 500 * time.Millisecond
)

func TestBig(*testing.T) {
	if testing.Short() {
		return
	}
	assert.That(nrows%nthreads == 0)
	fmt.Println("create tables")
	tables := createTables()
	defer os.Remove(dbfile)
	db, err := OpenDatabase(dbfile)
	ck(err)
	StartConcur(db, persistInterval)
	fmt.Println("create data")
	count := nrows / nthreads
	start := 0
	var wg sync.WaitGroup
	for i := 0; i < nthreads; i++ {
		wg.Add(1)
		go func(start int) {
			createData(db, tables, start, count)
			wg.Done()
		}(start)
		start += count
	}
	wg.Wait()
	fmt.Println("finished", ntrans.Load(), "transactions", db.Store.Size(), "bytes")

	db.ck.Stop()
	db.ck = nil

	db.MustCheck()

	nr := 0
	state := db.GetState()
	state.Meta.CheckAllMerged()
	state.Meta.ForEachInfo(func(ti *meta.Info) {
		nr += ti.Nrows
	})
	assert.This(nr).Is(nrows)
	db.MustCheck()

	db.Close()
	ck(CheckDatabase(dbfile))
}

func createTables() []string {
	db, err := CreateDatabase(dbfile)
	ck(err)
	tables := make([]string, ntables)
	randTable := str.UniqueRandom(4, 16)
	for i := 0; i < ntables; i++ {
		table := randTable()
		tables[i] = table
		randCol := str.UniqueRandom(4, 16)
		ncols := fromHash(table, maxcols)
		cols := make([]string, ncols)
		for j := 0; j < ncols; j++ {
			cols[j] = randCol()
		}
		nidxs := fromHash(table, maxidxs)
		idxSchema := make([]schema.Index, nidxs)
		idxInfo := make([]*index.Overlay, nidxs)
		for j := 0; j < nidxs; j++ {
			var idxcols []string
			var mode byte
			if j == 0 {
				idxcols = []string{cols[0]}
				mode = 'k'
			} else {
				nidxcols := fromHash(table, (maxidxcols))
				idxcols := make([]string, nidxcols)
				for k := 0; k < nidxcols; k++ {
					f := fromHash(table, ncols) - 1
					idxcols[k] = cols[f]
				}
				mode = 'i'
			}
			idxSchema[j] = schema.Index{Columns: idxcols, Mode: mode}
			idxInfo[j] = index.NewOverlay(db.Store, &ixkey.Spec{})
			idxInfo[j].Save()
		}
		schema := schema.Schema{Table: table, Columns: cols, Indexes: idxSchema}
		ts := &meta.Schema{Schema: schema}
		ts.SetBestKeys(0)
		ti := &meta.Info{Table: table, Indexes: idxInfo}
		db.AddNewTable(ts, ti)
	}
	db.Close()
	return tables
}

var ntrans atomic.Int32

var data = core.SuStr(str.Random(1024, 1024))

func createData(db *Database, tables []string, i, n int) {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	n += i
	for i < n {
		ut := db.NewUpdateTran()
		nt := 1 + rand.Intn(tablesPerTran)
		for j := 0; j < nt && i < n; j++ {
			table := tables[rand.Intn(ntables)]
			nr := 1 + rand.Intn(rowsPerTable)
			for k := 0; k <= nr && i < n; k++ {
				var b core.RecordBuilder
				b.Add(core.IntVal(i ^ 0x5555))
				b.Add(data)
				rec := b.Build()
				ut.Output(nil, table, rec)
				i++
			}
		}
		// time.Sleep(1 * time.Millisecond) // inside tran
		ut.Commit()
		ntrans.Add(1)
		// time.Sleep(1 * time.Millisecond) // between tran
	}
}

func fromHash(table string, max int) int {
	return int(1 + hash.String(table)%uint32(max))
}
