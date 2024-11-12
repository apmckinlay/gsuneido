// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"cmp"
	"fmt"
	"log"
	"slices"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/compile/lexer"
	"github.com/apmckinlay/gsuneido/core/types"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/hmap"
	"github.com/apmckinlay/gsuneido/util/generic/slc"
	"github.com/apmckinlay/gsuneido/util/pack"
	"github.com/apmckinlay/gsuneido/util/varint"
)

/*
WARNING: Lock is NOT reentrant
Methods that call Lock must not call other methods that call Lock.
RLock is reentrant, but must not call anything that might Lock.
The convention is that public methods should lock (if concurrent)
and private methods should not lock
*/

type HmapValue = hmap.Hmap[Value, Value, hmap.Meth[Value]]

// EmptyObject is a readonly empty SuObject
var EmptyObject = &SuObject{readonly: true}

// SuObject is a Suneido object
// i.e. a container with both list and named members.
// Zero value is a valid empty object.
type SuObject struct {
	ValueBase[SuObject]
	defval Value
	// copyCount is used to implement copy-on-write.
	// If there are no other copies the count will be zero (or nil).
	// The count is incremented when a lazy/deferred copy is made (in slice).
	// If count > 0 then the list and named must be copied before updating.
	// This is handled by startMutate.
	// It must be atomic because it may be shared by multiple objects
	// and therefore is not guarded by the object lock.
	copyCount *atomic.Int32
	list      []Value
	named     HmapValue
	rwMayLock
	// version is incremented by operations that change one of the sizes.
	// i.e. not by just updating a value in-place.
	// It is used to detect modification during iteration.
	version int32
	// clock is incremented by any modification, including in-place updates.
	// It is used to detect modification during packing.
	clock    uint32
	readonly bool
	sorting  bool
}

const obSizeLimit = 64_000 // ??? // LibLocate.list is 52,000

// NewSuObject creates an SuObject from a slice of Value's
func NewSuObject(args []Value) *SuObject {
	return &SuObject{list: args}
}

// SuObjectOf returns an SuObject from its arguments
func SuObjectOf(args ...Value) *SuObject {
	return &SuObject{list: args}
}

// SuObjectOfStrs returns an SuObject from a slice of strings
func SuObjectOfStrs(strs []string) *SuObject {
	list := make([]Value, len(strs))
	for i, s := range strs {
		list[i] = SuStr(s)
	}
	return NewSuObject(list)
}

// SuObjectOfArgs is used by massage and lrucache
func SuObjectOfArgs(args []Value, as *ArgSpec) *SuObject {
	assert.That(as.Each == 0)
	unnamed := int(as.Nargs) - len(as.Spec)
	ob := SuObjectOf(slc.Clone(args[:unnamed])...)
	for i, ni := range as.Spec {
		ob.Set(as.Names[ni], args[unnamed+i])
	}
	return ob
}

func (ob *SuObject) Copy() Container {
	if ob.RLock() {
		defer ob.RUnlock()
	}
	return ob.slice(0)
}

// slice returns a copy-on-write of an object, omitting the first n list values
func (ob *SuObject) slice(n int) *SuObject {
	ob.incCopyCount()
	var list []Value
	if n < len(ob.list) {
		list = ob.list[n:]
	}
	return &SuObject{copyCount: ob.copyCount,
		defval: ob.defval, named: ob.named, list: list}
}

func (ob *SuObject) incCopyCount() {
	if ob.copyCount == nil {
		// SetConcurrent initializes copyCount
		// so if it's nil we're not concurrent, so it's ok to set it.
		ob.copyCount = new(atomic.Int32)
	}
	ob.copyCount.Add(1)
}

var _ Container = (*SuObject)(nil) // includes Value and Lockable

func (ob *SuObject) ToObject() *SuObject {
	return ob
}

// Get returns the value associated with a key, or defval if not found
func (ob *SuObject) Get(_ *Thread, key Value) Value {
	// not read-only because of container default values
	if ob.Lock() {
		defer ob.Unlock()
	}
	return ob.get(key)
	// val := ob.get(key)
	// if ob.concurrent && !ob.readonly && IsConcurrent(val) == False {
	// 	log.Println("ERROR: non-concurrent value in concurrent object",
	// 		key, val)
	// 	th.PrintStack()
	// }
	// return val
}

