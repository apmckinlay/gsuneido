// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"sort"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

// makeUnodeFromNodeKnown builds a unode from a node's iterated keys,
// using the same strings node.search compares against (it.known).
func makeUnodeFromNodeKnown(nd node) unode {
	var u unode
	it := nd.iter()
	for it.next() {
		u = append(u, slot{key: string(it.known), off: it.offset})
	}
	return u
}

func TestUnodeSearchEquivalent(t *testing.T) {
	datas := []string{
		"a b c d",
		"xa xb xc xd",
		"a ab abc abcd",
		"a baac baacb c cbaba",
		"a ant anti ants b bun bunn bunnies bunny buns ca cat cats",
		"a aa aaa ab abc b ba bb bba bbb bbc bc c",
		"1000 1001 1002 1003",
	}

	for _, s := range datas {
		data := strings.Fields(s)
		nd := build(data)
		u := makeUnodeFromNodeKnown(nd)

		// Collect test keys: all knowns and all full data values,
		// plus a key larger than the last to check the upper boundary.
		var keys []string
		it := nd.iter()
		for it.next() {
			keys = append(keys, string(it.known))
		}
		keys = append(keys, data...)
		if len(data) > 0 {
			keys = append(keys, data[len(data)-1]+"\xff")
		}

		for _, k := range keys {
			nOff := nd.search(k)
			uOff := u.search(k)
			assert.T(t).This(uOff).Is(nOff)
		}
	}
}

func TestUnodeSearchEquivalentWords(t *testing.T) {
	data := append([]string(nil), words...)
	sort.Strings(data)
	nd := build(data)
	u := makeUnodeFromNodeKnown(nd)

	// Use all knowns and a sampling of full keys, plus an upper boundary key.
	var keys []string
	it := nd.iter()
	for it.next() {
		keys = append(keys, string(it.known))
	}
	for i := 0; i < len(data); i += 10 {
		keys = append(keys, data[i])
	}
	if len(data) > 0 {
		keys = append(keys, data[len(data)-1]+"\xff")
	}

	for _, k := range keys {
		nOff := nd.search(k)
		uOff := u.search(k)
		assert.T(t).This(uOff).Is(nOff)
	}
}

// Seek -------------------------------------------------------------

// buildNode constructs a node with strictly increasing keys and
// monotonically increasing offsets starting at 1.
func buildNode(keys []string) (node, unode) {
	var nb nodeBuilder
	off := uint64(1)
	for _, k := range keys {
		nb.Add(k, off, embedAll)
		off++
	}
	nd := nb.Entries()
	u := nd.toUnode()
	return nd, u
}

func offsetsEqual(t *testing.T, nd node, u unode, q string) {
	t.Helper()
	ni := nd.seek(q)
	ui := u.seek(q)
	assert.This(ui.off()).Is(ni.off())
	// eof should be consistent at the last entry
	assert.This(ui.eof()).Is(ni.eof())
}

func searchEqual(t *testing.T, nd node, u unode, q string) {
	t.Helper()
	assert.This(u.search(q)).Is(nd.search(q))
}

func TestUnodeSeekMatchesNodeSeek_Basic(t *testing.T) {
	// Keys include empty first-known, and a variety of prefixes
	keys := []string{"", "ant", "bat", "cat", "dog", "z", "zz"}
	nd, u := buildNode(keys)

	queries := []string{
		"", "a", "an", "ant", "antz",
		"b", "bat", "bb",
		"c", "cat", "caz",
		"d", "dog", "do~",
		"y", "z", "zz", "zzz", "~",
	}

	for _, q := range queries {
		offsetsEqual(t, nd, u, q)
		searchEqual(t, nd, u, q)
	}
}

func TestUnodeSeekMatchesNodeSeek_SingleEntry(t *testing.T) {
	keys := []string{""}
	nd, u := buildNode(keys)
	queries := []string{"", "a", "~"}

	for _, q := range queries {
		offsetsEqual(t, nd, u, q)
		searchEqual(t, nd, u, q)
	}
}

func TestUnodeSeekMatchesNodeSeek_TwoEntries(t *testing.T) {
	keys := []string{"", "k1"}
	nd, u := buildNode(keys)
	queries := []string{"", "j", "k", "k1", "k1x", "~"}

	for _, q := range queries {
		offsetsEqual(t, nd, u, q)
		searchEqual(t, nd, u, q)
	}
}

func TestUnodeSeekMatchesNodeSeek_Empty(t *testing.T) {
	var nd node
	u := nd.toUnode()
	queries := []string{"", "a", "~"}

	for _, q := range queries {
		ni := nd.seek(q)
		ui := u.seek(q)
		assert.This(ui.off()).Is(ni.off())
		// both should report not eof when empty since there's no last slot
		assert.This(ui.eof()).Is(ni.eof())
		searchEqual(t, nd, u, q)
	}
}

// Search benchmark

var BenchOff uint64

func buildNodeAndUnodeFromWords() (node, unode, []string) {
	// Use a sorted copy of words to satisfy builder requirements.
	data := append([]string(nil), words...)
	sort.Strings(data)
	nd := build(data)
	u := makeUnodeFromNodeKnown(nd)

	// Use the same comparison keys node.search uses (knowns).
	var keys []string
	it := nd.iter()
	for it.next() {
		keys = append(keys, string(it.known))
	}
	return nd, u, keys
}

func BenchmarkNodeSearch(b *testing.B) {
	nd, _, keys := buildNodeAndUnodeFromWords()
	i := 0
	for b.Loop() {
		BenchOff = nd.search(keys[i%len(keys)])
		i++
	}
}

func BenchmarkUnodeSearch(b *testing.B) {
	_, u, keys := buildNodeAndUnodeFromWords()
	i := 0
	for b.Loop() {
		BenchOff = u.search(keys[i%len(keys)])
		i++
	}
}

// toUnode benchmarks

func createBenchmarkNode() node {
	// Use sorted words to satisfy nodeBuilder requirements
	sortedWords := append([]string(nil), words...)
	sort.Strings(sortedWords)

	var nb nodeBuilder
	for i, w := range sortedWords {
		nb.Add(w, uint64(i+1), embedAll)
	}
	return nb.Entries()
}

func BenchmarkToUnode(b *testing.B) {
	nd := createBenchmarkNode()
	for b.Loop() {
		_ = nd.toUnode()
	}
}
