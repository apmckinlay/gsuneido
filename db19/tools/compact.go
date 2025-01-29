// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tools

import (
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/apmckinlay/gsuneido/core"
	. "github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/cksum"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/sortlist"
	"github.com/apmckinlay/gsuneido/util/system"
)

// Compact cleans up old records and index nodes that are no longer in use.
// It does this by copying live data to a new database file.
// In the process it concurrently does a full check of the database.
func Compact(dbfile string) (nTables, nViews int, oldSize, newSize uint64, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("compact failed: %v", e)
		}
	}()
	src, err := OpenDb(dbfile, stor.Read, false)
	ck(err)
	defer src.Close()
	oldSize = src.Store.Size()
	dst, tmpfile := tmpdb()
	defer func() { dst.Close(); os.Remove(tmpfile) }()

	state := src.GetState()
	type schemaSize struct {
		sc    *meta.Schema
		nrows int
	}

	nViews = copyViews(state, dst)

	schemas := make([]schemaSize, 0, 128)
	for sc := range state.Meta.Tables() {
		ti := state.Meta.GetRoInfo(sc.Table)
		ss := schemaSize{sc: sc, nrows: ti.Nrows}
		schemas = append(schemas, ss)
		nTables++
	}
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
	for range options.Nworkers {
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
	newSize = dst.Store.Size()
	dst.Close()
	src.Close()
	ck(system.RenameBak(tmpfile, dbfile))
	return nTables, nViews, oldSize, newSize, nil
}

func tmpdb() (*Database, string) {
	dst, err := os.CreateTemp(".", "gs*.tmp")
	ck(err)
	tmpfile := dst.Name()
	dst.Close()
	db, err := CreateDatabase(tmpfile)
	ck(err)
	return db, tmpfile
}

func copyViews(state *DbState, dst *Database) int {
	n := 0
	for name, def := range state.Meta.Views() {
		dst.AddView(name, def)
		n++
	}
	return n
}

func compactTable(state *DbState, src *Database, ts *meta.Schema, dst *Database) {
	defer func() {
		if e := recover(); e != nil {
			core.Fatal(ts.Table+":", e)
		}
	}()
	hasdel := ts.HasDeleted()
	info := state.Meta.GetRoInfo(ts.Table)
	sum := uint64(0)
	size := int64(0)
	list := sortlist.NewUnsorted(func(x uint64) bool { return x == 0 })
	var off2 uint64
	var dstbuf []byte
	nrows := info.Indexes[0].Check(func(off uint64) {
		sum += off // addition so order doesn't matter
		buf := src.Store.Data(off)
		n := core.RecLen(buf)
		buf = buf[:n+cksum.Len]
		cksum.MustCheck(buf)
		rec := core.Record(hacks.BStoS(buf[:n]))
		if hasdel || hasTrailingEmpty(rec) {
			rec = squeeze(rec, ts.Columns)
			n = len(rec)
			off2, dstbuf = dst.Store.Alloc(n + cksum.Len)
			copy(dstbuf, rec)
			cksum.Update(dstbuf)
		} else {
			off2, dstbuf = dst.Store.Alloc(len(buf))
			copy(dstbuf, buf)
		}
		list.Add(off2)
		size += int64(n)
	})
	list.Finish()
	assert.This(nrows).Is(info.Nrows)
	for i := 1; i < len(info.Indexes); i++ {
		CheckOtherIndex(ts.Indexes[i].Columns, info.Indexes[i], nrows, sum)
	}
	if hasdel {
		ts.Columns = slc.Without(ts.Columns, "-")
	}
	ovs := buildIndexes(ts, list, dst.Store, nrows) // same as load
	ti := &meta.Info{Table: ts.Table, Nrows: nrows, Size: size, Indexes: ovs}
	dst.AddNewTable(ts, ti)
}

func hasTrailingEmpty(r core.Record) bool {
	n := r.Count()
	return n > 0 && r.GetRaw(n-1) == ""
}
