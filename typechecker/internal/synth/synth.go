// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

// Package synth generates random Suneido class sources for property-based
package synth

import (
	"fmt"
	"maps"
	"slices"
	"sort"
	"strings"

	"pgregory.net/rapid"
)

type Config struct {
	WellTyped bool
}

// ---- IR ----

type LitKind int

const (
	LitNum LitKind = iota
	LitStr
	LitDate
	LitTrue
	LitFalse
	LitObj
)

type Lit struct {
	Kind LitKind
	N    int
	S    string
	D    string
}

type ExprKind int

const (
	ExLit ExprKind = iota
	ExVar
	ExMember   // .fN
	ExBin      // (a OP b)
	ExUnary    // (not a)
	ExTern     // (c ? a : b)
	ExCall     // .MN(args...)
	ExPred     // (recv.Pred?())
	ExFnPred   // Pred?(arg)  - function form (narrows a param inside a guard)
	ExRecvCall // recv.Method(args...)  - call on a named receiver
)

type Expr struct {
	Kind ExprKind
	Lit  Lit
	Name string // ExVar name / ExMember member / ExCall method
	Op   string // ExBin / ExUnary operator symbol
	Pred string // ExPred predicate name (e.g. "Number?")
	Args []Expr // sub-expressions
}

type StmtKind int

const (
	StAssign       StmtKind = iota // vN = expr
	StMemberAssign                 // .fN = expr
	StIf
	StWhile
	StReturn
	StExpr // a self-call statement
)

type Stmt struct {
	Kind StmtKind
	Name string // target var / member
	Expr Expr   // rhs / cond / return value
	Bare bool   // bare `return`
	Body []Stmt
	Else []Stmt
}

type Param struct {
	Name       string
	HasDefault bool
	Def        Lit
	Rest       bool // @rest
}

type Method struct {
	Name     string
	Params   []Param
	Body     []Stmt
	RetAnn   string // "" = no annotation
	Ret      GType  // intended return type (WellTyped mode ground truth)
	RetKnown bool   // Ret is meaningful
}

// CtorMode records the optional-member protocol for constructor testing:
// a member seeded `: false` that only New may assign, at type Payload.
type CtorMode int

const (
	CtorNone   CtorMode = iota // a regular member, not part of the protocol
	CtorAlways                 // New assigns it on every path -> must demote
	CtorMaybe                  // New assigns it conditionally -> must not demote
	CtorNever                  // New never assigns it -> must not demote
)

type Member struct {
	Name    string
	Seed    Lit
	Ctor    CtorMode // ground truth for demotion (wellTyped mode only)
	Payload GType    // assigned type when Ctor != CtorNone; never GtBool, whose
	// TBoolean-typed assignments contain False and block demotion (C1)
}

type Program struct {
	Name    string
	Base    string
	Members []Member
	Methods []Method
}

// ---- generators (general mode) ----

var (
	dateLits  = []string{"#20200101", "#19991231", "#20210615", "#20180302"}
	strLits   = []string{"a", "b", "hello", "x0"}
	binOps    = []string{"+", "-", "*", "/", "%", "$", "<", "<=", ">", ">=", "is", "isnt", "and", "or"}
	predNames = []string{"Number?", "String?", "Date?", "Object?", "Boolean?"}
)

func genExprLit(t *rapid.T) Lit {
	return Lit{Kind: LitNum, N: rapid.IntRange(1, 99).Draw(t, "exprnum")}
}

func genLit(t *rapid.T) Lit {
	switch rapid.IntRange(0, 5).Draw(t, "litkind") {
	case 0:
		return Lit{Kind: LitNum, N: rapid.IntRange(1, 99).Draw(t, "num")}
	case 1:
		return Lit{Kind: LitStr, S: rapid.SampledFrom(strLits).Draw(t, "str")}
	case 2:
		return Lit{Kind: LitDate, D: rapid.SampledFrom(dateLits).Draw(t, "date")}
	case 3:
		return Lit{Kind: LitTrue}
	case 4:
		return Lit{Kind: LitFalse}
	default:
		return Lit{Kind: LitObj}
	}
}

type genScope struct {
	locals      []string
	Params      []string
	Members     []string
	methodNames []string
}

