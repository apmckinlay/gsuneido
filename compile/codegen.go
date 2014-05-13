package compile

import (
	i "github.com/apmckinlay/gsuneido/interp"
	"github.com/apmckinlay/gsuneido/util/varint"
	"github.com/apmckinlay/gsuneido/value"
)

// TODO fold constant expressions

// codegen compiles a Function from an Ast
func codegen(ast Ast) *value.SuFunc {
	//fmt.Println("codegen", ast.String())
	cg := cgen{}
	cg.gen(ast)
	cg.emit(i.RETURN)
	return &value.SuFunc{
		Code:    cg.code,
		Values:  cg.values,
		Strings: cg.strings,
		Nlocals: len(cg.strings), // ultimately WRONG!
	}
}

type cgen struct {
	code    []byte
	values  []value.Value
	strings []string
}

func (cg *cgen) gen(ast Ast) {
	//fmt.Println("gen", ast.String())
	switch ast.KeyTok() {
	case FUNCTION:
		cg.function(ast)
	case RETURN:
		if len(ast.Children) == 1 {
			cg.gen(ast.first())
		}
		cg.emit(i.RETURN)
	case EXPRESSION:
		cg.gen(ast.first())
	case NUMBER:
		cg.number(ast.Value)
	case STRING:
		cg.emit(i.VALUE)
		i := cg.value(value.SuStr(ast.Value))
		cg.emitUint(i)
	case TRUE:
		cg.emit(i.TRUE)
	case FALSE:
		cg.emit(i.FALSE)
	case VALUE:
		cg.emit(i.VALUE)
		i := cg.value(ast.val)
		cg.emitUint(i)
	case NOT:
		cg.unop(ast, i.NOT)
	case ADD:
		cg.ubinop(ast, i.UPLUS, i.ADD)
	case SUB:
		cg.ubinop(ast, i.UMINUS, i.SUB)
	case IS, ISNT, MATCH, MATCHNOT, LT, LTE, GT, GTE,
		CAT, MUL, DIV, MOD, LSHIFT, RSHIFT, BITOR, BITAND, BITXOR:
		cg.binop(ast, i.IS+byte(ast.KeyTok()-IS))
	case BITNOT:
		cg.unop(ast, i.BITNOT)
	case IDENTIFIER:
		cg.rvalue(ast)
	case EQ:
		cg.gen(ast.second())
		cg.store(cg.lvalue(ast.first()))
	case ADDEQ, SUBEQ, CATEQ, MULEQ, DIVEQ, MODEQ,
		LSHIFTEQ, RSHIFTEQ, BITOREQ, BITANDEQ, BITXOREQ:
		ref := cg.lvalue(ast.first())
		cg.load(ref)
		cg.gen(ast.second())
		cg.emit(i.ADD + byte(ast.Token-ADDEQ))
		cg.store(ref)
	case INC, DEC:
		ref := cg.lvalue(ast.first())
		cg.load(ref)
		cg.emit(i.INT)
		cg.emitInt(1)
		cg.emit(i.ADD + byte(ast.Token-INC))
		cg.store(ref)
	case POSTINC, POSTDEC:
		ref := cg.lvalue(ast.first())
		cg.load(ref)
		cg.emit(i.DUP)
		cg.emit(i.INT)
		cg.emitInt(1)
		cg.emit(i.ADD + byte(ast.Token-POSTINC))
		cg.store(ref)
		cg.emit(i.POP)
	default:
		panic("not implemented" + ast.String())
	}
}

func (cg *cgen) function(ast Ast) {
	// TODO params
	stmts := ast.second().Children
	for si, stmt := range stmts {
		cg.gen(stmt)
		if si < len(stmts)-1 {
			cg.emit(i.POP)
		}
	}
	// TODO add return if that wasn't last statement
}

func (cg *cgen) number(s string) {
	val, err := value.NumFromString(s)
	if err != nil {
		panic("invalid number: " + s)
	}
	if si, ok := val.(value.SuInt); ok {
		cg.emit(i.INT)
		cg.emitInt(int(si))
	} else {
		cg.emit(i.VALUE)
		cg.emitUint(cg.value(val))
	}
}

// value returns an index for the value
// reusing if duplicate, adding otherwise
func (cg *cgen) value(v value.Value) int {
	for i, v2 := range cg.values {
		if v.Equals(v2) {
			return i
		}
	}
	i := len(cg.values)
	cg.values = append(cg.values, v)
	return i
}

func (cg *cgen) rvalue(ast Ast) {
	if isLocal(ast.Value) {
		cg.emit(i.LOAD)
		cg.emitUint(cg.local(ast.Value))
	} else {
		panic("not implemented")
	}
}

func (cg *cgen) lvalue(ast Ast) int {
	if ast.Token == IDENTIFIER && isLocal(ast.Value) {
		return cg.local(ast.Value)
	} else {
		panic("not implemented")
	}
}

func (cg *cgen) load(ref int) {
	cg.emit(i.LOAD)
	cg.emitUint(ref)
}

func (cg *cgen) store(ref int) {
	cg.emit(i.STORE)
	cg.emitUint(ref)
}

func isLocal(s string) bool {
	return 'a' <= s[0] && s[0] <= 'z'
}

// local returns the index for a local variable
func (cg *cgen) local(s string) int {
	for i, s2 := range cg.strings {
		if s == s2 {
			return i
		}
	}
	i := len(cg.strings)
	// TODO intern to avoid ref to source
	cg.strings = append(cg.strings, s)
	return i
}

// ubinop is called for ops that can be unary or binary
func (cg *cgen) ubinop(ast Ast, uop, bop byte) {
	if len(ast.Children) == 1 {
		cg.unop(ast, uop)
	} else {
		cg.binop(ast, bop)
	}
}

func (cg *cgen) unop(ast Ast, op byte) {
	cg.gen(ast.first())
	cg.emit(op)
}

func (cg *cgen) binop(ast Ast, op byte) {
	cg.gen(ast.first())
	cg.gen(ast.second())
	cg.emit(op)
}

func (cg *cgen) emit(b ...byte) {
	// TODO detect pop after side effect free instruction
	cg.code = append(cg.code, b...)
}

func (cg *cgen) emitUint(i int) {
	cg.code = varint.EncodeUint32(uint32(i), cg.code)
}

func (cg *cgen) emitInt(i int) {
	cg.code = varint.EncodeInt32(int32(i), cg.code)
}
