// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package testflathash

import (
	"math/rand"
	"testing"
)

//go:generate genny -in ../flathash/flathash.go -out pairhtbl.go -pkg testflathash gen "Key=int Item=Pair"

func TestPairHtbl(t *testing.T) {
	const n = 1200
	h := NewPairHtbl(0)
	data := make(map[int]int)
	check := func() {
		for k, v := range data {
			p := Pair{k, v}
			if *h.Get(k) != p {
				t.Error("expected", v, "for", k)
			}
		}
	}
	addRand := func() {
		key := rand.Intn(10000)
		val := rand.Int()
		h.Put(&Pair{key: key, val: val})
		data[key] = val
	}
	for i := 0; i < n; i++ {
		addRand()
		if i%11 == 0 {
			check()
		}
	}
	h = h.Dup()
	check()
}
