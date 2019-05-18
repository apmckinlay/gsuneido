package types

//go:generate stringer -type=Type

// to make stringer: go generate

type Type int

// must match Ord up to Object
const (
	Boolean Type = iota
	Number
	String
	Date
	Object
	Record
	Function
	Block
	BuiltinFunction
	Class
	Method
	Except
	Instance
	Iterator
	Transaction
	Query
	Cursor
)
