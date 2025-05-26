// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import (
	. "github.com/apmckinlay/gsuneido/core"
	"github.com/apmckinlay/gsuneido/core/types"
	"github.com/apmckinlay/gsuneido/util/lrucache"
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
	hits, misses := slc.Lc.Stats()
	return OpDiv(IntVal(misses), IntVal(hits+misses))
}

//-------------------------------------------------------------------

type suLruCache struct {
	ValueBase[*suLruCache]
	Fn Value
	Lc lrucache.Cache[Value, Value]
	MayLock
	okForResetAll bool
}

func newSuLruCache(size int, fn Value, okForResetAll bool) *suLruCache {
	return &suLruCache{Lc: *lrucache.New[Value, Value](size), Fn: fn, okForResetAll: okForResetAll}
}

func (slc *suLruCache) Fetch(key Value) Value {
	if slc.Lock() {
		defer slc.Unlock()
	}
	v, _ := slc.Lc.Get(key)
	return v
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
		for k, v := range slc.Lc.Entries() {
			k.SetConcurrent()
            v.SetConcurrent()
        }
	}
}

func (*suLruCache) Lookup(_ *Thread, method string) Value {
	return suLruCacheMethods[method]
}
