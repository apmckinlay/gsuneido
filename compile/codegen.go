package compile

/* TODO
- blocks
- try catch
- function calls
*/

// See also Disasm

import (
	"math"

	. "github.com/apmckinlay/gsuneido/lexer"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/op"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/verify"
)

// zeroFlags is shared/reused for all zero flags
var zeroFlags [MaxArgs]Flag

// codegen compiles an Ast to an SuFunc
func codegen(ast *Ast) *SuFunc {
	//fmt.Println("codegen", ast.String())
	cg := cgen{}
	var tmp [MaxArgs]Flag
	cg.flags = tmp[:0]
	cg.function(ast)
	if allZero(cg.flags) {
		cg.flags = zeroFlags[:len(cg.flags)]
	} else {
		// shrink flags to used size
		bigflags := cg.flags
		cg.flags = make([]Flag, len(bigflags))
		copy(cg.flags, bigflags)
	}
	return &SuFunc{
		Code:    cg.code,
		Nlocals: len(cg.names),
		ParamSpec: ParamSpec{
			Values:    cg.values,
			Strings:   cg.names,
			Nparams:   cg.nparams,
			Ndefaults: cg.ndefaults,
			Flags:     cg.flags},
	}
}

func allZero(flags []Flag) bool {
	for _, f := range flags {
		if f != 0 {
			return false
		}
	}
	return true
}

//TODO embed SuFunc
type cgen struct {
	nparams   int
	code      []byte
	values    []Value
	names     []string
	flags     []Flag
	ndefaults int
}

var tok2op = [Ntokens]byte{
	AND:      op.AND,
	OR:       op.OR,
	INC:      op.ADD,
	DEC:      op.SUB,
	ADDEQ:    op.ADD,
	SUBEQ:    op.SUB,
	CATEQ:    op.CAT,
	MULEQ:    op.MUL,
	DIVEQ:    op.DIV,
	MODEQ:    op.MOD,
	LSHIFTEQ: op.LSHIFT,
	RSHIFTEQ: op.RSHIFT,
	BITOREQ:  op.BITOR,
	BITANDEQ: op.BITAND,
	BITXOREQ: op.BITXOR,
	IS:       op.IS,
	ISNT:     op.ISNT,
	MATCH:    op.MATCH,
	MATCHNOT: op.MATCHNOT,
	LT:       op.LT,
	LTE:      op.LTE,
	GT:       op.GT,
	GTE:      op.GTE,
	CAT:      op.CAT,
	MOD:      op.MOD,
	LSHIFT:   op.LSHIFT,
	RSHIFT:   op.RSHIFT,
	BITOR:    op.BITOR,
	BITAND:   op.BITAND,
	BITXOR:   op.BITXOR,
}

func (cg *cgen) function(ast *Ast) {
	verify.That(ast.Keyword == FUNCTION)
	cg.params(ast.first())
	stmts := ast.second().Children
	for si, stmt := range stmts {
		cg.statement(stmt, nil, si == len(stmts)-1)
	}
}

func (cg *cgen) params(ast *Ast) {
	verify.That(ast.Text == "params")
	cg.nparams = len(ast.Children)
	for _, p := range ast.Children {
		name, flags := cg.param(p.Text)
		cg.names = append(cg.names, name) // no duplicate reuse
		cg.flags = append(cg.flags, flags)
		if p.value != nil {
			cg.ndefaults++
			cg.values = append(cg.values, p.value) // no duplicate reuse
		}
	}
}

func (cg *cgen) param(p string) (string, Flag) {
	if p[0] == '@' {
		return p[1:], AtParam
	}
	var flag Flag
	if p[0] == '.' {
		flag = DotParam
		p = p[1:]
	}
	if p[0] == '_' {
		flag |= DynParam
		p = p[1:]
	}
	if flag&DotParam == DotParam && str.Capitalized(p) {
		flag |= PubParam
		p = str.UnCapitalize(p)
	}
	return p, flag
}

