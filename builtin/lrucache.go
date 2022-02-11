// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"bytes"

	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/assert"
)

type suLruCacheGlobal struct {
	SuBuiltin
}

func init() {
	ps := params("(getfn, size=10, okForResetAll?=true)")
	Global.Builtin("LruCache", &suLruCacheGlobal{
		SuBuiltin{Fn: lruCacheCallClass,
			BuiltinParams: BuiltinParams{ParamSpec: ps}}})
}

var lruCacheCallClass = func(t *Thread, args []Value) Value {
	fn := args[0]
	size := ToInt(args[1])
	okForResetAll := ToBool(args[2])
	return newSuLruCache(size, fn, okForResetAll)
}

var lruCacheClassMethods = Methods{
	"ResetAll": method("()", func(th *Thread, this Value, _ []Value) Value {
		iter := ToContainer(Global.Find(th, GnSuneido)).Iter2(true, true)
		for _, v := iter(); v != nil; _, v = iter() {
			if lc, ok := v.(*suLruCache); ok && lc.okForResetAll {
				lc.Reset()
			}
		}
		return nil
	}),
}

func (d *suLruCacheGlobal) Lookup(t *Thread, method string) Callable {
	if f, ok := lruCacheClassMethods[method]; ok {
		return f
	}
	return d.SuBuiltin.Lookup(t, method) // for Params
}

func (d *suLruCacheGlobal) String() string {
	return "LruCache /* builtin class */"
}

//TODO merge GetN and GetN1 into Get using methodRaw

var suLruCacheMethods = Methods{
	// Get calls the getter with exactly the same arguments it receives.
	// If called with multiple arguments, the hash key is an @args object.
	//TODO after jSuneido is gone, we can replace Get,GetN,GetN1 with GetX
	"Get": methodRaw("(@args)", func(
		t *Thread, as *ArgSpec, this Value, args []Value) Value {
		if as.Nargs == 0 {
			panic("missing argument")
		}
		key := args[0]
		if as.Nargs > 1 {
			unnamed := int(as.Nargs) - len(as.Spec) // only valid if !atArg
			ob := &SuObject{}
			for i := 0; i < unnamed; i++ {
				ob.Add(args[i])
			}
			for i, ni := range as.Spec {
				ob.Set(as.Names[ni], args[unnamed+i])
			}
			key = ob
		}
		slc := this.(*suLruCache)
		val := slc.Fetch(key)
		if val == nil {
			val = slc.Fn.Call(t, nil, as) // call with existing stack args
			slc.Insert(key, val)
		}
		return val
	}),
	// "Get": method("(x)", func(t *Thread, this Value, args []Value) Value {
	// 	slc := this.(*suLruCache)
	// 	key := args[0]
	// 	val := slc.Fetch(key)
	// 	if val == nil {
	// 		val = t.Call(slc.Fn, key)
	// 		slc.Insert(key, val)
	// 	}
	// 	return val
	// }),
	"GetN": method("(@x)", func(t *Thread, this Value, args []Value) Value {
		slc := this.(*suLruCache)
		key := args[0]
		val := slc.Fetch(key)
		if val == nil {
			val = t.CallEach(slc.Fn, key)
			slc.Insert(key, val)
		}
		return val
	}),
	"GetN1": method("(x)", func(t *Thread, this Value, args []Value) Value {
		slc := this.(*suLruCache)
		key := args[0]
		val := slc.Fetch(key)
		if val == nil {
			val = t.CallEach(slc.Fn, key)
			slc.Insert(key, val)
		}
		return val
	}),
	"Reset": method0(func(this Value) Value {
		slc := this.(*suLruCache)
		slc.Reset()
		return nil
	}),
	"ResetOne": method1("(x)", func(this, arg Value) Value {
		slc := this.(*suLruCache)
		return SuBool(slc.Reset1(arg))
	}),
	"OkForResetAll?": method0(func(this Value) Value {
		slc := this.(*suLruCache)
		return SuBool(slc.okForResetAll)
	}),
	"GetMissRate": method0(func(this Value) Value {
		slc := this.(*suLruCache)
		misses := IntVal(slc.Lc.misses)
		gets := IntVal(slc.Lc.hits + slc.Lc.misses)
		return OpDiv(misses, gets)
	}),
}

//-------------------------------------------------------------------

type suLruCache struct {
	MayLock
	CantConvert
	Lc            lruCache
	Fn            Value
	okForResetAll bool
}

func newSuLruCache(size int, fn Value, okForResetAll bool) *suLruCache {
	return &suLruCache{Lc: *newLruCache(size), Fn: fn, okForResetAll: okForResetAll}
}

func (slc *suLruCache) Fetch(key Value) Value {
	if slc.Lock() {
		defer slc.Unlock()
	}
	return slc.Lc.Get(key)
}

func (slc *suLruCache) Insert(key, val Value) {
	if slc.Lock() {
		defer slc.Unlock()
	}
	slc.Lc.Put(key, val)
}

func (slc *suLruCache) Reset() {
	if slc.Lock() {
		defer slc.Unlock()
	}
	slc.Lc.Reset()
}

func (slc *suLruCache) Reset1(key Value) bool {
	if slc.Lock() {
		defer slc.Unlock()
	}
	return slc.Lc.Reset1(key)
}

// Value implementation

var _ Value = (*suLruCache)(nil)

