// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"math/rand"
	"sort"
	"strings"
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/verify"
)

func TestFAppendRead(t *testing.T) {
	type ent struct {
		offset uint64
		npre   int
		diff   string
	}
	var fn fNode
	var data []ent
	add := func(offset uint64, npre int, diff string) {
		fn = fAppend(fn, offset, npre, diff)
		data = append(data, ent{offset, npre, diff})
	}
	add(123, 2, "bar")
	add(456, 1, "foo")
	for _, e := range data {
		var npre int
		var diff string
		var off uint64
		npre, diff, off = fRead(fn)
		fn = fn[fLen(diff):]
		Assert(t).That(npre, Equals(e.npre))
		Assert(t).That(diff, Equals(e.diff))
		Assert(t).That(off, Equals(e.offset))
	}
}

func TestFnodeInsert(*testing.T) {
	datas := []string{
		"a b c d",
		"xa xb xc xd",
		"a ab abc abcd",
		"ant ants bun bunnies bunny buns cat a anti b bunn ca cats",
		"bbb bbc abc aa ab bc c aaa ba bba bb b a",
	}
	for _, s := range datas {
		data := strings.Fields(s)
		get := func(i uint64) string { return data[i] }
		// forward
		fn := fNode{}
		for i, d := range data {
			fn, _ = fn.insert(d, uint64(i), get)
			fn.checkUpTo(i, data, get)
		}
		verify.That(fn.check() == len(data))
		// reverse
		str.ListReverse(data)
		fn = nil
		for i, d := range data {
			fn, _ = fn.insert(d, uint64(i), get)
			fn.checkUpTo(i, data, get)
		}
		// builder
		fn = build(data)
		fn.checkData(data, get)
	}
}

func build(data []string) fNode {
	sort.Strings(data)
	b := fNodeBuilder{}
	for i, d := range data {
		b.Add(d, uint64(i))
	}
	return b.Entries()
}

func TestFnodeRandom(*testing.T) {
	const nData = 100
	var nGenerate = 20
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
			var fn fNode
			for i, d := range data {
				fn, _ = fn.insert(d, uint64(i), get)
				// fe.checkUpTo(i, data, get)
			}
			fn.checkData(data, get)
		}
	}
}

func TestDelete(t *testing.T) {
	var fn fNode
	const nData = 8 + 32
	var data = make([]string, nData)
	get := func(i uint64) string { return data[i] }
	randKey := str.UniqueRandomOf(1, 6, "abcdef")
	for i := 0; i < nData; i++ {
		data[i] = randKey()
	}
	sort.Strings(data)
	for i := 0; i < len(data); i++ {
		fn, _ = fn.insert(data[i], uint64(i), get)
	}
	// fn.printLeafNode(get)

	var ok bool

	// delete at end, simplest case
	for i := 0; i < 8; i++ {
		fn, ok = fn.delete(uint64(len(data) - 1))
		Assert(t).True(ok)
		data = data[:len(data)-1]
		fn.checkData(data, get)
	}
	// print("================================")
	// fn.printLeafNode(get)

	// delete at start
	const nStart = 8
	for i := 0; i < nStart; i++ {
		fn, ok = fn.delete(uint64(i))
		Assert(t).True(ok)
		data[i] = ""
		fn.checkData(data, get)
	}
	// print("================================")
	// fn.printLeafNode(get)

	for i := 0; i < len(data)-nStart; i++ {
		off := rand.Intn(len(data))
		for data[off] == "" {
			off = (off + 1) % len(data)
		}
		// print("================================ delete", data[off])
		fn, ok = fn.delete(uint64(off))
		Assert(t).True(ok)
		// fn.printLeafNode(get)
		data[off] = ""
		fn.checkData(data, get)
	}
}

func TestDelete2(t *testing.T) {
	data := []string{"a", "b", "c", "d", "e"}
	get := func(i uint64) string { return data[i] }
	var fn fNode
	for i := 0; i < len(data); i++ {
		fn, _ = fn.insert(data[i], uint64(i), get)
	}
	// fn.printLeafNode(get)

	var ok bool
	for i := 1; i < len(data); i++ {
		fn, ok = fn.delete(uint64(i))
		Assert(t).True(ok)
		// print("================================")
		// fn.printLeafNode(get)
		data[i] = ""
		fn.checkData(data, get)
	}
}

func TestWords(*testing.T) {
	data := words
	const nShuffle = 100
	get := func(i uint64) string { return data[i] }
	for si := 0; si < nShuffle; si++ {
		rand.Shuffle(len(data),
			func(i, j int) { data[i], data[j] = data[j], data[i] })
		var fn fNode
		for i, d := range data {
			fn, _ = fn.insert(d, uint64(i), get)
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
