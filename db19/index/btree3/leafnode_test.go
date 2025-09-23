// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"fmt"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestLeafNode_builder(t *testing.T) {
	assert := assert.This

	// Test single key
	b := &leafBuilder{}
	b.add("hello", 255)
	nd := b.finish()
	assert(nd.nkeys()).Is(1)
	assert(nd.offset(0)).Is(255)
	assert(nd.key(0)).Is("hello")
	assert(string(nd.prefix())).Is("")
	assert(nd.String()).Is("leaf{hello 255}")

	// Test multiple keys with offsets
	keys := []string{"apple", "banana", "cherry"}
	offsets := []uint64{100, 200, 300}
	b.reset()
	for i, key := range keys {
		b.add(key, offsets[i])
	}
	nd = b.finish()
	assert(nd.nkeys()).Is(3)
	assert(string(nd.prefix())).Is("")

	// Verify keys are stored correctly
	for i, expected := range keys {
		assert(nd.key(i)).Is(expected)
	}

	// Test maximum keys (255 limit)
	b.reset()
	for i := range 255 {
		b.add(fmt.Sprintf("key%03d", i), uint64(1000+i))
	}
	nd = b.finish()
	assert(nd.nkeys()).Is(255)
	assert(string(nd.prefix())).Is("key")

	// Verify first and last keys
	assert(nd.key(0)).Is("key000")
	assert(nd.key(254)).Is("key254")

	// Test panic on too many keys
	b.reset()
	for i := range 256 {
		b.add(fmt.Sprintf("key%03d", i), 123)
	}
	assert(func() { b.finish() }).Panics("too many keys")
}

func TestLeafPrefix(t *testing.T) {
	nd := makeLeaf("hello", 12, "helloa", 34, "hellfire", 56, "help", 78)
	assert.This(string(nd.prefix())).Is("hel")
	assert.This(nd.String()).Is("leaf{|hel| lo 12 loa 34 lfire 56 p 78}")
	assert.This(nd.key(0)).Is("hello")
	assert.This(nd.key(1)).Is("helloa")
	assert.This(nd.key(2)).Is("hellfire")
	assert.This(nd.key(3)).Is("help")
}

func TestLeafSearchEdgeCases(t *testing.T) {
	assert := assert.This

	var nd leafNode
	test := func(key string, expectedPos int, expectedFound bool) {
		t.Helper()
		pos, found := nd.search(key)
		assert(found).Is(expectedFound)
		assert(pos).Is(expectedPos)
	}

	// Test search with prefix compression
	nd = makeLeaf("prefix001", 1001, "prefix002", 1002,
		"prefix005", 1005, "prefix010", 1010)

	// Test exact matches
	test("prefix001", 0, true)
	test("prefix002", 1, true)
	test("prefix005", 2, true)
	test("prefix010", 3, true)

	// Test key smaller than prefix
	test("pre", 0, false) // would be at position 0

	// Test key larger than prefix but not matching any entry
	test("prefix003", 2, false) // would be at position 2
	test("prefix006", 3, false) // would be at position 3

	// Test key larger than all entries
	test("prefix999", 4, false) // would be at end

	// Test key with different prefix
	test("aaa", 0, false) // would be at start (aaa < prefix)
	test("zzz", 4, false) // would be at end (zzz > prefix)

	// Test no prefix compression
	nd = makeLeaf("apple", 100, "banana", 200, "cherry", 300)

	test("aaa", 0, false) // would be at start
	test("apple", 0, true)
	test("banana", 1, true)
	test("bear", 2, false)
	test("cherry", 2, true)
	test("zebra", 3, false) // would be at end

	// Test single key node
	nd = makeLeaf("single", 999)

	test("aaa", 0, false) // would be at start
	test("single", 0, true)
	test("zzz", 1, false) // would be at end
}

// makeLeaf takes offsets separated by keys
//
// e.g. makeLeaf("apple", 200, "banana", 300, "cherry", 400)
func makeLeaf(args ...any) leafNode {
	if len(args) == 0 {
		return leafNode{}
	}
	var b leafBuilder
	for i := 0; i < len(args)-1; i += 2 {
		b.add(args[i].(string), uint64(args[i+1].(int)))
	}
	return b.finish()
}

