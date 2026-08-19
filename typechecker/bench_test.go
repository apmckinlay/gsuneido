// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package typechecker

import (
	"fmt"
	"testing"
)

const benchClass = `class
	{
	New(.name, .limit = 10)
		{
		.items = Object()
		.total = 0
		}
	Add(item)
		{
		if not String?(item)
			return false
		.items.Add(item)
		.total += item.Size()
		return .total < .limit
		}
	Summary()
		{
		s = ""
		for x in .items
			s $= x.Upper() $ ", "
		return s.BeforeLast(", ")
		}
	Report(prefix = "")
		{
		n = .items.Size()
		if n is 0
			return prefix $ "empty"
		return prefix $ String(n) $ " items, " $ String(.total) $ " bytes"
		}
	Find(pat)
		{
		for x in .items
			if x.Has?(pat)
				return x
		return false
		}
	}`

func benchArgs(n int) []SourceEntry {
	out := make([]SourceEntry, n)
	for i := range out {
		out[i] = SourceEntry{Name: fmt.Sprintf("Bench%d", i), Src: benchClass}
	}
	return out
}

// one class per request - what TypeCheckerGui does for every row
func BenchmarkProcessOneClass(b *testing.B) {
	args := benchArgs(1)
	b.ResetTimer()
	for b.Loop() {
		if _, err := Process(Request{Method: "TypeInfer", Arguments: args}); err != nil {
			b.Fatal(err)
		}
	}
}

// a class with a lineage chain, as OrderedSrc produces
func BenchmarkProcessFourClasses(b *testing.B) {
	args := benchArgs(4)
	b.ResetTimer()
	for b.Loop() {
		if _, err := Process(Request{Method: "TypeInfer", Arguments: args}); err != nil {
			b.Fatal(err)
		}
	}
}

// references go through BuildReferenceRegistry, which builds its own PassCtx
func BenchmarkProcessWithReferences(b *testing.B) {
	args := benchArgs(1)
	refs := benchArgs(6)
	for i := range refs {
		refs[i].Name = fmt.Sprintf("Ref%d", i)
	}
	b.ResetTimer()
	for b.Loop() {
		if _, err := Process(Request{
			Method: "TypeInfer", Arguments: args, References: refs}); err != nil {
			b.Fatal(err)
		}
	}
}
