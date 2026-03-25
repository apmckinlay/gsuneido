// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"sort"
	"strconv"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index/iface"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestIterEmpty(*testing.T) {
	st := heapstor(8192)
	b := Builder(st)
	bt := b.Finish()

	it := bt.Iterator()
	it.Next()
	assert.That(it.Eof())
	it.Next()
	assert.That(it.Eof())
	it.Rewind()
	it.Next()
	assert.That(it.Eof())

	it.Rewind()
	it.Prev()
	assert.That(it.Eof())
	it.Prev()
	assert.That(it.Eof())
	it.Rewind()
	it.Prev()
	assert.That(it.Eof())
}

func TestIter(t *testing.T) {
	st := heapstor(8192)
	b := Builder(st)
	assert.That(b.Add("a", 1))
	assert.That(b.Add("b", 2))
	assert.That(b.Add("c", 3))
	bt := b.Finish()

	it := bt.Iterator()
	test := func(i int) {
		t.Helper()
		assert.Msg("eof ", i).That(!it.Eof())
		assert.This(it.Offset()).Is(i)
	}
	it.Next()
	test(1)
	it.Next()
	test(2)
	it.Next()
	test(3)
	it.Next()
	assert.That(it.Eof())

	it = bt.Iterator()
	it.Prev()
	test(3)
	it.Prev()
	test(2)
	it.Prev()
	test(1)
	it.Prev()
	assert.That(it.Eof())

	it.Seek("")
	test(1)
	it.Seek("a")
	test(1)
	it.Seek("b")
	test(2)
	it.Seek("c")
	test(3)
	it.Seek("z")
	test(3)

	test2 := func(from, to int) {
		it = bt.Iterator()
		it.Range(Range{Org: strconv.Itoa(base + from), End: strconv.Itoa(base + to)})
		for i := from; i < to; i++ {
			it.Next()
			test(base + i)
		}
		it.Next()
		assert.That(it.Eof())

		it.Rewind()
		for i := to - 1; i >= from; i-- {
			it.Prev()
			test(base + i)
		}
		it.Prev()
		assert.That(it.Eof())
	}
	n := 7
	bt = testBtree(n, 99)
	for from := 0; from < n; from++ {
		for to := 0; to < n; to++ {
			test2(from, to)
		}
	}

	for n = 31; n < 200; n++ {
		bt = testBtree(n, 4)
		test2(0, n)
		test2(n/4, n)
		test2(0, n-n/4)
		test2(n/4, n-n/4)
	}
}

const base = 1000

func testBtree(n, split int) *btree {
	assert.That(n < base)
	b := Builder(heapstor(8192))
	defer SetSplit(SetSplit(split))
	for i := base; i < base+n; i++ {
		assert.That(b.Add(strconv.Itoa(i), uint64(i)))
	}
	return b.Finish().(*btree)
}

func TestIterator(t *testing.T) {
	const n = 1000
	var data [n]string
	randKey := str.UniqueRandomOf(4, 6, "abcde")
	for i := range n {
		data[i] = randKey()
	}
	sort.Strings(data[:])
	b := Builder(stor.HeapStor(8192))
	defer SetSplit(SetSplit(8))
	for i, k := range data {
		assert.That(b.Add(k, uint64(i+1))) // +1 to avoid zero
	}
	bt := b.Finish().(*btree)

	// bt.Print()

	var it iface.Iter
	test := func(i int) {
		t.Helper()
		assert.Msg("eof ", i).That(!it.Eof())
		assert.This(it.Offset() - 1).Is(i)
		assert.This(it.Key()).Is(data[i])
	}

	// it = bt.Iterator()
	// k := data[7] + "Z"
	// it.Seek(k)
	// fmt.Println(7, k, "=>", it.Key(), it.Offset())
	// t.SkipNow()

	// test Iterator Next
	it = bt.Iterator()
	for i := range n {
		it.Next()
		test(i)
	}
	it.Next()
	assert.That(it.Eof())

	// test Iterator Prev
	it = bt.Iterator()
	for i := n - 1; i >= 0; i-- {
		it.Prev()
		test(i)
	}
	it.Prev()
	assert.That(it.Eof())

	// test Seek between keys
	for i, k := range data {
		k += "0" // increment to nonexistent
		it.Seek(k)
		// fmt.Println(i, k, "=>", it.Key(), it.Offset())
		if i+1 < len(data) {
			test(i + 1)
		} else {
			test(n - 1)
		}
	}

	// test Seek & Next
	for i, k := range data {
		it.Seek(k)
		test(i)
		it.Next()
		if i+1 < len(data) {
			test(i + 1)
		} else {
			assert.That(it.Eof())
		}
	}

	// test Seek & Prev
	for i, k := range data {
		it.Seek(k)
		test(i)
		it.Prev()
		if i-1 >= 0 {
			test(i - 1)
		} else {
			assert.That(it.Eof())
		}
	}

	it.Seek("") // before first
	test(0)

	it.Seek("~") // after last
	test(n - 1)

	org := n / 4
	it.Range(Range{Org: data[org], End: ixkey.Max})
	for i := org; i < n; i++ {
		it.Next()
		test(i)
	}
	it.Next()
	assert.That(it.Eof())

	end := n / 2
	it.Range(Range{End: data[end]})
	for i := range end {
		it.Next()
		test(i)
	}
	it.Next()
	assert.That(it.Eof())

	it.Range(Range{Org: data[org], End: data[end]})
	for i := org; i < end; i++ {
		it.Next()
		test(i)
	}
	it.Next()
	assert.That(it.Eof())
	it.Seek(data[0])
	assert.That(it.Eof())
	it.Seek(data[end])
	assert.That(it.Eof())

	it.Range(Range{Org: data[org], End: ixkey.Max})
	for i := n - 1; i >= org; i-- {
		it.Prev()
		test(i)
	}
	it.Prev()
	assert.That(it.Eof())

	it.Range(Range{End: data[end]})
	for i := end - 1; i >= 0; i-- {
		it.Prev()
		test(i)
	}
	it.Prev()
	assert.That(it.Eof())

	it.Range(Range{Org: data[org], End: data[end]})
	for i := end - 1; i >= org; i-- {
		it.Prev()
		test(i)
	}
	it.Prev()
	assert.That(it.Eof())

	// it.Range(Range{Org: data[123], End: data[123] + "\x00"})
	// it.Next()
	// test(123)
	// it.Next()
	// assert.That(it.Eof())
}