func (cg *cgen) statement(ast *Ast, labels *Labels, lastStmt bool) {
	switch ast.KeyTok() {
	case L_CURLY:
		for _, a := range ast.Children {
			cg.statement(a, labels, lastStmt)
		}
	case RETURN:
		if len(ast.Children) == 1 {
			cg.expr(ast.first())
		}
		if !lastStmt {
			cg.emit(op.RETURN)
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
		if ast.Text == "in" {
			panic("not implemented") // TODO
		} else {
			cg.forStmt(ast)
		}
	case THROW:
		cg.expr(ast.first())
		cg.emit(op.THROW)
	case BREAK:
		if labels == nil {
			panic("break can only be used within a loop")
		}
		labels.brk = cg.emitJump(op.JUMP, labels.brk)
	case CONTINUE:
		if labels == nil {
			panic("continue can only be used within a loop")
		}
		cg.emitBwdJump(op.JUMP, labels.cont)
	default: // expression
		cg.expr(ast)
		if !lastStmt {
			cg.emit(op.POP)
		}
	}
}

func (cg *cgen) ifStmt(ast *Ast, labels *Labels) {
	cg.expr(ast.first())
	f := cg.emitJump(op.FJUMP, -1)
	cg.statement(ast.second(), labels, false)
	if len(ast.Children) == 3 {
		end := cg.emitJump(op.JUMP, -1)
		cg.placeLabel(f)
		cg.statement(ast.third(), labels, false)
		cg.placeLabel(end)
	} else {
		cg.placeLabel(f)
	}
}

func (cg *cgen) switchStmt(ast *Ast, labels *Labels) {
	cg.expr(ast.first())
	end := -1
	for _, c := range ast.second().Children {
		caseBody, afterCase := -1, -1
		values := c.first().Children
		for v, val := range values {
			cg.expr(val)
			if v < len(values)-1 {
				caseBody = cg.emitJump(op.EQJUMP, -1)
			} else {
				afterCase = cg.emitJump(op.NEJUMP, -1)
			}
		}
		cg.placeLabel(caseBody)
		cg.statement(c.second(), labels, false)
		end = cg.emitJump(op.JUMP, end)
		cg.placeLabel(afterCase)
	}
	cg.emit(op.POP)
	if len(ast.Children) == 3 {
		cg.statement(ast.third(), labels, false)
	} else {
		cg.emitValue(SuStr("unhandled switch value"))
		cg.emit(op.THROW)
	}
	cg.placeLabel(end)
}

func (cg *cgen) foreverStmt(ast *Ast) {
	labels := cg.newLabels()
	cg.statement(ast.first(), labels, false)
	cg.emitJump(op.JUMP, labels.cont-len(cg.code)-3)
	cg.placeLabel(labels.brk)
}

func (cg *cgen) whileStmt(ast *Ast) {
	labels := cg.newLabels()
	cond := cg.emitJump(op.JUMP, -1)
	loop := cg.label()
	cg.statement(ast.second(), labels, false)
	cg.placeLabel(cond)
	cg.expr(ast.first())
	cg.emitBwdJump(op.TJUMP, loop)
	cg.placeLabel(labels.brk)
}

func (cg *cgen) dowhileStmt(ast *Ast) {
	labels := cg.newLabels()
	cg.statement(ast.first(), labels, false)
	cg.expr(ast.second())
	cg.emitBwdJump(op.TJUMP, labels.cont)
	cg.placeLabel(labels.brk)
}

func (cg *cgen) forStmt(ast *Ast) {
	cg.exprList(ast.first().Children) // init
	labels := cg.newLabels()
	cond := cg.emitJump(op.JUMP, -1)
	loop := cg.label()
	cg.statement(ast.fourth(), labels, false) // body
	cg.exprList(ast.third().Children)         // increment
	cg.placeLabel(cond)
	cg.expr(ast.second()) // condition
	cg.emitBwdJump(op.TJUMP, loop)
	cg.placeLabel(labels.brk)
}

func (cg *cgen) exprList(list []*Ast) {
	for _, expr := range list {
		cg.expr(expr)
		cg.emit(op.POP)
	}
}

// expressions -----------------------------------------------------------------

func (cg *cgen) expr(ast *Ast) {
	switch ast.KeyTok() {
	case NOT:
		cg.unary(ast, op.NOT)
	case ADD:
		if len(ast.Children) == 1 {
			cg.unary(ast, op.UPLUS)
		} else {
			cg.nary(ast, op.ADD)
		}
	case SUB:
		cg.unary(ast, op.UMINUS) // binary sub handled by add
	case MUL: // also handles div
		cg.nary(ast, op.MUL)
	case IS, ISNT, MATCH, MATCHNOT, LT, LTE, GT, GTE,
		CAT, MOD, LSHIFT, RSHIFT, BITOR, BITAND, BITXOR:
		cg.nary(ast, tok2op[ast.KeyTok()])
	case BITNOT:
		cg.unary(ast, op.BITNOT)
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
		cg.emit(tok2op[ast.Token])
		cg.store(ref)
	case INC, DEC:
		ref := cg.lvalue(ast.first())
		cg.dupLvalue(ref)
		cg.load(ref)
		if ast.Text == "post" {
			cg.dupUnderLvalue(ref)
			cg.emit(op.ONE)
			cg.emit(tok2op[ast.Token])
			cg.store(ref)
			cg.emit(op.POP)
		} else {
			cg.emit(op.ONE)
			cg.emit(tok2op[ast.Token])
			cg.store(ref)
		}
	case DOT: // a.b
		cg.expr(ast.first())
		cg.emitValue(SuStr(ast.second().Text))
		cg.emit(op.GET)
	case L_BRACKET: // a[b]
		cg.expr(ast.first())
		sub := ast.second()
		switch sub.Token {
		case RANGETO:
			cg.expr(sub.first())
			cg.expr(sub.second())
			cg.emit(op.RANGETO)
		case RANGELEN:
			cg.expr(sub.first())
			cg.expr(sub.second())
			cg.emit(op.RANGELEN)
		default:
			cg.expr(ast.second())
			cg.emit(op.GET)
		}
	case AND, OR:
		cg.andorExpr(ast)
	case Q_MARK:
		cg.qcExpr(ast)
	case IN:
		cg.inExpr(ast)
	case THIS:
		cg.emit(op.THIS)
	case FUNCTION:
		fn := codegen(ast)
		cg.emitValue(fn)
	default:
		if ast.Item == call {
			cg.call(ast)
		} else if ast.value != nil {
			cg.emitValue(ast.value)
		} else {
			panic("unhandled expression: " + ast.String())
		}
	}
}

func (cg *cgen) andorExpr(ast *Ast) {
	label := -1
	cg.expr(ast.first())
	for _, a := range ast.Children[1:] {
		label = cg.emitJump(tok2op[ast.Keyword], label)
		cg.expr(a)
	}
	cg.emit(op.BOOL)
	cg.placeLabel(label)
}

func (cg *cgen) qcExpr(ast *Ast) {
	f, end := -1, -1
	cg.expr(ast.first())
	f = cg.emitJump(op.Q_MARK, f)
	cg.expr(ast.second())
	end = cg.emitJump(op.JUMP, end)
	cg.placeLabel(f)
	cg.expr(ast.third())
	cg.placeLabel(end)
}

func (cg *cgen) inExpr(ast *Ast) {
	end := -1
	cg.expr(ast.first())
	for j, a := range ast.Children[1:] {
		cg.expr(a)
		if j < len(ast.Children)-2 {
			end = cg.emitJump(op.IN, end)
		} else {
			cg.emit(op.IS)
		}
	}
	cg.placeLabel(end)
}

func (cg *cgen) emitValue(val Value) {
	if val == True {
		cg.emit(op.TRUE)
	} else if val == False {
		cg.emit(op.FALSE)
	} else if val == SuInt(0) {
		cg.emit(op.ZERO)
	} else if val == SuInt(1) {
		cg.emit(op.ONE)
	} else if val == SuStr("") {
		cg.emit(op.EMPTYSTR)
	} else if i, ok := SmiToInt(val); ok {
		cg.emitInt16(op.INT, i)
	} else {
		cg.emitUint16(op.VALUE, cg.value(val))
	}
}

// value returns an index for the value
// reusing if duplicate, adding otherwise
func (cg *cgen) value(v Value) int {
	for i, v2 := range cg.values {
		if v.Equal(v2) {
			return i
		}
	}
	i := len(cg.values)
	cg.values = append(cg.values, v)
	return i
}

func (cg *cgen) identifier(ast *Ast) {
	if isLocal(ast.Text) {
		i := cg.name(ast.Text)
		if ast.Text[0] == '_' {
			cg.emitUint8(op.DYLOAD, i)
		} else {
			cg.emitUint8(op.LOAD, i)
		}
	} else {
		cg.emitUint16(op.GLOBAL, GlobalNum(ast.Text))
	}
}

const memRef = -1

func (cg *cgen) lvalue(ast *Ast) int {
	if ast.Token == IDENTIFIER && isLocal(ast.Text) {
		return cg.name(ast.Text)
	} else if ast.Token == DOT {
		cg.expr(ast.first())
		cg.emitValue(SuStr(ast.second().Text))
		return memRef
	} else if ast.Token == L_BRACKET {
		cg.expr(ast.first())
		cg.expr(ast.second())
		return memRef
	} else {
		panic("invalid lvalue: " + ast.String())
	}
}

func (cg *cgen) load(ref int) {
	if ref == memRef {
		cg.emit(op.GET)
	} else {
		if cg.names[ref][0] == '_' {
			cg.emitUint8(op.DYLOAD, ref)
		} else {
			cg.emitUint8(op.LOAD, ref)
		}
	}
}

func (cg *cgen) store(ref int) {
	if ref == memRef {
		cg.emit(op.PUT)
	} else {
		cg.emitUint8(op.STORE, ref)
	}
}

func (cg *cgen) dupLvalue(ref int) {
	if ref == memRef {
		cg.emit(op.DUP2)
	}
}

func (cg *cgen) dupUnderLvalue(ref int) {
	if ref == memRef {
		cg.emit(op.DUPX2)
	} else {
		cg.emit(op.DUP)
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
	cg.names = append(cg.names, s)
	return i
}

func (cg *cgen) unary(ast *Ast, op byte) {
	verify.That(len(ast.Children) == 1)
	cg.expr(ast.first())
	cg.emit(op)
}

func (cg *cgen) nary(ast *Ast, o byte) {
	cg.expr(ast.first())
	for _, a := range ast.Children[1:] {
		if o == op.ADD && a.Token == SUB && len(a.Children) == 1 {
			cg.expr(a.first())
			cg.emit(op.SUB)
		} else if o == op.MUL && a.Token == DIV && len(a.Children) == 1 {
			cg.expr(a.first())
			cg.emit(op.DIV)
		} else {
			cg.expr(a)
			cg.emit(o)
		}
	}
}

func (cg *cgen) call(ast *Ast) {
	fn := ast.first()
	method := fn.Token == DOT || fn.Token == L_BRACKET
	if method {
		cg.expr(fn.first()) // self
	}

	argspec := cg.args(ast.second())
	if method {
		if fn.Token == DOT {
			cg.emitValue(SuStr(fn.second().Text)) // method name
		} else { // L_BRACKET
			cg.expr(fn.second())
		}
		cg.emit(op.CALLMETH)
	} else {
		cg.expr(ast.first()) // function
		cg.emit(op.CALLFUNC)
	}
	cg.emit(argspec.Unnamed)
	named := len(argspec.Spec)
	verify.That(named <= math.MaxUint8)
	cg.emit(byte(named))
	cg.emit(argspec.Spec...)
}

func (cg *cgen) args(ast *Ast) ArgSpec {
	if ast.Item == atArg {
		cg.expr(ast.Children[0])
		return ArgSpec{Unnamed: EACH}
	} else if ast.Item == at1Arg {
		cg.expr(ast.Children[0])
		return ArgSpec{Unnamed: EACH1}
	}
	verify.That(ast.Item == argList)
	var spec []byte
	for _, arg := range ast.Children {
		if arg.Item != noKeyword {
			i := cg.name(arg.Item.Text)
			verify.That(i <= math.MaxUint8)
			spec = append(spec, byte(i))
		}
		cg.expr(arg.first())
	}
	verify.That(len(ast.Children) < int(EACH))
	return ArgSpec{Unnamed: byte(len(ast.Children) - len(spec)), Spec: spec}
}

// helpers ---------------------------------------------------------------------

// emit is used to append an op code
func (cg *cgen) emit(b ...byte) {
	cg.code = append(cg.code, b...)
}

func (cg *cgen) emitUint8(op byte, i int) {
	verify.That(0 <= i && i < math.MaxUint8)
	cg.emit(op, byte(i))
}

func (cg *cgen) emitInt16(op byte, i int) {
	verify.That(math.MinInt16 <= i && i < math.MaxInt16)
	cg.emit(op, byte(i>>8), byte(i))
}

func (cg *cgen) emitUint16(op byte, i int) {
	verify.That(0 <= i && i < math.MaxUint16)
	cg.emit(op, byte(i>>8), byte(i))
}

func (cg *cgen) emitJump(op byte, label int) int {
	adr := len(cg.code)
	verify.That(math.MinInt16 <= label && label <= math.MaxInt16)
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
		verify.That(math.MinInt16 <= adr && adr <= math.MaxInt16)
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
