// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func useNodes(t *testing.T, nodes map[string]entry) {
	t.Helper()
	savedMap, savedNext, savedStamps := nodeMap, nextId, stamps
	t.Cleanup(func() { nodeMap, nextId, stamps = savedMap, savedNext, savedStamps })
	nodeMap = nodes
	nextId = 0
	for _, e := range nodes {
		nextId = max(nextId, e.id+1)
	}
	stamps = make([]uint32, nextId)
	prev = make([]uint16, nextId)
}

func TestParseCommand(t *testing.T) {
	from, to, minCount, sortByCount, ok := parseCommand([]string{"roots", "to", "tables"})
	assert.T(t).That(ok)
	assert.T(t).This(from).Is("roots")
	assert.T(t).This(to).Is("tables")
	assert.T(t).This(minCount).Is(-1)
	assert.T(t).That(!sortByCount)

	from, to, minCount, sortByCount, ok = parseCommand([]string{"roots", "to", "counts^10"})
	assert.T(t).That(ok)
	assert.T(t).This(from).Is("roots")
	assert.T(t).This(to).Is("counts")
	assert.T(t).This(minCount).Is(10)
	assert.T(t).That(sortByCount)

	from, to, minCount, sortByCount, ok = parseCommand([]string{"roots", "to", "counts^"})
	assert.T(t).That(ok)
	assert.T(t).This(from).Is("roots")
	assert.T(t).This(to).Is("counts")
	assert.T(t).This(minCount).Is(-1)
	assert.T(t).That(sortByCount)

	_, _, _, _, ok = parseCommand([]string{"roots", "to", "counts^bad"})
	assert.T(t).That(!ok)

	from, to, minCount, sortByCount, ok = parseCommand([]string{"roots", "to", "tables^10"})
	assert.T(t).That(ok)
	assert.T(t).This(from).Is("roots")
	assert.T(t).This(to).Is("tables^10")
	assert.T(t).This(minCount).Is(-1)
	assert.T(t).That(!sortByCount)
}

func TestScanIdentifiers(t *testing.T) {
	useNodes(t, map[string]entry{"Foo": {id: 7}, "Bar": {id: 9}, "Baz": {id: 11}})

	src := `function () {
		Foo(Bar)
		Baz
		Foo
	}`
	assert.T(t).This(scanIdentifiers("Test", src)).Is([]uint16{7, 9, 11})
	assert.T(t).This(scanIdentifiers("Test", src)).Is([]uint16{7, 9, 11})
}

func TestScanIdentifiersClass(t *testing.T) {
	useNodes(t, map[string]entry{"Foo": {id: 7}, "Bar": {id: 9}})

	assert.T(t).This(scanIdentifiers("Test", `Baz {
		M() { Foo(Bar) }
	}`)).Is([]uint16{7, 9})
}

func TestBFS(t *testing.T) {
	useNodes(t, map[string]entry{
		"t0": {id: 0}, "t1": {id: 1}, "t2": {id: 2}, "t3": {id: 3}, "eta": {id: 4}, "lib": {id: 5},
	})
	savedTables, savedEdges := nTables, edges
	t.Cleanup(func() { nTables, edges = savedTables, savedEdges })
	nTables = 4
	edges = [][]uint16{
		0: nil,
		1: nil,
		2: nil,
		3: nil,
		4: {5, 0},
		5: {4, 2, 1},
	}
	assert.T(t).This(bfs(4)).Is([]uint16{0, 1, 2})
}

func TestPathTo(t *testing.T) {
	useNodes(t, map[string]entry{
		"t0": {id: 0}, "t1": {id: 1}, "t2": {id: 2}, "t3": {id: 3}, "eta": {id: 4}, "lib": {id: 5},
	})
	savedTables, savedEdges := nTables, edges
	t.Cleanup(func() { nTables, edges = savedTables, savedEdges })
	nTables = 4
	edges = [][]uint16{
		0: nil,
		1: nil,
		2: nil,
		3: nil,
		4: {5, 0},
		5: {4, 2, 1},
	}
	assert.T(t).This(pathTo(4, 2)).Is([]uint16{4, 5, 2})
	assert.T(t).This(pathTo(4, 4)).Is([]uint16{4})
	assert.T(t).This(pathTo(4, 3)).Is(nil)
}

