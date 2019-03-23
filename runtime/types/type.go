package types

//go:generate stringer -type=Type

// to make stringer: go generate

type Type int

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
)