func (sc *genScope) vars() []string {
	out := make([]string, 0, len(sc.locals)+len(sc.Params))
	out = append(out, sc.locals...)
	out = append(out, sc.Params...)
	return out
}

func genExpr(t *rapid.T, sc *genScope, depth int) Expr {
	kinds := []string{"lit"}
	if len(sc.vars()) > 0 {
		kinds = append(kinds, "var")
	}
	if len(sc.Members) > 0 {
		kinds = append(kinds, "member")
	}
	if depth > 0 {
		kinds = append(kinds, "bin", "unary", "tern")
		if len(sc.methodNames) > 0 {
			kinds = append(kinds, "call")
		}
		if len(sc.vars()) > 0 || len(sc.Members) > 0 {
			kinds = append(kinds, "pred")
		}
	}
	switch rapid.SampledFrom(kinds).Draw(t, "exprkind") {
	case "var":
		return Expr{Kind: ExVar, Name: rapid.SampledFrom(sc.vars()).Draw(t, "var")}
	case "member":
		return Expr{Kind: ExMember, Name: rapid.SampledFrom(sc.Members).Draw(t, "mem")}
	case "bin":
		return Expr{Kind: ExBin, Op: rapid.SampledFrom(binOps).Draw(t, "op"),
			Args: []Expr{genExpr(t, sc, depth-1), genExpr(t, sc, depth-1)}}
	case "unary":
		return Expr{Kind: ExUnary, Op: "not", Args: []Expr{genExpr(t, sc, depth-1)}}
	case "tern":
		return Expr{Kind: ExTern, Args: []Expr{
			genExpr(t, sc, depth-1), genExpr(t, sc, depth-1), genExpr(t, sc, depth-1)}}
	case "call":
		return genCall(t, sc, depth-1)
	case "pred":
		var recv Expr
		if len(sc.vars()) > 0 && (len(sc.Members) == 0 || rapid.Bool().Draw(t, "recvvar")) {
			recv = Expr{Kind: ExVar, Name: rapid.SampledFrom(sc.vars()).Draw(t, "predrecv")}
		} else {
			recv = Expr{Kind: ExMember, Name: rapid.SampledFrom(sc.Members).Draw(t, "predrecvm")}
		}
		return Expr{Kind: ExPred, Pred: rapid.SampledFrom(predNames).Draw(t, "pred"), Args: []Expr{recv}}
	default: // "lit"
		return Expr{Kind: ExLit, Lit: genExprLit(t)}
	}
}

func genCall(t *rapid.T, sc *genScope, depth int) Expr {
	m := rapid.SampledFrom(sc.methodNames).Draw(t, "callee")
	nargs := rapid.IntRange(0, 3).Draw(t, "nargs")
	args := make([]Expr, nargs)
	for i := range args {
		args[i] = genExpr(t, sc, depth)
	}
	return Expr{Kind: ExCall, Name: m, Args: args}
}

func genStmt(t *rapid.T, sc *genScope, depth int) Stmt {
	kinds := []string{"assign", "return", "expr"}
	if len(sc.Members) > 0 {
		kinds = append(kinds, "massign")
	}
	if depth > 0 {
		kinds = append(kinds, "if", "while")
	}
	switch rapid.SampledFrom(kinds).Draw(t, "stmtkind") {
	case "massign":
		return Stmt{Kind: StMemberAssign,
			Name: rapid.SampledFrom(sc.Members).Draw(t, "mtgt"), Expr: genExpr(t, sc, 2)}
	case "return":
		if rapid.Bool().Draw(t, "bareret") {
			return Stmt{Kind: StReturn, Bare: true}
		}
		return Stmt{Kind: StReturn, Expr: genExpr(t, sc, 2)}
	case "expr":
		return Stmt{Kind: StExpr, Expr: genCall(t, sc, 1)}
	case "if":
		st := Stmt{Kind: StIf, Expr: genExpr(t, sc, 2), Body: genStmts(t, sc, depth-1, 3)}
		if rapid.Bool().Draw(t, "haselse") {
			st.Else = genStmts(t, sc, depth-1, 3)
		}
		return st
	case "while":
		return Stmt{Kind: StWhile, Expr: genExpr(t, sc, 2), Body: genStmts(t, sc, depth-1, 3)}
	default: // "assign"
		name := fmt.Sprintf("v%d", rapid.IntRange(0, 4).Draw(t, "vidx"))
		e := genExpr(t, sc, 2)
		if !slices.Contains(sc.locals, name) {
			sc.locals = append(sc.locals, name)
		}
		return Stmt{Kind: StAssign, Name: name, Expr: e}
	}
}

