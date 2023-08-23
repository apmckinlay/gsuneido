// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"log"
	"math"
	"math/rand"
	"strconv"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/atomic"
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
const MaxTrans = 200

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
	// trans hold the outstanding/overlapping update transactions
	trans map[int]*CkTran
	// exclusive controls access to tables
	exclusive map[string]int
	seq       int
	oldest    int
	// clock is used to abort long transactions
	clock int
}

type CkTran struct {
	// failure is written by CkTran.abort and read by UpdateTran.
	// It is set to either a conflict or a timeout.
	failure   atomic.String
	tables    map[string]*cktbl
	state     *DbState
	start     int
	end       int
	birth     int
	readCount int
}

// readMax is higher than jSuneido to account for duplicate checks
const readMax = 20000

type cktbl struct {
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
	return &Check{db: db, trans: make(map[int]*CkTran), oldest: math.MaxInt,
		exclusive: make(map[string]int)}
}

func (ck *Check) Run(fn func() error) error {
	return fn()
}

func (ck *Check) StartTran() *CkTran {
	if len(ck.trans) >= MaxTrans {
		return nil
	}
	start := ck.next()
	var state *DbState
	if ck.db != nil {
		state = ck.db.GetState()
	}
	t := &CkTran{start: start, end: math.MaxInt, birth: ck.clock,
		tables: make(map[string]*cktbl), state: state}
	ck.trans[start] = t
	return t
}

func (ck *Check) next() int {
	ck.seq++
	return ck.seq
}

