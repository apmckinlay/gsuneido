// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tools

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sync"
	"sync/atomic"

	. "github.com/apmckinlay/gsuneido/db19"
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/options"
	rt "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

// DumpDatabase exports a dumped database to a file.
// In the process it concurrently does a full check of the database.
func DumpDatabase(dbfile, to string) (ntables int, err error) {
	db, err := OpenDb(dbfile, stor.READ, false)
	ck(err)
	defer db.Close()
	return Dump(db, to)
}

func Dump(db *Database, to string) (ntables int, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("dump failed: %v", e)
		}
	}()
	f, w := dumpOpen()
	tmpfile := f.Name()
	defer func() { db.Close(); f.Close(); os.Remove(tmpfile) }()
	ics := newIndexCheckers()
	defer ics.finish()

	state := db.GetState()
	dumpViews(state, w)
	state.Meta.ForEachSchema(func(sc *meta.Schema) {
		dumpTable2(db, sc, true, w, ics)
		ntables++
	})
	ck(w.Flush())
	f.Close()
	ics.finish()
	ck(RenameBak(tmpfile, to))
	return ntables, nil
}

// DumpTable exports a dumped table to a file.
// It returns the number of records dumped or panics on error.
func DumpTable(dbfile, table, to string) (nrecs int, err error) {
	db, err := OpenDb(dbfile, stor.READ, false)
	ck(err)
	defer db.Close()
	return DumpDbTable(db, table, to)
}

func DumpDbTable(db *Database, table, to string) (nrecs int, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("dump failed: %v", e)
		}
	}()
	f, w := dumpOpen()
	tmpfile := f.Name()
	defer func() { f.Close(); os.Remove(tmpfile) }()
	ics := newIndexCheckers()
	defer ics.finish()

	state := db.GetState()
	schema := state.Meta.GetRoSchema(table)
	if schema == nil {
		return 0, errors.New("dump failed: can't find " + table)
	}
	nrecs = dumpTable2(db, schema, false, w, ics)
	ck(w.Flush())
	f.Close()
	ics.finish()
	ck(RenameBak(tmpfile, to))
	return nrecs, nil

}

func dumpOpen() (*os.File, *bufio.Writer) {
	f, err := ioutil.TempFile(".", "gs*.tmp")
	ck(err)
	w := bufio.NewWriter(f)
	w.WriteString("Suneido dump 2\n")
	return f, w
}

func dumpTable2(db *Database, schema *meta.Schema, multi bool, w *bufio.Writer,
	ics *indexCheckers) int {
	state := db.GetState()
	w.WriteString("====== ")
	s := schema.String()
	if !multi {
		s = str.AfterFirst(s, " ")
	}
	w.WriteString(s + "\n")
	info := state.Meta.GetRoInfo(schema.Table)
	sum := uint64(0)
	count := info.Indexes[0].Check(func(off uint64) {
		sum += off                       // addition so order doesn't matter
		rec := OffToRecCk(db.Store, off) // verify data checksums
		writeInt(w, len(rec))
		w.WriteString(string(rec))
	})
	writeInt(w, 0) // end of table records
	assert.This(count).Is(info.Nrows)
	ics.checkOtherIndexes(info, count, sum) // concurrent
	return count
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
		var b rt.RecordBuilder
		b.Add(rt.SuStr(name))
		b.Add(rt.SuStr(def))
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
	ics := indexCheckers{work: make(chan indexCheck, 32)} // ???
	nw := options.Nworkers
	ics.wg.Add(nw)
	for i := 0; i < nw; i++ {
		go ics.worker()
	}
	return &ics
}

type void struct{}

type indexCheckers struct {
	err    atomic.Value
	work   chan indexCheck
	stop   chan void
	wg     sync.WaitGroup
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
			close(ics.stop) // notify main thread
		}
		ics.wg.Done()
	}()
	for ic := range ics.work {
		CheckOtherIndex(ic.index, ic.count, ic.sum)
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
