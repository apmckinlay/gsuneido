package runtime

import (
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/util/dnum"
)

type Callable interface {
	Call(t *Thread, as *ArgSpec) Value // raw args on stack
	Call0(t *Thread) Value
	Call1(t *Thread, a Value) Value
	Call2(t *Thread, a, b Value) Value
	Call3(t *Thread, a, b, c Value) Value
	Call4(t *Thread, a, b, c, d Value) Value
}

// Value is used to reference a Suneido value
type Value interface {
	// String returns a human readable string i.e. Suneido Display
	String() string

	// ToStr converts to a string
	ToStr() string

	ToInt() int

	ToDnum() dnum.Dnum

	Get(key Value) Value

	Put(key Value, val Value)

	RangeTo(i int, j int) Value
	RangeLen(i int, n int) Value

	Equal(other interface{}) bool

	Hash() uint32

	// Hash2 is used by object to shallow hash contents
	Hash2() uint32

	// TypeName returns the Suneido name for the type
	TypeName() string // or Value? (to avoid wrapping every time)

	Order() Ord

	// Compare returns -1 for less, 0 for equal, +1 for greater
	Compare(other Value) int

	Callable

	Lookup(method string) Callable
}

type Ord = int

const (
	ordBool Ord = iota
	ordNum      // SuInt, SuDnum
	ordStr      // SuStr, SuConcat
	ordDate
	ordObject
	OrdOther
)

var NilVal Value

func NumFromString(s string) Value {
	if strings.HasPrefix(s, "0x") {
		if n, err := strconv.ParseUint(s, 0, 32); err == nil {
			return NumFromInt(int(int32(n)))
		}
	}
	if n, err := strconv.ParseInt(s, 0, 32); err == nil {
		return NumFromInt(int(n))
	}
	return SuDnum{dnum.FromStr(s)}
}

// NumFromInt returns an SuInt if within range, else a SuDnum
func NumFromInt(n int) Value {
	if MinSuInt <= n && n <= MaxSuInt {
		return SuInt(n)
	}
	return SuDnum{dnum.FromInt(int64(n))}
}

type Showable interface {
	Show() string
}

// Show is .String() plus
// for classes it shows their contents
// for functions it shows their parameters
// for containers it sorts by member
func Show(v Value) string {
	if s,ok := v.(Showable); ok {
		return s.Show()
	}
	return v.String()
}
