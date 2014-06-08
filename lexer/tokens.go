package compile

// Token is returned by Lexer to identify the type of token
type Token uint8

// Keyword returns the token for a string it is a keyword, else NIL
func Keyword(s string) Token {
	return keywords[s]
}

func (t Token) Closing() Token {
	switch t {
	case L_PAREN:
		return R_PAREN
	case L_CURLY:
		return R_CURLY
	case L_BRACKET:
		return R_BRACKET
	default:
		panic("invalid closing")
	}
}

// String returns a name for tokens that do not have a string value
func (t Token) String() string {
	return tostring[t]
}

var tostring = map[Token]string{
	NIL:        "NIL",
	EOF:        "EOF",
	ERROR:      "ERROR",
	WHITESPACE: "WHITE",
	COMMENT:    "COMMENT",
	NEWLINE:    "NEWLINE",

	STATEMENTS: "STMTS",
	POSTINC:    "POSTINC",
	POSTDEC:    "POSTDEC",
	FOR_IN:     "FOR_IN",
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
	// operators and punctuation
	HASH
	COMMA
	COLON
	SEMICOLON
	Q_MARK
	AT
	DOT
	L_PAREN
	R_PAREN
	L_BRACKET
	R_BRACKET
	L_CURLY
	R_CURLY
	IS
	ISNT
	MATCH
	MATCHNOT
	LT
	LTE
	GT
	GTE
	ADD
	SUB
	CAT
	MUL
	DIV
	MOD
	LSHIFT
	RSHIFT
	BITOR
	BITAND
	BITXOR
	NOT
	INC
	DEC // must be after INC
	BITNOT
	EQ
	ADDEQ // opEQ's must be in same order as op's
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
	RANGETO
	RANGELEN
	// langauge keywords
	AND
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
	IN
	INT64
	LONG
	NEW
	OR
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
	// for AST
	PARAMS
	STATEMENTS
	POSTINC
	POSTDEC
	VALUE
	INTERNAL
	FOR_IN
)

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

var infix = map[Token]bool{
	AND:      true,
	OR:       true,
	Q_MARK:   true,
	MATCH:    true,
	MATCHNOT: true,
	ADD:      true,
	SUB:      true,
	CAT:      true,
	MUL:      true,
	DIV:      true,
	MOD:      true,
	LSHIFT:   true,
	RSHIFT:   true,
	BITOR:    true,
	BITAND:   true,
	BITXOR:   true,
}
