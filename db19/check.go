// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"log"
	"math"
	"math/rand/v2"
	"strconv"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/atomics"
	"github.com/apmckinlay/gsuneido/util/ordset"
	"github.com/apmckinlay/gsuneido/util/ranges"
)

/*
Reads are tracked as Ranges on specific indexes (ckreads).
Output and Deletes are tracked as Set's of key values for each index (ckwrites).
Deletes are also tracked by offset which is used when checking delete vs delete.

|        | committed                 ||| outstanding            |||
|        | read      | output | delete | read   | output | delete |
| ------ | --------- | ------ | ------ | ------ | ------ | ------ |
| read   | no (1)(3) | check  | check  | no (1) | check  | check  |
| output | no (3)    | no (2) | no (2) | check  | no (2) | no (2) |
| delete | no (3)    | no (2) | check  | check  | no (2) | check  |

1.  reads never conflict with reads
2.  outputs never conflict with other outputs or deletes
3.  don't need to check committed reads
*/

// MaxTrans is the maximum number of outstanding/overlapping update transactions
const MaxTrans = 500 // ???

// Set needs to be an ordered set so that reads can check for a range
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
	db *Database
	// actvtran hold the active transactions.
	// It is a map so we can delete from it.
	actvTran map[int]*CkTran
	// cmtdTran hold the committed transactions.
	// It is a map so we can delete from it.
	cmtdTran map[int]*CkTran
	// bytable holds the actions for each table for each transaction.
	// It is accessed by [tablename][tran.start]
	// It is a map so we can delete from it.
	bytable map[string]map[int]*actions
	// exclusive controls access to tables
	exclusive map[string]int
	seq       int
	// oldest is the start time of the oldest active transaction.
	// A value of math.MaxInt means it needs to be recomputed.
	oldest int
	// clock is used to abort long transactions
	clock int
}

type CkTran struct {
	// failure is written by CkTran.abort and read by UpdateTran.
	// It is set to either a conflict or a timeout.
	failure atomics.String
	tables  []string
	state   *DbState
	// start is the start time and also used as an identifier
	start int
	// end is math.MaxInt if the transaction is not ended
	end int
	// birth is used by tick to abort long transactions
	birth        int
	readCount    int
	hasUpdates   bool
	readConflict string
}

// readMax is higher than jSuneido to account for duplicate checks
// writeMax is handle in tran.go
const readMax = 20000

type actions struct {
	tran    *CkTran
	outputs ckwrites
	deletes ckwrites
	deloffs map[uint64]struct{}
	reads   ckreads
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
	return &Check{db: db,
		actvTran:  make(map[int]*CkTran),
		cmtdTran:  make(map[int]*CkTran),
		bytable:   make(map[string]map[int]*actions),
		oldest:    math.MaxInt,
		exclusive: make(map[string]int),
		seq:       1} // odd
}

func (ck *Check) Run(fn func() error) error {
	return fn()
}

var logAt = 100

func (ck *Check) StartTran() *CkTran {
	if ck.count() > logAt {
		log.Println("outstanding transactions reached", logAt)
		logAt += 100
	}
	if ck.count() >= MaxTrans {
		return nil
	}
	start := ck.next()
	var state *DbState
	if ck.db != nil {
		state = ck.db.GetState()
	}
	t := &CkTran{start: start, end: math.MaxInt, birth: ck.clock, state: state}
	ck.actvTran[start] = t
	return t
}

func (ck *Check) count() int {
	return len(ck.actvTran) + len(ck.cmtdTran)
}

func (ck *Check) next() int {
	ck.seq += 2 // odd
	return ck.seq
}

// AddExclusive is used for creating indexes on existing tables
func (ck *Check) AddExclusive(table string) bool {
	if end, ok := ck.exclusive[table]; ok && end == math.MaxInt {
		return false // already exclusive
	}
	if ts, ok := ck.bytable[table]; ok {
		for _, ta := range ts {
			if ta.outputs != nil || ta.deletes != nil {
				ck.abort(ta.tran.start, "preempted by exclusive ("+table+")")
			}
		}
	}
	ck.exclusive[table] = math.MaxInt
	return true
}

func (ck *Check) EndExclusive(table string) {
	if ck.exclusive[table] != math.MaxInt {
		return
	}
	end := ck.next()
	ck.exclusive[table] = end
	// after ending, we still block transactions that started previously
	ck.cleanEnded()
}

