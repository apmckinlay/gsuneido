// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"sort"
	"strings"
	"testing"
)

// ---- pipeline drivers --------------------------------------------------

type mapResolver map[string]string

func (m mapResolver) Resolve(name string) (*ClassObject, error) {
	src, ok := m[name]
	if !ok {
		return nil, errResolverMiss(name)
	}
	return NewClassObject(name, ParseClass(src)), nil
}

type errResolverMiss string

func (e errResolverMiss) Error() string { return "no class " + string(e) }

func runPasses(src, name string) (*ClassObject, TypeEnv) {
	return runPassesWith(src, name, nil)
}

func runPassesWith(src, name string, seed func(*TypeEnv)) (*ClassObject, TypeEnv) {
	cls := NewClassObject(name, ParseClass(src))
	env := NewTypeEnv()
	if seed != nil {
		seed(&env)
	}
	RunPipeline(cls, env, NewPassCtx(), nil)
	return cls, env
}

func runPassesWithClasses(src string,
	classes map[string]map[string]DynType) (*ClassObject, TypeEnv) {
	return runPassesWith(src, "T", func(e *TypeEnv) { e.Classes = classes })
}

func runPassesWithClassViews(src string,
	instance, seed map[string]map[string]DynType) (*ClassObject, TypeEnv) {
	return runPassesWith(src, "T", func(e *TypeEnv) {
		e.Classes = instance
		e.PreCtorClasses = seed
	})
}

func runPassesWithStatics(src string,
	classes, members map[string]map[string]DynType,
	bases map[string]string) (*ClassObject, TypeEnv) {
	return runPassesWith(src, "T", func(e *TypeEnv) {
		e.Classes = classes
		e.ClassMembers = members
		e.ClassBases = bases
	})
}

func runPassesWithRefSrcs(src string, refs map[string]string) (*ClassObject, TypeEnv) {
	names := make([]string, 0, len(refs))
	for rname := range refs {
		names = append(names, rname)
	}
	sort.Strings(names) // map order is random; keep the build deterministic
	list := make([]RefSource, 0, len(names))
	for _, rname := range names {
		list = append(list, RefSource{Name: rname, Src: refs[rname]})
	}
	return runPassesWithRefList(src, list)
}

// refs must arrive base-first, like the References a session is given
func runPassesWithRefList(src string, refs []RefSource) (*ClassObject, TypeEnv) {
	regs := BuildReferenceRegistry(refs)
	return runPassesWith(src, "T", func(e *TypeEnv) { regs.Seed(e) })
}

func mustRun(t *testing.T, src string) TypeEnv {
	t.Helper()
	cls, ok := safeParse(src, "T")
	if !ok {
		t.Fatal("source did not parse")
	}
	env, panicMsg := safeRun(cls)
	if panicMsg != "" {
		t.Fatalf("pipeline panicked: %s", panicMsg)
	}
	return env
}

// ---- diagnostics -------------------------------------------------------

func methodDiags(env TypeEnv, method string) []Diagnostic {
	var out []Diagnostic
	if env.Diagnostics == nil {
		return out
	}
	for _, d := range *env.Diagnostics {
		if d.Method == method {
			out = append(out, d)
		}
	}
	return out
}

func diagsByMethod(env TypeEnv) map[string][]Diagnostic {
	out := map[string][]Diagnostic{}
	if env.Diagnostics == nil {
		return out
	}
	for _, d := range *env.Diagnostics {
		out[d.Method] = append(out[d.Method], d)
	}
	return out
}

func condDiags(ds []Diagnostic, sev Severity) []Diagnostic {
	var out []Diagnostic
	for _, d := range ds {
		if d.Severity == sev {
			out = append(out, d)
		}
	}
	return out
}

func countSeverity(ds []Diagnostic, s Severity) int { return len(condDiags(ds, s)) }

func hasDiag(env TypeEnv, method string, sev Severity, substr string) bool {
	for _, d := range methodDiags(env, method) {
		if d.Severity == sev && strings.Contains(d.Msg, substr) {
			return true
		}
	}
	return false
}

func countDiagsWith(env TypeEnv, substr string) int {
	n := 0
	for _, d := range methodDiags(env, "Foo") {
		if strings.Contains(d.Msg, substr) {
			n++
		}
	}
	return n
}

func errOn(env TypeEnv, method, substr string) bool {
	return hasDiag(env, method, SeverityError, substr)
}

func errorCount(env TypeEnv, method string) int {
	return countSeverity(methodDiags(env, method), SeverityError)
}

func hasNotApplicableError(env TypeEnv, method string) bool {
	return errOn(env, method, "not applicable")
}

func hasOperatorError(env TypeEnv, method string) bool {
	return errOn(env, method, "expects number")
}

func hasMemberNotFound(env TypeEnv) bool {
	return errOn(env, "Foo", "no member")
}

// ---- arity / reference-registry helpers --------------------------------

func diagMsgs(src string, seed func(*TypeEnv), keep func(string) bool) []string {
	_, env := runPassesWith(src, "T", seed)
	var out []string
	if env.Diagnostics == nil {
		return out
	}
	for _, d := range *env.Diagnostics {
		if keep(d.Msg) {
			out = append(out, d.Msg)
		}
	}
	return out
}

func isArityMsg(msg string) bool {
	return strings.HasPrefix(msg, "too many arguments") ||
		strings.HasPrefix(msg, "missing argument")
}

func isOverrideWarn(msg string) bool {
	return strings.Contains(msg, "passed both positionally and by name")
}

// refSeed installs a reference registry built by refOf.
func refSeed(classes map[string]map[string]DynType,
	sigs map[string]map[string]*Signature) func(*TypeEnv) {
	return func(e *TypeEnv) {
		e.Classes = classes
		e.ClassMethodSigs = sigs
	}
}

func arityMsgs(src string) []string { return diagMsgs(src, nil, isArityMsg) }

func joined(src string) string { return strings.Join(arityMsgs(src), " | ") }

func arityMsgsRef(src string, classes map[string]map[string]DynType,
	sigs map[string]map[string]*Signature) []string {
	return diagMsgs(src, refSeed(classes, sigs), isArityMsg)
}

func overrideWarns(src string, classes map[string]map[string]DynType,
	sigs map[string]map[string]*Signature) []string {
	return diagMsgs(src, refSeed(classes, sigs), isOverrideWarn)
}

// a full registry, so base chains are visible; refs must arrive base-first
func arityMsgsRefList(src string, refs []RefSource) []string {
	regs := BuildReferenceRegistry(refs)
	return diagMsgs(src, func(e *TypeEnv) { regs.Seed(e) }, isArityMsg)
}

func refOf(name, src string) (map[string]map[string]DynType, map[string]map[string]*Signature) {
	cls, env := runPasses(src, name)
	return map[string]map[string]DynType{name: env.SnapshotReturns()},
		map[string]map[string]*Signature{name: ClassMethodSignatures(cls)}
}

// ---- type predicates ---------------------------------------------------

func unionContains(t DynType, want DynType) bool {
	if t == want {
		return true
	}
	if u, ok := t.(Union); ok {
		return u.Contains(want)
	}
	return false
}

func isUnionOf(t DynType, members ...DynType) bool {
	u, ok := t.(Union)
	if !ok || len(u.Types) != len(members) {
		return false
	}
	for _, m := range members {
		if !u.Contains(m) {
			return false
		}
	}
	return true
}
