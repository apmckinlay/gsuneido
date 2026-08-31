// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"github.com/apmckinlay/gsuneido/compile/ast"
)

func CallsiteResolutionPass(cls *ClassObject, env TypeEnv, pctx *PassCtx) bool {
	annotations := pctx.Annotations
	for _, fn := range cls.SortedMethods {
		for _, stmt := range fn.Body {
			if stmt != nil {
				annotateCallNodes(stmt, env, annotations)
			}
		}
	}
	return false
}

func annotateCallNodes(n ast.Node, env TypeEnv, annotations AnnotationSet) {
	if n == nil {
		return
	}
	if ep, ok := n.(*ast.ExprPos); ok {
		if ep.Expr != nil {
			annotateCallNodes(ep.Expr, env, annotations)
			if t := env.GetType(ep.Expr); t != TUnknown {
				env.SetType(ep, t)
			}
		}
		return
	}
	n.Children(func(c ast.Node) ast.Node {
		annotateCallNodes(c, env, annotations)
		return c
	})
	call, ok := n.(*ast.Call)
	if !ok {
		return
	}
	d := dispatchCall(call, env, annotations)
	env.SetGuessedCall(call, d.Guessed)
	env.SetCallSig(call, arityCalleeSig(call, &d, env))
	if d.Returns != TUnknown {
		env.SetType(call, d.Returns)
	}
}

func matchesReceiver(got, declared DynType) bool {
	if dynEqual(got, declared) {
		return true
	}
	if subtypeOf(got, declared) {
		return true
	}
	if declU, ok := declared.(Union); ok {
		for _, t := range declU.Types {
			if matchesReceiver(got, t) {
				return true
			}
		}
		return false
	}
	if gotU, ok := got.(Union); ok {
		if gotU.IsDirty {
			return false
		}
		for _, t := range gotU.Types {
			if !matchesReceiver(t, declared) {
				return false
			}
		}
		return len(gotU.Types) > 0
	}
	return false
}

var typeConstructorReturns = map[string]DynType{
	"Number":  TNumber,
	"String":  TString,
	"Object":  TObject,
	"Record":  TObject,
	"Display": TString,
	"Type":    TString,
	"Date":    Union{Types: []DynType{TDate, TFalse}}.Fold(),
}

// resolution consumes Returns; check turns these into diagnostics.
type dispatchKind int

const (
	dispatchSilent           dispatchKind = iota // nothing to report (this.X, unannotated, etc.)
	dispatchMatched                              // definite receiver matched a sig
	dispatchGuessSingle                          // unknown receiver, sole candidate
	dispatchGuessAgree                           // unknown receiver, multiple candidates with same return
	dispatchGuessDisagree                        // unknown receiver, multiple candidates with differing returns
	dispatchReceiverMismatch                     // definite receiver, sig exists for name but none accepts the receiver
	dispatchNoSuchMethod                         // unknown receiver, no sig exists with this name
	dispatchUnionMixed                           // union receiver, at least one element produced a mismatch
)

type callDispatch struct {
	Kind       dispatchKind
	Returns    DynType
	Sig        *Signature
	Method     string // the method name dispatched (for diagnostics)
	Receiver   DynType
	BadMembers []DynType
	Guessed    bool
	// the union receiver carries `?`: the mismatch is proven only for the
	// concrete arms, so the check caps the finding at warning
	DirtyReceiver bool
}

// shared by resolution (consumes Returns) and check (emits diagnostics).
func dispatchCall(call *ast.Call, env TypeEnv, annotations AnnotationSet) callDispatch {
	switch fn := call.Fn.(type) {
	case *ast.Ident:
		if fn.Name == "Date" {
			if len(call.Args) == 0 || env.DateProvenValid(call) {
				return callDispatch{Kind: dispatchSilent, Returns: TDate, Method: "Date"}
			}
			if env.DateProvenFalse(call) {
				return callDispatch{Kind: dispatchSilent, Returns: TFalse, Method: "Date"}
			}
		}
		return dispatchFreeCall(fn.Name, call.Args, env, annotations)
	case *ast.Mem:
		if className, ok := newConstructionClass(fn); ok {
			if className == "this" && env.ClassName != "" {
				className = env.ClassName
			}
			return callDispatch{
				Kind:    dispatchSilent,
				Returns: Instance{Class: className},
				Method:  "*new*",
			}
		}
		return dispatchMemCall(fn, call.Args, env, annotations)
	}
	return callDispatch{Kind: dispatchSilent, Returns: TUnknown}
}

