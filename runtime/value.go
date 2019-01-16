package runtime

import (
	"strconv"
	"strings"

	"github.com/apmckinlay/gsuneido/util/dnum"
)

// Value is used to reference a Suneido value
type Value interface {
	// String returns a human readable string i.e. Suneido Display
	String() string

	// ToStr converts to a string or panics
	// boolean and number are converted, other types are not
	ToStr() string

	// ToInt converts to an integer or panics
	// false and "" convert to 0 (but true does NOT convert to 1)
	ToInt() int

	// ToDnum converts to a dnum or panics
	// false and "" convert to 0 (but true does NOT convert to 1)
	ToDnum() dnum.Dnum

	// Get returns a member of an object/instance/class or a character of a string
	// returns nil if the member does not exist
	// The thread is necessary to call getters
	Get(t *Thread, key Value) Value

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

	Call(t *Thread, as *ArgSpec) Value

	Lookup(method string) Value
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
			return IntToValue(int(int32(n)))
		}
	}
	if n, err := strconv.ParseInt(s, 0, 32); err == nil {
		return IntToValue(int(n))
	}
	return SuDnum{dnum.FromStr(s)}
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
	SetName(name string)
}