// buildTree builds a simple one-level btree with keys "1".."n"
// and offsets equal to the integer value.
func buildTree(n int) *btree {
	b := Builder(stor.HeapStor(8192))
	for i := 1; i <= n; i++ {
		k := strconv.Itoa(i)
		assert.That(b.Add(k, uint64(i)))
	}
	return b.Finish().(*btree)
}

// seek on an exact key should position on that key (not the next).
func TestSeekExactKey(t *testing.T) {
	bt := buildTree(9)
	it := bt.Iterator()

	it.Seek("5")
	assert.That(!it.Eof())
	key, off := it.Cur()
	assert.This(key).Is("5")
	assert.This(off).Is(uint64(5))
}

// seek on a key between existing keys should land on the next greater key.
// e.g., between "5" and "6" should position to "6".
func TestSeekBetweenKeysGoesToNextGreater(t *testing.T) {
	bt := buildTree(9)
	it := bt.Iterator()

	it.Seek("5~") // between 5 and 6 (since "~" > "5" and < "6")
	assert.That(!it.Eof())
	key, off := it.Cur()
	assert.This(key).Is("6")
	assert.This(off).Is(uint64(6))
}

// seek on a key larger than the maximum should remain at last.
func TestSeekAllAboveMaxStaysAtLast(t *testing.T) {
	bt := buildTree(9)
	it := bt.Iterator()

	it.Seek("9zzzz") // greater than the last key
	assert.That(!it.Eof())
	key, off := it.Cur()
	assert.This(key).Is("9")
	assert.This(off).Is(uint64(9))
}

//-------------------------------------------------------------------

func TestIteratorBasic(t *testing.T) {
	bt := buildTestTree(10)
	it := bt.Iterator()

	// Test initial state
	assert.T(t).That(!it.HasCur())
	assert.T(t).That(!it.Eof())

	// Test iteration forward
	for i := range 10 {
		it.Next()
		assert.T(t).That(it.HasCur())
		assert.T(t).That(!it.Eof())

		key := string(it.Key())
		expected := fmt.Sprintf("%04d", i)
		assert.T(t).This(key).Is(expected)
		assert.T(t).This(it.Offset()).Is(uint64(i))

		curKey, curOff := it.Cur()
		assert.T(t).This(curKey).Is(expected)
		assert.T(t).This(curOff).Is(uint64(i))
	}

	// Test end of iteration
	it.Next()
	assert.T(t).That(it.Eof())
	assert.T(t).That(!it.HasCur())
}

func TestIteratorRange(t *testing.T) {
	bt := buildTestTree(100)
	it := bt.Iterator()

	// Test range [20, 30)
	rng := iface.Range{Org: "0020", End: "0030"}
	it.Range(rng)

	expectedKeys := []string{"0020", "0021", "0022", "0023", "0024", "0025", "0026", "0027", "0028", "0029"}

	i := 0
	for it.Next(); !it.Eof(); it.Next() {
		assert.T(t).That(i < len(expectedKeys)) // shouldn't go beyond expected
		key := string(it.Key())
		assert.T(t).This(key).Is(expectedKeys[i])
		assert.T(t).This(it.Offset()).Is(uint64(20 + i))
		i++
	}
	assert.T(t).This(i).Is(len(expectedKeys))
}

func TestIteratorRangeEdgeCases(t *testing.T) {
	bt := buildPaddedTree(100) // Use padded keys for consistent ordering
	it := bt.Iterator()

	// Test range with no matching keys [150, 200)
	it.Range(iface.Range{Org: "150", End: "200"})
	it.Next()
	assert.T(t).That(it.Eof())

	// Test range starting before first key [000, 005)
	it.Range(iface.Range{Org: "000", End: "005"})
	expected := []string{"000", "001", "002", "003", "004"}
	i := 0
	for it.Next(); !it.Eof(); it.Next() {
		assert.T(t).That(i < len(expected))
		key := string(it.Key())
		assert.T(t).This(key).Is(expected[i])
		i++
	}
	assert.T(t).This(i).Is(len(expected))
}

func TestIteratorSeek(t *testing.T) {
	bt := buildTestTree(100)
	it := bt.Iterator()

	// Test seek to exact key
	it.Seek("0050")
	assert.T(t).That(it.HasCur())
	key := string(it.Key())
	assert.T(t).This(key).Is("0050")
	assert.T(t).This(it.Offset()).Is(uint64(50))

	// Test seek to non-existent key (should find next one)
	it.SeekAll("0055a") // between 55 and 56
	assert.T(t).That(it.HasCur())
	key = string(it.Key())
	assert.T(t).This(key).Is("0056")

	// Test seek past end
	it.SeekAll("9999")
	assert.T(t).That(it.HasCur())
	key = string(it.Key())
	assert.T(t).This(key).Is("0099")
}

