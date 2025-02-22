// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package opcodes defines the bytecode instructions
// generated by compiler and executed by runtime
package opcodes

//go:generate stringer -type=Opcode

// Where applicable there are matching methods in thread.go or ops.go
// e.g. t.Pop or ops.Add

type Opcode byte

const (
	Nop Opcode = iota

	// stack --------------------------------------------------------

	// Pop pops the top of the stack
	Pop

	// push values --------------------------------------------------

	// Int <int16> pushes an integer
	Int
	// Value <uint8> pushes a literal Value
	Value
	// True pushes True
	True
	// False pushes False
	False
	// Zero pushes Zero (0)
	Zero
	// One pushes One (1)
	One
	// MinusOne pushes One (1)
	MinusOne
	// MaxInt pushes int32 max, used by RangeTo and RangeLen
	MaxInt
	// EmptyStr pushes EmptyStr ("")
	EmptyStr
	// PushReturn <uint8> pushes multiple return values onto the stack
	PushReturn

	// load and store -----------------------------------------------

	// Load <uint8> pushes a local variable onto the stack
	Load
	// Store <uint8> assigns the top value into a local variable (no pop)
	Store
	// LoadStore <local uint8> <op uint8> replaces the top value
	// with ob[m] op= val
	LoadStore
	// Dyload <uint8> pushes a dynamic variable onto the stack
	// It looks up the frame stack to find it, and copies it locally
	Dyload
	// Global <uint16> pushes the value of a global name
	Global
	// Get replaces the top two values (ob & mem) with ob.Get(mem)
	Get
	// Put pops the top three values (ob, mem, val) and does ob.Put(mem, val)
	Put
	// GetPut <uint8> replaces the top 3 values (ob, mem, val)
	// with ob[m] op= val
	GetPut
	// RangeTo replaces the top three values (x,i,j) with x.RangeTo(i,j)
	RangeTo
	// RangeLen replaces the top three values (x,i,n) with x.RangeLen(i,n)
	RangeLen
	// This pushes frame.this
	This

	// operations ---------------------------------------------------

	// Is replaces the top two values with x is y
	Is
	// Isnt replaces the top two values with x isnt y
	Isnt
	// Match replaces the top two values with x =~ y
	Match
	// MatchNot replaces the top two values with x !~ y
	MatchNot
	// Lt replaces the top two values with x < y
	Lt
	// Lte replaces the top two values with x <= y
	Lte
	// Gt replaces the top two values with x > y
	Gt
	// Gte replaces the top two values with x >= y
	Gte
	// Add replaces the top two values with x + y
	Add
	// Sub replaces the top two values with x - y
	Sub
	// Cat replaces the top two values with x $ y (strings)
	Cat
	// Mul replaces the top two values with x * y
	Mul
	// Div replaces the top two values with x / y
	Div
	// Mod replaces the top two values with x % y
	Mod
	// LeftShift replaces the top two values with x << y (integers)
	LeftShift
	// RightShift replaces the top two values with x >> y (unsigned)
	RightShift
	// BitOr replaces the top two values with x | y (integers)
	BitOr
	// BitAnd replaces the top two values with x | y (integers)
	BitAnd
	// BitXor replaces the top two values with x & y (integers)
	BitXor
	// BitNot replaces the top value with ^y (integer)
	BitNot
	// Not replaces the top value with not x (logical)
	Not
	// UnaryPlus converts the top value to a number
	UnaryPlus
	// UnaryMinus replaces the top value with -x
	UnaryMinus
	// InRange orgOp, orgVal, endOp, endVal replaces the top value with true or false
	InRange

	// control flow -------------------------------------------------

	// Or <int16> jumps if top is true, else it pops and continues
	// panics if top is not True or false
	Or
	// And <int16> jumps if top is false, else it pops and continues
	// panics if top is not True or false
	And
	// Bool checks that top is True or False, else it panics
	Bool
	// QMark pops and if false jumps, else it continues
	// panics if top is not True or false
	QMark
	// In <int16> pops the top value and compares it to the next value
	// if equal it pops the second value and pushes True,
	// else it leaves the second value on the stack
	In
	// Cover is used for coverage
	Cover
	// Jump <int16> jumps to a relative location in the code
	Jump
	// JumpTrue <int16> pops and if true jumps, else it continues
	// panics if top is not True or False
	JumpTrue
	// JumpFalse <int16> pops and if false jumps, else it continues
	// panics if top is not True or False
	JumpFalse
	// JumpIs <int16> pops the top value and compares it to the next value
	// if equal it pops the second value and jumps
	// else it leaves the second value on the stack and continues
	// panics if top is not True or False
	JumpIs
	// JumpIsnt <int16> pops the top value and compares it to the next value
	// if not equal it leaves the second value on the stack
	// else it pops the second value on the stack and continues
	// panics if top is not True or False
	JumpIsnt
	// JumpLt <int16> jumps if the top is less than the next value (no pop)
	JumpLt
	// Iter replaces the top with top.Iter()
	Iter
	// Iter2 replaces the top with top.Iter2()
	Iter2
	// ForIn <uint8> <int16> calls top.Next()
	// if the result is not nil, it assigns and jumps
	// else it continues
	ForIn
	// ForIn2 <uint8> <uint8> <int16> calls top.Next2()
	// if the result is not nil, it assigns and jumps
	// else it continues
	ForIn2
	// ForCount <int16> increments top and jumps if greater than second
	ForRange
	// ForRangeVar <uint8> <int16> increments top, stores it,
	// and jumps if greater than second
	ForRangeVar

	// exceptions ---------------------------------------------------

	// Throw pops and panics with that value
	Throw
	// Try <int16> <uint8> registers the catch jump and the catch pattern
	// so we will start catching
	Try
	// Catch <int16> clears the catch information to stop catching
	// and jumps past the catch code
	Catch

	// call and return ----------------------------------------------

	// CallFuncDiscard <uint8> calls the function popped from the stack
	// with the specified StdArgSpecs or frame.fn.ArgsSpecs
	// and discards the result
	CallFuncDiscard

	// CallFuncNoNil <uint8> calls the function popped from the stack
	// with the specified StdArgSpecs or frame.fn.ArgsSpecs
	// and pushes the result which must not be nil
	CallFuncNoNil

	// CallFuncNilOk<uint8> calls the function popped from the stack
	// with the specified StdArgSpecs or frame.fn.ArgsSpecs
	// and pushes the result which may be nil (return special case)
	CallFuncNilOk

	// CallMethDiscard <uint8> calls the method popped from the stack
	// with the specified StdArgSpecs or frame.fn.ArgsSpecs
	// and discards the result
	CallMethDiscard

	// CallMethNoNil <uint8> calls the method popped from the stack
	// with the specified StdArgSpecs or frame.fn.ArgsSpecs
	// and pushes the result which must not be nil
	CallMethNoNil

	// CallMethNilOk <uint8> calls the method popped from the stack
	// with the specified StdArgSpecs or frame.fn.ArgsSpecs
	// and pushes the result which may be nil (return special case)
	CallMethNilOk

	// Super <uint16> specifies where to start the method lookup
	// for the following CallMeth
	Super
	// Return returns the top of the stack
	Return
	// ReturnNil returns nil i.e. no return value
	ReturnNil
	// ReturnThrow, forces caller to check value
	ReturnThrow
	// ReturnMulti <uint8> returns multiple values from the stack
	ReturnMulti

	// blocks -------------------------------------------------------

	// Closure <uint8> pushes a new closure block instance
	Closure
	// BlockBreak panics "block:break" (handled by application code)
	BlockBreak
	// BlockContinue panics "block:continue" (handled by application code)
	BlockContinue
	// BlockReturn panics "block return" (handled by runtime)
	BlockReturn
	// BlockReturnNil pushes nil and then does BlockReturn
	BlockReturnNil
)
