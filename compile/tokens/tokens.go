// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package tokens defines the lexical tokens
package tokens

//go:generate stringer -type=Token

// Token is returned by Lexer to identify the type of token
type Token uint8

const (
	Nil Token = iota
	Eof
	Error
	Identifier // use IsIdent to include keywords
	Number
	String
	// Symbol is a #identifier string
	Symbol
	Whitespace
	Comment
	Newline
	// punctuation
	Hash
	Comma
	Semicolon
	At
	LParen
	RParen
	LBracket
	RBracket
	LCurly
	RCurly
	RangeTo
	RangeLen
	OpsStart // operators
	Not
	BitNot
	New
	Dot
	CompareStart // must be consecutive
	Is
	Isnt
	Match
	MatchNot
	Lt // must be consecutive
	Lte
	Gt
	Gte
	CompareEnd
	QMark
	Colon
	AssocStart // must be consecutive
	And
	Or
	BitOr
	BitAnd
	BitXor
	Add
	Sub
	Cat
	Mul
	Div
	AssocEnd
	Mod
	LShift
	RShift
	Pipe
	IncDecStart
	Inc
	PostInc
	Dec
	PostDec
	IncDecEnd
	AssignStart // must be consecutive
	Eq
	AddEq
	SubEq
	CatEq
	MulEq
	DivEq
	ModEq
	LShiftEq
	RShiftEq
	BitOrEq
	BitAndEq
	BitXorEq
	AssignEnd
	In
	// other language keywords
	Break
	Case
	Catch
	Class
	Continue
	Default
	Do
	Else
	False
	For
	Forever
	Function
	If
	Return
	Switch
	Super
	This
	Throw
	True
	Try
	While
	// query keywords
	QueryStart
	SummarizeStart
	Average
	Count
	List
	Max
	Min
	Total
	SummarizeEnd
	Alter
	By
	Cascade
	Create
	Delete
	Drop
	Ensure
	Extend
	History
	Index
	Insert
	Intersect
	Into
	Join
	Key
	Leftjoin
	Lower
	Minus
	Project
	Remove
	Rename
	Reverse
	Set
	Sort
	Summarize
	Sview
	TempIndex
	Times
	To
	Union
	Unique
	Update
	View
	Where
	Ntokens
)

var isIdent = [Ntokens]bool{ // note: array rather than map
	Identifier: true,
	And:        true,
	Break:      true,
	Case:       true,
	Catch:      true,
	Class:      true,
	Continue:   true,
	Default:    true,
	Do:         true,
	Else:       true,
	False:      true,
	For:        true,
	Forever:    true,
	Function:   true,
	If:         true,
	In:         true,
	Is:         true,
	Isnt:       true,
	New:        true,
	Not:        true,
	Or:         true,
	Return:     true,
	Switch:     true,
	Super:      true,
	This:       true,
	Throw:      true,
	True:       true,
	Try:        true,
	While:      true,
}

// IsIdent returns whether a token is an identifier.
// The token must be within the valid range.
func (token Token) IsIdent() bool {
	return token > QueryStart || isIdent[token]
}

func (token Token) IsOperator() bool {
	return OpsStart < token && token < AssignStart
}

func (token Token) IsAssign() bool {
	return AssignStart < token && token < AssignEnd
}

func (token Token) IsIncDec() bool {
	return IncDecStart < token && token < IncDecEnd
}
