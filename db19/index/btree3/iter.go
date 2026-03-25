// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"github.com/apmckinlay/gsuneido/db19/index/iface"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/util/assert"
)

// Iterator traverses a range of a btree.
type Iterator struct {
	bt               *btree
	rng              Range
	skipRng          Range               // suffix range used by skip-scan mode
	skipPrefixFields int                 // number of leading fields treated as prefix in skip-scan mode
	tree             [maxLevels]treeIter // tree[0] is root
	leaf             leafIter
	state            iterState
	noRange          bool // true if rng is iterator.All, bypasses checkRange
	curKeySet        bool
	curKey           string
	skipScan         bool   // true if SkipScan mode is active
	skipGroup        string // current prefix group in skip-scan traversal
}

const maxLevels = 8

type Range = iface.Range

type iterState byte

const (
	rewound iterState = iota
	within
	eof
)

func (bt *btree) Iterator() iface.Iter {
	return &Iterator{bt: bt, state: rewound, rng: iface.All, noRange: true,
		skipPrefixFields: 1}
}

// Key returns the current key or an empty string. It allocates.
func (it *Iterator) Key() string {
	if it.state != within {
		return ""
	}
	if !it.curKeySet {
		it.curKey = it.leaf.key()
		it.curKeySet = true
	}
	return it.curKey
}

// Offset returns the current offset or 0.
func (it *Iterator) Offset() uint64 {
	if it.state != within {
		return 0
	}
	return it.leaf.offset()
}

func (it *Iterator) Eof() bool {
	return it.state == eof
}

func (it *Iterator) Modified() bool {
	return false
}

// Cur returns the current key and offset. It allocates the key.
func (it *Iterator) Cur() (string, uint64) {
	return it.Key(), it.Offset()
}

// HasCur returns true if the iterator has a current item
func (it *Iterator) HasCur() bool {
	return it.state == within
}

// Rewind sets the iterator so Next goes to the first key in the range
// and Prev goes to the last key in the range
func (it *Iterator) Rewind() {
	it.state = rewound
	it.curKeySet = false
	it.skipGroup = ""
}

// Range sets the range and rewinds the iterator
func (it *Iterator) Range(rng Range) {
	it.rng = rng
	it.skipScan = false
	it.noRange = (rng == iface.All)
	it.Rewind()
}

// SkipScan enables btree3-only skip-scan mode.
// rng applies to suffix fields (excluding prefix fields).
// prefixFields defaults to 1.
func (it *Iterator) SkipScan(rng Range, prefixFields ...int) {
	n := 1
	if len(prefixFields) > 0 {
		n = prefixFields[0]
	}
	if n < 1 {
		panic("btree skip scan: prefixFields must be >= 1")
	}
	it.skipScan = true
	it.skipRng = rng
	it.skipPrefixFields = n
	it.rng = iface.All
	it.noRange = true
	it.skipGroup = ""
	it.Rewind()
}

//-------------------------------------------------------------------

// Next advances the iterator to the next key in the range or sets eof.
func (it *Iterator) Next() {
	it.curKeySet = false
	if it.skipScan {
		it.skipNext()
		return
	}
	switch it.state {
	case rewound:
		it.SeekAll(it.rng.Org)
		it.checkRange() // need to check both bounds after seek
	case within:
		it.next()
		it.checkRangeEnd() // only need to check end when moving forward
	case eof: // stick at eof
		return
	}
}

func (it *Iterator) next() {
	for {
		if it.leaf.next() {
			return
		} else if !it.nextLeaf() {
			it.state = eof
			return
		}
	}
}

func (it *Iterator) nextLeaf() bool {
	bt := it.bt
	i := bt.treeLevels - 1 // closest to leaf
	var nodeOff uint64
	// go up the tree until we can advance
	for ; i >= 0; i-- {
		if it.tree[i].next() {
			nodeOff = it.tree[i].offset()
			break
		} // else end of tree node, keep going up
	}
	if i < 0 {
		return false // end of root = eof
	}
	// then descend back down
	for i++; i < bt.treeLevels; i++ {
		it.tree[i] = bt.readTree(nodeOff).iter()
		assert.That(it.tree[i].next())
		nodeOff = it.tree[i].offset()
	}
	it.leaf = bt.readLeaf(nodeOff).iter()
	return true
}

