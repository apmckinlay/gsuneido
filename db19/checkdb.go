// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/cksum"
)

type dbcheck DbState

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

// CheckDatabase checks the integrity of the database.
func CheckDatabase(dbfile string) (ec error) {
	db, err := OpenDb(dbfile, stor.READ, false)
	if err != nil {
		return newErrCorrupt(err)
	}
	defer db.Close()
	return db.Check()
}

// Check is called by the builtin Database.Check()
func (db *Database) Check() (ec error) {
	defer func() {
		if e := recover(); e != nil {
			ec = newErrCorrupt(e)
		}
	}()
	state := db.Persist()
	runParallel(state, checkTable)
	return nil // may be overridden by defer/recover
}

func checkTable(state *DbState, table string) {
	info := state.Meta.GetRoInfo(table)
	if info == nil {
		panic("info missing for " + table)
	}
	count, sum := checkFirstIndex(state, info.Indexes[0])
	if count != info.Nrows {
		panic("count != nrows " + fmt.Sprint(count, info.Nrows))
	}
	for _, ix := range info.Indexes[1:] {
		CheckOtherIndex(ix, count, sum)
	}
}

func checkFirstIndex(state *DbState, ix *index.Overlay) (int, uint64) {
	sum := uint64(0)
	ix.CheckMerged()
	count := ix.Check(func(off uint64) {
		sum += off // addition so order doesn't matter
		buf := state.store.Data(off)
		size := runtime.RecLen(buf)
		cksum.MustCheck(buf[:size+cksum.Len])
	})
	return count, sum
}

func CheckOtherIndex(ix *index.Overlay, countPrev int, sumPrev uint64) {
	ix.CheckMerged()
	sum := uint64(0)
	count := ix.Check(func(off uint64) {
		sum += off // addition so order doesn't matter
	})
	if count != countPrev {
		panic("count mismatch " + fmt.Sprint(countPrev, count))
	}
	if sum != sumPrev {
		panic("checksum mismatch")
	}
}

//-------------------------------------------------------------------

type ErrCorrupt struct {
	err   interface{}
	table string
}

func (ec *ErrCorrupt) Error() string {
	s := "database corrupt"
	if ec.table != "" {
		s += ": " + ec.table
	}
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

func newErrCorrupt(e interface{}) *ErrCorrupt {
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
	err    atomic.Value
	fn     func(*DbState, string)
	state  *DbState
	work   chan string
	stop   chan void
	once   sync.Once
	wg     sync.WaitGroup
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
