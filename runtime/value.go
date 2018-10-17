package runtime

import (
	"strconv"

	"github.com/apmckinlay/gsuneido/util/dnum"
)

type Callable interface {
	Call(t *Thread, as *ArgSpec) Value // raw args on stack
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
	if n, err := strconv.ParseInt(s, 0, 16); err == nil {
		return SuInt(int(n))
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
