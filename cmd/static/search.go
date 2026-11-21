// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

var revEdges [][]uint16
var queue []uint16
var tables []uint16
var prev []uint16

func printCounts() {
	nEdges := 0
	for _, refs := range edges {
		nEdges += len(refs)
	}
	fmt.Printf("tables: %d, ETA roots: %d, library names: %d, edges: %d\n",
		nTables, nRoots, nLibNames, nEdges)
}

// run finds the tables reachable from from, or the path from from to to.
func run(from, to string, minCount int, sortByCount bool) {
	if to == "counts" {
		printCountsFor(from, minCount, sortByCount)
		return
	}
	starts, ok := startsFor(from)
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown node %q\n", from)
		return
	}
	if to == "tables" {
		if from == "roots" {
			printTablesByRoot(starts)
			return
		}
		ids := bfsAll(starts)
		sortByName(ids)
		fmt.Printf("Tables (%d): %s\n", len(ids), names(ids))
		return
	}
	target, ok := nodeMap[to]
	if !ok || target.id == alias {
		fmt.Fprintf(os.Stderr, "unknown node %q\n", to)
		return
	}
	if from == "roots" {
		roots := rootsTo(target.id)
		sortByName(roots)
		fmt.Printf("Roots (%d): %s\n", len(roots), names(roots))
		return
	}
	path := pathFrom(starts, target.id)
	if path == nil {
		fmt.Printf("no path from %s to %s\n", from, to)
		return
	}
	fmt.Printf("Path (%d): %s\n", len(path), names(path))
}

// printEdges prints the outgoing and incoming edges for name.
func printEdges(name string) {
	out, in, ok := edgesFor(name)
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown node %q\n", name)
		return
	}
	fmt.Printf("Incoming (%d): %s\n", len(in), names(in))
	fmt.Printf("Outgoing (%d): %s\n", len(out), names(out))
}

// edgesFor returns the sorted outgoing and incoming edge ids for name.
func edgesFor(name string) (out, in []uint16, ok bool) {
	e, ok := nodeMap[name]
	if !ok || e.id == alias {
		return nil, nil, false
	}
	out = slices.Clone(edges[e.id])
	sortByName(out)
	in = slices.Clone(reverseEdges()[e.id])
	sortByName(in)
	return out, in, true
}

// printCountsFor prints the number of tables reachable from each root,
// or the number of roots that can reach each table. When minCount is not -1,
// entries with a count less than minCount are omitted. When sortByCount is set,
// entries are sorted by descending count, then name, instead of by name.
func printCountsFor(from string, minCount int, sortByCount bool) {
	type countEntry struct {
		id    uint16
		count int
	}
	var entries []countEntry
	switch from {
	case "roots":
		for id := nTables; id < nTables+nRoots; id++ {
			entries = append(entries, countEntry{uint16(id), len(bfs(uint16(id)))})
		}
	case "tables":
		counts := make([]int, nextId)
		for id := nTables; id < nTables+nRoots; id++ {
			for _, tbl := range bfs(uint16(id)) {
				counts[tbl]++
			}
		}
		for id := range nTables {
			entries = append(entries, countEntry{uint16(id), counts[id]})
		}
	default:
		fmt.Fprintln(os.Stderr, `counts from must be "roots" or "tables"`)
		return
	}
	entries = slices.DeleteFunc(entries, func(e countEntry) bool {
		return e.count < minCount
	})
	slices.SortFunc(entries, func(a, b countEntry) int {
		if sortByCount && a.count != b.count {
			return b.count - a.count
		}
		return strings.Compare(nodeVec[a.id], nodeVec[b.id])
	})
	for _, e := range entries {
		fmt.Printf("%s: %d\n", nodeVec[e.id], e.count)
	}
}

// printTablesByRoot prints the tables reachable from each root.
func printTablesByRoot(roots []uint16) {
	sortByName(roots)
	for _, id := range roots {
		ids := bfs(id)
		sortByName(ids)
		fmt.Printf("%s (%d): %s\n", nodeVec[id], len(ids), names(ids))
	}
}

func sortByName(ids []uint16) {
	slices.SortFunc(ids, func(a, b uint16) int {
		return strings.Compare(nodeVec[a], nodeVec[b])
	})
}

// startsFor returns the starting nodes for a from reference.
func startsFor(from string) ([]uint16, bool) {
	if from == "roots" {
		starts := make([]uint16, 0, nRoots)
		for id := nTables; id < nTables+nRoots; id++ {
			starts = append(starts, uint16(id))
		}
		return starts, true
	}
	e, ok := nodeMap[from]
	if !ok || e.id == alias {
		return nil, false
	}
	return []uint16{e.id}, true
}

// names translates ids to their names, comma separated.
func names(ids []uint16) string {
	strs := make([]string, len(ids))
	for i, id := range ids {
		strs[i] = nodeVec[id]
	}
	return strings.Join(strs, ", ")
}

// bfs returns the sorted ids of the tables reachable from start.
func bfs(start uint16) []uint16 {
	return bfsAll([]uint16{start})
}

// bfsAll returns the sorted ids of the tables reachable from any start.
func bfsAll(starts []uint16) []uint16 {
	stamp++
	queue = queue[:0]
	for _, start := range starts {
		if markSeen(start) {
			queue = append(queue, start)
		}
	}
	tables = tables[:0]
	for i := 0; i < len(queue); i++ {
		id := queue[i]
		if int(id) < nTables {
			tables = append(tables, id)
		}
		for _, next := range edges[id] {
			if markSeen(next) {
				queue = append(queue, next)
			}
		}
	}
	slices.Sort(tables)
	return tables
}

// pathTo returns the shortest path from start to target, or nil.
func pathTo(start, target uint16) []uint16 {
	return pathFrom([]uint16{start}, target)
}

// pathFrom returns the ids of the shortest (first BFS) path from any start to
// target, or nil if target is not reachable. The path includes both ends.
func pathFrom(starts []uint16, target uint16) []uint16 {
	stamp++
	queue = queue[:0]
	for _, start := range starts {
		if markSeen(start) {
			prev[start] = start
			queue = append(queue, start)
		}
	}
	for i := 0; i < len(queue); i++ {
		id := queue[i]
		if id == target {
			path := []uint16{id}
			for prev[id] != id {
				id = prev[id]
				path = append(path, id)
			}
			slices.Reverse(path)
			return path
		}
		for _, next := range edges[id] {
			if markSeen(next) {
				prev[next] = id
				queue = append(queue, next)
			}
		}
	}
	return nil
}

// reverseEdges returns the transposed graph, building it on first use.
func reverseEdges() [][]uint16 {
	if revEdges != nil {
		return revEdges
	}
	revEdges = make([][]uint16, len(edges))
	for id, refs := range edges {
		for _, ref := range refs {
			revEdges[ref] = append(revEdges[ref], uint16(id))
		}
	}
	return revEdges
}

// rootsTo returns the ids of the ETA roots that can reach target.
func rootsTo(target uint16) []uint16 {
	rev := reverseEdges()
	stamp++
	queue = append(queue[:0], target)
	markSeen(target)
	var roots []uint16
	for i := 0; i < len(queue); i++ {
		id := queue[i]
		if int(id) >= nTables && int(id) < nTables+nRoots {
			roots = append(roots, id)
		}
		for _, from := range rev[id] {
			if markSeen(from) {
				queue = append(queue, from)
			}
		}
	}
	return roots
}
