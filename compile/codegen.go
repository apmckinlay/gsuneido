package compile

import (
	i "github.com/apmckinlay/gsuneido/interp"
	"github.com/apmckinlay/gsuneido/interp/globals"
	. "github.com/apmckinlay/gsuneido/lexer"
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
		cg.statement(stmt, nil, si == len(stmts)-1)
	}
}

func (cg *cgen) statement(ast Ast, labels *Labels, lastStmt bool) {
	switch ast.KeyTok() {
	case STATEMENTS:
		for _, a := range ast.Children {
			cg.statement(a, labels, lastStmt)
		}
	case NIL:
		// no code
	case RETURN:
		if len(ast.Children) == 1 {
			cg.expr(ast.first())
		}
		if !lastStmt {
			cg.emit(i.RETURN)
		}
	case IF:
		cg.ifStmt(ast, labels)
	case SWITCH:
		cg.switchStmt(ast, labels)
	case FOREVER:
		cg.foreverStmt(ast)
	case WHILE:
		cg.whileStmt(ast)
	case DO:
		cg.dowhileStmt(ast)
	case FOR:
		cg.forStmt(ast)
	case THROW:
		cg.expr(ast.first())
		cg.emit(i.THROW)
	case BREAK:
		if labels == nil {
			panic("break can only be used within a loop")
		}
		labels.brk = cg.emitJump(i.JUMP, labels.brk)
	case CONTINUE:
		if labels == nil {
			panic("continue can only be used within a loop")
		}
		cg.emitBwdJump(i.JUMP, labels.cont)
	default: // expression
		cg.expr(ast)
		if !lastStmt {
			cg.emit(i.POP)
		}
	}
}

func (cg *cgen) ifStmt(ast Ast, labels *Labels) {
	cg.expr(ast.first())
	f := cg.emitJump(i.FJUMP, -1)
	cg.statement(ast.second(), labels, false)
	if len(ast.Children) == 3 {
		end := cg.emitJump(i.JUMP, -1)
		cg.placeLabel(f)
		cg.statement(ast.third(), labels, false)
		cg.placeLabel(end)
	} else {
		cg.placeLabel(f)
	}
}

func (cg *cgen) switchStmt(ast Ast, labels *Labels) {
	cg.expr(ast.first())
	end := -1
	for _, c := range ast.second().Children {
		caseBody, afterCase := -1, -1
		values := c.first().Children
		for v, val := range values {
			cg.expr(val)
			if v < len(values)-1 {
				caseBody = cg.emitJump(i.EQJUMP, -1)
			} else {
				afterCase = cg.emitJump(i.NEJUMP, -1)
			}
		}
		cg.placeLabel(caseBody)
		cg.statement(c.second(), labels, false)
		end = cg.emitJump(i.JUMP, end)
		cg.placeLabel(afterCase)
	}
	cg.emit(i.POP)
	if len(ast.Children) == 3 {
		cg.statement(ast.third(), labels, false)
	} else {
		cg.emitValue(value.SuStr("unhandled switch value"))
		cg.emit(i.THROW)
	}
	cg.placeLabel(end)
}

func (cg *cgen) foreverStmt(ast Ast) {
	labels := cg.newLabels()
	cg.statement(ast.first(), labels, false)
	cg.emitJump(i.JUMP, labels.cont-len(cg.code)-3)
	cg.placeLabel(labels.brk)
}

func (cg *cgen) whileStmt(ast Ast) {
	labels := cg.newLabels()
	cond := cg.emitJump(i.JUMP, -1)
	loop := cg.label()
	cg.statement(ast.second(), labels, false)
	cg.placeLabel(cond)
	cg.expr(ast.first())
	cg.emitBwdJump(i.TJUMP, loop)
	cg.placeLabel(labels.brk)
}

func (cg *cgen) dowhileStmt(ast Ast) {
	labels := cg.newLabels()
	cg.statement(ast.first(), labels, false)
	cg.expr(ast.second())
	cg.emitBwdJump(i.TJUMP, labels.cont)
	cg.placeLabel(labels.brk)
}

func (cg *cgen) forStmt(ast Ast) {
	cg.exprList(ast.first().Children) // init
	labels := cg.newLabels()
	cond := cg.emitJump(i.JUMP, -1)
	loop := cg.label()
	cg.statement(ast.fourth(), labels, false) // body
	cg.exprList(ast.third().Children)         // increment
	cg.placeLabel(cond)
	cg.expr(ast.second()) // condition
	cg.emitBwdJump(i.TJUMP, loop)
	cg.placeLabel(labels.brk)
}

func (cg *cgen) exprList(list []Ast) {
	for _, expr := range list {
		cg.expr(expr)
		cg.emit(i.POP)
	}
}

// expressions -----------------------------------------------------------------

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
	case AND, OR:
		cg.andorExpr(ast)
	case Q_MARK:
		cg.qcExpr(ast)
	case IN:
		cg.inExpr(ast)
	default:
		panic("bad expression: " + ast.String())
	}
}

func (cg *cgen) andorExpr(ast Ast) {
	label := -1
	cg.expr(ast.first())
	for _, a := range ast.Children[1:] {
		label = cg.emitJump(tok2op[ast.Keyword], label)
		cg.expr(a)
	}
	cg.emit(i.BOOL)
	cg.placeLabel(label)
}

func (cg *cgen) qcExpr(ast Ast) {
	f, end := -1, -1
	cg.expr(ast.first())
	f = cg.emitJump(i.Q_MARK, f)
	cg.expr(ast.second())
	end = cg.emitJump(i.JUMP, end)
	cg.placeLabel(f)
	cg.expr(ast.third())
	cg.placeLabel(end)
}

func (cg *cgen) inExpr(ast Ast) {
	end := -1
	cg.expr(ast.first())
	for j, a := range ast.Children[1:] {
		cg.expr(a)
		if j < len(ast.Children)-2 {
			end = cg.emitJump(i.IN, end)
		} else {
			cg.emit(i.IS)
		}
	}
	cg.placeLabel(end)
}

var tok2op = map[Token]byte{
	AND: i.AND,
	OR:  i.OR,
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

// helpers ---------------------------------------------------------------------

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

func (cg *cgen) emitJump(op byte, label int) int {
	adr := len(cg.code)
	cg.emit(op, byte(label>>8), byte(label))
	return adr
}

func (cg *cgen) emitBwdJump(op byte, label int) {
	cg.emitJump(op, label-len(cg.code)-3)
}

func (cg *cgen) label() int {
	return len(cg.code)
}

func (cg *cgen) placeLabel(i int) {
	var adr, next int
	for ; i >= 0; i = next {
		next = int(cg.target(i))
		adr = len(cg.code) - (i + 3) // convert to relative offset
		cg.code[i+1] = byte(adr >> 8)
		cg.code[i+2] = byte(adr)
	}
}

func (cg *cgen) target(i int) int16 {
	return int16(uint16(cg.code[i+1])<<8 | uint16(cg.code[i+2]))
}

type Labels struct {
	brk  int // chained forward jump
	cont int // backward jump
}

func (cg *cgen) newLabels() *Labels {
	return &Labels{-1, cg.label()}
}
