// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		// Basic: replace function name, swap holes
		r = CombyReplace("foo(a + b,
	c.d)",
			"foo(:[first], :[second])",
			"bar(:[second], :[first])")
		Assert(r is "bar(c.d, a + b)")

		// Single hole
		r = CombyReplace("foo(hello)", "foo(:[x])", "bar(:[x])")
		Assert(r is "bar(hello)")

		// No match: returns source unchanged
		r = CombyReplace("bar(baz)", "foo(:[x])", "bar(:[x])")
		Assert(r is "bar(baz)")

		// Multiple non-overlapping matches
		r = CombyReplace("foo(a) foo(b)", "foo(:[x])", "bar(:[x])")
		Assert(r is "bar(a) bar(b)")

		// Replace with literal (no holes in replace)
		r = CombyReplace("foo(a + b)", "foo(:[x])", "removed")
		Assert(r is "removed")

		// Hole at start
		r = CombyReplace("a + b + c", ":[x] + :[y]", ":[y] - :[x]")
		Assert(r is "b + c - a")

		// Chained method calls: swap order
		r = CombyReplace("foo.Bar().Baz()",
			":[a].:[b]().:[c]()",
			":[a].:[c]().:[b]()")
		Assert(r is "foo.Baz().Bar()")

		// Empty search pattern: returns source unchanged
		Assert(CombyReplace("anything", "", "x") is "anything")

		// Hole in replace not in search: uses literal :[text]
		r = CombyReplace("foo(a)", "foo(:[x])", "bar(:[x]) :[y]")
		Assert(r is "bar(a) :[y]")

		// Text between matches preserved
		r = CombyReplace("prefix foo(a) middle foo(b) suffix",
			"foo(:[x])", "bar(:[x])")
		Assert(r is "prefix bar(a) middle bar(b) suffix")

		// Replace with empty string
		r = CombyReplace("foo(a) suffix", "foo(:[x])", "")
		Assert(r is " suffix")

		// Replace entire source
		r = CombyReplace("foo(a)", "foo(:[x])", ":[x]")
		Assert(r is "a")

		// Replace at start of source
		r = CombyReplace("foo(a) end", "foo(:[x])", "bar(:[x])")
		Assert(r is "bar(a) end")

		// Replace at end of source
		r = CombyReplace("start foo(a)", "foo(:[x])", "bar(:[x])")
		Assert(r is "start bar(a)")

		// Whitespace in template
		r = CombyReplace("foo ( bar )", "foo ( :[x] )", "baz(:[x])")
		Assert(r is "baz(bar)")
		}

	Test_suneido_patterns()
		{
		// Class rename
		r = CombyReplace("class : Test { Test_main() { } }",
			"class : :[base] { :[body] }",
			"class : :[base]Renamed { :[body] }")
		Assert(r is "class : TestRenamed { Test_main() { } }")

		// Return statement wrap
		r = CombyReplace("return a + b", "return :[expr]", "return (:[expr])")
		Assert(r is "return (a + b)")

		// if false isnt pattern: unwrap assignment
		r = CombyReplace(
			"if false isnt r = CombyScanner(source, pos, mode)",
			"if false isnt :[var] = :[expr]",
			":[var] = :[expr]")
		Assert(r is "r = CombyScanner(source, pos, mode)")

		// Method call rename
		r = CombyReplace("record.Rule_name(field)",
			":[obj].:[method](:[args])",
			":[obj].:[method]2(:[args])")
		Assert(r is "record.Rule_name2(field)")
		}

	Test_from_to()
		{
		src = "prefix foo(a) middle foo(b) suffix"
		// foo(a) at pos:7 end:13, foo(b) at pos:21 end:27
		pat = "foo(:[x])"
		repl = "bar(:[x])"

		// No match in range: whole source returned unchanged
		r = CombyReplace(src, pat, repl, from: 14, to: 20)
		Assert(r is src)

		// from/to exactly covering one match: prefix + replacement + rest
		r = CombyReplace(src, pat, repl, from: 7, to: 13)
		Assert(r is "prefix bar(a) middle foo(b) suffix")

		// from starts after match begins: prefix preserved, partial skip, next matched
		r = CombyReplace(src, pat, repl, from: 8)
		Assert(r is "prefix foo(a) middle bar(b) suffix")

		// to ends before match end: match skipped, whole source returned
		r = CombyReplace(src, pat, repl, from: 7, to: 12)
		Assert(r is src)

		// from/to spanning multiple matches: both replaced
		r = CombyReplace(src, pat, repl, from: 0)
		Assert(r is "prefix bar(a) middle bar(b) suffix")

		// from at exact match start: prefix preserved, match replaced
		r = CombyReplace(src, pat, repl, from: 7)
		Assert(r is "prefix bar(a) middle bar(b) suffix")

		// to at exact match end: match replaced (not >)
		r = CombyReplace(src, pat, repl, from: 0, to: 13)
		Assert(r is "prefix bar(a) middle foo(b) suffix")

		// to just before any match: no replacement, whole source returned
		r = CombyReplace(src, pat, repl, from: 0, to: 6)
		Assert(r is src)
		}
	}
