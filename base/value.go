package base

import (
	"errors"
	"math"
	"strconv"

	"github.com/apmckinlay/gsuneido/util/dnum"
)

// Value is used to reference a Suneido value
type Value interface {
	ToStr() string
	ToInt() int
	ToDnum() dnum.Dnum
	Get(key Value) Value
	Put(key Value, val Value)
	String() string
	Equals(other interface{}) bool
	Hash() uint32
	// hash2 is used by object to shallow hash contents
	hash2() uint32
	TypeName() string
	Order() ord
	// cmp returns -1 for <, 0 for ==, +1 for >
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

func NumFromString(s string) (Value, error) {
	n, err := strconv.ParseInt(s, 0, 16)
	if err == nil {
		return SuInt(int(n)), nil
	}
	dn, err := dnum.Parse(s)
	if err == nil {
		return DnumToValue(dn), nil
	}
	return NilVal, errors.New("invalid number: " + s)
}

// DnumToValue returns an SuInt if it fits, else a SuDnum
func DnumToValue(dn dnum.Dnum) Value {
	if dn.IsInt() {
		if n, err := dn.Int32(); err == nil &&
			math.MinInt16 <= n && n <= math.MaxInt16 {
			return SuInt(int(n))
		}
	}
	return SuDnum{dn}
}
