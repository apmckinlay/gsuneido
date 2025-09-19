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

// func TestTreeNode_insert(t *testing.T) {
// 	assert := assert.T(t).This

// 	nd := makeTree(123, "hello", 999)
// 	assert(nd.String()).Is("tree{123 hello 999}")

// 	// Test inserting at beginning
// 	nd = nd.insert("apple", 456)
// 	assert(nd.String()).Is("tree{456 apple 123 hello 999}")

// 	// Test inserting in middle
// 	nd = nd.insert("banana", 789)
// 	assert(nd.String()).Is("tree{456 apple 789 banana 123 hello 999}")

// 	// Test inserting at end
// 	nd = nd.insert("zebra", 888)
// 	assert(nd.String()).Is("tree{456 apple 789 banana 123 hello 888 zebra 999}")

// 	// Test inserting duplicate (should panic)
// 	assert(func() { nd.insert("banana", 000) }).Panics("duplicate key")

// 	// Test with makeTree for comparison
// 	nd2 := makeTree(456, "apple", 789, "banana", 123, "hello", 999, "zebra", 1000)
// 	assert(nd.String()).Is(nd2.String())

// 	// Test search functionality after inserts
// 	assert(nd.search("apple")).Is(1)
// 	assert(nd.search("banana")).Is(2)
// 	assert(nd.search("hello")).Is(3)
// 	assert(nd.search("zebra")).Is(4)
// 	assert(nd.search("aaa")).Is(0)
// 	assert(nd.search("cat")).Is(1)
// 	assert(nd.search("fox")).Is(2)
// 	assert(nd.search("zebraa")).Is(4)
// }
