// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"math"

	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/util/assert"
)

const smallRoot = 8 // ???

// RangeFrac returns a number from 0 to 1
// giving the fraction of the keys in the range >=org <end
func (bt *btree) RangeFrac(org, end string, _ int) float64 {
	if org >= end {
		return 0
	}
	if org == ixkey.Min && end == ixkey.Max {
		return 1
	}
	if bt.count == 0 {
		return .5 // ???
	}
	return bt.rangeFrac(org, end)
}

// rangeFrac descends the tree for both ends of the range in parallel
// until the descent diverges to different nodes, then calculates the fraction.
// For small ranges this will usually end up at the leaves
// and will give an exact result.
// For large ranges that diverge within the tree, more estimation is involved,
// but the result should still be within 2%
// NOTE: depends on btree.count being correct
func (bt *btree) rangeFrac(org, end string) float64 {
	_ = t && trace("=== rangeFrac", org, end)
	nkeys := bt.count

	// Descend in parallel until we diverge
	off := bt.root
	var fanout float64 // average fanout calculated from nrecs
	frac := 0.0
	div := 1.0
	atRoot := true
	for level := 0; level <= bt.treeLevels; level++ {
		nd := bt.readNode(level, off)
		n := nd.noffs()

		if level == 0 {
			// Calculate average fanout from nrecs
			// Subtract .5 from root noffs to allow for smaller rightmost node
			if bt.treeLevels == 1 {
				fanout = float64(nkeys) / (float64(n) - 0.5)
			} else if bt.treeLevels > 1 {
				// fanout = (nrecs / rootNoffs) ^ (1 / (treeLevels - 1))
				fanout = math.Pow(float64(nkeys)/(float64(n)-0.5), 1.0/float64(bt.treeLevels))
			}
		}

		orgPos, orgNext := search(nd, org)
		endPos, endNext := search(nd, end)

		if level == 0 && bt.treeLevels > 0 && n <= smallRoot {
			// if the root is small, use the next level
			nd, orgPos, orgNext, endPos, endNext =
				bt.fattenRoot(nd, org, orgPos, end, endPos)
			n = nd.noffs()
			level++
			// recalculate fanout given the new count from the root children
			fanout = math.Pow(float64(nkeys)/(float64(n)), 1.0/float64(bt.treeLevels-1))
		}

		if level == bt.treeLevels {
			// org and end are in the same leaf node, exact result
			return float64(endPos-orgPos) / float64(nkeys)
		}
		if orgPos != endPos {
			// Range spans multiple children at this level

			if level+1 == bt.treeLevels {
				// Next level is leaves - use nrecs for more accuracy
				return bt.leafRangeFrac(nd, orgPos, orgNext, org, endPos, endNext, end, fanout)
			}

			// diverged to two different tree nodes
			spread := endPos - orgPos
			if spread > 75 { // ???
				// positions are far enough apart, estimate based on current node
				return float64(spread) / float64(n) / div
			}
			// children are close, calculate using next level
			orgFrac := frac + float64(orgPos)/float64(n)/div
			endFrac := frac + float64(endPos)/float64(n)/div
			if atRoot {
				div = float64(n)
			} else {
				div *= fanout
			}

			orgChild := bt.readNode(level+1, orgNext)
			i, _ := search(orgChild, org)
			orgFrac += float64(i) / float64(orgChild.noffs()) / div

			endChild := bt.readNode(level+1, endNext)
			i, _ = search(endChild, end)
			endFrac += float64(i) / float64(endChild.noffs()) / div

			return endFrac - orgFrac
		}
		// Still in same child, descend further
		off = orgNext
		frac += float64(orgPos) / float64(n) / div
		if atRoot {
			atRoot = false
			div = float64(n)
		} else {
			div *= fanout
		}
	}
	panic(assert.ShouldNotReachHere())
}

func (bt *btree) fattenRoot(root node, org string, orgPos int,
	end string, endPos int) (fatRoot, int, uint64, int, uint64) {
	var m, k, orgI, endI int
	var orgNext, endNext uint64
	rn := root.noffs()
	nodes := make([]node, rn)
	for ni := range rn {
		node := bt.readNode(1, root.offset(ni))
		nodes[ni] = node
		no := node.noffs()
		m += no
		if ni < orgPos {
			orgI += no
		} else if ni == orgPos {
			k, orgNext = search(node, org)
			orgI += k
		}
		if ni < endPos {
			endI += no
		} else if ni == endPos {
			k, endNext = search(node, end)
			endI += k
		}
	}
	return fatRoot{nodes: nodes, no: m}, orgI, orgNext, endI, endNext
}

// fatRoot is a virtual node of the children of a small root
type fatRoot struct {
	nodes []node
	no    int
}

var _ node = fatRoot{}

func (fr fatRoot) offset(i int) uint64 {
	for _, nd := range fr.nodes {
		n := nd.noffs()
		if i < n {
			return nd.offset(i)
		}
		i -= n
	}
	panic(assert.ShouldNotReachHere())
}

func (fr fatRoot) noffs() int {
	return fr.no
}

func (fr fatRoot) size() int {
	panic(assert.ShouldNotReachHere())
}

// leafRangeFrac calculates the fraction when the range spans leaf nodes
// under the same parent tree node
func (bt *btree) leafRangeFrac(parent node, orgPos int, orgOff uint64, org string,
	endPos int, endOff uint64, end string, fanout float64) float64 {

	const maxToRead = 2 // maximum number of intermediate nodes to read for exact count

	// Read the org and end leaf nodes
	orgLeaf := bt.readLeaf(orgOff)
	orgI, _ := orgLeaf.search(org)

	endLeaf := bt.readLeaf(endOff)
	endI, _ := endLeaf.search(end)

	nodesBetween := endPos - orgPos - 1

	if nodesBetween == 0 {
		// adjacent leaf node - exact count
		keysInOrg := orgLeaf.noffs() - orgI
		keysInEnd := endI
		totalKeys := keysInOrg + keysInEnd
		return float64(totalKeys) / float64(bt.count)
	}

	// if org and end leaf nodes are close together
	// read all the nodes between for an exact count
	if nodesBetween <= maxToRead {
		totalKeys := orgLeaf.noffs() - orgI
		for i := orgPos + 1; i < endPos; i++ {
			leaf := bt.readLeaf(parent.offset(i))
			totalKeys += leaf.noffs()
		}
		totalKeys += endI
		return float64(totalKeys) / float64(bt.count)
	}

	// Many nodes between - use fanout estimate
	keysInOrg := orgLeaf.noffs() - orgI
	keysInEnd := endI
	estimatedBetween := float64(nodesBetween) * fanout
	totalKeys := float64(keysInOrg) + estimatedBetween + float64(keysInEnd)
	return totalKeys / float64(bt.count)
}
