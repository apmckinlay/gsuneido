// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package btree

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/btree/inter"
	"github.com/apmckinlay/gsuneido/db19/ixspec"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestFbMerge(*testing.T) {
	nMerges := 2000
	insertsPerMerge := 1000
	if testing.Short() {
		nMerges = 200
		insertsPerMerge = 200
	}
	d := newdat()
	GetLeafKey = func(_ *stor.Stor, _ *ixspec.T, i uint64) string {
		return d.o2k[i]
	}
	defer func(mns int) { MaxNodeSize = mns }(MaxNodeSize)
	MaxNodeSize = 64
	store := stor.HeapStor(8192)
	store.Alloc(1) // avoid offset 0
	fb := CreateFbtree(store, nil)

	for i := 0; i < nMerges; i++ {
		_ = T && trace("---")
		x := &inter.T{}
		for j := 0; j < insertsPerMerge; j++ {
			x.Insert(d.next(""))
		}
		fb = fb.MergeAndSave(x.Iter(false))
	}
	fb.check(nil)
	d.check(fb)
}

func (st *state) print() {
	fmt.Println("state:", st.fb.treeLevels)
	for _, m := range st.path {
		fmt.Println("   ", &m)
		fmt.Println("       ", m.node.knowns())
	}
}

func (m *merge) String() string {
	limit := m.limit
	if limit == "" {
		limit = `""`
	}
	mod := ""
	if m.modified {
		mod = " modified"
	}
	return fmt.Sprint("off ", m.off, " fi ", m.fi, " limit ", limit, mod)
}

type dat struct {
	keys []string
	o2k  map[uint64]string
	rand func() string
}

func newdat() *dat {
	return &dat{
		o2k:  map[uint64]string{},
		rand: str.UniqueRandom(4, 8),
	}
}

func (d *dat) next(prefix string) (string, uint64) {
	key := prefix + d.rand()
	off := uint64(len(d.keys))
	d.o2k[off] = key
	d.keys = append(d.keys, key)
	return key, off
}

func (d *dat) check(fb *fbtree) {
	sort.Strings(d.keys)
	i := 0
	iter := fb.Iter(true)
	for k, o, ok := iter(); ok; k, o, ok = iter() {
		assert.Msg("expect prefix of " + d.keys[i] + " got " + k).
			That(strings.HasPrefix(d.keys[i], k))
		assert.This(d.o2k[o]).Is(d.keys[i])
		i++
	}
	assert.This(i).Is(len(d.keys))
}
