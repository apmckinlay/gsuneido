// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"math/rand"

	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/ordset"
	"github.com/apmckinlay/gsuneido/util/ranges"
	"github.com/apmckinlay/gsuneido/util/verify"
)

//TODO abort long transactions

// Need to use an ordered set so that reads can check for a range
type Set = ordset.Set

type Ranges = ranges.Ranges

// Check holds the data for the transaction conflict checker.
// Checking is designed to be single threaded i.e. run in its own goroutine.
// It is intended to run asynchronously, i.e. callers not waiting for results.
// This allow more concurrency (overlap) with user code.
// Actions are checked as they are done, incrementally
// A conflict with a completed transaction aborts the current transaction.
// A conflict with an outstanding (not completed) transaction
// randomly aborts one of the two transactions.
// The checker serializes transaction commits.
// A single sequence counter is used to assign unique start and end values.
type Check struct {
	seq    int
	oldest int
	trans  map[int]*cktran
}

type cktran struct {
	start  int
	end    int
	tables map[string]*cktbl
}

type cktbl struct {
	// writes tracks outputs, updates, and deletes
	writes ckwrites
	reads  ckreads
}

type ckwrites []*Set
type ckreads []*Ranges

func NewCheck() *Check {
	return &Check{trans: make(map[int]*cktran), oldest: ints.MaxInt}
}

func (ck *Check) StartTran() int {
	start := ck.next()
	ck.trans[start] = &cktran{start: start, end: ints.MaxInt,
		tables: make(map[string]*cktbl)}
	return start
}

func (ck *Check) next() int {
	ck.seq++
	return ck.seq
}

// Read adds a read action.
// Will conflict if another transaction has a write within the range.
func (ck *Check) Read(tn int, table string, index int, from, to string) bool {
	trace("T", tn, "read", table, "index", index, "from", from, "to", to)
	t, ok := ck.trans[tn]
	if !ok {
		return false // it's gone, presumably aborted
	}
	verify.That(!t.isEnded())
	// check against overlapping transactions
	for _, t2 := range ck.trans {
		if t2 != t && overlap(t, t2) {
			if tbl, ok := t2.tables[table]; ok {
				if tbl.writes.anyInRange(index, from, to) {
					if ck.abort1of(t, t2) {
						return false // this transaction got aborted
					}
				}
			}
		}
	}
	t.saveRead(table, index, from, to)
	return true
}

func (t *cktran) saveRead(table string, index int, from,to string) {
	tbl, ok := t.tables[table]
	if !ok {
		tbl = &cktbl{}
		t.tables[table] = tbl
	}
		tbl.reads = tbl.reads.with(index, from, to)
}

func (cr ckreads) with(index int, from,to string) ckreads {
	for len(cr) <= index {
		cr = append(cr, nil)
	}
	if cr[index] == nil {
		cr[index] = &Ranges{}
	}
	cr[index].Insert(from, to)
	return cr
}

func (cr ckreads) contains(index int, key string) bool {
	return index < len(cr) && cr[index].Contains(key)
}

// Write adds output/update/delete actions.
// Updates require two calls, one with the from keys, another with the to keys.
// NOTE: Even if an update doesn't change a key, it still has to register it.
func (ck *Check) Write(tn int, table string, keys []string) bool {
	trace("T", tn, "write", table, "keys", keys)
	t, ok := ck.trans[tn]
	if !ok {
		return false // it's gone, presumably aborted
	}
	verify.That(!t.isEnded())
	// check against overlapping transactions
	for _, t2 := range ck.trans {
		if t2 != t && overlap(t, t2) {
			if tbl, ok := t2.tables[table]; ok {
				for i, key := range keys {
					// check against other writes
					if key != "" &&
						(tbl.writes.contains(i, key) ||
						tbl.reads.contains(i, key)) {
						if ck.abort1of(t, t2) {
							return false // this transaction got aborted
						}
					}
				}
			}
		}
	}
	t.saveWrite(table, keys)
	return true
}

func (t *cktran) saveWrite(table string, keys []string) {
	tbl, ok := t.tables[table]
	if !ok {
		tbl = &cktbl{}
		t.tables[table] = tbl
	}
	for i, key := range keys {
		tbl.writes = tbl.writes.with(i, key)
	}
}

func (cw ckwrites) contains(index int, key string) bool {
	return index < len(cw) && cw[index].Contains(key)
}

func (cw ckwrites) anyInRange(index int, from, to string) bool {
	if index >= len(cw) {
		return false
	}
	return cw[index].AnyInRange(from, to)
}

func (cw ckwrites) with(index int, key string) ckwrites {
	for len(cw) <= index {
		cw = append(cw, nil)
	}
	if cw[index] == nil {
		cw[index] = &Set{}
	}
	cw[index].Insert(key)
	return cw
}

// checkerAbortT1 is used by tests to avoid randomness
var checkerAbortT1 = false

// abort1of aborts one of t1 and t2.
// If t2 is committed, abort t1, otherwise choose randomly.
// It returns true if t1 is aborted, false if t2 is aborted.
func (ck *Check) abort1of(t1, t2 *cktran) bool {
	trace("conflict with", t2)
	if t2.isEnded() || checkerAbortT1 || rand.Intn(2) == 1 {
		ck.Abort(t1.start)
		return true
	}
	ck.Abort(t2.start)
	return false
}

func (t *cktran) isEnded() bool {
	return t.end != ints.MaxInt
}

// Abort cancels a transaction.
// It returns false if the transaction is not found (e.g. already aborted).
func (ck *Check) Abort(tn int) bool {
	trace("abort", tn)
	if _, ok := ck.trans[tn]; !ok {
		return false
	}
	delete(ck.trans, tn)
	if tn == ck.oldest {
		ck.oldest = ints.MaxInt // need to find the new oldest
	}
	ck.cleanEnded()
	return true
}

// Commit finishes a transaction.
// It returns false if the transaction is not found (e.g. already aborted).
// No additional checking required since actions have already been checked.
func (ck *Check) Commit(tn int) bool {
	trace("commit", tn)
	t, ok := ck.trans[tn]
	if !ok {
		return false // it's gone, presumably aborted
	}
	t.end = ck.next()
	if t.start == ck.oldest {
		ck.oldest = ints.MaxInt // need to find the new oldest
	}
	ck.cleanEnded()
	return true
}

func overlap(t1, t2 *cktran) bool {
	return t1.end > t2.start && t2.end > t1.start
}

// cleanEnded removes ended transactions
// that finished before the earliest outstanding start time.
func (ck *Check) cleanEnded() {
	// find oldest start of non-ended (would be faster with a heap)
	if ck.oldest == ints.MaxInt {
		for _, t := range ck.trans {
			if t.end == ints.MaxInt && t.start < ck.oldest {
				ck.oldest = t.start
			}
		}
		trace("OLDEST", ck.oldest)
	}
	// remove any ended transactions older than this
	for tn, t := range ck.trans {
		if t.end != ints.MaxInt && t.end < ck.oldest {
			trace("REMOVE", tn, "->", t.end)
			delete(ck.trans, tn)
		}
	}
}

func trace(args ...interface{}) {
	// fmt.Println(args...) // comment out to disable tracing
}
