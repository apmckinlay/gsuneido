// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package typechecker

import (
	"fmt"
	"maps"
	"sync"

	"github.com/apmckinlay/gsuneido/typechecker/annotations"
	"github.com/apmckinlay/gsuneido/typechecker/diagnostics"
	"github.com/apmckinlay/gsuneido/typechecker/internal/annotate"
	"github.com/apmckinlay/gsuneido/typechecker/internal/engine"
	"github.com/apmckinlay/gsuneido/typechecker/typealgebra"
)

// SourceEntry is one named Suneido source.
type SourceEntry struct {
	Name string
	Src  string
}

type Request struct {
	Method     string // "TypeInfer" or "TypeAnnotate"
	Arguments  []SourceEntry
	References []SourceEntry
	Config     map[string]string
}

type TypeInfo struct {
	Methods map[string]map[string]string
	Members map[string]string
}

type ResultDiagnostic struct {
	Class  string
	Method string
	Pos    int
	Line   int
	Col    int
	Msg    string
	Flag   diagnostics.Flag
}

type DiagnosticSet struct {
	Errors   []ResultDiagnostic
	Warnings []ResultDiagnostic
}

type Result struct {
	Method      string
	Results     []any
	Diagnostics DiagnosticSet
}

// registerMu serializes table rebuilds; readers of the table itself go
// through engine.Annotations, which is lock-free. It also guards the two
// layers the published table is rebuilt from.
var registerMu sync.Mutex
var builtinRegs []annotations.Registration
var suneidoRegs []annotations.Registration

// SetBuiltinSignatures publishes the signatures collected from the exe's own
// builtin declarations as the base layer of the table. builtin pushes them
// once at startup - it imports typechecker, not the reverse. Entries from
// RegisterSignatures are layered on top, so Suneido can override an exe
// signature without losing the rest.
func SetBuiltinSignatures(regs []annotations.Registration) error {
	registerMu.Lock()
	defer registerMu.Unlock()
	set, err := annotations.LoadLayered(regs, suneidoRegs)
	if err != nil {
		return err
	}
	builtinRegs = regs
	engine.SetAnnotations(set)
	return nil
}

// RegisterSignatures replaces the Suneido-supplied layer of the signature
// table: the exe's builtin declarations stay as the base and these entries
// go on top, one with the same name and receiver overriding the exe's.
// Replace-not-append makes the call idempotent, so re-registering after a
// library reload cannot leave stale entries behind. Any invalid entry fails
// the whole call, leaving the previous table in place.
func RegisterSignatures(regs []annotations.Registration) (int, error) {
	registerMu.Lock()
	defer registerMu.Unlock()
	set, err := annotations.LoadLayered(builtinRegs, regs)
	if err != nil {
		return 0, err
	}
	suneidoRegs = regs
	engine.SetAnnotations(set)
	return len(regs), nil
}

func buildConfig(raw map[string]string) (diagnostics.Config, error) {
	cfg := diagnostics.DefaultConfig()
	if v, ok := raw["strictStringConcat"]; ok {
		l, err := diagnostics.ParseLevel(v)
		if err != nil {
			return cfg, fmt.Errorf("config.strictStringConcat: %w", err)
		}
		cfg.StrictStringConcat = l
	}
	if v, ok := raw["strictCrossTypeCompares"]; ok {
		l, err := diagnostics.ParseLevel(v)
		if err != nil {
			return cfg, fmt.Errorf("config.strictCrossTypeCompares: %w", err)
		}
		cfg.StrictCrossTypeCompares = l
	}
	return cfg, nil
}

func stringifyTypes(in map[string]typealgebra.DynType) map[string]string {
	out := make(map[string]string, len(in))
	for name, ty := range in {
		out[name] = ty.String()
	}
	return out
}

func stringifyVarTypes(in map[string]map[string]typealgebra.DynType) map[string]map[string]string {
	out := make(map[string]map[string]string, len(in))
	for method, vars := range in {
		out[method] = stringifyTypes(vars)
	}
	return out
}

