// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tools

import (
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/stor"
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
	ics := newIndexCheckers()
	defer ics.finish()

	state := src.GetState()
	state.Meta.ForEachSchema(func(sc *meta.Schema) {
		compactTable(state, src, sc, dst, ics)
		ntables++
	})
	dst.GetState().Write(true)
	dst.Close()
	src.Close()
	ics.finish()
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

func compactTable(state *DbState, src *Database, ts *meta.Schema, dst *Database,
	ics *indexCheckers) {
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
	ics.checkOtherIndexes(info, count, sum) // concurrent
	dataSize := dst.Store.Size() - before
	ov := buildIndexes(ts, list, dst.Store, count) // same as load
	ti := &meta.Info{Table: ts.Table, Nrows: count, Size: dataSize, Indexes: ov}
	dst.LoadedTable(ts, ti)
}