func genStmts(t *rapid.T, sc *genScope, depth, maxStmts int) []Stmt {
	n := rapid.IntRange(0, maxStmts).Draw(t, "nstmts")
	out := make([]Stmt, 0, n)
	for range n {
		out = append(out, genStmt(t, sc, depth))
	}
	return out
}

func genMethod(t *rapid.T, name string, members, methodNames []string) Method {
	np := rapid.IntRange(0, 3).Draw(t, "nparams")
	params := make([]Param, np)
	paramNames := make([]string, np)
	sawDefault := false
	for i := range np {
		pn := fmt.Sprintf("p%d", i)
		paramNames[i] = pn
		if sawDefault || rapid.Bool().Draw(t, "hasdef") {
			sawDefault = true
			params[i] = Param{Name: pn, HasDefault: true, Def: genLit(t)}
		} else {
			params[i] = Param{Name: pn}
		}
	}
	sc := &genScope{Params: paramNames, Members: members, methodNames: methodNames}
	return Method{Name: name, Params: params, Body: genStmts(t, sc, 2, 5)}
}

func GenProgram(t *rapid.T, cfg Config) Program {
	if cfg.WellTyped {
		return genProgramWT(t)
	}
	nmem := rapid.IntRange(0, 3).Draw(t, "nmem")
	members := make([]Member, nmem)
	memberNames := make([]string, nmem)
	for i := range nmem {
		n := fmt.Sprintf("f%d", i)
		memberNames[i] = n
		members[i] = Member{Name: n, Seed: genLit(t)}
	}
	nmeth := rapid.IntRange(1, 3).Draw(t, "nmeth")
	methodNames := make([]string, nmeth)
	for i := range nmeth {
		methodNames[i] = fmt.Sprintf("M%d", i)
	}
	methods := make([]Method, nmeth)
	for i := range nmeth {
		methods[i] = genMethod(t, methodNames[i], memberNames, methodNames)
	}
	return Program{Name: "T", Members: members, Methods: methods}
}

// ---- renderer ----

func renderLit(l Lit) string {
	switch l.Kind {
	case LitNum:
		return fmt.Sprintf("%d", l.N)
	case LitStr:
		return fmt.Sprintf("%q", l.S)
	case LitDate:
		return l.D
	case LitTrue:
		return "true"
	case LitFalse:
		return "false"
	default: // LitObj
		return "#(1, 2, 3)"
	}
}

func renderExpr(e *Expr) string {
	switch e.Kind {
	case ExVar:
		return e.Name
	case ExMember:
		return "." + e.Name
	case ExBin:
		return "(" + renderExpr(&e.Args[0]) + " " + e.Op + " " + renderExpr(&e.Args[1]) + ")"
	case ExUnary:
		return "(" + e.Op + " " + renderExpr(&e.Args[0]) + ")"
	case ExTern:
		return "(" + renderExpr(&e.Args[0]) + " ? " + renderExpr(&e.Args[1]) + " : " + renderExpr(&e.Args[2]) + ")"
	case ExCall:
		parts := make([]string, len(e.Args))
		for i := range e.Args {
			parts[i] = renderExpr(&e.Args[i])
		}
		return "." + e.Name + "(" + strings.Join(parts, ", ") + ")"
	case ExPred:
		return "(" + renderExpr(&e.Args[0]) + "." + e.Pred + "())"
	case ExFnPred:
		return e.Pred + "(" + renderExpr(&e.Args[0]) + ")"
	case ExRecvCall:
		parts := make([]string, 0, len(e.Args)-1)
		for i := 1; i < len(e.Args); i++ {
			parts = append(parts, renderExpr(&e.Args[i]))
		}
		return renderExpr(&e.Args[0]) + "." + e.Name + "(" + strings.Join(parts, ", ") + ")"
	default: // ExLit
		return renderLit(e.Lit)
	}
}

