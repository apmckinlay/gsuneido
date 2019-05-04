package runtime

import (
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/runtime/types"
	"github.com/apmckinlay/gsuneido/util/dnum"
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
// - SuBlock
// - SuClass
// - SuInstance
// - SuMethod - not directly accessible, but returned for bound methods
// - SuIter - not directly accessible, but returned from e.g. object.Iter
type Value interface {
	// String returns a human readable string i.e. Suneido Display
	String() string

	// ToStr converts SuBool, SuInt, SuDnum, SuStr, SuConcat, SuExcept to string
	ToStr() (string, bool)

	// IfStr converts SuStr, SuConcat, SuExcept to string
	IfStr() (string, bool)

	// ToInt converts false (SuBool), "" (SuStr), SuInt, SuDnum to int
	ToInt() (int, bool)

	// IfInt converts SuInt, SuDnum to int
	IfInt() (int, bool)

	// ToDnum converts false (SuBool), "" (SuStr), SuInt, SuDnum to Dnum
	ToDnum() (dnum.Dnum, bool)

	// ToObject converts to an SuObject when applicable
	ToObject() (*SuObject, bool)

	// Get returns a member of an object/instance/class or a character of a string
	// returns nil if the member does not exist
	// The thread is necessary to call getters
	Get(t *Thread, key Value) Value

	Put(t *Thread, key Value, val Value)

	RangeTo(i int, j int) Value
	RangeLen(i int, n int) Value

	Equal(other interface{}) bool

	Hash() uint32

	// Hash2 is used by object to shallow hash contents
	Hash2() uint32

	// Type returns the Suneido name for the type
	Type() types.Type

	// Compare returns -1 for less, 0 for equal, +1 for greater
	Compare(other Value) int

	Callable

	Lookup(t *Thread, method string) Callable
}

// Callable is returned by Lookup
type Callable interface {
	Call(t *Thread, as *ArgSpec) Value
}

type Ord = int

// must match types
const (
	ordBool Ord = iota
	ordNum      // SuInt, SuDnum
	ordStr      // SuStr, SuConcat
	ordDate
	ordObject
	OrdOther
)

func Order(x Value) Ord {
	t := x.Type()
	if t <= types.Object {
		return Ord(t)
	} else if t == types.Record {
		return ordObject
	}
	return OrdOther
}

var NilVal Value

func NumFromString(s string) Value {
	if strings.HasPrefix(s, "0x") {
		if n, err := strconv.ParseUint(s, 0, 32); err == nil {
			return IntVal(int(int32(n)))
		}
	}
	if n, err := strconv.ParseInt(s, 0, 32); err == nil {
		return IntVal(int(n))
	}
	return SuDnum{Dnum: dnum.FromStr(s)}
}

type Showable interface {
	Show() string
}

// Show is .String() plus
// for classes it shows their contents
// for functions it shows their parameters
// for containers it sorts by member
func Show(v Value) string {
	if s, ok := v.(Showable); ok {
		return s.Show()
	}
	return v.String()
}

type Named interface {
	GetName() string
}

// ToStr converts SuBool, SuInt, SuDnum, SuStr, SuConcat, SuExcept to string.
// Calls Value.ToStr and panics if it fails
func ToStr(x Value) string {
	if s, ok := x.ToStr(); ok {
		return s
	}
	panic("can't convert " + x.Type().String() + " to String")
}

// IfStr converts SuStr, SuConcat, SuExcept to string.
// Calls Value.IfStr and panics if it fails
func IfStr(x Value) string {
	if s, ok := x.IfStr(); ok {
		return s
	}
	panic("can't convert " + x.Type().String() + " to String")
}

// ToStrOrString returns either IfStr() or String()
func ToStrOrString(x Value) string {
	if s, ok := x.IfStr(); ok {
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
	panic("can't convert " + errType(x) + " to integer")
}

// IfInt converts SuInt, SuDnum to int.
// Calls Value.IfInt and panics if it fails
func IfInt(x Value) int {
	if i, ok := x.IfInt(); ok {
		return i
	}
	panic("can't convert " + errType(x) + " to integer")
}

// ToDnum converts false (SuBool), "" (SuStr), SuInt, SuDnum to Dnum.
// Calls Value.ToDnum and panics if it fails
func ToDnum(x Value) dnum.Dnum {
	if dn, ok := x.ToDnum(); ok {
		return dn
	}
	panic("can't convert " + errType(x) + " to number")
}

// errType tweaks the TypeName to match cSuneido
func errType(x Value) string {
	if x == True {
		return "true"
	}
	t := x.Type().String()
	if t == "String" {
		return t
	}
	return strings.ToLower(t)
}

// ToObject converts to an SuObject or panics
func ToObject(x Value) *SuObject {
	if ob, ok := x.ToObject(); ok {
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

// Lookup looks for a method first in a methods map,
// and then in a global user defined class
// returning nil if not found in either place
func Lookup(t *Thread, methods Methods, gnUserDef int, method string) Callable {
	if m := methods[method]; m != nil {
		return m
	}
	if userdef := Global.Get(t, gnUserDef); userdef != nil {
		if c, ok := userdef.(*SuClass); ok {
			return c.get2(t, method)
		}
	}
	return nil
}

// CantConvert is embedded in Value types to supply default conversion methods
type CantConvert struct{}

func (CantConvert) ToInt() (int, bool) {
	return 0, false
}

func (CantConvert) IfInt() (int, bool) {
	return 0, false
}

func (CantConvert) ToDnum() (dnum.Dnum, bool) {
	return dnum.Zero, false
}

func (CantConvert) ToObject() (*SuObject, bool) {
	return nil, false
}

func (CantConvert) ToStr() (string, bool) {
	return "", false
}

func (CantConvert) IfStr() (string, bool) {
	return "", false
}

type ToStringable interface {
	ToString(*Thread) string
}
