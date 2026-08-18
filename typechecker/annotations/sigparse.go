// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package annotations

import (
	"fmt"
	"strings"

	"github.com/apmckinlay/gsuneido/compile"
	"github.com/apmckinlay/gsuneido/compile/ast"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"

	"github.com/apmckinlay/gsuneido/typechecker/typealgebra"
)

func parseSig(sig string) (Signature, error) {
	// like builtin params(), nil defaults are parsed as a placeholder string -
	// default values are discarded anyway, only HasDefault survives. Sigs
	// registered from Suneido may already carry the quoted form.
	sig = strings.ReplaceAll(sig, "'nil'", "nil")
	sig = strings.ReplaceAll(sig, "nil", "'nil'")
	src := "function " + sig + " {}"
	var result Signature
	var parseErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				parseErr = fmt.Errorf("parse error: %v", r)
			}
		}()
		p := compile.AstParser(src)
		val := p.Const()
		if p.Token != tok.Eof {
			parseErr = fmt.Errorf("trailing tokens after signature")
			return
		}
		fn, ok := val.(*ast.Function)
		if !ok {
			parseErr = fmt.Errorf("AstParser did not return *ast.Function (got %T)", val)
			return
		}
		s, err := SignatureFromAst(fn)
		if err != nil {
			parseErr = err
			return
		}
		result = s
	}()
	return result, parseErr
}

func SignatureFromAst(fn *ast.Function) (Signature, error) {
	var sig Signature
	if fn.ReturnAnnotation != "" {
		ret, err := typealgebra.ParseAnnotation(fn.ReturnAnnotation)
		if err != nil {
			return Signature{}, fmt.Errorf("return: %w", err)
		}
		sig.Returns = ret
	}
	if len(fn.Params) == 1 && strings.HasPrefix(fn.Params[0].Name.Name, "@") {
		sig.AtParam = true
		return sig, nil
	}
	for _, p := range fn.Params {
		typ := typealgebra.DynType(typealgebra.TUnknown)
		if p.Annotations != "" {
			t, err := typealgebra.ParseAnnotation(p.Annotations)
			if err != nil {
				return Signature{}, fmt.Errorf("param %s: %w", p.Name.Name, err)
			}
			typ = t
		}
		sig.Params = append(sig.Params, Param{
			Name:       paramBaseName(p.Name.Name),
			Typ:        typ,
			HasDefault: p.DefVal != nil || paramIsDynamic(p.Name.Name),
		})
	}
	return sig, nil
}

// ```suneido
// Foo(.Name) { }
// .Foo(name: "x")
//
//	^^^^^ callers use the un-capitalized name for a .Capitalized
//	      member-init param (mirrors gsuneido compile.param)
//
// ```
func paramBaseName(s string) string {
	dot := false
	if len(s) > 0 && s[0] == '.' {
		dot = true
		s = s[1:]
	}
	if len(s) > 0 && s[0] == '_' {
		s = s[1:]
	}
	if dot && len(s) > 0 && s[0] >= 'A' && s[0] <= 'Z' {
		s = strings.ToLower(s[:1]) + s[1:]
	}
	return s
}

func paramIsDynamic(s string) bool {
	if len(s) > 0 && s[0] == '.' {
		s = s[1:]
	}
	return len(s) > 0 && s[0] == '_'
}