func renderParams(ps []Param) string {
	parts := make([]string, len(ps))
	for i, p := range ps {
		switch {
		case p.Rest:
			parts[i] = "@" + p.Name
		case p.HasDefault:
			parts[i] = p.Name + " = " + renderLit(p.Def)
		default:
			parts[i] = p.Name
		}
	}
	return strings.Join(parts, ", ")
}

func writeIndent(b *strings.Builder, n int) {
	for range n {
		b.WriteByte('\t')
	}
}

func renderStmt(b *strings.Builder, s *Stmt, ind int) {
	switch s.Kind {
	case StMemberAssign:
		writeIndent(b, ind)
		fmt.Fprintf(b, ".%s = %s\n", s.Name, renderExpr(&s.Expr))
	case StReturn:
		writeIndent(b, ind)
		if s.Bare {
			b.WriteString("return\n")
		} else {
			fmt.Fprintf(b, "return %s\n", renderExpr(&s.Expr))
		}
	case StExpr:
		writeIndent(b, ind)
		b.WriteString(renderExpr(&s.Expr))
		b.WriteByte('\n')
	case StIf:
		writeIndent(b, ind)
		fmt.Fprintf(b, "if (%s)\n", renderExpr(&s.Expr))
		renderBlock(b, s.Body, ind)
		if len(s.Else) > 0 {
			writeIndent(b, ind)
			b.WriteString("else\n")
			renderBlock(b, s.Else, ind)
		}
	case StWhile:
		writeIndent(b, ind)
		fmt.Fprintf(b, "while (%s)\n", renderExpr(&s.Expr))
		renderBlock(b, s.Body, ind)
	default: // StAssign
		writeIndent(b, ind)
		fmt.Fprintf(b, "%s = %s\n", s.Name, renderExpr(&s.Expr))
	}
}

func renderBlock(b *strings.Builder, stmts []Stmt, ind int) {
	writeIndent(b, ind)
	b.WriteString("{\n")
	for i := range stmts {
		renderStmt(b, &stmts[i], ind+1)
	}
	writeIndent(b, ind)
	b.WriteString("}\n")
}

func renderMethod(b *strings.Builder, m *Method) {
	writeIndent(b, 1)
	fmt.Fprintf(b, "%s(%s)", m.Name, renderParams(m.Params))
	if m.RetAnn != "" {
		fmt.Fprintf(b, ": %s", m.RetAnn)
	}
	b.WriteByte('\n')
	renderBlock(b, m.Body, 1)
}

func Render(p *Program) string {
	var b strings.Builder
	if p.Base != "" {
		fmt.Fprintf(&b, "%s\n\t{\n", p.Base)
	} else {
		b.WriteString("class\n\t{\n")
	}
	for _, m := range p.Members {
		writeIndent(&b, 1)
		fmt.Fprintf(&b, "%s: %s\n", m.Name, renderLit(m.Seed))
	}
	for _, m := range p.Methods {
		renderMethod(&b, &m)
	}
	b.WriteString("\t}\n")
	return b.String()
}

type GType int

const (
	GtNum GType = iota
	GtStr
	GtDate
	GtObj
	GtBool
)

var (
	BaseTypes       = []GType{GtNum, GtStr, GtDate, GtObj, GtBool}
	comparableTypes = []GType{GtNum, GtStr, GtDate} // safe for < <= > >=
)

func GTypeLit(t *rapid.T, g GType) Lit {
	switch g {
	case GtNum:
		return Lit{Kind: LitNum, N: rapid.IntRange(1, 99).Draw(t, "wtnum")}
	case GtStr:
		return Lit{Kind: LitStr, S: rapid.SampledFrom(strLits).Draw(t, "wtstr")}
	case GtDate:
		return Lit{Kind: LitDate, D: rapid.SampledFrom(dateLits).Draw(t, "wtdate")}
	case GtObj:
		return Lit{Kind: LitObj}
	default: // GtBool
		if rapid.Bool().Draw(t, "wtbool") {
			return Lit{Kind: LitTrue}
		}
		return Lit{Kind: LitFalse}
	}
}