func (it *Iterator) skipNext() {
	switch it.state {
	case rewound:
		it.seekAllRaw(ixkey.Min)
		if it.state != within {
			return
		}
	case within:
		it.next()
	case eof:
		return
	}
	it.skipAdvanceToMatch()
}

func (it *Iterator) skipAdvanceToMatch() {
	for it.state == within {
		key := it.leaf.key()
		prefix, suffix := ixkey.SplitPrefixSuffix(key, it.skipPrefixFields)
		if prefix != it.skipGroup {
			// New first-field group: jump directly to this group's suffix lower bound.
			// This avoids scanning the beginning of every group when Org is selective.
			it.skipGroup = prefix
			it.skipSeekGroupOrg(prefix)
			if it.state != within {
				return
			}
			continue
		}
		if suffix >= it.skipRng.End {
			it.skipSeekNextGroup(prefix)
			continue
		}
		return
	}
}

func (it *Iterator) skipSeekGroupOrg(prefix string) {
	target := prefix
	if it.skipRng.Org != ixkey.Min {
		target = prefix + ixkey.Sep + it.skipRng.Org
	}
	it.seekAllRaw(target)
	// If we can't land at/after target, this group has no keys in range.
	if it.state != within || it.Key() < target {
		it.state = eof
	}
}

func (it *Iterator) skipSeekNextGroup(prefix string) {
	target := prefix + ixkey.Sep + ixkey.Max

	// First try to find the target in the current leaf. This is the fast path when
	// the next first-field group is still in the same leaf node.
	it.leaf = it.leaf.nd.seek(target)
	if !it.leaf.eof() {
		it.state = within
		return
	}

	bt := it.bt
	if bt.treeLevels == 0 {
		it.state = eof
		return
	}
	// treeNode.seek always returns a valid position (i < noffs),
	// so we always find the ancestor at the deepest tree level.
	level := bt.treeLevels - 1
	it.tree[level] = it.tree[level].nd.seek(target)
	nodeOff := it.tree[level].offset()
	it.leaf = bt.readLeaf(nodeOff).seek(target)
	if !it.leaf.eof() {
		it.state = within
		return
	}
	it.next()
}

//-------------------------------------------------------------------

// Prev moves the iterator to the previous key in the range or sets eof.
func (it *Iterator) Prev() {
	it.curKeySet = false
	if it.skipScan {
		it.skipPrev()
		return
	}
	switch it.state {
	case rewound:
		it.SeekAll(it.rng.End)
		if it.Eof() {
			return // empty tree
		}
		if it.leaf.key() >= it.rng.End {
			it.prev()
		}
		it.checkRange() // need to check both bounds after seek
	case within:
		it.prev()
		it.checkRangeOrg() // only need to check org when moving backward
	case eof: // stick at eof
		return
	}
}

func (it *Iterator) prev() {
	for {
		if it.leaf.prev() {
			it.state = within
			return
		} else if !it.prevLeaf() {
			it.state = eof
			return
		}
	}
}

func (it *Iterator) prevLeaf() bool {
	bt := it.bt
	i := bt.treeLevels - 1 // closest to leaf
	var nodeOff uint64
	// go up the tree until we can go back
	for ; i >= 0; i-- {
		if it.tree[i].prev() {
			nodeOff = it.tree[i].offset()
			break
		} // else beginning of tree node, keep going up
	}
	if i < 0 {
		return false // beginning of root = eof
	}
	// then descend back down to rightmost
	for i++; i < bt.treeLevels; i++ {
		it.tree[i] = bt.readTree(nodeOff).iter()
		it.tree[i].i = it.tree[i].nd.noffs() - 1 // position at end
		nodeOff = it.tree[i].offset()
	}
	it.leaf = bt.readLeaf(nodeOff).iter()
	it.leaf.i = it.leaf.nd.nkeys() // position at end
	return true
}

