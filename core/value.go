// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package core

import (
	"log"
	"strconv"
	"sync"

	"github.com/apmckinlay/gsuneido/core/types"
	"github.com/apmckinlay/gsuneido/util/dnum"
	"github.com/apmckinlay/gsuneido/util/str"
	// sync "github.com/sasha-s/go-deadlock"
)

// Value is a value visible to Suneido programmers
// The naming convention is to use a prefix of "Su"
// - SuBoolean
// - SuInt, SuDnum - numbers
// - SuStr, SuConcat, SuExcept - strings
// - SuDate
// - SuObject, SuRecord, SuSequence - objects
// - SuBuiltin*, SuBuiltinMethod*
// - SuFunc
// - SuClosure
// - SuClass
// - SuInstance
// - SuMethod - not directly accessible, but returned for bound methods
// - SuIter - not directly accessible, but returned from e.g. object.Iter
type Value interface {
	// String returns a human readable string i.e. Suneido Display
	// Note: strings will have quotes and be escaped
	String() string

	// AsStr converts SuBool, SuInt, SuDnum, SuStr, SuConcat, SuExcept to string
	AsStr() (string, bool)

	// ToStr converts SuStr, SuConcat, SuExcept to string
	ToStr() (string, bool)

	// ToInt converts false (SuBool), "" (SuStr), SuInt, SuDnum to int
	ToInt() (int, bool)

	// IfInt converts SuInt, SuDnum to int
	IfInt() (int, bool)

	// ToDnum converts false (SuBool), "" (SuStr), SuInt, SuDnum to Dnum
	ToDnum() (dnum.Dnum, bool)

	// ToContainer converts object,record,sequence to a Container
	ToContainer() (Container, bool)

	// Get returns a member of an object/instance/class or a character of a string
	// returns nil if the member does not exist
	// The thread is necessary to call getters
	Get(th *Thread, key Value) Value

	// Put sets its key member to val.
	// Implemented by SuObject, SuRecord, and SuInstance.
	// t is required by SuRecord to call observers.
	Put(th *Thread, key Value, val Value)

	// GetPut is used for update operations like += and ++ atomically
	GetPut(th *Thread, key Value, val Value,
		op func(x, y Value) Value, retOrig bool) Value

	RangeTo(i int, j int) Value
	RangeLen(i int, n int) Value

	Equal(other any) bool

	Hash() uint64

	// Hash2 is used by object to shallow hash contents
	Hash2() uint64

	// Type returns the Suneido name for the type
	Type() types.Type

	// Compare returns -1 for less, 0 for equal, +1 for greater
	Compare(other Value) int

	Call(th *Thread, this Value, as *ArgSpec) Value

	// Lookup returns a Value or nil if the method isn't found
	Lookup(th *Thread, method string) Value

	// SetConcurrent is called when a Value is
	// about to become reachable by multiple threads.
	// At the point where SetConcurrent is called it is still thread contained.
	// SetConcurrent should be called on any other Value's reachable from this.
	// Additional calls to SetConcurrent should be ignored.
	// NOTE: SetConcurrent cannot call abitrary code
	// because it is called when holding a lock.
	SetConcurrent()
}

type Ord int

// must match types
const (
	ordBool Ord = iota
	ordNum      // SuInt, SuDnum
	ordStr      // SuStr, SuConcat, SuExcept
	ordDate
	ordObject
	ordOther
)

const OrdStr = ordStr

func Order(x Value) Ord {
	t := x.Type()
	if t <= types.Object {
		return Ord(t)
	} else if t == types.Except {
		return OrdStr
	} else if t == types.Record {
		return ordObject
	}
	return ordOther
}

func (o Ord) String() string {
	return []string{"boolean", "number", "string", "date", "object", "other"}[o]
}

var NilVal Value

// NumFromString converts a string to an SuInt or SuDnum.
// It will panic for invalid input.
func NumFromString(s string) Value {
	base := 10
	if isHex(s) {
		base = 0
	}
	if n, err := strconv.ParseInt(s, base, 64); err == nil {
		return IntVal(int(n))
	}
	return SuDnum{Dnum: dnum.FromStr(s)}
}

func isHex(s string) bool {
	if len(s) < 3 {
		return false
	}
	if s[0] == '-' {
		s = s[1:]
	}
	return len(s) > 2 && s[1] != '0' && s[1] == 'x'
}

type Showable interface {
	Show() string
}

// Show is .String() plus
// for classes it shows their contents
// for functions it shows their parameters
// for containers it sorts by member
func Show(v Value) string {
	if v == nil {
		return "nil"
	}
	if s, ok := v.(Showable); ok {
		return s.Show()
	}
	return v.String()
}

type Named interface {
	GetName() string
}

// AsStr converts SuBool, SuInt, SuDnum, SuStr, SuConcat, SuExcept to string.
// Calls Value.AsStr and panics if it fails
func AsStr(x Value) string {
	if s, ok := x.AsStr(); ok {
		return s
	}
	panic("can't convert " + x.Type().String() + " to String")
}

// ToStr converts SuStr, SuConcat, SuExcept to string.
// Calls Value.ToStr and panics if it fails
func ToStr(x Value) string {
	if s, ok := x.ToStr(); ok {
		return s
	}
	panic("can't convert " + x.Type().String() + " to String")
}

