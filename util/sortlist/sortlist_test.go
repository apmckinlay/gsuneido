// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package sortlist

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"testing"

	"slices"

	"github.com/apmckinlay/gsuneido/util/assert"
	"golang.org/x/exp/constraints"
)

var z = func(x int) bool { return x == 0 }
var lt = func(x, y int) bool { return x < y }

func sz[T any](s []T) bool { return s == nil }
func slt[T constraints.Ordered](x, y []T) bool {
	return slices.Compare(x, y) < 0
}

func TestSortList_Unsorted(t *testing.T) {
	bldr := NewUnsorted(z)
	bldr.Add(2)
	bldr.Add(1)
	bldr.Sort(lt)
	bldrInorder(bldr, 2)
}

func TestSortList_Sorted(t *testing.T) {
	bldr := NewSorting(z, lt)
	bldr.Add(2)
	bldr.Add(1)
	list := bldr.Finish()
	listInorder(list, 2, z, lt)
}

func TestSortList_SortedSlice(t *testing.T) {
	bldr := NewSorting(sz[int], slt[int])
	bldr.Add([]int{2})
	bldr.Add([]int{1})
	list := bldr.Finish()
	listInorder(list, 2, sz[int], slt[int])
}

func FuzzSort(f *testing.F) {
	f.Fuzz(func(t *testing.T, nb uint8, n2 uint16) {
		testSorting(t, nb, n2)
		testUnsorted(t, nb, n2)
	})
}

func testSorting(_ *testing.T, nb uint8, n2 uint16) {
	n := int(nb)*blockSize + int(n2)
	if n > 100_000 {
		return
	}
	bldr := NewSorting(z, lt)
	for range n {
		bldr.Add(1 + int(rand.Int31())) // +1 so no zeros
	}
	bldr.Finish()
	bldrInorder(bldr, n)
}

func testUnsorted(_ *testing.T, nb uint8, n2 uint16) {
	n := int(nb)*blockSize + int(n2)
	if n > 100_000 {
		return
	}
	bldr := NewUnsorted(z)
	for range n {
		bldr.Add(1 + int(rand.Int31())) // +1 so no zeros
	}
	bldr.Finish()
	bldr.Sort(func(x, y int) bool { return x < y })
	bldrInorder(bldr, n)
}

func TestBuilder(*testing.T) {
	test(0)
	test(1)
	test(10)
	for n := 1; n <= 8; n++ {
		test(n * blockSize)
		test(n*blockSize + 5)
	}
}

func test(nitems int) {
	bldr := NewSorting(z, lt)
	for range nitems {
		bldr.Add(randint())
	}
	list := bldr.Finish()
	assert.This(list.size).Is(nitems)
	listInorder(list, nitems, z, lt)
	bldrInorder(bldr, nitems)

	bldr = NewUnsorted(z)
	for j := 1; j <= nitems; j++ {
		bldr.Add(int(j))
	}
	list = bldr.Finish()
	assert.This(list.size).Is(nitems)
	listInorder(list, nitems, z, lt)
	bldrInorder(bldr, nitems)

	less := func(x int, key []string) bool {
		y, _ := strconv.Atoi(key[0])
		return x < int(y)
	}
	it := list.Iter(less)
	it.Seek([]string{"0"})
	it.Seek([]string{"9999999999"})

	bldr.Sort(func(x, y int) bool { return y < x }) // reverse
}

var N int

func randint() int {
	// small delay to simulate work
	for range 200 {
		N++
	}
	return 1 + int(rand.Int31()) // +1 so no zeros
}

func listInorder[T any](li List[T], nitems int, z func(x T) bool, lt func(x, y T) bool) {
	n := 0
	var prev T
outer:
	for _, b := range li.blocks {
		for _, x := range b {
			if z(x) {
				break outer
			}
			assert.That(n == 0 || !lt(x, prev)) // prev <= x
			prev = x
			n++
		}
	}
	assert.This(n).Is(nitems)
}

func bldrInorder(b *Builder[int], nitems int) {
	n := 0
	prev := int(0)
	iter := b.Iter()
	for x := iter(); x != 0; x = iter() {
		assert.That(prev <= x)
		prev = x
		n++
	}
	assert.This(n).Is(nitems)
}

//-------------------------------------------------------------------

