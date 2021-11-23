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
)

const blockSize = 4096

type block [blockSize]uint64

type List struct {
	blocks []*block
	size   int
}

var zeroBlock block

type void struct{}

type Builder struct {
	less   func(x, y uint64) bool
	block  *block // current block
	i      int    // index in current block
	blocks []*block
	free   []*block
	work   chan void
	done   chan void
}

// NewSorting returns a new list Builder with incremental sorting.
func NewSorting(less func(x, y uint64) bool) *Builder {
	li := &Builder{less: less, blocks: make([]*block, 0, 4),
		work: make(chan void), done: make(chan void)}
	go li.worker()
	return li
}

// NewUnsorted returns a new list Builder without sorting.
func NewUnsorted() *Builder {
	li := &Builder{blocks: make([]*block, 0, 4)}
	return li
}

// Add adds a value to the list.
func (b *Builder) Add(x uint64) {
	assert.That(x != 0)
	if b.block == nil {
		b.block = new(block)
		b.i = 0
	}
	b.block[b.i] = x
	b.i++
	if b.i >= blockSize { // block full
		if b.done != nil {
			<-b.done // wait till processing of previous block is finished
		}
		b.blocks = append(b.blocks, b.block)
		b.block = nil
		if b.done != nil {
			b.work <- void{} // single worker to process this block
		}
	}
}

// Finish completes sorting and merging and returns the ordered List.
// No more values should be added after this.
func (b *Builder) Finish() List {
	size := len(b.blocks)*blockSize + b.i
	if b.done == nil {
		if b.block != nil { // partial last block
			b.block[b.i] = 0 // terminator
			b.blocks = append(b.blocks, b.block)
			b.block = nil
		}
		return List{blocks: b.blocks, size: size}
	}
	if b.done != nil {
		<-b.done
		close(b.work) // end worker goroutine
	}
	if b.block != nil { // partial last block
		if b.less != nil {
			sort.Sort(ablock{block: b.block, n: b.i, less: b.less})
		}
		b.block[b.i] = 0 // terminator
		b.blocks = append(b.blocks, b.block)
		b.block = nil
		b.merges(len(b.blocks))
	}
	b.finishMerges()
	return List{blocks: b.blocks, size: size}
}

// Sort sorts the list by the given compare function.
// It is used to re-sort by a different compare function.
func (b *Builder) Sort(less func(x, y uint64) bool) {
	b.less = less
	if b.block != nil { // partial last block
		b.block[b.i] = 0 // terminator
		b.blocks = append(b.blocks, b.block)
		b.block = nil
	}
	for i, block := range b.blocks {
		n := blockSize
		if i == len(b.blocks)-1 {
			n = b.i
		}
		sort.Sort(ablock{block: block, n: n, less: less})
		b.merges(i + 1) // merge as we sort for better cache use
	}
	b.finishMerges()
}

func (b *Builder) finishMerges() {
	nb := len(b.blocks)
	nb2 := bits.NextPow2(uint(nb))
	for len(b.blocks) < nb2 {
		b.blocks = append(b.blocks, &zeroBlock)
		b.merges(len(b.blocks))
	}
	b.blocks = b.blocks[:nb]
}

func (b *Builder) worker() {
	b.done <- void{}
	for range b.work {
		nb := len(b.blocks)
		sort.Sort(ablock{block: b.blocks[nb-1], n: blockSize, less: b.less})
		b.merges(nb)
		b.done <- void{}
	}
}

func (b *Builder) merges(nb int) {
	bi := nb - 1
	for mergeSize := 1; bi&mergeSize == mergeSize; mergeSize <<= 1 {
		b.merge(nb, mergeSize)
	}
}

func (b *Builder) merge(nb, size int) {
	leftLast := b.blocks[nb-size-1][blockSize-1]
	rightFirst := b.blocks[nb-size][0]
	if rightFirst == 0 || b.less(leftLast, rightFirst) {
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
		out.blocks[len(out.blocks)-1][out.i] = 0
	}
	copy(b.blocks[nb-2*size:], out.blocks)
}

func (b *Builder) iter(startBlock, nBlocks int) func() (uint64, bool) {
	blocks := b.blocks[startBlock : startBlock+nBlocks]
	bi := 0
	i := -1
	return func() (uint64, bool) {
		if i+1 < blockSize {
			i++
		} else {
			b.free = append(b.free, blocks[bi]) // recycle block
			if bi+1 < len(blocks) {
				bi++
				i = 0
			} else {
				return 0, false // finished
			}
		}
		if blocks[bi][i] == 0 {
			return 0, false
		}
		return blocks[bi][i], true
	}
}

// alloc is used by mergeOutput to recycle blocks
func (b *Builder) alloc() *block {
	nf := len(b.free)
	if nf > 0 {
		block := b.free[nf-1]
		b.free = b.free[:nf-1]
		return block
	}
	return new(block)
}

// mergeOutput is used to accumulate the result of a merge.
type mergeOutput struct {
	parent *Builder
	blocks []*block
	i      int // index in current/last block
}

func newMergeOutput(parent *Builder) *mergeOutput {
	return &mergeOutput{blocks: make([]*block, 0, 4), i: blockSize,
		parent: parent}
}

func (mo *mergeOutput) add(x uint64) {
	if mo.i >= blockSize {
		mo.blocks = append(mo.blocks, mo.parent.alloc())
		mo.i = 0
	}
	mo.blocks[len(mo.blocks)-1][mo.i] = x
	mo.i++
}

// ablock handles sorting a possibly partial block
type ablock struct {
	*block
	n    int
	less func(x, y uint64) bool
}

func (ab ablock) Len() int {
	return ab.n
}

func (ab ablock) Swap(i, j int) {
	b := ab.block
	b[i], b[j] = b[j], b[i]
}

func (ab ablock) Less(i, j int) bool {
	b := ab.block
	return ab.less(b[i], b[j])
}

func (b *Builder) Iter() func() uint64 {
	blocks := b.blocks
	if len(blocks) == 0 {
		return func() uint64 { return 0 }
	}
	bi := 0
	i := -1
	return func() uint64 {
		if i+1 < blockSize {
			i++
			if blocks[bi][i] == 0 {
				return 0 // finished
			}
		} else {
			if bi+1 < len(blocks) {
				bi++
				i = 0
			} else {
				return 0 // finished
			}
		}
		return blocks[bi][i]
	}
}

// Iter is used by tempindex
type Iter struct {
	blocks []*block
	size   int
	less   func(x uint64, key string) bool
	i      int
	state
}

type state int

const (
	rewound state = 1
	eof     state = 2
)

// Iter returns an iterator for the list.
// Warning: The less function (used by seek)
// must be consistent with the sort function.
func (list List) Iter(less func(x uint64, key string) bool) *Iter {
	return &Iter{blocks: list.blocks, size: list.size, state: rewound, less: less}
}

func (it *Iter) Rewind() {
	it.state = rewound
	it.i = -1
}

func (it *Iter) Next() {
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

func (it *Iter) Prev() {
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

func (it *Iter) Eof() bool {
	return it.state == eof
}

func (it *Iter) Cur() uint64 {
	assert.That(it.state == 0)
	return it.blocks[it.i/blockSize][it.i%blockSize]
}

func (it *Iter) Seek(key string) {
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

func (it *Iter) get(i int) uint64 {
	return it.blocks[i/blockSize][i%blockSize]
}
