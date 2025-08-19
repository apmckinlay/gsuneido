// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"bytes"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestVarint(t *testing.T) {
	for _, num := range []int{0, 1, 127, 128, 129, 255, 256, 999, 12345} {
		var nd node
		nd = addVarint(nd, num)
		_, n := getVarint(nd)
		assert.This(n).Is(num)
	}
}

func TestNodeAppendRead(t *testing.T) {
	type ent struct {
		offset uint64
		npre   int
		diff   string
	}
	var nd node
	var data []ent
	add := func(offset uint64, npre int, diff string) {
		nd = nd.append(offset, npre, diff)
		data = append(data, ent{offset: offset, npre: npre, diff: diff})
	}
	add(123, 999, "bar")
	add(456, 11, "foo")
	for _, e := range data {
		var npre int
		var diff []byte
		var off uint64
		npre, diff, off = nd.read()
		nd = nd[fLen(npre, diff):]
		assert.T(t).This(npre).Is(e.npre)
		assert.T(t).This(string(diff)).Is(e.diff)
		assert.T(t).This(off).Is(e.offset)
	}
}

func TestNodeSize(t *testing.T) {
	var nd node
	assert.T(t).This(nd.Size()).Is(0)
	nd = nd.append(123456, 123, "foo")
	assert.T(t).This(nd.Size()).Is(1)
	nd = nd.append(123, 456, "barbaz")
	assert.T(t).This(nd.Size()).Is(2)
}

func TestNodeInsert(*testing.T) {
	datas := []string{
		"a b c d",
		"xa xb xc xd",
		"a ab abc abcd",
		"a ant anti ants b bun bunn bunnies bunny buns ca cat cats",
		"a aa aaa ab abc b ba bb bba bbb bbc bc c",
		"1000 1001 1002 1003",
	}
	for _, s := range datas {
		// fmt.Println(s)
		data := strings.Fields(s)
		get := func(i uint64) string { return data[i] }

		// forward
		fwd := node{}
		for i, d := range data {
			fwd = fwd.update(d, uint64(i), get)
			fwd.checkUpTo(i, data, get)
		}
		// fmt.Println("forward")
		// fwd.printLeafNode(get)
		assert.That(fwd.check(get) == len(data))

		// reverse
		rev := node{}
		for i := len(data) - 1; i >= 0; i-- {
			rev = rev.update(data[i], uint64(i), get)
			// rev.checkUpTo(i, data, get)
			// fmt.Println()
			// rev.printLeafNode(get)
		}
		// fmt.Println("reverse")
		// rev.printLeafNode(get)
		assert.This(rev).Is(fwd)

		// builder
		bld := build(data)
		// fmt.Println("builder")
		// bld.printLeafNode(get)
		assert.This(bld).Is(fwd)

		// random
		const nperms = 100
		for range nperms {
			rnd := node{}
			perm := rand.Perm(len(data))
			for _, j := range perm {
				rnd = rnd.update(data[j], uint64(j), get)
			}
			assert.This(rnd).Is(fwd)
		}
	}
}

func build(data []string) node {
	b := nodeBuilder{}
	for i, d := range data {
		b.Add(d, uint64(i), 1)
	}
	return b.Entries()
}

func (nd node) checkData(data []string, get func(uint64) string) {
	nd.checkUpTo(len(data)-1, data, get)
}

// checkUpTo is used during inserting.
// It checks that inserted keys are present
// and uninserted keys are not present.
func (nd node) checkUpTo(i int, data []string, get func(uint64) string) {
	n := nd.check(get)
	nn := 0
	for j, d := range data {
		if (d != "" && j <= i) != nd.contains(d, get) {
			panic("can't find " + d)
		}
		if d != "" && j <= i {
			nn++
		}
	}
	if nn != n {
		panic("check count expected " + strconv.Itoa(n) +
			" got " + strconv.Itoa(nn))
	}
}

func (nd node) contains(s string, get func(uint64) string) bool {
	if len(nd) == 0 {
		return false
	}
	offset := nd.search(s)
	return s == get(offset)
}

func TestNodeInsertRandom(*testing.T) {
	const nData = 100
	var nGenerate = 1000
	var nShuffle = 20
	if testing.Short() {
		nGenerate = 1
		nShuffle = 4
	}
	var data = make([]string, nData)
	get := func(i uint64) string { return data[i] }
	for range nGenerate {
		data = data[0:nData]
		randKey := str.UniqueRandomOf(1, 6, "abcdef")
		for di := range nData {
			data[di] = randKey()
		}
		for range nShuffle {
			rand.Shuffle(len(data),
				func(i, j int) { data[i], data[j] = data[j], data[i] })
			var nd node
			for i, d := range data {
				nd = nd.update(d, uint64(i), get)
				// fe.checkUpTo(i, data, get)
			}
			nd.checkData(data, get)
		}
	}
}

