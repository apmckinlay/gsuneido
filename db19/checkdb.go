// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"errors"
	"fmt"
	"slices"
	"sync"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/index"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/meta/schema"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/options"
	"github.com/apmckinlay/gsuneido/util/cksum"
	"github.com/apmckinlay/gsuneido/util/hacks"
	"github.com/apmckinlay/gsuneido/util/str"
)

// quick check ------------------------------------------------------

// QuickCheck is the default partial checking done at start up.
func (db *Database) QuickCheck() *errCorrupt {
	return checkState(db.GetState(), quickCheckTable, "", nil)
}

func quickCheckTable(tcs *tableCheckers, table string) {
	info := tcs.state.Meta.GetRoInfo(table)
	for _, ix := range info.Indexes {
		ix.QuickCheck()
	}
}

// full check -------------------------------------------------------

// CheckDatabase is called by -check and -repair
func CheckDatabase(dbfile string) error {
	db, err := OpenDb(dbfile, stor.Read, false)
	if err != nil {
		return errCorruptWrap(err)
	}
	defer db.Close()
	if ec := checkState(db.GetState(), checkTable, "", nil); ec != nil {
		db.Corrupt()
		return ec
	}
	return nil
}

// Check is called by the builtin Database.Check()
func (db *Database) Check() (ec error) {
	state := db.Persist()
	if ec := checkState(state, checkTable, "", nil); ec != nil {
		db.Corrupt()
	}
	return ec
}

func (db *Database) MustCheck() {
	if err := db.Check(); err != nil {
		panic(err)
	}
}

func checkTable(tcs *tableCheckers, table string) {
	defer func() {
		if e := recover(); e != nil {
			if ec, ok := e.(*errCorrupt); ok {
				ec.table = table
				panic(ec)
			}
			panic(&errCorrupt{err: e, table: table})
		}
	}()
	info := tcs.state.Meta.GetRoInfo(table)
	sc := tcs.state.Meta.GetRoSchema(table)
	if info == nil {
		panic("info missing for " + table)
	}
	sc.Check(tcs.state.Meta.GetRoSchema)

	ifirst := 0
	if table == tcs.firstTable && tcs.firstIndex != nil {
		for i := range sc.Indexes {
			if slices.Equal(sc.Indexes[i].Columns, tcs.firstIndex) {
				ifirst = i
			}
		}
	}
	ix := &sc.Indexes[ifirst]
	nrows, size, sum := checkFirstIndex(tcs.state.store, ix, info.Indexes[ifirst])
	if nrows != info.Nrows {
		panic(&errCorrupt{ixcols: ix.Columns,
			err: fmt.Sprint("count ", nrows, " should equal info ", info.Nrows)})
	}
	if size != info.Size {
		panic(&errCorrupt{ixcols: ix.Columns,
			err: fmt.Sprint("size ", size, " should equal info ", info.Size)})
	}
	for i := range sc.Indexes {
		if i == ifirst {
			continue
		}
		if tcs.err.Load() != nil {
			break
		}
		CheckOtherIndex(tcs.state.store, &sc.Indexes[i], info.Indexes[i], nrows, sum)
	}
}

func checkFirstIndex(st *stor.Stor, ix *schema.Index, ov *index.Overlay) (int, int64, uint64) {
	defer func() {
		if e := recover(); e != nil {
			panic(&errCorrupt{err: e, ixcols: ix.Columns})
		}
	}()
	sum := uint64(0)
	size := int64(0)
	var buf []byte
	var n int
	ov.CheckMerged()
	base := func(off uint64) {
		sum += off // addition so order doesn't matter
		buf = st.Data(off)
		n = core.RecLen(buf)
		cksum.MustCheck(buf[:n+cksum.Len])
		size += int64(n)
	}
	var ck any = base
	if options.FullCheck {
		ck = func(key string, off uint64) {
			base(off)
			rec := core.Record(hacks.BStoS(buf[:n]))
			checkRecord(rec)
			checkKey(ix.Ixspec, key, rec)
		}
	}
	nrows := ov.CheckBtree(ck)
	return nrows, size, sum
}

func checkRecord(r core.Record) {
	for i := 0; i < r.Count(); i++ {
		r.GetVal(i)
	}
}

func checkKey(ixspec ixkey.Spec, key string, rec core.Record) {
	// theoretically, we could check this field by field
	// without building (allocating) the record key
	key2 := ixspec.Key(rec)
	if key != key2 {
		panic("key mismatch")
	}
}

