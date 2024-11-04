// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hamt

import (
	"iter"
	"math"
	"math/bits"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	. "github.com/apmckinlay/gsuneido/util/bits"
	"github.com/apmckinlay/gsuneido/util/cksum"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
)

type Item[K comparable] interface {
	Key() K
	Hash(K) uint64
	Cksum() uint32
	StorSize() int
	IsTomb() bool
	LastMod() int
	SetLastMod(mod int)
	Write(w *stor.Writer)
}

type Hamt[K comparable, E Item[K]] struct {
	root       *node[K, E]
	mutable    bool
	generation uint32 // if mutable, nodes with this generation are mutable
}

type node[K comparable, E Item[K]] struct {
	vals       []E
	ptrs       []*node[K, E]
	generation uint32
	bmVal      uint32
	bmPtr      uint32
}

const bitsPerItemNode = 5
const maskItem = 1<<bitsPerItemNode - 1

func (ht Hamt[K, E]) IsNil() bool {
	return ht.root == nil
}

func (ht Hamt[K, E]) SameAs(ht2 Hamt[K, E]) bool {
	return ht.root == ht2.root
}

func (ht Hamt[K, E]) MustGet(key K) E {
	it, ok := ht.Get(key)
	if !ok || it.IsTomb() {
		panic("MustGet failed")
	}
	return it
}

func (ht Hamt[K, E]) Get(key K) (E, bool) {
	it := ht.get(key)
	if it == nil {
		var zero E
		return zero, false
	}
	return *it, true
}

func (ht Hamt[K, E]) get(key K) *E {
	nd := ht.root
	if nd == nil {
		return nil
	}
	var z E
	hash := z.Hash(key)
	for shift := 0; shift < 32; shift += bitsPerItemNode { // iterative
		bit := hashbit(hash, shift)
		iv := bits.OnesCount32(nd.bmVal & (bit - 1))
		if (nd.bmVal & bit) != 0 {
			if nd.vals[iv].Key() == key {
				return &nd.vals[iv]
			}
		}
		if (nd.bmPtr & bit) == 0 {
			return nil
		}
		ip := bits.OnesCount32(nd.bmPtr & (bit - 1))
		nd = nd.ptrs[ip]
	}
	// overflow node, linear search
	for i := range nd.vals {
		if nd.vals[i].Key() == key {
			return &nd.vals[i]
		}
	}
	return nil // not found
}

func hashbit(hash uint64, shift int) uint32 {
	return 1 << ((hash >> shift) & maskItem)
}

//-------------------------------------------------------------------

func (ht Hamt[K, E]) Mutable() Hamt[K, E] {
	gen := ht.generation + 1
	nd := ht.root
	if nd == nil {
		nd = &node[K, E]{generation: gen}
	} else {
		nd = nd.dup()
	}
	nd.generation = gen
	return Hamt[K, E]{root: nd, mutable: true, generation: gen}
}

func (ht Hamt[K, E]) Put(item E) {
	if !ht.mutable {
		panic("can't modify an immutable Hamt")
	}
	key := item.Key()
	hash := item.Hash(key)
	ht.root.with(ht.generation, item, key, hash, 0)
}

func (nd *node[K, E]) with(gen uint32, item E, key K, hash uint64, shift int) *node[K, E] {
	// recursive
	if nd.generation != gen {
		// path copy on the way down the tree
		nd = nd.dup()
		nd.generation = gen // now mutable in this generation
	}
	if shift >= 32 {
		// overflow node
		for i := range nd.vals { // linear search
			if nd.vals[i].Key() == key {
				nd.vals[i] = item // update if found
				return nd
			}
		}
		nd.vals = append(nd.vals, item) // not found, add it
		return nd
	}
	bit := hashbit(hash, shift)
	iv := bits.OnesCount32(nd.bmVal & (bit - 1))
	if (nd.bmVal & bit) == 0 {
		// slot is empty, insert new value
		nd.bmVal |= bit
		var zero E
		nd.vals = append(nd.vals, zero)
		copy(nd.vals[iv+1:], nd.vals[iv:])
		nd.vals[iv] = item
		return nd
	}
	if nd.vals[iv].Key() == key {
		// already exists, update it
		nd.vals[iv] = item
		return nd
	}

	ip := bits.OnesCount32(nd.bmPtr & (bit - 1))
	if (nd.bmPtr & bit) != 0 {
		// recurse to child node
		nd.ptrs[ip] = nd.ptrs[ip].with(gen, item, key, hash, shift+bitsPerItemNode)
		return nd
	}
	// collision, push new value down to new child node
	child := &node[K, E]{generation: gen}
	child = child.with(gen, item, key, hash, shift+bitsPerItemNode)

	// point to new child node
	nd.ptrs = append(nd.ptrs, nil)
	copy(nd.ptrs[ip+1:], nd.ptrs[ip:])
	nd.ptrs[ip] = child
	nd.bmPtr |= bit

	return nd
}

func (nd *node[K, E]) dup() *node[K, E] {
	dup := *nd // shallow copy
	dup.vals = append(nd.vals[0:0:0], nd.vals...)
	dup.ptrs = append(nd.ptrs[0:0:0], nd.ptrs...)
	return &dup
}