func (it *Iterator) skipPrev() {
	// Start from physical end; skipRetreatToMatch applies suffix filtering.
	switch it.state {
	case rewound:
		it.seekAllRaw(ixkey.Max)
		if it.state != within {
			return
		}
	case within:
		it.prev()
	case eof:
		return
	}
	it.skipRetreatToMatch()
}

func (it *Iterator) skipRetreatToMatch() {
	// Reverse-direction mirror of skipAdvanceToMatch.
	for it.state == within {
		prefix, suffix := ixkey.SplitPrefixSuffix(it.leaf.key(), it.skipPrefixFields)
		if prefix != it.skipGroup {
			it.skipGroup = prefix
			it.skipSeekGroupEnd(prefix)
			if it.state != within {
				return
			}
			continue
		}
		if suffix >= it.skipRng.End {
			it.prev()
			continue
		}
		if suffix < it.skipRng.Org {
			it.skipSeekPrevGroup(prefix)
			continue
		}
		return
	}
}

func (it *Iterator) skipSeekGroupEnd(prefix string) {
	target := prefix + ixkey.Sep + it.skipRng.End
	it.seekAllRaw(target)
	// seekAllRaw positions >= target; back up until we're inside this group's range
	for it.state == within {
		p2, s2 := ixkey.SplitPrefixSuffix(it.leaf.key(), it.skipPrefixFields)
		if p2 > prefix || (p2 == prefix && s2 >= it.skipRng.End) {
			it.prev()
			it.curKeySet = false
			continue
		}
		// f2 < first: backed up past the group entirely (no keys in range).
		// Return with state=within in group f2; skipRetreatToMatch will
		// re-enter and process group f2 (skipFirst still == first, not f2).
		// f2 == first && s2 < End: in range, done.
		return
	}
}

func (it *Iterator) skipSeekPrevGroup(prefix string) {
	// Reuse already-loaded tree nodes: seek at the deepest tree level for 'first',
	// then descend to the leaf and back up one step.
	bt := it.bt
	if bt.treeLevels == 0 {
		// single-leaf tree: seek directly in the leaf
		it.leaf = bt.readLeaf(bt.root).seek(prefix)
	} else {
		// treeNode.seek always returns a valid position (i < noffs),
		// so we find the ancestor at level bt.treeLevels-1.
		level := bt.treeLevels - 1
		it.tree[level] = it.tree[level].nd.seek(prefix)
		nodeOff := it.tree[level].offset()
		it.leaf = bt.readLeaf(nodeOff).seek(prefix)
	}
	it.state = within
	if it.leaf.eof() || it.Key() >= prefix {
		it.prev()
	}
}

//-------------------------------------------------------------------

// Seek moves the iterator to the first position >= key.
// If the key is outside the current range, eof will be set.
func (it *Iterator) Seek(key string) {
	if !it.skipScan {
		it.SeekAll(key)
		it.checkRange()
		return
	}
	// In skip-scan mode, Seek should match regular Seek semantics on the
	// filtered view: find first visible key >= key, or stay on last visible key.
	visible := func(k string) (first string, ok bool) {
		p, s := ixkey.SplitPrefixSuffix(k, it.skipPrefixFields)
		return p, it.skipRng.Org <= s && s < it.skipRng.End
	}

	it.seekAllRaw(key)
	for it.state == within {
		k := it.leaf.key()
		if k < key {
			it.next()
			continue
		}
		if first, ok := visible(k); ok {
			it.skipGroup = first
			return
		}
		it.next()
	}

	// No visible key >= key; position to the last visible key.
	it.seekAllRaw(ixkey.Max)
	for it.state == within {
		if first, ok := visible(it.leaf.key()); ok {
			it.skipGroup = first
			return
		}
		it.prev()
	}
}