func TestEdgesFor(t *testing.T) {
	useNodes(t, map[string]entry{
		"a": {id: 0}, "b": {id: 1}, "c": {id: 2}, "d": {id: 3}, "Alias": {id: alias},
	})
	useNodeVec(t, []string{"a", "b", "c", "d"}, 4, 0)
	savedTables, savedEdges, savedRev := nTables, edges, revEdges
	t.Cleanup(func() { nTables, edges, revEdges = savedTables, savedEdges, savedRev })
	nTables = 4
	revEdges = nil
	edges = [][]uint16{
		0: {1, 2},
		1: {3},
		2: {1},
		3: nil,
	}

	out, in, ok := edgesFor("b")
	assert.T(t).That(ok)
	assert.T(t).This(out).Is([]uint16{3})
	assert.T(t).This(in).Is([]uint16{0, 2})

	_, _, ok = edgesFor("missing")
	assert.T(t).That(!ok)
	_, _, ok = edgesFor("Alias")
	assert.T(t).That(!ok)
}

func useNodeVec(t *testing.T, names []string, nTbl, nEta int) {
	t.Helper()
	savedVec, savedTables, savedEta, savedLibs := nodeVec, nTables, nRoots, nLibNames
	t.Cleanup(func() {
		nodeVec, nTables, nRoots, nLibNames = savedVec, savedTables, savedEta, savedLibs
	})
	nodeVec = names
	nTables, nRoots = nTbl, nEta
	nLibNames = len(names) - nTbl - nEta
}

func TestSetupVariants(t *testing.T) {
	useNodes(t, map[string]entry{
		"alpha": {id: 0}, "beta": {id: 1}, "Gamma": {id: 2}, "Delta": {id: 3},
		"Rule_alpha": {id: 4}, "Field_beta": {id: 5}, "GammaControl": {id: 6},
		"DeltaFormat": {id: 7}, "Rule_missing": {id: 8}, "Rule_late": {id: 9},
		"Field_alpha": {id: 10}, "Rule_Upper": {id: 11}, "fooControl": {id: 12},
	})
	// ETA nodes sit between tables and library names, so the scan range must
	// skip the ETA section without skipping any trailing library name.
	useNodeVec(t, []string{
		"T1", "T2",
		"/eta/FooControl",
		"Rule_alpha", "Rule_alpha", "Field_beta", "GammaControl", "DeltaFormat",
		"Rule_missing", "Field_alpha", "Rule_Upper", "fooControl", "Rule_late",
	}, 2, 1)

	setupVariants()

	assert.T(t).This(nodeMap["alpha"]).Is(entry{id: 0, bits: hasRule | hasField})
	assert.T(t).This(nodeMap["beta"]).Is(entry{id: 1, bits: hasField})
	assert.T(t).This(nodeMap["Gamma"]).Is(entry{id: 2, bits: hasControl})
	assert.T(t).This(nodeMap["Delta"]).Is(entry{id: 3, bits: hasFormat})
	assert.T(t).This(nodeMap["missing"]).Is(entry{id: alias, bits: hasRule})
	assert.T(t).This(nodeMap["late"]).Is(entry{id: alias, bits: hasRule})

	_, ok := nodeMap["Upper"]
	assert.T(t).That(!ok) // prefix requires a lowercase basename
	_, ok = nodeMap["foo"]
	assert.T(t).That(!ok) // suffix requires an uppercase basename
	_, ok = nodeMap["/eta/Foo"]
	assert.T(t).That(!ok) // ETA nodes are not library items
	_, ok = nodeMap[""]
	assert.T(t).That(!ok) // a bare "Control" has no basename
}

func TestScanIdentifiersIgnoresVariants(t *testing.T) {
	useNodes(t, map[string]entry{
		"Foo":        {id: 0, bits: hasControl | hasFormat},
		"FooControl": {id: 1},
		"FooFormat":  {id: 2},
		"bar":        {id: 3, bits: hasRule | hasField},
		"Rule_bar":   {id: 4},
		"Field_bar":  {id: 5},
	})
	assert.T(t).This(scanIdentifiers("Test", "function () { Foo }")).Is([]uint16{0})
}