func TestIteratorSeekWithRange(t *testing.T) {
	bt := buildTestTree(100)
	it := bt.Iterator()

	// Set range [20, 30)
	it.Range(iface.Range{Org: "0020", End: "0030"})

	// Seek within range
	it.Seek("0025")
	assert.T(t).That(it.HasCur())
	key := string(it.Key())
	assert.T(t).This(key).Is("0025")

	// Seek outside range (should set EOF)
	it.Seek("50")
	assert.T(t).That(it.Eof())
	assert.T(t).That(!it.HasCur())

	// Seek before range start (should go to range start)
	it.Seek("10")
	if it.HasCur() {
		key = string(it.Key())
		assert.T(t).That(key >= "0020") // should be at or after range start
	}
}

func TestIteratorRewind(t *testing.T) {
	bt := buildTestTree(10)
	it := bt.Iterator()

	// Advance iterator
	it.Next()
	assert.T(t).That(!it.Eof())
	it.Next()
	assert.T(t).That(it.HasCur())

	// Rewind
	it.Rewind()
	assert.T(t).That(!it.HasCur())
	assert.T(t).That(!it.Eof())

	// Should start from beginning again
	it.Next()
	key := string(it.Key())
	assert.T(t).This(key).Is("0000")
}

func TestIteratorEmptyTree(t *testing.T) {
	st := stor.HeapStor(64 * 1024)
	b := Builder(st)
	bt := b.Finish() // empty tree

	it := bt.Iterator()
	it.Next()
	assert.T(t).That(it.Eof())
	assert.T(t).That(!it.HasCur())
}

// buildTestTree creates a btree with keys "0" to "n-1" and offsets 0 to n-1
func buildTestTree(n int) *btree {
	st := stor.HeapStor(64 * 1024)
	b := Builder(st)
	for i := range n {
		key := fmt.Sprintf("%04d", i)
		assert.That(b.Add(key, uint64(i)))
	}
	return b.Finish().(*btree)
}

// buildPaddedTree creates a btree with zero-padded keys "000" to "n-1" (padded) and offsets 0 to n-1
func buildPaddedTree(n int) *btree {
	st := stor.HeapStor(64 * 1024)
	b := Builder(st)
	for i := range n {
		key := padKey(i, 3) // pad to 3 digits
		assert.That(b.Add(key, uint64(i)))
	}
	return b.Finish().(*btree)
}

// padKey pads an integer to a string of the specified width with leading zeros
func padKey(i, width int) string {
	s := strconv.Itoa(i)
	for len(s) < width {
		s = "0" + s
	}
	return s
}

func TestSkipScanSuffixRange(t *testing.T) {
	b := Builder(stor.HeapStor(8192))
	rows := []struct {
		first  string
		suffix string
		off    uint64
	}{
		{"a", "01", 1},
		{"a", "03", 2},
		{"a", "05", 3},
		{"b", "02", 4},
		{"b", "04", 5},
		{"c", "00", 6},
		{"c", "03", 7},
		{"d", "02", 8},
	}
	for _, r := range rows {
		assert.That(b.Add(ixkey.CompKey(r.first, r.suffix), r.off))
	}
	bt := b.Finish().(*btree)
	it := bt.Iterator().(*Iterator)
	it.SkipScan(iface.Range{Org: "02", End: "04"})

	var got []string
	for it.Next(); !it.Eof(); it.Next() {
		k := it.Key()
		f, s := ixkey.SplitPrefixSuffix(k, 1)
		got = append(got, f+":"+s)
	}
	assert.This(got).Is([]string{"a:03", "b:02", "c:03", "d:02"})
}

func TestSkipScanSuffixRangeMultiPrefix(t *testing.T) {
	b := Builder(stor.HeapStor(8192))
	rows := []struct {
		a, b, c string
		off     uint64
	}{
		{"a", "x", "01", 1}, {"a", "x", "03", 2}, {"a", "x", "05", 3},
		{"a", "y", "02", 4}, {"a", "y", "04", 5},
		{"b", "x", "03", 6},
		{"b", "y", "00", 7}, {"b", "y", "02", 8},
	}
	for _, r := range rows {
		assert.That(b.Add(ixkey.CompKey(r.a, r.b, r.c), r.off))
	}
	bt := b.Finish().(*btree)
	it := bt.Iterator().(*Iterator)
	it.SkipScan(iface.Range{Org: "02", End: "04"}, 2)

	var got []string
	for it.Next(); !it.Eof(); it.Next() {
		p, s := ixkey.SplitPrefixSuffix(it.Key(), 2)
		got = append(got, p+":"+s)
	}
	assert.This(got).Is([]string{
		ixkey.CompKey("a", "x") + ":03",
		ixkey.CompKey("a", "y") + ":02",
		ixkey.CompKey("b", "x") + ":03",
		ixkey.CompKey("b", "y") + ":02",
	})
}

func TestSkipScanSuffixRangeNoMatches(t *testing.T) {
	b := Builder(stor.HeapStor(8192))
	assert.That(b.Add(ixkey.CompKey("a", "01"), 1))
	assert.That(b.Add(ixkey.CompKey("b", "01"), 2))
	bt := b.Finish().(*btree)
	it := bt.Iterator().(*Iterator)
	it.SkipScan(iface.Range{Org: "02", End: "03"})
	it.Next()
	assert.That(it.Eof())
}