func (ob *SuObject) get(key Value) Value {
	if val := ob.getIfPresent(key); val != nil {
		return val
	}
	return ob.defaultValue(key)
}

func (ob *SuObject) defaultValue(key Value) Value {
	if ob.defval != nil {
		if d, ok := ob.defval.ToContainer(); ok {
			d = d.Copy()
			if !ob.readonly {
				ob.set(key, d)
			}
			return d
		}
	}
	return ob.defval
}

func (ob *SuObject) GetIfPresent(_ *Thread, key Value) Value {
	if ob.RLock() {
		defer ob.RUnlock()
	}
	return ob.getIfPresent(key)
}
func (ob *SuObject) getIfPresent(key Value) Value {
	if i, ok := key.IfInt(); ok && 0 <= i && i < len(ob.list) {
		return ob.list[i]
	}
	return ob.namedGet(key)
}

// HasKey returns true if the object contains the given key (not value)
func (ob *SuObject) HasKey(key Value) bool {
	if ob.RLock() {
		defer ob.RUnlock()
	}
	return ob.hasKey(key)
}
func (ob *SuObject) hasKey(key Value) bool {
	i, ok := key.IfInt()
	return (ok && 0 <= i && i < len(ob.list)) || ob.named.Get(key) != nil
}

// ListGet returns a value from the list, panics if index out of range
func (ob *SuObject) ListGet(i int) Value {
	if ob.RLock() {
		defer ob.RUnlock()
	}
	return ob.list[i]
}

// namedGet returns a named member or nil if it doesn't exist.
func (ob *SuObject) namedGet(key Value) Value {
	return ob.named.Get(key)
}

// Put adds or updates the given key and value
// The value will be added to the list if the key is the "next"
func (ob *SuObject) Put(_ *Thread, key, val Value) {
	if ob.Lock() {
		defer ob.Unlock()
	}
	ob.set(key, val)
}

func (ob *SuObject) GetPut(_ *Thread, m, v Value,
	op func(x, y Value) Value, retOrig bool) Value {
	if ob.Lock() {
		defer ob.Unlock()
	}
	orig := ob.get(m)
	if orig == nil {
		panic("uninitialized member: " + m.String())
	}
	v = op(orig, v)
	ob.set(m, v)
	if retOrig {
		return orig
	}
	return v
}

// Set implements Put, doesn't require thread.
// The value will be added to the list if the key is the "next"
func (ob *SuObject) Set(key, val Value) {
	if ob.Lock() {
		defer ob.Unlock()
	}
	ob.set(key, val)
}

func (ob *SuObject) CompareAndSet(key, newval, oldval Value) bool {
	if ob.Lock() {
		defer ob.Unlock()
	}
	if ob.getIfPresent(key) == oldval { // intentionally == not Equals
		ob.set(key, newval)
		return true
	}
	return false
}

// set implements Set without locking
func (ob *SuObject) set(key, val Value) {
	if ob.concurrent {
		key.SetConcurrent()
		val.SetConcurrent()
	}
	defer ob.endMutate(ob.startMutate())
	if i, ok := key.IfInt(); ok {
		if i == len(ob.list) {
			ob.add(val)
			return
		} else if 0 <= i && i < len(ob.list) {
			ob.list[i] = val
			return
		}
	}
	ob.named.Put(key, val)
	ob.ckSize()
}

func (ob *SuObject) ckSize() {
	if ob.size() > obSizeLimit {
		panic(fmt.Sprintf("object too large (> %d)", obSizeLimit))
	}
}

// Delete removes a key.
// If in the list, following list values are shifted over.
func (ob *SuObject) Delete(_ *Thread, key Value) bool {
	if ob.Lock() {
		defer ob.Unlock()
	}
	return ob.delete(key)
}
func (ob *SuObject) delete(key Value) bool {
	defer ob.endMutate(ob.startMutate())
	if i, ok := key.IfInt(); ok && 0 <= i && i < len(ob.list) {
		ob.listDelete(i)
		return true
	}
	return ob.named.Del(key) != nil
}

func (ob *SuObject) listDelete(i int) Value {
	x := ob.list[i]
	newlist := ob.list[:i+copy(ob.list[i:], ob.list[i+1:])]
	ob.list[len(ob.list)-1] = nil // aid garbage collection
	ob.list = newlist
	return x
}