// AddExclusive is used for creating indexes on existing tables
func (ck *Check) AddExclusive(table string) bool {
	if end, ok := ck.exclusive[table]; ok && end == math.MaxInt {
		return false // already exclusive
	}
	for _, t2 := range ck.trans {
		if tbl, ok := t2.tables[table]; ok {
			if len(tbl.outputs) > 0 || len(tbl.deletes) > 0 {
				ck.abort(t2.start, "preempted by exclusive ("+table+")")
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
	if !ck.AddExclusive(table) {
		return "already exclusive: " + table
	}
	return ck.RunEndExclusive(table, fn)
}

func (ck *Check) RunEndExclusive(table string, fn func()) (err any) {
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
	t, ok := ck.trans[t.start]
	if !ok {
		return false // it's gone, presumably aborted
	}
	assert.That(!t.ended())
	// check against overlapping transactions
	for _, t2 := range ck.trans {
		if t2 != t && overlap(t, t2) {
			if tbl, ok := t2.tables[table]; ok {
				if tbl.outputs.anyInRange(index, from, to) ||
					tbl.deletes.anyInRange(index, from, to) {
					if ck.abort1of(t, t2, "read", "write", table) {
						return false // this transaction got aborted
					}
				}
			}
		}
	}
	if !t.saveRead(table, index, from, to) {
		ck.abort(t.start, "too many reads in one update transaction")
	}
	return true
}

func (t *CkTran) saveRead(table string, index int, from, to string) bool {
	tbl, ok := t.tables[table]
	if !ok {
		tbl = &cktbl{}
		t.tables[table] = tbl
	}
	reads, inc := tbl.reads.with(index, from, to)
	if t.readCount += inc; t.readCount >= readMax {
		return false
	}
	tbl.reads = reads
	return true
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
	return ck.output(t, table, keys, nil)
}

// output adds an output action.
// oldkeys is used by Update so we can tell if the key changed.
func (ck *Check) output(t *CkTran, table string, keys, oldkeys []string) bool {
	traceln("T", t.start, "output", table, "keys", keys)
	t, ok := ck.trans[t.start]
	if !ok {
		return false // it's gone, presumably aborted
	}
	assert.That(!t.ended())
	if t.start < ck.exclusive[table] {
		ck.abort(t.start, "conflict with exclusive ("+table+")")
		return false
	}
	// check against overlapping transactions
	for _, t2 := range ck.trans {
		if t2 != t && overlap(t, t2) && !t2.ended() {
			if tbl, ok := t2.tables[table]; ok {
				for i, key := range keys {
					if (oldkeys == nil || key != oldkeys[i]) &&
						tbl.reads.contains(i, key) {
						if ck.abort1of(t, t2, "output", "read", table) {
							return false // this transaction got aborted
						}
					}
				}
			}
		}
	}
	if !t.saveOutput(table, keys) {
		ck.abort(t.start,
			"too many writes (output, update, or delete) in one transaction")
	}
	return true
}

func (t *CkTran) saveOutput(table string, keys []string) bool {
	tbl, ok := t.tables[table]
	if !ok {
		tbl = &cktbl{}
		t.tables[table] = tbl
	}
	for i, key := range keys {
		outs := tbl.outputs.with(i, key)
		if outs == nil {
			return false
		}
		tbl.outputs = outs
	}
	return true
}

// Delete adds a delete action.
// Outputs only need to be checked against outstanding reads.
// The keys are parallel with the indexes i.e. keys[i] is for indexes[i].
func (ck *Check) Delete(t *CkTran, table string, off uint64, keys []string) bool {
	traceln("T", t.start, "delete", table, "off", off, "keys", keys)
	t, ok := ck.trans[t.start]
	if !ok {
		return false // it's gone, presumably aborted
	}
	assert.That(!t.ended())
	if t.start < ck.exclusive[table] {
		ck.abort(t.start, "conflict with exclusive ("+table+")")
		return false
	}
	// check against overlapping transactions
	for _, t2 := range ck.trans {
		if t2 != t && overlap(t, t2) {
			if tbl, ok := t2.tables[table]; ok {
				if _, ok := tbl.deloffs[off]; ok {
					if ck.abort1of(t, t2, "delete", "delete", table) {
						return false // this transaction got aborted
					}
				}
				for i, key := range keys {
					if !t2.ended() && tbl.reads.contains(i, key) {
						if ck.abort1of(t, t2, "delete", "read", table) {
							return false // this transaction got aborted
						}
					}
				}
			}
		}
	}
	if !t.saveDelete(table, off, keys) {
		ck.abort(t.start,
			"too many writes (output, update, or delete) in one transaction")
	}
	return true
}

func (t *CkTran) saveDelete(table string, off uint64, keys []string) bool {
	tbl, ok := t.tables[table]
	if !ok {
		tbl = &cktbl{deloffs: make(map[uint64]struct{})}
		t.tables[table] = tbl
	} else if tbl.deloffs == nil {
		tbl.deloffs = make(map[uint64]struct{})
	}
	assert.That(off != 0)
	tbl.deloffs[off] = struct{}{}
	for i, key := range keys {
		dels := tbl.deletes.with(i, key)
		if dels == nil {
			return false
		}
		tbl.deletes = dels
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
	t, ok := ck.trans[t.start]
	if !ok {
		return -1 // it's gone, presumably aborted
	}
	return t.readCount
}

// checkerAbortT1 is used by tests to avoid randomness
var checkerAbortT1 = false

// abort1of aborts one of t1 and t2.
// If t2 is committed, abort t1, otherwise choose randomly.
// It returns true if t1 is aborted, false if t2 is aborted.
func (ck *Check) abort1of(t1, t2 *CkTran, act1, act2, table string) bool {
	traceln("conflict with", t2)
	if t2.ended() || checkerAbortT1 || rand.Intn(2) == 1 {
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
	t, ok := ck.trans[tn]
	if !ok {
		return false
	}
	if reason == "" {
		reason = "abort"
	}
	t.failure.Store(reason)
	delete(ck.trans, tn)
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
	tn := ut.ct.start
	traceln("commit", tn)
	t, ok := ck.trans[tn]
	if !ok {
		return nil // it's gone, presumably aborted
	}

	// reads := 0
	// for _, tbl := range t.tables {
	// 	for _, tr := range tbl.reads {
	// 		reads += tr.Count()
	// 	}
	// }
	// assert.This(t.readCount).Is(reads)

	t.end = ck.next()
	if t.start == ck.oldest {
		ck.oldest = math.MaxInt // need to find the new oldest
	}
	ck.cleanEnded()
	return t.tablesWritten()
}

func (t *CkTran) tablesWritten() []string {
	tw := make([]string, 0, 8)
	for table, tbl := range t.tables {
		if tbl.outputs != nil || tbl.deletes != nil {
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
	if ck.oldest == math.MaxInt {
		for _, t := range ck.trans {
			if t.end == math.MaxInt && t.start < ck.oldest {
				ck.oldest = t.start
			}
		}
		traceln("OLDEST", ck.oldest)
	}
	// remove any ended transactions older than this
	for tn, t := range ck.trans {
		if t.end != math.MaxInt && t.end < ck.oldest {
			traceln("REMOVE", tn, "->", t.end)
			delete(ck.trans, tn)
		}
	}
	for table, end := range ck.exclusive {
		if end < ck.oldest {
			traceln("REMOVE exclusive", table, end)
			delete(ck.exclusive, table)
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
	// traceln("tick", ck.clock)
	for tn, t := range ck.trans {
		if ck.clock-t.birth >= MaxAge {
			traceln("abort", tn, "age", ck.clock-t.birth)
			log.Println("aborted", t, "update transaction longer than", MaxAge, "seconds")
			ck.abort(tn, "transaction exceeded max age")
		}
	}
}

func (ck *Check) Stop() { // to satisfy Checker interface
	ck.db.persist(&execPersistSingle{}) // for tests
}

func traceln(...any) {
	// fmt.Println(args...) // comment out to disable tracing
}

// Transactions returns a list of the active update transactions
func (ck *Check) Transactions() []int {
	trans := make([]int, 0, 4)
	for _, t := range ck.trans {
		if !t.ended() {
			trans = append(trans, t.start)
		}
	}
	return trans
}

// Final returns the count of ended transactions overlapping with outstanding
func (ck *Check) Final() int {
	n := 0
	for _, t := range ck.trans {
		if t.ended() {
			n++
		}
	}
	return n
}
