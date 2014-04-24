package value

// Value is used to reference a Suneido value
type Value interface {
	ToStr() string
	ToInt() int // Q would it be better to make this int32 ?
	Get(key Value) Value
	Put(key Value, val Value)
}