func newConstructionClass(mem *ast.Mem) (string, bool) {
	if name, ok := memberName(mem.M); !ok || name != "*new*" {
		return "", false
	}
	if id, ok := mem.E.(*ast.Ident); ok {
		return id.Name, true
	}
	return "", false
}

// ```suneido
// SomeBuiltin(arg)
// ^^^^^^^^^^^      bare identifier - no receiver constraint
// ```
func dispatchFreeCall(name string, args []ast.Arg, env TypeEnv, annotations AnnotationSet) callDispatch {
	if t, ok := typeConstructorReturns[name]; ok {
		return callDispatch{Kind: dispatchSilent, Returns: t, Method: name}
	}
	if env.ClassKnown(name) {
		owner, resolved := callClassOwner(env, name)
		switch {
		case owner == name: // defined here, so its return is about this class
			if r, ok := summaryReturn(env, name, "CallClass", args); ok {
				return callDispatch{Kind: dispatchSilent, Returns: r, Method: name}
			}
			t, _ := env.ClassReturnSeed(name, "CallClass")
			return callDispatch{Kind: dispatchSilent, Returns: t, Method: name}
		case owner != "":
			// inherited: a base CallClass is written for every subclass, so
			// its return is a guess here - checks warn instead of erroring
			if r, ok := summaryReturn(env, owner, "CallClass", args); ok {
				return callDispatch{Kind: dispatchSilent, Returns: r,
					Method: name, Guessed: true}
			}
			if t, _ := env.ClassReturnSeed(owner, "CallClass"); t != TUnknown {
				return callDispatch{Kind: dispatchSilent, Returns: t,
					Method: name, Guessed: true}
			}
			// it told us nothing: the idiomatic base builds `new this`
			// (Singleton), so the instance guess is the better answer
			return callDispatch{Kind: dispatchSilent,
				Returns: Instance{Class: name}, Method: name, Guessed: true}
		}
		// no CallClass in the chain: the call constructs. When a base was
		// missing we cannot prove that, so flag it as a guess.
		return callDispatch{Kind: dispatchSilent, Returns: Instance{Class: name},
			Method: name, Guessed: !resolved}
	}

	for i, sig := range annotations[name] {
		if sig.Receiver == nil {
			s := &annotations[name][i]
			return callDispatch{
				Kind:    dispatchMatched,
				Returns: blockFormReturn(s, args),
				Sig:     s,
				Method:  name,
			}
		}
	}
	return callDispatch{Kind: dispatchSilent, Returns: TUnknown, Method: name}
}

const maxBaseChain = 32 // guards against a cyclic base chain

// bare Foo(...) dispatches to CallClass, which is inherited - SuClass.Call
// looks it up through the base chain and only falls back to creating an
// instance when no class in the chain defines one.
// owner is the class that supplies CallClass, or "" if the whole chain was
// walked without finding one. resolved is false when the walk hits a class we
// have no reference for, leaving both answers unproven.
func callClassOwner(env TypeEnv, class string) (owner string, resolved bool) {
	for c, i := class, 0; i < maxBaseChain; i++ {
		if _, found := env.ClassReturnSeed(c, "CallClass"); found {
			return c, true
		}
		base, hasBase := env.ClassBases[c]
		if !hasBase || base == "" {
			return "", true
		}
		if !env.ClassKnown(base) {
			return "", false
		}
		c = base
	}
	return "", false
}

