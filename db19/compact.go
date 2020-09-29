// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/stor"
	rt "github.com/apmckinlay/gsuneido/runtime"
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
	src, err := OpenDatabaseRead(dbfile)
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
	iter := info.Indexes[0].Iter(true)
	nrecs := 0
	for {
		_, off, ok := iter()
		if !ok {
			break
		}
		rec := offToBufCk(src.store, off) // verify data checksums
		off, buf := dst.store.Alloc(len(rec))
		copy(buf, rec)
		list.Add(off)
		nrecs++
	}
	list.Finish()
	dataSize := dst.store.Size() - before
	ov := buildIndexes(ts, list, dst.store, nrecs) // same as load
	ti := &meta.Info{Table: ts.Table, Nrows: nrecs, Size: dataSize, Indexes: ov}
	dst.LoadedTable(ts, ti)
}

func offToBufCk(store *stor.Stor, off uint64) []byte {
	buf := store.Data(off)
	size := rt.RecLen(buf)
	buf = buf[:size+cksum.Len]
	cksum.Check(buf)
	return buf
}