func PredName(g GType) string {
	switch g {
	case GtNum:
		return "Number?"
	case GtStr:
		return "String?"
	case GtDate:
		return "Date?"
	case GtObj:
		return "Object?"
	default:
		return "Boolean?"
	}
}

type wtSig struct {
	Name    string
	nparams int
	ret     GType
}

type wtScope struct {
	// localTypes is shared across all branch clones of a method: one declared
	// type per local, ever. sibling branches typing the same name differently
	// would make flow joins produce cross-type unions, so structural rewrites
	// (e.g. wrapping a statement in if(true)) could surface stale type arms.
	localTypes map[string]GType
	// assigned is per-branch: only locals assigned on this control path are
	// readable, so a first assignment inside a branch cannot leak into the
	// outer scope where the local may be unassigned at runtime.
	assigned map[string]bool
	narrowed map[string]GType // params narrowed to a clean type in this branch
	Members  map[string]GType
	Params   []string
	sigs     []wtSig
}

func (sc *wtScope) clone() *wtScope {
	return &wtScope{localTypes: sc.localTypes,
		assigned: maps.Clone(sc.assigned),
		narrowed: maps.Clone(sc.narrowed),
		Members:  sc.Members, Params: sc.Params, sigs: sc.sigs}
}

// visibleLocals lists locals readable in this scope at type want.
func (sc *wtScope) visibleLocals(want GType) []string {
	var out []string
	for n := range sc.assigned {
		if sc.localTypes[n] == want {
			out = append(out, n)
		}
	}
	sort.Strings(out)
	return out
}