// Erase removes a key.
// If in the list, following list values are NOT shifted over.
func (ob *SuObject) Erase(_ *Thread, key Value) bool {
	if ob.Lock() {
		defer ob.Unlock()
	}
	return ob.erase(key)
}
func (ob *SuObject) erase(key Value) bool {
	defer ob.endMutate(ob.startMutate())
	if i, ok := key.IfInt(); ok && 0 <= i && i < len(ob.list) {
		// migrate following list elements to named
		for j := len(ob.list) - 1; j > i; j-- {
			ob.named.Put(IntVal(j), ob.list[j])
			ob.list[j] = nil // aid garbage collection
		}
		ob.list = ob.list[:i]
		return true
	}
	return ob.named.Del(key) != nil
}

func (ob *SuObject) PopFirst() Value {
	if ob.Lock() {
		defer ob.Unlock()
	}
	if len(ob.list) < 1 {
		return nil
	}
	defer ob.endMutate(ob.startMutate())
	return ob.listDelete(0)
}

func (ob *SuObject) PopLast() Value {
	if ob.Lock() {
		defer ob.Unlock()
	}
	last := len(ob.list) - 1
	if last < 0 {
		return nil
	}
	defer ob.endMutate(ob.startMutate())
	return ob.listDelete(last)
}

// startMutate ensures the object is mutable (not readonly)
// and saves the list and named sizes (packed into an int).
// It also increments clock, assuming changes will be made.
func (ob *SuObject) startMutate() int {
	ob.mustBeMutable()
	ob.clock++
	return ob.sizes()
}

// endMutate increments the version if the sizes have changed
func (ob *SuObject) endMutate(nn int) {
	if nn != ob.sizes() {
		ob.version++
	}
}

func (ob *SuObject) sizes() int {
	return len(ob.list)<<32 | ob.named.Size()
}

// DeleteAll removes all the contents of the object, making it empty (size 0)
func (ob *SuObject) DeleteAll() {
	if ob.Lock() {
		defer ob.Unlock()
	}
	ob.deleteAll()
}
func (ob *SuObject) deleteAll() {
	defer ob.endMutate(ob.startMutate())
	ob.list = []Value{}
	ob.named = HmapValue{}
}

func (ob *SuObject) RangeTo(from int, to int) Value {
	if ob.RLock() {
		defer ob.RUnlock()
	}
	size := len(ob.list)
	from = prepFrom(from, size)
	to = prepTo(from, to, size)
	return ob.rangeTo(from, to)
}

func (ob *SuObject) RangeLen(from int, n int) Value {
	if ob.RLock() {
		defer ob.RUnlock()
	}
	size := len(ob.list)
	from = prepFrom(from, size)
	n = prepLen(n, size-from)
	return ob.rangeTo(from, from+n)
}

func (ob *SuObject) rangeTo(i int, j int) *SuObject {
	list := make([]Value, j-i)
	copy(list, ob.list[i:j])
	return &SuObject{list: list}
}

func (ob *SuObject) ToContainer() (Container, bool) {
	return ob, true
}

func (ob *SuObject) ListSize() int {
	if ob.RLock() {
		defer ob.RUnlock()
	}
	return len(ob.list)
}

func (ob *SuObject) NamedSize() int {
	if ob.RLock() {
		defer ob.RUnlock()
	}
	return ob.named.Size()
}

// Size returns the number of values in the object
func (ob *SuObject) Size() int {
	if ob.RLock() {
		defer ob.RUnlock()
	}
	return ob.size()
}

func (ob *SuObject) size() int {
	return len(ob.list) + ob.named.Size()
}

// Add appends a value to the list portion
func (ob *SuObject) Add(val Value) {
	if ob.concurrent {
		val.SetConcurrent()
	}
	if ob.Lock() {
		defer ob.Unlock()
	}
	ob.add(val)
}

// add implements Add without locking
func (ob *SuObject) add(val Value) {
	ob.mustBeMutable()
	ob.clock++
	ob.version++
	ob.list = append(ob.list, val)
	ob.migrate()
}

// Insert inserts at the given position.
// If the position is within the list, following values are moved over.
func (ob *SuObject) Insert(at int, val Value) {
	if ob.concurrent {
		val.SetConcurrent()
	}
	if ob.Lock() {
		defer ob.Unlock()
	}
	defer ob.endMutate(ob.startMutate())
	if 0 <= at && at <= len(ob.list) {
		// insert into list
		ob.list = append(ob.list, nil)
		copy(ob.list[at+1:], ob.list[at:])
		ob.list[at] = val
	} else {
		ob.set(IntVal(at), val)
	}
	ob.migrate()
}