// ```suneido
// this.Foo()          // -> env.LookupReturn(Foo), else builtin fallback
// ^^^^^^^^
// Date.Begin()        // -> static lookup "Date.Begin"
// ^^^^^^^^^^
// x.Replace("a", "b") // -> dispatch on x's inferred type (matched or guessed)
// ^^^^^^^^^
// ```
func dispatchMemCall(mem *ast.Mem, args []ast.Arg, env TypeEnv, annotations AnnotationSet) callDispatch {
	name, ok := memberName(mem.M)
	if !ok {
		return callDispatch{Kind: dispatchSilent, Returns: TUnknown}
	}

	if id, ok := mem.E.(*ast.Ident); ok {
		switch {
		case id.Name == "this":
			return dispatchThisCall(name, args, env, annotations)
		case id.Name == "super":
			return callDispatch{Kind: dispatchSilent, Returns: TUnknown, Method: name}
		case isGlobalIdent(id.Name):
			return dispatchGlobalStatic(id.Name, name, args, env, annotations)
		}
	}

	receiverType := env.GetType(mem.E)
	if inst, ok := receiverType.(Instance); ok {
		return dispatchInstance(inst, name, env)
	}
	return dispatchOnReceiver(receiverType, name, env, annotations)
}

func dispatchThisCall(name string, args []ast.Arg, env TypeEnv, annotations AnnotationSet) callDispatch {
	if fn, ok := env.Methods[name]; ok {
		if ct, ok := contextReturn(fn, args, env); ok {
			return callDispatch{Kind: dispatchSilent, Returns: ct, Method: name}
		}
	}
	if t := env.LookupReturn(name); t != TUnknown {
		return callDispatch{Kind: dispatchSilent, Returns: t, Method: name}
	}
	if env.ThisType != nil {
		return dispatchOnReceiver(env.ThisType, name, env, annotations)
	}
	d := scanByName(methodSigs(annotations[name]), name)
	d.Kind = dispatchSilent
	return d
}

func dispatchGlobalStatic(class, name string, args []ast.Arg, env TypeEnv, annotations AnnotationSet) callDispatch {
	if sigs := annotations[class+"."+name]; len(sigs) > 0 {
		return callDispatch{
			Kind:    dispatchSilent,
			Returns: sigs[0].Returns,
			Sig:     &sigs[0],
			Method:  class + "." + name,
		}
	}
	if r, ok := summaryReturn(env, class, name, args); ok {
		return callDispatch{Kind: dispatchSilent, Returns: r, Method: name}
	}
	if t, found := env.ClassReturnSeed(class, name); found {
		return callDispatch{Kind: dispatchSilent, Returns: t, Method: name}
	}
	return callDispatch{Kind: dispatchSilent, Returns: TUnknown, Method: name}
}

func dispatchInstance(inst Instance, name string, env TypeEnv) callDispatch {
	if inst.Class == env.ClassName {
		if t := env.LookupReturn(name); t != TUnknown {
			return callDispatch{Kind: dispatchSilent, Returns: t, Method: name, Receiver: inst}
		}
	}
	if t, found := env.ClassReturn(inst.Class, name); found {
		return callDispatch{Kind: dispatchSilent, Returns: t, Method: name, Receiver: inst}
	}
	return callDispatch{Kind: dispatchSilent, Returns: TUnknown, Method: name, Receiver: inst}
}

func dispatchOnReceiver(receiverType DynType, name string, env TypeEnv, annotations AnnotationSet) callDispatch {
	if u, ok := receiverType.(Union); ok {
		return dispatchUnionReceiver(u, name, env, annotations)
	}
	if inst, ok := receiverType.(Instance); ok {
		return dispatchInstance(inst, name, env)
	}

	sigs := methodSigs(annotations[name])

	if receiverType == TUnknown {
		return scanByName(sigs, name)
	}

	for i, sig := range sigs {
		if matchesReceiver(receiverType, sig.Receiver) {
			return callDispatch{
				Kind:     dispatchMatched,
				Returns:  sig.Returns,
				Sig:      &sigs[i],
				Method:   name,
				Receiver: receiverType,
			}
		}
	}

	if len(sigs) > 0 {
		return callDispatch{
			Kind:     dispatchReceiverMismatch,
			Returns:  TUnknown,
			Method:   name,
			Receiver: receiverType,
		}
	}
	return callDispatch{Kind: dispatchSilent, Returns: TUnknown, Method: name, Receiver: receiverType}
}

