// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"math"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/assert"
)

var _ = builtin(SeqQ, "(value)")

func SeqQ(val Value) Value {
	_, ok := val.(*SuSequence)
	return SuBool(ok)
}

var _ = builtin(Seq, "(from=false, to=false, by=1)")

func Seq(from, to, by Value) Value {
	if from == False {
		from = Zero
		to = MaxInt
	} else if to == False {
		to = from
		from = Zero
	}
	f := ToInt(from)
	return NewSuSequence(
		&seqIter{from: f, to: ToInt(to), by: ToInt(by), i: f})
}

type seqIter struct {
	MayLock
	from int
	to   int
	by   int
	i    int
}

func (seq *seqIter) Next() Value {
	if seq.Lock() {
		defer seq.Unlock()
	}
	assert.That(seq.by != 0)
	if seq.i >= seq.to {
		return nil
	}
	i := seq.i
	seq.i += seq.by
	return IntVal(i)
}

func (seq *seqIter) Dup() Iter {
	if seq.Lock() {
		defer seq.Unlock()
	}
	return &seqIter{from: seq.from, to: seq.to, by: seq.by, i: seq.from}
}

func (seq *seqIter) Infinite() bool {
	// to is read-only so no locking required
	return seq.to == math.MaxInt32 // has to match MaxInt opcode
}

func (seq *seqIter) Instantiate() *SuObject {
	n := (seq.to - seq.from + (seq.by - 1)) / seq.by
	InstantiateMax(n)
	list := make([]Value, n)
	i := seq.from
	for j := 0; i < seq.to; j++ {
		list[j] = IntVal(i)
		i += seq.by
	}
	return NewSuObject(list)
}