func (ob *SuObject) mustBeMutable() {
	if ob.readonly {
		panic("can't modify readonly objects")
	}
	if ob.sorting {
		panic("can't modify object during sort")
	}
	if ob.copyCount != nil && ob.copyCount.Load() > 0 {
		// copy on write i.e. unshare, disconnect
		// two threads could (rarely) get in here at the same time
		// (on different objects, since we lock)
		ob.list = slc.Clone(ob.list)
		ob.named = *ob.named.Copy()
		// must do this last because if count goes to zero
		// then another thread could modify
		ob.copyCount.Add(-1)
		ob.copyCount = new(atomic.Int32)
	}
}

func (ob *SuObject) migrate() {
	for {
		x := ob.named.Del(IntVal(len(ob.list)))
		if x == nil {
			break
		}
		ob.list = append(ob.list, x)
	}
	ob.ckSize()
}

// vstack is used by Display and Show
// to track what is in progress to detect self reference
type vstack []*SuObject

func (vs *vstack) Push(ob *SuObject) bool {
	for _, v := range *vs {
		if v == ob { // deliberately == (pointers) not Equal (contents)
			return false
		}
	}
	*vs = append(*vs, ob) // push
	return true
}

func (ob *SuObject) String() string {
	return ob.Display(nil)
}

func (ob *SuObject) Display(th *Thread) string {
	// locking is handled by ArgsIter
	buf := limitBuf{}
	ob.rstring(th, &buf, nil)
	return buf.String()
}

func (ob *SuObject) rstring(th *Thread, buf *limitBuf, inProgress vstack) {
	ob.rstring2(th, buf, "#(", ")", inProgress)
}

func (ob *SuObject) rstring2(th *Thread, buf *limitBuf, before, after string,
	inProgress vstack) {
	if !inProgress.Push(ob) {
		buf.WriteString("...")
		return
	} // no pop necessary because we pass inProgress vstack slice by value
	buf.WriteString(before)
	sep := ""
	iter := ob.ArgsIter()
	for k, v := iter(); v != nil; k, v = iter() {
		buf.WriteString(sep)
		sep = ", "
		if k == nil {
			valstr(th, buf, v, inProgress)
		} else {
			entstr(th, buf, k, v, inProgress)
		}
	}
	buf.WriteString(after)
}

func entstr(th *Thread, buf *limitBuf, k Value, v Value, inProgress vstack) {
	if ks := Unquoted(k); ks != "" {
		buf.WriteString(string(ks))
	} else {
		valstr(th, buf, k, inProgress)
	}
	buf.WriteString(":")
	if v != True {
		buf.WriteString(" ")
		valstr(th, buf, v, inProgress)
	}
}

type recursable interface {
	rstring(th *Thread, buf *limitBuf, inProgress vstack)
}

var _ recursable = (*SuObject)(nil)

func valstr(th *Thread, buf *limitBuf, v Value, inProgress vstack) {
	if r, ok := v.(recursable); ok {
		r.rstring(th, buf, inProgress)
	} else {
		buf.WriteString(Display(th, v))
	}
}

func Unquoted(k Value) string {
	if ss, ok := k.(SuStr); ok {
		s := string(ss)
		// want true/false to be quoted to avoid ambiguity
		if (s != "true" && s != "false") && lexer.IsIdentifier(s) {
			return s
		}
	}
	return ""
}

// Show is like Display except that it sorts named members.
// It is used for tests.
func (ob *SuObject) Show() string {
	// locking is handled by ArgsIter
	return ob.show("#(", ")", nil)
}
func (ob *SuObject) show(before, after string, inProgress vstack) string {
	buf := &limitBuf{}
	buf.WriteString(before)
	type kv struct{ k, v Value }
	mems := []kv{}
	sep := ""
	iter := ob.ArgsIter()
	for k, v := iter(); v != nil; k, v = iter() {
		if k == nil {
			buf.WriteString(sep)
			sep = ", "
			valstr(nil, buf, v, inProgress)
		} else {
			mems = append(mems, kv{k, v})
		}
	}
	sort.Slice(mems,
		func(i, j int) bool { return mems[i].k.Compare(mems[j].k) < 0 })
	for _, p := range mems {
		buf.WriteString(sep)
		sep = ", "
		entstr(nil, buf, p.k, p.v, inProgress)
	}
	buf.WriteString(after)
	return buf.String()
}

