// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package index

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index/btree"
	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/ranges"
	"github.com/apmckinlay/gsuneido/util/str"
)

type testTran struct {
	getIndex func() *Overlay
	reads    ranges.Ranges
}

func (t *testTran) GetIndexI(string, int) *Overlay {
	return t.getIndex()
}

func (t *testTran) Read(_ string, _ int, from, to string) {
	t.reads.Insert(from, to)
}

func TestOverIter(t *testing.T) {
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
	bt := btree.CreateBtree(stor.HeapStor(8192), nil)
	it := NewOverIter("", 0)
	test := func(expected int) {
		t.Helper()
		if expected == -1 {
			assert.Msg("expected Eof").That(it.Eof())
		} else {
			assert.That(!it.Eof())
			key, off := it.Cur()
			assert.This(key).Is(strconv.Itoa(expected))
			assert.This(off).Is(uint64(expected))
		}
	}
	newTran := func() *testTran {
		return &testTran{getIndex: func() *Overlay {
			return &Overlay{bt: bt, layers: []*ixbuf.T{even, odd}}
		}}
	}
	tran := newTran()
	testNext := func(expected int) { it.Next(tran); t.Helper(); test(expected) }
	testPrev := func(expected int) { it.Prev(tran); t.Helper(); test(expected) }
	for i := 0; i < 10; i++ {
		testNext(i)
		if i == 5 {
			tran = newTran()
		}
	}
	testNext(-1)

	it.Rewind()
	for i := 9; i >= 0; i-- {
		testPrev(i)
		if i == 5 {
			tran = newTran()
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
	tran = newTran()
	testNext(11)
	testNext(2)
	even.Delete("11", 11)
	tran = newTran()
	testPrev(1) // modified AND changed direction

	it.Rewind()
	testPrev(9)
	testPrev(8)
	odd.Insert("77", 77)
	tran = newTran()
	testPrev(77)
	testPrev(7)

	it.Range(Range{Org: "3", End: "6"})
	testNext(3)
	testNext(4)
	testPrev(3)
	testPrev(-1)
	it.Rewind()
	testPrev(5)
	testPrev(4)
	testNext(5)
	testNext(-1)
}

func TestOverIterDeletePrevBug(*testing.T) {
	btree.GetLeafKey = func(_ *stor.Stor, _ *ixkey.Spec, i uint64) string {
		return strconv.Itoa(int(i))
	}
	bldr := btree.Builder(stor.HeapStor(8192))
	for i := 1; i <= 9; i++ {
		bldr.Add(strconv.Itoa(i), uint64(i))
	}
	bt := bldr.Finish()
	ib := &ixbuf.T{}
	ov := &Overlay{bt: bt, layers: []*ixbuf.T{}, mut: ib}
	t := &testTran{getIndex: func() *Overlay { return ov }}

	// it := NewOverIter("", 0)
	// it.Prev(t)
	// assert.That(!it.Eof())
	// key, off := it.Cur()
	// assert.This(key).Is("9")
	// assert.This(off).Is(9)

	ov.Delete("9", uint64(9))
	it := NewOverIter("", 0)
	it.Prev(t)
	assert.That(!it.Eof())
	key, off := it.Cur()
	assert.This(key).Is("8")
	assert.This(off).Is(8)
}

func TestOverIterReads(*testing.T) {
	bt := btree.CreateBtree(stor.HeapStor(8192), nil)
	ib := &ixbuf.T{}
	for i := 0; i < 10; i++ {
		ib.Insert(strconv.Itoa(i), uint64(i))
	}
	t := &testTran{getIndex: func() *Overlay {
		return &Overlay{bt: bt, layers: []*ixbuf.T{ib}}
	}}
	assert.This(t.reads.String()).Is("")
	it := NewOverIter("", 0)
	it.Next(t)
	assert.This(t.reads.String()).Is("->0")
	it.Next(t)
	assert.This(t.reads.String()).Is("->1")
	it.Prev(t)
	it.Prev(t)
	assert.That(it.Eof())
	assert.This(t.reads.String()).Is("->1")
	it.Rewind()
	it.Prev(t)
	it.Prev(t)
	assert.This(t.reads.String()).Is("->1 8->\xff\xff\xff\xff\xff\xff\xff\xff")
	it.Next(t)
	it.Next(t)
	assert.This(t.reads.String()).Is("->1 8->\xff\xff\xff\xff\xff\xff\xff\xff")
	it.Rewind()
	it.Next(t)
	it.Next(t)
	it.Next(t)
	assert.This(t.reads.String()).Is("->2 8->\xff\xff\xff\xff\xff\xff\xff\xff")

	t.reads = ranges.Ranges{} // reset
	for it.Rewind(); !it.Eof(); it.Next(t) {
	}
	assert.This(t.reads.String()).Is("->\xff\xff\xff\xff\xff\xff\xff\xff")

	t.reads = ranges.Ranges{} // reset
	for it.Rewind(); !it.Eof(); it.Prev(t) {
	}
	assert.This(t.reads.String()).Is("->\xff\xff\xff\xff\xff\xff\xff\xff")
}

func TestOverIterCombine(*testing.T) {
	var data []string
	defer func(mns int) { btree.MaxNodeSize = mns }(btree.MaxNodeSize)
	btree.MaxNodeSize = 64
	bt := btree.CreateBtree(stor.HeapStor(8192), nil)
	mut := &ixbuf.T{}
	u := &ixbuf.T{}
	ov := &Overlay{bt: bt, layers: []*ixbuf.T{u}, mut: mut}
	checkIter(data, ov)

	const n = 100
	randKey := str.UniqueRandom(3, 7)

	data = insert(data, n, randKey, mut)
	checkIterator(data, ov)

	data = insert(data, n, randKey, u)
	checkIterator(data, ov)

	count := len(data)
	assert.This(count).Is(n * 2)

	for i := 0; i < n/2; i++ {
		j := rand.Intn(len(data))
		if data[j] != "" {
			ov.Delete(data[j], key2off(data[j]))
			data[j] = ""
			count--
		}
	}
	count2 := checkIterator(data, ov)
	assert.This(count2).Is(count)
}

func checkIterator(data []string, ov *Overlay) int {
	sort.Strings(data)
	count := 0
	it := NewOverIter("", 0)
	tran := &testTran{getIndex: func() *Overlay { return ov }}
	for _, k := range data {
		if k == "" {
			continue
		}
		it.Next(tran)
		k2, o2 := it.Cur()
		assert.This(k2).Is(k)
		assert.This(o2).Is(key2off(k))
		count++
	}
	it.Next(tran)
	assert.True(it.Eof())
	return count
}

func TestOverIterRandom(*testing.T) {
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
	bt := btree.CreateBtree(stor.HeapStor(8192), nil)
	ibs := []*ixbuf.T{{}, {}, {}}
	ov := &Overlay{bt: bt, layers: ibs}
	tran := &testTran{getIndex: func() *Overlay { return ov }}
	mi := NewOverIter("", 0)
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
			traceln("===== merge")
			mi.Rewind()
			for mi.Next(tran); !mi.Eof(); mi.Next(tran) {
				key, _ := mi.Cur()
				traceln(mi.curIter, key)
			}
			traceln("===== dumb")
			it.Rewind()
			for it.Next(); !it.Eof(); it.Next() {
				key := it.Cur()
				traceln(key)
			}
			traceln("=====")
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
			tran = &testTran{getIndex: func() *Overlay { return ov }}
		case 3, 4: // next
			it.Next()
			mi.Next(tran)
			trace("next ")
			check()
		case 5, 6: // prev
			it.Prev()
			mi.Prev(tran)
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

func TestOverIterDups(*testing.T) {
	bt := btree.CreateBtree(stor.HeapStor(8192), nil)
	u := &ixbuf.T{}
	u.Insert("", 1)
	mut := &ixbuf.T{}
	mut.Insert("", 2|ixbuf.Update)
	ov := &Overlay{bt: bt, layers: []*ixbuf.T{u}, mut: mut}
	tran := &testTran{getIndex: func() *Overlay { return ov }}

	it := NewOverIter("", 0)
	it.Next(tran)
	assert.That(!it.Eof())
	key, off := it.Cur()
	assert.This(key).Is("")
	assert.This(off).Is(2)
	it.Next(tran)
	assert.That(it.Eof())

	it = NewOverIter("", 0)
	it.Prev(tran)
	assert.That(!it.Eof())
	key, off = it.Cur()
	assert.This(key).Is("")
	assert.This(off).Is(2)
	it.Prev(tran)
	assert.That(it.Eof())
}

func TestOverIterBug(*testing.T) {
	bt := btree.CreateBtree(stor.HeapStor(8192), nil)
	layers := []*ixbuf.T{{}, {}, {}, {}}
	layers[0].Insert("z", 1)
	layers[1].Insert("z", 1 | ixbuf.Update)
	layers[2].Insert("a", 2)
	layers[3].Insert("a", 2 | ixbuf.Delete)
	layers[3].Insert("z", 1 | ixbuf.Delete)
	ov := &Overlay{bt: bt, layers: layers}
	tran := &testTran{getIndex: func() *Overlay { return ov }}

	it := NewOverIter("", 0)
	it.Next(tran)
	assert.That(it.Eof())

	layers = []*ixbuf.T{{}, {}, {}, {}}
	layers[0].Insert("a", 1)
	layers[1].Insert("a", 1 | ixbuf.Update)
	layers[2].Insert("z", 2)
	layers[3].Insert("z", 2 | ixbuf.Delete)
	layers[3].Insert("a", 1 | ixbuf.Delete)
	ov = &Overlay{bt: bt, layers: layers}
	tran = &testTran{getIndex: func() *Overlay { return ov }}

	it = NewOverIter("", 0)
	it.Prev(tran)
	assert.That(it.Eof())
}

func TestOverIterBug2(*testing.T) {
	b := btree.Builder(stor.HeapStor(8192))
	b.Add("1111", 1111)
	b.Add("2222", 2222)
	bt := b.Finish()
	layers := []*ixbuf.T{{}}
	layers[0].Insert("1111", 1111 | ixbuf.Delete)
	layers[0].Insert("2222", 2222 | ixbuf.Delete)
	ov := &Overlay{bt: bt, layers: layers}
	btree.GetLeafKey = func(_ *stor.Stor, _ *ixkey.Spec, i uint64) string {
		return strconv.Itoa(int(i))
	}
	tran := &testTran{getIndex: func() *Overlay { return ov }}
	it := NewOverIter("", 0)
	it.Next(tran)
	assert.That(it.Eof())
}

func TestOverIterBug3(*testing.T) {
	b := btree.Builder(stor.HeapStor(8192))
	b.Add("1111", 1111)
	bt := b.Finish()
	layers := []*ixbuf.T{{}}
	layers[0].Insert("1111", 1111 | ixbuf.Delete)
	ov := &Overlay{bt: bt, layers: layers}
	btree.GetLeafKey = func(_ *stor.Stor, _ *ixkey.Spec, i uint64) string {
		return strconv.Itoa(int(i))
	}
	tran := &testTran{getIndex: func() *Overlay { return ov }}
	
	it := NewOverIter("", 0)
	it.Range(Range{Org: "5555", End: "55559999"})
	it.Next(tran)
	assert.That(it.Eof())

	it = NewOverIter("", 0)
	it.Range(Range{Org: "5555", End: "55559999"})
	it.Prev(tran)
	assert.That(it.Eof())
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
// They are used by TestOverIterRandom to verify the behavior of OverIter.
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
		it.state = front
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
		it.state = back
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
