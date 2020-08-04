// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package sortlist

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/verify"
)

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
	bldr := NewBuilder(ints.Compare)
	for j := 0; j < nitems; j++ {
		bldr.Add(randint())
	}
	list := bldr.List()
	list.ckblocks()
}

func randint() int {
	time.Sleep(10)
	return 1 + int(rand.Int31()) // +1 so no zeros
}

func (li *List) ckblocks() {
	prev := 0
	for bi, b := range li.blocks {
		for i, x := range b {
			if x == 0 {
				break
			}
			if x < prev {
				fmt.Println("ck", bi, i, "prev", prev, "cur", x)
			}
			verify.That(prev <= x)
			prev = x
		}
	}
}

func TestNextPow2(t *testing.T) {
	Assert(t).That(nextPow2(0), Equals(0))
	Assert(t).That(nextPow2(1), Equals(1))
	Assert(t).That(nextPow2(2), Equals(2))
	Assert(t).That(nextPow2(3), Equals(4))
	Assert(t).That(nextPow2(100), Equals(128))
	Assert(t).That(nextPow2(1000), Equals(1024))
}

//-------------------------------------------------------------------

const nitems = 4 * blockSize // number of blocks must be power of 2 for merging

var G int

func BenchmarkSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		slice := mksimple()
		G = slice[0]
	}
}

func TestSimple(*testing.T) {
	slice := mksimple()
	for i := 1; i < nitems; i++ {
		verify.That(slice[i-1] <= slice[i])
	}
}

func mksimple() []int {
	slice := []int{}
	for j := 0; j < nitems; j++ {
		slice = append(slice, randint())
	}
	sort.Ints(slice)
	return slice
}

//-------------------------------------------------------------------

func BenchmarkChunked(b *testing.B) {
	for i := 0; i < b.N; i++ {
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
	for j := 0; j < nitems; j++ {
		list.Add(randint())
	}
	sort.Sort(list)
	return list
}

func ckblocks(blocks []*block) {
	prev := 0
	for bi, b := range blocks {
		for i, x := range b {
			if x == 0 {
				return
			}
			if x < prev {
				fmt.Println("ck", bi, i, "prev", prev, "cur", x)
			}
			verify.That(prev <= x)
			prev = x
		}
	}
}

//-------------------------------------------------------------------

func BenchmarkMerged(b *testing.B) {
	for i := 0; i < b.N; i++ {
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
	for j := 0; j < nitems; j++ {
		list.Add(randint())
	}
	return list
}

//-------------------------------------------------------------------

func BenchmarkConc(b *testing.B) {
	for i := 0; i < b.N; i++ {
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
	for j := 0; j < nitems; j++ {
		list.Add(randint())
	}
	list.End()
	return list
}

//-------------------------------------------------------------------

type chunked struct {
	blocks []*block
	i      int // index in current/last block
}

func newchunked() *chunked {
	return &chunked{blocks: make([]*block, 0, 4), i: blockSize}
}

func (li *chunked) Add(x int) {
	if li.i >= blockSize {
		li.blocks = append(li.blocks, new(block))
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
	blocks []*block
	i      int // index in current/last block
	free   []*block
}

func newmerged() *merged {
	return &merged{blocks: make([]*block, 0, 4), i: blockSize,
		free: make([]*block, 0, 4)}
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
	verify.That(len(out.blocks) == 2*size)
	verify.That(out.i == blockSize)
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
	blocks []*block
	i      int // index in current/last block
	parent *merged
	// prev   int
}

func newchunked2(parent *merged) *chunked2 {
	return &chunked2{blocks: make([]*block, 0, 4), i: blockSize,
		parent: parent}
}

func (li *chunked2) Add(x int) {
	// verify.That(li.prev <= x)
	// li.prev = x
	if li.i >= blockSize {
		li.blocks = append(li.blocks, li.parent.alloc())
		li.i = 0
	}
	li.blocks[len(li.blocks)-1][li.i] = x
	li.i++
}

func (li *merged) alloc() *block {
	nf := len(li.free)
	if nf > 0 {
		// fmt.Println("using free")
		b := li.free[nf-1]
		li.free = li.free[:nf-1]
		return b
	}
	// fmt.Println("alloc block")
	return new(block)
}

// ablock handles sorting a possibly partial block
type ablock2 struct {
	*block
	n   int
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
	blocks []*block
	i      int // index in current/last block
	free   []*block
	work   chan int
	done   chan void
}

func newconc() *conc {
	li := &conc{blocks: make([]*block, 1, 4), i: blockSize,
		work: make(chan int), done: make(chan void)}
	go li.worker()
	return li
}

func (li *conc) Add(x int) {
	if li.i >= blockSize {
		li.blocks[len(li.blocks)-1] = new(block)
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
	verify.That(len(out.blocks) == 2*size)
	verify.That(out.i == blockSize)
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
	blocks []*block
	i      int // index in current/last block
	parent *conc
	prev   int
}

func newchunked3(parent *conc) *chunked3 {
	return &chunked3{blocks: make([]*block, 0, 4), i: blockSize,
		parent: parent}
}

func (li *chunked3) Add(x int) {
	verify.That(li.prev <= x)
	li.prev = x
	if li.i >= blockSize {
		li.blocks = append(li.blocks, li.parent.alloc())
		li.i = 0
	}
	li.blocks[len(li.blocks)-1][li.i] = x
	li.i++
}

func (li *conc) alloc() *block {
	nf := len(li.free)
	if nf > 0 {
		// fmt.Println("using free")
		b := li.free[nf-1]
		li.free = li.free[:nf-1]
		return b
	}
	// fmt.Println("alloc block")
	return new(block)
}
