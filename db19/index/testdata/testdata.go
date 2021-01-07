// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package testdata

import (
	"math/rand"
	"sort"
	"strings"

	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/db19/stor"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/str"
)

type T = dat

type dat struct {
	keys []string
	o2k  map[uint64]string
	k2o  map[string]uint64
	off  uint64
}

func New() *dat {
	return &dat{
		o2k: map[uint64]string{},
		k2o: map[string]uint64{},
	}
}

func (d *dat) GetLeafKey(_ *stor.Stor, _ *ixkey.Spec, i uint64) string {
	return d.o2k[i]
}

func (d *dat) Len() int {
	return len(d.keys)
}

func (d *dat) Gen() (string, uint64) {
	for i := 0; i < 10; i++ {
		key := str.Random(4, 8)
		if _, ok := d.k2o[key]; !ok {
			off := d.NextOff()
			d.k2o[key] = off
			d.o2k[off] = key
			d.keys = append(d.keys, key)
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
	i := rand.Intn(len(d.keys))
	key := d.keys[i]
	off := d.k2o[key]
	return i, key, off
}

func (d *dat) Delete(i int) {
	last := len(d.keys) - 1
	d.keys[i] = d.keys[last]
	d.keys = d.keys[:last]
}

func (d *dat) Update(key string, off uint64) {
	// oldoff := d.k2o[key]
	d.k2o[key] = off
	// delete(d.o2k, oldoff)
	d.o2k[off] = key
}

type iter = func() (string, uint64, bool)

type tree interface {
	Lookup(key string) uint64
	Iter(bool) iter
}

func (d *dat) Check(fb tree) {
	for _, key := range d.keys {
		assert.Msg(key).
			This(fb.Lookup(key)).Is(d.k2o[key])
	}
	d.CheckIter(fb.Iter(true))
}

func (d *dat) CheckIter(it iter) {
	sort.Strings(d.keys)
	i := 0
	for k, o, ok := it(); ok; k, o, ok = it() {
		assert.Msg("expect prefix of " + d.keys[i] + " got " + k).
			That(strings.HasPrefix(d.keys[i], k))
		assert.This(d.o2k[o]).Is(d.keys[i])
		i++
	}
	assert.This(i).Is(len(d.keys))
}
