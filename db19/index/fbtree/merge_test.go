// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package fbtree

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"testing"

	"github.com/apmckinlay/gsuneido/db19/index/ixbuf"
	"github.com/apmckinlay/gsuneido/db19/index/ixspec"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

func TestMerge(*testing.T) {
	nMerges := 2000
	opsPerMerge := 1000
	if testing.Short() {
		nMerges = 200
		opsPerMerge = 200
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
		_ = t && trace("---")
		x := &ixbuf.T{}
		for j := 0; j < opsPerMerge; j++ {
			k := rand.Intn(4)
			switch {
			case k == 0 || k == 1 || d.Len() == 0:
				x.Insert(d.gen())
			case k == 2:
				_, key, _ := d.rand()
				off := d.nextOff()
				x.Update(key, off)
				d.update(key, off)
			case k == 3:
				i, key, off := d.rand()
				x.Delete(key, off)
				d.delete(i)
			}
		}
		fb = fb.MergeAndSave(x.Iter(false))
	}
	fb.Check(nil)
	d.check(fb)
}

//-------------------------------------------------------------------

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
	k2o  map[string]uint64
	off  uint64
}

func newdat() *dat {
	return &dat{
		o2k: map[uint64]string{},
		k2o: map[string]uint64{},
	}
}

func (d *dat) Len() int {
	return len(d.keys)
}

func (d *dat) gen() (string, uint64) {
	for i := 0; i < 10; i++ {
		key := str.Random(4, 8)
		if _, ok := d.k2o[key]; !ok {
			off := d.nextOff()
			d.k2o[key] = off
			d.o2k[off] = key
			d.keys = append(d.keys, key)
			return key, off
		}
	}
	panic("too many duplicates")
}

func (d *dat) nextOff() uint64 {
	d.off++
	return d.off
}

func (d *dat) rand() (int, string, uint64) {
	i := rand.Intn(len(d.keys))
	key := d.keys[i]
	off := d.k2o[key]
	return i, key, off
}

func (d *dat) delete(i int) {
	last := len(d.keys) - 1
	d.keys[i] = d.keys[last]
	d.keys = d.keys[:last]
}

func (d *dat) update(key string, off uint64) {
	// oldoff := d.k2o[key]
	d.k2o[key] = off
	// delete(d.o2k, oldoff)
	d.o2k[off] = key
}

func (d *dat) check(fb *fbtree) {
	for _, key := range d.keys {
		assert.Msg(key).
			This(fb.Search(key)).Is(d.k2o[key])
	}

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