const maxbuf = 64 * 1024

type limitBuf struct {
	sb strings.Builder
}

func (buf *limitBuf) WriteString(s string) {
	if buf.sb.Len()+len(s) > maxbuf {
		log.Panicln("buffer overflow displaying object >", maxbuf)
	}
	buf.sb.WriteString(s)
}

func (buf *limitBuf) String() string {
	return buf.sb.String()
}

func (ob *SuObject) Hash() uint64 {
	if ob.RLock() {
		defer ob.RUnlock()
	}
	hash := ob.hash2()
	if len(ob.list) > 0 {
		hash = 31*hash + ob.list[0].Hash2()
		if len(ob.list) > 1 {
			hash = 31*hash + ob.list[1].Hash2()
		}
	}
	if 0 < ob.named.Size() && ob.named.Size() <= 4 {
		iter := ob.named.Iter()
		for {
			k, v := iter()
			if k == nil {
				break
			}
			hash = 31*hash + k.Hash2()
			hash = 31*hash + v.Hash2()
		}
	}
	return hash
}

// Hash2 is shallow to prevents infinite recursion and deadlock
func (ob *SuObject) Hash2() uint64 {
	if ob.RLock() {
		defer ob.RUnlock()
	}
	return ob.hash2()
}

func (ob *SuObject) hash2() uint64 {
	hash := uint64(17)
	hash = 31*hash + uint64(ob.named.Size())
	hash = 31*hash + uint64(len(ob.list))
	return hash
}

func (ob *SuObject) Equal(other any) bool {
	// locking is handled by deepEqual
	val, ok := other.(Value)
	return ok && deepEqual(ob, val)
}

func (*SuObject) Type() types.Type {
	return types.Object
}

// Compare compares only list values (not named)
func (ob *SuObject) Compare(other Value) int {
	return deepCompare(ob, other)
}

func deepCompare(x Value, y Value) int {
	var tx, ty types.Type
	inProgress := make(inProgressStack, 0, 8) // 8 should handle most cases
	stack := make([]Value, 0, 32)             // 32 should handle most cases
	for {
		if x == y || inProgress.has(x, y) {
			goto endOfLoop
		}
		if x == nil {
			return -1
		}
		if y == nil {
			return +1
		}
		tx = order[x.Type()]
		ty = order[y.Type()]
		if tx != ty {
			return 2 * cmp.Compare(tx, ty)
		}
		switch tx {
		case types.Object:
			xTodo := children(ToContainer(x).ToObject(), &stack)
			yTodo := children(ToContainer(y).ToObject(), &stack)
			inProgress.push(x, y, xTodo, yTodo)
		default:
			cmp := x.Compare(y)
			if cmp != 0 {
				return cmp
			}
		}
	endOfLoop:
		x, y = inProgress.next()
		if x == nil && y == nil {
			return 0 // equal
		}
	}
}

func children(ob *SuObject, stack *[]Value) []Value {
	if !ob.RLock() {
		// not concurrent, don't need to copy
		return ob.list
	}
	defer ob.RUnlock()
	start := len(*stack)
	expand(stack, len(ob.list))
	copy((*stack)[start:], ob.list)
	return (*stack)[start:]
}

func (ob *SuObject) Call(th *Thread, _ Value, as *ArgSpec) Value {
	args := th.Args(&ParamSpec1, as)
	if x := ob.Get(th, args[0]); x != nil {
		return x
	}
	return args[0]
}

// ObjectMethods is initialized by the builtin package
var ObjectMethods Methods

var gnObjects = Global.Num("Objects")

func (*SuObject) Lookup(th *Thread, method string) Value {
	return Lookup(th, ObjectMethods, gnObjects, method)
}

// Slice returns a copy of the object, omitting the first n list values
func (ob *SuObject) Slice(n int) Container {
	if ob.RLock() {
		defer ob.RUnlock()
	}
	return ob.slice(n)
}

