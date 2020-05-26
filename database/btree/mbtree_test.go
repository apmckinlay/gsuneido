// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"math/rand"
	"sort"
	"testing"

	. "github.com/apmckinlay/gsuneido/util/hamcrest"
)

func TestMbtree(t *testing.T) {
	x := newMbtree()
	mCompare(t, x, mLeafSlots{})
	data := mLeafSlots{
		{"hello", 456},
		{"andrew", 123},
		{"zorro", 789},
	}
	for _, v := range data {
		x.Insert(v.key, v.rec)
	}
	mCompare(t, x, data)
}

func mCompare(t *testing.T, x *mbtree, data mLeafSlots) {
	sort.Sort(data)
	iter := x.Iterator()
	for _, v := range data {
		key, rec, ok := iter()
		Assert(t).That(ok, Equals(true))
		Assert(t).That(key, Equals(v.key))
		Assert(t).That(rec, Equals(v.rec))
	}
	_, _, ok := iter()
	Assert(t).That(ok, Equals(false))
}

type mLeafSlots []mLeafSlot

func (a mLeafSlots) Len() int      { return len(a) }
func (a mLeafSlots) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a mLeafSlots) Less(i, j int) bool {
	return a[i].key < a[j].key ||
		(a[i].key == a[j].key && a[i].rec < a[j].rec)
}

func TestMbtreeRandom(t *testing.T) {
	const n = mSize * 87
	data := make(mLeafSlots, n)
	x := newMbtree()
	for i := uint64(0); i < n; i++ {
		key := randomString() + randomString()
		x.Insert(key, i)
		data[i] = mLeafSlot{key, i}
	}
	mCompare(t, x, data)
}

const chars = "abcdefghijklmnopqrstuvwxyz"
const maxchars = 8

func randomString() string {
	n := 1 + rand.Int63()%maxchars
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[rand.Int63()%int64(len(chars))]
	}
	return string(b)
}

func TestMbtreeOrdered(t *testing.T) {
	const n = mSize * 87
	data := make(mLeafSlots, n)
	x := newMbtree()
	for i := uint64(0); i < n; i++ {
		key := randomString() + randomString()
		data[i] = mLeafSlot{key, i}
	}
	sort.Sort(data)
	for _, v := range data {
		x.Insert(v.key, v.rec)
	}
	mCompare(t, x, data)
}