func dispatchUnionReceiver(u Union, name string, env TypeEnv, annotations AnnotationSet) callDispatch {
	var returns []DynType
	var firstSig *Signature
	var badMembers []DynType
	sameSig := true
	mixed := false
	for _, elt := range u.Types {
		sub := dispatchOnReceiver(elt, name, env, annotations)
		if sub.Kind == dispatchReceiverMismatch || sub.Kind == dispatchNoSuchMethod {
			mixed = true
			badMembers = append(badMembers, elt)
			continue
		}
		if sub.Returns == TUnknown {
			return callDispatch{Kind: dispatchSilent, Returns: TUnknown, Method: name, Receiver: u}
		}
		returns = append(returns, sub.Returns)
		if firstSig == nil {
			firstSig = sub.Sig
		} else if sub.Sig != firstSig {
			sameSig = false
		}
	}
	if mixed {
		if u.IsDirty {
			// dirt only adds possible arms - it never removes a concrete bad
			// one - so with a concrete good arm alongside, the mismatch is
			// still worth a warning. All-bad concrete arms more likely mean
			// broken inference upstream; those stay silent.
			if len(returns) == 0 {
				return callDispatch{Kind: dispatchSilent, Returns: TUnknown, Method: name, Receiver: u}
			}
			return callDispatch{Kind: dispatchUnionMixed, Returns: TUnknown, Method: name,
				Receiver: u, BadMembers: badMembers, DirtyReceiver: true}
		}
		return callDispatch{Kind: dispatchUnionMixed, Returns: TUnknown, Method: name, Receiver: u, BadMembers: badMembers}
	}
	folded := foldReturns(returns)
	if u.IsDirty {
		return callDispatch{Kind: dispatchSilent, Returns: folded, Method: name, Receiver: u}
	}
	var chosen *Signature
	if sameSig {
		chosen = firstSig
	}
	return callDispatch{
		Kind:     dispatchMatched,
		Returns:  folded,
		Sig:      chosen,
		Method:   name,
		Receiver: u,
	}
}

func methodSigs(sigs []Signature) []Signature {
	for i := range sigs {
		if sigs[i].Receiver == nil {
			out := make([]Signature, 0, len(sigs)-1)
			for j := range sigs {
				if sigs[j].Receiver != nil {
					out = append(out, sigs[j])
				}
			}
			return out
		}
	}
	return sigs
}

// receiver unknown: guess by scanning builtin method names. results in GuessedCalls, so checks warn, never error
// GuessTaintPass carries the provenance through locals and chained receivers
func scanByName(sigs []Signature, name string) callDispatch {
	switch len(sigs) {
	case 0:
		return callDispatch{Kind: dispatchNoSuchMethod, Returns: TUnknown, Method: name}
	case 1:
		return callDispatch{
			Kind:    dispatchGuessSingle,
			Returns: sigs[0].Returns,
			Sig:     &sigs[0],
			Method:  name,
			Guessed: true,
		}
	default:
		if allReturnsAgree(sigs) {
			return callDispatch{
				Kind:    dispatchGuessAgree,
				Returns: sigs[0].Returns,
				Method:  name,
				Guessed: true,
			}
		}
		rs := make([]DynType, 0, len(sigs))
		for _, s := range sigs {
			rs = append(rs, s.Returns)
		}
		return callDispatch{
			Kind:    dispatchGuessDisagree,
			Returns: markDirty(foldReturns(rs)),
			Method:  name,
			Guessed: true,
		}
	}
}

func allReturnsAgree(sigs []Signature) bool {
	if len(sigs) == 0 {
		return true
	}
	first := sigs[0].Returns
	for _, s := range sigs[1:] {
		if !dynEqual(s.Returns, first) {
			return false
		}
	}
	return true
}

func foldReturns(rs []DynType) DynType {
	if len(rs) == 0 {
		return TUnknown
	}
	out := rs[0]
	for _, r := range rs[1:] {
		out = U(out, r)
	}
	return out
}
