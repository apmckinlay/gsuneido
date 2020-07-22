// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package hamt

import (
	"math/bits"

	"github.com/cheekybits/genny/generic"
)

type Item generic.Type
type KeyType generic.Type

type ItemHamt struct {
	root       *nodeItem
	mutable    bool
	generation uint32 // if mutable, nodes with this generation are mutable
}

type nodeItem struct {
	generation uint32
	bmVal      uint32
	bmPtr      uint32
	vals       []Item
	ptrs       []*nodeItem
}

const bitsPerItemNode = 5
const maskItem = 1<<bitsPerItemNode - 1

func (ht ItemHamt) IsNil() bool {
	return ht.root == nil
}

func (ht ItemHamt) MustGet(key KeyType) Item {
	it, ok := ht.Get(key)
	if !ok {
		panic("Hamt MustGet failed")
	}
	return it
}

func (ht ItemHamt) Get(key KeyType) (Item, bool) {
	it := ht.get(key)
	if it == nil {
		var zero Item
		return zero, false
	}
	return *it, true
}

func (ht ItemHamt) get(key KeyType) *Item {
	nd := ht.root
	if nd == nil {
		return nil
	}
	hash := ItemHash(key)
	for shift := 0; shift < 32; shift += bitsPerItemNode { // iterative
		bit := nd.bit(hash, shift)
		iv := bits.OnesCount32(nd.bmVal & (bit - 1))
		if (nd.bmVal & bit) != 0 {
			if ItemKey(nd.vals[iv]) == key {
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
		if ItemKey(nd.vals[i]) == key {
			return &nd.vals[i]
		}
	}
	return nil // not found
}

func (*nodeItem) bit(hash uint32, shift int) uint32 {
	return 1 << ((hash >> shift) & maskItem)
}

//-------------------------------------------------------------------

func (ht ItemHamt) Mutable() ItemHamt {
	gen := ht.generation + 1
	nd := ht.root
	if nd == nil {
		nd = &nodeItem{generation: gen}
	}
	nd = nd.dup()
	nd.generation = gen
	return ItemHamt{root: nd, mutable: true, generation: gen}
}

func (ht ItemHamt) Put(item Item) {
	if !ht.mutable {
		panic("can't modify an immutable Hamt")
	}
	key := ItemKey(item)
	hash := ItemHash(key)
	ht.root.with(ht.generation, item, key, hash, 0)
}

func (nd *nodeItem) with(gen uint32, item Item, key KeyType, hash uint32, shift int) *nodeItem {
	// recursive
	if nd.generation != gen {
		// path copy on the way down the tree
		nd = nd.dup()
		nd.generation = gen // now mutable in this generation
	}
	if shift >= 32 {
		// overflow node
		for i := range nd.vals { // linear search
			if ItemKey(nd.vals[i]) == key {
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
		var zero Item
		nd.vals = append(nd.vals, zero)
		copy(nd.vals[iv+1:], nd.vals[iv:])
		nd.vals[iv] = item
		return nd
	}
	if ItemKey(nd.vals[iv]) == key {
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
	child := &nodeItem{generation: gen}
	child = child.with(gen, item, key, hash, shift+bitsPerItemNode)

	// point to new child node
	nd.ptrs = append(nd.ptrs, nil)
	copy(nd.ptrs[ip+1:], nd.ptrs[ip:])
	nd.ptrs[ip] = child
	nd.bmPtr |= bit

	return nd
}

func (nd *nodeItem) dup() *nodeItem {
	dup := *nd // shallow copy
	dup.vals = append(nd.vals[0:0:0], nd.vals...)
	dup.ptrs = append(nd.ptrs[0:0:0], nd.ptrs...)
	return &dup
}

func (ht ItemHamt) Freeze() ItemHamt {
	return ItemHamt{root: ht.root, generation: ht.generation}
}

//-------------------------------------------------------------------

// Delete removes an item.
func (ht ItemHamt) Delete(key KeyType) bool {
	if !ht.mutable {
		panic("can't modify an immutable Hamt")
	}
	hash := ItemHash(key)
	_, ok := ht.root.without(ht.generation, key, hash, 0)
	return ok
}

func (nd *nodeItem) without(gen uint32, key KeyType, hash uint32, shift int) (*nodeItem, bool) {
	// recursive
	if nd.generation != gen {
		// path copy on the way down the tree
		nd = nd.dup()
		nd.generation = gen // now mutable in this generation
	}
	if shift >= 32 {
		// overflow node
		for i := range nd.vals { // linear search
			if ItemKey(nd.vals[i]) == key {
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
		if ItemKey(nd.vals[iv]) == key {
			nd.bmVal &^= bit
			nd.vals = append(nd.vals[:iv], nd.vals[iv+1:]...) // preserve order
			if nd.bmVal == 0 && nd.bmPtr == 0 {               // node emptied
				nd = nil
			}
			return nd, true
		}
	}
	if (nd.bmPtr & bit) == 0 {
		return nd, false
	}
	ip := bits.OnesCount32(nd.bmPtr & (bit - 1))
	child, ok := nd.ptrs[ip].without(gen, key, hash, shift+bitsPerItemNode) // recurse
	if child != nil {
		nd.ptrs[ip] = child
	} else { // child emptied
		nd.bmPtr &^= bit
		nd.ptrs = append(nd.ptrs[:ip], nd.ptrs[ip+1:]...) // preserve order
		if nd.bmPtr == 0 && nd.bmVal == 0 {               // this node emptied
			nd = nil
		}
	}
	return nd, ok
}

//-------------------------------------------------------------------

func (ht ItemHamt) ForEach(fn func(Item)) {
	if ht.root != nil {
		ht.root.forEach(fn)
	}
}

func (nd *nodeItem) forEach(fn func(Item)) {
	for i := range nd.vals {
		fn(nd.vals[i])
	}
	for _, p := range nd.ptrs {
		p.forEach(fn)
	}
}
