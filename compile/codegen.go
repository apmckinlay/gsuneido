package compile

// See also Disasm

import (
	"fmt"
	"math"
	"sync/atomic"
	"unsafe"

	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/lexer/tokens"
	. "github.com/apmckinlay/gsuneido/runtime"
	op "github.com/apmckinlay/gsuneido/runtime/opcodes"
	"github.com/apmckinlay/gsuneido/util/ints"
	"github.com/apmckinlay/gsuneido/util/str"
	"github.com/apmckinlay/gsuneido/util/verify"
)

//TODO source statement position for errors

// cgen is the context/results for compiling a function or block
type cgen struct {
	outerFn *ast.Function
	// outerVars is used to determine if block is closure
	// it is lazily derived from outerFn
	outerVars      map[string]bool
	code           []byte
	argspecs       []*ArgSpec
	base           Gnum
	isNew          bool
	isBlock        bool
	firstStatement bool
	ParamSpec
	// srcPos contains pairs of source and code position deltas
	srcPos []byte
	// srcBase is the starting point for the srcPos source deltas
	srcBase  int
	srcPrev  int
	codePrev int
}

// codegen compiles an Ast to an SuFunc
func codegen(fn *ast.Function) *SuFunc {
	return codegen2(fn, false)
}

func codegen2(fn *ast.Function, isBlock bool) *SuFunc {
	cg := cgen{outerFn: fn, base: fn.Base, isNew: fn.IsNewMethod, isBlock: isBlock}
	return cg.codegen(fn)
}

func (cg *cgen) codegen(fn *ast.Function) *SuFunc {
	cg.function(fn)
	cg.finishParamSpec()
	for _, as := range cg.argspecs {
		as.Names = cg.Values
	}

	return &SuFunc{
		Code:      *(*string)(unsafe.Pointer(&cg.code)),
		Nlocals:   uint8(len(cg.Names)),
		ParamSpec: cg.ParamSpec,
		ArgSpecs:  cg.argspecs,
		Id:        fn.Id,
		SrcPos:    *(*string)(unsafe.Pointer(&cg.srcPos)),
		SrcBase:   cg.srcBase,
	}
}

func codegenBlock(ast *ast.Function,
	outerFn *ast.Function, outerNames []string) (*SuFunc, []string) {
	base := len(outerNames)
	cg := cgen{outerFn: outerFn, base: ast.Base, isBlock: true}
	cg.Names = outerNames

	f := cg.codegen(ast)

	// hide parameters from outer function
	outerNames = f.Names
	f.Names = make([]string, len(outerNames))
	verify.That(base <= math.MaxUint8)
	f.Offset = uint8(base)
	copy(f.Names, outerNames)
	for i := 0; i < int(f.Nparams); i++ {
		outerNames[base+i] = ""
	}
	return f, outerNames
}