func TestIterEmpty(t *testing.T) {
	b := NewSorting[int](nil, nil)
	list := b.Finish() // empty

	it := list.Iter(nil)
	it.Next()
	assert.T(t).That(it.Eof())
	it.Rewind()
	it.Next()
	assert.T(t).That(it.Eof())

	it = list.Iter(nil)
	it.Prev()
	assert.T(t).That(it.Eof())
	it.Rewind()
	it.Prev()
	assert.T(t).That(it.Eof())

	it.Seek(nil)
}

func TestIterOne(t *testing.T) {
	b := NewUnsorted(z)
	for i := range blockSize {
		b.Add(int(i + 1))
	}
	list := b.Finish() // empty
	less := func(x int, key []string) bool {
		y, _ := strconv.Atoi(key[0])
		return x < int(y)
	}
	it := list.Iter(less)
	it.Seek([]string{"0"})
	it.Seek([]string{"2222"})
	it.Seek([]string{"999999"})
}

func TestIter(t *testing.T) {
	b := NewSorting(z, lt)
	for j := 1; j <= 10; j++ {
		b.Add(int(j))
	}
	list := b.Finish()
	less := func(x int, key []string) bool {
		y, _ := strconv.Atoi(key[0])
		return x < int(y)
	}
	const eof = -1
	it := list.Iter(less)
	test := func(expected int) {
		t.Helper()
		if expected == eof {
			assert.Msg(expected, "should be eof").That(it.Eof())
		} else {
			assert.Msg(expected, "should not be eof").That(!it.Eof())
			assert.This(it.Cur()).Is(int(expected))
		}
	}
	testNext := func(expected int) { it.Next(); t.Helper(); test(expected) }
	testPrev := func(expected int) { it.Prev(); t.Helper(); test(expected) }

	for i := 1; i <= 10; i++ {
		testNext(i)
	}
	testNext(eof)

	it.Rewind()
	for i := 10; i >= 1; i-- {
		testPrev(i)
	}
	testPrev(eof)

	it.Rewind()
	testNext(1)
	testPrev(eof) // stick at eof
	testPrev(eof)
	testNext(eof)

	it.Rewind()
	testPrev(10)
	testPrev(9)
	testPrev(8)
	testNext(9)
	testNext(10) // last
	testPrev(9)

	for i := 1; i <= 10; i++ {
		it.Seek([]string{strconv.Itoa(i)})
		test(i)
	}
}

func TestIterSeek(t *testing.T) {
	b := NewSorting(z, lt)
	for j := 2; j <= 20; j += 2 {
		b.Add(int(j))
	}
	list := b.Finish()
	less := func(x int, key []string) bool {
		y, _ := strconv.Atoi(key[0])
		return x < int(y)
	}
	it := list.Iter(less)
	it.Seek([]string{"5"})
	assert.This(it.Cur()).Is(6)
	it.Seek([]string{"6"})
	assert.This(it.Cur()).Is(6)
}

//-------------------------------------------------------------------

const nitems = 4 * blockSize // number of blocks must be power of 2 for merging

var G int

func BenchmarkSimple(b *testing.B) {
	for range b.N {
		slice := mksimple()
		G = slice[0]
	}
}

func TestSimple(*testing.T) {
	slice := mksimple()
	for i := 1; i < nitems; i++ {
		assert.That(slice[i-1] <= slice[i])
	}
}

func mksimple() []int {
	slice := []int{}
	for range nitems {
		slice = append(slice, randint())
	}
	sort.Sort(uint64Slice(slice))
	return slice
}

type uint64Slice []int

func (p uint64Slice) Len() int { return len(p) }

func (p uint64Slice) Less(i, j int) bool { return p[i] < p[j] }

