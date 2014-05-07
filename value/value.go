/*
Package value implements the value types for Suneido

The naming convention is that Suneido Value types start with "Su"
e.g. SuBool, SuInt, SuStr, etc.
*/
package value

import "github.com/apmckinlay/gsuneido/util/dnum"

// Value is used to reference a Suneido value
type Value interface {
	ToStr() string
	ToInt() int32
	ToDnum() dnum.Dnum
	Get(key Value) Value
	Put(key Value, val Value)
	String() string
	Equals(other interface{}) bool
	Hash() uint32
	// hash2 is used by object to shallow hash contents
	hash2() uint32
	TypeName() string
	order() ordering
	// cmp returns -1 for <, 0 for ==, +1 for >
	cmp(other Value) int // ops Cmp ensures other has same ordering

	// TODO add lookup that returns method
}

type ordering int

const (
	OrdBool ordering = iota
	OrdNum           // SuInt, SuDnum
	OrdStr           // SuStr, SuConcat
	OrdDate
	OrdObject
	OrdOther
)
