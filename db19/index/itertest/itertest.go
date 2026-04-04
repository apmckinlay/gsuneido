// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package itertest

import (
	"math/rand/v2"
	"slices"
	"sort"
	"testing"

	"github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/db19/index/iface"
	"github.com/apmckinlay/gsuneido/db19/index/ixkey"
	"github.com/apmckinlay/gsuneido/util/assert"
)

type KeyOff struct {
	Key string
	Off uint64
}

type Iter interface {
	Next()
	Prev()
	Rewind()
	SkipScan(prefixRn, suffixRng iface.Range, prefixLen int)
	HasCur() bool
	Eof() bool
	Cur() (key string, off uint64)
}

const npool = 16

var pool = func() []string {
	pool := make([]string, npool)
	pool[0] = "" // include empty string for better skip-scan coverage
	for i := 1; i < npool; i++ {
		pool[i] = core.Pack(core.IntVal(i - 1))
	}
	return pool
}()

func SkipScanTest(f *testing.F, makeIter func([]KeyOff) Iter) {
	// add some random cases so it does something as a test (when not fuzzing)
	for range 100 {
		f.Add(rand.Uint64(), randUint8(), rand.Uint32(), randBytes(100))
	}
	f.Fuzz(func(t *testing.T, seed uint64, u8 uint8, u32 uint32, steps []byte) {
		rng := rand.New(rand.NewPCG(seed, seed))

		nkeys := int(u8) // 0 to 255
		nfields, prefixLen, prefixRng, suffixRng := genCriteria(u32)
		allKeys := genData(rng, nfields, nkeys)
		it := makeIter(allKeys)
		// remove any zeroed (deleted) entries and sort to get the visible keys
		keys := slices.DeleteFunc(allKeys, func(k KeyOff) bool { return k.Off == 0 })
		sort.Slice(keys, func(i, j int) bool { return keys[i].Key < keys[j].Key })

		runSkipScan(t, it, keys, prefixRng, suffixRng, prefixLen, steps)
	})
}

func runSkipScan(t *testing.T, it Iter, keys []KeyOff,
	prefixRng, suffixRng iface.Range, prefixLen int, steps []byte) {
	t.Helper()
	it.SkipScan(prefixRng, suffixRng, prefixLen)
	or := newOracle(keys, prefixRng, suffixRng, prefixLen)

	for i := 0; !or.Eof(); i++ {
		assertSame(t, i, it, or)
		it.Next()
		or.Next()
	}
	assertSame(t, 9999, it, or)

	for i := 0; !or.Eof(); i++ {
		assertSame(t, i+10000, it, or)
		it.Prev()
		or.Prev()
	}
	assertSame(t, 19999, it, or)

	for step := range steps {
		switch step % 3 {
		case 0:
			it.Next()
			or.Next()
		case 1:
			it.Prev()
			or.Prev()
		case 2:
			it.Rewind()
			or.Rewind()
		}
		assertSame(t, step+20000, it, or)
	}
}

func randUint8() uint8 {
	return uint8(rand.N(256))
}

func randBytes(maxLen int) []byte {
	n := rand.N(maxLen)
	b := make([]byte, n)
	for i := 0; i < len(b); i += 8 {
		val := rand.Uint64()
		for j := 0; j < 8 && i+j < len(b); j++ {
			b[i+j] = byte(val)
			val >>= 8
		}
	}
	return b
}

func genData(rng *rand.Rand, nfields, nkeys int) []KeyOff {
	seen := map[string]struct{}{}
	data := make([]KeyOff, 0, nkeys)
	off := uint64(1)
	var enc ixkey.Encoder
	for len(data) < nkeys {
		for range nfields {
			enc.Add(pool[rng.IntN(len(pool))])
		}
		k := enc.String()
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		data = append(data, KeyOff{Key: k, Off: off})
		off++
	}
	if _, ok := seen[""]; !ok && nkeys > 0 && rng.IntN(2) == 1 {
		data[rng.IntN(nkeys)].Key = ""
	}
	return data
}

func genCriteria(x uint32) (nfields, preLen int, preRng iface.Range, sufRng iface.Range) {
	preEq := next(&x, 3)
	preCond := next(&x, 2)
	gap := next(&x, 3)
	if preEq+preCond+gap == 0 {
		gap = 1
	}
	sufEq := next(&x, 3)
	sufCond := next(&x, 2)
	if sufEq+sufCond == 0 {
		sufEq = 1
	}
	rem := next(&x, 2)

	nfields = int(preEq + preCond + gap + sufEq + sufCond + rem)
	preLen = int(preEq + preCond + gap)

	preRng = makeRange(&x, preEq, preCond)
	sufRng = makeRange(&x, sufEq, sufCond)
	return
}

func makeRange(x *uint32, eq, cond int) iface.Range {
	orgEnc := ixkey.Encoder{}
	endEnc := ixkey.Encoder{}
	for range eq {
		v := pool[next(x, npool)]
		orgEnc.Add(v)
		endEnc.Add(v)
	}
	if cond == 1 {
		orgEnc.Add(pool[next(x, npool/2)])
		endEnc.Add(pool[npool/2+next(x, npool/2)])
	}
	return iface.Range{Org: orgEnc.String(), End: endEnc.String()}
}

