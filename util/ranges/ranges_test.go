// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ranges

import (
	"fmt"
	stditer "iter"
	"math/rand"
	"sort"
	"strconv"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestContains(t *testing.T) {
	assert := assert.T(t)
	rs := &Ranges{}
	rs.Insert("b", "e")
	rs.Insert("i", "k")
	assert.This(rs.check()).Is(2)

	assert.False(rs.Contains("a"))
	assert.True(rs.Contains("b"))
	assert.True(rs.Contains("c"))
	assert.True(rs.Contains("e"))
	assert.False(rs.Contains("f"))
	assert.False(rs.Contains("h"))
	assert.True(rs.Contains("i"))
	assert.True(rs.Contains("j"))
	assert.True(rs.Contains("k"))
	assert.False(rs.Contains("z"))
}

func TestRanges(t *testing.T) {
	test := func(from, to, expected string) {
		t.Helper()
		rs := &Ranges{}
		rs.Insert("c", "e")
		rs.Insert("i", "k")
		assert.T(t).This(rs.String()).Is("c->e i->k")
		rs.Insert(from, to)
		assert.T(t).This(rs.String()).Is(expected)
		rs.check()
	}
	// overlap both
	test("a", "m", "a->m")
	test("d", "j", "c->k")
	// contained
	test("c", "d", "c->e i->k")
	test("c", "e", "c->e i->k")
	test("d", "e", "c->e i->k")
	test("i", "j", "c->e i->k")
	test("i", "k", "c->e i->k")
	test("j", "k", "c->e i->k")
	// overlap one
	test("a", "d", "a->e i->k")
	test("d", "f", "c->f i->k")
	test("a", "f", "a->f i->k")
	test("h", "j", "c->e h->k")
	test("j", "m", "c->e i->m")
	test("h", "m", "c->e h->m")
	// no overlap
	test("a", "b", "a->b c->e i->k")
	test("f", "g", "c->e f->g i->k")
	test("m", "n", "c->e i->k m->n")
}

func TestRandom(t *testing.T) {
	assert := assert.T(t).This
	var nums [90000]bool
	random := func(rlen int) (string, string) {
		from := 10000 + rand.Intn(90000-100)
		to := from + rand.Intn(rlen)
		for i := from; i <= to; i++ {
			nums[i-10000] = true
		}
		return strconv.Itoa(from), strconv.Itoa(to)
	}
	incs := 0
	rs := &Ranges{}
	for rs.tree == nil || rs.tree.size < 100 {
		from, to := random(7)
		incs += rs.Insert(from, to)
	}
	count := rs.check()
	assert(incs).Is(count)
	for n, in := range nums {
		assert(rs.Contains(strconv.Itoa(n + 10000))).Is(in)
	}

	if !testing.Short() {
		for rs.tree.size > 50 {
			from, to := random(61)
			incs += rs.Insert(from, to)
		}
		count := rs.check()
		assert(incs).Is(count)
		for n, in := range nums {
			assert(rs.Contains(strconv.Itoa(n + 10000))).Is(in)
		}
	}

	incs += rs.Insert("10000", "99999")
	assert(rs.check()).Is(1)
	assert(incs).Is(1)
}

func TestRandomNonOverlapping(t *testing.T) {
	const n = 2 * nodeSize * 80
	data := make([]string, n)
	randKey := str.UniqueRandom(3, 10)
	for i := range n {
		data[i] = randKey()
	}
	sort.Strings(data)
	rs := &Ranges{}
	for i := 0; i < n; i += 2 {
		rs.Insert(data[i], data[i+1])
	}
	rs.check()
	expect := func(expected bool, val string) {
		if expected != rs.Contains(val) {
			t.Error("expected", expected, "for", val)
		}
	}
	for i := 0; i < n; i += 2 {
		expect(false, smaller(data[i]))
		expect(true, bigger(data[i]))
		expect(true, smaller(data[i+1]))
		expect(false, bigger(data[i+1]))
	}
}

func bigger(s string) string {
	return s + "+"
}

func smaller(s string) string {
	last := len(s) - 1
	return s[:last] + string(s[last]-1) + "~"
}

func (rs *Ranges) check() int {
	n := 0
	prevTo := ""
	first := true
	for from, to := range rs.All() {
		if !first && prevTo >= from {
			panic("check: out of order " + prevTo + ", " + from)
		}
		if from > to {
			panic("from > to" + from + ", " + to)
		}
		prevTo = to
		first = false
		n++
	}
	return n
}

func TestOverflow(t *testing.T) {
	const n = 12000 * 3
	data := make([]string, n)
	randKey := str.UniqueRandom(4, 10)
	for i := range n {
		data[i] = randKey()
	}
	sort.Strings(data)
	order := rand.Perm(n / 2)
	rs := &Ranges{}
	for i, o := range order {
		if rs.Insert(data[o*2], data[o*2+1]) == Full {
			fmt.Println(i)
			return
		}
	}
	t.Fatal("should have overflowed")
}

func TestAddReturn(t *testing.T) {
	rs := &Ranges{}
	n := 0
	n += rs.Insert(" ", "a")
	n += rs.Insert(" ", "b")
	n += rs.Insert(" ", "c")
	assert.This(n).Is(rs.check())
}

func TestEmptyStringRange(t *testing.T) {
	assert := assert.T(t)
	rs := &Ranges{}
	
	// Test that empty Ranges doesn't contain empty string
	assert.False(rs.Contains("")) // Should return false for empty Ranges
	assert.This(rs.String()).Is("") // Should be empty string
	assert.This(rs.check()).Is(0) // Should have 0 ranges
	
	// Test inserting empty string range into empty Ranges
	inc := rs.Insert("", "")
	assert.This(inc).Is(Added) // Should return Added (1), not Existed (0)
	assert.This(rs.String()).Is("->") // Should contain the empty range
	assert.True(rs.Contains("")) // Should contain empty string
	assert.This(rs.check()).Is(1) // Should have 1 range
	
	// Test inserting the same empty string range again
	inc2 := rs.Insert("", "")
	assert.This(inc2).Is(Existed) // Should return Existed (0) now
	assert.This(rs.String()).Is("->") // Should still contain just the empty range
	assert.This(rs.check()).Is(1) // Should still have 1 range
}

//-------------------------------------------------------------------

func (rs *Ranges) All() stditer.Seq2[string, string] {
	return func(yield func(from, to string) bool) {
		if rs.tree == nil {
			rs.leaf.forEach(yield)
		} else {
			for i := range rs.tree.size {
				if !rs.tree.slots[i].leaf.forEach(yield) {
					return
				}
			}
		}
	}
}

func (leaf *leafNode) forEach(yield func(from, to string) bool) bool {
	for i := range leaf.size {
		if !yield(leaf.slots[i].from, leaf.slots[i].to) {
			return false
		}
	}
	return true
}