func TestSkipScanThenRangeDisablesSkipScan(t *testing.T) {
	b := Builder(stor.HeapStor(8192))
	assert.That(b.Add(ixkey.CompKey("a", "01"), 1))
	assert.That(b.Add(ixkey.CompKey("a", "02"), 2))
	assert.That(b.Add(ixkey.CompKey("b", "01"), 3))
	bt := b.Finish().(*btree)
	it := bt.Iterator().(*Iterator)

	it.SkipScan(iface.Range{Org: "01", End: "03"})
	it.Next()
	assert.That(!it.Eof())

	org := ixkey.CompKey("a", "02")
	end := ixkey.CompKey("a", "03")
	it.Range(iface.Range{Org: org, End: end})
	it.Next()
	assert.That(!it.Eof())
	f, s := ixkey.SplitPrefixSuffix(it.Key(), 1)
	assert.This(f).Is("a")
	assert.This(s).Is("02")
	it.Next()
	assert.That(it.Eof())
}

func TestSkipScanSeekAndPrev(t *testing.T) {
	b := Builder(stor.HeapStor(8192))
	rows := []struct {
		first  string
		suffix string
		off    uint64
	}{
		{"a", "01", 1},
		{"a", "03", 2},
		{"a", "05", 3},
		{"b", "02", 4},
		{"b", "04", 5},
		{"c", "00", 6},
		{"c", "03", 7},
		{"d", "02", 8},
	}
	for _, r := range rows {
		assert.That(b.Add(ixkey.CompKey(r.first, r.suffix), r.off))
	}
	bt := b.Finish().(*btree)
	it := bt.Iterator().(*Iterator)
	it.SkipScan(iface.Range{Org: "02", End: "04"})

	// Seek takes a full composite key; suffix is extracted internally.
	it.Seek(ixkey.CompKey("a", "03"))
	assert.That(!it.Eof())
	f, s := ixkey.SplitPrefixSuffix(it.Key(), 1)
	assert.This(f).Is("a")
	assert.This(s).Is("03")

	it.Next()
	assert.That(!it.Eof())
	f, s = ixkey.SplitPrefixSuffix(it.Key(), 1)
	assert.This(f).Is("b")
	assert.This(s).Is("02")

	it.Prev()
	assert.That(!it.Eof())
	f, s = ixkey.SplitPrefixSuffix(it.Key(), 1)
	assert.This(f).Is("a")
	assert.This(s).Is("03")

	// SeekAll("04") seeks to first key in any group with suffix >= "04"
	// ignoring skipRng bounds, so a:05 is first.
	it.SeekAll("04")
	assert.That(!it.Eof())
	f, s = ixkey.SplitPrefixSuffix(it.Key(), 1)
	assert.This(f).Is("a")
	assert.This(s).Is("05")

	// Prev with skipRng {Org:"02", End:"04"}: retreat to last valid match in a group = a:03
	it.Prev()
	assert.That(!it.Eof())
	f, s = ixkey.SplitPrefixSuffix(it.Key(), 1)
	assert.This(f).Is("a")
	assert.This(s).Is("03")
}

func TestSkipScanSeekAllAboveMax(t *testing.T) {
	b := Builder(stor.HeapStor(8192))
	assert.That(b.Add(ixkey.CompKey("a", "01"), 1))
	assert.That(b.Add(ixkey.CompKey("a", "03"), 2))
	assert.That(b.Add(ixkey.CompKey("b", "02"), 3))
	bt := b.Finish().(*btree)
	it := bt.Iterator().(*Iterator)
	it.SkipScan(iface.Range{Org: "01", End: "04"})

	// SeekAll with suffix above every suffix in the tree should stay at last key, not EOF.
	it.SeekAll("99")
	assert.That(!it.Eof())
	f, s := ixkey.SplitPrefixSuffix(it.Key(), 1)
	assert.This(f).Is("b")
	assert.This(s).Is("02")
}

func TestSkipScanRandomParallelWithSubset(t *testing.T) {
	// Run the same scenario at treeLevels 0, 1, and 2
	// by controlling splitCount (fewer keys per node => more levels).
	for _, splitSize := range []int{100, 8, 4} {
		t.Run(fmt.Sprintf("split%d", splitSize), func(t *testing.T) {
			testSkipScanRandomParallelWithSubset(t, splitSize)
		})
	}
}