func (cg *cgen) finishParamSpec() {
	if !allZero(cg.Flags) {
		return
	}
	cg.Flags = zeroFlags[:len(cg.Flags)]
	if 0 <= cg.Nparams && cg.Nparams <= 4 {
		cg.Signature = ^(1 + cg.Nparams)
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

// zeroFlags is shared/reused for all zero flags
var zeroFlags [MaxArgs]Flag

// binary and nary ast node token to operation
var tok2op = [tok.Ntokens]op.Opcode{
	tok.Add:      op.Add,
	tok.Sub:      op.Sub,
	tok.Cat:      op.Cat,
	tok.Mul:      op.Mul,
	tok.Div:      op.Div,
	tok.Mod:      op.Mod,
	tok.LShift:   op.LeftShift,
	tok.RShift:   op.RightShift,
	tok.BitOr:    op.BitOr,
	tok.BitAnd:   op.BitAnd,
	tok.BitXor:   op.BitXor,
	tok.AddEq:    op.Add,
	tok.SubEq:    op.Sub,
	tok.CatEq:    op.Cat,
	tok.MulEq:    op.Mul,
	tok.DivEq:    op.Div,
	tok.ModEq:    op.Mod,
	tok.LShiftEq: op.LeftShift,
	tok.RShiftEq: op.RightShift,
	tok.BitOrEq:  op.BitOr,
	tok.BitAndEq: op.BitAnd,
	tok.BitXorEq: op.BitXor,
	tok.Is:       op.Is,
	tok.Isnt:     op.Isnt,
	tok.Match:    op.Match,
	tok.MatchNot: op.MatchNot,
	tok.Lt:       op.Lt,
	tok.Lte:      op.Lte,
	tok.Gt:       op.Gt,
	tok.Gte:      op.Gte,
	tok.And:      op.And,
	tok.Or:       op.Or,
}

func (cg *cgen) function(fn *ast.Function) {
	cg.params(fn.Params)
	cg.chainNew(fn)
	stmts := fn.Body
	cg.firstStatement = true
	for si, stmt := range stmts {
		cg.statement(stmt, nil, si == len(stmts)-1)
		cg.firstStatement = false
	}
}

func (cg *cgen) params(params []ast.Param) {
	cg.Nparams = uint8(len(params))
	for _, p := range params {
		name, flags := param(p.Name)
		if flags == AtParam && len(params) != 1 {
			panic("@param must be the only parameter")
		}
		cg.Names = append(cg.Names, name) // no duplicate reuse
		cg.Flags = append(cg.Flags, flags)
		if p.DefVal != nil {
			cg.Ndefaults++
			cg.Values = append(cg.Values, p.DefVal) // no duplicate reuse
		}
	}
}

func (cg *cgen) chainNew(fn *ast.Function) {
	if !fn.IsNewMethod || hasSuperCall(fn.Body) || cg.base <= 0 {
		return
	}
	cg.emit(op.This)
	cg.emitValue(SuStr("New"))
	cg.emitUint16(op.Super, cg.base)
	cg.emitUint8(op.CallMeth, 0)
}

func hasSuperCall(stmts []ast.Statement) bool {
	if len(stmts) < 1 {
		return false
	}
	expr, ok := stmts[0].(*ast.ExprStmt)
	if !ok {
		return false
	}
	call, ok := expr.E.(*ast.Call)
	if !ok {
		return false
	}
	fn, ok := call.Fn.(*ast.Ident)
	return ok && fn.Name == "super"
}

func param(p string) (string, Flag) {
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

// savePos saves source to code position information (for call stacks)
func (cg *cgen) savePos(sp int) {
	if cg.srcPos == nil {
		cg.srcBase = sp
		cg.srcPrev = sp
		cg.codePrev = len(cg.code)
		cg.srcPos = make([]byte, 0, 8)
	} else {
		ds := sp - cg.srcPrev
		dc := len(cg.code) - cg.codePrev
		for ds > 0 || dc > 0 {
			ns := ints.Min(math.MaxUint8, ds)
			nc := ints.Min(math.MaxUint8, dc)
			cg.srcPos = append(cg.srcPos, byte(ns), byte(nc))
			ds -= ns
			dc -= nc
		}
		cg.srcPrev = sp
		cg.codePrev = len(cg.code)
	}
}

func (cg *cgen) statement(node ast.Node, labels *Labels, lastStmt bool) {
	cg.savePos(node.(ast.Statement).Position())
	switch node := node.(type) {
	case *ast.Compound:
		for _, stmt := range node.Body {
			cg.statement(stmt, labels, lastStmt)
		}
	case *ast.Return:
		cg.returnStmt(node, lastStmt)
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
	case *ast.ForIn:
		cg.forInStmt(node)
	case *ast.Throw:
		cg.expr(node.E)
		cg.emit(op.Throw)
	case *ast.TryCatch:
		cg.tryCatchStmt(node, labels)
	case *ast.Break:
		cg.breakStmt(labels)
	case *ast.Continue:
		cg.continueStmt(labels)
	case *ast.ExprStmt:
		cg.expr(node.E)
		if !lastStmt {
			cg.emit(op.Pop)
		}
	default:
		panic("unexpected statement type " + fmt.Sprintf("%T", node))
	}
}

func (cg *cgen) statements(stmts []ast.Statement, labels *Labels, lastStmt bool) {
	for _, stmt := range stmts {
		cg.statement(stmt, labels, lastStmt)
	}
}

func (cg *cgen) returnStmt(node *ast.Return, lastStmt bool) {
	if node.E != nil {
		cg.expr(node.E)
	}
	if cg.isBlock {
		if node.E == nil {
			cg.emit(op.BlockReturnNil)
		} else {
			cg.emit(op.BlockReturn)
		}
	} else {
		if !lastStmt {
			if node.E == nil {
				cg.emit(op.ReturnNil)
			} else {
				cg.emit(op.Return)
			}
		}
	}
}

func (cg *cgen) breakStmt(labels *Labels) {
	if labels != nil {
		labels.brk = cg.emitJump(op.Jump, labels.brk)
	} else if cg.isBlock {
		cg.emit(op.BlockBreak)
	} else {
		panic("break can only be used within a loop")
	}
}

func (cg *cgen) continueStmt(labels *Labels) {
	if labels != nil {
		cg.emitBwdJump(op.Jump, labels.cont)
	} else if cg.isBlock {
		cg.emit(op.BlockContinue)
	} else {
		panic("continue can only be used within a loop")
	}
}

func (cg *cgen) ifStmt(node *ast.If, labels *Labels) {
	cg.expr(node.Cond)
	f := cg.emitJump(op.JumpFalse, -1)
	cg.statement(node.Then, labels, false)
	if node.Else != nil {
		end := cg.emitJump(op.Jump, -1)
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
				caseBody = cg.emitJump(op.JumpIs, -1)
			} else {
				afterCase = cg.emitJump(op.JumpIsnt, -1)
			}
		}
		cg.placeLabel(caseBody)
		cg.statements(c.Body, labels, false)
		end = cg.emitJump(op.Jump, end)
		cg.placeLabel(afterCase)
	}
	cg.emit(op.Pop)
	if node.Default != nil {
		cg.statements(node.Default, labels, false)
	} else {
		cg.emitValue(SuStr("unhandled switch value"))
		cg.emit(op.Throw)
	}
	cg.placeLabel(end)
}

func (cg *cgen) foreverStmt(node *ast.Forever) {
	labels := cg.newLabels()
	cg.statement(node.Body, labels, false)
	cg.emitJump(op.Jump, labels.cont-len(cg.code)-3)
	cg.placeLabel(labels.brk)
}

func (cg *cgen) whileStmt(node *ast.While) {
	labels := cg.newLabels()
	cond := cg.emitJump(op.Jump, -1)
	loop := cg.label()
	cg.statement(node.Body, labels, false)
	cg.placeLabel(cond)
	cg.expr(node.Cond)
	cg.emitBwdJump(op.JumpTrue, loop)
	cg.placeLabel(labels.brk)
}

func (cg *cgen) dowhileStmt(node *ast.DoWhile) {
	labels := cg.newLabels()
	cg.statement(node.Body, labels, false)
	cg.expr(node.Cond)
	cg.emitBwdJump(op.JumpTrue, labels.cont)
	cg.placeLabel(labels.brk)
}

func (cg *cgen) forStmt(node *ast.For) {
	cg.exprList(node.Init)
	labels := cg.newLabels()
	cond := -1
	if node.Cond != nil {
		cond = cg.emitJump(op.Jump, -1)
	}
	loop := cg.label()
	cg.statement(node.Body, labels, false)
	cg.exprList(node.Inc) // increment
	if node.Cond == nil {
		cg.emitBwdJump(op.Jump, loop)
	} else {
		cg.placeLabel(cond)
		cg.expr(node.Cond)
		cg.emitBwdJump(op.JumpTrue, loop)
	}
	cg.placeLabel(labels.brk)
}

func (cg *cgen) forInStmt(node *ast.ForIn) {
	cg.expr(node.E)
	cg.emit(op.Iter)
	labels := cg.newLabels()
	cg.emitForIn(node.Var, labels)
	cg.statement(node.Body, labels, false)
	cg.emitBwdJump(op.Jump, labels.cont)
	cg.placeLabel(labels.brk)
	cg.emit(op.Pop)
}

func (cg *cgen) emitForIn(name string, labels *Labels) {
	i := cg.name(name)
	adr := len(cg.code)
	cg.emit(op.ForIn, byte(labels.brk>>8), byte(labels.brk), byte(i))
	labels.brk = adr
}

func (cg *cgen) tryCatchStmt(node *ast.TryCatch, labels *Labels) {
	catch := cg.emitJump(op.Try, -1)
	cg.emitMore(byte(cg.value(SuStr(node.CatchFilter))))
	cg.statement(node.Try, labels, false)
	after := cg.emitJump(op.Catch, -1)
	cg.placeLabel(catch)
	if node.CatchVar != "" {
		cg.emit(op.Store, byte(cg.name(node.CatchVar)))
	}
	cg.emit(op.Pop)
	if node.Catch != nil {
		cg.statement(node.Catch, labels, false)
	}
	cg.placeLabel(after)
}

func (cg *cgen) exprList(list []ast.Expr) {
	for _, expr := range list {
		cg.expr(expr)
		cg.emit(op.Pop)
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
		cg.emit(op.Get)
	case *ast.RangeTo:
		cg.expr(node.E)
		cg.exprOr(node.From, op.Zero)
		cg.exprOr(node.To, op.MaxInt)
		cg.emit(op.RangeTo)
	case *ast.RangeLen:
		cg.expr(node.E)
		cg.exprOr(node.From, op.Zero)
		cg.exprOr(node.Len, op.MaxInt)
		cg.emit(op.RangeLen)
	case *ast.In:
		cg.inExpr(node)
	case *ast.Call:
		cg.call(node)
	case *ast.Function:
		fn := codegen(node)
		cg.emitValue(fn)
	case *ast.Block:
		cg.block(node)
	default:
		panic("unhandled expression type: " + fmt.Sprintf("%T", node))
	}
}

func (cg *cgen) exprOr(expr ast.Expr, op op.Opcode) {
	if expr == nil {
		cg.emit(op)
	} else {
		cg.expr(expr)
	}
}

func (cg *cgen) identifier(node *ast.Ident) {
	if node.Name == "this" {
		cg.emit(op.This)
	} else if isLocal(node.Name) {
		i := cg.name(node.Name)
		if node.Name[0] == '_' {
			cg.emitUint8(op.Dyload, i)
		} else {
			cg.emitUint8(op.Load, i)
		}
	} else {
		cg.emitUint16(op.Global, Global.Num(node.Name))
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
	if i > math.MaxUint8 {
		panic("too many local variables (>255)")
	}
	cg.Names = append(cg.Names, s)
	return i
}

func (cg *cgen) unary(node *ast.Unary) {
	if node.Tok == tok.LParen {
		cg.expr(node.E)
		return
	}
	o := utok2op[node.Tok]
	if tok.Inc <= node.Tok && node.Tok <= tok.PostDec {
		ref := cg.lvalue(node.E)
		cg.dupLvalue(ref)
		cg.load(ref)
		if node.Tok == tok.PostInc || node.Tok == tok.PostDec {
			cg.dupUnderLvalue(ref)
			cg.emit(op.One)
			cg.emit(o)
			cg.store(ref)
			cg.emit(op.Pop)
		} else {
			cg.emit(op.One)
			cg.emit(o)
			cg.store(ref)
		}
	} else {
		cg.expr(node.E)
		cg.emit(o)
	}
}

// Unary ast expr node token to operation
var utok2op = [tok.Ntokens]op.Opcode{
	tok.Add:     op.UnaryPlus,
	tok.Sub:     op.UnaryMinus,
	tok.Not:     op.Not,
	tok.BitNot:  op.BitNot,
	tok.Inc:     op.Add,
	tok.PostInc: op.Add,
	tok.Dec:     op.Sub,
	tok.PostDec: op.Sub,
}

func (cg *cgen) binary(node *ast.Binary) {
	switch node.Tok {
	case tok.Eq:
		ref := cg.lvalue(node.Lhs)
		cg.expr(node.Rhs)
		cg.store(ref)
	case tok.AddEq, tok.SubEq, tok.CatEq, tok.MulEq, tok.DivEq, tok.ModEq,
		tok.LShiftEq, tok.RShiftEq, tok.BitOrEq, tok.BitAndEq, tok.BitXorEq:
		ref := cg.lvalue(node.Lhs)
		cg.dupLvalue(ref)
		cg.load(ref)
		cg.expr(node.Rhs)
		cg.emit(tok2op[node.Tok])
		cg.store(ref)
	case tok.Is, tok.Isnt, tok.Match, tok.MatchNot, tok.Mod,
		tok.LShift, tok.RShift, tok.Lt, tok.Lte, tok.Gt, tok.Gte:
		cg.expr(node.Lhs)
		cg.expr(node.Rhs)
		cg.emit(tok2op[node.Tok])
	default:
		panic("unhandled binary operation " + node.Tok.String())
	}
}

func (cg *cgen) nary(node *ast.Nary) {
	if node.Tok == tok.And || node.Tok == tok.Or {
		cg.andorExpr(node)
	} else if node.Tok == tok.Mul {
		cg.muldivExpr(node)
	} else {
		o := tok2op[node.Tok]
		cg.expr(node.Exprs[0])
		for _, e := range node.Exprs[1:] {
			if node.Tok == tok.Add && isUnary(e, tok.Sub) {
				cg.expr(e.(*ast.Unary).E)
				cg.emit(op.Sub)
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
	lastExpr := node.Exprs[len(node.Exprs)-1]
	if !isCompare(lastExpr) {
		cg.emit(op.Bool) // not needed if last expr was comparison
	}
	cg.placeLabel(label)
}

func isCompare(e ast.Expr) bool {
	bin, ok := e.(*ast.Binary)
	return ok && tok.CompareStart < bin.Tok && bin.Tok < tok.CompareEnd
}

func (cg *cgen) muldivExpr(node *ast.Nary) {
	var divs []ast.Expr
	cg.expr(node.Exprs[0])
	for _, e := range node.Exprs[1:] {
		if isUnary(e, tok.Div) {
			divs = append(divs, e.(*ast.Unary).E)
		} else {
			cg.expr(e)
			cg.emit(op.Mul)
		}
	}
	if len(divs) > 0 {
		cg.expr(divs[0])
		for _, e := range divs[1:] {
			cg.expr(e)
			cg.emit(op.Mul)
		}
		cg.emit(op.Div)
	}
}

func isUnary(e ast.Expr, tok tok.Token) bool {
	u, ok := e.(*ast.Unary)
	return ok && u.Tok == tok
}

func (cg *cgen) qcExpr(node *ast.Trinary) {
	f, end := -1, -1
	cg.expr(node.Cond)
	f = cg.emitJump(op.QMark, f)
	cg.expr(node.T)
	end = cg.emitJump(op.Jump, end)
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
			end = cg.emitJump(op.In, end)
		} else {
			cg.emit(op.Is)
		}
	}
	cg.placeLabel(end)
}

func (cg *cgen) emitValue(val Value) {
	if val == True {
		cg.emit(op.True)
	} else if val == False {
		cg.emit(op.False)
	} else if val == Zero {
		cg.emit(op.Zero)
	} else if val == One {
		cg.emit(op.One)
	} else if val == EmptyStr {
		cg.emit(op.EmptyStr)
	} else if i, ok := SuIntToInt(val); ok {
		cg.emitInt16(op.Int, i)
	} else {
		cg.emitUint8(op.Value, cg.value(val))
	}
}

// value returns an index for the constant value
// reusing if duplicate, adding otherwise
func (cg *cgen) value(v Value) int {
	for i, v2 := range cg.Values {
		if v.Equal(v2) {
			return i
		}
	}
	i := len(cg.Values)
	if i > math.MaxUint8 {
		panic("too many constants (>255)")
	}
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
		cg.emit(op.Get)
	} else {
		if cg.Names[ref][0] == '_' {
			cg.emitUint8(op.Dyload, ref)
		} else {
			cg.emitUint8(op.Load, ref)
		}
	}
}

func (cg *cgen) store(ref int) {
	if ref == memRef {
		cg.emit(op.Put)
	} else {
		cg.emitUint8(op.Store, ref)
	}
}

func (cg *cgen) dupLvalue(ref int) {
	if ref == memRef {
		cg.emit(op.Dup2)
	}
}

func (cg *cgen) dupUnderLvalue(ref int) {
	if ref == memRef {
		cg.emit(op.Dupx2)
	} else {
		cg.emit(op.Dup)
	}
}

var superNew = &ast.Mem{
	E: &ast.Ident{Name: "super"},
	M: &ast.Constant{Val: SuStr("New")}}

func (cg *cgen) call(node *ast.Call) {
	fn := node.Fn

	if id, ok := fn.(*ast.Ident); ok && id.Name == "super" {
		if !cg.isNew {
			panic("super(...) only valid in New method")
		}
		if !cg.firstStatement {
			panic("super(...) must be first statement")
		}
		fn = superNew // super(...) => super.New(...)
	}

	mem, method := fn.(*ast.Mem)
	superCall := false
	if method {
		if x, ok := mem.E.(*ast.Ident); ok && x.Name == "super" {
			superCall = true
			if cg.base <= 0 {
				panic("super requires parent")
			}
			cg.emit(op.This)
		} else {
			cg.expr(mem.E)
		}
	}
	argspec := cg.args(node.Args)
	if method {
		if fn != superNew {
			if c, ok := mem.M.(*ast.Constant); ok && c.Val == SuStr("New") {
				panic("cannot explicitly call New method")
			}
		}
		cg.expr(mem.M)
		if superCall {
			cg.emitUint16(op.Super, cg.base)
		}
		cg.emit(op.CallMeth)
	} else {
		cg.expr(fn)
		cg.emit(op.CallFunc)
	}
	verify.That(argspec < math.MaxUint8)
	cg.emitMore(byte(argspec))
}

// generates code to push the arguments and returns an ArgSpec index
func (cg *cgen) args(args []ast.Arg) int {
	if len(args) == 1 {
		if args[0].Name == SuStr("@") {
			cg.expr(args[0].E)
			return AsEach
		} else if args[0].Name == SuStr("@+1") {
			cg.expr(args[0].E)
			return AsEach1
		}
	}
	var spec []byte
	for _, arg := range args {
		if arg.Name != nil {
			i := cg.value(arg.Name)
			spec = append(spec, byte(i))
		}
		cg.expr(arg.E)
	}
	verify.That(len(args) < math.MaxUint8)
	return cg.argspec(&ArgSpec{Nargs: byte(len(args)), Spec: spec})
}

func (cg *cgen) argspec(as *ArgSpec) int {
	as.Names = cg.Values // not final, but needed for Equal
	for i, a := range StdArgSpecs {
		if as.Equal(a) {
			return i
		}
	}
	for i, a := range cg.argspecs {
		if cg.argSpecEq(a, as) {
			return i + len(StdArgSpecs)
		}
	}
	cg.argspecs = append(cg.argspecs, as)
	return len(cg.argspecs) - 1 + len(StdArgSpecs)
}

// argSpecEq checks if two ArgSpec's are equal
// using cg.Values instead of the ArgSpec Names
// We can't set argspec.Names = cg.Values yet
// because cg.Values is still growing and may be reallocated.
func (cg *cgen) argSpecEq(a1, a2 *ArgSpec) bool {
	if a1.Nargs != a2.Nargs || a1.Each != a2.Each || len(a1.Spec) != len(a2.Spec) {
		return false
	}
	for i := range a1.Spec {
		if !cg.Values[a1.Spec[i]].Equal(cg.Values[a2.Spec[i]]) {
			return false
		}
	}
	return true
}

var funcId uint32

func (cg *cgen) block(b *ast.Block) {
	f := &b.Function
	blockVars := ast.VarSet(f)
	itParam(f, blockVars)
	var fn *SuFunc
	if cg.blockIsFunction(blockVars) {
		fn = codegen2(f, true)
		cg.emitValue(fn)
	} else {
		// closure
		fn, cg.Names = codegenBlock(f, cg.outerFn, cg.Names)
		i := cg.value(fn)
		cg.emitUint8(op.Block, i)
	}
	if cg.outerFn.Id == 0 {
		cg.outerFn.Id = atomic.AddUint32(&funcId, 1)
	}
	fn.OuterId = cg.outerFn.Id
}

func itParam(f *ast.Function, blockVars map[string]bool) {
	if len(f.Params) == 0 && blockVars["it"] {
		// automatic "it" parameter
		f.Params = []ast.Param{{Name: "it"}}
	}
}

// blockIsFunction returns false if the block is a closure
// i.e. if it shares any variables with the enclosing function
func (cg *cgen) blockIsFunction(blockVars map[string]bool) bool {
	// a block cannot have its own "this"
	// so any reference means it must be a closure (not a function)
	if blockVars["this"] {
		return false
	}
	if cg.outerVars == nil {
		cg.outerVars = ast.VarSet(cg.outerFn) // cache
	}
	for v := range blockVars {
		if cg.outerVars[v] {
			return false
		}
	}
	return true
}

// helpers ---------------------------------------------------------------------

// emit is used to append an op code
func (cg *cgen) emit(op op.Opcode, b ...byte) {
	cg.code = append(append(cg.code, byte(op)), b...)
}

func (cg *cgen) emitMore(b ...byte) {
	cg.code = append(cg.code, b...)
}

func (cg *cgen) emitUint8(op op.Opcode, i int) {
	verify.That(0 <= i && i < math.MaxUint8)
	cg.emit(op, byte(i))
}

func (cg *cgen) emitInt16(op op.Opcode, i int) {
	verify.That(math.MinInt16 <= i && i <= math.MaxInt16)
	cg.emit(op, byte(i>>8), byte(i))
}

func (cg *cgen) emitUint16(op op.Opcode, i int) {
	verify.That(0 <= i && i < math.MaxUint16)
	cg.emit(op, byte(i>>8), byte(i))
}

func (cg *cgen) emitJump(op op.Opcode, label int) int {
	adr := len(cg.code)
	verify.That(math.MinInt16 <= label && label <= math.MaxInt16)
	cg.emit(op, byte(label>>8), byte(label))
	return adr
}

func (cg *cgen) emitBwdJump(op op.Opcode, label int) {
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

// newLabels should be called where continue should go
func (cg *cgen) newLabels() *Labels {
	return &Labels{brk: -1, cont: cg.label()}
}
