package compile

import (
	i "github.com/apmckinlay/gsuneido/interp"
	"github.com/apmckinlay/gsuneido/interp/globals"
	"github.com/apmckinlay/gsuneido/util/varint"
	"github.com/apmckinlay/gsuneido/util/verify"
	"github.com/apmckinlay/gsuneido/value"
)

// codegen compiles a Function from an Ast
func codegen(ast Ast) *value.SuFunc {
	//fmt.Println("codegen", ast.String())
	cg := cgen{}
	cg.function(ast)
	return &value.SuFunc{
		Code:    cg.code,
		Values:  cg.values,
		Strings: cg.names,
		Nlocals: len(cg.names),
	}
}

type cgen struct {
	code   []byte
	values []value.Value
	names  []string
}

func (cg *cgen) function(ast Ast) {
	verify.That(ast.Keyword == FUNCTION)
	// TODO params
	stmts := ast.second().Children
	for si, stmt := range stmts {
		cg.statement(stmt, si == len(stmts)-1)
	}
}

func (cg *cgen) statement(ast Ast, lastStmt bool) {
	switch ast.KeyTok() {
	case RETURN:
		if len(ast.Children) == 1 {
			cg.expr(ast.first())
		}
		if !lastStmt {
			cg.emit(i.RETURN)
		}
	case EXPRESSION:
		cg.expr(ast.first())
		if !lastStmt {
			cg.emit(i.POP)
		}
	default:
		panic("bad statement: " + ast.String())
	}
}

func (cg *cgen) expr(ast Ast) {
	switch ast.KeyTok() {
	case VALUE:
		cg.emitValue(ast.value)
	case NOT:
		cg.unary(ast, i.NOT)
	case ADD:
		if len(ast.Children) == 1 {
			cg.unary(ast, i.UPLUS)
		} else {
			cg.nary(ast, i.ADD)
		}
	case SUB:
		cg.unary(ast, i.UMINUS) // binary sub handled by add
	case MUL: // also handles div
		cg.nary(ast, i.MUL)
	case IS, ISNT, MATCH, MATCHNOT, LT, LTE, GT, GTE,
		CAT, MOD, LSHIFT, RSHIFT, BITOR, BITAND, BITXOR:
		cg.nary(ast, i.IS+byte(ast.KeyTok()-IS))
	case BITNOT:
		cg.unary(ast, i.BITNOT)
	case IDENTIFIER:
		cg.identifier(ast)
	case EQ:
		ref := cg.lvalue(ast.first())
		cg.expr(ast.second())
		cg.store(ref)
	case ADDEQ, SUBEQ, CATEQ, MULEQ, DIVEQ, MODEQ,
		LSHIFTEQ, RSHIFTEQ, BITOREQ, BITANDEQ, BITXOREQ:
		ref := cg.lvalue(ast.first())
		cg.dupLvalue(ref)
		cg.load(ref)
		cg.expr(ast.second())
		cg.emit(i.ADD + byte(ast.Token-ADDEQ))
		cg.store(ref)
	case INC, DEC:
		ref := cg.lvalue(ast.first())
		cg.dupLvalue(ref)
		cg.load(ref)
		cg.emit(i.ONE)
		cg.emit(i.ADD + byte(ast.Token-INC))
		cg.store(ref)
	case POSTINC, POSTDEC:
		ref := cg.lvalue(ast.first())
		cg.dupLvalue(ref)
		cg.load(ref)
		cg.dupUnderLvalue(ref)
		cg.emit(i.ONE)
		cg.emit(i.ADD + byte(ast.Token-POSTINC))
		cg.store(ref)
		cg.emit(i.POP)
	case DOT: // a.b
		cg.expr(ast.first())
		cg.emitValue(value.SuStr(ast.second().Text))
		cg.emit(i.GET)
	case L_BRACKET: // a[b]
		cg.expr(ast.first())
		cg.expr(ast.second())
		cg.emit(i.GET)
	default:
		panic("bad expression: " + ast.String())
	}
}

func (cg *cgen) emitValue(val value.Value) {
	if val == value.True {
		cg.emit(i.TRUE)
	} else if val == value.False {
		cg.emit(i.FALSE)
	} else if val == value.SuInt(0) {
		cg.emit(i.ZERO)
	} else if val == value.SuInt(1) {
		cg.emit(i.ONE)
	} else if val == value.SuStr("") {
		cg.emit(i.EMPTYSTR)
	} else if si, ok := val.(value.SuInt); ok {
		cg.emit(i.INT)
		cg.emitInt(int(si))
	} else {
		cg.emit(i.VALUE)
		vi := cg.value(val)
		cg.emitUint(vi)
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

func (cg *cgen) identifier(ast Ast) {
	if isLocal(ast.Text) {
		if ast.Text[0] == '_' {
			cg.emit(i.DYLOAD)
		} else {
			cg.emit(i.LOAD)
		}
		cg.emitUint(cg.name(ast.Text))
	} else {
		cg.emit(i.GLOBAL)
		cg.emitUint(globals.NameNum(ast.Text))
	}
}

const MEM_REF = -1

func (cg *cgen) lvalue(ast Ast) int {
	if ast.Token == IDENTIFIER && isLocal(ast.Text) {
		return cg.name(ast.Text)
	} else if ast.Token == DOT {
		cg.expr(ast.first())
		cg.emitValue(value.SuStr(ast.second().Text))
		return MEM_REF
	} else if ast.Token == L_BRACKET {
		cg.expr(ast.first())
		cg.expr(ast.second())
		return MEM_REF
	} else {
		panic("invalid lvalue: " + ast.String())
	}
}

func (cg *cgen) load(ref int) {
	if ref == MEM_REF {
		cg.emit(i.GET)
	} else {
		if cg.names[ref][0] == '_' {
			cg.emit(i.DYLOAD)
		} else {
			cg.emit(i.LOAD)
		}
		cg.emitUint(ref)
	}
}

func (cg *cgen) store(ref int) {
	if ref == MEM_REF {
		cg.emit(i.PUT)
	} else {
		cg.emit(i.STORE)
		cg.emitUint(ref)
	}
}

func (cg *cgen) dupLvalue(ref int) {
	if ref == MEM_REF {
		cg.emit(i.DUP2)
	}
}

func (cg *cgen) dupUnderLvalue(ref int) {
	if ref == MEM_REF {
		cg.emit(i.DUPX2)
	} else {
		cg.emit(i.DUP)
	}
}

// includes dynamic
func isLocal(s string) bool {
	return ('a' <= s[0] && s[0] <= 'z') || s[0] == '_'
}

// name returns the index for a name variable
func (cg *cgen) name(s string) int {
	for i, s2 := range cg.names {
		if s == s2 {
			return i
		}
	}
	i := len(cg.names)
	// TODO intern to avoid ref to source
	cg.names = append(cg.names, s)
	return i
}

func (cg *cgen) unary(ast Ast, op byte) {
	verify.That(len(ast.Children) == 1)
	cg.expr(ast.first())
	cg.emit(op)
}

func (cg *cgen) nary(ast Ast, op byte) {
	cg.expr(ast.first())
	for _, a := range ast.Children[1:] {
		if op == i.ADD && a.Token == SUB && len(a.Children) == 1 {
			cg.expr(a.first())
			cg.emit(i.SUB)
		} else if op == i.MUL && a.Token == DIV && len(a.Children) == 1 {
			cg.expr(a.first())
			cg.emit(i.DIV)
		} else {
			cg.expr(a)
			cg.emit(op)
		}
	}
}

// emit is used to append an op code
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