// Find returns the key of the first occurrence of the value else False.
// Lock to avoid object-modified-during-iteration.
func (ob *SuObject) Find(val Value) Value {
	if ob.RLock() {
		defer ob.RUnlock()
	}
	for i, v := range ob.list {
		if v.Equal(val) {
			return IntVal(i)
		}
	}
	named := ob.named.Iter()
	for k, v := named(); k != nil; k, v = named() {
		if v.Equal(val) {
			return k
		}
	}
	return False
}

// ArgsIter is similar to Iter2 but it returns a nil key for list elements
func (ob *SuObject) ArgsIter() func() (Value, Value) {
	if ob.RLock() {
		defer ob.RUnlock()
	}
	version := ob.version
	next := 0
	named := ob.named.Iter()
	return func() (Value, Value) {
		if ob.RLock() {
			defer ob.RUnlock()
		}
		ob.versionCheck(version)
		i := next
		next++
		if i < len(ob.list) {
			return nil, ob.list[i]
		}
		key, val := named()
		if key == nil {
			return nil, nil
		}
		return key, val
	}
}

// Iter2 iterates through list and named elements.
// List elements are returned with their numeric key index.
func (ob *SuObject) Iter2(list, named bool) func() (Value, Value) {
	if ob.RLock() {
		defer ob.RUnlock()
	}
	return ob.iter2(list, named)
}
func (ob *SuObject) iter2(list, named bool) func() (Value, Value) {
	version := ob.version
	next := 0
	if list && !named {
		return func() (Value, Value) {
			if ob.RLock() {
				defer ob.RUnlock()
			}
			ob.versionCheck(version)
			i := next
			if i < len(ob.list) {
				next++
				return IntVal(i), ob.list[i]
			}
			return nil, nil
		}
	}
	namedIter := ob.named.Iter()
	if named && !list {
		return func() (Value, Value) {
			if ob.RLock() {
				defer ob.RUnlock()
			}
			ob.versionCheck(version)
			key, val := namedIter()
			if key == nil {
				return nil, nil
			}
			return key, val
		}
	}
	// else named && list
	return func() (Value, Value) {
		if ob.RLock() {
			defer ob.RUnlock()
		}
		ob.versionCheck(version)
		i := next
		if i < len(ob.list) {
			next++
			return IntVal(i), ob.list[i]
		}
		key, val := namedIter()
		if key == nil {
			return nil, nil
		}
		return key, val
	}
}

func (ob *SuObject) versionCheck(version int32) {
	if ob.version != version {
		panic("object modified during iteration")
	}
}

func (ob *SuObject) clockCheck(clock uint32, op string) {
	if ob.clock != clock {
		panic("object modified during " + op)
	}
}

func (ob *SuObject) Iter() Iter {
	return &obIter{ob: ob, iter: ob.Iter2(true, true),
		result: func(k, v Value) Value { return v }}
}

func (ob *SuObject) ToRecord(th *Thread, hdr *Header) Record {
	if ob.RLock() {
		defer ob.RUnlock()
	}
	assert.That(len(hdr.Fields) == 1)
	fields := hdr.Fields[0]
	rb := RecordBuilder{}
	var tsField string
	var ts PackableValue
	for _, f := range fields {
		if strings.HasSuffix(f, "_TS") { // also done in SuRecord ToRecord
			tsField = f
			ts = th.Timestamp()
			rb.Add(ts)
		} else {
			x := ob.namedGet(SuStr(f))
			if x == nil {
				rb.AddRaw("")
			} else {
				rb.AddRaw(PackValue(x))
			}
		}
	}
	if tsField != "" && !ob.readonly {
		ob.set(SuStr(tsField), ts)
	}
	return rb.Trim().Build()
}

func (ob *SuObject) Sort(th *Thread, lt Value) {
	if ob.Lock() {
		defer ob.Unlock()
	}
	ob.mustBeMutable()
	ob.clock++
	ob.version++
	if lt == False {
		slices.SortStableFunc(ob.list, func(x, y Value) int {
			return x.Compare(y)
		})
	} else {
		func() {
			ob.sorting = true
			defer func() { ob.sorting = false }()
			ob.Unlock() // can't hold lock while calling arbitrary code
			defer ob.Lock()
			sort.SliceStable(ob.list, func(i, j int) bool {
				return ToBool(th.Call(lt, ob.list[i], ob.list[j]))
			})
			// note: could become concurrent while unlocked
		}()
	}
}

