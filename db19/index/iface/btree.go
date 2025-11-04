// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package iface

import (
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/stor"
)

type Btree interface {
	Cksum() uint32
	SetIxspec(is *ixkey.Spec)
	TreeLevels() int
	Lookup(key string) uint64
	QuickCheck()
	Check(fn any) (count, size, nnodes int)
	Write(*stor.Writer)
	RangeFrac(org, end string, nrecs int) float64
	Iterator() Iter
	MergeAndSave(iter IterFn) Btree
	SetSplit(ndsize int)
}

type BtreeBuilder interface {
	Add(key string, off uint64) bool
	Finish() Btree
}

type IterFn = func() (key string, off uint64, ok bool)
