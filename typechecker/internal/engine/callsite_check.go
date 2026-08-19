// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"fmt"
	"sort"
	"strings"

	"github.com/apmckinlay/gsuneido/compile/ast"
	"github.com/apmckinlay/gsuneido/core"
)

func CallsiteCheckPass(cls *ClassObject, env TypeEnv, pctx *PassCtx) bool {
	annotations := pctx.Annotations
	for name, fn := range cls.SortedMethods {
		ctx := callCheckCtx{env: env, annotations: annotations, method: name}
		for _, stmt := range fn.Body {
			if stmt != nil {
				ctx.walk(stmt, 0)
			}
		}
	}
	return false
}

type callCheckCtx struct {
	env         TypeEnv
	annotations AnnotationSet
	method      string
}

func (c *callCheckCtx) walk(n ast.Node, pos int) {
	if n == nil {
		return
	}
	if ep, ok := n.(*ast.ExprPos); ok {
		newPos := int(ep.Pos)
		if ep.Expr != nil {
			c.walk(ep.Expr, newPos)
		}
		return
	}
	if s, ok := n.(ast.Statement); ok {
		pos = s.Position()
	}
	n.Children(func(child ast.Node) ast.Node {
		c.walk(child, pos)
		return child
	})
	if call, ok := n.(*ast.Call); ok {
		c.checkCall(call, pos)
	}
}

func (c *callCheckCtx) checkCall(call *ast.Call, pos int) {
	d := dispatchCall(call, c.env, c.annotations)
	switch d.Kind {
	case dispatchReceiverMismatch:
		msg := fmt.Sprintf("method %q%s not applicable to receiver of type %v",
			d.Method, onReceiver(call), d.Receiver)
		c.env.Report(c.receiverDiagnostic(call, pos, msg))
		return
	case dispatchUnionMixed:
		msg := fmt.Sprintf("method %q%s not applicable on at least one path of union receiver %v%s",
			d.Method, onReceiver(call), d.Receiver, badMembersNote(d.BadMembers))
		if d.DirtyReceiver {
			// proven only for the concrete arms; the `?` arm may be fine.
			// Got carries the dirty union so confidence sinks accordingly.
			c.env.Report(&Diagnostic{
				Severity: SeverityWarning,
				Method:   c.method,
				Pos:      pos,
				Got:      []DynType{d.Receiver},
				Msg:      msg + " (receiver union contains ?; unknown arms not checked)",
			})
			return
		}
		c.env.Report(c.receiverDiagnostic(call, pos, msg))
		return
	case dispatchNoSuchMethod:
		msg := fmt.Sprintf("no built-in method %q%s and receiver type unknown",
			d.Method, onReceiver(call))
		c.env.Report(&Diagnostic{
			Severity: SeverityWarning,
			Method:   c.method,
			Pos:      pos,
			Msg:      msg,
		})
		return
	case dispatchGuessSingle:
		c.env.Report(&Diagnostic{
			Severity: SeverityWarning,
			Method:   c.method,
			Pos:      pos,
			Msg: fmt.Sprintf("receiver type unknown; assuming built-in %q on %v",
				d.Method, d.Sig.Receiver),
		})
	case dispatchGuessAgree:
		c.env.Report(&Diagnostic{
			Severity: SeverityWarning,
			Method:   c.method,
			Pos:      pos,
			Msg: fmt.Sprintf("receiver type unknown; %q has multiple overloads agreeing on return %v",
				d.Method, d.Returns),
		})
	case dispatchGuessDisagree:
		c.env.Report(&Diagnostic{
			Severity: SeverityWarning,
			Method:   c.method,
			Pos:      pos,
			Msg: fmt.Sprintf("receiver type unknown; %q has ambiguous overloads with differing returns (folded to %v)",
				d.Method, d.Returns),
		})
	}
	if sig := c.env.CallSig(call); sig != nil {
		c.checkArgs(call, sig, pos)
	}
}

// a receiver whose type rests on a guess cannot prove a mismatch, so the finding caps at warning
func (c *callCheckCtx) receiverDiagnostic(call *ast.Call, pos int, msg string) *Diagnostic {
	sev := SeverityError
	if callReceiverGuessed(call, c.method, c.env) {
		sev = SeverityWarning
		msg += " (receiver type guessed from builtin overloads)"
	}
	return &Diagnostic{Severity: sev, Method: c.method, Pos: pos, Msg: msg}
}

func onReceiver(call *ast.Call) string {
	if mem, ok := call.Fn.(*ast.Mem); ok && mem.E != nil {
		return onExpr(mem.E)
	}
	return ""
}

func onExpr(e ast.Expr) string {
	if s := echoExpr(e); s != "" {
		return " on `" + s + "`"
	}
	return ""
}

func echoExpr(e ast.Expr) (s string) {
	defer func() { _ = recover() }()
	return e.Echo()
}

func badMembersNote(bad []DynType) string {
	if len(bad) == 0 {
		return ""
	}
	return "; no overload for " + joinTypes(bad)
}

func joinTypes(ts []DynType) string {
	parts := make([]string, len(ts))
	for i, t := range ts {
		parts[i] = fmt.Sprint(t)
	}
	sort.Strings(parts)
	return strings.Join(parts, ", ")
}

func (c *callCheckCtx) checkArgs(call *ast.Call, sig *Signature, callPos int) {
	if sig.AtParam {
		return
	}
	sigGuessed := c.env.GuessedCall(call) || callReceiverGuessed(call, c.method, c.env)
	named := make(map[string]bool, len(call.Args))
	for i := range call.Args {
		arg := &call.Args[i]
		if name, ok := argName(arg); ok && !isAtArg(arg) {
			named[name] = true
		}
	}
	for i := range call.Args {
		arg := &call.Args[i]
		if isAtArg(arg) {
			return
		}
		c.checkCallArg(call, sig, arg, i, named, sigGuessed, callPos)
	}
}

