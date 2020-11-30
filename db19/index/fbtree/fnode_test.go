// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package fbtree

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

func TestFAppendRead(t *testing.T) {
	type ent struct {
		offset uint64
		npre   int
		diff   string
	}
	var fn fnode
	var data []ent
	add := func(offset uint64, npre int, diff string) {
		fn = fn.append(offset, npre, diff)
		data = append(data, ent{offset, npre, diff})
	}
	add(123, 2, "bar")
	add(456, 1, "foo")
	for _, e := range data {
		var npre int
		var diff []byte
		var off uint64
		npre, diff, off = fn.read()
		fn = fn[fLen(diff):]
		assert.T(t).This(npre).Is(e.npre)
		assert.T(t).This(string(diff)).Is(e.diff)
		assert.T(t).This(off).Is(e.offset)
	}
}

func TestFnodeInsert(*testing.T) {
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
		fwd := fnode{}
		for i, d := range data {
			fwd = fwd.insert(d, uint64(i), get)
			fwd.checkUpTo(i, data, get)
		}
		// fmt.Println("forward")
		// fwd.printLeafNode(get)
		assert.That(fwd.check(get) == len(data))

		// reverse
		rev := fnode{}
		for i := len(data) - 1; i >= 0; i-- {
			rev = rev.insert(data[i], uint64(i), get)
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
		for i := 0; i < nperms; i++ {
			rnd := fnode{}
			perm := rand.Perm(len(data))
			for _, j := range perm {
				rnd = rnd.insert(data[j], uint64(j), get)
			}
			assert.This(rnd).Is(fwd)
		}
	}
}

func build(data []string) fnode {
	b := fNodeBuilder{}
	for i, d := range data {
		b.Add(d, uint64(i), 1)
	}
	return b.Entries()
}

func (fn fnode) checkData(data []string, get func(uint64) string) {
	fn.checkUpTo(len(data)-1, data, get)
}

// checkUpTo is used during inserting.
// It checks that inserted keys are present
// and uninserted keys are not present.
func (fn fnode) checkUpTo(i int, data []string, get func(uint64) string) {
	n := fn.check(get)
	nn := 0
	for j, d := range data {
		if (d != "" && j <= i) != fn.contains(d, get) {
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

func TestFnodeInsertRandom(*testing.T) {
	const nData = 100
	var nGenerate = 1000
	var nShuffle = 20
	if testing.Short() {
		nGenerate = 1
		nShuffle = 4
	}
	var data = make([]string, nData)
	get := func(i uint64) string { return data[i] }
	for gi := 0; gi < nGenerate; gi++ {
		data = data[0:nData]
		randKey := str.UniqueRandomOf(1, 6, "abcdef")
		for di := 0; di < nData; di++ {
			data[di] = randKey()
		}
		for si := 0; si < nShuffle; si++ {
			rand.Shuffle(len(data),
				func(i, j int) { data[i], data[j] = data[j], data[i] })
			var fn fnode
			for i, d := range data {
				fn = fn.insert(d, uint64(i), get)
				// fe.checkUpTo(i, data, get)
			}
			fn.checkData(data, get)
		}
	}
}

func TestFnodeInsertWords(*testing.T) {
	data := words
	const nShuffle = 100
	get := func(i uint64) string { return data[i] }
	for si := 0; si < nShuffle; si++ {
		rand.Shuffle(len(data),
			func(i, j int) { data[i], data[j] = data[j], data[i] })
		var fn fnode
		for i, d := range data {
			fn = fn.insert(d, uint64(i), get)
			// fe.checkUpto(i, data, get)
		}
		fn.checkData(data, get)
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

func TestFnodeDelete(*testing.T) {
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
	var without, del fnode
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
	const ntimes = 100000
	for i := 0; i < ntimes; i++ {
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
		for i := 0; i < len(data); i++ {
			all := build(data)
			datawo = append(datawo[:0], data...)
			datawo[i] = ""
			// if i == 2 {
			// 	fmt.Println("===========================")
			// 	all.printLeafNode(get)
			// }
			// fmt.Println("---------------------------", i, data[i])

			b := fNodeBuilder{}
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

func compare(f, g fnode) string {
	fit := f.iter()
	git := g.iter()
	for {
		fok := fit.next()
		gok := git.next()
		if fok != gok {
			return "DIFFERENT lengths"
		}
		if !fok {
			return "" // ok
		}
		switch {
		case fit.offset != git.offset:
			return "DIFFERENT offsets"
		case fit.npre != git.npre:
			return "DIFFERENT npre"
		case !bytes.HasPrefix(git.known, fit.known):
			return "DIFFERENT known not prefix"
		}
	}
}

type slot struct {
	key string
	off uint64
}

func TestFnodeInsertDelete(*testing.T) {
	var ok bool
	var fn fnode
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
	const N = 1_000_000
	for i := 0; i < N; i++ {
		if rand.Intn(13)+1 > len(data) {
			// insert
			key := r()
			for ; dup(key); key = r() {
			}
			off := uint64(i)
			// print("insert", key, off)
			fn = fn.insert(key, off, get)
			data = append(data, slot{key: key, off: off})
		} else {
			// delete
			i := rand.Intn(len(data))
			// print("delete", data[i].key, data[i].off)
			fn, ok = fn.delete(data[i].off)
			assert.That(ok)
			data[i] = data[len(data)-1]
			data = data[:len(data)-1]
		}
		assert.This(fn.check(get)).Is(len(data))
		for _, s := range data {
			if !fn.contains(s.key, get) {
				print(data)
				fn.printLeafNode(get)
				panic("lookup failed " + s.key)
			}
		}
	}
}

var S1 []byte
var S2 []byte
var FN fnode

func BenchmarkFnode(b *testing.B) {
	get := func(i uint64) string { return words[i] }
	var fn fnode
	for i, d := range words {
		fn = fn.insert(d, uint64(i), get)
	}
	FN = fn

	for i := 0; i < b.N; i++ {
		iter := fn.iter()
		for iter.next() {
			S1 = iter.known
			S2 = iter.diff
		}
	}
}

func ExampleFnodeBuilderSplit() {
	var fb fNodeBuilder
	fb.Add("1234xxxx", 1234, 1)
	fb.Add("1235xxxx", 1235, 1)
	fb.Add("1299xxxx", 1299, 1)
	fb.Add("1300xxxx", 1300, 1)
	fb.Add("1305xxxx", 1305, 1)
	store := stor.HeapStor(8192)
	leftOff, splitKey := fb.Split(store)
	// assert.T(t).This(splitKey).Is("13")
	fmt.Println("splitKey", splitKey)
	fmt.Println("LEFT ---")
	readNode(store, leftOff).print()
	fmt.Println("RIGHT ---")
	fb.fe.print()

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

func ExampleFbmergeSplit() {
	var fb fNodeBuilder
	fb.Add("1234xxxx", 1234, 1)
	fb.Add("1235xxxx", 1235, 1)
	fb.Add("1299xxxx", 1299, 1)
	fb.Add("1300xxxx", 1300, 1)
	fb.Add("1305xxxx", 1305, 1)
	m := merge{node: fb.fe, modified: true}
	left, right, splitKey := m.split()
	// assert.T(t).This(splitKey).Is("13")
	fmt.Println("splitKey", splitKey)
	fmt.Println("LEFT ---")
	left.print()
	fmt.Println("RIGHT ---")
	right.print()

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