func CheckOtherIndex(st *stor.Stor, ix *schema.Index, ov *index.Overlay, nrows int, sumPrev uint64) {
	defer func() {
		if e := recover(); e != nil {
			panic(&errCorrupt{err: e, ixcols: ix.Columns})
		}
	}()
	ov.CheckMerged()
	sum := uint64(0)
	ov.CheckMerged()
	var ck any = func(off uint64) {
		sum += off // addition so order doesn't matter
	}
	if options.FullCheck {
		ck = func(key string, off uint64) {
			sum += off
			buf := st.Data(off)
			n := core.RecLen(buf)
			rec := core.Record(hacks.BStoS(buf[:n]))
			checkKey(ix.Ixspec, key, rec)
		}
	}
	nr := ov.CheckBtree(ck)
	if nr != nrows {
		panic(fmt.Sprint("count ", nr, " should equal info ", nrows))
	}
	if sum != sumPrev {
		panic("checksum mismatch")
	}
}

//-------------------------------------------------------------------

// errCorrupt is used to record the table and index
// so that repair can check these first (to speed up the checking)
type errCorrupt struct {
	err    any
	table  string
	ixcols []string
}

func (ec *errCorrupt) Error() string {
	s := "database corrupt"
	if ec.table != "" {
		s += " " + ec.table
	}
	if ec.ixcols != nil {
		s += " " + str.Join("[,]", ec.ixcols)
	}
	if ec.err != nil {
		s += ": " + fmt.Sprint(ec.err)
	}
	return s
}

func (ec *errCorrupt) Table() string {
	if ec == nil {
		return ""
	}
	return ec.table
}

func (ec *errCorrupt) Ixcols() []string {
	if ec == nil {
		return nil
	}
	return ec.ixcols
}

func errCorruptWrap(e any) *errCorrupt {
	if e == nil {
		return nil
	}
	if e2, ok := e.(*errCorrupt); ok {
		return e2
	}
	if e2, ok := e.(error); ok {
		return &errCorrupt{err: e2}
	}
	return &errCorrupt{err: errors.New(fmt.Sprint(e))}
}

//-------------------------------------------------------------------

// checkState runs fn for all tables in state using a worker pool.
// fn is either quickCheckTable (for startup) or checkTable (for check & repair).
func checkState(state *DbState, fn func(*tableCheckers, string),
	firstTable string, firstIndex []string) (ec *errCorrupt) {
	tcs := newTableCheckers(state, fn)
	tcs.firstTable = firstTable
	tcs.firstIndex = firstIndex
	tcs.sendWork()
	return tcs.finish()
}

type tableCheckers struct {
	err        atomic.Pointer[errCorrupt]
	fn         func(*tableCheckers, string)
	state      *DbState
	work       chan string
	stop       chan void
	wg         sync.WaitGroup
	firstTable string
	firstIndex []string
}

func newTableCheckers(state *DbState, fn func(*tableCheckers, string)) *tableCheckers {
	tcs := tableCheckers{
		state: state,
		fn:    fn,
		work:  make(chan string), // no buffer to stop quickly
		stop:  make(chan void),
	}
	nw := options.Nworkers // more doesn't seem to help
	for range nw {
		tcs.wg.Go(tcs.worker)
	}
	return &tcs
}

func (tcs *tableCheckers) sendWork() {
	defer func() {
		if e := recover(); e != nil {
			tcs.error(e, "")
		}
	}()
	if tcs.firstTable != "" {
		tcs.work <- tcs.firstTable
	}
	for ts := range tcs.state.Meta.Tables() {
		if ts.Table == tcs.firstTable {
			continue
		}
		select {
		case tcs.work <- ts.Table:
		case <-tcs.stop: // stop is closed if a worker gets an error
			return
		}
	}
}

func (tcs *tableCheckers) worker() {
	var table string
	defer func() {
		if e := recover(); e != nil {
			tcs.error(e, table)
		}
	}()
	for table = range tcs.work {
		tcs.fn(tcs, table)
	}
}

func (tcs *tableCheckers) error(e any, table string) {
	ec, _ := e.(*errCorrupt)
	if ec == nil {
		ec = &errCorrupt{err: e, table: table}
	} else {
		ec.table = table
	}
	if tcs.err.CompareAndSwap(nil, ec) { // save first error
		close(tcs.stop) // notify main thread, once only
	}
}

func (tcs *tableCheckers) finish() *errCorrupt {
	close(tcs.work)
	// Theoretically, we don't need to wait if we have an error
	// but then you can get errors when you close the store.
	// Could defer the waits to the end but that's difficult to arrange.
	// if ec := tcs.err.Load(); ec != nil {
	// 	return ec // if error then don't need to wait
	// }
	tcs.wg.Wait()
	return tcs.err.Load()
}
