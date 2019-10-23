package runtime

import (
	"sort"
	"strings"

	"sync"
	"sync/atomic"

	"github.com/apmckinlay/gsuneido/lexer"
	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/hmap"
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/pack"
	"github.com/apmckinlay/gsuneido/util/varint"
	// sync "github.com/sasha-s/go-deadlock"
)

/*
WARNING: sync.Mutex lock is NOT reentrant
Methods that lock must not call other methods that lock.
Public methods must lock (if concurrent)
Private methods must NOT lock
*/

// EmptyObject is a readonly empty SuObject
var EmptyObject = emptyOb()

func emptyOb() *SuObject {
	ob := NewSuObject()
	ob.SetReadOnly()
	return ob
}

// SuObject is a Suneido object
// i.e. a container with both list and named members.
// Zero value is a valid empty object.
//
// If concurrent is 0, no locking, assumed to be thread contained
// If concurrent is 1, guarded by lock, assumed to be shared
type SuObject struct {
	list       []Value
	named      hmap.Hmap
	readonly   bool
	concurrent int32 // access atomically
	version    int32
	clock      int32
	lock       sync.Mutex
	defval     Value
	CantConvert
}

// NewSuObject creates an SuObject from its arguments
func NewSuObject(args ...Value) *SuObject {
	ob := &SuObject{list: make([]Value, len(args))}
	for i, arg := range args {
		ob.list[i] = arg
	}
	return ob
}

func (ob *SuObject) Copy() Container {
	if ob.Lock() {
		defer ob.lock.Unlock()
	}
	ob2 := ob.slice(0)
	return &ob2
}

func (ob *SuObject) slice(n int) (ob2 SuObject) {
	ob2.named = *ob.named.Copy()
	ob2.defval = ob.defval
	if n < len(ob.list) {
		ob2.list = append(ob.list[:0:0], ob.list[n:]...)
	}
	return
}

var _ Container = (*SuObject)(nil)

func (ob *SuObject) ToObject() *SuObject {
	return ob
}

// Get returns the value associated with a key, or defval if not found
func (ob *SuObject) Get(t *Thread, key Value) Value {
	if val := ob.GetIfPresent(t, key); val != nil {
		return val
	}
	return ob.defaultValue(key)
}

func (ob *SuObject) defaultValue(key Value) Value {
	if ob.Lock() {
		defer ob.lock.Unlock()
	}
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
		defer ob.lock.Unlock()
	}
	if i, ok := key.IfInt(); ok && 0 <= i && i < len(ob.list) {
		return ob.list[i]
	}
	return ob.namedGet(key)
}

// Has returns true if the object contains the given key (not value)
func (ob *SuObject) HasKey(key Value) bool {
	if ob.Lock() {
		defer ob.lock.Unlock()
	}
	i, ok := key.IfInt()
	return (ok && 0 <= i && i < len(ob.list)) || ob.named.Get(key) != nil
}

// ListGet returns a value from the list, panics if index out of range
func (ob *SuObject) ListGet(i int) Value {
	if ob.Lock() {
		defer ob.lock.Unlock()
	}
	return ob.list[i]
}

// NamedGet returns a value from the named or nil if not found
func (ob *SuObject) NamedGet(key Value) Value {
	if ob.Lock() {
		defer ob.lock.Unlock()
	}
	return ob.namedGet(key)
}

func (ob *SuObject) namedGet(key Value) Value {
	val := ob.named.Get(key)
	if val == nil {
		return nil
	}
	return val.(Value)
}

// Put adds or updates the given key and value
// The value will be added to the list if the key is the "next"
func (ob *SuObject) Put(_ *Thread, key Value, val Value) {
	ob.Set(key, val)
}

// Set implements Put, doesn't require thread.
// The value will be added to the list if the key is the "next"
func (ob *SuObject) Set(key Value, val Value) {
	if ob.Lock() {
		defer ob.lock.Unlock()
	}
	ob.set(key, val)
}

func (ob *SuObject) set(key Value, val Value) {
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
		defer ob.lock.Unlock()
	}
	defer ob.endMutate(ob.startMutate())
	if i, ok := key.IfInt(); ok && 0 <= i && i < len(ob.list) {
		newlist := ob.list[:i+copy(ob.list[i:], ob.list[i+1:])]
		ob.list[len(ob.list)-1] = nil // aid garbage collection
		ob.list = newlist
		return true
	}
	return ob.named.Del(key) != nil
}

// Erase removes a key.
// If in the list, following list values are NOT shifted over.
func (ob *SuObject) Erase(_ *Thread, key Value) bool {
	if ob.Lock() {
		defer ob.lock.Unlock()
	}
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

// Clear removes all the contents of the object, making it empty (size 0)
func (ob *SuObject) Clear() {
	if ob.Lock() {
		defer ob.lock.Unlock()
	}
	defer ob.endMutate(ob.startMutate())
	ob.list = []Value{}
	ob.named = hmap.Hmap{}
}

func (ob *SuObject) RangeTo(from int, to int) Value {
	if ob.Lock() {
		defer ob.lock.Unlock()
	}
	size := ob.size()
	from = prepFrom(from, size)
	to = prepTo(from, to, size)
	return ob.rangeTo(from, to)
}

func (ob *SuObject) RangeLen(from int, n int) Value {
	if ob.Lock() {
		defer ob.lock.Unlock()
	}
	size := ob.size()
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
		defer ob.lock.Unlock()
	}
	return len(ob.list)
}