func BenchmarkLeafSearch(b *testing.B) {
	// Test with prefix compression (6-char prefix + 4-char suffix)
	b.Run("Prefix", func(b *testing.B) {
		builder := &leafBuilder{}
		// Create 75 keys with 6-char prefix + 4-char suffix
		for i := range 75 {
			key := fmt.Sprintf("prefix%04d", i) // 10 chars total
			builder.add(key, uint64(1000+i))
		}
		nd := builder.finish()

		// Search for keys at different positions
		searchKeys := []string{
			// "aaa",
			"prefix0000", // first key
			"prefix0037", // middle key
			"prefix0074", // last key
			"prefix0050", // another middle key
			// "zzz",
		}

		for b.Loop() {
			for _, key := range searchKeys {
				_, _ = nd.search(key)
			}
		}
	})

	// Test with no prefix compression (diverse 10-char keys)
	b.Run("NoPrefix", func(b *testing.B) {
		builder := &leafBuilder{}
		// Create 75 diverse keys with no common prefix
		keys := make([]string, 75)
		for i := range 75 {
			// Generate diverse 10-character keys
			key := fmt.Sprintf("%c%c%c%c%c%05d",
				'a'+byte(i%26),
				'b'+byte((i*7)%26),
				'c'+byte((i*13)%26),
				'd'+byte((i*19)%26),
				'e'+byte((i*23)%26),
				i)
			keys[i] = key
			builder.add(key, uint64(2000+i))
		}
		nd := builder.finish()

		// Search for keys at different positions
		searchKeys := []string{
			keys[0],  // first key
			keys[37], // middle key
			keys[74], // last key
			keys[50], // another middle key
		}

		for b.Loop() {
			for _, key := range searchKeys {
				_, _ = nd.search(key)
			}
		}
	})
}

func TestLeafNode_seek(t *testing.T) {
	assert := assert.T(t).This

	// Test single key node
	nd := makeLeaf("hello", 200)

	// Key smaller than existing key - should return iterator with i=-1
	it := nd.seek("apple")
	assert(it.i).Is(0)

	// Exact match - should position at key index
	it = nd.seek("hello")
	assert(it.i).Is(0)

	// Key larger than existing key
	it = nd.seek("zebra")
	assert(it.i).Is(1)

	// Test multiple keys without prefix compression
	nd = makeLeaf("apple", 100, "banana", 200, "cherry", 300, "date", 400)

	// Test exact matches at different positions
	it = nd.seek("apple")
	assert(it.i).Is(0)

	it = nd.seek("banana")
	assert(it.i).Is(1)

	it = nd.seek("cherry")
	assert(it.i).Is(2)

	it = nd.seek("date")
	assert(it.i).Is(3)

	// Test keys between existing keys - should find last key <= search key
	it = nd.seek("avocado") // between "apple" and "banana"
	assert(it.i).Is(1)

	it = nd.seek("blueberry") // between "banana" and "cherry"
	assert(it.i).Is(2)

	it = nd.seek("coconut") // between "cherry" and "date"
	assert(it.i).Is(3)

	// Test key smaller than all keys
	it = nd.seek("aaa")
	assert(it.i).Is(0)

	// Test key larger than all keys
	it = nd.seek("zebra")
	assert(it.i).Is(4)

	// Test with prefix compression
	nd = makeLeaf("prefix001", 1001, "prefix002", 1002, "prefix005", 1005, "prefix010", 1010)

	// Test exact matches with prefix
	it = nd.seek("prefix001")
	assert(it.i).Is(0)

	it = nd.seek("prefix002")
	assert(it.i).Is(1)

	it = nd.seek("prefix005")
	assert(it.i).Is(2)

	it = nd.seek("prefix010")
	assert(it.i).Is(3)

	// Test key smaller than prefix
	it = nd.seek("pre")
	assert(it.i).Is(0)

	it = nd.seek("aaa")
	assert(it.i).Is(0)

	// Test key between prefix entries
	it = nd.seek("prefix003") // between "prefix002" and "prefix005"
	assert(it.i).Is(2)

	it = nd.seek("prefix006") // between "prefix005" and "prefix010"
	assert(it.i).Is(3)

	// Test key larger than prefix but not matching entries
	it = nd.seek("prefix999")
	assert(it.i).Is(4)

	// Test key with different prefix but lexicographically >= prefix
	it = nd.seek("zzz")
	assert(it.i).Is(4) // all entries are < "zzz"

	// Test iterator navigation from seek position
	nd = makeLeaf("apple", 100, "banana", 200, "cherry", 300)
	it = nd.seek("banana")
	assert(it.i).Is(1)

	// Test next/prev from seek position
	assert(it.next()).Is(true)
	assert(it.i).Is(2)

	assert(it.prev()).Is(true)
	assert(it.i).Is(1)

	// Test seek to position -1 and then next
	it = nd.seek("aaa")
	assert(it.i).Is(0)
	assert(it.next()).Is(true)
	assert(it.i).Is(1)
}

