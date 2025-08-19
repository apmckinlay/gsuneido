// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import "github.com/apmckinlay/gsuneido/util/hacks"

// unode is an uncompressed node.
// btree keeps a unode of the root.
// Prev also converts nodes to unode's since node can only be iterated forward.
// unode's normally come from node.toUnode()
// and contain known values rather than full keys.
type unode []slot

type slot struct {
	key string
	off uint64
}

type unodeIter struct {
	u unode
	i int
}

func (ui *unodeIter) eof() bool {
	return ui.i+1 >= len(ui.u)
}

func (ui *unodeIter) next() bool {
	ui.i++
	return ui.i < len(ui.u)
}

func (ui *unodeIter) prev() bool {
	if ui.i == -1 {
		ui.i = len(ui.u)
	}
	ui.i--
	return ui.i >= 0
}

func (ui *unodeIter) off() uint64 {
	if ui.i < 0 || ui.i >= len(ui.u) {
		return 0
	}
	return ui.u[ui.i].off
}

func (ui *unodeIter) toUnodeIter(*btree) iNodeIter {
	return ui
}

// search returns the offset of the last entry with key <= the given key.
// Returns 0 if there is no such entry.
// It must find the same slot as node.search
func (u unode) search(key string) uint64 {
	if len(u) == 0 {
		return 0
	}
	lo, hi := 0, len(u)-1
	idx := -1
	for lo <= hi {
		mid := (lo + hi) / 2
		if u[mid].key <= key {
			idx = mid
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	if idx >= 0 {
		return u[idx].off
	}
	return 0
}

// seek returns an iterator positioned at the last entry with key <= the given key.
// If there is no such entry, the iterator index will be -1.
// It must find the same slot as node.seek
func (u unode) seek(key string) *unodeIter {
	lo, hi := 0, len(u)-1
	idx := -1
	for lo <= hi {
		mid := (lo + hi) / 2
		if u[mid].key <= key {
			idx = mid
			lo = mid + 1
		} else {
			hi = mid - 1
		}
	}
	return &unodeIter{u: u, i: idx}
}

// toUnode converts a node to a unode
func (nd node) toUnode() unode {
	nslots := 0
	keylen := 0
	it := nd.iter()
	for it.next() {
		nslots++
		keylen += len(it.known)
	}
	u := make(unode, nslots)
	keys := make([]byte, keylen)
	it.rewind()
	j := 0
	for i := 0; it.next(); i++ {
		n := copy(keys[j:], it.known)
		u[i] = slot{key: hacks.BStoS(keys[j : j+n]), off: it.offset}
		j += n
	}
	return u
}