func (p uint64Slice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

//-------------------------------------------------------------------

func BenchmarkChunked(b *testing.B) {
	for range b.N {
		list := mkchunked()
		G = list.blocks[0][0]
	}
}

func TestChunked(*testing.T) {
	list := mkchunked()
	ckblocks(list.blocks)
}

func mkchunked() *chunked {
	list := newchunked()
	for range nitems {
		list.Add(randint())
	}
	sort.Sort(list)
	return list
}

func ckblocks(blocks []*block[int]) {
	prev := int(0)
	for bi, b := range blocks {
		for i, x := range b {
			if x == 0 {
				return
			}
			if x < prev {
				fmt.Println("ck", bi, i, "prev", prev, "cur", x)
			}
			assert.That(prev <= x)
			prev = x
		}
	}
}

//-------------------------------------------------------------------

func BenchmarkMerged(b *testing.B) {
	for range b.N {
		list := mkmerged()
		G = list.blocks[0][0]
	}
}

func TestMerged(*testing.T) {
	list := mkmerged()
	ckblocks(list.blocks)
}

func mkmerged() *merged {
	list := newmerged()
	for range nitems {
		list.Add(randint())
	}
	return list
}

//-------------------------------------------------------------------

func BenchmarkConc(b *testing.B) {
	for range b.N {
		list := mkconc()
		G = list.blocks[0][0]
	}
}

func TestConc(*testing.T) {
	list := mkconc()
	ckblocks(list.blocks)
}

func mkconc() *conc {
	list := newconc()
	for range nitems {
		list.Add(randint())
	}
	list.End()
	return list
}

//-------------------------------------------------------------------

type chunked struct {
	blocks []*block[int]
	i      int // index in current/last block
}

func newchunked() *chunked {
	return &chunked{blocks: make([]*block[int], 0, 4), i: blockSize}
}

func (li *chunked) Add(x int) {
	if li.i >= blockSize {
		li.blocks = append(li.blocks, new(block[int]))
		li.i = 0
	}
	li.blocks[len(li.blocks)-1][li.i] = x
	li.i++
}

func (li *chunked) Len() int {
	return li.i + blockSize*(len(li.blocks)-1)
}

func (li *chunked) Swap(i, j int) {
	li.blocks[i/blockSize][i%blockSize],
		li.blocks[j/blockSize][j%blockSize] =
		li.blocks[j/blockSize][j%blockSize],
		li.blocks[i/blockSize][i%blockSize]
}

func (li *chunked) Less(i, j int) bool {
	return li.blocks[i/blockSize][i%blockSize] <
		li.blocks[j/blockSize][j%blockSize]
}

//-------------------------------------------------------------------

type merged struct {
	blocks []*block[int]
	i      int // index in current/last block
	free   []*block[int]
}

func newmerged() *merged {
	return &merged{blocks: make([]*block[int], 0, 4), i: blockSize,
		free: make([]*block[int], 0, 4)}
}

func (li *merged) Add(x int) {
	if li.i >= blockSize {
		li.blocks = append(li.blocks, li.alloc())
		li.i = 0
	}
	li.blocks[len(li.blocks)-1][li.i] = x
	li.i++
	if li.i >= blockSize {
		li.sortMerge()
	}
}

func (li *merged) sortMerge() {
	bi := len(li.blocks) - 1
	// fmt.Printf("bi %b\n", bi)
	sort.Sort(ablock2{block: li.blocks[bi], n: blockSize})
	for mergeSize := 1; bi&mergeSize == mergeSize; mergeSize <<= 1 {
		li.merge(mergeSize)
	}
}

func (li *merged) merge(size int) {
	out := newchunked2(li)
	nb := len(li.blocks)
	// fmt.Println("merge size", size, "from", nb-2*size)
	aiter := li.iter(nb-size, size)
	biter := li.iter(nb-2*size, size)
	aval, aok := aiter()
	bval, bok := biter()
	for aok && bok {
		if aval <= bval {
			out.Add(aval)
			aval, aok = aiter()
		} else {
			out.Add(bval)
			bval, bok = biter()
		}
	}
	for aok {
		out.Add(aval)
		aval, aok = aiter()
	}
	for bok {
		out.Add(bval)
		bval, bok = biter()
	}
	assert.This(len(out.blocks)).Is(2 * size)
	assert.This(out.i).Is(blockSize)
	ckblocks(out.blocks)
	// copy blocks from out
	dest := nb - 2*size
	for i, b := range out.blocks {
		li.blocks[dest+i] = b
	}
}

func (li *merged) iter(startBlock, nBlocks int) func() (int, bool) {
	blocks := li.blocks[startBlock : startBlock+nBlocks]
	bi := 0
	i := -1
	return func() (int, bool) {
		if i+1 < blockSize {
			i++
		} else {
			li.free = append(li.free, blocks[bi])
			if bi+1 < len(blocks) {
				bi++
				i = 0
			} else {
				return 0, false // finished
			}
		}
		return blocks[bi][i], true
	}
}

type chunked2 struct {
	blocks []*block[int]
	i      int // index in current/last block
	parent *merged
	// prev   int
}

func newchunked2(parent *merged) *chunked2 {
	return &chunked2{blocks: make([]*block[int], 0, 4), i: blockSize,
		parent: parent}
}

func (li *chunked2) Add(x int) {
	// assert.That(li.prev <= x)
	// li.prev = x
	if li.i >= blockSize {
		li.blocks = append(li.blocks, li.parent.alloc())
		li.i = 0
	}
	li.blocks[len(li.blocks)-1][li.i] = x
	li.i++
}

func (li *merged) alloc() *block[int] {
	nf := len(li.free)
	if nf > 0 {
		// fmt.Println("using free")
		b := li.free[nf-1]
		li.free = li.free[:nf-1]
		return b
	}
	// fmt.Println("alloc block")
	return new(block[int])
}

// ablock handles sorting a possibly partial block
type ablock2 struct {
	*block[int]
	n int
}

func (ab ablock2) Len() int {
	return ab.n
}

func (ab ablock2) Swap(i, j int) {
	b := ab.block
	b[i], b[j] = b[j], b[i]
}

func (ab ablock2) Less(i, j int) bool {
	b := ab.block
	return b[i] < b[j]
}

//-------------------------------------------------------------------

type conc struct {
	blocks []*block[int]
	i      int // index in current/last block
	free   []*block[int]
	work   chan int
	done   chan void
}

func newconc() *conc {
	li := &conc{blocks: make([]*block[int], 1, 4), i: blockSize,
		work: make(chan int), done: make(chan void)}
	go li.worker()
	return li
}

func (li *conc) Add(x int) {
	if li.i >= blockSize {
		li.blocks[len(li.blocks)-1] = new(block[int])
		li.i = 0
	}
	li.blocks[len(li.blocks)-1][li.i] = x
	li.i++
	if li.i >= blockSize {
		nb := len(li.blocks)
		<-li.done
		li.blocks = append(li.blocks, nil)
		li.work <- nb
	}
}

func (li *conc) End() {
	<-li.done
	close(li.work)
	li.blocks = li.blocks[:len(li.blocks)-1]
}

func (li *conc) worker() {
	li.done <- void{}
	for nb := range li.work {
		li.sortMerge(nb)
		li.done <- void{}
	}
}

func (li *conc) sortMerge(nb int) {
	bi := nb - 1
	// fmt.Printf("sortMerge bi %b\n", bi)
	sort.Sort(ablock2{block: li.blocks[bi], n: blockSize})
	for mergeSize := 1; bi&mergeSize == mergeSize; mergeSize <<= 1 {
		li.merge(nb, mergeSize)
	}
}

func (li *conc) merge(nb, size int) {
	out := newchunked3(li)
	// fmt.Println("merge nb", nb, "size", size, "from", nb-2*size)
	aiter := li.iter(nb-size, size)
	biter := li.iter(nb-2*size, size)
	aval, aok := aiter()
	bval, bok := biter()
	for aok && bok {
		if aval <= bval {
			out.Add(aval)
			aval, aok = aiter()
		} else {
			out.Add(bval)
			bval, bok = biter()
		}
	}
	for aok {
		out.Add(aval)
		aval, aok = aiter()
	}
	for bok {
		out.Add(bval)
		bval, bok = biter()
	}
	assert.This(len(out.blocks)).Is(2 * size)
	assert.This(out.i).Is(blockSize)
	ckblocks(out.blocks)
	// copy blocks from out
	dest := nb - 2*size
	for i, b := range out.blocks {
		li.blocks[dest+i] = b
	}
}

func (li *conc) iter(startBlock, nBlocks int) func() (int, bool) {
	blocks := li.blocks[startBlock : startBlock+nBlocks]
	bi := 0
	i := -1
	return func() (int, bool) {
		if i+1 < blockSize {
			i++
		} else {
			li.free = append(li.free, blocks[bi])
			if bi+1 < len(blocks) {
				bi++
				i = 0
			} else {
				return 0, false // finished
			}
		}
		return blocks[bi][i], true
	}
}

type chunked3 struct {
	blocks []*block[int]
	i      int // index in current/last block
	parent *conc
	// prev   int
}

func newchunked3(parent *conc) *chunked3 {
	return &chunked3{blocks: make([]*block[int], 0, 4), i: blockSize,
		parent: parent}
}

func (li *chunked3) Add(x int) {
	// assert.That(li.prev <= x)
	// li.prev = x
	if li.i >= blockSize {
		li.blocks = append(li.blocks, li.parent.alloc())
		li.i = 0
	}
	li.blocks[len(li.blocks)-1][li.i] = x
	li.i++
}

func (li *conc) alloc() *block[int] {
	nf := len(li.free)
	if nf > 0 {
		// fmt.Println("using free")
		b := li.free[nf-1]
		li.free = li.free[:nf-1]
		return b
	}
	// fmt.Println("alloc block")
	return new(block[int])
}