func TestNodeInsertWords(*testing.T) {
	data := words
	const nShuffle = 100
	get := func(i uint64) string { return data[i] }
	for range nShuffle {
		rand.Shuffle(len(data),
			func(i, j int) { data[i], data[j] = data[j], data[i] })
		var nd node
		for i, d := range data {
			nd = nd.update(d, uint64(i), get)
			// fe.checkUpto(i, data, get)
		}
		nd.checkData(data, get)
	}
}

var words = []string{
	"tract",
	"pluck",
	"rumor",
	"choke",
	"abbey",
	"robot",
	"north",
	"dress",
	"pride",
	"dream",
	"judge",
	"coast",
	"frank",
	"suite",
	"merit",
	"chest",
	"youth",
	"throw",
	"drown",
	"power",
	"ferry",
	"waist",
	"moral",
	"woman",
	"swipe",
	"straw",
	"shell",
	"class",
	"claim",
	"tired",
	"stand",
	"chaos",
	"shame",
	"thigh",
	"bring",
	"lodge",
	"amuse",
	"arrow",
	"charm",
	"swarm",
	"serve",
	"world",
	"raise",
	"means",
	"honor",
	"grand",
	"stock",
	"model",
	"greet",
	"basic",
	"fence",
	"fight",
	"level",
	"title",
	"knife",
	"wreck",
	"agony",
	"white",
	"child",
	"sport",
	"cheat",
	"value",
	"marsh",
	"slide",
	"tempt",
	"catch",
	"valid",
	"study",
	"crack",
	"swing",
	"plead",
	"flush",
	"awful",
	"house",
	"stage",
	"fever",
	"equal",
	"fault",
	"mouth",
	"mercy",
	"colon",
	"belly",
	"flash",
	"style",
	"plant",
	"quote",
	"pitch",
	"lobby",
	"gloom",
	"patch",
	"crime",
	"anger",
	"petty",
	"spend",
	"strap",
	"novel",
	"sword",
	"match",
	"tasty",
	"stick",
}

func TestNodeDelete(*testing.T) {
	datas := []string{
		"a b c d",
		"xa xb xc xd",
		"a ab abc abcd",
		"a baac baacb c cbaba",
		"a ant anti ants b bun bunn bunnies bunny buns ca cat cats",
		"a aa aaa ab abc b ba bb bba bbb bbc bc c",
		"1000 1001 1002 1003",
	}
	var data []string
	get := func(i uint64) string { return data[i] }
	var without, del node
	finished := false
	defer func() {
		if !finished {
			build(data).printLeafNode(get)
			fmt.Println("WITHOUT")
			without.printLeafNode(get)
			fmt.Println("DELETED")
			del.printLeafNode(get)
		}
	}()
	var ntimes = 100000
	if testing.Short() {
		ntimes = 10000
	}
	for i := range ntimes {
		if i < len(datas) {
			data = strings.Fields(datas[i])
		} else {
			// random data
			r := str.UniqueRandomOf(1, 6, "abc")
			data = make([]string, 5)
			for j := range data {
				data[j] = r()
			}
			sort.Strings(data)
		}
		var datawo []string
		for i := range len(data) {
			all := build(data)
			datawo = append(datawo[:0], data...)
			datawo[i] = ""
			// if i == 2 {
			// 	fmt.Println("===========================")
			// 	all.printLeafNode(get)
			// }
			// fmt.Println("---------------------------", i, data[i])

			b := nodeBuilder{}
			for j, d := range data {
				if j != i {
					b.Add(d, uint64(j), 1)
				}
			}
			without = b.Entries()
			// fmt.Println("WITHOUT")
			// without.printLeafNode(get)
			without.check(get)
			without.checkData(datawo, get)

			// fmt.Println("DELETED")
			// del.printLeafNode(get)
			del, _ = all.delete(uint64(i))
			del.check(get)
			del.checkData(datawo, get)

			if err := compare(without, del); err != "" {
				panic(err)
			}
		}
	}
	finished = true
}

func compare(nd1, nd2 node) string {
	it1 := nd1.iter()
	it2 := nd2.iter()
	for {
		ok1 := it1.next()
		ok2 := it2.next()
		if ok1 != ok2 {
			return "DIFFERENT lengths"
		}
		if !ok1 {
			return "" // ok
		}
		switch {
		case it1.offset != it2.offset:
			return "DIFFERENT offsets"
		case it1.npre != it2.npre:
			return "DIFFERENT npre"
		case !bytes.HasPrefix(it2.known, it1.known):
			return "DIFFERENT known not prefix"
		}
	}
}

