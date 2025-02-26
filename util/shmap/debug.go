// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package shmap

import (
	"fmt"

	"github.com/apmckinlay/gsuneido/util/assert"
	"golang.org/x/exp/constraints"
)

func (m *Map[K, V, H]) summary() {
	fmt.Println("==============================================")
	fmt.Println("MAP count", m.count, "dirlen", len(m.dir), "depth", m.depth)
	for ti, tbl := range m.dir {
		fmt.Print("TABLE ", ti)
		if ti > 0 && m.dir[ti-1] == tbl {
			fmt.Println(" same")
		} else {
			fmt.Println(" ngroups", len(tbl.groups), "growthLeft", tbl.growthLeft,
				"depth", tbl.depth)
		}
	}
}

func (m *Map[K, V, H]) print() {
	fmt.Println("==============================================")
	fmt.Println("MAP count", m.count, "dirlen", len(m.dir), "depth", m.depth)
	for ti, tbl := range m.dir {
		fmt.Print("TABLE ", ti)
		if ti > 0 && m.dir[ti-1] == tbl {
			fmt.Println(" same")
		} else {
			tbl.print()
		}
	}
}

func (tbl *table[K, V]) print() {
	fmt.Println(" ngroups", len(tbl.groups), "growthLeft", tbl.growthLeft,
		"depth", tbl.depth)
	for gi := range tbl.groups {
		fmt.Println("group", gi)
		grp := &tbl.groups[gi]
		ctrls := grp.control
		for i := range groupSize {
			c := uint8(ctrls)
			ctrls >>= 8
			fmt.Printf("%d %2x: ", i, c)
			if c == empty {
				fmt.Println("empty")
			} else if c == deleted {
				fmt.Println("deleted")
			} else {
				fmt.Printf("%v => %v\n", grp.keys[i], grp.vals[i])
			}
		}
	}
}

func (m *Map[K, V, H]) check() {
	for _, tbl := range m.dir {
		entries := 0
		deletes := 0
		hasEmpty := false
		for gi := range tbl.groups {
			grp := &tbl.groups[gi]
			ctrls := grp.control
			for i := range groupSize {
				c := uint8(ctrls)
				ctrls >>= 8
				if c == deleted {
					deletes++
				} else if c == empty {
					hasEmpty = true
				} else {
					entries++
					k := grp.keys[i]
					h := m.help.Hash(k)
					h2 := uint8(h & 0x7f)
					assert.Msg("control").This(c).Is(h2 | 0x80)
					v, ok := m.Get(k)
					assert.Msg("Get", k).That(ok)
					assert.Msg("Get", k).This(v).Is(grp.vals[i])
				}
			}
		}
		assert.Msg("entries").This(entries).Is(tbl.count)
		assert.Msg("has empty").That(hasEmpty)
		growthLeft := len(tbl.groups)*loadFactor - (entries + deletes)
		assert.Msg("growthLeft").This(growthLeft).Is(tbl.growthLeft)
	}
}

// NewMapInt returns a map for integer keys
// NOTE: this is intended for testing, to give a consistent result
func NewMapInt[K constraints.Integer, V any]() *Map[K, V, integer[K]] {
	return &Map[K, V, integer[K]]{help: integer[K]{}}
}

type integer[K constraints.Integer] struct{}


func (m integer[K]) Hash(k K) uint64 {
	x := uint64(k)
    x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
    x = (x ^ (x >> 27)) * 0x94d049bb133111eb
    x = x ^ (x >> 31)
    return x
}
func (m integer[K]) Equal(x, y K) bool {
	return x == y
}
