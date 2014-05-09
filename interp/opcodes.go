package interp

const (
	RETURN = iota
	PUSHINT
	PUSHVAL
	ADD
	SUB
	CAT
	MUL
	DIV
	MOD
	STORE
	LOAD
	UPLUS
	UMINUS
)
