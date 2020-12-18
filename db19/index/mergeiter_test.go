// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package index

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestMergeIter(t *testing.T) {
	assert := assert.T(t)
	from := func(args ...int) *ixbuf.T {
		ib := &ixbuf.T{}
		for _, n := range args {
			ib.Insert(strconv.Itoa(n), uint64(n))
		}
		return ib
	}
	even := from(0, 2, 4, 6, 8)
	odd := from(1, 3, 5, 7, 9)
	modCount := 0
	callback := func(mc int) (int, []iterator) {
		if mc == modCount {
			return mc, nil
		}
		return modCount, []iterator{even.Iterator(), odd.Iterator()}
	}
	it := NewMergeIter(callback)
	test := func(expected int) {
		t.Helper()
		if expected == -1 {
			assert.That(it.Eof())
		} else {
			key, off := it.Cur()
			assert.This(key).Is(strconv.Itoa(expected))
			assert.This(off).Is(uint64(expected))
		}
	}
	testNext := func(expected int) { it.Next(); t.Helper(); test(expected) }
	testPrev := func(expected int) { it.Prev(); t.Helper(); test(expected) }
	for i := 0; i < 10; i++ {
		testNext(i)
		if i == 5 {
			modCount++
		}
	}
	testNext(-1)

	it.Rewind()
	for i := 9; i >= 0; i-- {
		testPrev(i)
		if i == 5 {
			modCount++
		}
	}
	testPrev(-1)

	it.Rewind()
	testNext(0)
	testPrev(-1) // stick at eof
	testPrev(-1)
	testNext(-1)

	it.Rewind()
	testPrev(9)
	testPrev(8)
	testPrev(7)
	testNext(8)
	testNext(9) // last
	testPrev(8)

	it.Rewind()
	testNext(0)
	testNext(1)
	even.Insert("11", 11)
	testNext(11)
	testNext(2)
	even.Delete("11", 11)
	testPrev(1) // modified AND changed direction

	it.Rewind()
	testPrev(9)
	testPrev(8)
	odd.Insert("77", 77)
	testPrev(77)
	modCount++
	testPrev(7)
}

func TestMergeIterRandom(*testing.T) {
	trace := func(args ...interface{}) {
		// fmt.Print(args...)
	}
	traceln := func(args ...interface{}) {
		// fmt.Println(args...)
	}
	type val struct {
		layer int
		key   string
	}
	gen := make(dat)
	data := new(dummy)
	it := &dumIter{d: data}
	which := map[string]int{}
	ibs := [3]ixbuf.T{}
	modCount := 0
	callback := func(mc int) (int, []iterator) {
		if mc == modCount {
			return mc, nil
		}
		return modCount, []iterator{
			ibs[0].Iterator(), ibs[1].Iterator(), ibs[2].Iterator()}
	}
	mi := NewMergeIter(callback)
	check := func() {
		if it.Eof() {
			traceln("EOF")
			assert.That(mi.Eof())
			it.Rewind()
			mi.Rewind()
		} else {
			actual, _ := mi.Cur()
			traceln(actual)
			assert.This(actual).Is(it.Cur())
		}
	}
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("+ + +")
			mi.Rewind()
			for mi.Next(); !mi.Eof(); mi.Next() {
				key, _ := mi.Cur()
				traceln(mi.curIter, key)
			}
			fmt.Println("+ + +")
			it.Rewind()
			for it.Next(); !it.Eof(); it.Next() {
				key := it.Cur()
				traceln(key)
			}
			fmt.Println("+ + +")
			panic(e)
		}
	}()
	const sizing = 15
	const N = 1_000_000
	for n := 0; n < N; n++ {
		trace(n, " ")
		switch rand.Intn(7) {
		case 0:
			if rand.Intn(sizing)+1 > data.Len() {
				// insert
				key := gen.randKey()
				data.Insert(key)
				w := rand.Intn(len(ibs))
				ibs[w].Insert(key, 1)
				which[key] = w
				traceln("insert", w, key, "len", data.Len())
			} else {
				// delete
				j := rand.Intn(data.Len())
				key := data.Get(j)
				w := which[key]
				traceln("delete", w, key, "len", data.Len()-1)
				data.Delete(key)
				gen.delete(key)
				ibs[w].Delete(key, 1)
			}
		case 1: // rewind
			traceln("rewind")
			it.Rewind()
			mi.Rewind()
		case 2:
			traceln("invalidate")
			modCount++
		case 3, 4: // next
			it.Next()
			mi.Next()
			trace("next ")
			check()
		case 5, 6: // prev
			it.Prev()
			mi.Prev()
			trace("prev ")
			check()
		}
		size := 0
		for _, ib := range ibs {
			size += ib.Len()
		}
		assert.This(size).Is(data.Len())
	}
}

//-------------------------------------------------------------------

type dat map[string]struct{}

func (d dat) randKey() string {
	for {
		key := str.Random(3, 3)
		if _, ok := d[key]; !ok {
			d[key] = struct{}{}
			return key
		}
	}
}

func (d dat) delete(key string) {
	delete(d, key)
}

//-------------------------------------------------------------------

// dummy and dumIter are a very simple, and therefore hopefully correct,
// version of an ordered iterable container.
// They are used by TestMergeIterRandom to verify the behavior of MergeIter.
type dummy struct {
	keys []string
}

func (d *dummy) Len() int {
	return len(d.keys)
}

func (d *dummy) Get(i int) string {
	return d.keys[i]
}

func (d *dummy) Insert(key string) {
	for i, k := range d.keys {
		if key < k {
			d.keys = append(d.keys, "")
			copy(d.keys[i+1:], d.keys[i:])
			d.keys[i] = key
			return
		}
	}
	d.keys = append(d.keys, key)
}

func (d *dummy) Delete(key string) {
	for i, k := range d.keys {
		if key == k {
			copy(d.keys[i:], d.keys[i+1:])
			d.keys = d.keys[:len(d.keys)-1]
			return
		}
	}
	panic("key not found")
}

func (d *dummy) print() {
	fmt.Println("+ + +")
	for _, k := range d.keys {
		fmt.Println(k)
	}
	fmt.Println("+ + +")
}

type dumIter struct {
	d   *dummy
	cur string
	state
}

func (it *dumIter) Rewind() {
	it.state = rewound
}

func (it *dumIter) Eof() bool {
	return it.state == eof
}

func (it *dumIter) Cur() string {
	return it.cur
}

func (it *dumIter) Next() {
	if it.state == eof {
		return
	}
	if it.state == rewound {
		if len(it.d.keys) == 0 {
			it.state = eof
			return
		}
		it.cur = it.d.keys[0]
		it.state = within
		return
	}
	for _, k := range it.d.keys {
		if k > it.cur {
			it.cur = k
			return
		}
	}
	it.state = eof
}

func (it *dumIter) Prev() {
	if it.state == eof {
		return
	}
	if it.state == rewound {
		if len(it.d.keys) == 0 {
			it.state = eof
			return
		}
		it.cur = it.d.keys[len(it.d.keys)-1]
		it.state = within
		return
	}
	for i := len(it.d.keys) - 1; i >= 0; i-- {
		k := it.d.keys[i]
		if k < it.cur {
			it.cur = k
			return
		}
	}
	it.state = eof
}
