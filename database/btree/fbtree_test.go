// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
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
	var fe fEntries
	var data []ent
	add := func(offset uint64, npre int, diff string) {
		fe = fAppend(fe, offset, npre, diff)
		data = append(data, ent{offset, npre, diff})
	}
	add(123, 2, "bar")
	add(456, 1, "foo")
	for _, e := range data {
		var offset uint64
		var npre int
		var diff string
		fe, offset, npre, diff = fRead(fe)
		Assert(t).That(offset, Equals(e.offset))
		Assert(t).That(npre, Equals(e.npre))
		Assert(t).That(diff, Equals(e.diff))
	}
}

func TestInsert(*testing.T) {
	datas := []string{
		"a b c d",
		"xa xb xc xd",
		"a ab abc abcd",
		"ant ants bun bunnies bunny buns cat a anti b bunn ca cats",
		"bbb bbc abc aa ab bc c aaa ba bba bb b a",
	}
	for _, s := range datas {
		data := strings.Fields(s)
		// print()
		// print()
		get := func(i uint64) string { return data[i] }
		var fe fEntries
		// forward
		for i, d := range data {
			fe = fe.insert(d, uint64(i), get)
			fe.checkUpTo(i, data, get)
		}
		// print("------------------------")
		// fe.printRaw(get)
		// print("------------------------")
		verify.That(fe.check() == len(data))
		// reverse
		str.ListReverse(data)
		fe = nil
		for i, d := range data {
			fe = fe.insert(d, uint64(i), get)
			fe.checkUpTo(i, data, get)
		}
		// print("------------------------")
		// fe.printRaw(get)
		// print("------------------------")
	}
}

func TestRandom(*testing.T) {
	const nData = 100
	var nGenerate = 40
	var nShuffle = 40
	if testing.Short() {
		nGenerate = 1
		nShuffle = 4
	}
	var data = make([]string, nData)
	defer func() {
		if e := recover(); e != nil {
			print(data)
			panic(e)
		}
	}()
	get := func(i uint64) string { return data[i] }
	for gi := 0; gi < nGenerate; gi++ {
		data = data[0:nData]
		for di := 0; di < nData; di++ {
			data[di] = str.RandomOf(1, 6, "abcdef")
		}
		data = str.ListUnique(data)
		// print(data)
		for si := 0; si < nShuffle; si++ {
			rand.Shuffle(len(data),
				func(i, j int) { data[i], data[j] = data[j], data[i] })
			var fe fEntries
			for i, d := range data {
				fe = fe.insert(d, uint64(i), get)
				// fe.checkUpTo(i, data, get)
			}
			fe.checkData(data, get)
			// print("------------------------")
			// fe.printRaw(get)
			// print("------------------------")
			// if si == 0 {
			// 	fe.stats()
			// }
		}

	}
}

func TestFileData(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping test in short mode")
    }
	data := fileData("../../../bizpartnername.txt")
	get := func(i uint64) string { return data[i] }
	const nShuffle = 4
	for si := 0; si < nShuffle; si++ {
		rand.Shuffle(len(data),
			func(i, j int) { data[i], data[j] = data[j], data[i] })
		var fe fEntries
		for i, d := range data {
			fe = fe.insert(d, uint64(i), get)
			// fe.checkData(i, data, get)
		}
		fe.checkData(data, get)
		// print("------------------------")
		// fe.printRaw(get)
		// print("------------------------")
		if si == 0 {
			// fe.printRaw(get)
			// fe.stats()
		}
	}
}

func fileData(filename string) []string {
	file, _ := os.Open(filename)
	defer file.Close()
	data := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}
	return data
}

//-------------------------------------------------------------------

func (fe fEntries) stats() {
	n := fe.check()
	avg := float32(len(fe)-7*n) / float32(n)
	print("    n", n, "len", len(fe), "avg", avg)
}

func (fe fEntries) checkData(data []string, get func(uint64) string) {
	n := len(data)
	fe.checkUpTo(n-1, data, get)
}

// checkUpTo is used during inserting.
// It checks that inserted keys are present
// and uninserted keys are not present.
func (fe fEntries) checkUpTo(i int, data []string, get func(uint64) string) {
	verify.That(fe.check() == i+1)
	for j, d := range data {
		if j <= i != fe.contains(d, get) {
			panic("can't find " + d)
		}
	}
}

func (fe fEntries) check() int {
	n := 0
	prev := ""
	it := fe.Iter()
	for it.next() {
		if it.known < prev {
			panic("fEntries out of order")
		}
		prev = it.known
		n++
	}
	return n
}

func (fe fEntries) print() {
	it := fe.Iter()
	for it.next() {
		print(it.offset, it.known)
	}
}

func (fe fEntries) printRaw(get func(uint64) string) {
	it := fe.Iter()
	for it.next() {
		print(it.fi, "{", it.offset, it.npre, it.diff, "}", it.known, "=", get(it.offset))
	}
}

func print(args ...interface{}) {
	for i, x := range args {
		switch x := x.(type) {
		case string:
			if x == "" {
				args[i] = "'" + x + "'"
			}
		}
	}
	fmt.Println(args...)
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