func (ob *SuObject) Unique() {
	if ob.Lock() {
		defer ob.Unlock()
	}
	defer ob.endMutate(ob.startMutate())
	if !ob.concurrent {
		ob.list = unique(ob.list)
	} else { // concurrent
		func() {
			ob.sorting = true
			defer func() { ob.sorting = false }()
			ob.Unlock() // can't hold lock while calling Equal
			defer ob.Lock()
			ob.list = unique(ob.list)
			// note: could become concurrent while unlocked
		}()
	}
}

func unique(list []Value) []Value {
	n := len(list)
	if n < 2 {
		return list
	}
	dst := 1
	for src := 1; src < n; src++ {
		if list[src].Equal(list[src-1]) {
			continue
		}
		if dst < src {
			list[dst] = list[src]
		}
		dst++
	}
	for i := dst; i < n; i++ {
		list[i] = nil // for gc
	}
	return list[:dst]
}

// SetConcurrent when readonly does not need to set shouldLock.
// But it still needs to set its children to concurrent
// because they may be SuRecords that require locking when readonly
// or other objects that require copyCount initialization.
func (ob *SuObject) SetConcurrent() {
	if !ob.concurrent {
		ob.concurrent = true
		if !ob.readonly {
			ob.shouldLock = true
		}
		ob.SetChildConc()
		if ob.copyCount == nil {
			ob.copyCount = new(atomic.Int32)
		}
	}
}

func (ob *SuObject) SetChildConc() {
	// recursive, deep
	for i := range len(ob.list) {
		ob.list[i].SetConcurrent()
	}
	iter := ob.named.Iter()
	for k, v := iter(); k != nil; k, v = iter() {
		k.SetConcurrent()
		v.SetConcurrent()
	}
	if ob.defval != nil {
		ob.defval.SetConcurrent()
	}
}

func (ob *SuObject) IsConcurrent() Value {
	if ob.concurrent {
		if !ob.shouldLock {
			return EmptyStr
		}
		return True
	}
	return False
}

func (ob *SuObject) SetReadOnly() {
	if !ob.setReadOnly() {
		return
	}
	// can't hold lock while accessing other objects
	iter := ob.Iter2(true, true)
	for k, v := iter(); k != nil; k, v = iter() {
		if x, ok := v.ToContainer(); ok {
			x.SetReadOnly() // recursive, deep
		}
	}
}

func (ob *SuObject) setReadOnly() bool {
	if ob.Lock() {
		defer ob.Unlock()
	}
	if ob.readonly {
		return false
	}
	ob.readonly = true
	return true
}

func (ob *SuObject) IsReadOnly() bool {
	if ob.Lock() {
		defer ob.Unlock()
	}
	return ob.isReadOnly()
}

func (ob *SuObject) isReadOnly() bool {
	return ob.readonly
}

func (ob *SuObject) SetDefault(def Value) {
	if def != nil && ob.concurrent {
		def.SetConcurrent()
	}
	if ob.Lock() {
		defer ob.Unlock()
	}
	ob.mustBeMutable()
	ob.defval = def
}

func (ob *SuObject) Reverse() {
	if ob.Lock() {
		defer ob.Unlock()
	}
	ob.mustBeMutable()
	ob.clock++
	ob.version++
	for lo, hi := 0, len(ob.list)-1; lo < hi; lo, hi = lo+1, hi-1 {
		ob.list[lo], ob.list[hi] = ob.list[hi], ob.list[lo]
	}
}

// BinarySearch does a binary search with default comparisons
func (ob *SuObject) BinarySearch(value Value) int {
	if ob.RLock() {
		defer ob.RUnlock()
	}
	return sort.Search(len(ob.list), func(i int) bool {
		return ob.list[i].Compare(value) >= 0
	})
}

// BinarySearch2 does a binary search with a user specified less than function
func (ob *SuObject) BinarySearch2(th *Thread, value, lt Value) int {
	ob.RLock()
	defer ob.RUnlock()
	defer ob.clockCheck(ob.clock, "BinarySearch")
	list := ob.list
	return sort.Search(len(list), func(i int) bool {
		ob.RUnlock() // can't hold lock while calling arbitrary code
		defer ob.RLock()
		return True != th.Call(lt, list[i], value)
		// note: could become concurrent during lt
	})
}

// Packable ---------------------------------------------------------

var _ Packable = (*SuObject)(nil)