// converts a byte offset into 1-based line/column for the result.
func offsetToLineCol(src string, off int) (line, col int) {
	if off < 0 {
		off = 0
	}
	if off > len(src) {
		off = len(src)
	}
	line, col = 1, 1
	for i := 0; i < off; i++ {
		if src[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return
}

func parseArgument(a SourceEntry) (co *engine.ClassObject, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%s: %v", a.Name, e)
		}
	}()
	return engine.NewClassObject(a.Name, engine.ParseClass(a.Src)), nil
}

// a bug in a pass must not escape into the interpreter as a Go panic - report
// it against the class and carry on with the next argument
func checkArgument(a SourceEntry, method string, cls *engine.ClassObject,
	env engine.TypeEnv, pctx *engine.PassCtx,
	parentReturns map[string]typealgebra.DynType) (result any, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%s: %v", a.Name, e)
		}
	}()
	engine.RunPipeline(cls, env, pctx, parentReturns)
	if method == "TypeAnnotate" {
		return annotate.AnnotateClass(a.Src, env, cls), nil
	}
	return TypeInfo{
		Methods: stringifyVarTypes(engine.MethodVarTypes(cls, env)),
		Members: stringifyTypes(engine.MemberTypes(cls, env)),
	}, nil
}

// the class still needs a result even when checking it failed
func failedResult(method, src string) any {
	if method == "TypeAnnotate" {
		return src
	}
	return TypeInfo{
		Methods: map[string]map[string]string{},
		Members: map[string]string{},
	}
}

func Process(req Request) (Result, error) {
	if req.Method != "TypeInfer" && req.Method != "TypeAnnotate" {
		return Result{}, fmt.Errorf("unknown method: %q (expected TypeInfer or TypeAnnotate)", req.Method)
	}
	if !engine.SignaturesRegistered() {
		return Result{}, fmt.Errorf(
			"no signatures registered - call TypeChecker.RegisterSignatures first")
	}

	cfg, err := buildConfig(req.Config)
	if err != nil {
		return Result{}, err
	}

	confFilter, err := parseConfidenceFilter(req.Config)
	if err != nil {
		return Result{}, err
	}

	refs := make([]engine.RefSource, len(req.References))
	for i, r := range req.References {
		refs[i] = engine.RefSource{Name: r.Name, Src: r.Src}
	}
	regs := engine.BuildReferenceRegistry(refs)

	parsed := make([]*engine.ClassObject, len(req.Arguments))
	for i, a := range req.Arguments {
		if parsed[i], err = parseArgument(a); err != nil {
			return Result{}, err
		}
	}

	env := engine.NewTypeEnv()
	regs.Seed(&env)
	pctx := engine.NewPassCtx()
	parentReturns := map[string]typealgebra.DynType{}
	results := make([]any, len(parsed))
	var collected []rankedDiag
	for i, c := range parsed {
		class := req.Arguments[i].Name
		result, err := checkArgument(req.Arguments[i], req.Method, c, env, pctx,
			parentReturns)
		if err != nil {
			result = failedResult(req.Method, req.Arguments[i].Src)
			collected = append(collected, rankedDiag{
				severity:   diagnostics.SeverityError,
				confidence: 1,
				entry: ResultDiagnostic{
					Class: class,
					Msg:   "typechecker internal error: " + err.Error(),
				},
			})
		}
		results[i] = result
		if env.Diagnostics != nil {
			filtered := diagnostics.FilterDiagnostics(*env.Diagnostics, cfg)
			for _, d := range filtered {
				line, col := offsetToLineCol(req.Arguments[i].Src, d.Pos)
				collected = append(collected, rankedDiag{
					severity:   d.Severity,
					confidence: diagnostics.ScoreConfidence(&d),
					entry: ResultDiagnostic{
						Class:  class,
						Method: d.Method,
						Pos:    d.Pos,
						Line:   line,
						Col:    col,
						Msg:    d.Msg,
						Flag:   d.Flag,
					},
				})
			}
			*env.Diagnostics = (*env.Diagnostics)[:0]
		}
		parentReturns = maps.Clone(env.Returns)
	}

	return Result{
		Method:      req.Method,
		Results:     results,
		Diagnostics: rankDiagnostics(collected, confFilter),
	}, nil
}
