package lexer

//go:generate stringer -type=Token

// to make stringer: go generate

// Token is returned by Lexer to identify the type of token
type Token uint8

// Keyword returns the token for a string it is a keyword, else NIL
func Keyword(s string) Token {
	return keywords[s]
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

var keywords = map[string]Token{
	"and":      AND,
	"break":    BREAK,
	"case":     CASE,
	"catch":    CATCH,
	"class":    CLASS,
	"continue": CONTINUE,
	"default":  DEFAULT,
	"do":       DO,
	"else":     ELSE,
	"false":    FALSE,
	"for":      FOR,
	"forever":  FOREVER,
	"function": FUNCTION,
	"if":       IF,
	"in":       IN,
	"is":       IS,
	"isnt":     ISNT,
	"new":      NEW,
	"not":      NOT,
	"or":       OR,
	"return":   RETURN,
	"string":   STRING,
	"super":    SUPER,
	"switch":   SWITCH,
	"this":     THIS,
	"throw":    THROW,
	"true":     TRUE,
	"try":      TRY,
	"while":    WHILE,
	"xor":      ISNT,
}