func (ht Hamt[K, E]) Freeze() Hamt[K, E] {
	return Hamt[K, E]{root: ht.root, generation: ht.generation}
}

//-------------------------------------------------------------------

// Delete removes an item. It returns whether the item was found.
func (ht Hamt[K, E]) Delete(key K) bool {
	if !ht.mutable {
		panic("can't modify an immutable Hamt[string,E]")
	}
	var z E
	hash := z.Hash(key)
	_, ok := ht.root.without(ht.generation, key, hash, 0)
	return ok
}

func (nd *node[K, E]) without(gen uint32, key K, hash uint64, shift int) (*node[K, E], bool) {
	// recursive
	if nd.generation != gen {
		// path copy on the way down the tree
		nd = nd.dup()
		nd.generation = gen // now mutable in this generation
	}
	if shift >= 32 {
		// overflow node
		for i := range nd.vals { // linear search
			if nd.vals[i].Key() == key {
				nd.vals[i] = nd.vals[len(nd.vals)-1]
				nd.vals = nd.vals[:len(nd.vals)-1]
				if len(nd.vals) == 0 { // node emptied
					nd = nil
				}
				return nd, true
			}
		}
		return nd, false
	}
	bit := hashbit(hash, shift)
	iv := bits.OnesCount32(nd.bmVal & (bit - 1))
	if (nd.bmVal & bit) != 0 {
		if nd.vals[iv].Key() == key {
			// found it
			if (nd.bmPtr & bit) == 0 { // no child
				nd.bmVal &^= bit
				nd.vals = append(nd.vals[:iv], nd.vals[iv+1:]...) // preserve order
				if nd.bmVal == 0 && nd.bmPtr == 0 {               // node emptied
					nd = nil
				}
			} else {
				// pull up child value
				ip := bits.OnesCount32(nd.bmPtr & (bit - 1))
				child, item := nd.ptrs[ip].pullUp(gen)
				nd.vals[iv] = item // replace the item
				if child != nil {
					nd.ptrs[ip] = child
				} else { // child emptied
					nd.bmPtr &^= bit
					nd.ptrs = append(nd.ptrs[:ip], nd.ptrs[ip+1:]...) // preserve order
				}
			}
			return nd, true
		}
	}
	if (nd.bmPtr & bit) == 0 {
		return nd, false
	}
	ip := bits.OnesCount32(nd.bmPtr & (bit - 1))
	child, ok := nd.ptrs[ip].without(gen, key, hash, shift+bitsPerItemNode) // RECURSE
	if child != nil {
		nd.ptrs[ip] = child
	} else { // child emptied
		nd.bmPtr &^= bit
		nd.ptrs = append(nd.ptrs[:ip], nd.ptrs[ip+1:]...) // preserve order
	}
	return nd, ok
}

func (nd *node[K, E]) pullUp(gen uint32) (*node[K, E], E) {
	// recursive
	if nd.generation != gen {
		// path copy on the way down the tree
		nd = nd.dup()
		nd.generation = gen // now mutable in this generation
	}
	if nd.bmPtr != 0 { // have children
		assert.That(nd.bmVal != 0)
		ip := len(nd.ptrs) - 1
		child, item := nd.ptrs[ip].pullUp(gen) // RECURSE
		if child != nil {
			nd.ptrs[ip] = child
		} else {
			nd.ptrs = nd.ptrs[:ip] // drop empty child node
			// clear highest one bit
			nd.bmPtr = clearHighestOneBit(nd.bmPtr)
		}
		return nd, item
	}
	// no children
	iv := len(nd.vals) - 1
	item := nd.vals[iv]
	if iv == 0 { // last value in node
		return nil, item
	}
	nd.vals = nd.vals[:iv]
	if nd.bmVal != 0 { // not an overflow node
		// clear highest one bit
		nd.bmVal = clearHighestOneBit(nd.bmVal)
	}
	return nd, item
}

func clearHighestOneBit(n uint32) uint32 {
	return n &^ (1 << (31 - bits.LeadingZeros32(n)))
}

//-------------------------------------------------------------------

func (ht Hamt[K, E]) All() iter.Seq[E] {
	return func(yield func(E) bool) {
		if ht.root != nil {
			ht.root.forEach(yield)
		}
	}
}

func (nd *node[K, E]) forEach(fn func(E) bool) bool {
	for i := range nd.vals {
		if !fn(nd.vals[i]) {
			return false
		}
	}
	for _, p := range nd.ptrs {
		if !p.forEach(fn) {
			return false
		}
	}
	return true
}

//-------------------------------------------------------------------

type Chain[K comparable, E Item[K]] struct {
	Hamt[K, E]
	// offs are the offsets in the database file
	// of the item chunks in the current chain, oldest first.
	// These are used for "merging" chunks to manage chain size.
	Offs []uint64
	// ages are the oldest/min ages in the chunks.
	// They are parallel to offs (same len).
	Ages []int
	// clock counts persists.
	// lastMod is set to the current clock to mark an item as modified.
	Clock int
}

