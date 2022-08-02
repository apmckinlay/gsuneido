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
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/sortlist"
	"github.com/apmckinlay/gsuneido/util/system"
)

// Compact cleans up old records and index nodes that are no longer in use.
// It does this by copying live data to a new database file.
// In the process it concurrently does a full check of the database.
func Compact(dbfile string) (nTables, nViews int, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("compact failed: %v", e)
		}
	}()
	src, err := OpenDb(dbfile, stor.Read, false)
	ck(err)
	defer src.Close()
	dst, tmpfile := tmpdb()
	defer func() { dst.Close(); os.Remove(tmpfile) }()

	state := src.GetState()
	type schemaSize struct {
		sc    *meta.Schema
		nrows int
	}

	nViews = copyViews(state, dst)

	schemas := make([]schemaSize, 0, 128)
	state.Meta.ForEachSchema(func(sc *meta.Schema) {
		ti := state.Meta.GetRoInfo(sc.Table)
		ss := schemaSize{sc: sc, nrows: ti.Nrows}
		schemas = append(schemas, ss)
		nTables++
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
	dst.GetState().Write()
	dst.Close()
	src.Close()
	ck(system.RenameBak(tmpfile, dbfile))
	return nTables, nViews, nil
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

func copyViews(state *DbState, dst *Database) int {
	n := 0
	state.Meta.ForEachView(func(name, def string) {
		dst.AddView(name, def)
		n++
	})
	return n
}

func compactTable(state *DbState, src *Database, ts *meta.Schema, dst *Database) {
	defer func() {
		if e := recover(); e != nil {
			runtime.Fatal(ts.Table+":", e)
		}
	}()
	hasdel := ts.HasDeleted()
	info := state.Meta.GetRoInfo(ts.Table)
	sum := uint64(0)
	size := uint64(0)
	list := sortlist.NewUnsorted()
	var off2 uint64
	var buf []byte
	var n int
	count := info.Indexes[0].Check(func(off uint64) {
		sum += off // addition so order doesn't matter
		if hasdel {
			rec := OffToRecCk(src.Store, off) // verify data checksums
			rec = squeeze(rec, ts.Columns)
			n = len(rec)
			off2, buf = dst.Store.Alloc(n + cksum.Len)
			copy(buf, rec)
			cksum.Update(buf)
		} else {
			rec := src.Store.Data(off)
			n = runtime.RecLen(rec)
			rec = rec[:n+cksum.Len]
			cksum.MustCheck(rec)
			off2, buf = dst.Store.Alloc(len(rec))
			copy(buf, rec)
		}
		list.Add(off2)
		size += uint64(n)
	})
	list.Finish()
	assert.This(count).Is(info.Nrows)
	for i := 1; i < len(info.Indexes); i++ {
		CheckOtherIndex(info.Indexes[i], count, sum, i)
	}
	if hasdel {
		ts.Columns = slc.Without(ts.Columns, "-")
	}
	ovs := buildIndexes(ts, list, dst.Store, count) // same as load
	ti := &meta.Info{Table: ts.Table, Nrows: count, Size: size, Indexes: ovs}
	dst.LoadedTable(ts, ti)
}
