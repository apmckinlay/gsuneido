// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	"bytes"

	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/types"
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

var lruCacheCallClass = func(th *Thread, args []Value) Value {
	fn := args[0]
	size := ToInt(args[1])
	okForResetAll := ToBool(args[2])
	return newSuLruCache(size, fn, okForResetAll)
}

var lruStaticMethods = methods("lruStatic")

var _ = staticMethod(lruStatic_ResetAll, "()")

func lruStatic_ResetAll(th *Thread, _ []Value) Value {
	iter := ToContainer(Global.Find(th, GnSuneido)).Iter2(true, true)
	for _, v := iter(); v != nil; _, v = iter() {
		if lc, ok := v.(*suLruCache); ok && lc.okForResetAll {
			lc.Reset()
		}
	}
	return nil
}

var _ = staticMethod(lruStatic_Members, "()")

func lruStatic_Members() Value {
	return lruStatic_members
}

var lruStatic_members = methodList(lruStaticMethods)

func (lc *suLruCacheGlobal) Lookup(th *Thread, method string) Value {
	if f, ok := lruStaticMethods[method]; ok {
		return f
	}
	return lc.SuBuiltin.Lookup(th, method) // for Params
}

func (*suLruCacheGlobal) String() string {
	return "LruCache /* builtin class */"
}

var suLruCacheMethods = methods("lru")

// Get calls the getter with exactly the same arguments it receives.
// If called with multiple arguments, the hash key is an @args object.
var _ = method(lru_Get, "(@args)")

func lru_Get(th *Thread, as *ArgSpec, this Value, args []Value) Value {
	if as.Nargs == 0 {
		panic("missing argument")
	}
	key := args[0]
	if as.Nargs > 1 {
		key = SuObjectOfArgs(args, as)
	}
	slc := this.(*suLruCache)
	val := slc.Fetch(key)
	if val == nil {
		val = slc.Fn.Call(th, nil, as) // call with existing stack args
		slc.Insert(key, val)
	}
	return val
}

var _ = method(lru_Reset, "()")

func lru_Reset(this Value) Value {
	slc := this.(*suLruCache)
	slc.Reset()
	return nil
}

var _ = method(lru_OkForResetAllQ, "()")

func lru_OkForResetAllQ(this Value) Value {
	slc := this.(*suLruCache)
	return SuBool(slc.okForResetAll)
}

var _ = method(lru_GetMissRate, "()")

func lru_GetMissRate(this Value) Value {
	slc := this.(*suLruCache)
	misses := IntVal(slc.Lc.misses)
	gets := IntVal(slc.Lc.hits + slc.Lc.misses)
	return OpDiv(misses, gets)
}

//-------------------------------------------------------------------

type suLruCache struct {
	ValueBase[*suLruCache]
	Fn Value
	Lc lruCache
	MayLock
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
		key.SetConcurrent()
		val.SetConcurrent()
	}
	slc.Lc.Put(key, val)
}

func (slc *suLruCache) Reset() {
	if slc.Lock() {
		defer slc.Unlock()
	}
	slc.Lc.Reset()
}

// Value implementation

var _ Value = (*suLruCache)(nil)

func (*suLruCache) Type() types.Type {
	return types.LruCache
}

func (slc *suLruCache) Equal(other any) bool {
	return slc == other
}

func (slc *suLruCache) SetConcurrent() {
	if slc.SetConc() {
		for _, e := range slc.Lc.entries {
			e.key.SetConcurrent()
			e.val.SetConcurrent()
		}
	}
}

func (*suLruCache) Lookup(_ *Thread, method string) Value {
	return suLruCacheMethods[method]
}

//-------------------------------------------------------------------

type lruCache struct {
	lru     []uint8 // uint8 means max size of 256
	entries []entry
	hm      HmapValue
	size    int
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
	v, ok := lc.hm.Get(key)
	if !ok {
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