// ToStrOrString returns either ToStr() or String()
// i.e. strings won't have quotes
func ToStrOrString(x Value) string {
	if s, ok := x.ToStr(); ok {
		return s
	}
	return x.String()
}

// ToInt converts false (SuBool), "" (SuStr), SuInt, SuDnum to int.
// Calls Value.ToInt and panics if it fails
func ToInt(x Value) int {
	if i, ok := x.ToInt(); ok {
		return i
	}
	panic("can't convert " + ErrType(x) + " to integer")
}

// ToInt64 does ToDnum and ToInt64 and panics if it fails
func ToInt64(x Value) int64 {
	if i, ok := ToDnum(x).ToInt64(); ok {
		return i
	}
	panic("can't convert " + ErrType(x) + " to integer")
}

// IfInt converts SuInt, SuDnum to int.
// Calls Value.IfInt and panics if it fails
func IfInt(x Value) int {
	if i, ok := x.IfInt(); ok {
		return i
	}
	panic("can't convert " + ErrType(x) + " to integer")
}

// SuIntToInt converts to int if its argument is *smi or SuInt64.
// It is used by Equal methods.
func SuIntToInt(x any) (int, bool) {
	if si, ok := x.(*smi); ok {
		return si.toInt(), true
	}
	if si, ok := x.(SuInt64); ok {
		return int(si.int64), true
	}
	return 0, false
}

// ToDnum converts false (SuBool), "" (SuStr), SuInt, SuDnum to Dnum.
// Calls Value.ToDnum and panics if it fails
func ToDnum(x Value) dnum.Dnum {
	if dn, ok := x.ToDnum(); ok {
		return dn
	}
	panic("can't convert " + ErrType(x) + " to number")
}

// ErrType tweaks the Type to match cSuneido
func ErrType(x Value) string {
	if x == nil {
		return "nil"
	}
	if x == True {
		return "true"
	}
	if _, ok := x.(*SuSequence); ok {
		return "sequence"
	}
	t := x.Type().String()
	if t == "String" {
		return t
	}
	return str.ToLower(t)
}

// ToContainer converts to a Container or panics
func ToContainer(x Value) Container {
	if ob, ok := x.ToContainer(); ok {
		return ob
	}
	panic("can't convert " + x.Type().String() + " to Object")
}

func ToBool(x Value) bool {
	if x == True {
		return true
	}
	if x == False {
		return false
	}
	panic("can't convert " + x.Type().String() + " to Boolean")
}

func AsDate(v Value) (SuDate, bool) {
	if d, ok := v.(SuDate); ok {
		return d, true
	}
	if t, ok := v.(SuTimestamp); ok {
		return t.SuDate, true
	}
	return NilDate, false
}

// Lookup looks for a method first in a methods map,
// and then in a global user defined class
// returning nil if not found in either place
func Lookup(th *Thread, methods Methods, gnUserDef int, method string) Value {
	if m := methods[method]; m != nil {
		return m
	}
	return UserDef(th, gnUserDef, method)
}

func UserDef(th *Thread, gnUserDef int, method string) Value {
	if userdef := Global.Find(th, gnUserDef); userdef != nil {
		if c, ok := userdef.(*SuClass); ok {
			return c.get2(th, method, nil)
		}
	}
	return nil
}

type ToStringable interface {
	ToString(*Thread) string
}

// PackValue packs a Value if it is Packable, else it panics
func PackValue(v Value) string {
	if p, ok := v.(Packable); ok {
		return Pack(p)
	}
	panic("can't pack " + ErrType(v))
}

// PackSize returns the pack size of the value if it is Packable, else it panics
func PackSize(v Value) int {
	if p, ok := v.(Packable); ok {
		hash := uint64(17)
		return p.PackSize(&hash)
	}
	panic("can't pack " + ErrType(v))
}

type PackableValue interface {
	Value
	Packable
}

// IntVal returns an SuInt if it fits, else a SuDnum
func IntVal(n int) PackableValue {
	if MinSuInt <= n && n <= MaxSuInt {
		return SuInt(n)
	}
	return SuInt64{int64: int64(n)}
}

// Int64Val returns an SuInt if it fits, else a SuDnum
func Int64Val(n int64) PackableValue {
	if MinSuInt < n && n < MaxSuInt {
		return SuInt(int(n))
	}
	return SuInt64{int64: int64(n)}
}

// MayLock can be embedded to provide locking.
// concurrent is set *before* an object is shared
// and doesn't change after that
// so it should not require atomic or locked access.
// Before concurrent is set, no locking is done.
type MayLock struct {
	lock       sync.Mutex
	concurrent bool
}

func (x *MayLock) SetConcurrent() {
	x.concurrent = true
}

// SetConc returns true if it sets concurrent
func (x *MayLock) SetConc() bool {
	if x.concurrent {
		return false
	}
	x.concurrent = true
	return true
}

func (x *MayLock) Lock() bool {
	if x == nil {
		log.Fatal("Lock nil")
	}
	if x.concurrent {
		x.lock.Lock()
		return true
	}
	return false
}

func (x *MayLock) Unlock() bool {
	if x.concurrent {
		x.lock.Unlock()
		return true
	}
	return false
}

func (x *MayLock) IsConcurrent() Value {
	return SuBool(x.concurrent)
}

func IsConcurrent(x any) Value {
	if ic, ok := x.(interface{ IsConcurrent() Value }); ok {
		return ic.IsConcurrent()
	}
	return EmptyStr
}
