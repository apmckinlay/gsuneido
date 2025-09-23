// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"fmt"
	"testing"

	"github.com/apmckinlay/gsuneido/util/assert"
)

func TestTreeNode_builder(t *testing.T) {
	assert := assert.T(t).This

	// Test single key
	builder := &treeBuilder{}
	builder.add(1000, "hello")
	nd := builder.finish(2000)
	assert(nd.nkeys()).Is(1)
	assert(nd.size()).Is(len(nd))
	assert(nd.offset(0)).Is(1000)
	assert(string(nd.key(0))).Is("hello")
	assert(nd.offset(1)).Is(2000)

	// Test multiple keys with offsets
	builder = &treeBuilder{}
	keys := []string{"apple", "banana", "cherry"}
	offsets := []uint64{100, 200, 300}

	for i, key := range keys {
		builder.add(offsets[i], key)
	}

	nd = builder.finish(400)
	assert(nd.nkeys()).Is(3)
	assert(nd.size()).Is(len(nd))

	// Verify keys are stored correctly
	for i, expected := range keys {
		assert(string(nd.key(i))).Is(expected)
	}

	// Test maximum keys (255 limit)
	builder = &treeBuilder{}
	for i := range 255 {
		builder.add(uint64(i*1000), fmt.Sprintf("key%03d", i))
	}

	nd = builder.finish(1)
	assert(nd.nkeys()).Is(255)
	assert(nd.size()).Is(len(nd))

	// Verify first and last keys
	assert(string(nd.key(0))).Is("key000")
	assert(string(nd.key(254))).Is("key254")

	// Test panic on too many keys
	builder = &treeBuilder{}
	for i := range 256 {
		builder.add(123, fmt.Sprintf("key%03d", i))
	}
	assert(func() { builder.finish(123) }).Panics("too many keys")
}

func TestTreeNode_search(t *testing.T) {
	assert := assert.T(t).This

	// Test single key
	nd := makeTree(123, "hello", 456)
	assert(nd.String()).Is("tree{123 <hello> 456}")
	assert(nd.search("a")).Is(123)
	assert(nd.search("hello")).Is(456)
	assert(nd.search("z")).Is(456)

	// Test multiple keys
	data := []any{11, "apple", 22, "banana", 33, "cherry", 44, "date", 55}
	nd = makeTree(data...)
	assert(nd.nkeys()).Is(4)
	for i := 1; i < len(data)-1; i += 2 {
		assert(string(nd.key(i / 2))).Is(data[i])
		assert(nd.search(data[i].(string))).Is(data[i+1])
	}
	assert(nd.search("zebra")).Is(55)
	assert(nd.search("aaa")).Is(11)
	assert(nd.search("ban")).Is(22)
}

// makeTree takes offsets separated by keys
//
// e.g. makeTree(100, "apple", 200, "banana", 300, "cherry", 400)
func makeTree(args ...any) treeNode {
	if len(args) == 0 {
		return treeNode{}
	}
	var b treeBuilder
	for i := 0; i < len(args)-1; i += 2 {
		b.add(uint64(args[i].(int)), args[i+1].(string))
	}
	return b.finish(uint64(args[len(args)-1].(int)))
}

func TestTreeNode_seek(t *testing.T) {
	assert := assert.T(t).This

	// Test single key tree node
	nd := makeTree(100, "hello", 200)

	// Key smaller than existing key - should return iterator with i=-1
	it := nd.seek("apple")
	assert(it.i).Is(0)

	// Exact match - should position at key index
	it = nd.seek("hello")
	assert(it.i).Is(1)

	// Key larger than existing key - should position at last key
	it = nd.seek("zebra")
	assert(it.i).Is(1)

	// Test multiple keys
	nd = makeTree(11, "apple", 22, "banana", 33, "cherry", 44, "date", 55)

	// Test exact matches at different positions
	it = nd.seek("apple")
	assert(it.i).Is(1)

	it = nd.seek("banana")
	assert(it.i).Is(2)

	it = nd.seek("cherry")
	assert(it.i).Is(3)

	it = nd.seek("date")
	assert(it.i).Is(4)

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

	// Test iterator navigation from seek position
	it = nd.seek("banana")
	assert(it.i).Is(2)
	assert(it.next()).Is(true)
	assert(it.i).Is(3)

	assert(it.prev()).Is(true)
	assert(it.i).Is(2)

	// Test seek with prefix-like keys
	nd2 := makeTree(100, "test", 200, "testing", 300, "tests", 400)

	it = nd2.seek("test")
	assert(it.i).Is(1)

	it = nd2.seek("testing")
	assert(it.i).Is(2)

	it = nd2.seek("tests")
	assert(it.i).Is(3)

	// Test with key between "test" and "testing"
	it = nd2.seek("testg")
	assert(it.i).Is(1)

	// Test with key between "testing" and "tests"
	it = nd2.seek("testj")
	assert(it.i).Is(2)
}

