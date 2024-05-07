// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package testdata

import (
	rand "math/rand/v2"
	"sort"
	"strings"

	"github.com/apmckinlay/gsuneido/db19/index/iterator"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

type T = dat

type dat struct {
	Keys []string
	O2k  map[uint64]string
	K2o  map[string]uint64
	off  uint64
}

func New() *dat {
	return &dat{
		O2k: map[uint64]string{},
		K2o: map[string]uint64{},
	}
}

func (d *dat) GetLeafKey(_ *stor.Stor, _ *ixkey.Spec, i uint64) string {
	return d.O2k[i]
}

func (d *dat) Len() int {
	return len(d.Keys)
}

func (d *dat) Gen() (string, uint64) {
	for i := 0; i < 10; i++ {
		key := str.Random(4, 8)
		if _, ok := d.K2o[key]; !ok {
			off := d.NextOff()
			d.K2o[key] = off
			d.O2k[off] = key
			d.Keys = append(d.Keys, key)
			return key, off
		}
	}
	panic("too many duplicates")
}

func (d *dat) NextOff() uint64 {
	d.off++
	return d.off
}

func (d *dat) Rand() (int, string, uint64) {
	i := rand.IntN(len(d.Keys))
	key := d.Keys[i]
	off := d.K2o[key]
	return i, key, off
}

func (d *dat) Delete(i int) {
	last := len(d.Keys) - 1
	d.Keys[i] = d.Keys[last]
	d.Keys = d.Keys[:last]
}

func (d *dat) Update(key string, off uint64) {
	// oldoff := d.k2o[key]
	d.K2o[key] = off
	// delete(d.o2k, oldoff)
	d.O2k[off] = key
}

type tree interface {
	Lookup(key string) uint64
}

func (d *dat) Check(t tree) {
	for _, key := range d.Keys {
		assert.Msg(key).This(t.Lookup(key)).Is(d.K2o[key])
	}
}

func (d *dat) CheckIter(it iterator.T) {
	sort.Strings(d.Keys)
	i := 0
	it.Rewind()
	for it.Next(); !it.Eof(); it.Next() {
		k, o := it.Cur()
		assert.Msg("expect prefix of " + d.Keys[i] + " got " + k).
			That(strings.HasPrefix(d.Keys[i], k))
		assert.This(d.O2k[o]).Is(d.Keys[i])
		i++
	}
	assert.This(i).Is(len(d.Keys))
}
