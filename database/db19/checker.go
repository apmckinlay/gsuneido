// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/ordset"
)

type Set = ordset.Set

// Check holds the data for the transaction conflict checker.
// Checking is designed to be single threaded i.e. run in its own goroutine.
// It is intended to run asynchronously, i.e. not waiting for results.
// This allow more concurrency (overlap) with user code.
// Actions are checked as they are done, incrementally
// A conflict with a completed transaction aborts the current transaction.
// A conflict with an outstanding (not completed) transaction
// does not immediately abort either transaction.
// The two transactions are marked as mutually conflicting.
// Whichever one commits first will succeed and the other will be aborted.
// Or if one aborts, the other can commit.
// The checker serializes transaction commits.
// A single sequence counter is used to assign unique start and end values.
type Check struct {
	seq    int
	oldest int
	trans  map[int]*cktran
}

type cktran struct {
	start     int
	end       int
	tables    map[string]*cktbl
	conflicts []*cktran
}

type cktbl struct {
	// writes tracks outputs, updates, and deletes
	writes ckwrites
	//TODO reads
}

type ckwrites []*Set

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

// Write adds output/update/delete actions.
// Updates require two calls, one with the from keys, another with the to keys.
func (ck *Check) Write(tn int, table string, keys []string) bool {
	trace("T", tn, "output", table, "keys", keys)
	// check overlapping transactions
	t, ok := ck.trans[tn]
	if !ok {
		return false // it's gone, presumably aborted
	}
	for _, t2 := range ck.trans {
		if overlap(t, t2) {
			if tbl, ok := t2.tables[table]; ok {
				for i, key := range keys {
					if key != "" && tbl.writes.Contains(i, key) {
						if !ck.addConflict(t, t2) {
							return false // other tran is already committed
						}
					}
				}
			}
		}
	}
	// save action in this transaction
	tbl, ok := t.tables[table]
	if !ok {
		tbl = &cktbl{}
		t.tables[table] = tbl
	}
	for i, key := range keys {
		tbl.writes = tbl.writes.With(i, key)
	}
	return true
}

func (o ckwrites) Contains(index int, key string) bool {
	return index < len(o) && o[index].Contains(key)
}

func (o ckwrites) With(index int, key string) ckwrites {
	for len(o) <= index {
		o = append(o, nil)
	}
	if o[index] == nil {
		o[index] = &Set{}
	}
	o[index].Insert(key)
	return o
}

// addConflict makes t1 and t2 mutually conflicting
func (ck *Check) addConflict(t1, t2 *cktran) bool {
	trace("conflict with", t2)
	if t2.isEnded() {
		ck.Abort(t1.start)
		return false
	}
	t1.conflicts = append(t1.conflicts, t2)
	t2.conflicts = append(t2.conflicts, t1)
	return true
}

func (t *cktran) isEnded() bool {
	return t.end != ints.MaxInt
}

// delConflict deletes t2 from t1's conflicts
func delConflict(t1, t2 *cktran) {
	for i, t := range t1.conflicts {
		if t == t2 && t.end == ints.MaxInt {
			last := len(t1.conflicts) - 1
			t1.conflicts[i] = t1.conflicts[last]
			t1.conflicts = t1.conflicts[:last]
		}
	}
}

// Abort cancels a transaction.
// It returns false if the transaction is not found (e.g. already aborted).
// Conflicts with this transaction are removed from the conflicting transactions.
func (ck *Check) Abort(tn int) bool {
	trace("abort", tn)
	t, ok := ck.trans[tn]
	if !ok {
		return false
	}
	delete(ck.trans, tn)
	if tn == ck.oldest {
		ck.oldest = ints.MaxInt // need to find the new oldest
	}
	ck.cleanEnded()
	// remove conflicts
	for _, t2 := range t.conflicts {
		delConflict(t2, t)
	}
	return true
}

// Commit finishes a transaction.
// It returns false if the transaction is not found (e.g. already aborted).
// Conflicting transactions are aborted.
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
	// abort conflicting
	for _, t2 := range t.conflicts {
		ck.Abort(t2.start)
	}
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
