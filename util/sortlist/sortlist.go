// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package sortlist implements a sorted list of uint64 that is built incrementally.
// The list is in blocks of blockSize.
// Blocks are individually sorted once they are full.
// Blocks are two-way merged.
// A background goroutine is used to do the sorting and merging concurrently.
//
// Blocks are recycled by merges so we use at most two extra blocks.
//
// Note: Zero is used as a terminator, it must not be added as a value.
package sortlist

import (
	"sort"

	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/bits"
	"golang.org/x/exp/slices"
)

const blockSize = 4096

type block[T any] [blockSize]T

type List[T any] struct {
	blocks []*block[T]
	size   int
}

type void struct{}

type Builder[T any] struct {
	zero      func(x T) bool
	less      func(x, y T) bool
	block     *block[T] // current block
	i         int       // index in current block
	blocks    []*block[T]
	free      []*block[T]
	work      chan void // send to tell worker there is something to do
	done      chan any  // worker sends nil or error when finished
	zeroBlock *block[T]
}

// NewSorting returns a new list Builder with incremental sorting.
func NewSorting[T any](zero func(x T) bool, less func(x, y T) bool) *Builder[T] {
	li := &Builder[T]{less: less, zero: zero, blocks: make([]*block[T], 0, 4),
		work: make(chan void), done: make(chan any), zeroBlock: &block[T]{}}
	go li.worker()
	return li
}

// NewUnsorted returns a new list Builder without sorting.
func NewUnsorted[T any](zero func(x T) bool) *Builder[T] {
	li := &Builder[T]{blocks: make([]*block[T], 0, 4), zero: zero,
		zeroBlock: &block[T]{}}
	return li
}

// Add adds a value to the list.
func (b *Builder[T]) Add(x T) {
	assert.That(!b.zero(x))
	if b.block == nil {
		b.block = new(block[T])
	}
	b.block[b.i] = x
	b.i++
	if b.i >= blockSize { // block full
		b.i = 0
		if b.done != nil {
			// wait for worker to finish previous work
			if err := <-b.done; err != nil {
				panic(err)
			}
		}
		b.blocks = append(b.blocks, b.block)
		b.block = nil
		if b.done != nil {
			b.work <- void{} // signal worker to process this block
		}
	}
}

// Finish completes sorting and merging and returns the ordered List.
// No more values should be added after this.
func (b *Builder[T]) Finish() List[T] {
	var zero T
	size := len(b.blocks)*blockSize + b.i
	if b.done == nil {
		if b.block != nil { // partial last block
			b.block[b.i] = zero // terminator
			b.blocks = append(b.blocks, b.block)
			b.block = nil
		}
		return List[T]{blocks: b.blocks, size: size}
	} else {
		// wait for worker to finish previous work
		if err := <-b.done; err != nil {
			panic(err)
		}
		close(b.work) // end worker
	}
	if b.block != nil { // partial last block
		if b.less != nil {
			sort.Sort(ablock[T]{block: b.block, n: b.i, less: b.less})
		}
		b.block[b.i] = zero // terminator
		b.blocks = append(b.blocks, b.block)
		b.block = nil
		b.merges(len(b.blocks))
	}
	b.finishMerges()
	return List[T]{blocks: b.blocks, size: size}
}

// Sort sorts the list by the given compare function.
// It is used to re-sort by a different compare function.
func (b *Builder[T]) Sort(less func(x, y T) bool) {
	var zero T
	b.less = less
	if b.block != nil { // partial last block
		b.block[b.i] = zero // terminator
		b.blocks = append(b.blocks, b.block)
		b.block = nil
	}
	for bi, block := range b.blocks {
		n := blockSize
		if bi == len(b.blocks)-1 {
			if i := slices.IndexFunc(block[:], b.zero); i >= 0 {
				n = i
			}
		}
		sort.Sort(ablock[T]{block: block, n: n, less: less})
		b.merges(bi + 1) // merge as we sort for better cache use
	}
	b.finishMerges()
}

func (b *Builder[T]) finishMerges() {
	nb := len(b.blocks)
	nb2 := bits.NextPow2(uint(nb))
	for len(b.blocks) < nb2 {
		b.blocks = append(b.blocks, b.zeroBlock)
		b.merges(len(b.blocks))
	}
	b.blocks = b.blocks[:nb]
}

func (b *Builder[T]) worker() {
	defer func() {
		if e := recover(); e != nil {
			b.done <- e
		}
	}()
	b.done <- nil
	for range b.work {
		nb := len(b.blocks)
		sort.Sort(ablock[T]{block: b.blocks[nb-1], n: blockSize, less: b.less})
		b.merges(nb)
		b.done <- nil
	}
}

func (b *Builder[T]) merges(nb int) {
	bi := nb - 1
	for mergeSize := 1; bi&mergeSize == mergeSize; mergeSize <<= 1 {
		b.merge(nb, mergeSize)
	}
}

