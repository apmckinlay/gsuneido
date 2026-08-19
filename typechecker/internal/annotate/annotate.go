// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package annotate

import (
	"fmt"
	"sort"
	"strings"

	"github.com/apmckinlay/gsuneido/compile/ast"
	"github.com/apmckinlay/gsuneido/compile/lexer"
	tok "github.com/apmckinlay/gsuneido/compile/tokens"
	"github.com/apmckinlay/gsuneido/typechecker/internal/engine"
	"github.com/apmckinlay/gsuneido/typechecker/typealgebra"
)

func targetNodeSet(cls *engine.ClassObject) map[ast.Node]struct{} {
	set := map[ast.Node]struct{}{}
	var visit func(n ast.Node)
	visit = func(n ast.Node) {
		if n == nil {
			return
		}
		set[n] = struct{}{}
		if ep, ok := n.(*ast.ExprPos); ok {
			visit(ep.Expr)
			return
		}
		if fi, ok := n.(*ast.ForIn); ok {
			if fi.Var.Name != "" {
				set[&fi.Var] = struct{}{}
			}
			if fi.Var2.Name != "" {
				set[&fi.Var2] = struct{}{}
			}
		}
		switch x := n.(type) {
		case ast.Statement:
			x.Children(func(c ast.Node) ast.Node { visit(c); return c })
		case ast.Expr:
			x.Children(func(c ast.Node) ast.Node { visit(c); return c })
		}
	}
	for _, fn := range cls.SortedMethods {
		for _, stmt := range fn.Body {
			visit(stmt)
		}
	}
	return set
}

// alreadyAnnotated reports whether src already carries this exact annotation at
// pos, so re-running the annotator over its own output does not stack copies
func alreadyAnnotated(src string, pos int, text string) bool {
	if pos < 0 || pos > len(src) {
		return false
	}
	return strings.HasPrefix(strings.TrimLeft(src[pos:], " \t"), strings.TrimSpace(text))
}

func exprEnd(n ast.Node) (int, bool) {
	switch x := n.(type) {
	case *ast.Ident:
		if !x.Implicit {
			return x.GetEnd(), true
		}
	case *ast.Call:
		if x.End > 0 {
			return int(x.End), true
		}
	}
	return 0, false
}

type ann struct {
	start int
	end   int
	text  string
}

type editor struct {
	src     string
	anns    []ann
	usedPos map[int]bool
}

func newEditor(src string) *editor {
	return &editor{src: src, usedPos: make(map[int]bool)}
}

func (e *editor) claim(pos int) bool {
	if pos <= 0 || pos > len(e.src) || e.usedPos[pos] {
		return false
	}
	e.usedPos[pos] = true
	return true
}

func (e *editor) add(start, end int, text string) {
	if start == end && alreadyAnnotated(e.src, start, text) {
		return
	}
	if e.src[start:end] == text {
		return
	}
	e.anns = append(e.anns, ann{start: start, end: end, text: text})
}

func (e *editor) insert(pos int, text string) { e.add(pos, pos, text) }

func (e *editor) respell(at int, decl string) {
	if end, ok := typealgebra.AnnotationSpan(e.src, at); ok {
		e.add(at, end, " :"+typealgebra.CanonicalAnnotation(decl))
	}
}

func (e *editor) declare(pos int, ty engine.DynType, commentFmt string) {
	if native, ok := typealgebra.NativeAnnotation(ty); ok {
		e.insert(pos, " :"+native)
	} else {
		e.insert(pos, fmt.Sprintf(commentFmt, ty))
	}
}

func (e *editor) apply() string {
	sort.Slice(e.anns, func(i, j int) bool {
		return e.anns[i].start > e.anns[j].start
	})
	result := e.src
	lastStart := len(e.src)
	for _, a := range e.anns {
		if a.end > lastStart {
			continue
		}
		result = result[:a.start] + a.text + result[a.end:]
		lastStart = a.start
	}
	return result
}

func paramNameEnd(p *ast.Param) (int, bool) {
	raw := p.Name.Name
	if raw == "" || raw[0] == '@' {
		return 0, false
	}
	identLen := len(raw)
	if raw[0] == '.' {
		identLen--
	}
	return int(p.Name.Pos) + identLen, true
}

func (e *editor) exprs(env engine.TypeEnv) {
	for node, ty := range env.Nodes {
		if ty == engine.TUnknown {
			continue
		}
		if end, ok := exprEnd(node); ok && e.claim(end) {
			e.insert(end, fmt.Sprintf(" /* %s */", ty))
		}
	}
}

