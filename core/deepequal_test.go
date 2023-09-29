// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestDeepEqual(t *testing.T) {
	test := func(xkind, ykind, conc int) {
		t.Helper()
		x := mk(xkind, conc&xconc == xconc)
		y := mk(ykind, conc&yconc == yconc)
		defer func() {
			if err := recover(); err != nil {
				t.Error(err, xkind, ykind, conc)
			}
		}()
		if deepEqual(x, y) != (xkind == ykind) {
			t.Error("failed", xkind, ykind, conc)
		}
	}
	for xk := 0; xk < nkinds; xk++ {
		for yk := 0; yk < nkinds; yk++ {
			for c := 0; c <= xyconc; c++ {
				test(xk, yk, c)
			}
		}
	}
}

const (
	xconc  = 1
	yconc  = 2
	xyconc = xconc | yconc
)

const (
	empty = iota
	zero
	one
	few
	nest
	loop
	justNamed
	listAndNamed
	instance
	emptyInstance
	instanceNest
	nkinds
)

func mk(kind int, concurrent bool) Value {
	if kind < instance {
		return mkObject(kind, concurrent)
	}
	return mkInstance(kind, concurrent)
}

func mkObject(kind int, concurrent bool) *SuObject {
	ob := &SuObject{}
	switch kind {
	case zero:
		ob.Add(Zero)
	case one:
		ob.Add(One)
	case few:
		ob.Add(Zero)
		ob.Add(One)
	case nest:
		ob.Add(mk(few, false))
		ob.Add(mk(one, false))
	case loop: // ob => x => ob
		x := mkObject(few, false)
		x.Add(ob)
		ob.Add(x)
	case justNamed:
		ob.Set(SuStr("foo"), Zero)
		ob.Set(SuStr("bar"), One)
	case listAndNamed:
		ob.Add(True)
		ob.Add(False)
		ob.Set(SuStr("foo"), Zero)
		ob.Set(SuStr("bar"), One)
	}
	if concurrent {
		ob.SetConcurrent()
	}
	return ob
}

func mkInstance(kind int, concurrent bool) *SuInstance {
	ob := &SuInstance{MemBase: NewMemBase(), useDeepEquals: true}
	switch kind {
	case instanceNest:
		ob.Data["ob"] = mkObject(nest, false)
		ob.Data["in"] = mkInstance(instance, false)
		fallthrough
	case instance:
		ob.Data["foo"] = Zero
		ob.Data["bar"] = One
	}
	if concurrent {
		ob.SetConcurrent()
	}
	return ob
}

func TestExpand(t *testing.T) {
	assert := assert.T(t).This
	var slice []Value
	expand(&slice, 1)
	assert(len(slice)).Is(1)
	expand(&slice, 10)
	assert(len(slice)).Is(11)
	expand(&slice, 1)
	assert(len(slice)).Is(12)
	expand(&slice, 100)
	assert(len(slice)).Is(112)
}