func TestLeafNodeDelete(t *testing.T) {
	assert := assert.T(t).This

	// Test deleting from a single-key node
	nd := makeLeaf("hello", 123)
	assert(nd.nkeys()).Is(1)

	result := nd.delete(0)
	assert(len(result)).Is(0) // empty node

	// WITHOUT PREFIX
	// Test deleting first
	nd = makeLeaf("apple", 100, "banana", 200, "cherry", 300)
	nd = nd.delete(0) // delete "apple"
	assert(nd.String()).Is("leaf{banana 200 cherry 300}")

	// Test deleting from middle
	nd = makeLeaf("apple", 100, "banana", 200, "cherry", 300)
	nd = nd.delete(1) // delete "banana"
	assert(nd.String()).Is("leaf{apple 100 cherry 300}")

	// Test deleting from end
	nd = makeLeaf("apple", 100, "banana", 200, "cherry", 300)
	nd = nd.delete(2) // delete "cherry"
	assert(nd.String()).Is("leaf{apple 100 banana 200}")

	// WITH PREFIX
	// Test deleting first
	nd = makeLeaf("prefix1", 1, "prefix2", 2, "prefix5", 5, "prefix9", 9)
	nd = nd.delete(0) // delete "prefix1"
	assert(nd.String()).Is("leaf{|prefix| 2 2 5 5 9 9}")

	// Test deleting from middle
	nd = makeLeaf("prefix1", 1, "prefix2", 2, "prefix5", 5, "prefix9", 9)
	nd = nd.delete(1) // delete "prefix2"
	assert(nd.String()).Is("leaf{|prefix| 1 1 5 5 9 9}")

	// Test deleting from end
	nd = makeLeaf("prefix1", 1, "prefix2", 2, "prefix5", 5, "prefix9", 9)
	nd = nd.delete(3) // delete "prefix9"
	assert(nd.String()).Is("leaf{|prefix| 1 1 2 2 5 5}")

	// Test deleting second last - prefix should be removed
	nd = makeLeaf("prefix1", 1, "prefix2", 2)
	result = nd.delete(0) // delete "prefix1", leaving only "prefix2"
	assert(result.String()).Is("leaf{prefix2 2}")
}

func TestLeafNodeUpdate(t *testing.T) {
	assert := assert.T(t).This

	// Test updating single key node
	nd := makeLeaf("hello", 123)
	nd = nd.update(0, 456)
	assert(nd.String()).Is("leaf{hello 456}")

	// Test updating without prefix compression
	nd = makeLeaf("apple", 100, "banana", 200, "cherry", 300)

	// Update first entry
	nd = nd.update(0, 999)
	assert(nd.String()).Is("leaf{apple 999 banana 200 cherry 300}")

	// Update middle entry
	nd = makeLeaf("apple", 100, "banana", 200, "cherry", 300)
	nd = nd.update(1, 777)
	assert(nd.String()).Is("leaf{apple 100 banana 777 cherry 300}")

	// Update last entry
	nd = makeLeaf("apple", 100, "banana", 200, "cherry", 300)
	nd = nd.update(2, 555)
	assert(nd.String()).Is("leaf{apple 100 banana 200 cherry 555}")

	// Test updating with prefix compression
	// Update first entry with prefix
	nd = makeLeaf("prefix1", 1001, "prefix2", 1002, "prefix5", 1005, "prefix9", 1009)
	nd = nd.update(0, 2001)
	assert(nd.String()).Is("leaf{|prefix| 1 2001 2 1002 5 1005 9 1009}")

	// Update middle entry with prefix
	nd = makeLeaf("prefix1", 1001, "prefix2", 1002, "prefix5", 1005, "prefix9", 1009)
	nd = nd.update(2, 3005)
	assert(nd.String()).Is("leaf{|prefix| 1 1001 2 1002 5 3005 9 1009}")

	// Update last entry with prefix
	nd = makeLeaf("prefix1", 1001, "prefix2", 1002, "prefix5", 1005, "prefix9", 1009)
	nd = nd.update(3, 4009)
	assert(nd.String()).Is("leaf{|prefix| 1 1001 2 1002 5 1005 9 4009}")
}

