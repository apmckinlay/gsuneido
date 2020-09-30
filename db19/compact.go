// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/cksum"
	"github.com/apmckinlay/gsuneido/util/sortlist"
)

// Compact cleans up old records and index nodes that are no longer used.
// It does this by copying to a new database file.
func Compact(dbfile string) int {
	defer func() {
		if e := recover(); e != nil {
			panic("compact failed: " + fmt.Sprint(e))
		}
	}()
	src, err := openDatabase(dbfile, stor.READ, false)
	ck(err)
	defer src.Close()
	dst, tmpfile := tmpdb()
	defer func() { dst.Close(); os.Remove(tmpfile) }()

	state := src.GetState()
	ntables := 0
	state.meta.ForEachSchema(func(sc *meta.Schema) {
		compactTable(state, src, sc, dst)
		ntables++
	})
	dst.GetState().Write()
	dst.Close()
	src.Close()
	ck(renameBak(tmpfile, dbfile))

	return ntables
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
	info := state.meta.GetRoInfo(ts.Table)
	before := dst.store.Size()
	list := sortlist.NewUnsorted()
	sum := uint64(0)
	count := info.Indexes[0].Check(func(off uint64) {
		sum += off // addition so order doesn't matter
		rec := src.store.Data(off)
		size := runtime.RecLen(rec)
		rec = rec[:size+cksum.Len]
		if !cksum.Check(rec) {
			panic(&ErrCorrupt{table: info.Table})
		}
		off2, buf := dst.store.Alloc(len(rec))
		copy(buf, rec)
		list.Add(off2)
	})
	list.Finish()
	assert.This(count).Is(info.Nrows)
	dataSize := dst.store.Size() - before
	checkOtherIndexes(info, count, sum)            //TODO concurrent
	ov := buildIndexes(ts, list, dst.store, count) // same as load
	ti := &meta.Info{Table: ts.Table, Nrows: count, Size: dataSize, Indexes: ov}
	dst.LoadedTable(ts, ti)
}

func checkOtherIndexes(info *meta.Info, count int, sum uint64) {
	for i := 1; i < len(info.Indexes); i++ {
		count, sum = checkOtherIndex(info.Table, info.Indexes[i], count, sum)
	}
}
