package compile

// See also Disasm

import (
	"fmt"
	"math"

	"github.com/apmckinlay/gsuneido/compile/ast"
	. "github.com/apmckinlay/gsuneido/lexer"
	. "github.com/apmckinlay/gsuneido/runtime"
	"github.com/apmckinlay/gsuneido/runtime/op"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/verify"
)

// zeroFlags is shared/reused for all zero flags
var zeroFlags [MaxArgs]Flag

// codegen compiles an Ast to an SuFunc
func codegen(fn *ast.Function) *SuFunc {
	//fmt.Println("codegen", ast.String())
	cg := cgen{}
	cg.function(fn)
	if allZero(cg.Flags) {
		cg.Flags = zeroFlags[:len(cg.Flags)]
	}
	return &SuFunc{
		Code:      cg.code,
		Nlocals:   len(cg.Names),
		ParamSpec: cg.ParamSpec,
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

type cgen struct {
	ParamSpec
	code []byte
}

// binary and nary ast node token to operation
var tok2op = [Ntokens]byte{
	ADD:      op.ADD,
	SUB:      op.SUB,
	CAT:      op.CAT,
	MUL:      op.MUL,
	DIV:      op.DIV,
	MOD:      op.MOD,
	LSHIFT:   op.LSHIFT,
	RSHIFT:   op.RSHIFT,
	BITOR:    op.BITOR,
	BITAND:   op.BITAND,
	BITXOR:   op.BITXOR,
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
	AND:      op.AND,
	OR:       op.OR,
}

func (cg *cgen) function(fn *ast.Function) {
	cg.params(fn.Params)
	stmts := fn.Body
	for si, stmt := range stmts {
		cg.statement(stmt, nil, si == len(stmts)-1)
	}
}

func (cg *cgen) params(params []ast.Param) {
	cg.Nparams = len(params)
	for _, p := range params {
		name, flags := cg.param(p.Name)
		cg.Names = append(cg.Names, name) // no duplicate reuse
		cg.Flags = append(cg.Flags, flags)
		if p.DefVal != nil {
			cg.Ndefaults++
			cg.Values = append(cg.Values, p.DefVal) // no duplicate reuse
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

func (cg *cgen) statement(node ast.Node, labels *Labels, lastStmt bool) {
	switch node := node.(type) {
	case *ast.Compound:
		for _, stmt := range node.Body {
			cg.statement(stmt, labels, lastStmt)
		}
	case *ast.Return:
		if node.E != nil {
			cg.expr(node.E)
		}
		if !lastStmt {
			cg.emit(op.RETURN)
		}
	case *ast.If:
		cg.ifStmt(node, labels)
	case *ast.Switch:
		cg.switchStmt(node, labels)
	case *ast.Forever:
		cg.foreverStmt(node)
	case *ast.While:
		cg.whileStmt(node)
	case *ast.DoWhile:
		cg.dowhileStmt(node)
	case *ast.For:
		cg.forStmt(node)
	case *ast.Throw:
		cg.expr(node.E)
		cg.emit(op.THROW)
	case *ast.TryCatch:
		cg.emit(op.TRY)
		//TODO
	case *ast.Break:
		if labels == nil {
			panic("break can only be used within a loop")
		}
		labels.brk = cg.emitJump(op.JUMP, labels.brk)
	case *ast.Continue:
		if labels == nil {
			panic("continue can only be used within a loop")
		}
		cg.emitBwdJump(op.JUMP, labels.cont)
	case *ast.Expression:
		cg.expr(node.E)
		if !lastStmt {
			cg.emit(op.POP)
		}
	default:
		panic("unexpected statement type")
	}
}

func (cg *cgen) statements(stmts []ast.Statement, labels *Labels, lastStmt bool) {
	for _, stmt := range stmts {
		cg.statement(stmt, labels, lastStmt)
	}
}

func (cg *cgen) ifStmt(node *ast.If, labels *Labels) {
	cg.expr(node.Cond)
	f := cg.emitJump(op.FJUMP, -1)
	cg.statement(node.Then, labels, false)
	if node.Else != nil {
		end := cg.emitJump(op.JUMP, -1)
		cg.placeLabel(f)
		cg.statement(node.Else, labels, false)
		cg.placeLabel(end)
	} else {
		cg.placeLabel(f)
	}
}

func (cg *cgen) switchStmt(node *ast.Switch, labels *Labels) {
	cg.expr(node.E)
	end := -1
	for _, c := range node.Cases {
		caseBody, afterCase := -1, -1
		for v, e := range c.Exprs {
			cg.expr(e)
			if v < len(c.Exprs)-1 {
				caseBody = cg.emitJump(op.EQJUMP, -1)
			} else {
				afterCase = cg.emitJump(op.NEJUMP, -1)
			}
		}
		cg.placeLabel(caseBody)
		cg.statements(c.Body, labels, false)
		end = cg.emitJump(op.JUMP, end)
		cg.placeLabel(afterCase)
	}
	cg.emit(op.POP)
	if node.Default != nil {
		cg.statements(node.Default, labels, false)
	} else {
		cg.emitValue(SuStr("unhandled switch value"))
		cg.emit(op.THROW)
	}
	cg.placeLabel(end)
}

func (cg *cgen) foreverStmt(node *ast.Forever) {
	labels := cg.newLabels()
	cg.statement(node.Body, labels, false)
	cg.emitJump(op.JUMP, labels.cont-len(cg.code)-3)
	cg.placeLabel(labels.brk)
}

func (cg *cgen) whileStmt(node *ast.While) {
	labels := cg.newLabels()
	cond := cg.emitJump(op.JUMP, -1)
	loop := cg.label()
	cg.statement(node.Body, labels, false)
	cg.placeLabel(cond)
	cg.expr(node.Cond)
	cg.emitBwdJump(op.TJUMP, loop)
	cg.placeLabel(labels.brk)
}

func (cg *cgen) dowhileStmt(node *ast.DoWhile) {
	labels := cg.newLabels()
	cg.statement(node.Body, labels, false)
	cg.expr(node.Cond)
	cg.emitBwdJump(op.TJUMP, labels.cont)
	cg.placeLabel(labels.brk)
}

func (cg *cgen) forStmt(node *ast.For) {
	cg.exprList(node.Init)
	labels := cg.newLabels()
	cond := cg.emitJump(op.JUMP, -1)
	loop := cg.label()
	cg.statement(node.Body, labels, false)
	cg.exprList(node.Inc) // increment
	cg.placeLabel(cond)
	cg.expr(node.Cond)
	cg.emitBwdJump(op.TJUMP, loop)
	cg.placeLabel(labels.brk)
}

func (cg *cgen) exprList(list []ast.Expr) {
	for _, expr := range list {
		cg.expr(expr)
		cg.emit(op.POP)
	}
}

// expressions -----------------------------------------------------------------

func (cg *cgen) expr(node ast.Expr) {
	switch node := node.(type) {
	case *ast.Constant:
		cg.emitValue(node.Val)
	case *ast.Ident:
		cg.identifier(node)
	case *ast.Unary:
		cg.unary(node)
	case *ast.Binary:
		cg.binary(node)
	case *ast.Nary:
		cg.nary(node)
	case *ast.Trinary:
		cg.qcExpr(node)
	case *ast.Mem:
		cg.expr(node.E)
		cg.expr(node.M)
		cg.emit(op.GET)
	case *ast.RangeTo:
		cg.expr(node.E)
		cg.exprOr(node.From, op.ZERO)
		cg.exprOr(node.To, op.MAXINT)
		cg.emit(op.RANGETO)
	case *ast.RangeLen:
		cg.expr(node.E)
		cg.exprOr(node.From, op.ZERO)
		cg.exprOr(node.Len, op.MAXINT)
		cg.emit(op.RANGELEN)
	case *ast.In:
		cg.inExpr(node)
	case *ast.Function:
		fn := codegen(node)
		cg.emitValue(fn)
	case *ast.Call:
		cg.call(node)
	case *ast.Block:
		//TODO
	default:
		panic("unhandled expression: " + fmt.Sprint(node))
	}
}

func (cg *cgen) exprOr(expr ast.Expr, op byte) {
	if expr == nil {
		cg.emit(op)
	} else {
		cg.expr(expr)
	}
}

func (cg *cgen) identifier(node *ast.Ident) {
	if node.Name == "this" {
		cg.emit(op.THIS)
	} else if isLocal(node.Name) {
		i := cg.name(node.Name)
		if node.Name[0] == '_' {
			cg.emitUint8(op.DYLOAD, i)
		} else {
			cg.emitUint8(op.LOAD, i)
		}
	} else {
		cg.emitUint16(op.GLOBAL, GlobalNum(node.Name))
	}
}

// includes dynamic
func isLocal(s string) bool {
	return ('a' <= s[0] && s[0] <= 'z') || s[0] == '_'
}

// name returns the index for a name variable
func (cg *cgen) name(s string) int {
	for i, s2 := range cg.Names {
		if s == s2 {
			return i
		}
	}
	i := len(cg.Names)
	cg.Names = append(cg.Names, s)
	return i
}

func (cg *cgen) unary(node *ast.Unary) {
	o := utok2op[node.Tok]
	if INC <= node.Tok && node.Tok <= POSTDEC {
		ref := cg.lvalue(node.E)
		cg.dupLvalue(ref)
		cg.load(ref)
		if node.Tok == POSTINC || node.Tok == POSTDEC {
			cg.dupUnderLvalue(ref)
			cg.emit(op.ONE)
			cg.emit(o)
			cg.store(ref)
			cg.emit(op.POP)
		} else {
			cg.emit(op.ONE)
			cg.emit(o)
			cg.store(ref)
		}
	} else {
		cg.expr(node.E)
		cg.emit(o)
	}
}

// Unary ast expr node token to operation
var utok2op = [Ntokens]byte{
	ADD:     op.UPLUS,
	SUB:     op.UMINUS,
	NOT:     op.NOT,
	BITNOT:  op.BITNOT,
	INC:     op.ADD,
	POSTINC: op.ADD,
	DEC:     op.SUB,
	POSTDEC: op.SUB,
}

func (cg *cgen) binary(node *ast.Binary) {
	switch node.Tok {
	case EQ:
		ref := cg.lvalue(node.Lhs)
		cg.expr(node.Rhs)
		cg.store(ref)
	case ADDEQ, SUBEQ, CATEQ, MULEQ, DIVEQ, MODEQ,
		LSHIFTEQ, RSHIFTEQ, BITOREQ, BITANDEQ, BITXOREQ:
		ref := cg.lvalue(node.Lhs)
		cg.dupLvalue(ref)
		cg.load(ref)
		cg.expr(node.Rhs)
		cg.emit(tok2op[node.Tok])
		cg.store(ref)
	case IS, ISNT, MATCH, MATCHNOT, MOD, LSHIFT, RSHIFT,
		LT, LTE, GT, GTE:
		cg.expr(node.Lhs)
		cg.expr(node.Rhs)
		cg.emit(tok2op[node.Tok])
	default:
		panic("unhandled binary operation " + node.Tok.String())
	}
}

func (cg *cgen) nary(node *ast.Nary) {
	if node.Tok == AND || node.Tok == OR {
		cg.andorExpr(node)
	} else {
		o := tok2op[node.Tok]
		cg.expr(node.Exprs[0])
		for _, e := range node.Exprs[1:] {
			if node.Tok == ADD && isUnary(e, SUB) {
				cg.expr(e.(*ast.Unary).E)
				cg.emit(op.SUB)
			} else if node.Tok == MUL && isUnary(e, DIV) {
				cg.expr(e.(*ast.Unary).E)
				cg.emit(op.DIV)
			} else {
				cg.expr(e)
				cg.emit(o)
			}
		}
	}
}

func (cg *cgen) andorExpr(node *ast.Nary) {
	label := -1
	cg.expr(node.Exprs[0])
	for _, e := range node.Exprs[1:] {
		label = cg.emitJump(tok2op[node.Tok], label)
		cg.expr(e)
	}
	cg.emit(op.BOOL)
	cg.placeLabel(label)
}

func isUnary(e ast.Expr, tok Token) bool {
	u, ok := e.(*ast.Unary)
	return ok && u.Tok == tok
}

func (cg *cgen) qcExpr(node *ast.Trinary) {
	f, end := -1, -1
	cg.expr(node.Cond)
	f = cg.emitJump(op.Q_MARK, f)
	cg.expr(node.T)
	end = cg.emitJump(op.JUMP, end)
	cg.placeLabel(f)
	cg.expr(node.F)
	cg.placeLabel(end)
}

func (cg *cgen) inExpr(node *ast.In) {
	end := -1
	cg.expr(node.E)
	for j, e := range node.Exprs {
		cg.expr(e)
		if j < len(node.Exprs)-1 {
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
	for i, v2 := range cg.Values {
		if v.Equal(v2) {
			return i
		}
	}
	i := len(cg.Values)
	cg.Values = append(cg.Values, v)
	return i
}

const memRef = -1

func (cg *cgen) lvalue(node ast.Expr) int {
	switch node := node.(type) {
	case *ast.Ident:
		if isLocal(node.Name) {
			return cg.name(node.Name)
		}
	case *ast.Mem:
		cg.expr(node.E)
		cg.expr(node.M)
		return memRef
	}
	panic("invalid lvalue: " + fmt.Sprint(node))
}

func (cg *cgen) load(ref int) {
	if ref == memRef {
		cg.emit(op.GET)
	} else {
		if cg.Names[ref][0] == '_' {
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

func (cg *cgen) call(node *ast.Call) {
	fn := node.Fn
	mem, method := fn.(*ast.Mem)
	if method {
		cg.expr(mem.E)
	}
	argspec := cg.args(node.Args)
	simple := byte(0)
	named := len(argspec.Spec)
	if method {
		if named == 0 && argspec.Unnamed <= 3 {
			simple = argspec.Unnamed + 1
		}
		if c, ok := mem.M.(*ast.Constant); ok && c.Val == SuStr("New") {
			panic("cannot explicitly call New method")
		}
		cg.expr(mem.M)
		cg.emit(op.CALLMETH + simple)
	} else {
		if named == 0 && argspec.Unnamed <= 4 {
			simple = argspec.Unnamed + 1
		}
		cg.expr(fn)
		cg.emit(op.CALLFUNC + simple)
	}
	if simple == 0 {
		cg.emit(argspec.Unnamed)
		verify.That(named <= math.MaxUint8)
		cg.emit(byte(named))
		cg.emit(argspec.Spec...)
	}
}

// generates code to push the arguments and returns an ArgSpec
// for non-string argument names we build an object and use @args
func (cg *cgen) args(args []ast.Arg) ArgSpec {
	if len(args) == 1 {
		if args[0].Name == SuStr("@") {
			cg.expr(args[0].E)
			return ArgSpec{Unnamed: EACH}
		} else if args[0].Name == SuStr("@+1") {
			cg.expr(args[0].E)
			return ArgSpec{Unnamed: EACH1}
		}
	}
	if !stringNames(args) {
		return cg.objectArgs(args)
	}
	var spec []byte
	for _, arg := range args {
		if arg.Name != nil {
			i := cg.name(arg.Name.ToStr())
			verify.That(i <= math.MaxUint8)
			spec = append(spec, byte(i))
		}
		cg.expr(arg.E)
	}
	verify.That(len(args) < int(EACH))
	return ArgSpec{Unnamed: byte(len(args) - len(spec)), Spec: spec}
}

// stringNames returns true if there are no non-string argument names
// might want to limit to identifiers?
func stringNames(args []ast.Arg) bool {
	for _, arg := range args {
		if arg.Name != nil {
			if _, ok := arg.Name.(SuStr); !ok {
				return false
			}
		}
	}
	return true
}

func (cg *cgen) objectArgs(args []ast.Arg) ArgSpec {
	var vis []int
	for _, arg := range args {
		if arg.Name != nil {
			vi := cg.value(arg.Name)
			verify.That(vi <= math.MaxUint16)
			vis = append(vis, vi)
		}
		cg.expr(arg.E)
	}
	cg.emit(op.OBJECT)
	cg.emit(byte(len(args) - len(vis))) // unnamed
	cg.emit(byte(len(vis)))             // named
	for _, vi := range vis {
		cg.emit(byte(vi>>8), byte(vi))
	}
	return ArgSpec{Unnamed: EACH}
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