func TestScanIdentifiersIgnoresAlias(t *testing.T) {
	useNodes(t, map[string]entry{
		"Foo":        {id: alias, bits: hasControl},
		"FooControl": {id: 0},
	})
	assert.T(t).This(scanIdentifiers("Test", "function () { Foo }")).Is(nil)
	assert.T(t).This(scanIdentifiers("Test", "function () { FooControl }")).Is([]uint16{0})
}

func TestScanIdentifiersExcludesNames(t *testing.T) {
	useNodes(t, map[string]entry{
		"Params": {id: 0},
		"Foo":    {id: 1},
	})
	assert.T(t).This(scanIdentifiers("Test", "function () { Params; Foo }")).Is([]uint16{1})
}

func TestAddEdges(t *testing.T) {
	savedEdges := edges
	t.Cleanup(func() { edges = savedEdges })
	edges = make([][]uint16, 2)

	addEdges(0, []uint16{1, 2})
	addEdges(0, []uint16{2, 3})
	assert.T(t).This(edges[0]).Is([]uint16{1, 2, 3})
}

func TestStartsForAlias(t *testing.T) {
	useNodes(t, map[string]entry{
		"Foo": {id: alias},
		"Bar": {id: 1},
	})
	_, ok := startsFor("Foo")
	assert.T(t).That(!ok)
	starts, ok := startsFor("Bar")
	assert.T(t).That(ok)
	assert.T(t).This(starts).Is([]uint16{1})
}

func TestScanStringTables(t *testing.T) {
	useNodes(t, map[string]entry{
		"customer": {id: 0}, "orders": {id: 1}, "NotATable": {id: 2},
	})
	savedTables := nTables
	t.Cleanup(func() { nTables = savedTables })
	nTables = 2

	src := `function () {
		Query("select * from customer join orders where x = NotATable")
	}`
	assert.T(t).This(scanIdentifiers("Test", src)).Is([]uint16{0, 1})
}

func TestScanStringInContainer(t *testing.T) {
	useNodes(t, map[string]entry{
		"wo_workdesc": {id: 0}, "Browse": {id: 1},
	})
	savedTables := nTables
	t.Cleanup(func() { nTables = savedTables })
	nTables = 1

	assert.T(t).This(scanIdentifiers("Test",
		`function () { #(Browse "wo_workdesc") }`)).Is([]uint16{0})

	assert.T(t).This(scanIdentifiers("Test",
		`#(Browse "wo_workdesc")`)).Is([]uint16{0})

	assert.T(t).This(scanIdentifiers("Test",
		`Foo { Data: #(Browse "wo_workdesc") }`)).Is([]uint16{0})
}

func TestScanStringControlFormat(t *testing.T) {
	useNodes(t, map[string]entry{
		"Foo":        {id: 0, bits: hasControl | hasFormat},
		"FooControl": {id: 1},
		"FooFormat":  {id: 2},
		"Bar":        {id: 3},
	})
	savedTables := nTables
	t.Cleanup(func() { nTables = savedTables })
	nTables = 0

	assert.T(t).This(scanIdentifiers("Test", `function () { "Foo Bar" }`)).
		Is([]uint16{1, 2})
}

func TestScanStringExcludesVariant(t *testing.T) {
	useNodes(t, map[string]entry{
		"Foo":        {id: 0, bits: hasControl | hasFormat},
		"FooControl": {id: 1},
		"FooFormat":  {id: 2},
	})
	savedTables := nTables
	t.Cleanup(func() { nTables = savedTables })
	nTables = 0
	excludeNames["FooControl"] = struct{}{}
	t.Cleanup(func() { delete(excludeNames, "FooControl") })

	assert.T(t).This(scanIdentifiers("Test", `function () { "Foo" }`)).
		Is([]uint16{2})
}

func TestSkipName(t *testing.T) {
	assert.T(t).That(skipName("foo.js"))
	assert.T(t).That(skipName("foo.css"))
	assert.T(t).That(skipName("Foo__webgui"))
	assert.T(t).That(skipName("Foo__webguiBar"))
	assert.T(t).That(skipName("FooTest"))
	assert.T(t).That(skipName("Test"))
	assert.T(t).That(!skipName("Foo"))
	assert.T(t).That(!skipName("js"))
	assert.T(t).That(!skipName("webgui"))
	assert.T(t).That(!skipName("TestFoo"))
}