func (ob *SuObject) NamedSize() int {
	if ob.Lock() {
		defer ob.lock.Unlock()
	}
	return ob.named.Size()
}

// Size returns the number of values in the object
func (ob *SuObject) Size() int {
	if ob.Lock() {
		defer ob.lock.Unlock()
	}
	return ob.size()
}

func (ob *SuObject) size() int {
	return len(ob.list) + ob.named.Size()
}

// Add appends a value to the list portion
func (ob *SuObject) Add(val Value) {
	if ob.Lock() {
		defer ob.lock.Unlock()
	}
	ob.add(val)
}

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
		defer ob.lock.Unlock()
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
		ob.list = append(ob.list, x.(Value))
	}
}

func (ob *SuObject) String() string {
	if ob.Lock() {
		defer ob.lock.Unlock()
	}
	buf, sep := ob.vecstr()
	iter := ob.named.Iter()
	for {
		k, v := iter()
		if k == nil {
			break
		}
		sep = entstr(buf, k, v, sep)
	}
	buf.WriteString(")")
	return buf.String()
}

const maxbuf = 16 * 1024

func (ob *SuObject) vecstr() (*strings.Builder, string) {
	buf := strings.Builder{}
	sep := ""
	buf.WriteString("#(")
	for _, v := range ob.list {
		if buf.Len() > maxbuf {
			panic("buffer overflow displaying object")
		}
		buf.WriteString(sep)
		sep = ", "
		buf.WriteString(v.String())
	}
	return &buf, sep
}

func entstr(buf *strings.Builder, k interface{}, v interface{}, sep string) string {
	if buf.Len() > maxbuf {
		panic("buffer overflow displaying object")
	}
	buf.WriteString(sep)
	sep = ", "
	if ks, ok := k.(SuStr); ok && unquoted(string(ks)) {
		buf.WriteString(string(ks))
	} else {
		buf.WriteString(k.(Value).String())
	}
	buf.WriteString(":")
	if v != True {
		buf.WriteString(" ")
		buf.WriteString(v.(Value).String())
	}
	return sep
}

func unquoted(s string) bool {
	// want true/false to be quoted to avoid ambiguity
	return (s != "true" && s != "false") && lexer.IsIdentifier(s)
}

func (ob *SuObject) Show() string {
	if ob.Lock() {
		defer ob.lock.Unlock()
	}
	buf, sep := ob.vecstr()
	mems := []Value{}
	iter := ob.named.Iter()
	for {
		k, _ := iter()
		if k == nil {
			break
		}
		mems = append(mems, k.(Value))
	}
	sort.Slice(mems,
		func(i, j int) bool { return mems[i].Compare(mems[j]) < 0 })
	for _, k := range mems {
		v := ob.named.Get(k)
		sep = entstr(buf, k, v, sep)
	}
	buf.WriteString(")")
	return buf.String()
}

func (ob *SuObject) Hash() uint32 {
	if ob.Lock() {
		defer ob.lock.Unlock()
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
			hash = 31*hash + k.(Value).Hash2()
			hash = 31*hash + v.(Value).Hash2()
		}
	}
	return hash
}

