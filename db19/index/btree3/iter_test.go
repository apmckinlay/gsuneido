// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"fmt"
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
	b.shouldSplit = func(nd node) bool {
		return nd.noffs() >= split
	}
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
	b.shouldSplit = func(nd node) bool {
		return nd.noffs() >= 8
	}
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
	for i := 0; i < 10; i++ {
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

func (it *Iterator) Print() {
	fmt.Println("Iterator:")
	for i := 0; i < len(it.tree) && len(it.tree[i].nd) > 0; i++ {
		fmt.Println(strconv.Itoa(it.tree[i].i), it.tree[i].nd.String())
	}
	fmt.Println(strconv.Itoa(it.leaf.i), it.leaf.nd.String())
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