func sortedNamesByType(m map[string]GType, want GType) []string {
	var out []string
	for k, v := range m {
		if v == want {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	return out
}

func sortedKeys(m map[string]GType) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// genTypedExpr returns an expression that infers to want (clean single type).
func genTypedExpr(t *rapid.T, sc *wtScope, want GType, depth int) Expr {
	var opts []func() Expr
	opts = append(opts, func() Expr { return Expr{Kind: ExLit, Lit: GTypeLit(t, want)} })
	opts = append(opts, wtNameOpts(sc, want)...)
	if depth > 0 {
		opts = append(opts, wtDeepOpts(t, sc, want, depth)...)
	}
	return opts[rapid.IntRange(0, len(opts)-1).Draw(t, "wtexpr")]()
}

// wtNameOpts are the in-scope names readable at type want.
func wtNameOpts(sc *wtScope, want GType) []func() Expr {
	var opts []func() Expr
	for _, n := range sc.visibleLocals(want) {
		nm := n
		opts = append(opts, func() Expr { return Expr{Kind: ExVar, Name: nm} })
	}
	for _, n := range sortedNamesByType(sc.narrowed, want) {
		nm := n
		opts = append(opts, func() Expr { return Expr{Kind: ExVar, Name: nm} })
	}
	for _, n := range sortedNamesByType(sc.Members, want) {
		nm := n
		opts = append(opts, func() Expr { return Expr{Kind: ExMember, Name: nm} })
	}
	return opts
}

// wtDeepOpts are the compound shapes allowed while depth remains.
func wtDeepOpts(t *rapid.T, sc *wtScope, want GType, depth int) []func() Expr {
	var opts []func() Expr
	switch want {
	case GtNum:
		opts = append(opts, func() Expr {
			return Expr{Kind: ExBin, Op: rapid.SampledFrom([]string{"+", "-", "*"}).Draw(t, "wtnumop"),
				Args: []Expr{genTypedExpr(t, sc, GtNum, depth-1), genTypedExpr(t, sc, GtNum, depth-1)}}
		})
	case GtStr:
		opts = append(opts, func() Expr {
			return Expr{Kind: ExBin, Op: "$",
				Args: []Expr{genTypedExpr(t, sc, GtStr, depth-1), genTypedExpr(t, sc, GtStr, depth-1)}}
		})
	case GtBool:
		opts = append(opts, wtBoolOpts(t, sc, depth)...)
	}
	// ternary: boolean condition, both arms of type want
	opts = append(opts, func() Expr {
		return Expr{Kind: ExTern, Args: []Expr{
			genTypedExpr(t, sc, GtBool, depth-1),
			genTypedExpr(t, sc, want, depth-1),
			genTypedExpr(t, sc, want, depth-1)}}
	})
	// self-call to a method whose return type is want, exact arity
	var callable []wtSig
	for _, s := range sc.sigs {
		if s.ret == want {
			callable = append(callable, s)
		}
	}
	if len(callable) > 0 {
		opts = append(opts, func() Expr {
			s := callable[rapid.IntRange(0, len(callable)-1).Draw(t, "wtcallee")]
			args := make([]Expr, s.nparams)
			for i := range args {
				args[i] = genTypedExpr(t, sc, rapid.SampledFrom(BaseTypes).Draw(t, "wtargt"), depth-1)
			}
			return Expr{Kind: ExCall, Name: s.Name, Args: args}
		})
	}
	return opts
}

// wtBoolOpts are the boolean-producing expression shapes.
func wtBoolOpts(t *rapid.T, sc *wtScope, depth int) []func() Expr {
	return []func() Expr{
		// same-typed comparison (avoids the cross-type-compare warning)
		func() Expr {
			ct := rapid.SampledFrom(comparableTypes).Draw(t, "wtcmpt")
			return Expr{Kind: ExBin, Op: rapid.SampledFrom([]string{"<", "<=", ">", ">="}).Draw(t, "wtcmp"),
				Args: []Expr{genTypedExpr(t, sc, ct, depth-1), genTypedExpr(t, sc, ct, depth-1)}}
		},
		// equality on same-typed operands
		func() Expr {
			et := rapid.SampledFrom(BaseTypes).Draw(t, "wteqt")
			return Expr{Kind: ExBin, Op: rapid.SampledFrom([]string{"is", "isnt"}).Draw(t, "wteq"),
				Args: []Expr{genTypedExpr(t, sc, et, depth-1), genTypedExpr(t, sc, et, depth-1)}}
		},
		// and / or on booleans
		func() Expr {
			return Expr{Kind: ExBin, Op: rapid.SampledFrom([]string{"and", "or"}).Draw(t, "wtlogic"),
				Args: []Expr{genTypedExpr(t, sc, GtBool, depth-1), genTypedExpr(t, sc, GtBool, depth-1)}}
		},
		// not
		func() Expr {
			return Expr{Kind: ExUnary, Op: "not", Args: []Expr{genTypedExpr(t, sc, GtBool, depth-1)}}
		},
	}
}

func genTypedStmt(t *rapid.T, sc *wtScope, depth int, ret GType, allowGuard bool) Stmt {
	kinds := []string{"assign", "return"}
	if len(sc.Members) > 0 {
		kinds = append(kinds, "massign")
	}
	if depth > 0 {
		kinds = append(kinds, "ifplain", "while")
		if allowGuard && len(sc.Params) > 0 {
			kinds = append(kinds, "guard")
		}
	}
	switch rapid.SampledFrom(kinds).Draw(t, "wtstmt") {
	case "massign":
		fn := rapid.SampledFrom(sortedKeys(sc.Members)).Draw(t, "wtmtgt")
		return Stmt{Kind: StMemberAssign, Name: fn, Expr: genTypedExpr(t, sc, sc.Members[fn], 2)}
	case "return":
		return Stmt{Kind: StReturn, Expr: genTypedExpr(t, sc, ret, 2)}
	case "ifplain":
		st := Stmt{Kind: StIf, Expr: genTypedExpr(t, sc, GtBool, 2),
			Body: genTypedStmts(t, sc.clone(), depth-1, 3, ret, allowGuard)}
		if rapid.Bool().Draw(t, "wtelse") {
			st.Else = genTypedStmts(t, sc.clone(), depth-1, 3, ret, allowGuard)
		}
		return st
	case "while":
		return Stmt{Kind: StWhile, Expr: genTypedExpr(t, sc, GtBool, 2),
			Body: genTypedStmts(t, sc.clone(), depth-1, 3, ret, allowGuard)}
	case "guard":
		p := rapid.SampledFrom(sc.Params).Draw(t, "wtguardp")
		g := rapid.SampledFrom(BaseTypes).Draw(t, "wtguardt")
		branch := sc.clone()
		branch.narrowed[p] = g
		st := Stmt{Kind: StIf,
			Expr: Expr{Kind: ExFnPred, Pred: PredName(g), Args: []Expr{{Kind: ExVar, Name: p}}},
			Body: genTypedStmts(t, branch, depth-1, 3, ret, false)}
		if rapid.Bool().Draw(t, "wtguardelse") {
			st.Else = genTypedStmts(t, sc.clone(), depth-1, 3, ret, false)
		}
		return st
	default: // "assign"
		vn := fmt.Sprintf("v%d", rapid.IntRange(0, 4).Draw(t, "wtvidx"))
		ty, declared := sc.localTypes[vn]
		if !declared {
			ty = rapid.SampledFrom(BaseTypes).Draw(t, "wtvtype")
			sc.localTypes[vn] = ty
		}
		if !sc.assigned[vn] {
			// build the rhs before marking vn readable: a first assignment on
			// this path must not read itself (use-before-def at runtime)
			e := genTypedExpr(t, sc, ty, 2)
			sc.assigned[vn] = true
			return Stmt{Kind: StAssign, Name: vn, Expr: e}
		}
		return Stmt{Kind: StAssign, Name: vn, Expr: genTypedExpr(t, sc, ty, 2)}
	}
}

func genTypedStmts(t *rapid.T, sc *wtScope, depth, max int, ret GType, allowGuard bool) []Stmt {
	n := rapid.IntRange(0, max).Draw(t, "wtnstmts")
	out := make([]Stmt, 0, n)
	for range n {
		out = append(out, genTypedStmt(t, sc, depth, ret, allowGuard))
	}
	return out
}

func genTypedMethod(t *rapid.T, sig wtSig, members map[string]GType, sigs []wtSig) Method {
	params := make([]Param, sig.nparams)
	paramNames := make([]string, sig.nparams)
	for i := 0; i < sig.nparams; i++ {
		pn := fmt.Sprintf("p%d", i)
		paramNames[i] = pn
		params[i] = Param{Name: pn}
	}
	sc := &wtScope{localTypes: map[string]GType{}, assigned: map[string]bool{},
		narrowed: map[string]GType{}, Members: members, Params: paramNames, sigs: sigs}
	body := genTypedStmts(t, sc, 2, 4, sig.ret, true)
	// guarantee a terminating return of the method's fixed type
	body = append(body, Stmt{Kind: StReturn, Expr: genTypedExpr(t, sc, sig.ret, 2)})
	return Method{Name: sig.Name, Params: params, Body: body, Ret: sig.ret, RetKnown: true}
}

func genProgramWT(t *rapid.T) Program {
	nmem := rapid.IntRange(0, 3).Draw(t, "wtnmem")
	members := map[string]GType{}
	memberList := make([]Member, nmem)
	for i := range nmem {
		n := fmt.Sprintf("f%d", i)
		ty := rapid.SampledFrom(BaseTypes).Draw(t, "wtmemtype")
		members[n] = ty
		memberList[i] = Member{Name: n, Seed: GTypeLit(t, ty)}
	}
	nmeth := rapid.IntRange(1, 3).Draw(t, "wtnmeth")
	sigs := make([]wtSig, nmeth)
	for i := range nmeth {
		sigs[i] = wtSig{Name: fmt.Sprintf("M%d", i),
			nparams: rapid.IntRange(0, 2).Draw(t, "wtnparams"),
			ret:     rapid.SampledFrom(BaseTypes).Draw(t, "wtret")}
	}
	methods := make([]Method, nmeth)
	for i := range nmeth {
		// DAG calls only (earlier methods): a self-call cycle can be valueless
		// (M1() { return .M1() }), which voids the generator's return-type
		// ground truth. general mode still generates recursion.
		methods[i] = genTypedMethod(t, sigs[i], members, sigs[:i])
	}
	if rapid.Bool().Draw(t, "wtctor") {
		optionals := genOptionalMembers(t)
		memberList = append(memberList, optionals...)
		methods = append(methods, genCtor(t, optionals, members))
	}
	return Program{Name: "T", Members: memberList, Methods: methods}
}

var ctorPayloads = []GType{GtNum, GtStr, GtDate, GtObj}

func genOptionalMembers(t *rapid.T) []Member {
	nopt := rapid.IntRange(0, 2).Draw(t, "wtnopt")
	out := make([]Member, nopt)
	for i := range nopt {
		out[i] = Member{Name: fmt.Sprintf("o%d", i), Seed: Lit{Kind: LitFalse},
			Ctor:    rapid.SampledFrom([]CtorMode{CtorAlways, CtorMaybe, CtorNever}).Draw(t, "wtoptmode"),
			Payload: rapid.SampledFrom(ctorPayloads).Draw(t, "wtoptpayload")}
	}
	return out
}

// genCtor builds New: optional-member assignments first (so a filler return
// cannot precede them and break the always-assigned guarantee), then ordinary
// typed statements. no params, no guards, no calls - every assigned value is
// a clean expression, keeping the demotion ground truth decidable.
func genCtor(t *rapid.T, optionals []Member, members map[string]GType) Method {
	sc := &wtScope{localTypes: map[string]GType{}, assigned: map[string]bool{},
		narrowed: map[string]GType{}, Members: members}
	// assigned VALUES come from a member-free scope: a member read may be
	// dirty (param-fed elsewhere), and a possibly-false value tripping C2
	// would void the CtorAlways ground truth. conditions may read members.
	valSc := &wtScope{localTypes: map[string]GType{}, assigned: map[string]bool{},
		narrowed: map[string]GType{}, Members: map[string]GType{}}
	var body []Stmt
	hasMaybe := false
	for _, o := range optionals {
		if o.Ctor == CtorMaybe {
			hasMaybe = true
		}
	}
	// maybe-conditions must stay undecidable to ConstructorExecPass: its
	// feasibility check proves is/isnt/not/and/or over disjoint singletons
	// (e.g. `false isnt true` is always true, making the branch effectively
	// unconditional). a `<` comparison against a fresh local is opaque to it
	// and to the parser's constant folder.
	condT := GtNum
	if hasMaybe {
		condT = rapid.SampledFrom(comparableTypes).Draw(t, "wtctorcondt")
		sc.localTypes["v7"] = condT
		sc.assigned["v7"] = true
		body = append(body, Stmt{Kind: StAssign, Name: "v7",
			Expr: Expr{Kind: ExLit, Lit: GTypeLit(t, condT)}})
	}
	for _, o := range optionals {
		switch o.Ctor {
		case CtorAlways:
			body = append(body, Stmt{Kind: StMemberAssign, Name: o.Name,
				Expr: genTypedExpr(t, valSc, o.Payload, 2)})
		case CtorMaybe:
			cond := Expr{Kind: ExBin, Op: rapid.SampledFrom([]string{"<", "<=", ">", ">="}).Draw(t, "wtctorcmp"),
				Args: []Expr{{Kind: ExVar, Name: "v7"}, {Kind: ExLit, Lit: GTypeLit(t, condT)}}}
			body = append(body, Stmt{Kind: StIf, Expr: cond,
				Body: []Stmt{{Kind: StMemberAssign, Name: o.Name,
					Expr: genTypedExpr(t, valSc, o.Payload, 2)}}})
		}
	}
	filler := rapid.SampledFrom(BaseTypes).Draw(t, "wtctorret")
	body = append(body, genTypedStmts(t, sc, 1, 3, filler, false)...)
	return Method{Name: "New", Body: body}
}

// LitGType recovers the GType of a seeded literal.
func LitGType(l Lit) GType {
	switch l.Kind {
	case LitNum:
		return GtNum
	case LitStr:
		return GtStr
	case LitDate:
		return GtDate
	case LitTrue, LitFalse:
		return GtBool
	default:
		return GtObj
	}
}

// AnnName is the source annotation spelling of a GType (": number" etc).
func AnnName(g GType) string {
	switch g {
	case GtNum:
		return "number"
	case GtStr:
		return "string"
	case GtDate:
		return "date"
	case GtObj:
		return "object"
	default:
		return "boolean"
	}
}
