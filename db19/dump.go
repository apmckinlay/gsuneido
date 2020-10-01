// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/db19/btree"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/ints"
)

// DumpDatabase exports a dumped database to a file.
// In the process it does a full check of the database.
// It returns the number of tables dumped or panics on error.
func DumpDatabase(dbfile, to string) int {
	defer func() {
		if e := recover(); e != nil {
			os.Remove(to)
			panic("dump failed: " + fmt.Sprint(e))
		}
	}()
	db, f, w := dumpOpen(dbfile, to)
	defer db.Close()
	defer f.Close()
	ics := NewIndexCheckers()
	defer ics.finish("dump")

	state := db.GetState()
	ntables := 0
	state.meta.ForEachSchema(func(sc *meta.Schema) {
		dumpTable(db, sc, true, w, ics)
		ntables++
		if atomic.LoadInt32(&ics.err) != 0 {
			panic("dump failed: database corrupt?")
		}
	})
	ck(w.Flush())
	return ntables
}

// DumpTable exports a dumped table to a file.
// It returns the number of records dumped or panics on error.
func DumpTable(dbfile, table, to string) int {
	defer func() {
		if e := recover(); e != nil {
			os.Remove(to)
			panic("dump failed: " + table + " " + fmt.Sprint(e))
		}
	}()
	db, f, w := dumpOpen(dbfile, to)
	defer db.Close()
	defer f.Close()
	ics := NewIndexCheckers()
	defer ics.finish("dump")

	state := db.GetState()
	schema := state.meta.GetRoSchema(table)
	if schema == nil {
		panic("can't find " + table)
	}
	n := dumpTable(db, schema, false, w, ics)
	ck(w.Flush())
	return n
}

func dumpOpen(dbfile, to string) (*Database, *os.File, *bufio.Writer) {
	db, err := openDatabase(dbfile, stor.READ, false)
	ck(err)
	os.Remove(to) //TODO bak
	f, err := os.Create(to)
	ck(err)
	w := bufio.NewWriter(f)
	w.WriteString("Suneido dump 2\n")
	return db, f, w
}

func dumpTable(db *Database, schema *meta.Schema, multi bool, w *bufio.Writer,
	ics *indexCheckers) int {
	state := db.GetState()
	w.WriteString("====== ")
	if multi {
		w.WriteString(schema.Table + " ")
	}
	w.WriteString(schema.String() + "\n")
	info := state.meta.GetRoInfo(schema.Table)
	sum := uint64(0)
	count := info.Indexes[0].Check(func(off uint64) {
		sum += off                       // addition so order doesn't matter
		rec := offToRecCk(db.store, off) // verify data checksums
		writeInt(w, len(rec))
		w.WriteString(string(rec))
		//TODO squeeze records when table has deleted fields
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

// ------------------------------------------------------------------
// Concurrent checking of additional indexes. Also used by compact.

func NewIndexCheckers() *indexCheckers {
	var ics indexCheckers
	ics.in = make(chan indexCheck, 32)                       // ???
	nworkers := ints.Min(8, ints.Max(1, runtime.NumCPU()-1)) // ???
	ics.wg.Add(nworkers)
	for i := 0; i < nworkers; i++ {
		go ics.worker()
	}
	return &ics
}

type indexCheckers struct {
	wg sync.WaitGroup
	in chan indexCheck
	// indexError is set to non-zero when an error is detected.
	// It must be accessed atomically.
	err int32
}

type indexCheck struct {
	index *btree.Overlay
	count int
	sum   uint64
}

func (ics *indexCheckers) checkOtherIndexes(info *meta.Info, count int, sum uint64) {
	for i := 1; i < len(info.Indexes); i++ {
		ics.in <- indexCheck{index: info.Indexes[i], count: count, sum: sum}
	}
}

func (ics *indexCheckers) worker() {
	defer func() {
		if e := recover(); e != nil {
			atomic.StoreInt32(&ics.err, 1)
		}
		ics.wg.Done()
	}()
	for ic := range ics.in {
		checkOtherIndex("", ic.index, ic.count, ic.sum)
		if atomic.LoadInt32(&ics.err) != 0 {
			break
		}
	}
}

func (ics *indexCheckers) finish(which string) {
	if ics.in != nil {
		close(ics.in)
		ics.in = nil
	}
	ics.wg.Wait()
	if atomic.LoadInt32(&ics.err) != 0 {
		panic(which + " failed: database corrupt?")
	}
}
