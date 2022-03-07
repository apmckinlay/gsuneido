// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package runtime

import (
	"log"
	"sort"
	"strings"

	"sync/atomic"

	"github.com/apmckinlay/gsuneido/compile/lexer"
	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/assert"
	"github.com/apmckinlay/gsuneido/util/generic/ord"
	"github.com/apmckinlay/gsuneido/util/pack"
	"github.com/apmckinlay/gsuneido/util/varint"
)

/*
WARNING: sync.Mutex lock is NOT reentrant
Methods that lock must not call other methods that lock.
The convention is that public methods should lock (if concurrent)
and private methods should not lock
*/

// EmptyObject is a readonly empty SuObject
var EmptyObject = &SuObject{readonly: true}

// SuObject is a Suneido object
// i.e. a container with both list and named members.
// Zero value is a valid empty object.
type SuObject struct {
	ValueBase[SuObject]
	named  Hmap
	list   []Value
	defval Value
	MayLock
	// version is incrmented by operations that change one of the sizes.
	// i.e. not by just updating a value in-place.
	// It is used to detect modification during iteration.
	version int32
	// clock is incremented by any modification, including in-place updates.
	// It is used to detect modification during packing.
	clock    int32
	readonly bool
}

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

func (ob *SuObject) Copy() Container {
	if ob.Lock() {
		defer ob.Unlock()
	}
	return ob.slice(0)
}

// slice returns a copy of an object, omitting the first n list values
func (ob *SuObject) slice(n int) *SuObject {
	var list []Value
	if n < len(ob.list) {
		list = append(ob.list[:0:0], ob.list[n:]...) // copy
	}
	return &SuObject{defval: ob.defval, named: *ob.named.Copy(), list: list}
}

var _ Container = (*SuObject)(nil) // includes Value and Lockable

func (ob *SuObject) ToObject() *SuObject {
	return ob
}

// Get returns the value associated with a key, or defval if not found
func (ob *SuObject) Get(_ *Thread, key Value) Value {
	if ob.Lock() {
		defer ob.Unlock()
	}
	return ob.get(key)
	// val := ob.get(key)
	// if ob.concurrent && !ob.readonly && IsConcurrent(val) == False {
	// 	log.Println("ERROR: non-concurrent value in concurrent object",
	// 		key, val)
	// 	t.PrintStack()
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
	if ob.Lock() {
		defer ob.Unlock()
	}
	return ob.getIfPresent(key)
}
func (ob *SuObject) getIfPresent(key Value) Value {
	if i, ok := key.IfInt(); ok && 0 <= i && i < len(ob.list) {
		return ob.list[i]
	}
	return ob.namedGet(key)
}

// Has returns true if the object contains the given key (not value)
func (ob *SuObject) HasKey(key Value) bool {
	if ob.Lock() {
		defer ob.Unlock()
	}
	return ob.hasKey(key)
}
func (ob *SuObject) hasKey(key Value) bool {
	i, ok := key.IfInt()
	return (ok && 0 <= i && i < len(ob.list)) || ob.named.Get(key) != nil
}

// ListGet returns a value from the list, panics if index out of range
func (ob *SuObject) ListGet(i int) Value {
	if ob.Lock() {
		defer ob.Unlock()
	}
	return ob.list[i]
}

