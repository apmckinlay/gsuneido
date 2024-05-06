// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/cksum"
)

// quick check ------------------------------------------------------

// QuickCheck is the default partial checking done at start up.
func (db *Database) QuickCheck() (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("check failed: %v", e)
		}
	}()
	// t := time.Now()
	runParallel(db.GetState(), quickCheckTable)
	// fmt.Println("quick checked in", time.Since(t).Round(time.Millisecond))
	return nil
}

func quickCheckTable(state *DbState, table string) {
	info := state.Meta.GetRoInfo(table)
	for _, ix := range info.Indexes {
		ix.QuickCheck()
	}
}

// full check -------------------------------------------------------

// CheckDatabase is called by -check and -repair
func CheckDatabase(dbfile string) (ec error) {
	db, err := OpenDb(dbfile, stor.Read, false)
	if err != nil {
		return newErrCorrupt(err)
	}
	defer db.Close()
	defer func() {
		if e := recover(); e != nil {
			db.Corrupt()
			ec = newErrCorrupt(e)
		}
	}()
	runParallel(db.GetState(), checkTable)
	return nil // may be overridden by defer/recover
}

// Check is called by the builtin Database.Check()
func (db *Database) Check() (ec error) {
	defer func() {
		if e := recover(); e != nil {
			db.Corrupt()
			ec = newErrCorrupt(e)
		}
	}()
	state := db.Persist()
	runParallel(state, checkTable)

	if state.Off != 0 {
		state2 := ReadState(db.Store, state.Off)
		assert.This(state.Meta.CksumData()).Is(state2.Meta.CksumData())
	}

	return nil // may be overridden by defer/recover
}

func (db *Database) MustCheck() {
	if err := db.Check(); err != nil {
		panic(err)
	}
}

func checkTable(state *DbState, table string) {
	info := state.Meta.GetRoInfo(table)
	sc := state.Meta.GetRoSchema(table)
	if info == nil {
		panic("info missing for " + table)
	}
	count, size, sum := checkFirstIndex(state, table, sc.Indexes[0].Columns,
		info.Indexes[0])
	if count != info.Nrows {
		panic(fmt.Sprint(table, " ", sc.Indexes[0].Columns,
			" count ", count, " should equal info ", info.Nrows))
	}
	if size != info.Size {
		panic(fmt.Sprint(table, " size ", size, " should equal info ", info.Size))
	}
	for i := 1; i < len(info.Indexes); i++ {
		CheckOtherIndex(table, sc.Indexes[i].Columns, info.Indexes[i], count, sum)
	}
}

func checkFirstIndex(state *DbState, table string, ixcols []string,
	ix *index.Overlay) (int, uint64, uint64) {
	defer func() {
		if e := recover(); e != nil {
			panic(fmt.Sprintln(table, ixcols, e))
		}
	}()
	sum := uint64(0)
	size := uint64(0)
	ix.CheckMerged()
	count := ix.Check(func(off uint64) {
		sum += off // addition so order doesn't matter
		buf := state.store.Data(off)
		n := core.RecLen(buf)
		cksum.MustCheck(buf[:n+cksum.Len])
		size += uint64(n)
	})
	return count, size, sum
}

func CheckOtherIndex(table string, ixcols []string, ix *index.Overlay, nrows int, sumPrev uint64) {
	defer func() {
		if e := recover(); e != nil {
			panic(fmt.Sprintln(table, ixcols, e))
		}
	}()
	ix.CheckMerged()
	sum := uint64(0)
	count := ix.Check(func(off uint64) {
		sum += off // addition so order doesn't matter
	})
	if count != nrows {
		panic(fmt.Sprint("count ", count, " should equal info ", nrows))
	}
	if sum != sumPrev {
		panic("checksum mismatch")
	}
}

//-------------------------------------------------------------------

type ErrCorrupt struct {
	err   any
	table string
}

func (ec *ErrCorrupt) Error() string {
	s := "database corrupt"
	if ec.err != nil {
		s += ": " + fmt.Sprint(ec.err)
	}
	return s
}

func (ec *ErrCorrupt) Table() string {
	if ec == nil { // used by repair
		return ""
	}
	return ec.table
}

func newErrCorrupt(e any) *ErrCorrupt {
	if e == nil {
		return nil
	}
	if e2, ok := e.(*ErrCorrupt); ok {
		return e2
	}
	if e2, ok := e.(error); ok {
		return &ErrCorrupt{err: e2}
	}
	return &ErrCorrupt{err: errors.New(fmt.Sprint(e))}
}

//-------------------------------------------------------------------

func runParallel(state *DbState, fn func(*DbState, string)) {
	tcs := newTableCheckers(state, fn)
	defer tcs.finish()
	tcs.state.Meta.ForEachSchema(func(ts *meta.Schema) {
		select {
		case tcs.work <- ts.Table:
		case <-tcs.stop:
			panic("") // overridden by finish
		}
	})
}

func newTableCheckers(state *DbState, fn func(*DbState, string)) *tableCheckers {
	tcs := tableCheckers{
		state: state,
		fn:    fn,
		work:  make(chan string, 1), // ???
		stop:  make(chan void),
	}
	nw := options.Nworkers
	tcs.wg.Add(nw)
	for i := 0; i < nw; i++ {
		go tcs.worker()
	}
	return &tcs
}

type tableCheckers struct {
	err    atomic.Pointer[ErrCorrupt]
	fn     func(*DbState, string)
	state  *DbState
	work   chan string
	stop   chan void
	wg     sync.WaitGroup
	once   sync.Once
	closed bool
}

func (tcs *tableCheckers) worker() {
	var table string
	defer func() {
		if e := recover(); e != nil {
			tcs.err.Store(&ErrCorrupt{err: e, table: table})
			tcs.once.Do(func() { close(tcs.stop) }) // notify main thread
		}
		tcs.wg.Done()
	}()
	for table = range tcs.work {
		tcs.fn(tcs.state, table)
	}
}

func (tcs *tableCheckers) finish() {
	if !tcs.closed {
		close(tcs.work)
		tcs.closed = true
	}
	tcs.wg.Wait()
	if err := tcs.err.Load(); err != nil {
		panic(err)
	}
}
