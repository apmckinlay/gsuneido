// Package tokens defines the lexical tokens
package tokens

//go:generate stringer -type=Token

// to make stringer: go generate

// Token is returned by Lexer to identify the type of token
type Token uint8

const (
	Nil Token = iota
	Eof
	Error
	Identifier
	Number
	String
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
	// operators
	Not
	BitNot
	New
	Dot
	Is
	Isnt
	Match
	MatchNot
	Lt
	Lte
	Gt
	Gte
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
	Inc
	PostInc
	Dec
	PostDec
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
)

const Ntokens = int(While + 1)

var isIdent = [Ntokens]bool{
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

// IsIdent returns whether a token is an identifier
// token must be within the valid range
func (token Token) IsIdent() bool {
	return isIdent[token]
}