// Persist is just for tests, it doesn't actually persist
func (ck *Check) Persist() *DbState {
	return ck.db.GetState()
}

func (ck *Check) RunExclusive(table string, fn func()) any {
	// only for tests
	if !ck.AddExclusive(table) {
		return "already exclusive: " + table
	}
	return ck.RunEndExclusive(table, fn)
}

func (ck *Check) RunEndExclusive(table string, fn func()) (err any) {
	// only for tests
	defer ck.EndExclusive(table)
	defer func() {
		if e := recover(); e != nil {
			err = e
		}
	}()
	fn()
	return nil
}

//-------------------------------------------------------------------

// Read adds a read action.
// Will conflict if another transaction has a write within the range.
func (ck *Check) Read(t *CkTran, table string, index int, from, to string) bool {
	traceln("T", t.start, "read", table, "index", index, "from", from, "to", to)
	if t.readConflict != "" {
		// once we have a conflict, we don't need to check any more reads
		// either it will be read-only (no conflicts) or it will abort
		return true
	}
	if _, ok := ck.actvTran[t.start]; !ok {
		return false // it's gone, presumably aborted
	}
	assert.That(!t.ended())
	// check against overlapping transactions
	if ts, ok := ck.bytable[table]; ok {
		for _, ta := range ts {
			if ta.tran != t && overlap(t, ta.tran) {
				if ta.outputs.anyInRange(index, from, to) ||
					ta.deletes.anyInRange(index, from, to) {
					if ck.abort1of(t, ta.tran, "read", "write", table) {
						return false // this transaction got aborted
					}
				}
			}
		}
	}
	if !ck.saveRead(t, table, index, from, to) {
		ck.abort(t.start, "too many reads in one update transaction")
	}
	return true
}

func (ck *Check) saveRead(t *CkTran, table string, index int, from, to string) bool {
	ta := ck.getActs(t, table)
	reads, inc := ta.reads.with(index, from, to)
	if t.readCount += inc; t.readCount >= readMax {
		return false
	}
	ta.reads = reads
	return true
}

func (ck *Check) getActs(t *CkTran, table string) *actions {
	ts, ok := ck.bytable[table]
	if !ok {
		ts = make(map[int]*actions)
		ck.bytable[table] = ts
	}
	ta, ok := ts[t.start]
	if !ok {
		ta = &actions{tran: t}
		ts[t.start] = ta
		t.tables = append(t.tables, table)
	}
	return ta
}

func (cr ckreads) with(index int, from, to string) (ckreads, int) {
	for len(cr) <= index {
		cr = append(cr, nil)
	}
	if cr[index] == nil {
		cr[index] = &Ranges{}
	}
	inc := cr[index].Insert(from, to)
	return cr, inc
}

func (cr ckreads) contains(index int, key string) bool {
	return index < len(cr) && cr[index].Contains(key)
}

// Update adds a delete of oldoff/oldkeys and an output of newkeys.
func (ck *Check) Update(t *CkTran,
	table string, oldoff uint64, oldkeys, newkeys []string) bool {
	return ck.Delete(t, table, oldoff, oldkeys) &&
		ck.output(t, table, newkeys, oldkeys)
}

// Output adds an output action.
// Outputs only need to be checked against outstanding reads.
// The keys are parallel with the indexes i.e. keys[i] is for indexes[i].
func (ck *Check) Output(t *CkTran, table string, keys []string) bool {
	if !ck.gotUpdate(t) {
		return false
	}
	return ck.output(t, table, keys, nil)
}

// gotUpdate returns false if it aborts the transaction
func (ck *Check) gotUpdate(t *CkTran) bool {
	if !t.hasUpdates {
		if t.readConflict != "" {
			ck.abort(t.start, t.readConflict)
			return false
		}
		t.hasUpdates = true
	}
	return true
}

