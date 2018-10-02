package lexer

//go:generate stringer -type=Token

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
	DEC
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
	BOOL
	BREAK
	BUFFER
	CALLBACK
	CASE
	CATCH
	CHAR
	CLASS
	CONTINUE
	CREATE
	DEFAULT
	DLL
	DO
	DOUBLE
	ELSE
	FALSE
	FLOAT
	FOR
	FOREVER
	FUNCTION
	GDIOBJ
	HANDLE
	IF
	INT64
	LONG
	RESOURCE
	RETURN
	SHORT
	STRUCT
	SWITCH
	SUPER
	THIS
	THROW
	TRUE
	TRY
	VOID
	WHILE
	// query keywords
	ALTER
	AVERAGE
	CASCADE
	COUNT
	DELETE
	DROP
	ENSURE
	EXTEND
	HISTORY
	INDEX
	INSERT
	INTERSECT
	INTO
	JOIN
	KEY
	LEFTJOIN
	LIST
	MAX
	MIN
	MINUS
	PROJECT
	REMOVE
	RENAME
	REVERSE
	SET
	SORT
	SUMMARIZE
	SVIEW
	TIMES
	TO
	TOTAL
	UNION
	UNIQUE
	UPDATE
	UPDATES
	VIEW
	WHERE
)

const Ntokens = int(WHERE + 1)

var keywords = map[string]Token{
	"and":      AND,
	"bool":     BOOL,
	"break":    BREAK,
	"buffer":   BUFFER,
	"callback": CALLBACK,
	"case":     CASE,
	"catch":    CATCH,
	"char":     CHAR,
	"class":    CLASS,
	"continue": CONTINUE,
	"default":  DEFAULT,
	"dll":      DLL,
	"do":       DO,
	"double":   DOUBLE,
	"else":     ELSE,
	"false":    FALSE,
	"float":    FLOAT,
	"for":      FOR,
	"forever":  FOREVER,
	"function": FUNCTION,
	"gdiobj":   GDIOBJ,
	"handle":   HANDLE,
	"if":       IF,
	"in":       IN,
	"int64":    INT64,
	"is":       IS,
	"isnt":     ISNT,
	"long":     LONG,
	"new":      NEW,
	"not":      NOT,
	"or":       OR,
	"resource": RESOURCE,
	"return":   RETURN,
	"short":    SHORT,
	"string":   STRING,
	"struct":   STRUCT,
	"super":    SUPER,
	"switch":   SWITCH,
	"this":     THIS,
	"throw":    THROW,
	"true":     TRUE,
	"try":      TRY,
	"void":     VOID,
	"while":    WHILE,
	"xor":      ISNT,
}
