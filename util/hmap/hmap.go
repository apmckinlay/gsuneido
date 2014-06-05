/*
Package hmap implements a hash map similar to Go's map
but that allows any key that has Hash and Equals methods.

Based on the Go map implementation:
http://golang.org/src/pkg/runtime/hashmap.c
Without incremental resizing.
*/
package hmap

import (
	"github.com/apmckinlay/gsuneido/util/bits"
	"github.com/apmckinlay/gsuneido/util/verify"
)

type Key interface {
	Hash() uint32
	Equals(interface{}) bool
}

type Hmap struct {
	size    int
	version int
	buckets []bucket // size must be power of two
}

const bucketsize = 8
const load = 6

// individual arrays to avoid padding on tophash
type bucket struct {
	keys     [bucketsize]Key
	vals     [bucketsize]interface{}
	tophash  [bucketsize]uint8 // zero means empty slot
	overflow *bucket
}

// NewHmap returns a new Hmap with at least the specified capacity
func NewHmap(cap int) *Hmap {
	cap = int(bits.Clp2(uint64(cap/load) + 1))
	return &Hmap{0, 0, make([]bucket, cap)}
}

func (hm *Hmap) hash(key Key) (b uint32, top uint8) {
	h := key.Hash()
	b = h & uint32(len(hm.buckets)-1)
	top = uint8(h >> 24) // based on 32 bit hash and 8 bit tophash
	if top == 0 {
		top = 1
	}
	return
}

// Size returns the current number of key/value entries
func (hm *Hmap) Size() int {
	return hm.size
}

// Get returns the value associated with the given key
// or nil if it isn't present
func (hm *Hmap) Get(key Key) interface{} {
	if hm.size == 0 {
		return nil
	}
	b, top := hm.hash(key)
	for buck := &hm.buckets[b]; buck != nil; buck = buck.overflow {
		for i := 0; i < bucketsize; i++ {
			if buck.tophash[i] == top && key.Equals(buck.keys[i]) {
				return buck.vals[i]
			}
		}
	}
	return nil
}

// Put inserts or updates the key to have the specified value
func (hm *Hmap) Put(key Key, val interface{}) {
	hm.ensureBuckets()
restart:
	b, top := hm.hash(key)
	inserti := -1
	var insertb *bucket = nil
	buck := &hm.buckets[b]
	for {
		for i := 0; i < bucketsize; i++ {
			if buck.tophash[i] != top {
				if buck.tophash[i] == 0 && inserti == -1 {
					inserti = i
					insertb = buck
				}
				continue
			}
			// found one where tophash matches
			if key.Equals(buck.keys[i]) {
				// update existing entry
				buck.vals[i] = val
				return
			}
		}
		if buck.overflow == nil {
			break
		}
		buck = buck.overflow
	}
	// not found, add it
	if hm.size > bucketsize && hm.size > load*len(hm.buckets) {
		hm.grow()
		goto restart
	}
	if inserti == -1 {
		// buckets are full, make a new one
		insertb = &bucket{}
		inserti = 0
		buck.overflow = insertb
	}
	insertb.tophash[inserti] = top
	insertb.keys[inserti] = key
	insertb.vals[inserti] = val
	hm.version++
	hm.size++
}

func (hm *Hmap) ensureBuckets() {
	if len(hm.buckets) == 0 {
		hm.buckets = make([]bucket, 1)
	}
}

// grow doubles the capacity
func (hm *Hmap) grow() {
	oldsize := hm.size
	oldbuckets := hm.buckets
	hm.buckets = make([]bucket, 2*len(hm.buckets))
	hm.size = 0
	for b := 0; b < len(oldbuckets); b++ {
		for buck := &oldbuckets[b]; buck != nil; buck = buck.overflow {
			for i := 0; i < bucketsize; i++ {
				if buck.tophash[i] != 0 {
					hm.Put(buck.keys[i], buck.vals[i])
					// NOTE: could use a specialized insert without dup checking
				}
			}
		}
	}
	verify.That(hm.size == oldsize)
}

// Del removes a key (and its value)
// returning true if the key was found or false if it wasn't
// NOTE: Does not shrink the Hmap
func (hm *Hmap) Del(key Key) (val interface{}) {
	if hm.size == 0 {
		return
	}
	b, top := hm.hash(key)
	for buck := &hm.buckets[b]; buck != nil; buck = buck.overflow {
		for i := 0; i < bucketsize; i++ {
			if buck.tophash[i] == top && key.Equals(buck.keys[i]) {
				buck.tophash[i] = 0
				val = buck.vals[i]
				buck.vals[i] = nil
				hm.version++
				hm.size--
				return
			}
		}
	}
	return
}

type Iter struct {
	hmap    *Hmap
	version int
	b       int
	buck    *bucket
	i       int
}

func (hm *Hmap) Iter() *Iter {
	return &Iter{hm, hm.version, 0, &hm.buckets[0], -1}
}

func (it *Iter) Next() (key Key, val interface{}) {
	if it.b >= len(it.hmap.buckets) {
		return nil, nil // end
	}
	if it.version != it.hmap.version {
		panic("hmap modified during iteration")
	}
	for {
		for it.i++; it.i < bucketsize && it.buck.tophash[it.i] == 0; it.i++ {
		}
		if it.i < bucketsize {
			return it.buck.keys[it.i], it.buck.vals[it.i]
		}
		it.i = -1
		it.buck = it.buck.overflow
		if it.buck != nil {
			continue
		}
		it.b++
		if it.b >= len(it.hmap.buckets) {
			return nil, nil // end
		}
		it.buck = &it.hmap.buckets[it.b]
	}
}
