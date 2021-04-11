// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tools

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"sync"

	. "github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/cksum"
	"github.com/apmckinlay/gsuneido/util/sortlist"
)

// Compact cleans up old records and index nodes that are no longer in use.
// It does this by copying live data to a new database file.
// In the process it concurrently does a full check of the database.
func Compact(dbfile string) (ntables int, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("compact failed: %v", e)
		}
	}()
	src, err := OpenDb(dbfile, stor.READ, false)
	ck(err)
	defer src.Close()
	dst, tmpfile := tmpdb()
	defer func() { dst.Close(); os.Remove(tmpfile) }()

	state := src.GetState()
	type schemaSize struct {
		sc    *meta.Schema
		nrows int
	}
	schemas := make([]schemaSize, 0, 128)
	state.Meta.ForEachSchema(func(sc *meta.Schema) {
		ti := state.Meta.GetRoInfo(sc.Table)
		ss := schemaSize{sc: sc, nrows: ti.Nrows}
		schemas = append(schemas, ss)
		ntables++
	})
	// sort reverse to start largest first
	sort.Slice(schemas, func(i, j int) bool {
		return schemas[i].nrows > schemas[j].nrows
	})
	type compactJob struct {
		state *DbState
		src   *Database
		ts    *meta.Schema
		dst   *Database
	}
	var wg sync.WaitGroup
	channel := make(chan compactJob)
	for i := 0; i < options.Nworkers; i++ {
		wg.Add(1)
		go func() {
			for job := range channel {
				compactTable(job.state, job.src, job.ts, job.dst)
			}
			wg.Done()
		}()
	}
	for _, sc := range schemas {
		channel <- compactJob{state: state, src: src, ts: sc.sc, dst: dst}
	}
	close(channel)
	wg.Wait()
	dst.GetState().Write(true)
	dst.Close()
	src.Close()
	ck(RenameBak(tmpfile, dbfile))
	return ntables, nil
}

func tmpdb() (*Database, string) {
	dst, err := ioutil.TempFile(".", "gs*.tmp")
	ck(err)
	tmpfile := dst.Name()
	dst.Close()
	db, err := CreateDatabase(tmpfile)
	ck(err)
	return db, tmpfile
}

func compactTable(state *DbState, src *Database, ts *meta.Schema, dst *Database) {
	info := state.Meta.GetRoInfo(ts.Table)
	before := dst.Store.Size()
	list := sortlist.NewUnsorted()
	sum := uint64(0)
	count := info.Indexes[0].Check(func(off uint64) {
		sum += off // addition so order doesn't matter
		rec := src.Store.Data(off)
		size := runtime.RecLen(rec)
		rec = rec[:size+cksum.Len]
		cksum.MustCheck(rec)
		off2, buf := dst.Store.Alloc(len(rec))
		copy(buf, rec)
		//TODO squeeze records when table has deleted fields
		list.Add(off2)
	})
	list.Finish()
	assert.This(count).Is(info.Nrows)
	for i := 1; i < len(info.Indexes); i++ {
		func() {
			defer func() {
				if e := recover(); e != nil {
					fmt.Println(ts.Table, ts.Indexes[i].Columns, e)
				}
			}()
			CheckOtherIndex(info.Indexes[i], count, sum)
		}()
	}
	dataSize := dst.Store.Size() - before
	ov := buildIndexes(ts, list, dst.Store, count) // same as load
	ti := &meta.Info{Table: ts.Table, Nrows: count, Size: dataSize, Indexes: ov}
	dst.LoadedTable(ts, ti)
}