func TestLeafNodeInsert(t *testing.T) {
	assert := assert.T(t).This

	// Test inserting into empty node
	nd := leafNode{}
	nd = nd.insert(0, "hello", 123)
	assert(nd.String()).Is("leaf{hello 123}")

	// Test inserting into single key node
	nd = makeLeaf("hello", 123)

	// Insert at beginning
	nd = nd.insert(0, "apple", 100)
	assert(nd.String()).Is("leaf{apple 100 hello 123}")

	// Insert at end
	nd = makeLeaf("hello", 123)
	nd = nd.insert(1, "world", 456)
	assert(nd.String()).Is("leaf{hello 123 world 456}")

	// WITHOUT PREFIX
	// Test inserting at beginning
	nd = makeLeaf("banana", 200, "cherry", 300, "date", 400)
	nd = nd.insert(0, "apple", 100)
	assert(nd.String()).Is("leaf{apple 100 banana 200 cherry 300 date 400}")

	// Test inserting in middle
	nd = makeLeaf("apple", 100, "cherry", 300, "date", 400)
	nd = nd.insert(1, "banana", 200)
	assert(nd.String()).Is("leaf{apple 100 banana 200 cherry 300 date 400}")

	// Test inserting at end
	nd = makeLeaf("apple", 100, "banana", 200, "cherry", 300)
	nd = nd.insert(3, "date", 400)
	assert(nd.String()).Is("leaf{apple 100 banana 200 cherry 300 date 400}")

	// WITH PREFIX
	// Test inserting at beginning with same prefix
	nd = makeLeaf("prefix2", 2, "prefix5", 5, "prefix9", 9)
	nd = nd.insert(0, "prefix1", 1)
	assert(nd.String()).Is("leaf{|prefix| 1 1 2 2 5 5 9 9}")

	// Test inserting in middle with same prefix
	nd = makeLeaf("prefix1", 1, "prefix5", 5, "prefix9", 9)
	nd = nd.insert(1, "prefix2", 2)
	assert(nd.String()).Is("leaf{|prefix| 1 1 2 2 5 5 9 9}")

	// Test inserting at end with same prefix
	nd = makeLeaf("prefix1", 1, "prefix2", 2, "prefix5", 5)
	nd = nd.insert(3, "prefix9", 9)
	assert(nd.String()).Is("leaf{|prefix| 1 1 2 2 5 5 9 9}")

	// PREFIX CHANGE CASES (should use leafBuilder)
	// Test inserting at beginning that changes prefix
	nd = makeLeaf("prefix1", 1, "prefix2", 2, "prefix5", 5)
	nd = nd.insert(0, "aaa", 0)
	assert(nd.String()).Is("leaf{aaa 0 prefix1 1 prefix2 2 prefix5 5}")

	// Test inserting at end that changes prefix
	nd = makeLeaf("prefix1", 1, "prefix2", 2, "prefix5", 5)
	nd = nd.insert(3, "zzz", 999)
	assert(nd.String()).Is("leaf{prefix1 1 prefix2 2 prefix5 5 zzz 999}")

	// Test inserting key that shortens prefix
	nd = makeLeaf("prefixlong1", 1, "prefixlong2", 2)
	nd = nd.insert(0, "prebar", 3)
	assert(nd.String()).Is("leaf{|pre| bar 3 fixlong1 1 fixlong2 2}")

	// Test maximum keys limit
	nd = leafNode{}
	for i := 0; i < 255; i++ {
		nd = nd.insert(i, fmt.Sprintf("key%03d", i), uint64(i))
	}
	assert(nd.nkeys()).Is(255)

	// Test panic on too many keys
	assert(func() { nd.insert(255, "overflow", 999) }).Panics("too many keys")
}

// func (nd leafNode) print() {
// 	if len(nd) == 0 {
// 		fmt.Println("Empty leafNode")
// 		return
// 	}
// 	fmt.Println("leafNode nkeys:", nd.nkeys(), "prelen:", int(nd[1]))
// 	for i := 0; i < nd.nkeys(); i++ {
// 		base := 2 + i*7
// 		fldoff := uint16(nd[base])<<8 | uint16(nd[base+1])
// 		fmt.Println("fldoff:", fldoff, "dboff:", nd.offset(i))
// 	}
// 	fmt.Println("Prefix:", string(nd.prefix()))
// 	for i := range nd.nkeys() {
// 		fmt.Println("Key:", nd.key(i))
// 	}
// }