func TestNodeInsertDelete(*testing.T) {
	var ok bool
	var nd node
	var data []slot
	dup := func(key string) bool {
		for i := range data {
			if data[i].key == key {
				return true
			}
		}
		return false
	}
	get := func(off uint64) string {
		for i := range data {
			if data[i].off == off {
				return data[i].key
			}
		}
		panic("get: offset not found")
	}
	r := func() string { return str.RandomOf(1, 6, "abc") }
	var N = 1_000_000
	if testing.Short() {
		N = 100_000
	}
	for i := range N {
		if rand.Intn(13)+1 > len(data) {
			// insert
			key := r()
			for ; dup(key); key = r() {
			}
			off := uint64(i)
			// print("insert", key, off)
			nd = nd.update(key, off, get)
			data = append(data, slot{key: key, off: off})
		} else {
			// delete
			i := rand.Intn(len(data))
			// print("delete", data[i].key, data[i].off)
			nd, ok = nd.delete(data[i].off)
			assert.That(ok)
			data[i] = data[len(data)-1]
			data = data[:len(data)-1]
		}
		assert.This(nd.check(get)).Is(len(data))
		for _, s := range data {
			if !nd.contains(s.key, get) {
				print(data)
				nd.printLeafNode(get)
				panic("lookup failed " + s.key)
			}
		}
	}
}

var S1 []byte
var S2 string

func BenchmarkNodeIter(b *testing.B) {
	get := func(i uint64) string { return words[i] }
	var nd node
	for i, d := range words {
		nd = nd.update(d, uint64(i), get)
	}
	for b.Loop() {
		iter := nd.iter()
		for iter.next() {
			S1 = iter.known
		}
	}
}

func BenchmarkUnodeIter(b *testing.B) {
	var u unode
	for i, w := range words {
		u = append(u, slot{key: w, off: uint64(i)})
	}
	for b.Loop() {
		ui := &unodeIter{u: u, i: -1}
		for ui.next() {
			S2 = u[ui.i].key
		}
	}
}

func Example_node_BuilderSplit() {
	var b nodeBuilder
	b.Add("1234xxxx", 1234, 1)
	b.Add("1235xxxx", 1235, 1)
	b.Add("1299xxxx", 1299, 1)
	b.Add("1300xxxx", 1300, 1)
	b.Add("1305xxxx", 1305, 1)
	st := stor.HeapStor(8192)
	leftOff, splitKey := b.Split(st)
	// assert.T(t).This(splitKey).Is("13")
	fmt.Println("splitKey", splitKey)
	fmt.Println("LEFT ---")
	readNode(st, leftOff).print()
	fmt.Println("RIGHT ---")
	b.node.print()

	// Output:
	// splitKey 13
	// LEFT ---
	// 1234 ''
	// 1235 1235
	// 1299 129
	// RIGHT ---
	// 1300 ''
	// 1305 1305
}

func Example_merge_split() {
	var b nodeBuilder
	b.Add("1234xxxx", 1234, 1)
	b.Add("1235xxxx", 1235, 1)
	b.Add("1299xxxx", 1299, 1)
	b.Add("1300xxxx", 1300, 1)
	b.Add("1305xxxx", 1305, 1)
	m := merge{node: b.node, modified: true}
	left, right, splitKey := m.split(b.node.Size())
	fmt.Println("splitKey", splitKey)
	fmt.Println("LEFT ---")
	left.print()
	fmt.Println("RIGHT ---")
	right.print()

	// Output:
	// splitKey 129
	// LEFT ---
	// 1234 ''
	// 1235 1235
	// RIGHT ---
	// 1299 ''
	// 1300 13
	// 1305 1305
}

func Test_split_bug(*testing.T) {
	var b nodeBuilder
	b.Add("a", 1, 1)
	b.Add("b", 2, 1)
	b.Add("c", 3, 1)
	b.Add("d", 4, 1)
	b.Add("e", 5, 1)
	b.Add("f", 6, 1)
	b.Add("g", 7, 1)
	b.Add(strings.Repeat("z", 64), 8, 99)
	m := merge{node: b.node, modified: true}
	left, right, splitKey := m.split(b.node.Size())
	fmt.Println("splitKey", splitKey)
	fmt.Println("LEFT ---")
	left.print()
	fmt.Println("RIGHT ---")
	right.print()
}
