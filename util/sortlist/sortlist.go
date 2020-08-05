// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package sortlist implements a sorted list of int's that is built incrementally.
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
	"math/bits"
	"sort"
)

const blockSize = 4096

type block [blockSize]uint64

type List struct {
	blocks []*block
}

var zeroBlock block

type void struct{}

type Builder struct {
	cmp    func(x, y uint64) int
	block  *block // current block
	i      int    // index in current block
	blocks []*block
	free   []*block
	work   chan void
	done   chan void
}

// NewBuilder returns a new List Builder.
// It starts the worker goroutine.
func NewBuilder(cmp func(x, y uint64) int) *Builder {
	li := &Builder{cmp: cmp, blocks: make([]*block, 0, 4),
		work: make(chan void), done: make(chan void)}
	go li.worker()
	return li
}

// Add appends a value to the list.
// When a block is full it signals the worker goroutine to process it.
func (b *Builder) Add(x uint64) {
	if b.block == nil {
		b.block = new(block)
		b.i = 0
	}
	b.block[b.i] = x
	b.i++
	if b.i >= blockSize {
		<-b.done
		b.blocks = append(b.blocks, b.block)
		b.block = nil
		b.work <- void{}
	}
}

// List finishes sorting and merging and returns the ordered List.
// The Builder should not be used after this.
// The worker goroutine is stopped.
func (b *Builder) List() List {
	<-b.done
	close(b.work)
	if b.block != nil { // partial last block
		sort.Sort(ablock{block: b.block, n: b.i, cmp: b.cmp})
		b.block[b.i] = 0 // terminator
		b.blocks = append(b.blocks, b.block)
		b.merges()
	}
	nb := len(b.blocks)
	nb2 := nextPow2(nb)
	for len(b.blocks) < nb2 {
		b.blocks = append(b.blocks, &zeroBlock)
		b.merges()
	}
	return List{b.blocks[:nb]}
}

// nextPow2 returns the smallest power of 2 >= n
func nextPow2(n int) int {
	x := uint32(n)
	x = uint32(1) << (32 - bits.LeadingZeros32(x-1))
	return int(x)
}

func (b *Builder) worker() {
	b.done <- void{}
	for range b.work {
		nb := len(b.blocks)
		sort.Sort(ablock{block: b.blocks[nb-1], n: blockSize, cmp: b.cmp})
		b.merges()
		b.done <- void{}
	}
}

func (b *Builder) merges() {
	nb := len(b.blocks)
	bi := nb - 1
	for mergeSize := 1; bi&mergeSize == mergeSize; mergeSize <<= 1 {
		b.merge(nb, mergeSize)
	}
}

func (b *Builder) merge(nb, size int) {
	leftLast := b.blocks[nb-size-1][blockSize-1]
	rightFirst := b.blocks[nb-size][0]
	if rightFirst == 0 || b.cmp(leftLast, rightFirst) <= 0 {
		return // nothing to do
	}
	out := newMergeOutput(b)
	aiter := b.iter(nb-2*size, size)
	biter := b.iter(nb-size, size)
	aval, aok := aiter()
	bval, bok := biter()
	for aok && bok {
		if b.cmp(aval, bval) <= 0 {
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
		out.blocks[len(out.blocks)-1][out.i] = 0 // terminator
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
			if blocks[bi][i] == 0 {
				return 0, false
			}
		} else {
			b.free = append(b.free, blocks[bi]) // recycle block
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
	n   int
	cmp func(x, y uint64) int
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
	return ab.cmp(b[i], b[j]) < 0
}
