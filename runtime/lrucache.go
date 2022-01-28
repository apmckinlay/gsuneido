// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

type LruCache struct {
	size int
	hm    Hmap
	getfn func(Value) Value
	// entry is the head/tail of the lru link list.
	// The values are not used.
	entry
}

type entry struct {
	Value
	key Value
	next *entry
	prev *entry
}

var _ Value = (*entry)(nil)

func NewLruCache(size int, getfn func(Value) Value) *LruCache {
	lru := &LruCache{size: size, getfn: getfn}
	lru.next = &lru.entry
	lru.prev = &lru.entry
	return lru
}

func (lru *LruCache) Get(key Value) Value {
	e := lru.hm.Get(key).(*entry)
	if e != nil { // in cache
		lru.moveToFront(e)
		return e.Value
	}
	// not in cache
	if lru.hm.Size() > lru.size {
		// discard oldest (lru.prev)
		lru.hm.Del(lru.prev.key)
		lru.unlink(lru.prev)
	}
	// get/calc it and add it
	val := lru.getfn(key)
	e = &entry{Value: val, key: key}
	lru.hm.Put(key, e)
	lru.insert(e)
	return val
}

func (lru *LruCache) moveToFront(e *entry) {
	lru.unlink(e)
	lru.insert(e)
}

func (lru *LruCache) unlink(e *entry) {
	e.prev.next = e.next
	e.next.prev = e.prev
}

func (lru *LruCache) insert(e *entry) {
	e.prev = &lru.entry
	e.next = lru.next
	e.next.prev = e
	lru.next = e
}
