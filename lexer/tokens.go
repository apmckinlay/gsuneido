package lexer

//go:generate stringer -type=Token

// to make stringer: go generate

// Token is returned by Lexer to identify the type of token
type Token uint8

// Keyword returns the token for a string it is a keyword
// otherwise IDENTIFIER and a copy of the string
func Keyword(s string) (Token, string) {
	if 2 <= len(s) && len(s) <= 8 && s[0] >= 'a' {
		for _, pair := range keywords {
			if pair.kw == s {
				return pair.tok, pair.kw
			}
		}
	}
	return IDENTIFIER, dup(s)
}

// Str returns a name for tokens that do not have a string value
func (t Token) Str() string {
	return tostring[t]
}

var tostring = map[Token]string{
	EOF:        "EOF",
	ERROR:      "ERROR",
	WHITESPACE: "WHITE",
	COMMENT:    "COMMENT",
	NEWLINE:    "NEWLINE",
}

const (
	NIL Token = iota
	EOF
	ERROR
	IDENTIFIER
	NUMBER
	STRING
	WHITESPACE
	COMMENT
	NEWLINE
	// punctuation
	HASH
	COMMA
	SEMICOLON
	AT
	L_PAREN
	R_PAREN
	L_BRACKET
	R_BRACKET
	L_CURLY
	R_CURLY
	RANGETO
	RANGELEN
	// operators
	NOT
	BITNOT
	NEW
	DOT
	IS
	ISNT
	MATCH
	MATCHNOT
	LT
	LTE
	GT
	GTE
	Q_MARK
	COLON
	ASSOC_START // must be consecutive
	AND
	OR
	BITOR
	BITAND
	BITXOR
	ADD
	SUB
	CAT
	MUL
	DIV
	ASSOC_END
	MOD
	LSHIFT
	RSHIFT
	INC
	POSTINC
	DEC
	POSTDEC
	ASSIGN_START // must be consecutive
	EQ
	ADDEQ
	SUBEQ
	CATEQ
	MULEQ
	DIVEQ
	MODEQ
	LSHIFTEQ
	RSHIFTEQ
	BITOREQ
	BITANDEQ
	BITXOREQ
	ASSIGN_END
	IN
	// other language keywords
	BREAK
	CASE
	CATCH
	CLASS
	CONTINUE
	DEFAULT
	DO
	ELSE
	FALSE
	FOR
	FOREVER
	FUNCTION
	IF
	RETURN
	SWITCH
	SUPER
	THIS
	THROW
	TRUE
	TRY
	WHILE
)

const Ntokens = int(WHILE + 1)

// keywords doesn't use a map because we want to use the keyword string literals
// ordered by frequency of use to optimize successful searches
var keywords = []struct {
	kw  string
	tok Token
}{
	{"return", RETURN},
	{"if", IF},
	{"false", FALSE},
	{"is", IS},
	{"true", TRUE},
	{"isnt", ISNT},
	{"and", AND},
	{"function", FUNCTION},
	{"for", FOR},
	{"in", IN},
	{"not", NOT},
	{"super", SUPER},
	{"or", OR},
	{"else", ELSE},
	{"class", CLASS},
	{"this", THIS},
	{"case", CASE},
	{"new", NEW},
	{"continue", CONTINUE},
	{"throw", THROW},
	{"try", TRY},
	{"catch", CATCH},
	{"while", WHILE},
	{"break", BREAK},
	{"switch", SWITCH},
	{"default", DEFAULT},
	{"do", DO},
	{"forever", FOREVER},
}

var IsIdent = [Ntokens]bool{
	IDENTIFIER: true,
	AND:        true,
	BREAK:      true,
	CASE:       true,
	CATCH:      true,
	CLASS:      true,
	CONTINUE:   true,
	DEFAULT:    true,
	DO:         true,
	ELSE:       true,
	FALSE:      true,
	FOR:        true,
	FOREVER:    true,
	FUNCTION:   true,
	IF:         true,
	IN:         true,
	IS:			true,
	ISNT:		true,
	NEW:		true,
	NOT:        true,
	OR:         true,
	RETURN:     true,
	SWITCH:     true,
	SUPER:      true,
	THIS:       true,
	THROW:      true,
	TRUE:       true,
	TRY:        true,
	WHILE:      true,
}