// Hash2 is shallow so prevents infinite recursion
func (ob *SuObject) Hash2() uint32 {
	if ob.Lock() {
		defer ob.lock.Unlock()
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
	if val, ok := other.(Value); ok {
		return deepEqual(ob, val)
	}
	return false
}

func (*SuObject) Type() types.Type {
	return types.Object
}

// Compare compares only list values (not named)
func (ob *SuObject) Compare(other Value) int {
	if cmp := ints.Compare(ordObject, Order(other)); cmp != 0 {
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
			return ints.Compare(int(tx), int(ty))
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
	n := len(ob.list)
	expand(stack, n)
	start := len(*stack)
	for i := 0; i < n; i++ {
		(*stack)[start+i] = ob.list[i]
	}
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

// Slice returns a copy of the object, with the first n list elements removed
func (ob *SuObject) Slice(n int) Container {
	if ob.Lock() {
		defer ob.lock.Unlock()
	}
	ob2 := ob.slice(n)
	return &ob2
}

// ArgsIter is similar to Iter2 but it returns a nil key for list elements
func (ob *SuObject) ArgsIter() func() (Value, Value) {
	version := atomic.LoadInt32(&ob.version)
	next := 0
	named := ob.named.Iter()
	return func() (Value, Value) {
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
		return key.(Value), val.(Value)
	}
}

// Iter2 iterates through list and named elements.
// List elements are returned with their numeric key index.
func (ob *SuObject) Iter2(list, named bool) func() (Value, Value) {
	version := atomic.LoadInt32(&ob.version)
	next := 0
	if list && !named {
		return func() (Value, Value) {
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
			ob.modificationCheck(version)
			key, val := namedIter()
			if key == nil {
				return nil, nil
			}
			return key.(Value), val.(Value)
		}
	}
	return func() (Value, Value) {
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
		return key.(Value), val.(Value)
	}
}

func (ob *SuObject) modificationCheck(version int32) {
	if atomic.LoadInt32(&ob.version) != version {
		panic("object modified during iteration")
	}
}

func (ob *SuObject) Iter() Iter {
	return &obIter{ob: ob, iter: ob.Iter2(true, true),
		result: func(k, v Value) Value { return v }}
}

func (ob *SuObject) ToRecord(t *Thread, hdr *Header) Record {
	if ob.Lock() {
		defer ob.lock.Unlock()
	}
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
	return rb.Build()
}

func (ob *SuObject) Sort(t *Thread, lt Value) {
	if ob.Lock() {
		defer ob.lock.Unlock()
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
		defer ob.lock.Unlock()
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
	if atomic.LoadInt32(&ob.concurrent) == 1 ||
		ob.readonly { // don't need concurrent if readonly
		return
	}
	atomic.StoreInt32(&ob.concurrent, 1)
	iter := ob.Iter2(true, true)
	// recursive, deep
	for k, v := iter(); k != nil; k, v = iter() {
		k.SetConcurrent()
		v.SetConcurrent()
	}
}

func (ob *SuObject) SetReadOnly() {
	ob.Lock()
	if ob.readonly {
		ob.Unlock()
		return
	}
	ob.readonly = true
	ob.Unlock() // don't hold multiple locks
	iter := ob.Iter2(true, true)
	for k, v := iter(); k != nil; k, v = iter() {
		if x, ok := v.ToContainer(); ok {
			x.SetReadOnly() // recursive, deep
		}
	}
	atomic.StoreInt32(&ob.concurrent, 0) // don't need concurrent if readonly
}

func (ob *SuObject) IsReadOnly() bool {
	if ob.Lock() {
		defer ob.lock.Unlock()
	}
	return ob.readonly
}

func (ob *SuObject) SetDefault(def Value) {
	if ob.Lock() {
		defer ob.lock.Unlock()
	}
	ob.mustBeMutable()
	ob.defval = def
}

func (ob *SuObject) Reverse() {
	if ob.Lock() {
		defer ob.lock.Unlock()
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
		defer ob.lock.Unlock()
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

func (ob *SuObject) Lock() bool {
	if atomic.LoadInt32(&ob.concurrent) == 1 {
		ob.lock.Lock()
		return true
	}
	return false
}

func (ob *SuObject) Unlock() {
	if atomic.LoadInt32(&ob.concurrent) == 1 {
		ob.lock.Unlock()
	}
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
		defer func() {
			ob.lock.Unlock()
		}()
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

func packSize(x interface{}, clock int32, stack packStack) int {
	if p, ok := x.(Packable); ok {
		n := p.PackSize2(clock, stack)
		return varint.Len(uint64(n)) + n
	}
	panic("can't pack " + ErrType(x.(Value)))
}

func (ob *SuObject) PackSize3() int {
	if ob.Lock() {
		defer ob.lock.Unlock()
	}
	if ob.size() == 0 {
		return 1 // just tag
	}
	ps := 1 // tag
	ps += varint.Len(uint64(len(ob.list)))
	for _, v := range ob.list {
		ps += packSize3(v)
	}
	ps += varint.Len(uint64(ob.named.Size()))
	iter := ob.named.Iter()
	for k, v := iter(); k != nil; k, v = iter() {
		ps += packSize3(k) + packSize3(v)
	}
	return ps
}

func packSize3(x interface{}) int {
	if p, ok := x.(Packable); ok {
		n := p.PackSize3()
		return varint.Len(uint64(n)) + n
	}
	panic("can't pack " + ErrType(x.(Value)))
}

func (ob *SuObject) Pack(clock int32, buf *pack.Encoder) {
	if ob.Lock() {
		defer ob.lock.Unlock()
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

func packValue(x interface{}, clock int32, buf *pack.Encoder) {
	n := x.(Packable).PackSize3()
	buf.VarUint(uint64(n))
	x.(Packable).Pack(clock, buf)
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

// old format

func UnpackObjectOld(s string) *SuObject {
	return unpackObjectOld(s, &SuObject{})
}

func unpackObjectOld(s string, ob *SuObject) *SuObject {
	if len(s) <= 1 {
		return ob
	}
	buf := pack.NewDecoder(s[1:])
	var v Value
	n := buf.Int32()
	for i := 0; i < n; i++ {
		v = unpackValueOld(buf)
		ob.add(v)
	}
	var k Value
	n = buf.Int32()
	for i := 0; i < n; i++ {
		k = unpackValueOld(buf)
		v = unpackValueOld(buf)
		ob.set(k, v)
	}
	return ob
}

func unpackValueOld(buf *pack.Decoder) Value {
	size := buf.Int32()
	return UnpackOld(buf.Get(size))
}