// output adds an output action.
// oldkeys is used by Update so we can tell if the key changed.
func (ck *Check) output(t *CkTran, table string, keys, oldkeys []string) bool {
	traceln("T", t.start, "output", table, "keys", keys)
	if _, ok := ck.actvTran[t.start]; !ok {
		return false // it's gone, presumably aborted
	}
	assert.That(!t.ended())
	if t.start < ck.exclusive[table] {
		ck.abort(t.start, "conflict with exclusive ("+table+")")
		return false
	}
	// check against overlapping transactions
	if ts, ok := ck.bytable[table]; ok {
		for _, ta := range ts {
			if ta.tran != t && overlap(t, ta.tran) && !ta.tran.ended() {
				for i, key := range keys {
					if (oldkeys == nil || key != oldkeys[i]) &&
						ta.reads.contains(i, key) {
						if ck.abort1of(t, ta.tran, "output", "read", table) {
							return false // this transaction got aborted
						}
					}
				}
			}
		}
	}
	if !ck.saveOutput(t, table, keys) {
		ck.abort(t.start,
			"too many writes (output, update, or delete) in one transaction")
	}
	return true
}

func (ck *Check) saveOutput(t *CkTran, table string, keys []string) bool {
	ta := ck.getActs(t, table)
	for i, key := range keys {
		outs := ta.outputs.with(i, key)
		if outs == nil {
			return false
		}
		ta.outputs = outs
	}
	return true
}

// Delete adds a delete action.
// Outputs only need to be checked against outstanding reads.
// The keys are parallel with the indexes i.e. keys[i] is for indexes[i].
func (ck *Check) Delete(t *CkTran, table string, off uint64, keys []string) bool {
	traceln("T", t.start, "delete", table, "off", off, "keys", keys)
	if !ck.gotUpdate(t) {
		return false
	}
	if _, ok := ck.actvTran[t.start]; !ok {
		return false // it's gone, presumably aborted
	}
	assert.That(!t.ended())
	if t.start < ck.exclusive[table] {
		ck.abort(t.start, "conflict with exclusive ("+table+")")
		return false
	}
	// check against overlapping transactions
	if ts, ok := ck.bytable[table]; ok {
		for _, ta := range ts {
			if ta.tran != t && overlap(t, ta.tran) {
				if _, ok := ta.deloffs[off]; ok {
					if ck.abort1of(t, ta.tran, "delete", "delete", table) {
						return false // this transaction got aborted
					}
				}
				for i, key := range keys {
					if !ta.tran.ended() && ta.reads.contains(i, key) {
						if ck.abort1of(t, ta.tran, "delete", "read", table) {
							return false // this transaction got aborted
						}
					}
				}
			}
		}
	}
	if !ck.saveDelete(t, table, off, keys) {
		ck.abort(t.start,
			"too many writes (output, update, or delete) in one transaction")
	}
	return true
}

func (ck *Check) saveDelete(t *CkTran, table string, off uint64, keys []string) bool {
	ta := ck.getActs(t, table)
	if ta.deloffs == nil {
		ta.deloffs = make(map[uint64]struct{})
	}
	assert.That(off != 0)
	ta.deloffs[off] = struct{}{}
	for i, key := range keys {
		dels := ta.deletes.with(i, key)
		if dels == nil {
			return false
		}
		ta.deletes = dels
	}
	return true
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
	if !cw[index].Insert(key) {
		return nil
	}
	return cw
}

func (ck *Check) ReadCount(t *CkTran) int {
	if _, ok := ck.actvTran[t.start]; !ok {
		return -1 // it's gone, presumably aborted
	}
	return t.readCount
}

// checkerAbortT1 is used by tests to avoid randomness
var checkerAbortT1 = false

// abort1of aborts one of t1 and t2.
// If t2 is committed, abort t1, otherwise choose randomly.
// It returns true if t1 is aborted, false otherwise.
func (ck *Check) abort1of(t1, t2 *CkTran, act1, act2, table string) bool {
	traceln("conflict with", t2)
	if !t1.hasUpdates {
		t1.readConflict = act1 + " in this transaction conflicted with " +
			act2 + " in another transaction (" + table + ")"
		traceln(t1, "readConflict =", t1.readConflict)
		return false // t1 not aborted
	}
	if !t2.hasUpdates {
		t2.readConflict = act2 + " in this transaction conflicted with " +
			act1 + " in another transaction (" + table + ")"
		traceln(t2, "readConflict =", t2.readConflict)
		return false // t1 not aborted
	}
	if t2.ended() || checkerAbortT1 || rand.IntN(2) == 1 {
		ck.abort(t1.start, act1+" in this transaction conflicted with "+
			act2+" in another transaction ("+table+")")
		return true
	}
	ck.abort(t2.start, act2+" in this transaction conflicted with "+
		act1+" in another transaction ("+table+")")
	return false
}

