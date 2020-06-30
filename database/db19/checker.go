// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package db19

import (
	"github.com/apmckinlay/gsuneido/util/ints"
)

type Check struct {
	seq    int
	oldest int
	trans  map[int]*cktran
}

type cktran struct {
	start int
	end   int
}

func NewCheck() *Check {
	return &Check{trans: make(map[int]*cktran), oldest: ints.MaxInt}
}

func (ck *Check) StartTran() int {
	start := ck.next()
	ck.trans[start] = &cktran{start: start, end: ints.MaxInt}
	return start
}

func (ck *Check) next() int {
	ck.seq++
	return ck.seq
}

func (ck *Check) Abort(tn int) {
	delete(ck.trans, tn)
	if tn == ck.oldest {
		ck.oldest = ints.MaxInt // need to find the new oldest
	}
	ck.cleanEnded()
}

func (ck *Check) Commit(tn int) {
	t := ck.trans[tn]
	t.end = ck.next()
	if t.start == ck.oldest {
		ck.oldest = ints.MaxInt // need to find the new oldest
	}
	ck.cleanEnded()
}

func overlap(t1, t2 *cktran) bool {
	return t1.end > t2.start && t2.end > t1.start
}

func (ck *Check) cleanEnded() {
	// find oldest start of non-ended (would be faster with a heap)
	if ck.oldest == ints.MaxInt {
		for _, t := range ck.trans {
			if t.end == ints.MaxInt && t.start < ck.oldest {
				ck.oldest = t.start
			}
		}
		// fmt.Println("OLDEST", ck.oldest)
	}
	// remove any ended transactions older than this
	for tn, t := range ck.trans {
		if t.end != ints.MaxInt && t.end < ck.oldest {
			// fmt.Println("REMOVE", tn, "->", t.end)
			delete(ck.trans, tn)
		}
	}
}
