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
	var fe fNode
	var data []ent
	add := func(offset uint64, npre int, diff string) {
		fe = fAppend(fe, offset, npre, diff)
		data = append(data, ent{offset, npre, diff})
	}
	add(123, 2, "bar")
	add(456, 1, "foo")
	for _, e := range data {
		var npre int
		var diff string
		fe, npre, diff = fRead(fe)
		Assert(t).That(npre, Equals(e.npre))
		Assert(t).That(diff, Equals(e.diff))
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
		fe := fNode{}
		for i, d := range data {
			fe,_ = fe.insert(d, uint64(i), get)
			fe.checkUpTo(i, data, get)
		}
		verify.That(fe.check() == len(data))
		// reverse
		str.ListReverse(data)
		fe = nil
		for i, d := range data {
			fe,_ = fe.insert(d, uint64(i), get)
			fe.checkUpTo(i, data, get)
		}
		// builder
		fe = build(data)
		fe.checkData(data, get)
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
	var nGenerate = 40
	var nShuffle = 40
	if testing.Short() {
		nGenerate = 1
		nShuffle = 4
	}
	var data = make([]string, nData)
	get := func(i uint64) string { return data[i] }
	for gi := 0; gi < nGenerate; gi++ {
		data = data[0:nData]
		for di := 0; di < nData; di++ {
			data[di] = str.RandomOf(1, 6, "abcdef")
		}
		data = str.ListUnique(data)
		for si := 0; si < nShuffle; si++ {
			rand.Shuffle(len(data),
				func(i, j int) { data[i], data[j] = data[j], data[i] })
			var fe fNode
			for i, d := range data {
				fe,_ = fe.insert(d, uint64(i), get)
				// fe.checkUpTo(i, data, get)
			}
			fe.checkData(data, get)
		}

	}
}

func TestWords(*testing.T) {
	data := words
	const nShuffle = 100
	get := func(i uint64) string { return data[i] }
	for si := 0; si < nShuffle; si++ {
		rand.Shuffle(len(data),
			func(i, j int) { data[i], data[j] = data[j], data[i] })
		var fe fNode
		for i, d := range data {
			fe,_ = fe.insert(d, uint64(i), get)
			// fe.checkUpto(i, data, get)
		}
		fe.checkData(data, get)
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
