// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package index

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"testing"

	btree "github.com/apmckinlay/gsuneido/db19/index/btree3"
	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
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
	from := func(args ...int) *ixbuf.T {
		ib := &ixbuf.T{}
		for _, n := range args {
			ib.Insert(strconv.Itoa(n), uint64(n))
		}
		return ib
	}
	even := from(2, 4, 6, 8)
	odd := from(1, 3, 5, 7, 9)
	bt := btree.CreateBtree(stor.HeapStor(8192), nil)
	it := NewOverIter("", 0)
	test := func(expected int) {
		t.Helper()
		if expected == -1 {
			if !it.Eof() {
				key, off := it.Cur()
				panic(fmt.Sprintln("expected Eof, got", key, off))
			}
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
	for i := 1; i < 10; i++ {
		testNext(i)
		if i == 5 {
			tran = newTran()
		}
	}
	testNext(-1)

	it.Rewind()
	for i := 9; i > 0; i-- {
		testPrev(i)
		if i == 5 {
			tran = newTran()
		}
	}
	testPrev(-1)

	it.Rewind()
	testNext(1)
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
	bldr := btree.Builder(stor.HeapStor(8192))
	for i := 1; i <= 9; i++ {
		assert.That(bldr.Add(strconv.Itoa(i), uint64(i)))
	}
	bt := bldr.Finish()
	ib := &ixbuf.T{}
	ov := &Overlay{bt: bt, layers: []*ixbuf.T{}, mut: ib}
	t := &testTran{getIndex: func() *Overlay { return ov }}

	ov.Delete("9", 9)
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
	for i := 1; i < 10; i++ {
		ib.Insert(strconv.Itoa(i), uint64(i))
	}
	t := &testTran{getIndex: func() *Overlay {
		return &Overlay{bt: bt, layers: []*ixbuf.T{ib}}
	}}

	it := NewOverIter("", 0)

	// Test incremental read tracking
	assert.This(t.reads.String()).Is("")

	// First Next() reads from range start to "1"
	it.Next(t)
	assert.This(t.reads.String()).Is("->1")

	// Second Next() reads from "1" to "2", merges to "->2"
	it.Next(t)
	assert.This(t.reads.String()).Is("->2")

	// Prev() reads from "2" to "1", already covered by existing range
	it.Prev(t)
	assert.This(t.reads.String()).Is("->2")

	// Prev() reads from "1" to range start, already covered
	it.Prev(t)
	assert.That(it.Eof())
	assert.This(t.reads.String()).Is("->2")

	// Reset and test Prev from rewind
	t.reads = ranges.Ranges{}
	it.Rewind()

	// First Prev() from rewind reads from "9" to range end
	it.Prev(t)
	assert.This(t.reads.String()).Is("9->\xff\xff\xff\xff\xff\xff\xff\xff")

	// Second Prev() reads from "8" to "9", merges to "8->end"
	it.Prev(t)
	assert.This(t.reads.String()).Is("8->\xff\xff\xff\xff\xff\xff\xff\xff")

	// Test full forward iteration - should read entire range incrementally
	t.reads = ranges.Ranges{}
	for it.Rewind(); !it.Eof(); it.Next(t) {
	}
	assert.This(t.reads.String()).Is("->\xff\xff\xff\xff\xff\xff\xff\xff")

	// Test full backward iteration - should read entire range incrementally
	t.reads = ranges.Ranges{}
	for it.Rewind(); !it.Eof(); it.Prev(t) {
	}
	assert.This(t.reads.String()).Is("->\xff\xff\xff\xff\xff\xff\xff\xff")
}

func TestOverIterReadsWithRange(*testing.T) {
	bt := btree.CreateBtree(stor.HeapStor(8192), nil)
	ib := &ixbuf.T{}
	for i := 1; i < 10; i++ {
		ib.Insert(strconv.Itoa(i), uint64(i))
	}
	t := &testTran{getIndex: func() *Overlay {
		return &Overlay{bt: bt, layers: []*ixbuf.T{ib}}
	}}

	it := NewOverIter("", 0)
	// Set explicit range from "2" to "7"
	it.Range(Range{Org: "2", End: "7"})

	// Test incremental read tracking with explicit range
	assert.This(t.reads.String()).Is("")

	// First Next() reads from range start "2" to "2"
	it.Next(t)
	assert.This(t.reads.String()).Is("2->2")

	// Second Next() reads from "2" to "3", merges to "2->3"
	it.Next(t)
	assert.This(t.reads.String()).Is("2->3")

	// Third Next() reads from "3" to "4", merges to "2->4"
	it.Next(t)
	assert.This(t.reads.String()).Is("2->4")

	// Prev() reads from "4" to "3", already covered
	it.Prev(t)
	assert.This(t.reads.String()).Is("2->4")

	// Prev() reads from "3" to "2", already covered
	it.Prev(t)
	assert.This(t.reads.String()).Is("2->4")

	// Prev() reads from "2" to range start "2", hits EOF
	it.Prev(t)
	assert.That(it.Eof())
	assert.This(t.reads.String()).Is("2->4")

	// Reset and test Prev from rewind with explicit range
	t.reads = ranges.Ranges{}
	it.Rewind()

	// First Prev() from rewind reads from "6" to range end "7"
	it.Prev(t)
	assert.This(t.reads.String()).Is("6->7")

	// Second Prev() reads from "5" to "6", merges to "5->7"
	it.Prev(t)
	assert.This(t.reads.String()).Is("5->7")

	// Test full forward iteration within explicit range
	t.reads = ranges.Ranges{}
	for it.Rewind(); !it.Eof(); it.Next(t) {
	}
	// Should read from range start to range end
	assert.This(t.reads.String()).Is("2->7")

	// Test full backward iteration within explicit range
	t.reads = ranges.Ranges{}
	for it.Rewind(); !it.Eof(); it.Prev(t) {
	}
	// Should read from range start to range end
	assert.This(t.reads.String()).Is("2->7")
}

func TestOverIterCombine(*testing.T) {
	var data []string
	bt := btree.CreateBtree(stor.HeapStor(8192), nil)
	bt.SetSplit(64)
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

	for range n / 2 {
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
	trace := func(args ...any) {
		// fmt.Print(args...)
	}
	traceln := func(args ...any) {
		// fmt.Println(args...)
	}
	gen := make(dat)
	data := new(dummy)
	it := &dumIter{d: data}
	keyoff := map[string]uint64{}
	bt := btree.CreateBtree(stor.HeapStor(8192), nil)
	ibs := []*ixbuf.T{{}, {}, {}}
	ov := &Overlay{bt: bt, layers: ibs}
	tran := &testTran{getIndex: func() *Overlay { return ov }}
	oi := NewOverIter("", 0)
	check := func() {
		if it.Eof() {
			traceln("EOF")
			assert.That(oi.Eof())
			it.Rewind()
			oi.Rewind()
		} else {
			key, off := oi.Cur()
			traceln(key, off)
			assert.This(key).Is(it.Cur())
			assert.This(off).Is(keyoff[key])
		}
	}
	defer func() {
		if e := recover(); e != nil {
			traceln("===== merge")
			oi.Rewind()
			for oi.Next(tran); !oi.Eof(); oi.Next(tran) {
				key, _ := oi.Cur()
				traceln(key)
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
	var N = 5_000_000
	if testing.Short() {
		N = 200_000
	}
	for n := 1; n < N; n++ { // reserve 0 for deleted
		trace(n, " ")
		switch rand.Intn(7) {
		case 0:
			if rand.Intn(sizing)+1 > data.Len() {
				// insert
				key := gen.randKey()
				data.Insert(key)
				ibs[2].Insert(key, uint64(n))
				keyoff[key] = uint64(n)
				traceln("insert", key, "len", data.Len())
			} else {
				// delete
				j := rand.Intn(data.Len())
				key := data.Get(j)
				traceln("delete", key, "len", data.Len()-1)
				data.Delete(key)
				gen.delete(key)
				ibs[2].Delete(key, keyoff[key])
				keyoff[key] = 0
			}
		case 1: // rewind
			traceln("rewind")
			it.Rewind()
			oi.Rewind()
		case 2:
			traceln("invalidate")
			tmp := *ov
			ov = &tmp
		case 3, 4: // next
			it.Next()
			oi.Next(tran)
			if !oi.Eof() {
				oi.Cur()
			}
			trace("next ")
			check()
		case 5, 6: // prev
			it.Prev()
			oi.Prev(tran)
			if !oi.Eof() {
				oi.Cur()
			}
			trace("prev ")
			check()
		}
		size := 0
		for _, ib := range ibs {
			size += ib.Len()
		}
		assert.This(size).Is(data.Len())
	}
	// for _, ib := range ibs {
	// 	fmt.Println(">>>", ib.Len())
	// }
}

func TestOverIterRandom2(t *testing.T) {
	const nlayers = 10
	const nkeys = 100
	const steps = 100_000

	dum := new(dummy)
	it := &dumIter{d: dum}
	keyoff := map[string]uint64{}
	bt := btree.CreateBtree(stor.HeapStor(8192), nil)
	ibs := []*ixbuf.T{}
	for range nlayers {
		ibs = append(ibs, &ixbuf.T{})
	}
	ov := &Overlay{bt: bt, layers: ibs}
	tran := &testTran{getIndex: func() *Overlay { return ov }}
	oi := NewOverIter("", 0)

	keys := make([]string, nkeys)
	randkey := str.UniqueRandom(3, 3)
	for i := range nkeys {
		keys[i] = randkey()
	}

	nextoff := uint64(1)
	for i := range nlayers {
		for range nkeys / 5 {
			key := keys[rand.Intn(nkeys)]
			if off := keyoff[key]; off != 0 {
				if rand.Intn(3) == 1 {
					ibs[i].Delete(key, off)
					dum.Delete(key)
					keyoff[key] = 0
				} else {
					ibs[i].Update(key, nextoff)
					keyoff[key] = nextoff
					nextoff++
				}
			} else {
				ibs[i].Insert(key, nextoff)
				dum.Insert(key)
				keyoff[key] = nextoff
				nextoff++
			}
		}
	}
	// ov.Print()

	// forward
	for it.Next(); !it.Eof(); it.Next() {
		oi.Next(tran)
		assert.False(oi.Eof())
		key, off := oi.Cur()
		assert.This(key).Is(it.Cur())
		assert.This(off).Is(keyoff[key])
	}
	oi.Next(tran)
	assert.True(oi.Eof())

	// reverse
	it.Rewind()
	oi.Rewind()
	for it.Prev(); !it.Eof(); it.Prev() {
		oi.Prev(tran)
		assert.False(oi.Eof())
		key, off := oi.Cur()
		assert.This(key).Is(it.Cur())
		assert.This(off).Is(keyoff[key])
	}
	oi.Prev(tran)
	assert.True(oi.Eof())

	// random walk
	it.Rewind()
	oi.Rewind()
	for range steps {
		if rand.Intn(2) == 1 {
			it.Next()
			oi.Next(tran)
		} else {
			it.Prev()
			oi.Prev(tran)
		}
		if it.Eof() {
			assert.That(oi.Eof())
			it.Rewind()
			oi.Rewind()
		} else {
			key, off := oi.Cur()
			assert.This(key).Is(it.Cur())
			assert.This(off).Is(keyoff[key])
		}
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
	layers[1].Insert("z", 1|ixbuf.Update)
	layers[2].Insert("a", 2)
	layers[3].Insert("a", 2|ixbuf.Delete)
	layers[3].Insert("z", 1|ixbuf.Delete)
	ov := &Overlay{bt: bt, layers: layers}
	tran := &testTran{getIndex: func() *Overlay { return ov }}

	it := NewOverIter("", 0)
	it.Next(tran)
	assert.That(it.Eof())

	layers = []*ixbuf.T{{}, {}, {}, {}}
	layers[0].Insert("a", 1)
	layers[1].Insert("a", 1|ixbuf.Update)
	layers[2].Insert("z", 2)
	layers[3].Insert("z", 2|ixbuf.Delete)
	layers[3].Insert("a", 1|ixbuf.Delete)
	ov = &Overlay{bt: bt, layers: layers}

	it = NewOverIter("", 0)
	it.Prev(tran)
	assert.That(it.Eof())
}

func TestOverIterBug2(*testing.T) {
	b := btree.Builder(stor.HeapStor(8192))
	assert.That(b.Add("1111", 1111))
	assert.That(b.Add("2222", 2222))
	bt := b.Finish()
	layers := []*ixbuf.T{{}}
	layers[0].Insert("1111", 1111|ixbuf.Delete)
	layers[0].Insert("2222", 2222|ixbuf.Delete)
	ov := &Overlay{bt: bt, layers: layers}
	tran := &testTran{getIndex: func() *Overlay { return ov }}
	it := NewOverIter("", 0)
	it.Next(tran)
	assert.That(it.Eof())
}

func TestOverIterBug3(*testing.T) {
	b := btree.Builder(stor.HeapStor(8192))
	assert.That(b.Add("1111", 1111))
	bt := b.Finish()
	layers := []*ixbuf.T{{}}
	layers[0].Insert("1111", 1111|ixbuf.Delete)
	ov := &Overlay{bt: bt, layers: layers}
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

func TestOverIterBug4(*testing.T) {
	b := btree.Builder(stor.HeapStor(8192))
	assert.That(b.Add("1", 1))
	bt := b.Finish()
	// bt.Print()
	layers := []*ixbuf.T{{}}
	layers[0].Update("1", 2)
	// layers[0].Print()
	ov := &Overlay{bt: bt, layers: layers}
	tran := &testTran{getIndex: func() *Overlay { return ov }}
	it := NewOverIter("", 0)
	it.Prev(tran)
	assert.That(!it.Eof())
	key, off := it.Cur()
	assert.This(key).Is("1")
	assert.This(off).Is(2)
	it.Next(tran)
	assert.That(it.Eof())
}

func TestOverIterInvalidCombine(t *testing.T) {
	bt := btree.CreateBtree(stor.HeapStor(8192), nil)
	layers := make([]*ixbuf.T, 3)
	for i := range layers {
		layers[i] = &ixbuf.T{}
	}
	layers[0].Insert("k0", 2)
	layers[0].Insert("k1", 1)
	layers[1].Insert("k1", 1|ixbuf.Delete)
	layers[1].Insert("k2", 3)
	layers[2].Insert("k0", 5|ixbuf.Update)
	layers[2].Insert("k1", 4)

	ov := &Overlay{bt: bt, layers: layers}
	// fmt.Println(ov)
	tran := &testTran{getIndex: func() *Overlay { return ov }}

	oi := NewOverIter("", 0)
	oi.Range(Range{Org: "k2", End: "k2z"})

	oi.Next(tran)

	// trigger OverIter modified
	tmp2 := *ov
	ov = &tmp2

	oi.Prev(tran) // => invalid Combine add,add +1 +4
}

func TestOverIterCurDeletedBug(t *testing.T) {
	bt := btree.CreateBtree(stor.HeapStor(8192), nil)
	layers := make([]*ixbuf.T, 3)
	for i := range layers {
		layers[i] = &ixbuf.T{}
	}
	layers[0].Insert("k2", 2)
	layers[0].Insert("k4", 1)
	layers[1].Insert("k3", 3)
	layers[1].Insert("k4", 1|ixbuf.Delete)
	layers[2].Insert("k1", 4)
	layers[2].Insert("k2", 2|ixbuf.Delete)

	ov := &Overlay{bt: bt, layers: layers}
	// fmt.Println(ov)
	tran := &testTran{getIndex: func() *Overlay { return ov }}

	oi := NewOverIter("", 0)
	oi.Range(Range{Org: "k2", End: "k4"})

	oi.Prev(tran)
	if !oi.Eof() {
		oi.Cur()
	}

	// trigger OverIter modified
	tmp2 := *ov
	ov = &tmp2

	// try to trigger the OverIter Cur deleted bug
	oi.Prev(tran)
	if !oi.Eof() {
		oi.Cur()
	}
}

func TestOverIterRandom3(t *testing.T) {
	// NOTE: this is just a smoke test, it does not check the data at all
	n := 400_000
	if testing.Short() {
		n = 10_000
	}
	var ops []string
	var ov *Overlay
	var keys []string
	var keyOffsets map[string]uint64
	_ = keyOffsets
	var nextOffset uint64
	var oi *OverIter
	var i int
	seed := rand.Int63()
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("seed", seed, "loop", i)
			fmt.Println(ops)
			fmt.Println(ov)
			panic(r)
		}
	}()
	localRand := rand.New(rand.NewSource(seed))
	for i = range n {
		ops = ops[:0]
		ov, keys, keyOffsets, nextOffset = generateRandomLayers(localRand)
		oi = NewOverIter("", 0)
		tran := &testTran{getIndex: func() *Overlay { return ov }}
		for range 100 {
			oi.Check()
			switch localRand.Intn(6) {
			case 0:
				ops = append(ops, "next")
				oi.Next(tran)
				if !oi.Eof() {
					oi.Cur()
				}
			case 1:
				ops = append(ops, "prev")
				oi.Prev(tran)
				if !oi.Eof() {
					oi.Cur()
				}
			case 2:
				ops = append(ops, "rewind")
				oi.Rewind()
			case 3:
				_ = nextOffset
				topLayer := ov.layers[len(ov.layers)-1]
				mod := randomModifyLayer(topLayer, keys, keyOffsets, &nextOffset, localRand)
				ops = append(ops, "modify "+mod)
				ov.Check()
			case 4: // start a new transaction
				if localRand.Intn(3) == 0 {
					// Simulate concurrent modification: generate completely new overlay
					// This tests OverIter's update detection when overlay structure changes
					ov, keys, keyOffsets, nextOffset = generateRandomLayers(localRand)
					ops = append(ops, "newtran+")
				} else {
					// Normal case: just start new transaction with same overlay
					tmp := *ov
					ov = &tmp
					ops = append(ops, "newtran-")
				}
			case 5: // test Range functionality with new iterator
				oi = createRandomRangedIterator(keys, localRand)
				ops = append(ops, "range"+oi.rng.String())
			}
		}
	}
}

func generateRandomLayers(localRand *rand.Rand) (*Overlay, []string, map[string]uint64, uint64) {
	bt := btree.CreateBtree(stor.HeapStor(8192), nil)

	nlayers := localRand.Intn(7) + 2
	layers := make([]*ixbuf.T, nlayers)
	for i := range layers {
		layers[i] = &ixbuf.T{}
	}

	keys := []string{"k0", "k1", "k2", "k3", "k4", "k5", "k6", "k7"}
	nkeys := localRand.Intn(len(keys)-1) + 1
	keyOffsets := make(map[string]uint64) // track key offset, 0 if non-existent
	nextOffset := uint64(1)

	for _, layer := range layers {
		nmods := localRand.Intn(5)
		for range nmods {
			randomModifyLayer(layer, keys[:nkeys], keyOffsets, &nextOffset, localRand)
		}
	}

	ov := &Overlay{bt: bt, layers: layers}
	ov.Check()
	return ov, keys[:nkeys], keyOffsets, nextOffset
}

func randomModifyLayer(layer *ixbuf.T, keys []string, keyOffsets map[string]uint64, nextOffset *uint64, localRand *rand.Rand) string {
	key := keys[localRand.Intn(len(keys))]
	if keyOffsets[key] == 0 {
		// Key doesn't exist - can only do a plain add (no flags)
		layer.Insert(key, *nextOffset)
		keyOffsets[key] = *nextOffset
		*nextOffset++
		return key + "+" + strconv.Itoa(int(*nextOffset-1))
	} else {
		// key exists
		if localRand.Intn(2) == 1 {
			layer.Insert(key, *nextOffset|ixbuf.Update)
			keyOffsets[key] = *nextOffset
			*nextOffset++
			return key + "=" + strconv.Itoa(int(*nextOffset-1))
		} else {
			layer.Insert(key, keyOffsets[key]|ixbuf.Delete)
			off := keyOffsets[key]
			keyOffsets[key] = 0 // key no longer exists
			return key + "-" + strconv.Itoa(int(off))
		}
	}
}

func createRandomRangedIterator(keys []string, localRand *rand.Rand) *OverIter {
	it := NewOverIter("", 0)
	switch localRand.Intn(4) {
	case 0:
		// Full range (default)
		return it
	case 1:
		// Range with two different keys
		i := localRand.Intn(len(keys))
		j := localRand.Intn(len(keys))
		if i == j {
			j = (j + 1) % len(keys)
		}
		x, y := keys[i], keys[j]
		it.Range(Range{Org: min(x, y), End: max(x, y)})
	case 2:
		// Range with single key as both start and end
		key := keys[localRand.Intn(len(keys))]
		it.Range(Range{Org: key, End: key})
	case 3:
		// Range with key prefixes or variations
		key := keys[localRand.Intn(len(keys))]
		it.Range(Range{Org: key, End: key + "z"})
	}
	return it
}

func TestCheckOverlay(t *testing.T) {
	assert := assert.T(t)
	bt := btree.CreateBtree(stor.HeapStor(8192), nil)

	// Test valid overlay
	layer1 := &ixbuf.T{}
	layer1.Insert("a", 1)
	layer1.Insert("b", 2)

	layer2 := &ixbuf.T{}
	layer2.Update("a", 3) // valid update
	layer2.Insert("c", 4)

	layer3 := &ixbuf.T{}
	layer3.Delete("b", 2) // valid delete with correct offset

	ov := &Overlay{bt: bt, layers: []*ixbuf.T{layer1, layer2, layer3}}

	// Should not panic
	ov.Check()

	// Test invalid overlay - update without add
	// Create two layers: first has no "x", second tries to update "x"
	baseLayer := &ixbuf.T{}
	baseLayer.Insert("a", 1) // some other key

	invalidLayer := &ixbuf.T{}
	invalidLayer.Insert("x", 99|ixbuf.Update) // manually create update without prior add

	invalidOv := &Overlay{bt: bt, layers: []*ixbuf.T{baseLayer, invalidLayer}}
	assert.This(func() { invalidOv.Check() }).
		Panics("update of non-existent key")

	// Test invalid delete offset
	invalidDeleteLayer1 := &ixbuf.T{}
	invalidDeleteLayer1.Insert("d", 5)

	invalidDeleteLayer2 := &ixbuf.T{}
	invalidDeleteLayer2.Insert("d", 4|ixbuf.Delete) // wrong offset for delete

	invalidDeleteOv := &Overlay{bt: bt, layers: []*ixbuf.T{invalidDeleteLayer1, invalidDeleteLayer2}}

	assert.This(func() { invalidDeleteOv.Check() }).
		Panics("delete offset mismatch")
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

// func (d *dummy) print() {
// 	fmt.Println("+ + +")
// 	for _, k := range d.keys {
// 		fmt.Println(k)
// 	}
// 	fmt.Println("+ + +")
// }

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

// BenchmarkOverIterSingle benchmarks iteration when singleIter optimization applies
// i.e., when there are no layers and mut is empty (read-only transactions)
func BenchmarkOverIterSingle(b *testing.B) {
	const nrecs = 1000
	bldr := btree.Builder(stor.HeapStor(8192))
	for i := 1; i <= nrecs; i++ {
		// Zero-pad to ensure lexicographic order matches numeric order
		key := strconv.Itoa(i*10 + 100000) // e.g., "100010", "100020", etc.
		bldr.Add(key, uint64(i))
	}
	bt := bldr.Finish()
	ov := &Overlay{bt: bt, layers: []*ixbuf.T{{}}} // single empty layer, no mut
	tran := &testTran{getIndex: func() *Overlay { return ov }}

	b.Run("Next single", func(b *testing.B) {
		for b.Loop() {
			it := NewOverIter("", 0)
			for it.Next(tran); !it.Eof(); it.Next(tran) {
			}
			assert.That(it.singleIter)
		}
	})

	b.Run("Prev single", func(b *testing.B) {
		for b.Loop() {
			it := NewOverIter("", 0)
			for it.Prev(tran); !it.Eof(); it.Prev(tran) {
			}
			assert.That(it.singleIter)
		}
	})

	b.Run("Next merge", func(b *testing.B) {
		for b.Loop() {
			it := NewOverIter("", 0)
			it.update(tran)       // sets singleIter = true
			it.singleIter = false // force unoptimized path
			for it.Next(tran); !it.Eof(); it.Next(tran) {
			}
			assert.That(!it.singleIter)
		}
	})

	b.Run("Prev merge", func(b *testing.B) {
		for b.Loop() {
			it := NewOverIter("", 0)
			it.update(tran)       // sets singleIter = true
			it.singleIter = false // force unoptimized path
			for it.Prev(tran); !it.Eof(); it.Prev(tran) {
			}
			assert.That(!it.singleIter)
		}
	})
}
