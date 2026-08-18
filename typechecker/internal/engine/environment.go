// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"maps"
	"strings"

	"github.com/apmckinlay/gsuneido/compile/ast"
)

type stores struct {
	Nodes map[ast.Node]DynType
	// keyed by *ast.Param because Param does not implement Node
	Params           map[*ast.Param]DynType
	Members          map[string]DynType // seed-inclusive union view; PostCtorMembers is the post-New demoted view
	PostCtorMembers  map[string]DynType // written only by ConstructorExecPass
	Returns          map[string]DynType // written only by ReturnUnionPass
	PreCtorReturns   map[string]DynType // written only by CaptureSeedReturns
	AnnotatedReturns map[string]DynType
	Diagnostics      *[]Diagnostic
	Classes          map[string]map[string]DynType
	PreCtorClasses   map[string]map[string]DynType
	ValidDateCalls   map[ast.Node]bool
	FalseDateCalls   map[ast.Node]bool
	GuessedCalls     map[ast.Node]bool // provenance: this call's type is a guess; check passes cap it at warning
	GuessedVars      map[string]bool
	AssertedMembers  map[string]AssertFact
	Summaries        map[string]ReturnSummary
	ClassSummaries   map[string]map[string]ReturnSummary
	ClassMembers     map[string]map[string]DynType
	ClassBases       map[string]string
	CallSigs         map[ast.Node]*Signature
	ClassMethodSigs  map[string]map[string]*Signature
	ReqConflicts     map[string]bool
}

// the value part holds only the per-class scope; copying a TypeEnv scopes it
// to a class without touching what other copies see
type TypeEnv struct {
	*stores
	ClassName  string
	ClassBase  string
	ThisType   DynType
	Methods    map[string]*ast.Function
	MethodSigs map[string]*Signature
}

func (env TypeEnv) ClassReturn(class, method string) (DynType, bool) {
	if env.Classes == nil {
		return TUnknown, false
	}
	if m, ok := env.Classes[class]; ok {
		if t, ok := m[method]; ok {
			return t, true
		}
	}
	return TUnknown, false
}

func (env TypeEnv) ClassReturnSeed(class, method string) (DynType, bool) {
	if env.PreCtorClasses != nil {
		if m, ok := env.PreCtorClasses[class]; ok {
			if t, ok := m[method]; ok {
				return t, true
			}
		}
	}
	return env.ClassReturn(class, method)
}

func unmangledPrivate(class, member string) string {
	prefix := class + "_"
	if rest, ok := strings.CutPrefix(member, prefix); ok &&
		rest != "" && rest[0] >= 'a' && rest[0] <= 'z' {
		return rest
	}
	return ""
}

func looksPrivatized(member string) bool {
	if member == "" || member[0] < 'A' || member[0] > 'Z' {
		return false
	}
	for i := strings.IndexByte(member, '_'); i > 0 && i+1 < len(member); {
		if member[i+1] >= 'a' && member[i+1] <= 'z' {
			return true
		}
		j := strings.IndexByte(member[i+1:], '_')
		if j < 0 {
			return false
		}
		i += 1 + j
	}
	return false
}

func (env TypeEnv) ClassStaticType(class, member string) DynType {
	for c := class; ; {
		mt, known := env.ClassMembers[c]
		if !known {
			return TUnknown
		}
		if t, ok := mt[member]; ok {
			return t
		}
		if priv := unmangledPrivate(c, member); priv != "" {
			if t, ok := mt[priv]; ok {
				return t
			}
		}
		if t, ok := env.ClassReturnSeed(c, "Getter_"); ok {
			return t
		}
		if t, ok := env.ClassReturnSeed(c, "Getter_"+member); ok {
			return t
		}
		if c = env.ClassBases[c]; c == "" {
			return TUnknown
		}
	}
}

func (env TypeEnv) ClassStaticAccessible(class, member string) (accessible, resolved bool) {
	for c := class; ; {
		mt, known := env.ClassMembers[c]
		if !known {
			return false, false
		}
		if _, ok := mt[member]; ok {
			return true, true
		}
		if priv := unmangledPrivate(c, member); priv != "" {
			if _, ok := mt[priv]; ok {
				return true, true
			}
		}
		if _, ok := env.ClassReturnSeed(c, member); ok {
			return true, true
		}
		if _, ok := env.ClassReturnSeed(c, "Getter_"); ok {
			return true, true // generic getter accepts any name
		}
		if _, ok := env.ClassReturnSeed(c, "Getter_"+member); ok {
			return true, true
		}
		if c = env.ClassBases[c]; c == "" {
			return false, true
		}
	}
}

