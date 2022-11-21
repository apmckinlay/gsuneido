// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package ranges

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"testing"
	"time"

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
	rs := &Ranges{}
	for rs.tree == nil || rs.tree.size < 100 {
		from, to := random(7)
		rs.Insert(from, to)
	}
	rs.check()
	for n, in := range nums {
		assert(rs.Contains(strconv.Itoa(n + 10000))).Is(in)
	}

	if !testing.Short() {
		for rs.tree.size > 50 {
			from, to := random(61)
			rs.Insert(from, to)
		}
		rs.check()
		for n, in := range nums {
			assert(rs.Contains(strconv.Itoa(n + 10000))).Is(in)
		}
	}

	rs.Insert("10000", "99999")
	assert(rs.check()).Is(1)
}

func TestRandomNonOverlapping(t *testing.T) {
	const n = 2 * nodeSize * 80
	data := make([]string, n)
	randKey := str.UniqueRandom(3, 10)
	for i := 0; i < n; i++ {
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
	rs.ForEach(func(from, to string) {
		if prevTo >= from {
			panic("check: out of order " + prevTo + ", " + from)
		}
		if from > to {
			panic("from > to" + from + ", " + to)
		}
		prevTo = to
		n++
	})
	return n
}

func TestOverflow(t *testing.T) {
	const n = 12000 * 2
	data := make([]string, n)
	rand.Seed(time.Now().UnixNano())
	randKey := str.UniqueRandom(4, 10)
	for i := 0; i < n; i++ {
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

//-------------------------------------------------------------------

type visitor func(from, to string)

func (rs *Ranges) ForEach(fn visitor) {
	if rs.tree == nil {
		rs.leaf.forEach(fn)
	} else {
		for i := 0; i < rs.tree.size; i++ {
			rs.tree.slots[i].leaf.forEach(fn)
		}
	}
}

func (leaf *leafNode) forEach(fn visitor) {
	for i := 0; i < leaf.size; i++ {
		fn(leaf.slots[i].from, leaf.slots[i].to)
	}
}