func (e *editor) params(env engine.TypeEnv, cls *engine.ClassObject) {
	for _, fn := range cls.SortedMethods {
		for i := range fn.Params {
			p := &fn.Params[i]
			ty, ok := env.Params[p]
			if !ok || ty == engine.TUnknown {
				continue
			}
			end, ok := paramNameEnd(p)
			if !ok || !e.claim(end) {
				continue
			}
			if p.Annotations != "" {
				e.respell(end, p.Annotations)
			} else {
				e.declare(end, ty, " /* %s */")
			}
		}
	}
}

func (e *editor) returns(env engine.TypeEnv, cls *engine.ClassObject) {
	for mname, ty := range env.Returns {
		fn, ok := cls.Methods[mname]
		if !ok || !e.claim(int(fn.Pos1)) {
			continue
		}
		if fn.ReturnAnnotation != "" {
			e.respell(int(fn.Pos1), fn.ReturnAnnotation)
		} else {
			e.declare(int(fn.Pos1), ty, " /* -> %s */")
		}
	}
}

func (e *editor) members(env engine.TypeEnv) {
	for mname, end := range memberValueEnds(e.src, env) {
		ty, ok := env.Members[mname]
		if !ok || ty == engine.TUnknown || !e.claim(end) {
			continue
		}
		e.insert(end, fmt.Sprintf(" /* :: %s */", ty))
	}
}

func annotatedSource(src string, env engine.TypeEnv, cls *engine.ClassObject) string {
	e := newEditor(src)
	e.exprs(env)
	e.params(env, cls)
	e.returns(env, cls)
	e.members(env)
	return e.apply()
}

func memberValueEnds(src string, env engine.TypeEnv) map[string]int {
	ends := make(map[string]int)
	if len(env.Members) == 0 {
		return ends
	}

	type item struct {
		tk  tok.Token
		pos int
		txt string
	}
	var items []item
	lxr := lexer.NewLexer(src)
	for {
		it := lxr.Next()
		if it.Token == tok.Eof {
			break
		}
		items = append(items, item{tk: it.Token, pos: int(it.Pos), txt: it.Text})
	}

	isSkip := func(t tok.Token) bool {
		return t == tok.Whitespace || t == tok.Newline || t == tok.Comment
	}
	nextSig := func(from int) int {
		for j := from; j < len(items); j++ {
			if !isSkip(items[j].tk) {
				return j
			}
		}
		return -1
	}

	depth := 0
	expectMember := false
	for i := 0; i < len(items); i++ {
		it := items[i]
		switch it.tk {
		case tok.LCurly, tok.LParen, tok.LBracket:
			depth++
			expectMember = depth == 1 && it.tk == tok.LCurly
			continue
		case tok.RCurly, tok.RParen, tok.RBracket:
			depth--
			expectMember = false
			continue
		}

		if depth != 1 {
			continue
		}

		if it.tk == tok.Newline || it.tk == tok.Comma || it.tk == tok.Semicolon {
			expectMember = true
			continue
		}
		if isSkip(it.tk) {
			continue
		}
		if !expectMember {
			continue
		}
		if it.tk != tok.Identifier {
			expectMember = false
			continue
		}

		name := it.txt
		if _, ok := env.Members[name]; !ok {
			expectMember = false
			continue
		}
		ni := nextSig(i + 1)
		if ni < 0 || items[ni].tk != tok.Colon {
			expectMember = false
			continue
		}

		valueDepth := 0
		endPos := -1
		endIdx := -1
		for k := ni + 1; k < len(items); k++ {
			t := items[k].tk
			done := false
			switch {
			case (t == tok.Newline || t == tok.Comma || t == tok.Semicolon) && valueDepth == 0:
				endPos = items[k].pos
				endIdx = k
				done = true
			case t == tok.LCurly || t == tok.LParen || t == tok.LBracket:
				valueDepth++
			case t == tok.RCurly || t == tok.RParen || t == tok.RBracket:
				if valueDepth == 0 {
					endPos = items[k].pos
					endIdx = k
					done = true
				} else {
					valueDepth--
				}
			}
			if done {
				e := k
				for e > ni+1 && isSkip(items[e-1].tk) {
					e--
				}
				endPos = items[e].pos
				break
			}
		}

		if endPos > 0 {
			ends[name] = endPos
		}
		if endIdx >= 0 {
			i = endIdx - 1
		} else {
			break
		}
		expectMember = false
	}

	return ends
}

func AnnotateClass(src string, env engine.TypeEnv, cls *engine.ClassObject) string {
	scoped := env.AnnotationView(targetNodeSet(cls))
	return annotatedSource(src, scoped, cls)
}