const All = math.MinInt

// WriteChain writes a new chunk of items
// containing at least the newly modified items
// plus the contents of zero or more older chunks.
// It returns a new ItemChain. It does not modify the original.
// Conceptually we merge chunks,
// but actually we write a new chunk containing old and new items
// and unlink/abandon the old chunk(s).
func (c *Chain[K, E]) WriteChain(store *stor.Stor) (uint64, Chain[K, E]) {
	assert.That(!c.mutable)
	no := len(c.Offs)
	merge := nmerge(no, c.Clock)
	oldest := c.Clock
	if merge > 0 {
		oldest = c.Ages[no-merge]
	}
	lastMod := oldest
	if merge == no {
		lastMod = All
	}
	prevOff := uint64(0)
	if no > 0 && merge < no {
		prevOff = c.Offs[no-merge-1]
	}
	off := c.Write(store, prevOff, lastMod)
	if off == 0 {
		if no > 0 {
			off = c.Offs[no-1] // nothing written, return current chain
		}
		return off, *c
	}
	n := no - merge
	c2 := Chain[K, E]{
		Hamt:  c.Hamt,
		Clock: c.Clock + 1,
		Offs:  slc.With(c.Offs[:n:n], off),
		Ages:  slc.With(c.Ages[:n:n], oldest),
	}
	return off, c2
}

const maxChain = 7 // ???

func nmerge(no, clock int) int {
	if no >= maxChain {
		return no
	}
	return min(no, TrailingOnes(clock))
}

func (ht Hamt[K, E]) Write(st *stor.Stor, prevOff uint64, lastMod int) uint64 {
	size := 0
	ck := uint32(0)
	for it := range ht.All() {
		if it.IsTomb() {
			if lastMod == All {
				continue
			}
		} else {
			ck += it.Cksum()
		}
		if it.LastMod() >= lastMod {
			size += it.StorSize()
		}
	}
	if size == 0 && (prevOff == 0 || lastMod != All) {
		return 0
	}
	size += 3 + 5 + cksum.Len + 4
	off, buf := st.Alloc(size)
	w := stor.NewWriter(buf)
	w.Put3(size)
	w.Put5(prevOff)
	w.Put4(int(ck))
	for it := range ht.All() {
		if lastMod != All || !it.IsTomb() {
			if it.LastMod() >= lastMod {
				it.Write(w)
			}
		}
	}
	assert.That(w.Len() == size-cksum.Len)
	cksum.Update(buf)
	return off
}

//-------------------------------------------------------------------

func ReadChain[K comparable, E Item[K]](st *stor.Stor, off uint64,
	rdfn func(st *stor.Stor, r *stor.Reader) E) Chain[K, E] {
	if off == 0 {
		return Chain[K, E]{}
	}
	offs := make([]uint64, 0, 8)
	ht := Hamt[K, E]{}.Mutable()
	offs = append(offs, off)
	lastMod := -1
	// the checksum of the most recent chunk is the checksum of the Hamt
	next, ck := ht.read(st, off, lastMod, rdfn)
	for lastMod--; next != 0; lastMod-- {
		offs = append(offs, next)
		next, _ = ht.read(st, next, lastMod, rdfn)
	}
	ck2 := uint32(0)
	for it := range ht.All() {
		if !it.IsTomb() {
			ck2 += it.Cksum()
		}
	}
	if ck != ck2 {
		panic("metadata checksum mismatch")
	}
	ages := make([]int, len(offs))
	for i := range ages {
		lastMod++
		ages[i] = lastMod
	}
	return Chain[K, E]{
		Hamt: ht.Freeze(),
		Offs: slc.Reverse(offs),
		Ages: ages,
	}
}

func (ht Hamt[K, E]) read(st *stor.Stor, off uint64,
	lastMod int, rdfn func(st *stor.Stor, r *stor.Reader) E) (uint64, uint32) {
	initial := ht.IsNil() // optimization
	buf := st.Data(off)
	size := stor.NewReader(buf).Get3()
	cksum.MustCheck(buf[:size])
	r := stor.NewReader(buf[3 : size-cksum.Len])
	prevOff := r.Get5()
	ck := uint32(r.Get4())
	for r.Remaining() > 0 {
		it := rdfn(st, r)
		// reading newest first, so ignore older versions
		if initial || ht.get(it.Key()) == nil {
			it.SetLastMod(lastMod)
			ht.Put(it)
		}
	}
	return prevOff, ck
}

//-------------------------------------------------------------------

func (c *Chain[K, E]) Cksum() uint32 {
	cksum := uint32(c.Clock)
	for i := range c.Offs {
		cksum += uint32(c.Offs[i]) + uint32(c.Ages[i])
	}
	for it := range c.All() {
		cksum += it.Cksum()
	}
	return cksum
}

// Cksum on Hamt is for the logical state, whereas Cksum on Chain is physical.
func (ht Hamt[K, E]) Cksum() uint32 {
	cksum := uint32(0)
	for it := range ht.All() {
		if !it.IsTomb() {
			cksum += it.Cksum()
		}
	}
	return cksum
}