func next(src *uint32, mod uint32) int {
	result := *src % uint32(mod)
	*src /= mod
	return int(result)
}

//-------------------------------------------------------------------

type state byte

const (
	rewound state = iota
	within
	eof
)

type oracleIter struct {
	data []KeyOff
	cur  int
	state
}

func newOracle(keys []KeyOff, prefixRng, suffixRng iface.Range,
	prefixLen int) *oracleIter {
	out := make([]KeyOff, 0, len(keys))
	for _, k := range keys {
		pfx, sfx := ixkey.SplitPrefixSuffix(k.Key, prefixLen)
		if prefixRng.Org <= pfx && pfx < prefixRng.End &&
			suffixRng.Org <= sfx && sfx < suffixRng.End {
			out = append(out, k)
		}
	}
	return &oracleIter{data: out}
}

func (it *oracleIter) Rewind() {
	it.cur, it.state = 0, rewound
}

func (it *oracleIter) HasCur() bool {
	return it.state == within
}

func (it *oracleIter) Eof() bool {
	return it.state == eof
}

func (it *oracleIter) Cur() (string, uint64) {
	return it.data[it.cur].Key, it.data[it.cur].Off
}

func (it *oracleIter) Next() {
	if it.state == eof {
		return
	}
	if len(it.data) == 0 {
		it.state = eof
		return
	}
	if it.state == rewound {
		it.cur = 0
		it.state = within
		return
	}
	if it.cur+1 >= len(it.data) {
		it.state = eof
		return
	}
	it.cur++
}

func (it *oracleIter) Prev() {
	if it.state == eof {
		return
	}
	if len(it.data) == 0 {
		it.state = eof
		return
	}
	if it.state == rewound {
		it.cur = len(it.data) - 1
		it.state = within
		return
	}
	if it.cur == 0 {
		it.state = eof
		return
	}
	it.cur--
}

//-------------------------------------------------------------------

func assertSame(t *testing.T, step int, it Iter, or *oracleIter) {
	t.Helper()

	itHas := it.HasCur()
	orHas := or.HasCur()
	assert.T(t).Msg("step", step, "hascur or").This(itHas).Is(orHas)

	itEof := it.Eof()
	orEof := or.Eof()
	assert.T(t).Msg("step", step, "eof or", step).This(itEof).Is(orEof)

	if !(itHas && orHas) {
		return
	}

	itk, ito := it.Cur()
	ork, oro := or.Cur()
	assert.T(t).Msg("step", step, "key or").This(itk).Is(ork)
	assert.T(t).Msg("step", step, "off or").This(ito).Is(oro)
}

func SkipScanEmptyPrefixTest(t *testing.T, makeIter func([]KeyOff) Iter) {
	t.Helper()
	keys := []KeyOff{
		{Key: ixkey.CompKey("", "One"), Off: 1},
		{Key: ixkey.CompKey("", "Two"), Off: 2},
		{Key: ixkey.CompKey("", "Three"), Off: 3},
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].Key < keys[j].Key })
	runSkipScan(t, makeIter(keys), keys, iface.All, iface.Range{Org: "Two", End: "Two\x00"}, 1, nil)
}

func SkipScanEmptyPrefixPrevTest(t *testing.T, makeIter func([]KeyOff) Iter) {
	t.Helper()
	keys := []KeyOff{
		{Key: ixkey.CompKey("", "01"), Off: 1},
		{Key: ixkey.CompKey("", "03"), Off: 2},
		{Key: ixkey.CompKey("", "05"), Off: 3},
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].Key < keys[j].Key })
	runSkipScan(t, makeIter(keys), keys, iface.All, iface.Range{Org: "01", End: "04"}, 1, nil)
}

func SkipScanEmptyStringSuffixTest(t *testing.T, makeIter func([]KeyOff) Iter) {
	t.Helper()
	packStr := func(s string) string {
		if s == "" {
			return ""
		}
		return string(rune(core.PackString)) + s
	}
	key := func(name, path string) string {
		return ixkey.CompKey(packStr(name), packStr(path))
	}
	const groups = 50
	const extras = 20
	keys := make([]KeyOff, 0, groups*(extras+1))
	off := uint64(1)
	for i := range groups {
		name := "n" + twoDigits(i)
		keys = append(keys, KeyOff{Key: key(name, ""), Off: off})
		off++
		for j := range extras {
			path := "/x" + twoDigits(j)
			keys = append(keys, KeyOff{Key: key(name, path), Off: off})
			off++
		}
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i].Key < keys[j].Key })
	runSkipScan(t, makeIter(keys), keys, iface.All, iface.Range{Org: "", End: "\x00"}, 1, nil)
}

func twoDigits(n int) string {
	if n < 10 {
		return "0" + string(rune('0'+n))
	}
	return string([]rune{'0' + rune(n/10), '0' + rune(n%10)})
}