// namedGet returns a named member or nil if it doesn't exist.
func (ob *SuObject) namedGet(key Value) Value {
	val := ob.named.Get(key)
	if val == nil {
		return nil
	}
	return val
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

// set implements Set without locking
func (ob *SuObject) set(key, val Value) {
	if ob.concurrent {
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
// and saves the list and named sizes (packed into an int)
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
	ob.named = Hmap{}
}

func (ob *SuObject) RangeTo(from int, to int) Value {
	if ob.Lock() {
		defer ob.Unlock()
	}
	size := len(ob.list)
	from = prepFrom(from, size)
	to = prepTo(from, to, size)
	return ob.rangeTo(from, to)
}

func (ob *SuObject) RangeLen(from int, n int) Value {
	if ob.Lock() {
		defer ob.Unlock()
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
	if ob.Lock() {
		defer ob.Unlock()
	}
	return len(ob.list)
}

func (ob *SuObject) NamedSize() int {
	if ob.Lock() {
		defer ob.Unlock()
	}
	return ob.named.Size()
}

// Size returns the number of values in the object
func (ob *SuObject) Size() int {
	if ob.Lock() {
		defer ob.Unlock()
	}
	return ob.size()
}

func (ob *SuObject) size() int {
	return len(ob.list) + ob.named.Size()
}

// Add appends a value to the list portion
func (ob *SuObject) Add(val Value) {
	if ob.Lock() {
		defer ob.Unlock()
		val.SetConcurrent()
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
	if ob.Lock() {
		defer ob.Unlock()
		val.SetConcurrent()
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
}

func (ob *SuObject) migrate() {
	for {
		x := ob.named.Del(IntVal(len(ob.list)))
		if x == nil {
			break
		}
		ob.list = append(ob.list, x)
	}
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

func (ob *SuObject) Display(t *Thread) string {
	buf := limitBuf{}
	ob.rstring(t, &buf, nil)
	return buf.String()
}

func (ob *SuObject) rstring(t *Thread, buf *limitBuf, inProgress vstack) {
	ob.rstring2(t, buf, "#(", ")", inProgress)
}

func (ob *SuObject) rstring2(t *Thread, buf *limitBuf, before, after string,
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
			valstr(t, buf, v, inProgress)
		} else {
			entstr(t, buf, k, v, inProgress)
		}
	}
	buf.WriteString(after)
}

func (ob *SuObject) vecstr(t *Thread, buf *limitBuf, inProgress vstack) string {
	sep := ""
	for _, v := range ob.list {
		buf.WriteString(sep)
		sep = ", "
		valstr(t, buf, v, inProgress)
	}
	return sep
}

func entstr(t *Thread, buf *limitBuf, k Value, v Value, inProgress vstack) {
	if ks := Unquoted(k); ks != "" {
		buf.WriteString(string(ks))
	} else {
		valstr(t, buf, k, inProgress)
	}
	buf.WriteString(":")
	if v != True {
		buf.WriteString(" ")
		valstr(t, buf, v, inProgress)
	}
}

type recursable interface {
	rstring(t *Thread, buf *limitBuf, inProgress vstack)
}

var _ recursable = (*SuObject)(nil)

func valstr(t *Thread, buf *limitBuf, v Value, inProgress vstack) {
	if r, ok := v.(recursable); ok {
		r.rstring(t, buf, inProgress)
	} else {
		buf.WriteString(Display(t, v))
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

func (ob *SuObject) Hash() uint32 {
	if ob.Lock() {
		defer ob.Unlock()
	}
	hash := ob.hash2()
	if len(ob.list) > 0 {
		hash = 31*hash + ob.list[0].Hash()
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

// Hash2 is shallow so prevents infinite recursion
func (ob *SuObject) Hash2() uint32 {
	if ob.Lock() {
		defer ob.Unlock()
	}
	return ob.hash2()
}

func (ob *SuObject) hash2() uint32 {
	hash := uint32(17)
	hash = 31*hash + uint32(ob.named.Size())
	hash = 31*hash + uint32(len(ob.list))
	return hash
}

func (ob *SuObject) Equal(other interface{}) bool {
	val, ok := other.(Value)
	return ok && deepEqual(ob, val)
}

func (*SuObject) Type() types.Type {
	return types.Object
}

// Compare compares only list values (not named)
func (ob *SuObject) Compare(other Value) int {
	if cmp := ord.Compare(ordObject, Order(other)); cmp != 0 {
		return cmp
	}
	// now know other is an object so ToContainer won't panic
	return cmp2(ob, ToContainer(other).ToObject())
}

func cmp2(x Value, y Value) int {
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
			return ord.Compare(tx, ty)
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
	if !ob.Lock() {
		// not concurrent, don't need to copy
		return ob.list
	}
	defer ob.Unlock()
	start := len(*stack)
	expand(stack, len(ob.list))
	copy((*stack)[start:], ob.list)
	return (*stack)[start:]
}

func (*SuObject) Call(*Thread, Value, *ArgSpec) Value {
	panic("can't call Object")
}

// ObjectMethods is initialized by the builtin package
var ObjectMethods Methods

var gnObjects = Global.Num("Objects")

func (*SuObject) Lookup(t *Thread, method string) Callable {
	return Lookup(t, ObjectMethods, gnObjects, method)
}

// Slice returns a copy of the object, omitting the first n list values
func (ob *SuObject) Slice(n int) Container {
	if ob.Lock() {
		defer ob.Unlock()
	}
	return ob.slice(n)
}

// Find returns the key of the first occurence of the value else False.
// Lock to avoid object-modified-during-iteration.
func (ob *SuObject) Find(val Value) Value {
	if ob.Lock() {
		defer ob.Unlock()
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
	if ob.Lock() {
		defer ob.Unlock()
	}
	version := ob.version
	next := 0
	named := ob.named.Iter()
	return func() (Value, Value) {
		if ob.Lock() {
			defer ob.Unlock()
		}
		ob.modificationCheck(version)
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
	if ob.Lock() {
		defer ob.Unlock()
	}
	return ob.iter2(list, named)
}
func (ob *SuObject) iter2(list, named bool) func() (Value, Value) {
	version := ob.version
	next := 0
	if list && !named {
		return func() (Value, Value) {
			if ob.Lock() {
				defer ob.Unlock()
			}
			ob.modificationCheck(version)
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
			if ob.Lock() {
				defer ob.Unlock()
			}
			ob.modificationCheck(version)
			key, val := namedIter()
			if key == nil {
				return nil, nil
			}
			return key, val
		}
	}
	// else named && list
	return func() (Value, Value) {
		if ob.Lock() {
			defer ob.Unlock()
		}
		ob.modificationCheck(version)
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

func (ob *SuObject) modificationCheck(version int32) {
	if ob.version != version {
		panic("object modified during iteration")
	}
}

func (ob *SuObject) Iter() Iter {
	return &obIter{ob: ob, iter: ob.Iter2(true, true),
		result: func(k, v Value) Value { return v }}
}

func (ob *SuObject) ToRecord(t *Thread, hdr *Header) Record {
	if ob.Lock() {
		defer ob.Unlock()
	}
	assert.That(len(hdr.Fields) == 1)
	fields := hdr.Fields[0]
	rb := RecordBuilder{}
	var tsField string
	var ts SuDate
	for _, f := range fields {
		if strings.HasSuffix(f, "_TS") { // also done in SuRecord ToRecord
			tsField = f
			ts = t.Dbms().Timestamp()
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

func (ob *SuObject) Sort(t *Thread, lt Value) {
	if ob.Lock() {
		defer ob.Unlock()
	}
	ob.mustBeMutable()
	ob.clock++
	ob.version++
	if lt == False {
		sort.SliceStable(ob.list, func(i, j int) bool {
			return ob.list[i].Compare(ob.list[j]) < 0
		})
	} else {
		sort.SliceStable(ob.list, func(i, j int) bool {
			return True == t.Call(lt, ob.list[i], ob.list[j])
		})
	}
}

func (ob *SuObject) Unique() {
	if ob.Lock() {
		defer ob.Unlock()
	}
	defer ob.endMutate(ob.startMutate())
	n := len(ob.list)
	if n < 2 {
		return
	}
	dst := 1
	for src := 1; src < n; src++ {
		if ob.list[src].Equal(ob.list[src-1]) {
			continue
		}
		if dst < src {
			ob.list[dst] = ob.list[src]
		}
		dst++
	}
	for i := dst; i < n; i++ {
		ob.list[i] = nil // for gc
	}
	ob.list = ob.list[:dst]
}

func (ob *SuObject) SetConcurrent() {
	if ob.concurrent ||
		ob.readonly { // don't need concurrent if readonly
		return
	}
	ob.concurrent = true
	// recursive, deep
	for i := 0; i < len(ob.list); i++ {
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
		return True
	}
	if ob.readonly {
		return EmptyStr
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
	if ob.Lock() {
		defer ob.Unlock()
		if def != nil {
			def.SetConcurrent()
		}
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
	if ob.Lock() {
		defer ob.Unlock()
	}
	return sort.Search(len(ob.list), func(i int) bool {
		return ob.list[i].Compare(value) >= 0
	})
}

// BinarySearch2 does a binary search with a user specified less than function
func (ob *SuObject) BinarySearch2(t *Thread, value, lt Value) int {
	return sort.Search(ob.ListSize(), func(i int) bool {
		return True != t.Call(lt, ob.ListGet(i), value)
	})
}

// Packable ---------------------------------------------------------

var _ Packable = (*SuObject)(nil)

const packNestLimit = 20

func (ob *SuObject) PackSize(clock *int32) int {
	*clock = atomic.AddInt32(&packClock, 1)
	return ob.PackSize2(*clock, newPackStack())
}

func (ob *SuObject) PackSize2(clock int32, stack packStack) int {
	// must check stack before locking to avoid recursive deadlock
	stack.push(ob)
	if ob.Lock() {
		defer ob.Unlock()
	}
	ob.clock = clock
	if ob.size() == 0 {
		return 1 // just tag
	}
	ps := 1 // tag
	ps += varint.Len(uint64(len(ob.list)))
	for _, v := range ob.list {
		ps += packSize(v, clock, stack)
	}
	ps += varint.Len(uint64(ob.named.Size()))
	iter := ob.named.Iter()
	for k, v := iter(); k != nil; k, v = iter() {
		ps += packSize(k, clock, stack) + packSize(v, clock, stack)
	}
	return ps
}

func packSize(x Value, clock int32, stack packStack) int {
	if p, ok := x.(Packable); ok {
		n := p.PackSize2(clock, stack)
		return varint.Len(uint64(n)) + n
	}
	panic("can't pack " + ErrType(x))
}

func (ob *SuObject) Pack(clock int32, buf *pack.Encoder) {
	if ob.Lock() {
		defer ob.Unlock()
	}
	ob.pack(clock, buf, PackObject)
}

func (ob *SuObject) pack(clock int32, buf *pack.Encoder, tag byte) {
	if ob.clock != clock {
		panic("object modified during packing")
	}
	buf.Put1(tag)
	if ob.size() == 0 {
		return
	}
	buf.VarUint(uint64(len(ob.list)))
	for _, v := range ob.list {
		packValue(v, clock, buf)
	}
	buf.VarUint(uint64(ob.named.Size()))
	iter := ob.named.Iter()
	for k, v := iter(); k != nil; k, v = iter() {
		packValue(k, clock, buf)
		packValue(v, clock, buf)
	}
}

func packValue(x Value, clock int32, buf *pack.Encoder) {
	buf0 := *buf
	buf.Put1(0) // 99% of the time we only need one byte for the size
	x.(Packable).Pack(clock, buf)
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
	var v Value
	n := int(buf.VarUint())
	for i := 0; i < n; i++ {
		v = unpackValue(buf)
		ob.add(v)
	}
	var k Value
	n = int(buf.VarUint())
	for i := 0; i < n; i++ {
		k = unpackValue(buf)
		v = unpackValue(buf)
		ob.set(k, v)
	}
	return ob
}

func unpackValue(buf *pack.Decoder) Value {
	size := int(buf.VarUint())
	return Unpack(buf.Get(size))
}
