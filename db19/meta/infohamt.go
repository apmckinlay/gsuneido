// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package meta

import (
	"math/bits"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/cksum"
)

type InfoHamt struct {
	root       *nodeInfo
	mutable    bool
	generation uint32 // if mutable, nodes with this generation are mutable
}

type nodeInfo struct {
	generation uint32
	bmVal      uint32
	bmPtr      uint32
	vals       []*Info
	ptrs       []*nodeInfo
}

const bitsPerInfoNode = 5
const maskInfo = 1<<bitsPerInfoNode - 1

func (ht InfoHamt) IsNil() bool {
	return ht.root == nil
}

func (ht InfoHamt) Get(key string) (*Info, bool) {
	it := ht.get(key)
	if it == nil {
		var zero *Info
		return zero, false
	}
	return *it, true
}

func (ht InfoHamt) get(key string) **Info {
	nd := ht.root
	if nd == nil {
		return nil
	}
	hash := InfoHash(key)
	for shift := 0; shift < 32; shift += bitsPerInfoNode { // iterative
		bit := nd.bit(hash, shift)
		iv := bits.OnesCount32(nd.bmVal & (bit - 1))
		if (nd.bmVal & bit) != 0 {
			if InfoKey(nd.vals[iv]) == key {
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
		if InfoKey(nd.vals[i]) == key {
			return &nd.vals[i]
		}
	}
	return nil // not found
}

func (*nodeInfo) bit(hash uint32, shift int) uint32 {
	return 1 << ((hash >> shift) & maskInfo)
}

//-------------------------------------------------------------------

func (ht InfoHamt) Mutable() InfoHamt {
	gen := ht.generation + 1
	nd := ht.root
	if nd == nil {
		nd = &nodeInfo{generation: gen}
	}
	nd = nd.dup()
	nd.generation = gen
	return InfoHamt{root: nd, mutable: true, generation: gen}
}

func (ht InfoHamt) Put(item *Info) {
	if !ht.mutable {
		panic("can't modify an immutable Hamt")
	}
	key := InfoKey(item)
	hash := InfoHash(key)
	ht.root.with(ht.generation, item, key, hash, 0)
}

func (nd *nodeInfo) with(gen uint32, item *Info, key string, hash uint32, shift int) *nodeInfo {
	// recursive
	if nd.generation != gen {
		// path copy on the way down the tree
		nd = nd.dup()
		nd.generation = gen // now mutable in this generation
	}
	if shift >= 32 {
		// overflow node
		for i := range nd.vals { // linear search
			if InfoKey(nd.vals[i]) == key {
				nd.vals[i] = item // update if found
				return nd
			}
		}
		nd.vals = append(nd.vals, item) // not found, add it
		return nd
	}
	bit := nd.bit(hash, shift)
	iv := bits.OnesCount32(nd.bmVal & (bit - 1))
	if (nd.bmVal & bit) == 0 {
		// slot is empty, insert new value
		nd.bmVal |= bit
		var zero *Info
		nd.vals = append(nd.vals, zero)
		copy(nd.vals[iv+1:], nd.vals[iv:])
		nd.vals[iv] = item
		return nd
	}
	if InfoKey(nd.vals[iv]) == key {
		// already exists, update it
		nd.vals[iv] = item
		return nd
	}

	ip := bits.OnesCount32(nd.bmPtr & (bit - 1))
	if (nd.bmPtr & bit) != 0 {
		// recurse to child node
		nd.ptrs[ip] = nd.ptrs[ip].with(gen, item, key, hash, shift+bitsPerInfoNode)
		return nd
	}
	// collision, push new value down to new child node
	child := &nodeInfo{generation: gen}
	child = child.with(gen, item, key, hash, shift+bitsPerInfoNode)

	// point to new child node
	nd.ptrs = append(nd.ptrs, nil)
	copy(nd.ptrs[ip+1:], nd.ptrs[ip:])
	nd.ptrs[ip] = child
	nd.bmPtr |= bit

	return nd
}

func (nd *nodeInfo) dup() *nodeInfo {
	dup := *nd // shallow copy
	dup.vals = append(nd.vals[0:0:0], nd.vals...)
	dup.ptrs = append(nd.ptrs[0:0:0], nd.ptrs...)
	return &dup
}

func (ht InfoHamt) Freeze() InfoHamt {
	return InfoHamt{root: ht.root, generation: ht.generation}
}

//-------------------------------------------------------------------

// Delete removes an item. It returns whether the item was found.
func (ht InfoHamt) Delete(key string) bool {
	if !ht.mutable {
		panic("can't modify an immutable Hamt")
	}
	hash := InfoHash(key)
	_, ok := ht.root.without(ht.generation, key, hash, 0)
	return ok
}

func (nd *nodeInfo) without(gen uint32, key string, hash uint32, shift int) (*nodeInfo, bool) {
	// recursive
	if nd.generation != gen {
		// path copy on the way down the tree
		nd = nd.dup()
		nd.generation = gen // now mutable in this generation
	}
	if shift >= 32 {
		// overflow node
		for i := range nd.vals { // linear search
			if InfoKey(nd.vals[i]) == key {
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
	bit := nd.bit(hash, shift)
	iv := bits.OnesCount32(nd.bmVal & (bit - 1))
	if (nd.bmVal & bit) != 0 {
		if InfoKey(nd.vals[iv]) == key {
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
	child, ok := nd.ptrs[ip].without(gen, key, hash, shift+bitsPerInfoNode) // RECURSE
	if child != nil {
		nd.ptrs[ip] = child
	} else { // child emptied
		nd.bmPtr &^= bit
		nd.ptrs = append(nd.ptrs[:ip], nd.ptrs[ip+1:]...) // preserve order
	}
	return nd, ok
}

func (nd *nodeInfo) pullUp(gen uint32) (*nodeInfo, *Info) {
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
			nd.bmPtr = nd.clearHighestOneBit(nd.bmPtr)
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
		nd.bmVal = nd.clearHighestOneBit(nd.bmVal)
	}
	return nd, item
}

func (*nodeInfo) clearHighestOneBit(n uint32) uint32 {
	return n &^ (1 << (31 - bits.LeadingZeros32(n)))
}

//-------------------------------------------------------------------

func (ht InfoHamt) ForEach(fn func(*Info)) {
	if ht.root != nil {
		ht.root.forEach(fn)
	}
}

func (nd *nodeInfo) forEach(fn func(*Info)) {
	for i := range nd.vals {
		fn(nd.vals[i])
	}
	for _, p := range nd.ptrs {
		p.forEach(fn)
	}
}

//-------------------------------------------------------------------

func (ht InfoHamt) Write(st *stor.Stor, prevOff uint64,
	filter func(it *Info) bool) uint64 {
	size := 0
	ck := uint32(0)
	ht.ForEach(func(it *Info) {
		if filter(it) {
			size += it.storSize()
		}
		if !it.isTomb() {
			ck += it.Cksum()
		}
	})
	if size == 0 {
		return 0
	}
	size += 3 + 5 + cksum.Len + 4
	off, buf := st.Alloc(size)
	w := stor.NewWriter(buf)
	w.Put3(size)
	w.Put5(prevOff)
	w.Put4(int(ck))
	ht.ForEach(func(it *Info) {
		if filter(it) {
			it.Write(w)
		}
	})
	assert.That(w.Len() == size-cksum.Len)
	cksum.Update(buf)
	return off
}

func ReadInfoChain(st *stor.Stor, off uint64) (InfoHamt, []uint64) {
	offs := make([]uint64, 0, 8)
	ht := InfoHamt{}.Mutable()
	tomb := make(map[string]struct{}, 16)
	offs = append(offs, off)
	var ck uint32
	off, ck = ht.read(st, off, tomb)
	for off != 0 {
		offs = append(offs, off)
		off, _ = ht.read(st, off, tomb)
	}
	ck2 := uint32(0)
	ht.ForEach(func(it *Info) {
		ck2 += it.Cksum()
	})
	if ck != ck2 {
		panic("metadata checksum mismatch")
	}
	return ht.Freeze(), offs
}

func (ht InfoHamt) read(st *stor.Stor, off uint64, tomb map[string]struct{}) (uint64, uint32) {
	initial := ht.IsNil() // optimization
	buf := st.Data(off)
	size := stor.NewReader(buf).Get3()
	cksum.MustCheck(buf[:size])
	r := stor.NewReader(buf[3 : size-cksum.Len])
	prevOff := r.Get5()
	ck := uint32(r.Get4())
	for r.Remaining() > 0 {
		it := ReadInfo(st, r)
		if initial || ht.get(InfoKey(it)) == nil {
			// doesn't exist yet
			key := InfoKey(it)
			if it.isTomb() {
				tomb[key] = struct{}{}
			} else if _, ok := tomb[key]; !ok {
				ht.Put(it)
			}
		}
	}
	return prevOff, ck
}
