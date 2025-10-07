// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package tests

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"testing"

	btree1 "github.com/apmckinlay/gsuneido/db19/index/btree"
	btree3 "github.com/apmckinlay/gsuneido/db19/index/btree3"
	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/bits"
	"github.com/apmckinlay/gsuneido/util/hacks"
)

const nEntries = 100_000
const heapSize = 1024 * 1024

func BenchmarkBtreeBuilder(b *testing.B) {
	for b.Loop() {
		st := stor.HeapStor(heapSize)
		builder := btree1.Builder(st)
		for i := range nEntries {
			key := fmt.Sprintf("%05d", i)
			builder.Add(key, uint64(i))
		}
		builder.Finish()
		assert.That(st.Size() < heapSize)
	}
}

func BenchmarkBtreeBuilder3(b *testing.B) {
	for b.Loop() {
		st := stor.HeapStor(heapSize)
		builder := btree3.Builder(st)
		for i := range nEntries {
			key := fmt.Sprintf("%05d", i)
			builder.Add(key, uint64(i))
		}
		builder.Finish()
		assert.That(st.Size() < heapSize)
	}
}

func BenchmarkBtreeBuilderBase(b *testing.B) {
	for b.Loop() {
		for i := range nEntries {
			_ = fmt.Sprintf("%05d", i)
		}
	}
}

func BenchmarkBtreeLookup(b *testing.B) {
	defer func(prev func(*stor.Stor, *ixkey.Spec, uint64) string) {
		btree1.GetLeafKey = prev
	}(btree1.GetLeafKey)
	n := 0
	btree1.GetLeafKey = func(st *stor.Stor, _ *ixkey.Spec, off uint64) string {
		n++
		buf := st.Data(off)
		return hacks.BStoS(buf[:5])
	}
	st := stor.HeapStor(heapSize)
	builder := btree1.Builder(st)
	for i := range nEntries {
		key := fmt.Sprintf("%05d", i)
		off, buf := st.Alloc(100)
		copy(buf, key)
		builder.Add(key, off)
	}
	bt := builder.Finish()
	for b.Loop() {
		bt.Lookup(fmt.Sprintf("%05d", rand.IntN(nEntries)))
	}
	fmt.Println("GetLeafKey", n)
}

func BenchmarkBtreeLookup3(b *testing.B) {
	st := stor.HeapStor(heapSize)
	builder := btree3.Builder(st)
	for i := range nEntries {
		key := fmt.Sprintf("%05d", i)
		builder.Add(key, uint64(i))
	}
	bt := builder.Finish()
	for b.Loop() {
		bt.Lookup(fmt.Sprintf("%05d", rand.IntN(nEntries)))
	}
}

func BenchmarkBtreeLookupBase(b *testing.B) {
	for b.Loop() {
		_ = fmt.Sprintf("%05d", rand.IntN(nEntries))
	}
}

func BenchmarkBtreeIter(b *testing.B) {
	defer func(prev func(*stor.Stor, *ixkey.Spec, uint64) string) {
		btree1.GetLeafKey = prev
	}(btree1.GetLeafKey)
	glf := 0
	btree1.GetLeafKey = func(st *stor.Stor, _ *ixkey.Spec, off uint64) string {
		glf++
		buf := st.Data(off)
		return string(buf[:5])
	}
	st := stor.HeapStor(heapSize)
	builder := btree1.Builder(st)
	for i := range nEntries {
		key := fmt.Sprintf("%05d", i)
		off, buf := st.Alloc(100)
		copy(buf, key)
		builder.Add(key, off)
	}
	loops := 0
	bt := builder.Finish()
	for b.Loop() {
		loops++
		iter := bt.Iterator()
		iter.Range(btree3.Range{Org: "01000", End: "99000"})
		for iter.Next(); !iter.Eof(); iter.Next() {
			// _, _ = iter.Cur()
		}
	}
	fmt.Println("GetLeafKey", glf/loops)
}