func (*suLruCache) Get(*Thread, Value) Value {
	panic("LruCache does not support get")
}

func (*suLruCache) Put(*Thread, Value, Value) {
	panic("LruCache does not support put")
}

func (*suLruCache) GetPut(*Thread, Value, Value, func(x, y Value) Value, bool) Value {
	panic("LruCache does not support update")
}

func (*suLruCache) RangeTo(int, int) Value {
	panic("LruCache does not support range")
}

func (*suLruCache) RangeLen(int, int) Value {
	panic("LruCache does not support range")
}

func (*suLruCache) Hash() uint32 {
	panic("LruCache hash not implemented")
}

func (*suLruCache) Hash2() uint32 {
	panic("LruCache hash not implemented")
}

func (*suLruCache) Compare(Value) int {
	panic("LruCache compare not implemented")
}

func (*suLruCache) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call LruCache")
}

func (slc *suLruCache) String() string {
	return "anLruCache"
}

func (*suLruCache) Type() types.Type {
	return types.LruCache
}

func (slc *suLruCache) Equal(other interface{}) bool {
	slc2, ok := other.(*suLruCache)
	return ok && slc == slc2
}

func (slc *suLruCache) SetConcurrent() {
	if slc.SetConc() {
		for _, e := range slc.Lc.entries {
			e.key.SetConcurrent()
			e.val.SetConcurrent()
		}
	}
}

func (*suLruCache) Lookup(_ *Thread, method string) Callable {
	return suLruCacheMethods[method]
}

//-------------------------------------------------------------------

type lruCache struct {
	size    int
	hm      Hmap
	lru     []uint8 // uint8 means max size of 256
	entries []entry
	hits    int
	misses  int
}

type entry struct {
	key Value
	val Value
}

func newLruCache(req int) *lruCache {
	var sizes = []int{6, 13, 27, 55, 111, 223} // 7/8 of ^2
	size := 223                                // max
	for _, n := range sizes {
		if req <= n {
			size = n
			break
		}
	}
	return &lruCache{size: size,
		lru:     make([]uint8, 0, size),
		entries: make([]entry, 0, size),
	}
}

func (lc *lruCache) Get(key Value) Value {
	v := lc.hm.Get(key)
	if v == nil {
		// not in cache
		lc.misses++
		return nil
	}
	// in cache
	lc.hits++
	ei, _ := v.ToInt()
	li := bytes.IndexByte(lc.lru, uint8(ei))
	if li < lc.size-lc.size/8 {
		// move to the newest (the end)
		copy(lc.lru[li:], lc.lru[li+1:])
		lc.lru[len(lc.lru)-1] = uint8(ei)
	}
	return lc.entries[ei].val
}

func (lc *lruCache) Put(key, val Value) {
	ei := len(lc.entries)
	if ei < lc.size {
		lc.entries = append(lc.entries, entry{key: key, val: val})
		lc.lru = append(lc.lru, uint8(ei))
	} else { // full
		// replace oldest entry lru[0]
		ei = int(lc.lru[0])
		lc.hm.Del(lc.entries[ei].key)
		lc.entries[ei] = entry{key: key, val: val}
		copy(lc.lru, lc.lru[1:])
		lc.lru[lc.size-1] = uint8(ei)
	}
	lc.hm.Put(key, SuInt(ei))
}

func (lc *lruCache) GetPut(key Value, getfn func(key Value) Value) Value {
	val := lc.Get(key)
	if val == nil {
		val = getfn(key)
		lc.Put(key, val)
	}
	return val
}

func (lc *lruCache) Reset() {
	lc.hm.Clear()
	lc.lru = lc.lru[:0]
	lc.entries = lc.entries[:0]
	lc.hits = 0
	lc.misses = 0
}

func (lc *lruCache) Reset1(key Value) bool {
	v := lc.hm.Get(key)
	if v == nil {
		return false
	}
	lc.hm.Del(key)
	// delete from entries
	ei, _ := v.ToInt()
	copy(lc.entries[ei:], lc.entries[ei+1:])
	lc.entries = lc.entries[:len(lc.entries)-1]
	// delete from lru
	li := bytes.IndexByte(lc.lru, uint8(ei))
	copy(lc.lru[li:], lc.lru[li+1:])
	lc.lru = lc.lru[:len(lc.lru)-1]
	return true
}

// check is used by the test
func (lc *lruCache) check() {
	for _, ei := range lc.lru {
		e := lc.entries[ei]
		x := lc.hm.Get(e.key)
		assert.That(x != nil)
		xi, _ := x.ToInt()
		assert.That(xi == int(ei))
	}
	for ei, e := range lc.entries {
		x := lc.hm.Get(e.key)
		assert.That(x != nil)
		xi, _ := x.ToInt()
		assert.That(xi == int(ei))
	}
}

// func (lc *lruCache) print() {
// 	fmt.Println("lru")
// 	for li, ei := range lc.lru {
// 		fmt.Println(li, ei)
// 	}
// 	fmt.Println("entries")
// 	for ei, e := range lc.entries {
// 		fmt.Println(ei, e.key, e.val)
// 	}
// 	fmt.Println("hmap")
// 	it := lc.hm.Iter()
// 	for k, x := it(); k != nil; k, x = it() {
// 		fmt.Println(k, x)
// 	}
// }
