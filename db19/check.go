// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"log"
	"math/rand"
	"strconv"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/ordset"
	"github.com/apmckinlay/gsuneido/util/ranges"
)

const maxTrans = 200

// Need to use an ordered set so that reads can check for a range
type Set = ordset.Set

type Ranges = ranges.Ranges

// Check holds the data for the transaction conflict checker.
// Checking is designed to be single threaded i.e. run in its own goroutine.
// It is intended to run asynchronously, i.e. callers not waiting for results.
// This allow more concurrency (overlap) with user code.
// Actions are checked as they are done, incrementally.
// A conflict with a completed transaction aborts the current transaction.
// A conflict with an outstanding (not completed) transaction
// randomly aborts one of the two transactions.
// The checker serializes transaction commits.
// A single sequence counter is used to assign unique start and end values.
// See CheckCo for the concurrent channel based interface to Check.
// See Checker for the common interface to Check and CheckCo
type Check struct {
	db     *Database
	seq    int
	oldest int
	// clock is used to abort long transactions
	clock int
	// trans hold the outstanding/overlapping update transactions
	trans map[int]*CkTran
}

type CkTran struct {
	start    int
	end      int
	birth    int
	tables   map[string]*cktbl
	state    *DbState
	conflict atomic.Value // string
}

type cktbl struct {
	// writes tracks outputs, updates, and deletes
	writes ckwrites
	reads  ckreads
}

// ckwrites and ckreads have one element per index
// in parallel with the schema and info indexes
type (
	ckwrites []*Set
	ckreads  []*Ranges
)

func (t *CkTran) String() string {
	return "ut" + strconv.Itoa(t.start)
}

func NewCheck(db *Database) *Check {
	return &Check{db: db, trans: make(map[int]*CkTran), oldest: ints.MaxInt}
}

func (ck *Check) StartTran() *CkTran {
	if len(ck.trans) >= maxTrans {
		return nil
	}
	start := ck.next()
	var state *DbState
	if ck.db != nil {
		state = ck.db.GetState()
	}
	t := &CkTran{start: start, end: ints.MaxInt, birth: ck.clock,
		tables: make(map[string]*cktbl), state: state}
	ck.trans[start] = t
	return t
}

func (ck *Check) next() int {
	ck.seq++
	return ck.seq
}

// Read adds a read action.
// Will conflict if another transaction has a write within the range.
func (ck *Check) Read(t *CkTran, table string, index int, from, to string) bool {
	trace("T", t.start, "read", table, "index", index, "from", from, "to", to)
	t, ok := ck.trans[t.start]
	if !ok {
		return false // it's gone, presumably aborted
	}
	assert.That(!t.isEnded())
	// check against overlapping transactions
	for _, t2 := range ck.trans {
		if t2 != t && overlap(t, t2) {
			if tbl, ok := t2.tables[table]; ok {
				if tbl.writes.anyInRange(index, from, to) {
					if ck.abort1of(t, t2, "read", "write") {
						return false // this transaction got aborted
					}
				}
			}
		}
	}
	t.saveRead(table, index, from, to)
	return true
}

func (t *CkTran) saveRead(table string, index int, from, to string) {
	tbl, ok := t.tables[table]
	if !ok {
		tbl = &cktbl{}
		t.tables[table] = tbl
	}
	tbl.reads = tbl.reads.with(index, from, to)
}

func (cr ckreads) with(index int, from, to string) ckreads {
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
// Will conflict with another write to the same index/key or an overlapping read.
// The keys are parallel with the indexes i.e. keys[i] is for indexes[i].
// Updates require two calls, one with the old keys, another with the new keys.
// NOTE: Even if an update doesn't change a key, it still has to register it.
func (ck *Check) Write(t *CkTran, table string, keys []string) bool {
	trace("T", t.start, "write", table, "keys", keys)
	t, ok := ck.trans[t.start]
	if !ok {
		return false // it's gone, presumably aborted
	}
	assert.That(!t.isEnded())
	// check against overlapping transactions
	for _, t2 := range ck.trans {
		if t2 != t && overlap(t, t2) {
			if tbl, ok := t2.tables[table]; ok {
				for i, key := range keys {
					if key != "" {
						act2 := ""
						if tbl.writes.contains(i, key) {
							act2 = "write"
						} else if tbl.reads.contains(i, key) {
							act2 = "read"
						} else {
							continue
						}
						if ck.abort1of(t, t2, "write", act2) {
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

func (t *CkTran) saveWrite(table string, keys []string) {
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
func (ck *Check) abort1of(t1, t2 *CkTran, act1, act2 string) bool {
	trace("conflict with", t2)
	if t2.isEnded() || checkerAbortT1 || rand.Intn(2) == 1 {
		ck.abort(t1.start, act1+" in this transaction conflicted with "+
			act2+" in another transaction")
		return true
	}
	ck.abort(t2.start, act2+" in this transaction conflicted with "+
		act1+" in another transaction")
	return false
}

func (t *CkTran) isEnded() bool {
	return t.end != ints.MaxInt
}

// Abort cancels a transaction.
// It returns false if the transaction is not found (e.g. already aborted).
func (ck *Check) Abort(t *CkTran, reason string) bool {
	return ck.abort(t.start, reason)
}

func (ck *Check) abort(tn int, reason string) bool {
	trace("abort", tn)
	t, ok := ck.trans[tn]
	if !ok {
		return false
	}
	t.conflict.Store(reason)
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
func (ck *Check) Commit(ut *UpdateTran) bool {
	return ck.commit(ut) != nil
}

func (ck *Check) commit(ut *UpdateTran) []string {
	tn := ut.num()
	trace("commit", tn)
	t, ok := ck.trans[tn]
	if !ok {
		return nil // it's gone, presumably aborted
	}
	t.end = ck.next()
	if t.start == ck.oldest {
		ck.oldest = ints.MaxInt // need to find the new oldest
	}
	ck.cleanEnded()
	return t.tablesWritten()
}

func (t *CkTran) tablesWritten() []string {
	tw := make([]string, 0, 8)
	for table, tbl := range t.tables {
		if tbl.writes != nil && !tbl.writes[0].Empty() {
			tw = append(tw, table)
		}
	}
	return tw
}

func overlap(t1, t2 *CkTran) bool {
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

// MaxAge is the maximum number of ticks that a transaction can be outstanding.
// Transactions are aborted if they exceed this limit.
var MaxAge = 20

// tick should be called regularly e.g. once per second
// to abort transactions older than MaxAge.
func (ck *Check) tick() {
	ck.clock++
	trace("tick", ck.clock)
	for tn, t := range ck.trans {
		if ck.clock-t.birth >= MaxAge {
			trace("abort", tn, "age", ck.clock-t.birth)
			log.Println("aborted", t, "update transaction longer than", MaxAge, "seconds")
			ck.abort(tn, "transaction exceeded max age")
		}
	}
}

func (ck *Check) Stop() { // to satisfy Checker interface
}

func trace(...interface{}) {
	// fmt.Println(args...) // comment out to disable tracing
}
