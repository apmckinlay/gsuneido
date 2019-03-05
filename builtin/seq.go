package builtin

import (
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/util/ints"
)

var _ = builtin3("Seq(from, to=false, by=1)",
	func(from, to, by Value) Value {
		if from == False {
			from = Zero
			to = MaxInt
		} else if to == False {
			to = from
			from = Zero
		}
		f := ToInt(from)
		return NewSuSequence(&seqIter{f, ToInt(to), ToInt(by), f})
	})

type seqIter struct {
	from int
	to   int
	by   int
	i    int
}

var _ SeqIter = (*seqIter)(nil)

func (seq *seqIter) Next() Value {
	if seq.i >= seq.to {
		return nil
	}
	i := seq.i
	seq.i += seq.by
	return IntToValue(i)
}

func (seq *seqIter) Dup() SeqIter {
	return &seqIter{seq.from, seq.to, seq.by, 0}
}

func (seq *seqIter) Infinite() bool {
	return seq.to == ints.MaxInt
}
