/*
Package value implements the value types for Suneido
*/
package value

import "github.com/apmckinlay/gsuneido/util/hmap"

// Value is used to reference a Suneido value
type Value interface {
	ToStr() string
	ToInt() int // Q would it be better to make this int32 ?
	Get(key Value) Value
	Put(key Value, val Value)
	String() string
	hmap.Key
	// Hash2 is used by object to shallow hash contents
	Hash2() uint32
}