func BenchmarkBtreeIter3(b *testing.B) {
	st := stor.HeapStor(heapSize)
	builder := btree3.Builder(st)
	for i := range nEntries {
		key := fmt.Sprintf("%05d", i)
		builder.Add(key, uint64(i))
	}
	bt := builder.Finish()
	for b.Loop() {
		iter := bt.Iterator()
		iter.Range(btree3.Range{Org: "01000", End: "99000"})
		for iter.Next(); !iter.Eof(); iter.Next() {
			// _, _ = iter.Cur()
		}
	}
}

var X byte

func BenchmarkBtreeMergeAdd(b *testing.B) {
	rng := rand.New(rand.NewPCG(123, 456))
	st := stor.HeapStor(heapSize)
	st.Alloc(1) // avoid offset 0
	defer func(prev func(*stor.Stor, *ixkey.Spec, uint64) string) {
		btree1.GetLeafKey = prev
	}(btree1.GetLeafKey)
	glf := 0
	btree1.GetLeafKey = func(st *stor.Stor, _ *ixkey.Spec, off uint64) string {
		glf++
		buf := st.Data(rand.Uint64N(st.Size()))
		X = buf[0]
		return strconv.Itoa(int(off - 1))
	}
	bt := btree1.Builder(st).Finish()

	i := uint32(0)
	x := &ixbuf.T{}
	for b.Loop() {
		batchSize := rng.IntN(100) + 1
		for range batchSize {
			keyNum := bits.Shuffle32(i)
			key := strconv.Itoa(int(keyNum))
			x.Insert(key, uint64(keyNum)+1) // +1 to avoid 0
			i++
			assert.That(i != 0)
		}
		bt = bt.MergeAndSave(x.Iter())
		x.Clear()
	}
	fmt.Println(i, glf)
}

func BenchmarkBtreeMergeAdd3(b *testing.B) {
	rng := rand.New(rand.NewPCG(123, 456))
	st := stor.HeapStor(heapSize)
	st.Alloc(1) // avoid offset 0
	empty := btree3.Builder(st).Finish()
	bt := empty

	i := uint32(0)
	x := &ixbuf.T{}
	for b.Loop() {
		batchSize := rng.IntN(100) + 1
		for range batchSize {
			// shuffle the bits to randomize the order
			// and interleave to get some dense updates
			keyNum := interleave3(int(bits.Shuffle32(i>>3)), int(i&7))
			key := strconv.Itoa(int(keyNum))
			x.Insert(key, uint64(keyNum)+1) // +1 to avoid 0
			i++
			assert.That(i != 0)
		}
		bt = bt.MergeAndSave(x.Iter())
		x.Clear()
	}
	fmt.Println(i, st.Size())
}

func BenchmarkBtreeMergeMix(b *testing.B) {
	rng := rand.New(rand.NewPCG(123, 456))
	st := stor.HeapStor(heapSize)
	st.Alloc(1) // avoid offset 0
	defer func(prev func(*stor.Stor, *ixkey.Spec, uint64) string) {
		btree1.GetLeafKey = prev
	}(btree1.GetLeafKey)
	glf := 0
	btree1.GetLeafKey = func(st *stor.Stor, _ *ixkey.Spec, off uint64) string {
		glf++
		buf := st.Data(rand.Uint64N(st.Size()))
		X = buf[0]
		return strconv.Itoa(int(off - 1))
	}
	bt := btree1.Builder(st).Finish()

	i := uint32(0) // insert counter
	ib := &ixbuf.T{}

	for range 1000 {
		keyNum := interleave3(int(bits.Shuffle32(i>>3)), int(i&7))
		key := strconv.Itoa(int(keyNum))
		ib.Insert(key, uint64(keyNum)+1) // +1 to avoid 0
		i++
	}
	bt = bt.MergeAndSave(ib.Iter())

	u := uint32(500) // update counter, start 500 behind insert
	d := uint32(0)   // delete counter, start 1000 behind insert
	for b.Loop() {
		ib.Clear()
		n := rng.IntN(30) + 1
		for range n + 1 { // 1 extra insert to gradually grow btree
			// shuffle the bits to randomize the order
			// and interleave to get some dense updates
			keyNum := interleave3(int(bits.Shuffle32(i>>3)), int(i&7))
			key := strconv.Itoa(int(keyNum))
			ib.Insert(key, uint64(keyNum)+1) // +1 to avoid 0
			i++
		}
		for range n {
			keyNum := interleave3(int(bits.Shuffle32(u>>3)), int(u&7))
			key := strconv.Itoa(int(keyNum))
			ib.Update(key, uint64(keyNum)+1) // +1 to avoid 0
			u++
		}
		for range n {
			keyNum := interleave3(int(bits.Shuffle32(d>>3)), int(d&7))
			key := strconv.Itoa(int(keyNum))
			ib.Delete(key, uint64(keyNum)+1) // +1 to avoid 0
			d++
		}
		bt = bt.MergeAndSave(ib.Iter())
	}
	fmt.Println(i, st.Size())
}