func NewTypeEnv() TypeEnv {
	diags := make([]Diagnostic, 0, 4)
	return TypeEnv{stores: &stores{
		Nodes:            make(map[ast.Node]DynType, 10),
		Params:           make(map[*ast.Param]DynType, 10),
		Members:          make(map[string]DynType, 10),
		PostCtorMembers:  make(map[string]DynType, 4),
		Returns:          make(map[string]DynType, 10),
		PreCtorReturns:   make(map[string]DynType, 4),
		AnnotatedReturns: make(map[string]DynType, 4),
		Diagnostics:      &diags,
		ValidDateCalls:   make(map[ast.Node]bool, 4),
		FalseDateCalls:   make(map[ast.Node]bool, 4),
		GuessedCalls:     make(map[ast.Node]bool, 8),
		GuessedVars:      make(map[string]bool, 8),
		AssertedMembers:  make(map[string]AssertFact, 4),
		Summaries:        make(map[string]ReturnSummary, 4),
		CallSigs:         make(map[ast.Node]*Signature, 8),
		ReqConflicts:     make(map[string]bool, 2),
	}}
}

func (env TypeEnv) WithClass(cls *ClassObject, sigs map[string]*Signature) TypeEnv {
	env.ClassName = cls.Name
	env.ClassBase = cls.Base
	env.ThisType = extensionThisType(cls.Name)
	env.Methods = cls.Methods
	env.MethodSigs = sigs
	return env
}

func (env TypeEnv) AnnotationView(keep map[ast.Node]struct{}) TypeEnv {
	nodes := make(map[ast.Node]DynType, len(keep))
	for n, ty := range env.Nodes {
		if _, ok := keep[n]; ok {
			nodes[n] = ty
		}
	}
	return TypeEnv{stores: &stores{
		Nodes:   nodes,
		Params:  env.Params,
		Members: env.Members,
		Returns: env.Returns,
	}}
}

func (env TypeEnv) SetCallSig(n ast.Node, sig *Signature) {
	if env.CallSigs == nil {
		return
	}
	if sig == nil {
		delete(env.CallSigs, n)
		return
	}
	env.CallSigs[n] = sig
}

func (env TypeEnv) CallSig(n ast.Node) *Signature {
	return env.CallSigs[n]
}

func (env TypeEnv) MarkValidDate(n ast.Node) {
	env.ValidDateCalls[n] = true
}

func (env TypeEnv) DateProvenValid(n ast.Node) bool {
	return env.ValidDateCalls[n]
}

func (env TypeEnv) MarkFalseDate(n ast.Node) {
	env.FalseDateCalls[n] = true
}

func (env TypeEnv) DateProvenFalse(n ast.Node) bool {
	return env.FalseDateCalls[n]
}

func (env TypeEnv) SetGuessedCall(n ast.Node, guessed bool) {
	if env.GuessedCalls == nil {
		return
	}
	if guessed {
		env.GuessedCalls[n] = true
	} else {
		delete(env.GuessedCalls, n)
	}
}

func (env TypeEnv) GuessedCall(n ast.Node) bool {
	return env.GuessedCalls[n]
}

func (env TypeEnv) SetGuessedVar(method, name string) {
	if env.GuessedVars != nil {
		env.GuessedVars[method+"\x00"+name] = true
	}
}

func (env TypeEnv) GuessedVar(method, name string) bool {
	return env.GuessedVars[method+"\x00"+name]
}

func (env TypeEnv) Report(d *Diagnostic) {
	if env.Diagnostics == nil {
		return
	}
	*env.Diagnostics = append(*env.Diagnostics, *d)
}

func (env TypeEnv) GetType(node ast.Node) DynType {
	if ty, ok := env.Nodes[node]; ok {
		return ty
	}
	return TUnknown
}

func (env TypeEnv) SetType(node ast.Node, ty DynType) {
	env.Nodes[node] = ty
}

func (env TypeEnv) SetParam(p *ast.Param, ty DynType) {
	env.Params[p] = ty
}

func (env TypeEnv) LookupMember(name string) (DynType, bool) {
	ty, ok := env.Members[name]
	return ty, ok
}

// overwrites any existing value
func (env TypeEnv) SeedMember(name string, ty DynType) {
	env.Members[name] = ty
}

func (env TypeEnv) UnionMember(name string, ty DynType) {
	if existing, ok := env.Members[name]; ok {
		ty = U(existing, ty)
	}
	if u, ok := ty.(Union); ok {
		env.Members[name] = u.Fold()
	} else {
		env.Members[name] = ty
	}
}

func (env TypeEnv) LookupReturn(method string) DynType {
	if ty, ok := env.Returns[method]; ok {
		return ty
	}
	return TUnknown
}

func (env TypeEnv) PublishReturn(method string, ty DynType) {
	env.Returns[method] = ty
}

// independent copy
func (env TypeEnv) SnapshotReturns() map[string]DynType {
	return maps.Clone(env.Returns)
}

// independent copy
func (env TypeEnv) SnapshotSummaries() map[string]ReturnSummary {
	return maps.Clone(env.Summaries)
}

func (env TypeEnv) SnapshotPreCtorReturns() map[string]DynType {
	return maps.Clone(env.PreCtorReturns)
}

func (env TypeEnv) CaptureSeedReturns() {
	maps.Copy(env.PreCtorReturns, env.Returns)
}

func (env TypeEnv) ClassKnown(class string) bool {
	if env.Classes != nil {
		if _, ok := env.Classes[class]; ok {
			return true
		}
	}
	if env.PreCtorClasses != nil {
		if _, ok := env.PreCtorClasses[class]; ok {
			return true
		}
	}
	return false
}
