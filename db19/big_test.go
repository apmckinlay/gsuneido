// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/apmckinlay/gsuneido/db19/btree"
	"github.com/apmckinlay/gsuneido/db19/ixspec"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/hash"
	"github.com/apmckinlay/gsuneido/util/str"
)

const (
	dbfile        = "tmp.db"
	ntables       = 1009
	maxcols       = 211
	maxidxs       = 19
	maxidxcols    = 5
	nrows         = 500_000
	nthreads      = 5
	tablesPerTran = 53
)

func TestBig(*testing.T) {
	if testing.Short() {
		return
	}
	fmt.Println("create tables")
	tables := createTables()
	// defer os.Remove(dbfile)
	db, err := OpenDatabase(dbfile)
	ck(err)
	StartConcur(db, 777*time.Millisecond)
	fmt.Println("create data")
	var wg sync.WaitGroup
	count := nrows / nthreads
	start := 0
	for i := 0; i < nthreads; i++ {
		wg.Add(1)
		go func(start int) {
			createData(db, tables, start, count)
			wg.Done()
		}(start)
		start += count
	}
	wg.Wait()

	db.ck.Stop()
	db.ck = nil
	nr := 0
	state := db.GetState()
	state.meta.ForEachInfo(func(ti *meta.Info) {
		nr += ti.Nrows
	})
	assert.This(nr).Is(nrows)

	db.Close()
	fmt.Println("check")
	err = CheckDatabase(dbfile)
	ck(err)

	fmt.Println("merge1", mergeTime1/time.Duration(mergeCount))
	fmt.Println("merge2", mergeTime2/time.Duration(mergeCount))
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
		idxInfo := make([]*btree.Overlay, nidxs)
		for j := 0; j < nidxs; j++ {
			var idxcols []string
			var mode int
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
			idxInfo[j] = btree.NewOverlay(db.store, &ixspec.T{})
			idxInfo[j].Save(false)
		}
		schema := schema.Schema{Table: table, Columns: cols, Indexes: idxSchema}
		ts := &meta.Schema{Schema: schema}
		ti := &meta.Info{Table: table, Indexes: idxInfo}
		db.LoadedTable(ts, ti)
	}
	db.Persist(true)
	db.Close()
	return tables
}

func createData(db *Database, tables []string, i, n int) {
	n += i
	for i < n {
		ut := db.NewUpdateTran()
		ntables := 1 + rand.Intn(tablesPerTran)
		for j := 0; j < ntables && i < n; j++ {
			table := tables[rand.Intn(ntables)]
			ncols := fromHash(table, maxcols)
			var b rt.RecordBuilder
			b.Add(rt.IntVal(i).(rt.Packable))
			for k := 1; k < ncols; k++ {
				b.Add(rt.SuStr(str.Random(0, 20)))
			}
			rec := b.Build()
			ut.Output(table, rec)
			i++
		}
		// time.Sleep(4 * time.Millisecond) // inside tran
		ut.Commit()
		// time.Sleep(4 * time.Millisecond) // between tran
	}
}

func fromHash(table string, max int) int {
	return int(1 + hash.HashString(table)%uint32(max))
}
