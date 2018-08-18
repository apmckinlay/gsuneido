package base

import (
	"strconv"

	"github.com/apmckinlay/gsuneido/util/dnum"
)

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

	Equals(other interface{}) bool

	Hash() uint32

	// hash2 is used by object to shallow hash contents
	hash2() uint32

	// TypeName returns the Suneido name for the type
	TypeName() string

	Order() ord

	// Cmp returns -1 for <, 0 for ==, +1 for >
	Cmp(other Value) int // ops Cmp ensures other has same ordering
}

type ord int

const (
	ordBool ord = iota
	ordNum      // SuInt, SuDnum
	ordStr      // SuStr, SuConcat
	ordDate
	ordObject
	ordOther
)

var NilVal Value

func NumFromString(s string) Value {
	if n, err := strconv.ParseInt(s, 0, 16); err == nil {
		return SuInt(int(n))
	}
	return SuDnum{dnum.FromStr(s)}
}

// Index converts a value to an integer or else panics if not convertible
func Index(v Value) int {
	if i, ok := SmiToInt(v); ok {
		return i
	}
	if dn, ok := v.(SuDnum); ok {
		return dn.ToInt()
	}
	panic("indexes must be integers")
}