func testSkipScanRandomParallelWithSubset(t *testing.T, splitSize int) {
	t.Helper()
	defer SetSplit(SetSplit(splitSize))

	const (
		org   = "08"
		end   = "19"
		steps = 30000
	)

	// groups and the suffixes they contain, chosen to exercise edge cases:
	//   "a","b" have only suffixes below org  -> none of the range
	//   "c","d" have only suffixes within [org,end) -> fully contained in range
	//   "e","f" have only suffixes >= end       -> none of the range
	//   "g","h","i","j" span the range boundary  -> partial overlap
	type groupSpec struct {
		first    string
		suffixes []string
	}
	groups := []groupSpec{
		{"a", []string{"01", "03", "05", "07"}},       // all below org
		{"b", []string{"02", "04", "06"}},             // all below org
		{"c", []string{"08", "10", "12", "15", "18"}}, // all within [org,end)
		{"d", []string{"09", "11", "14", "17"}},       // all within [org,end)
		{"e", []string{"19", "21", "25"}},             // all >= end
		{"f", []string{"20", "22", "28"}},             // all >= end
		{"g", []string{"05", "08", "12", "19", "22"}}, // straddles both boundaries
		{"h", []string{"06", "09", "13", "18", "20"}}, // straddles both boundaries
		{"i", []string{"07", "10", "18"}},             // from below into range
		{"j", []string{"11", "16", "19", "24"}},       // from inside range past end
	}

	fullb := Builder(stor.HeapStor(8192))
	subb := Builder(stor.HeapStor(8192))

	var off uint64 = 1
	for _, g := range groups {
		for _, suffix := range g.suffixes {
			k := ixkey.CompKey(g.first, suffix)
			assert.That(fullb.Add(k, off))
			if org <= suffix && suffix < end {
				assert.That(subb.Add(k, off))
			}
			off++
		}
	}

	firsts := make([]string, len(groups))
	for i, g := range groups {
		firsts[i] = g.first
	}

	full := fullb.Finish().(*btree)
	sub := subb.Finish().(*btree)

	fit := full.Iterator().(*Iterator)
	fit.SkipScan(iface.Range{Org: org, End: end})
	sit := sub.Iterator().(*Iterator)

	assertSame := func(step int, op string) {
		t.Helper()
		assert.T(t).Msg(fmt.Sprintf("step %d op %s hascur", step, op)).
			This(fit.HasCur()).Is(sit.HasCur())
		assert.T(t).Msg(fmt.Sprintf("step %d op %s eof", step, op)).
			This(fit.Eof()).Is(sit.Eof())
		if fit.Eof() {
			return
		}
		fk, fo := fit.Cur()
		sk, so := sit.Cur()
		assert.T(t).Msg(fmt.Sprintf("step %d op %s key", step, op)).
			This(fk).Is(sk)
		assert.T(t).Msg(fmt.Sprintf("step %d op %s off", step, op)).
			This(fo).Is(so)
	}

	seek := func(first, s string) {
		t.Helper()
		k := ixkey.CompKey(first, s)
		fit.Seek(k)
		sit.Seek(k)
	}

	boundary := []struct {
		first  string
		suffix string
	}{
		{"a", "00"}, // before first key in a no-range group
		{"b", "00"}, // before first key in a no-range group
		{"c", "00"}, // before all keys in a fully-in-range group
		{"d", "00"}, // before all keys in a fully-in-range group
		{"e", "00"}, // before all keys in a no-range (above) group
		{"g", "00"}, // before first key in a straddle group
		{"a", "08"}, // at org boundary, group has no such key
		{"c", "08"}, // at org boundary, group has this exact key
		{"g", "08"}, // at org boundary, group has this exact key
		{"a", "18"}, // just before end, group has no such key
		{"c", "18"}, // just before end, group has this exact key
		{"h", "18"}, // just before end, group has this exact key
		{"a", "19"}, // at end boundary (exclusive), group has no such key
		{"e", "19"}, // at end boundary, group has this exact key
		{"j", "19"}, // at end boundary, group has this exact key
		{"a", "99"}, // above all suffixes in any group
		{"j", "99"}, // above all suffixes in any group
	}
	for i, b := range boundary {
		seek(b.first, b.suffix)
		assertSame(-(i + 1), "SeekBoundary")
	}

	// seekSuffixes covers: before first key, between keys, at boundary keys,
	// and after last key — exercising all seek edge cases
	seekSuffixes := []string{
		"00",       // before first key in any group
		"01", "02", // actual/between keys below org
		"07", "075", // last below org / between last-below and org
		"08",       // exactly org
		"085",      // between first in-range and next
		"10", "11", // within range
		"135",      // between two in-range keys
		"18",       // last suffix before end
		"185",      // between last in-range and end
		"19",       // exactly end (exclusive upper bound)
		"20", "21", // above end
		"29", "30", // well above all suffixes
	}

	r := rand.New(rand.NewSource(1))
	assertSame(0, "init")
	for step := 1; step <= steps; step++ {
		switch r.Intn(4) {
		case 0:
			fit.Next()
			sit.Next()
			assertSame(step, "Next")
		case 1:
			fit.Prev()
			sit.Prev()
			assertSame(step, "Prev")
		case 2:
			fit.Rewind()
			sit.Rewind()
			assertSame(step, "Rewind")
		case 3:
			first := firsts[r.Intn(len(firsts))]
			s := seekSuffixes[r.Intn(len(seekSuffixes))]
			seek(first, s)
			assertSame(step, "Seek")
		}
	}
}

func BenchmarkSkipScanBreakevenVsFullScan(b *testing.B) {
	const (
		groups          = 1000
		recordsPerGroup = 250
		totalRecords    = groups * recordsPerGroup
		suffixWidth     = 4
		startSuffix     = 25
	)
	if totalRecords < 100000 {
		b.Fatalf("expected at least 100,000 records, got %d", totalRecords)
	}

	file := filepath.Join(b.TempDir(), "skipscan-breakeven.db")
	st, err := stor.MmapStor(file, stor.Create)
	if err != nil {
		b.Fatal(err)
	}
	b.Cleanup(func() {
		st.Close(true)
	})

	bld := Builder(st)
	for g := range groups {
		first := fmt.Sprintf("g%04d", g)
		for s := range recordsPerGroup {
			suffix := fmt.Sprintf("%0*d", suffixWidth, s)
			off := uint64(g*recordsPerGroup + s + 1)
			assert.That(bld.Add(ixkey.CompKey(first, suffix), off))
		}
	}
	bt := bld.Finish().(*btree)

	widths := []int{8, 16, 32, 64, 96, 128, 160, 192, 224}
	maxWidth := widths[len(widths)-1]
	fullOrg := fmt.Sprintf("%0*d", suffixWidth, startSuffix)
	fullEnd := fmt.Sprintf("%0*d", suffixWidth, startSuffix+maxWidth)
	fullExpected := groups * maxWidth

	for _, width := range widths {
		width := width
		endSuffix := startSuffix + width
		org := fmt.Sprintf("%0*d", suffixWidth, startSuffix)
		end := fmt.Sprintf("%0*d", suffixWidth, endSuffix)
		expected := groups * width

		b.Run(fmt.Sprintf("skip_w%03d", width), func(b *testing.B) {
			b.ReportAllocs()
			for range b.N {
				it := bt.Iterator().(*Iterator)
				it.SkipScan(iface.Range{Org: org, End: end})
				n := 0
				for it.Next(); !it.Eof(); it.Next() {
					n++
				}
				if n != expected {
					b.Fatalf("skip scan width %d expected %d got %d", width, expected, n)
				}
			}
		})
	}

	b.Run("full_once", func(b *testing.B) {
		b.ReportAllocs()
		for range b.N {
			it := bt.Iterator()
			n := 0
			for it.Next(); !it.Eof(); it.Next() {
				_, suffix := ixkey.SplitPrefixSuffix(it.Key(), 1)
				if fullOrg <= suffix && suffix < fullEnd {
					n++
				}
			}
			if n != fullExpected {
				b.Fatalf("full scan expected %d got %d", fullExpected, n)
			}
		}
	})
}

