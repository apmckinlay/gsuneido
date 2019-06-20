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

	// ToContainer converts to a Container when applicable
	ToContainer() (Container, bool)

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
	Call(t *Thread, this Value, as *ArgSpec) Value
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

// ToDnum converts false (SuBool), "" (SuStr), SuInt, SuDnum to Dnum.
// Calls Value.ToDnum and panics if it fails
func ToDnum(x Value) dnum.Dnum {
	if dn, ok := x.ToDnum(); ok {
		return dn
	}
	panic("can't convert " + ErrType(x) + " to number")
}

// ErrType tweaks the TypeName to match cSuneido
func ErrType(x Value) string {
	if x == nil {
		return "nil"
	}
	if x == True {
		return "true"
	}
	t := x.Type().String()
	if t == "String" {
		return t
	}
	return strings.ToLower(t)
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

// Lookup looks for a method first in a methods map,
// and then in a global user defined class
// returning nil if not found in either place
func Lookup(t *Thread, methods Methods, gnUserDef int, method string) Callable {
	if m := methods[method]; m != nil {
		return m
	}
	return UserDef(t, gnUserDef, method)
}

func UserDef(t *Thread, gnUserDef int, method string) Callable {
	if userdef := Global.Get(t, gnUserDef); userdef != nil {
		if c, ok := userdef.(*SuClass); ok {
			return c.get2(t, method)
		}
	}
	return nil
}

// deepEqual is used by Container and SuInstance Equal
func deepEqual(x Value, y Value, inProgress pairs) bool {
	if xo, ok := x.ToContainer(); ok {
		if yo, ok := y.ToContainer(); ok {
			return containerEqual(xo, yo, inProgress)
		}
	}
	if xi, ok := x.(*SuInstance); ok {
		if yi, ok := y.(*SuInstance); ok {
			return instanceEqual(xi, yi, inProgress)
		}
	}
	return x.Equal(y)
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

func (CantConvert) ToContainer() (Container, bool) {
	return nil, false
}

func (CantConvert) AsStr() (string, bool) {
	return "", false
}

func (CantConvert) ToStr() (string, bool) {
	return "", false
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

// IntVal returns an SuInt if it fits, else a SuDnum
func IntVal(n int) Value {
	if MinSuInt <= n && n <= MaxSuInt {
		return SuInt(n)
	}
	return SuDnum{Dnum: dnum.FromInt(int64(n))}
}

// Int64Val returns an SuInt if it fits, else a SuDnum
func Int64Val(n int64) Value {
	if MinSuInt < n && n < MaxSuInt {
		return SuInt(int(n))
	}
	return SuDnum{Dnum: dnum.FromInt(n)}
}