func (b *Builder[T]) merge(nb, size int) {
	leftLast := b.blocks[nb-size-1][blockSize-1]
	rightFirst := b.blocks[nb-size][0]
	if b.zero(rightFirst) || b.less(leftLast, rightFirst) {
		return // nothing to do
	}
	out := newMergeOutput(b)
	aiter := b.iter(nb-2*size, size)
	biter := b.iter(nb-size, size)
	aval, aok := aiter()
	bval, bok := biter()
	for aok && bok {
		if b.less(aval, bval) {
			out.add(aval)
			aval, aok = aiter()
		} else {
			out.add(bval)
			bval, bok = biter()
		}
	}
	for aok {
		out.add(aval)
		aval, aok = aiter()
	}
	for bok {
		out.add(bval)
		bval, bok = biter()
	}
	if out.i < blockSize {
		// add terminator since recycled blocks aren't zeroed
		var zero T
		out.blocks[len(out.blocks)-1][out.i] = zero
	}
	copy(b.blocks[nb-2*size:], out.blocks)
}

func (b *Builder[T]) iter(startBlock, nBlocks int) func() (T, bool) {
	var zero T
	blocks := b.blocks[startBlock : startBlock+nBlocks]
	bi := 0
	i := -1
	return func() (T, bool) {
		if i+1 < blockSize {
			i++
		} else {
			assert.That(blocks[bi] != b.zeroBlock)
			b.free = append(b.free, blocks[bi]) // recycle block
			if bi+1 < len(blocks) {
				bi++
				i = 0
			} else {
				return zero, false // finished
			}
		}
		if b.zero(blocks[bi][i]) {
			return zero, false
		}
		return blocks[bi][i], true
	}
}

// alloc is used by mergeOutput to recycle blocks
func (b *Builder[T]) alloc() *block[T] {
	nf := len(b.free)
	if nf > 0 {
		block := b.free[nf-1]
		b.free = b.free[:nf-1]
		return block
	}
	return new(block[T])
}

// mergeOutput is used to accumulate the result of a merge.
type mergeOutput[T any] struct {
	parent *Builder[T]
	blocks []*block[T]
	i      int // index in current/last block
}

func newMergeOutput[T any](parent *Builder[T]) *mergeOutput[T] {
	return &mergeOutput[T]{blocks: make([]*block[T], 0, 4), i: blockSize,
		parent: parent}
}

func (mo *mergeOutput[T]) add(x T) {
	if mo.i >= blockSize {
		mo.blocks = append(mo.blocks, mo.parent.alloc())
		mo.i = 0
	}
	mo.blocks[len(mo.blocks)-1][mo.i] = x
	mo.i++
}

// ablock handles sorting a possibly partial block
type ablock[T any] struct {
	*block[T]
	n    int
	less func(x, y T) bool
}

func (ab ablock[T]) Len() int {
	return ab.n
}

func (ab ablock[T]) Swap(i, j int) {
	b := ab.block
	b[i], b[j] = b[j], b[i]
}

func (ab ablock[T]) Less(i, j int) bool {
	b := ab.block
	return ab.less(b[i], b[j])
}

// Iter from Builder returns a function that returns 0 when finished
func (b *Builder[T]) Iter() func() T {
	var zero T
	blocks := b.blocks
	if len(blocks) == 0 {
		return func() T { return zero }
	}
	bi := 0
	i := -1
	return func() T {
		if i+1 < blockSize {
			i++
			if b.zero(blocks[bi][i]) {
				return zero // finished
			}
		} else {
			if bi+1 < len(blocks) {
				bi++
				i = 0
			} else {
				return zero // finished
			}
		}
		return blocks[bi][i]
	}
}

// Iter is used by tempindex
type Iter[T any] struct {
	blocks []*block[T]
	size   int
	less   iterLess[T]
	i      int
	state
}

type iterLess[T any] func(x T, key []string) bool

type state int

const (
	rewound state = 1
	eof     state = 2
)

// Iter returns an iterator for the list.
// Warning: The less function (used by seek)
// must be consistent with the sort function.
func (li List[T]) Iter(less iterLess[T]) *Iter[T] {
	return &Iter[T]{blocks: li.blocks, size: li.size, state: rewound, less: less}
}

func (it *Iter[T]) Rewind() {
	it.state = rewound
	it.i = -1
}

func (it *Iter[T]) Next() {
	switch it.state {
	case rewound:
		it.i = 0
		it.state = 0
	case eof:
		// stick
	default:
		it.i++
	}
	if it.i >= it.size {
		it.state = eof
	}
}

func (it *Iter[T]) Prev() {
	switch it.state {
	case rewound:
		it.i = it.size - 1
		it.state = 0
	case eof:
		// stick
	default:
		it.i--
	}
	if it.i < 0 {
		it.state = eof
	}
}

func (it *Iter[T]) Eof() bool {
	return it.state == eof
}

func (it *Iter[T]) Cur() T {
	assert.That(it.state == 0)
	return it.blocks[it.i/blockSize][it.i%blockSize]
}

// Seek does a binary search
func (it *Iter[T]) Seek(key []string) {
	first := 0
	n := it.size
	for n > 0 {
		half := n >> 1
		middle := first + half
		if it.less(it.get(middle), key) {
			first = middle + 1
			n -= half + 1
		} else {
			n = half
		}
	}
	if first >= it.size {
		it.state = eof
		it.i = -1
	} else {
		it.i = first
		it.state = 0
	}
}

func (it *Iter[T]) get(i int) T {
	return it.blocks[i/blockSize][i%blockSize]
}
