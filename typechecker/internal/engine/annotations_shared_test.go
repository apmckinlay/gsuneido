// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

import (
	"testing"

	"pgregory.net/rapid"

	"github.com/apmckinlay/gsuneido/typechecker/internal/synth"
)

func deepCopyAnnotations(s AnnotationSet) AnnotationSet {
	out := make(AnnotationSet, len(s))
	for name, sigs := range s {
		cp := make([]Signature, len(sigs))
		for i := range sigs {
			cp[i] = sigs[i]
			cp[i].Params = append([]Param(nil), sigs[i].Params...)
		}
		out[name] = cp
	}
	return out
}

func dynStr(t DynType) string {
	if t == nil {
		return "<nil>"
	}
	return t.String()
}

func annotationsDiff(a, b AnnotationSet) string {
	if len(a) != len(b) {
		return "table size changed"
	}
	for name, as := range a {
		bs, ok := b[name]
		if !ok {
			return name + ": entry disappeared"
		}
		if len(as) != len(bs) {
			return name + ": overload count changed"
		}
		for i := range as {
			if dynStr(as[i].Returns) != dynStr(bs[i].Returns) ||
				dynStr(as[i].Receiver) != dynStr(bs[i].Receiver) ||
				as[i].AtParam != bs[i].AtParam {
				return name + ": signature changed"
			}
			if len(as[i].Params) != len(bs[i].Params) {
				return name + ": param count changed"
			}
			for j := range as[i].Params {
				p, q := as[i].Params[j], bs[i].Params[j]
				if p.Name != q.Name || dynStr(p.Typ) != dynStr(q.Typ) ||
					p.HasDefault != q.HasDefault ||
					p.Inferred != q.Inferred || p.Why != q.Why {
					return name + "." + p.Name + ": param changed"
				}
			}
		}
	}
	return ""
}

func TestPropSharedAnnotationsNotMutated(t *testing.T) {
	// any non-trivial table will do - the property is that a pipeline run
	// never writes through the shared pointer, whatever is in it
	withSigs(t,
		Reg{Receiver: "string", Name: "Size", Sig: "() :number"},
		Reg{Receiver: "string", Name: "Extract", Sig: "(pattern, part=false) :false|string"},
		Reg{Receiver: "object", Name: "Add", Sig: "(@args) :object"},
		Reg{Kind: "free", Name: "Object", Sig: "(@args) :object"},
		Reg{Kind: "static", Class: "Database", Name: "Info", Sig: "() :object"},
	)
	shared := Annotations()
	if len(shared) == 0 {
		t.Fatal("annotations not loaded - test would be vacuous")
	}
	before := deepCopyAnnotations(shared)
	rapid.Check(t, func(rt *rapid.T) {
		src := genSource(rt, synth.Config{})
		cls, ok := safeParse(src, "T")
		if !ok {
			rt.Skip("unparseable render")
		}
		env := NewTypeEnv()
		RunPipeline(cls, env, &PassCtx{Annotations: shared}, nil)
	})
	if d := annotationsDiff(before, Annotations()); d != "" {
		t.Fatalf("shared annotations mutated by a pipeline run: %s", d)
	}
}