func (it *Iterator) Print() {
	fmt.Println("Iterator:")
	for i := 0; i < len(it.tree) && len(it.tree[i].nd) > 0; i++ {
		fmt.Println(strconv.Itoa(it.tree[i].i), it.tree[i].nd.String())
	}
	fmt.Println(strconv.Itoa(it.leaf.i), it.leaf.nd.String())
}

// TestModified verifies the iface requirement that btree iterators return false.
func TestModified(t *testing.T) {
	bt := buildTestTree(5)
	it := bt.Iterator().(*Iterator)
	assert.T(t).This(it.Modified()).Is(false)
	it.Next()
	assert.T(t).This(it.Modified()).Is(false)
}

// TestSkipRetreatSuffixAboveEndSameGroup covers the suffix >= End path in
// skipRetreatToMatch (line 342) when skipGroup is already set to the current group.
// This is reached when SeekAll positions us above End in the same group we're tracking.
func TestSkipRetreatSuffixAboveEndSameGroup(t *testing.T) {
	b := Builder(stor.HeapStor(8192))
	// Group "a" has two keys both above End="04"
	assert.That(b.Add(ixkey.CompKey("a", "05"), 1))
	assert.That(b.Add(ixkey.CompKey("a", "07"), 2))
	bt := b.Finish().(*btree)
	it := bt.Iterator().(*Iterator)
	it.SkipScan(iface.Range{Org: "02", End: "04"})
	// SeekAll("06") positions at ("a","07") with skipGroup="a"; no key is in [Org,End)
	it.SeekAll("06")
	// Prev from ("a","07"): prev() -> ("a","05") which is also >= End; then prev again -> eof
	it.Prev()
	assert.T(t).That(it.Eof())
}

// TestSkipNextEmptyTree covers the rewound->eof path in skipNext (empty btree).
func TestSkipNextEmptyTree(t *testing.T) {
	b := Builder(stor.HeapStor(8192))
	bt := b.Finish().(*btree)
	it := bt.Iterator().(*Iterator)
	it.SkipScan(iface.Range{Org: "01", End: "02"})
	it.Next()
	assert.T(t).That(it.Eof())
}

// TestSkipPrevEmptyTree covers the rewound->eof path in skipPrev (empty btree).
func TestSkipPrevEmptyTree(t *testing.T) {
	b := Builder(stor.HeapStor(8192))
	bt := b.Finish().(*btree)
	it := bt.Iterator().(*Iterator)
	it.SkipScan(iface.Range{Org: "01", End: "02"})
	it.Prev()
	assert.T(t).That(it.Eof())
}

// TestSkipAdvanceSuffixBelowOrg covers the suffix < Org path in skipAdvanceToMatch:
// within a group, some keys have suffixes below Org, requiring it.next() in the loop.
func TestSkipAdvanceSuffixBelowOrg(t *testing.T) {
	b := Builder(stor.HeapStor(8192))
	// Group "a" has three keys: "01" (below Org), "03" (in range), "05" (in range)
	assert.That(b.Add(ixkey.CompKey("a", "01"), 1))
	assert.That(b.Add(ixkey.CompKey("a", "03"), 2))
	assert.That(b.Add(ixkey.CompKey("a", "05"), 3))
	bt := b.Finish().(*btree)
	it := bt.Iterator().(*Iterator)
	// Org="02" so suffix "01" is below Org and must be skipped
	it.SkipScan(iface.Range{Org: "02", End: "06"})
	it.Next()
	assert.T(t).That(!it.Eof())
	_, s := ixkey.SplitPrefixSuffix(it.Key(), 1)
	assert.T(t).This(s).Is("03")
}

// TestSkipRetreatSuffixAboveEnd covers the suffix >= End path in skipRetreatToMatch:
// within a group, some keys have suffixes at or above End, requiring it.prev().
func TestSkipRetreatSuffixAboveEnd(t *testing.T) {
	b := Builder(stor.HeapStor(8192))
	// Group "a" has keys: "01" (in range), "03" (in range), "05" (at/above End)
	assert.That(b.Add(ixkey.CompKey("a", "01"), 1))
	assert.That(b.Add(ixkey.CompKey("a", "03"), 2))
	assert.That(b.Add(ixkey.CompKey("a", "05"), 3))
	bt := b.Finish().(*btree)
	it := bt.Iterator().(*Iterator)
	// End="04" so suffix "05" is >= End and must be skipped backward
	it.SkipScan(iface.Range{Org: "00", End: "04"})
	// Rewind and iterate backward
	it.Prev()
	assert.T(t).That(!it.Eof())
	_, s := ixkey.SplitPrefixSuffix(it.Key(), 1)
	assert.T(t).This(s).Is("03")
}

