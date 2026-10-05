// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"math/rand"
	"sync"
	"testing"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/race"
)

// NOTE: these tests depend on the race detector to find problems

func TestConcurrentSeq(t *testing.T) {
	if !race.Enabled {
		t.Skip("RACE NOT ENABLED")
	}
	to := 1_000_000
	if testing.Short() {
		to = 50_000
	}
	it := seqIter{by: 1, to: to}
	it.SetConcurrent()
	fn := func() {
		for it.Next() != nil {
			it.Infinite()
			it.Dup().Next()
		}
	}
	const nthreads = 6
	for range nthreads {
		go fn()
	}
	fn()
}

func TestConcurrentSequencePerms(t *testing.T) {
	th := permsTestSetup(t)
	makeIter := compile.Constant(`function () {
		return class {
			Next() { return QueryFirst('builtin_perm_private sort id').id }
			Dup() { return this }
			Infinite?() { return false }
			}()
		}`)
	seq := Sequence(th, []core.Value{th.Call(makeIter)}).(*core.SuSequence)
	th.SetPerms(core.AllPerms)
	seq.SetConcurrent()
	it := seq.Iter()
	assert.T(t).This(func() { it.Next() }).Panics("not authorized: builtin_perm_private")
	assert.T(t).This(func() { it.Dup().Next() }).Panics("not authorized: builtin_perm_private")

	seq = Sequence(th, []core.Value{th.Call(makeIter)}).(*core.SuSequence)
	th.SetPerms(nil)
	seq.SetConcurrent()
	assert.T(t).This(seq.Iter().Next()).Is(core.IntVal(42))
	assert.T(t).This(seq.Iter().Dup().Next()).Is(core.IntVal(42))

	makeIter = compile.Constant(`function () {
		return class {
			Next() { Permit.Table('private', 'write') }
			Dup() { return this }
			Infinite?() { return false }
			}()
		}`)
	th.SetPerms(core.AllPerms)
	seq = Sequence(th, []core.Value{th.Call(makeIter)}).(*core.SuSequence)
	seq.SetConcurrent()
	assert.T(t).This(func() { seq.Iter().Next() }).
		Panics("Permit.Table can only be called from Auth")
}

func TestConcurrentSequence(t *testing.T) {
	if !race.Enabled {
		t.Skip("RACE NOT ENABLED")
	}
	size := 100_000
	if testing.Short() {
		size = 10_000
	}
	sq := core.NewSuSequence(&seqIter{by: 1, to: 1_000})
	sq.SetConcurrent()
	var wg sync.WaitGroup
	const nthreads = 6
	for range nthreads {
		wg.Go(func() {
			n := rand.Intn(size)
			for i := range n {
				sq.Infinite()
				sq.Instantiated()
				sq.Iter().Next()
				if i == n/2 {
					sq.ToContainer() // instantiates
				}
			}
			assert.That(sq.Instantiated())
		})
	}
	wg.Wait()
}
