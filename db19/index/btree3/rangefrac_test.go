// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"fmt"
	"math"
	"math/rand/v2"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/util/assert"
)

func BenchmarkBtreeRangeFrac(b *testing.B) {
	rng := rand.New(rand.NewPCG(123, 456))
	key := func(i int) string {
		return fmt.Sprintf("%05d", i)
	}
	const n = 11_000
	bldr := Builder(heapstor(64 * 1024))
	for i := range n {
		assert.That(bldr.Add(key(i), 1))
	}
	bt := bldr.Finish().(*btree)
	for b.Loop() {
		x := key(rng.IntN(n))
		y := key(rng.IntN(n))
		org := min(x, y)
		end := max(x, y)
		bt.RangeFrac(org, end, n)
	}
}

func TestBtreeRangeFrac(t *testing.T) {
	var n int
	var bt *btree
	f := "%05d"
	key := func(i int) string {
		return fmt.Sprintf(f, i)
	}
	makeBtree := func(m int) {
		n = m
		b := Builder(heapstor(64 * 1024))
		for i := range n {
			assert.That(b.Add(key(i), 1))
		}
		bt = b.Finish().(*btree)
		// fmt.Println(bt.Stats())
	}
	var one, all, over, under int
	test := func(org, end string, expected float64) {
		t.Helper()
		frac := bt.rangeFrac(org, end)
		diff := frac - expected
		if math.Abs(diff) < .01 {
			one++
		}
		all++
		if math.Abs(diff) > .025 {
			t.Fatal(org, end,
				"got", frac, "expected", expected, "difference", diff)
		}
		if diff >= 0 {
			over++
		} else {
			under++
		}
	}

	// treeLevels 0, full root leaf node
	makeBtree(99)
	assert.Msg("tree levels").This(bt.treeLevels).Is(0)
	test(ixkey.Min, ixkey.Max, 1)
	test(ixkey.Min, ixkey.Min, 0)
	test(ixkey.Max, ixkey.Max, 0)
	for i := 0; i <= n; i++ {
		for j := i; j <= n; j++ {
			test(key(i), key(j), float64(j-i)/float64(n))
		}
	}

	// treeLevels 1, small root
	makeBtree(420)
	assert.Msg("tree levels").This(bt.treeLevels).Is(1)
	assert.This(bt.readTree(bt.root).noffs()).Is(5)
	test(ixkey.Min, ixkey.Max, 1)
	test(ixkey.Min, ixkey.Min, 0)
	test(ixkey.Max, ixkey.Max, 0)
	for i := 0; i <= n; i++ {
		for j := i; j <= n; j++ {
			test(key(i), key(j), float64(j-i)/float64(n))
		}
	}

	// treeLevels 1, small root
	makeBtree(2000)
	assert.Msg("tree levels").This(bt.treeLevels).Is(1)
	assert.This(bt.readTree(bt.root).noffs()).Is(20)
	test(ixkey.Min, ixkey.Max, 1)
	test(ixkey.Min, ixkey.Min, 0)
	test(ixkey.Max, ixkey.Max, 0)
	for i := 0; i <= n; i += 7 {
		for j := i; j <= n; j += 7 {
			test(key(i), key(j), float64(j-i)/float64(n))
		}
	}

	// treeLevels 2, small root, right edge empty
	makeBtree(11000)
	assert.Msg("tree levels").This(bt.treeLevels).Is(2)
	assert.This(bt.readTree(bt.root).noffs()).Is(2)
	test(ixkey.Min, ixkey.Max, 1)
	test(ixkey.Min, ixkey.Min, 0)
	test(ixkey.Max, ixkey.Max, 0)
	for i := 0; i <= n; i += 17 {
		for j := i; j <= n; j += 17 {
			test(key(i), key(j), float64(j-i)/float64(n))
		}
	}

	if testing.Short() {
		t.Skip("Skipping long test")
	}

	// treeLevels 2, small root, right edge full
	makeBtree(20000)
	assert.Msg("tree levels").This(bt.treeLevels).Is(2)
	assert.This(bt.readTree(bt.root).noffs()).Is(2)
	test(ixkey.Min, ixkey.Max, 1)
	test(ixkey.Min, ixkey.Min, 0)
	test(ixkey.Max, ixkey.Max, 0)
	for i := 0; i <= n; i += 17 {
		for j := i; j <= n; j += 17 {
			test(key(i), key(j), float64(j-i)/float64(n))
		}
	}

	// treeLevels 3, small root, right edge empty
	f = "%07d"
	makeBtree(1_040_000)
	assert.Msg("tree levels").This(bt.treeLevels).Is(3)
	assert.This(bt.readTree(bt.root).noffs()).Is(2)
	test(ixkey.Min, ixkey.Max, 1)
	test(ixkey.Min, ixkey.Min, 0)
	test(ixkey.Max, ixkey.Max, 0)
	for i := 0; i <= n; i += 293 {
		for j := i + 1; j <= n; j += 293 {
			test(key(i), key(j), float64(j-i)/float64(n))
		}
	}

	// fmt.Println("one", one, "all", all, float64(one)/float64(all))
	// fmt.Println("over", over, "under", under, "ratio", float64(over)/float64(under))
}
