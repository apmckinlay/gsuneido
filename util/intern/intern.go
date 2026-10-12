// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// package intern stores a set of strings.
// Each string is separately allocated in strp format
// (length prefixed) so it can be used directly as a Go string
// or as a *SuStrP, sharing the same bytes.
// NOTE: strings are never removed or garbage collected (other than Clear)
package intern

import (
	"strings"
	"sync"

	"github.com/apmckinlay/gsuneido/util/hash"
	"github.com/apmckinlay/gsuneido/util/shmap"
	"github.com/apmckinlay/gsuneido/util/strp"
)

var (
	htbl shmap.Map[*byte, struct{}, helper]
	lock       sync.Mutex
	bytesAlloc int
)

type helper struct{}

func (helper) Hash(p *byte) uint64 {
	return hash.String(strp.String(p))
}

func (helper) Equal(x, y *byte) bool {
	return strp.String(x) == strp.String(y)
}

// lookup supports probing the table with a Go string
type lookup struct{}

func (lookup) Hash(s string) uint64 {
	return hash.String(s)
}

func (lookup) Equal(s string, p *byte) bool {
	return s == strp.String(p)
}

// String returns the interned copy of s.
// Strings longer than strp.MaxLen are cloned but not stored.
func String(s string) string {
	if p := StrP(s); p != nil {
		return strp.String(p)
	}
	return strings.Clone(s)
}

// StrP returns a length-prefixed pointer to interned storage for s.
// It returns nil if len(s) > strp.MaxLen.
// The returned pointer remains valid even after Clear.
func StrP(s string) *byte {
	if len(s) > strp.MaxLen {
		return nil
	}
	lock.Lock()
	defer lock.Unlock()
	if p, ok := htbl.GetKeyAs(s, lookup{}); ok {
		return p
	}
	p := strp.From(s)
	htbl.Put(p, struct{}{})
	bytesAlloc += len(s) + 1
	return p
}

// Count returns the number of unique strings stored
func Count() int {
	lock.Lock()
	defer lock.Unlock()
	return htbl.Size()
}

// Bytes returns the total allocated size of the strings
func Bytes() int {
	lock.Lock()
	defer lock.Unlock()
	return bytesAlloc
}

// Clear discards all the strings and starts fresh
func Clear() {
	lock.Lock()
	defer lock.Unlock()
	htbl.Clear()
	bytesAlloc = 0
}