// SeekAll moves the iterator to the first position >= key.
// If the key is larger than the largest key,
// it will be positioned at the largest key.
// The state will be set to within
// unless the btree is empty in which case it will be set to eof.
// It does *not* apply the current range.
func (it *Iterator) SeekAll(key string) {
	if !it.skipScan {
		it.seekAllRaw(key)
		return
	}
	// In skip-scan mode, SeekAll finds the first key where suffix >= key,
	// without applying skipRng bounds. This mirrors normal SeekAll (no range).
	it.seekAllRaw(ixkey.Min)
	it.skipGroup = ""
	startedWithin := it.state == within
	it.skipSuffixSeekUnbounded(key)
	// Mirror seekAllRaw: if no suffix matched (but tree is non-empty), back up to
	// the last physical key rather than leaving the iterator at EOF.
	if startedWithin && it.state == eof {
		it.prev()
		it.state = within
	}
}

func (it *Iterator) seekAllRaw(key string) {
	it.curKeySet = false
	bt := it.bt
	off := bt.root
	rightEdge := true
	for i := range bt.treeLevels {
		it.tree[i] = bt.readTree(off).seek(key)
		off = it.tree[i].offset()
		rightEdge = rightEdge && it.tree[i].i >= it.tree[i].nd.nkeys()
	}
	leaf := bt.readLeaf(off)
	if leaf.nkeys() == 0 {
		assert.That(bt.treeLevels == 0) // only root can be empty
		it.state = eof
		return
	}
	it.leaf = leaf.seek(key)
	if it.leaf.eof() {
		if rightEdge {
			it.prev()
		} else {
			it.next()
		}
	}
	it.state = within
}

// skipSuffixSeekUnbounded advances to the first key with suffix >= minSuffix,
// across all first-field groups, without applying any upper bound.
func (it *Iterator) skipSuffixSeekUnbounded(minSuffix string) {
	for it.state == within {
		prefix, suffix := ixkey.SplitPrefixSuffix(it.leaf.key(), it.skipPrefixFields)
		if prefix != it.skipGroup {
			it.skipGroup = prefix
			target := prefix + ixkey.Sep + minSuffix
			it.seekAllRaw(target)
			it.skipGroup = prefix
			continue
		}
		if suffix < minSuffix {
			it.next()
			continue
		}
		return
	}
}

// checkRange changes state from within to eof
// if the current key is outside the range
func (it *Iterator) checkRange() {
	if it.noRange || it.state != within {
		return
	}
	prefix := it.leaf.prefix()
	suffix := it.leaf.suffix()
	if gte(prefix, suffix, it.rng.End) || !gte(prefix, suffix, it.rng.Org) {
		it.state = eof
	}
}

// checkRangeEnd changes state from within to eof
// if the current key is >= rng.End (used for Next)
func (it *Iterator) checkRangeEnd() {
	if it.noRange || it.state != within {
		return
	}
	prefix := it.leaf.prefix()
	suffix := it.leaf.suffix()
	if gte(prefix, suffix, it.rng.End) {
		it.state = eof
	}
}

// checkRangeOrg changes state from within to eof
// if the current key is < rng.Org (used for Prev)
func (it *Iterator) checkRangeOrg() {
	if it.noRange || it.state != within {
		return
	}
	prefix := it.leaf.prefix()
	suffix := it.leaf.suffix()
	if !gte(prefix, suffix, it.rng.Org) {
		it.state = eof
	}
}

// gte returns true if prefix+suffix >= target
// without concatenating prefix and suffix
func gte(prefix, suffix []byte, bound string) bool {
	plen := len(prefix)
	slen := len(suffix)
	tlen := len(bound)

	// Compare the prefix portion first
	cmpLen := min(tlen, plen)
	for i := 0; i < cmpLen; i++ {
		if prefix[i] < bound[i] {
			return false
		}
		if prefix[i] > bound[i] {
			return true
		}
	}

	// If bound is entirely within prefix length, check if we have more data
	if tlen <= plen {
		return plen+slen >= tlen
	}

	// Compare the suffix portion
	boundOffset := plen
	remainingBound := tlen - plen
	cmpLen = min(remainingBound, slen)
	for i := 0; i < cmpLen; i++ {
		if suffix[i] < bound[boundOffset+i] {
			return false
		}
		if suffix[i] > bound[boundOffset+i] {
			return true
		}
	}

	return plen+slen >= tlen
}