func (c *callCheckCtx) checkCallArg(call *ast.Call, sig *Signature, arg *ast.Arg, index int,
	named map[string]bool, sigGuessed bool, callPos int) {
	param := matchParam(sig, arg, index)
	if param == nil {
		c.checkDiscardedNamed(arg, call, sigGuessed)
		return
	}
	if arg.Name == nil && named[param.Name] {
		return
	}
	if param.Typ == nil || param.Typ == TUnknown {
		return
	}
	argType := c.env.GetType(arg.E)
	if argType == TUnknown {
		return
	}
	argPos := arg.GetPos()
	if argPos == 0 {
		argPos = callPos
	}
	guessed := sigGuessed || exprGuessed(arg.E, c.method, c.env)
	c.checkOneArg(argType, param, argPos, guessed, call, arg)
}

func (c *callCheckCtx) checkDiscardedNamed(arg *ast.Arg, call *ast.Call, sigGuessed bool) {
	if arg.Name == nil || sigGuessed {
		return
	}
	name, ok := argName(arg)
	if !ok || name == "block" {
		return
	}
	pos := arg.GetPos()
	if pos == 0 {
		pos = nodePos(call)
	}
	c.env.Report(&Diagnostic{
		Severity: SeverityWarning,
		Method:   c.method,
		Pos:      pos,
		Msg: fmt.Sprintf("named argument %q matches no parameter of %s and is discarded at runtime",
			name, calleeName(call)),
	})
}

// ```suneido
// ob.Sort!(block, dirty: true)
//
//	^^^^^  ^^^^^^^^^^^
//	|      named arg - looked up by Name "dirty"
//	positional - binds to Params[0]
//
// ```
// nil when the arg lands outside the declared params.
func matchParam(sig *Signature, arg *ast.Arg, index int) *Param {
	if arg.Name != nil {
		name, ok := argName(arg)
		if !ok {
			return nil // a non-string name binds to no parameter
		}
		for i := range sig.Params {
			if sig.Params[i].Name == name {
				return &sig.Params[i]
			}
		}
		return nil
	}
	if index < len(sig.Params) {
		return &sig.Params[index]
	}
	return nil
}

// ```suneido
// .Foo(@args)     // args-spread - positional reasoning dies
//
//	^^^^^^^^^      // encoded by parser as a named arg with Name "@" (or "@+1")
//
// ```
// see compile/expression.go's argList handling.
func isAtArg(arg *ast.Arg) bool {
	if arg.Name == nil {
		return false
	}
	if s, ok := arg.Name.(core.SuStr); ok {
		return string(s) == "@" || string(s) == "@+1"
	}
	return false
}

// a named argument's name is a Value, not a string - `f(1: x)` parses and gives
// Arg{Name: SuInt(1)}. The runtime matches names against the SuStr param names,
// so a non-string name binds to no parameter. ok=false means there is nothing to
// look up.
func argName(arg *ast.Arg) (string, bool) {
	if arg.Name == nil {
		return "", false
	}
	return arg.Name.ToStr()
}

func (c *callCheckCtx) checkOneArg(argType DynType, param *Param, pos int, guessed bool,
	call *ast.Call, arg *ast.Arg) {
	members, dirty := decomposeForCheck(argType)
	var bad []DynType
	for _, m := range members {
		if !memberFits(m, param.Typ) {
			bad = append(bad, m)
		}
	}
	msg := func(suffix string) string {
		if param.Inferred {
			return fmt.Sprintf("%s to %s has type %v, but %q is passed on to %s which requires %v%s",
				argRef(arg, param), calleeDisplay(call), argType, param.Name, requirementRoot(param), param.Typ, suffix)
		}
		return fmt.Sprintf("argument %q: type %v not assignable to declared %v%s",
			param.Name, argType, param.Typ, suffix)
	}
	confidence := 0.0
	if param.Inferred {
		confidence = 0.75
	}
	switch {
	case len(bad) > 0 && guessed:
		c.env.Report(&Diagnostic{
			Severity: SeverityWarning,
			Method:   c.method,
			Pos:      pos,
			Msg:      msg(" (type guessed from builtin overloads)"),
		})
	case len(bad) > 0:
		c.env.Report(&Diagnostic{
			Severity:   SeverityError,
			Method:     c.method,
			Pos:        pos,
			Confidence: confidence,
			Msg:        msg(""),
		})
	case dirty && !param.Inferred:
		c.env.Report(&Diagnostic{
			Severity: SeverityWarning,
			Method:   c.method,
			Pos:      pos,
			Msg: fmt.Sprintf("argument %q: type %v contains unknown, cannot prove assignable to declared %v",
				param.Name, argType, param.Typ),
		})
	}
}

func requirementRoot(param *Param) string {
	if param.Why != "" {
		return param.Why
	}
	return "a typed callee"
}

func argRef(arg *ast.Arg, param *Param) string {
	if s := argText(arg); s != "" {
		return "argument `" + s + "`"
	}
	return fmt.Sprintf("argument %q", param.Name)
}

func argText(arg *ast.Arg) (s string) {
	defer func() { _ = recover() }()
	if arg.E == nil {
		return ""
	}
	s = arg.E.Echo()
	if len(s) > 40 {
		s = s[:37] + "..."
	}
	return s
}
