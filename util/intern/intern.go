// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package intern handles sharing common strings
package intern

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/apmckinlay/gsuneido/util/hash"
)

// redundant to store string as key and value
// but no way to get key that matched lookup

// buf holds the actual data of the strings
var buffer []byte

// strings points into buffer
var strings = []string{""} // hash set can't handle zero

// hashes stores/caches the high bits of the string hashes
// to speed up resizing the hash set
var hashes = []uint16{0}

var set hashSet16

var lock sync.Mutex

// Index returns the index a shared copy of a string,
// adding it if not present.
func Index(s string) uint16 {
	lock.Lock()
	defer lock.Unlock()
	h := hash16(s)
	pi := set.slotFor(s, h)
	if *pi == 0 {
		t := addStringToBuffer(s)
		strings = append(strings, t)
		hashes = append(hashes, h)
		*pi = uint16(len(strings) - 1)
	}
	return *pi
}

func String(i uint16) string {
	lock.Lock()
	defer lock.Unlock()
	return strings[i]
}

func addStringToBuffer(s string) string {
	i := len(buffer)
	buffer = append(buffer, s...)
	p := buffer[i:]
	return *(*string)(unsafe.Pointer(&p))
}

// hashSet16 stores a set of int16 indexes into strings
type hashSet16 struct {
	tbl      []uint16
	size     int
	capShift uint32
	capMask  uint16
}

const bighop = 11 // should be prime

func (hs *hashSet16) slotFor(s string, h uint16) *uint16 {
	if hs.size >= len(hs.tbl)*5/8 { // roughly 60%
		hs.grow()
	}
	return hs.slotFor_(s, h)
}
func (hs *hashSet16) slotFor_(s string, h uint16) *uint16 {
	i := h >> hs.capShift
	hops := 0
	for hs.tbl[i] != 0 {
		if strings[hs.tbl[i]] == s {
			return &hs.tbl[i]
		}
		if hops < bighop {
			i++
		} else {
			i += bighop
		}
		i &= hs.capMask
		hops++
	}
	hs.size++
	return &hs.tbl[i]
}

const phi32 = 2654435769

func hash16(s string) uint16 {
	h := hash.HashString(s)
	return uint16((h * phi32) >> 16)
}

const initialSize = 2          //1024        // should be power of 2
const initialCapShift = 16 - 1 //16 - 10 // must match initialSize
const initialCapMask = 1

func (hs *hashSet16) grow() {
	if len(hs.tbl) == 0 {
		hs.tbl = make([]uint16, initialSize)
		hs.capShift = initialCapShift
		hs.capMask = initialCapMask
		return
	}
	hs.tbl = make([]uint16, len(hs.tbl)*2)
	hs.size = 0
	hs.capShift--
	hs.capMask = hs.capMask<<1 + 1
	for i, s := range strings {
		*hs.slotFor(s, hashes[i]) = uint16(i)
	}
}

func (hs *hashSet16) print() {
	fmt.Println("size", hs.size, "cap", len(hs.tbl), "lim", len(hs.tbl)*5/8)
	for _, x := range hs.tbl {
		fmt.Printf("%v ", x)
	}
	fmt.Println()
}