func BenchmarkBtreeMergeMix3(b *testing.B) {
	rng := rand.New(rand.NewPCG(123, 456))
	st := stor.HeapStor(heapSize)
	st.Alloc(1) // avoid offset 0
	bt := btree3.Builder(st).Finish()

	i := uint32(0) // insert counter
	ib := &ixbuf.T{}

	for range 1000 {
		keyNum := interleave3(int(bits.Shuffle32(i>>3)), int(i&7))
		key := strconv.Itoa(int(keyNum))
		ib.Insert(key, uint64(keyNum)+1) // +1 to avoid 0
		i++
	}
	bt = bt.MergeAndSave(ib.Iter())

	u := uint32(500) // update counter, start 500 behind insert
	d := uint32(0)   // delete counter, start 1000 behind insert
	for b.Loop() {
		ib.Clear()
		n := rng.IntN(30) + 1
		for range n + 1 { // 1 extra insert to gradually grow btree
			// shuffle the bits to randomize the order
			// and interleave to get some dense updates
			keyNum := interleave3(int(bits.Shuffle32(i>>3)), int(i&7))
			key := strconv.Itoa(int(keyNum))
			ib.Insert(key, uint64(keyNum)+1) // +1 to avoid 0
			i++
		}
		for range n {
			keyNum := interleave3(int(bits.Shuffle32(u>>3)), int(u&7))
			key := strconv.Itoa(int(keyNum))
			ib.Update(key, uint64(keyNum)+1) // +1 to avoid 0
			u++
		}
		for range n {
			keyNum := interleave3(int(bits.Shuffle32(d>>3)), int(d&7))
			key := strconv.Itoa(int(keyNum))
			ib.Delete(key, uint64(keyNum)+1) // +1 to avoid 0
			d++
		}
		bt = bt.MergeAndSave(ib.Iter())
	}
	fmt.Println(i, st.Size())
}

// interleave3 interleaves the bottom 3 bits of b into n.
// n is shifted left 3 places and the bottom 6 bits become n2 b2 n1 b1 n0 b0
func interleave3(n int, b int) int {
	// Shift upper bits of n (bit 3 and above) left by 3
	result := (n >> 3) << 6

	// Interleave bottom 3 bits of n and b
	result |= ((n >> 2) & 1) << 5 // n2 to bit 5
	result |= ((b >> 2) & 1) << 4 // b2 to bit 4
	result |= ((n >> 1) & 1) << 3 // n1 to bit 3
	result |= ((b >> 1) & 1) << 2 // b1 to bit 2
	result |= (n & 1) << 1        // n0 to bit 1
	result |= b & 1               // b0 to bit 0

	return result
}

func TestInterleave3(t *testing.T) {
	result := interleave3(0, 0)
	assert.This(result).Is(0)
	result = interleave3(0, 1)
	assert.This(result).Is(1)
	result = interleave3(1, 0)
	assert.This(result).Is(2)
	result = interleave3(1, 1)
	assert.This(result).Is(3)
	result = interleave3(2, 0)
	assert.This(result).Is(8)
	result = interleave3(2, 1)
	assert.This(result).Is(9)
	result = interleave3(3, 0)
	assert.This(result).Is(10)
	result = interleave3(3, 1)
	assert.This(result).Is(11)
}