func TestTreeNode_insert(t *testing.T) {
	assert := assert.T(t).This

	// Test inserting into empty node
	var nd treeNode
	nd = nd.insert(0, "key1", 100)
	assert(nd.String()).Is("tree{100 <key1> 0}")

	// Test inserting at beginning
	nd = nd.insert(0, "aaa", 50)
	assert(nd.String()).Is("tree{50 <aaa> 100 <key1> 0}")

	// Test inserting at end
	nd = nd.insert(2, "zzz", 200)
	assert(nd.String()).Is("tree{50 <aaa> 100 <key1> 200 <zzz> 0}")

	// Test inserting in middle
	nd = nd.insert(1, "middle", 75)
	assert(nd.String()).Is("tree{50 <aaa> 75 <middle> 100 <key1> 200 <zzz> 0}")

	// Test inserting into single key node
	nd = makeTree(123, "hello", 456)
	// Insert at beginning
	nd = nd.insert(0, "apple", 111)
	assert(nd.String()).Is("tree{111 <apple> 123 <hello> 456}")

	// Insert at end
	nd = makeTree(123, "hello", 456)
	nd = nd.insert(1, "zebra", 999)
	assert(nd.String()).Is("tree{123 <hello> 999 <zebra> 456}")
}

func TestTreeNode_update(t *testing.T) {
	assert := assert.T(t).This

	// Test updating single key node
	nd := makeTree(123, "hello", 456)
	nd = nd.update(0, 999)
	assert(nd.String()).Is("tree{999 <hello> 456}")

	// Test updating without modification to structure
	nd = makeTree(100, "apple", 200, "banana", 300, "cherry", 400)

	// Update first entry
	nd = nd.update(0, 150)
	assert(nd.String()).Is("tree{150 <apple> 200 <banana> 300 <cherry> 400}")

	// Update middle entry
	nd = makeTree(100, "apple", 200, "banana", 300, "cherry", 400)
	nd = nd.update(1, 250)
	assert(nd.String()).Is("tree{100 <apple> 250 <banana> 300 <cherry> 400}")

	// Update last entry
	nd = makeTree(100, "apple", 200, "banana", 300, "cherry", 400)
	nd = nd.update(2, 350)
	assert(nd.String()).Is("tree{100 <apple> 200 <banana> 350 <cherry> 400}")
}

func TestTreeNode_delete(t *testing.T) {
	assert := assert.T(t).This

	// Test deleting from a single-key node
	nd := makeTree(123, "hello", 456)
	result := nd.delete(0)
	assert(len(result)).Is(0) // empty node

	// Test deleting first
	nd = makeTree(100, "apple", 200, "banana", 300, "cherry", 400)
	nd = nd.delete(0) // delete "apple"
	assert(nd.String()).Is("tree{200 <banana> 300 <cherry> 400}")

	// Test deleting from middle
	nd = makeTree(100, "apple", 200, "banana", 300, "cherry", 400)
	nd = nd.delete(1) // delete "banana"
	assert(nd.String()).Is("tree{100 <apple> 300 <cherry> 400}")

	// Test deleting from end
	nd = makeTree(100, "apple", 200, "banana", 300, "cherry", 400)
	nd = nd.delete(2) // delete "cherry"
	assert(nd.String()).Is("tree{100 <apple> 200 <banana> 400}")

	// Test deleting down to single key
	nd = makeTree(100, "apple", 200, "banana", 300)
	nd = nd.delete(0) // delete "apple"
	assert(nd.String()).Is("tree{200 <banana> 300}")
}
