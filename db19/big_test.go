// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	"math/rand"
	"os"
	"slices"
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
	nrows           = 2_000_000
	nthreads        = 8 // must divide evenly into nrows
	tablesPerTran   = 7
	rowsPerTable    = 7
	persistInterval = 500 * time.Millisecond
)

func TestBig(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
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
	var wg sync.WaitGroup
	for i := range nthreads {
		wg.Go(func() {
			createData(db, tables, i*count, count)
		})
	}
	wg.Wait()
	fmt.Println("finished", ntrans.Load(), "transactions", db.Store.Size(), "bytes")

	db.ck.Stop()
	db.ck = nil

	db.MustCheck()

	nr := 0
	state := db.GetState()
	state.Meta.CheckAllMerged()
	for ti := range state.Meta.Infos() {
		nr += ti.Nrows
	}
	assert.This(nr).Is(nrows)
	db.MustCheck()

	db.Close()
	ck(CheckDatabase(dbfile, true))
	PrintStates(dbfile, true)
}

func createTables() []string {
	db, err := CreateDatabase(dbfile)
	ck(err)
	tables := make([]string, ntables)
	randTable := str.UniqueRandom(4, 16)
	for i := range ntables {
		table := randTable()
		tables[i] = table
		randCol := str.UniqueRandom(4, 16)
		ncols := fromHash(table, maxcols)
		cols := make([]string, ncols)
		for j := range ncols {
			cols[j] = randCol()
		}
		nidxs := fromHash(table, maxidxs)
		var idxSchema []schema.Index
		var idxInfo []*index.Overlay
		for j := range nidxs {
			var idxcols []string
			var mode byte
			if j == 0 {
				idxcols = []string{cols[0]}
				mode = 'k'
			} else {
				nidxcols := fromHash(table, maxidxcols)
				idxcols = make([]string, nidxcols)
				for k := range nidxcols {
					f := fromHash(table, ncols) - 1
					idxcols[k] = cols[f]
				}
				mode = 'i'
				if isdup(idxcols, idxSchema) {
					continue
				}
			}
			idxSchema = append(idxSchema, schema.Index{Columns: idxcols, Mode: mode})
			ov := index.NewOverlay(db.Store, &ixkey.Spec{})
			idxInfo = append(idxInfo, ov)
		}
		schema := schema.Schema{Table: table, Columns: cols, Indexes: idxSchema}
		schema.Check()
		ts := &meta.Schema{Schema: schema}
		ts.SetBestKeys(0)
		ti := &meta.Info{Table: table, Indexes: idxInfo}
		db.AddNewTable(ts, ti)
	}
	db.PersistClose()
	return tables
}

func isdup(cols []string, idxSchema []schema.Index) bool {
	for _, ix := range idxSchema {
		if slices.Equal(ix.Columns, cols) {
			return true
		}
	}
	return false
}

var ntrans atomic.Int32

var data = core.SuStr(str.Random(512, 512))

func createData(db *Database, tables []string, i, n int) {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	n += i
	for i < n {
		ut := db.NewUpdateTran()
		nt := 1 + rand.Intn(tablesPerTran)
		for range min(nt, n) {
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
	return int(1 + hash.String(table)%uint64(max))
}
