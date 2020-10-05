// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/apmckinlay/gsuneido/db19/btree"
	"github.com/apmckinlay/gsuneido/db19/meta"
	"github.com/apmckinlay/gsuneido/db19/stor"
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
	t := time.Now()
	runParallel(db.GetState(), quickCheckTable)
	fmt.Println("quick checked in", time.Since(t).Round(time.Millisecond))
	return nil
}

func quickCheckTable(state *DbState, ts *meta.Schema) {
	info := state.meta.GetRoInfo(ts.Table)
	for _, ix := range info.Indexes {
		ix.QuickCheck()
	}
}

// full check -------------------------------------------------------

// CheckDatabase checks the integrity of the database.
func CheckDatabase(dbfile string) (ec error) {
	db, err := openDatabase(dbfile, stor.READ, false)
	if err != nil {
		return newErrCorrupt(err)
	}
	defer db.Close()
	return db.Check()
}

func(db *Database) Check() (ec error) {
	defer func() {
		if e := recover(); e != nil {
			ec = newErrCorrupt(e)
		}
	}()

	runParallel(db.GetState(), checkTable)
	return nil // may be overridden by defer/recover
}

func checkTable(state *DbState, sc *meta.Schema) {
	info := state.meta.GetRoInfo(sc.Table)
	count, sum := checkFirstIndex(state, info.Indexes[0])
	if count != info.Nrows {
		panic("count != nrows")
	}
	for i := 1; i < len(info.Indexes); i++ {
		ix := info.Indexes[i]
		count, sum = checkOtherIndex(ix, count, sum)
	}
}

func checkFirstIndex(state *DbState, ix *btree.Overlay) (int, uint64) {
	sum := uint64(0)
	count := ix.Check(func(off uint64) {
		sum += off // addition so order doesn't matter
		buf := state.store.Data(off)
		size := runtime.RecLen(buf)
		cksum.MustCheck(buf[:size+cksum.Len])
	})
	return count, sum
}

func checkOtherIndex(ix *btree.Overlay,
	countPrev int, sumPrev uint64) (int, uint64) {
	sum := uint64(0)
	count := ix.Check(func(off uint64) {
		sum += off // addition so order doesn't matter
	})
	if count != countPrev || sum != sumPrev {
		panic("count/sum mismatch")
	}
	return count, sum
}

//-------------------------------------------------------------------

type ErrCorrupt struct {
	err   interface{}
	table string
}

func (ec *ErrCorrupt) Error() string {
	if ec.err == nil {
		return "database corrupt"
	}
	return fmt.Sprint("database corrupt: ", ec.err)
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

func runParallel(state *DbState, fn func(*DbState, *meta.Schema)) {
	tcs := newTableCheckers(state, fn)
	defer tcs.finish()
	tcs.state.meta.ForEachSchema(func(ts *meta.Schema) {
		select {
		case tcs.work <- ts:
		case <-tcs.stop:
			panic("") // overridden by finish
		}
	})
}

func newTableCheckers(state *DbState, fn func(*DbState, *meta.Schema)) *tableCheckers {
	tcs := tableCheckers{
		state: state,
		fn:    fn,
		work:  make(chan *meta.Schema, 1), // ???
		stop:  make(chan void),
	}
	nw := nworkers()
	tcs.wg.Add(nw)
	for i := 0; i < nw; i++ {
		go tcs.worker()
	}
	return &tcs
}

type tableCheckers struct {
	fn     func(*DbState, *meta.Schema)
	wg     sync.WaitGroup
	state  *DbState
	work   chan *meta.Schema
	stop   chan void
	err    atomic.Value
	closed bool
}

func (tcs *tableCheckers) worker() {
	var table string
	defer func() {
		if e := recover(); e != nil {
			tcs.err.Store(&ErrCorrupt{err: e, table: table})
			close(tcs.stop) // notify main thread
		}
		tcs.wg.Done()
	}()
	for ts := range tcs.work {
		table = ts.Table
		tcs.fn(tcs.state, ts)
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