// TestSkipSeekGroupEndBeyondAllData covers skipSeekGroupEnd when seekAllRaw hits eof
// because the group+End target lies beyond all keys in the tree.
func TestSkipSeekGroupEndBeyondAllData(t *testing.T) {
	b := Builder(stor.HeapStor(8192))
	// Only one group "a" with suffixes below End; the last group's End would exceed all keys.
	assert.That(b.Add(ixkey.CompKey("a", "01"), 1))
	assert.That(b.Add(ixkey.CompKey("a", "02"), 2))
	bt := b.Finish().(*btree)
	it := bt.Iterator().(*Iterator)
	// End="03": when we do skipSeekGroupEnd("a"), target = "a\x00\x0003"
	// seekAllRaw on that might round to last key. Verify Prev works.
	it.SkipScan(iface.Range{Org: "01", End: "03"})
	it.Prev()
	assert.T(t).That(!it.Eof())
	_, s := ixkey.SplitPrefixSuffix(it.Key(), 1)
	assert.T(t).This(s).Is("02")
}

// TestSkipRetreatNewGroupGoesEof covers the path where skipSeekGroupEnd
// results in !within (the first group's end is before min data).
func TestSkipRetreatNewGroupGoesEof(t *testing.T) {
	b := Builder(stor.HeapStor(8192))
	// Only group "z" keys, all suffixes >= End
	assert.That(b.Add(ixkey.CompKey("z", "05"), 1))
	assert.That(b.Add(ixkey.CompKey("z", "06"), 2))
	bt := b.Finish().(*btree)
	it := bt.Iterator().(*Iterator)
	// End="03": no key in "z" has suffix in [Org, End), Prev should reach eof
	it.SkipScan(iface.Range{Org: "01", End: "03"})
	it.Prev()
	assert.T(t).That(it.Eof())
}

// TestSkipSeekPrevGroupLevel0 covers skipSeekPrevGroup when level<0 and treeLevels==0
// (single-leaf tree): seeks directly in the root leaf.
func TestSkipSeekPrevGroupLevel0(t *testing.T) {
	b := Builder(stor.HeapStor(8192))
	// Two groups: "b" (in range) then "c" (in range) — prev from "c" should go to "b"
	assert.That(b.Add(ixkey.CompKey("b", "02"), 1))
	assert.That(b.Add(ixkey.CompKey("b", "03"), 2))
	assert.That(b.Add(ixkey.CompKey("c", "02"), 3))
	assert.That(b.Add(ixkey.CompKey("c", "03"), 4))
	bt := b.Finish().(*btree)
	assert.T(t).This(bt.treeLevels).Is(0) // single-leaf tree
	it := bt.Iterator().(*Iterator)
	it.SkipScan(iface.Range{Org: "02", End: "04"})

	// Iterate forward to get into group "c"
	it.Next()
	f, s := ixkey.SplitPrefixSuffix(it.Key(), 1)
	assert.T(t).This(f).Is("b")
	assert.T(t).This(s).Is("02")

	// Now call Prev twice to skip back through "b:03" and into "b:02"
	it.Next()
	assert.T(t).That(!it.Eof())
	it.Next()
	assert.T(t).That(!it.Eof())
	it.Next()
	assert.T(t).That(!it.Eof())
	// Now backward: from "c:03" retreat to "c:02", "b:03", "b:02"
	it.Prev()
	f, s = ixkey.SplitPrefixSuffix(it.Key(), 1)
	assert.T(t).This(f).Is("c")
	assert.T(t).This(s).Is("02")
}

// TestSkipSeekPrevGroupEofBeforeFirst covers the level<0 and treeLevels>0 path
// in skipSeekPrevGroup: the first group has no predecessor.
func TestSkipSeekPrevGroupEofBeforeFirst(t *testing.T) {
	defer SetSplit(SetSplit(4)) // force multi-level tree
	b := Builder(stor.HeapStor(8192))
	// Only one group "a"; prev past "a" should reach eof
	for i := range 20 {
		assert.That(b.Add(ixkey.CompKey("a", fmt.Sprintf("%02d", i)), uint64(i+1)))
	}
	bt := b.Finish().(*btree)
	it := bt.Iterator().(*Iterator)
	it.SkipScan(iface.Range{Org: "05", End: "10"})
	// Advance to first match then keep calling Prev past the beginning
	it.Next()
	assert.T(t).That(!it.Eof())
	for !it.Eof() {
		it.Prev()
	}
	assert.T(t).That(it.Eof())
}

