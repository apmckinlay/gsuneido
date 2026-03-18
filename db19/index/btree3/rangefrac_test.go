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

// TestBtreeRangeFracNonNegative verifies that rangeFrac never returns negative
// when the tree is sparse (fanout < n at some level).
// This can happen when a tree built with many keys has most keys deleted,
// leaving a tree structure where count << n_fat * n^2.
// The test manually constructs such a tree with treeLevels=4.
//
// Structure (treeLevels=4):
//   root -> [child0, child1]  (small root, triggers fattenRoot)
//   child0 -> [t200, empty, empty, empty, empty]  (n_fat=10 after fattenRoot)
//   child1 -> [empty, empty, empty, empty, empty]
//   t200 -> [t3a, t3c, t3ef, empty, empty]  (n=5 at level 2)
//   t3a -> [leaf_a(a1,a2,a3), leaf_b(b1)]
//   t3c -> [leaf_c(c1), leaf_d(d1)]
//   ...
//
// rangeFrac("b1", "c1"):
//   fat root: orgPos=0, endPos=0 (both in t200) -> descend with atRoot=false
//   level 2 (t200): orgPos=0, endPos=1 -> diverge with atRoot=false
//   level+1=3 != treeLevels=4 -> close-children path (not leafRangeFrac)
//   fanout=(8/10)^(1/3)=0.928 < n=5 -> result is negative
func TestBtreeRangeFracNonNegative(t *testing.T) {
	hs := heapstor(64 * 1024)

	leaf_a := makeLeaf("a1", 1, "a2", 2, "a3", 3)
	leaf_b := makeLeaf("b1", 4)
	leaf_c := makeLeaf("c1", 5)
	leaf_d := makeLeaf("d1", 6)
	leaf_ef := makeLeaf("e1", 7, "f1", 8)
	emptyLeaf := makeLeaf()

	la := leaf_a.write(hs)
	lb := leaf_b.write(hs)
	lc := leaf_c.write(hs)
	ld := leaf_d.write(hs)
	lef := leaf_ef.write(hs)
	emptyO := emptyLeaf.write(hs)

	t3a := makeTree(la, "b1", lb)
	t3c := makeTree(lc, "d1", ld)
	t3ef := makeTree(lef, "~", emptyO)
	emptyT3 := makeTree(emptyO, "~", emptyO)

	t3ao := t3a.write(hs)
	t3co := t3c.write(hs)
	t3efo := t3ef.write(hs)
	emptyT3o := emptyT3.write(hs)

	t200 := makeTree(t3ao, "c1", t3co, "e1", t3efo, "~", emptyT3o, "~~", emptyT3o)
	emptyT2 := makeTree(emptyT3o, "~", emptyT3o, "~~", emptyT3o, "~~~", emptyT3o, "~~~~", emptyT3o)

	t200o := t200.write(hs)
	emptyT2o := emptyT2.write(hs)

	child0 := makeTree(t200o, "~", emptyT2o, "~~", emptyT2o, "~~~", emptyT2o, "~~~~", emptyT2o)
	child1 := makeTree(emptyT2o, "~~~~~", emptyT2o, "~~~~~~", emptyT2o, "~~~~~~~", emptyT2o, "~~~~~~~~", emptyT2o)

	child0o := child0.write(hs)
	child1o := child1.write(hs)

	root := makeTree(child0o, "j1", child1o)
	rooto := root.write(hs)

	bt := &btree{stor: hs, root: rooto, treeLevels: 4, count: 8}

	// rangeFrac must never return negative
	// "b1" is near the end of t3a (i_org=1, noffs=2)
	// "c1" is near the beginning of t3c (i_end=0, noffs=2)
	frac := bt.rangeFrac("b1", "c1")
	if frac < 0 {
		t.Errorf("rangeFrac(b1, c1) = %g < 0 (should be non-negative)", frac)
	}
}