func (ob *SuObject) PackSize(hash *uint64) int {
	return ob.PackSize2(hash, newPackStack())
}

func (ob *SuObject) PackSize2(hash *uint64, stack packStack) int {
	// must check stack before locking to avoid recursive deadlock
	stack.push(ob)
	if ob.RLock() {
		defer ob.RUnlock()
	}
	*hash = *hash*31 + uint64(ob.clock)
	if ob.size() == 0 {
		return 1 // just tag
	}
	ps := 1 // tag
	ps += varint.Len(uint64(len(ob.list)))
	for _, v := range ob.list {
		ps += packSize(v, hash, stack)
	}
	ps += varint.Len(uint64(ob.named.Size()))
	iter := ob.named.Iter()
	for k, v := iter(); k != nil; k, v = iter() {
		ps += packSize(k, hash, stack) + packSize(v, hash, stack)
	}
	return ps
}

func packSize(x Value, hash *uint64, stack packStack) int {
	if p, ok := x.(Packable); ok {
		n := p.PackSize2(hash, stack)
		return varint.Len(uint64(n)) + n
	}
	panic("can't pack " + ErrType(x))
}

func (ob *SuObject) Pack(hash *uint64, buf *pack.Encoder) {
	if ob.RLock() {
		defer ob.RUnlock()
	}
	ob.pack(hash, buf, PackObject)
}

func (ob *SuObject) pack(hash *uint64, buf *pack.Encoder, tag byte) {
	*hash = *hash*31 + uint64(ob.clock)
	buf.Put1(tag)
	if ob.size() == 0 {
		return
	}
	buf.VarUint(uint64(len(ob.list)))
	for _, v := range ob.list {
		packValue(v, hash, buf)
	}
	buf.VarUint(uint64(ob.named.Size()))
	iter := ob.named.Iter()
	for k, v := iter(); k != nil; k, v = iter() {
		packValue(k, hash, buf)
		packValue(v, hash, buf)
	}
}

func packValue(x Value, hash *uint64, buf *pack.Encoder) {
	buf0 := *buf
	buf.Put1(0) // 99% of the time we only need one byte for the size
	x.(Packable).Pack(hash, buf)
	n := len(buf.Buffer()) - len(buf0.Buffer()) - 1
	varlen := varint.Len(uint64(n))
	if varlen > 1 {
		// move what we just packed to make room for larger varint
		buf.Move(n, varlen-1)
	}
	buf0.VarUint(uint64(n))
}

func UnpackObject(s string) *SuObject {
	return unpackObject(s, &SuObject{})
}

func unpackObject(s string, ob *SuObject) *SuObject {
	if len(s) <= 1 {
		return ob
	}
	buf := pack.NewDecoder(s[1:])
	n := int(buf.VarUint())
	ob.list = make([]Value, n)
	for i := range n {
		ob.list[i] = unpackValue(buf)
	}
	n = int(buf.VarUint())
	for range n {
		k := unpackValue(buf)
		v := unpackValue(buf)
		ob.named.Put(k, v)
	}
	return ob
}

func unpackValue(buf *pack.Decoder) Value {
	size := int(buf.VarUint())
	return Unpack(buf.Get(size))
}

//-------------------------------------------------------------------

// rwMayLock is similar to MayLock except that it has a read-write lock.
type rwMayLock struct {
	lock sync.RWMutex
	// concurrent is whether it's accessible to multiple goroutines.
	// It is needed so SetChildConc is only called once.
	concurrent bool
	// shouldLock is whether we need to lock or not.
	// This is usually true when concurrent, except for SuObject's
	// that are already readonly when SetConcurrent is called
	shouldLock bool
}

func (x *rwMayLock) RLock() bool {
	if x == nil {
		log.Fatal("Lock nil")
	}
	if x.shouldLock {
		x.lock.RLock()
		return true
	}
	return false
}

func (x *rwMayLock) RUnlock() bool {
	if x.shouldLock {
		x.lock.RUnlock()
		return true
	}
	return false
}

func (x *rwMayLock) Lock() bool {
	if x == nil {
		log.Fatal("Lock nil")
	}
	if x.shouldLock {
		x.lock.Lock()
		return true
	}
	return false
}

func (x *rwMayLock) Unlock() bool {
	if x.shouldLock {
		x.lock.Unlock()
		return true
	}
	return false
}
