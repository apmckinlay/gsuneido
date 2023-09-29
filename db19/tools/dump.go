// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tools

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/core"
	. "github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/system"
)

const dumpVersion = "Suneido dump 3\n"
const dumpVersionPrev = "Suneido dump 2\n"
const dumpVersionBase = "Suneido dump"

// DumpDatabase exports a dumped database to a file.
// In the process it concurrently does a full check of the database.
func DumpDatabase(dbfile, to string) (nTables, nViews int, err error) {
	db, err := OpenDb(dbfile, stor.Read, false)
	ck(err)
	defer db.Close()
	return Dump(db, to)
}

func Dump(db *Database, to string) (nTables, nViews int, err error) {
	defer func() {
		if e := recover(); e != nil {
			db.Corrupt()
			err = fmt.Errorf("dump failed: %v", e)
		}
	}()
	f, w, err := dumpOpen()
	if err != nil {
		return 0, 0, err
	}
	tmpfile := f.Name()
	defer func() { f.Close(); os.Remove(tmpfile) }()
	nTables, nViews = dump(db, w)
	if err := w.Flush(); err != nil {
		return 0, 0, err
	}
	f.Close()
	ck(system.RenameBak(tmpfile, to))
	return nTables, nViews, nil
}

func dump(db *Database, w *bufio.Writer) (nTables, nViews int) {
	ics := newIndexCheckers()
	defer ics.finish()
	state := db.Persist()
	nViews = dumpViews(state, w)
	tables := make([]string, 0, 512)
	state.Meta.ForEachSchema(func(sc *meta.Schema) {
		tables = append(tables, sc.Table)
	})
	sort.Strings(tables)
	for _, table := range tables {
		dumpTable2(db, state, table, true, w, ics)
	}
	return len(tables), nViews
}

// DumpTable exports a dumped table to a file.
// It returns the number of records dumped or panics on error.
func DumpTable(dbfile, table, to string) (nrecs int, err error) {
	db, err := OpenDb(dbfile, stor.Read, false)
	ck(err)
	defer db.Close()
	return DumpDbTable(db, table, to)
}

func DumpDbTable(db *Database, table, to string) (nrecs int, err error) {
	defer func() {
		if e := recover(); e != nil {
			db.Corrupt()
			err = fmt.Errorf("dump failed: %v", e)
		}
	}()
	f, w, err := dumpOpen()
	if err != nil {
		return 0, err
	}
	tmpfile := f.Name()
	defer func() { f.Close(); os.Remove(tmpfile) }()
	nrecs = dumpDbTable(db, nrecs, table, w, f)
	if err := w.Flush(); err != nil {
		return 0, err
	}
	f.Close()
	ck(system.RenameBak(tmpfile, to))
	return nrecs, nil
}

func dumpDbTable(db *Database, nrecs int, table string, w *bufio.Writer, f *os.File) int {
	ics := newIndexCheckers()
	defer ics.finish()
	state := db.Persist()
	return dumpTable2(db, state, table, false, w, ics)
}

func dumpOpen() (*os.File, *bufio.Writer, error) {
	f, err := os.CreateTemp(".", "gs*.tmp")
	if err != nil {
		return nil, nil, err
	}
	w := bufio.NewWriter(f)
	w.WriteString(dumpVersion)
	return f, w, nil
}

func dumpTable2(db *Database, state *DbState, table string, multi bool,
	w *bufio.Writer, ics *indexCheckers) int {
	w.WriteString("====== ")
	sc := state.Meta.GetRoSchema(table)
	if sc == nil {
		panic("can't find " + table)
	}
	hasdel := sc.HasDeleted()
	schema := sc.DumpString()
	if !multi {
		schema = str.AfterFirst(schema, " ")
	}
	w.WriteString(schema + "\n")
	info := state.Meta.GetRoInfo(table)
	sum := uint64(0)
	count := info.Indexes[0].Check(func(off uint64) {
		sum += off                       // addition so order doesn't matter
		rec := OffToRecCk(db.Store, off) // verify data checksums
		if hasdel {
			rec = squeeze(rec, sc.Columns)
		}
		writeInt(w, len(rec))
		w.WriteString(string(rec))
	})
	writeInt(w, 0) // end of table records
	assert.This(count).Is(info.Nrows)
	ics.checkOtherIndexes(info, count, sum) // concurrent
	return count
}

func squeeze(rec core.Record, cols []string) core.Record {
	var rb core.RecordBuilder
	for i, col := range cols {
		if col != "-" {
			rb.AddRaw(rec.GetRaw(i))
		}
	}
	return rb.Build()
}

func writeInt(w *bufio.Writer, n int) {
	assert.That(0 <= n && n <= math.MaxUint32)
	w.WriteByte(byte(n >> 24))
	w.WriteByte(byte(n >> 16))
	w.WriteByte(byte(n >> 8))
	w.WriteByte(byte(n))
}

func dumpViews(state *DbState, w *bufio.Writer) int {
	w.WriteString("====== views (view_name,view_definition) key(view_name)\n")
	nrecs := 0
	state.Meta.ForEachView(func(name, def string) {
		var b core.RecordBuilder
		b.Add(core.SuStr(name))
		b.Add(core.SuStr(def))
		rec := b.Trim().Build()
		writeInt(w, len(rec))
		w.WriteString(string(rec))
		nrecs++
	})
	writeInt(w, 0) // end of table records
	return nrecs
}

// ------------------------------------------------------------------
// Concurrent checking of additional indexes. Also used by compact.

func newIndexCheckers() *indexCheckers {
	ics := indexCheckers{work: make(chan indexCheck, 32), // ???
		stop: make(chan void)}
	nw := options.Nworkers
	ics.wg.Add(nw)
	for i := 0; i < nw; i++ {
		go ics.worker()
	}
	return &ics
}

type void struct{}

type indexCheckers struct {
	err    atomic.Value // any
	work   chan indexCheck
	stop   chan void
	wg     sync.WaitGroup
	once   sync.Once
	closed bool
}

type indexCheck struct {
	index *index.Overlay
	count int
	sum   uint64
}

func (ics *indexCheckers) checkOtherIndexes(info *meta.Info, count int, sum uint64) {
	for i := 1; i < len(info.Indexes); i++ {
		select {
		case ics.work <- indexCheck{index: info.Indexes[i], count: count, sum: sum}:
		case <-ics.stop:
			panic("") // overridden by finish
		}
	}
}

func (ics *indexCheckers) worker() {
	defer func() {
		if e := recover(); e != nil {
			ics.err.Store(e)
			ics.once.Do(func() { close(ics.stop) }) // notify main thread
		}
		ics.wg.Done()
	}()
	for ic := range ics.work {
		CheckOtherIndex(ic.index, ic.count, ic.sum, -1)
	}
}

func (ics *indexCheckers) finish() {
	if !ics.closed {
		close(ics.work)
		ics.closed = true
	}
	ics.wg.Wait()
	if err := ics.err.Load(); err != nil {
		panic(err)
	}
}