func (t *CkTran) ended() bool {
	return t.end != math.MaxInt
}

func (t *CkTran) Failed() bool {
	return t.failure.Load() != ""
}

// Abort cancels a transaction.
// It returns false if the transaction is not found (e.g. already aborted).
func (ck *Check) Abort(t *CkTran, reason string) bool {
	return ck.abort(t.start, reason)
}

func (ck *Check) abort(tn int, reason string) bool {
	traceln("abort", tn)
	t, ok := ck.actvTran[tn]
	if !ok {
		return false
	}
	if reason == "" {
		reason = "abort"
	}
	t.failure.Store(reason)
	ck.removeByTable(t)
	delete(ck.actvTran, tn)
	if tn == ck.oldest {
		ck.oldest = math.MaxInt // need to find the new oldest
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
	tw := []string{} // not nil
	if ck.db.IsCorrupted() {
		ck.Abort(ut.ct, "database is locked")
		return tw // not nil/failure
	}
	tn := ut.ct.start
	traceln("commit", tn)
	t, ok := ck.actvTran[tn]
	if !ok {
		return nil // it's gone, presumably aborted
	}
	t.end = ck.next()
	if t.start == ck.oldest {
		ck.oldest = math.MaxInt // need to find the new oldest
	}
	// move transaction from active to committed
	delete(ck.actvTran, tn)
	if ut.ct.hasUpdates {
		assert.That(ut.ct.readConflict == "")
		ck.cmtdTran[tn] = t
		tw = ck.tablesWritten(t)
	} else {
		traceln(t, "commit with no updates")
		ck.removeByTable(t)
	}
	ck.cleanEnded()
	return tw
}

func (ck *Check) tablesWritten(t *CkTran) []string {
	tw := make([]string, 0, 8) // ???
	for _, table := range t.tables {
		ta := ck.bytable[table][t.start]
		if ta.outputs != nil || ta.deletes != nil {
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
// It is called by commit and abort.
func (ck *Check) cleanEnded() {
	// find oldest start of non-ended (would be faster with a heap)
	if ck.oldest == math.MaxInt {
		for _, t := range ck.actvTran {
			// assert.That(!t.ended())
			if t.start < ck.oldest {
				ck.oldest = t.start
			}
		}
		traceln("OLDEST", ck.oldest)
	}
	// remove any ended transactions older than this
	for tn, t := range ck.cmtdTran {
		// assert.That(t.ended())
		if t.end < ck.oldest {
			traceln("REMOVE", tn, "->", t.end)
			ck.removeByTable(t)
			delete(ck.cmtdTran, tn)
		}
	}
	for table, end := range ck.exclusive {
		if end < ck.oldest {
			traceln("REMOVE exclusive", table, end)
			delete(ck.exclusive, table)
		}
	}
}

// removeByTable is called by abort and cleanEnded to removeByTable from bytable.
// Note: the caller must still delete from actvTran or cmtdTran.
func (ck *Check) removeByTable(t *CkTran) {
	for _, table := range t.tables {
		ts := ck.bytable[table]
		delete(ts, t.start)
		// keep ck.bytable[table] even if empty
	}
}

// MaxAge is the maximum number of ticks that a transaction can be outstanding.
// Transactions are aborted if they exceed this limit.
var MaxAge = 20

// tick should be called regularly e.g. once per second
// to abort transactions older than MaxAge.
func (ck *Check) tick() {
	ck.clock++
	// traceln("tick", ck.clock)
	for tn, t := range ck.actvTran {
		if ck.clock-t.birth >= MaxAge {
			traceln("abort", tn, "age", ck.clock-t.birth)
			log.Println("aborted", t, "update transaction longer than", MaxAge, "seconds")
			ck.abort(tn, "transaction exceeded max age")
		}
	}
}

func (ck *Check) Stop() { // to satisfy Checker interface
	ck.db.PersistSync() // for tests
}

func traceln(...any) {
	// fmt.Println(args...) // comment out to disable tracing
}

// Transactions returns a list of the active update transactions
func (ck *Check) Transactions() []int {
	trans := make([]int, 0, len(ck.actvTran))
	for _, t := range ck.actvTran {
		// assert.That(!t.ended())
		trans = append(trans, t.start)
	}
	return trans
}

// Final returns the count of committed transactions overlapping with outstanding
func (ck *Check) Final() int {
	return len(ck.cmtdTran)
}