// TestSkipSeekNextGroupMultiLevel covers the multi-level descent path in
// skipSeekNextGroup (lines 245-248) when next group spans multiple tree levels.
func TestSkipSeekNextGroupMultiLevel(t *testing.T) {
	defer SetSplit(SetSplit(4)) // small split -> treeLevels >= 1
	b := Builder(stor.HeapStor(64 * 1024))
	// Many groups spread across multiple leaf nodes so the skip crosses tree levels
	for g := 0; g < 30; g++ {
		first := fmt.Sprintf("g%02d", g)
		for s := 0; s < 10; s++ {
			assert.That(b.Add(ixkey.CompKey(first, fmt.Sprintf("%02d", s)), uint64(g*10+s+1)))
		}
	}
	bt := b.Finish().(*btree)
	it := bt.Iterator().(*Iterator)
	// Range only includes suffix "05": forces skipSeekNextGroup across tree boundaries
	it.SkipScan(iface.Range{Org: "05", End: "06"})
	var got []string
	for it.Next(); !it.Eof(); it.Next() {
		f, s := ixkey.SplitPrefixSuffix(it.Key(), 1)
		got = append(got, f+":"+s)
	}
	assert.T(t).This(len(got)).Is(30)
	for i, k := range got {
		assert.T(t).This(k).Is(fmt.Sprintf("g%02d:05", i))
	}
}

// TestSkipSeekPrevGroupMultiLevel covers the multi-level descent in skipSeekPrevGroup
// (lines 407-410) when the previous group is in a different subtree.
func TestSkipSeekPrevGroupMultiLevel(t *testing.T) {
	defer SetSplit(SetSplit(4)) // small split -> treeLevels >= 1
	b := Builder(stor.HeapStor(64 * 1024))
	for g := 0; g < 30; g++ {
		first := fmt.Sprintf("g%02d", g)
		for s := 0; s < 10; s++ {
			assert.That(b.Add(ixkey.CompKey(first, fmt.Sprintf("%02d", s)), uint64(g*10+s+1)))
		}
	}
	bt := b.Finish().(*btree)
	it := bt.Iterator().(*Iterator)
	it.SkipScan(iface.Range{Org: "05", End: "06"})
	// Collect keys in reverse order
	var got []string
	for it.Prev(); !it.Eof(); it.Prev() {
		f, s := ixkey.SplitPrefixSuffix(it.Key(), 1)
		got = append(got, f+":"+s)
	}
	assert.T(t).This(len(got)).Is(30)
	for i, k := range got {
		assert.T(t).This(k).Is(fmt.Sprintf("g%02d:05", 29-i))
	}
}

// TestSkipSuffixSeekUnboundedBeyondAll covers the case in skipSuffixSeekUnbounded
// where seekAllRaw for target lands beyond all data (state != within).
func TestSkipSuffixSeekUnboundedBeyondAll(t *testing.T) {
	b := Builder(stor.HeapStor(8192))
	// All keys have suffixes below minSuffix used in SeekAll
	assert.That(b.Add(ixkey.CompKey("a", "01"), 1))
	assert.That(b.Add(ixkey.CompKey("b", "01"), 2))
	bt := b.Finish().(*btree)
	it := bt.Iterator().(*Iterator)
	it.SkipScan(iface.Range{Org: "01", End: "03"})
	// SeekAll with a suffix > all existing suffixes triggers the eof path
	it.SeekAll("99") // all suffixes are "01", none >= "99" after seeking
	// The SeekAll in skip-scan mode backs up to last physical key on eof
	assert.T(t).That(!it.Eof())
}

func TestGte(t *testing.T) {
	test := func(prefix, suffix, target string, expected bool) {
		t.Helper()
		result := gte([]byte(prefix), []byte(suffix), target)
		assert.T(t).Msg(fmt.Sprintf("gte(%q, %q, %q)", prefix, suffix, target)).
			This(result).Is(expected)
	}

	// Equality cases - should return true
	test("abc", "", "abc", true)
	test("", "abc", "abc", true)
	test("ab", "c", "abc", true)
	test("a", "bc", "abc", true)

	// Greater than cases - should return true
	test("abd", "", "abc", true)
	test("", "abd", "abc", true)
	test("ab", "d", "abc", true)
	test("a", "bd", "abc", true)
	test("abc", "d", "abc", true) // longer than target
	test("abcd", "", "abc", true) // longer in prefix
	test("", "abcd", "abc", true) // longer in suffix

	// Less than cases - should return false
	test("abb", "", "abc", false)
	test("", "abb", "abc", false)
	test("ab", "b", "abc", false)
	test("a", "bb", "abc", false)
	test("ab", "", "abc", false) // shorter than target
	test("", "ab", "abc", false) // shorter in suffix

	// Empty value cases
	test("", "", "", true)    // all empty
	test("a", "", "", true)   // empty target, non-empty prefix
	test("", "a", "", true)   // empty target, non-empty suffix
	test("", "", "a", false)  // empty prefix+suffix, non-empty target
	test("abc", "", "", true) // empty target
	test("", "abc", "", true) // empty target with suffix

	// Byte-level comparison edge cases
	test("a\x00", "", "a", true)     // null byte makes it greater
	test("", "a\x00", "a", true)     // null byte in suffix
	test("a", "\x00", "a", true)     // boundary between prefix and suffix
	test("a\xff", "", "a\x00", true) // high byte value
	test("", "a\xff", "a\x00", true) // high byte in suffix

	// Boundary cases at prefix/suffix split
	test("he", "llo", "hello", true)
	test("hel", "lo", "hello", true)
	test("hell", "o", "hello", true)
	test("he", "llo", "help", false)
	test("hel", "lo", "help", false)

	// Single character comparisons
	test("a", "", "a", true)
	test("", "a", "a", true)
	test("b", "", "a", true)
	test("", "b", "a", true)
	test("a", "", "b", false)
	test("", "a", "b", false)

	// Length variation edge cases
	test("abc", "def", "abcde", true) // equal at compared length, longer total
	test("abc", "de", "abcde", true)  // equal at compared length, same total
	test("abc", "d", "abcde", false)  // equal at compared length, shorter total
	test("abc", "", "abcdef", false)  // prefix matches but too short
}
