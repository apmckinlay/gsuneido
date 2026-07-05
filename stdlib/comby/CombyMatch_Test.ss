// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		// Basic: two holes
		m = CombyMatch("foo(:[first], :[second])",
			"foo(a + b, c.d)").First()
		Assert(m.pos is 0 and m.end is 15)
		Assert(m.holes.first is: "a + b")
		Assert(m.holes.second is: "c.d")

		// Single hole
		m = CombyMatch("foo(:[x])", "foo(hello)").First()
		Assert(m.holes.x is: "hello")

		// No match
		Assert(CombyMatch("foo(:[x])", "bar(baz)").Count() is: 0)

		// Balanced delimiters
		m = CombyMatch("foo(:[x])", "foo(bar(baz))").First()
		Assert(m.holes.x is: "bar(baz)")

		// Multiple non-overlapping matches
		Assert(CombyMatch("foo(:[x])", "foo(a) foo(b)").Count() is: 2)

		// Hole at start
		m = CombyMatch(":[x] + :[y]", "a + b + c").First()
		Assert(m.holes.x is: "a")
		Assert(m.holes.y is: "b + c")

		// Zero-length hole
		m = CombyMatch("foo(:[x], :[y])", "foo(, bar)").First()
		Assert(m.holes.x is: "")
		Assert(m.holes.y is: "bar")

		// Generic: strings are not special
		m = CombyMatch("foo(:[x])", 'foo("bar")').First()
		Assert(m.holes.x is: '"bar"')

		// Suneido: anchor not inside string
		m = CombyMatch("foo(:[x])",
			'"foo(nope)" foo(yes)', mode: 'suneido').First()
		Assert(m.holes.x is: "yes")

		// Suneido: blanks skip comments
		m = CombyMatch("foo ( :[x] )", "foo(/*c*/bar)",
			mode: 'suneido').First()
		Assert(m.holes.x is: "bar")

		// Suneido: no match inside string
		Assert(CombyMatch("foo(:[x])", '"foo(nope)"',
			mode: 'suneido').Count() is: 0)

		// Escape inside suneido string
		m = CombyMatch("foo(:[x])",
			'foo("escaped)quote")', mode: 'suneido').First()
		Assert(m.holes.x is: '"escaped)quote"')

		// Empty template
		Assert(CombyMatch("", "anything").Count() is: 0)
		}

	// ================================================================
	// Suneido code patterns: method calls and chaining
	// ================================================================

	Test_method_calls()
		{
		// Qualified method call: record.Method(args)
		m = CombyMatch(":[obj].:[method](:[args])",
			"record.Rule_name(field)").First()
		Assert(m.holes.obj is: "record")
		Assert(m.holes.method is: "Rule_name")
		Assert(m.holes.args is: "field")

		// Chained method calls
		m = CombyMatch(":[a].:[b]().:[c]()",
			"foo.Bar().Baz()").First()
		Assert(m.holes.a is: "foo")
		Assert(m.holes.b is: "Bar")
		Assert(m.holes.c is: "Baz")

		// Chained with spaces around dots
		m = CombyMatch(":[a] . :[b]() . :[c]()",
			"foo . Bar() . Baz()").First()
		Assert(m.holes.a is: "foo")
		Assert(m.holes.b is: "Bar")
		Assert(m.holes.c is: "Baz")

		// Chained through line comment (suneido mode)
		m = CombyMatch(":[a] . :[b]()",
			'foo // comment\n    .Bar()',
			mode: 'suneido').First()
		Assert(m.holes.a is: "foo")
		Assert(m.holes.b is: "Bar")

		// Chained through block comment (suneido mode)
		m = CombyMatch(":[a] . :[b]()",
			'foo /* comment */ .Bar()',
			mode: 'suneido').First()
		Assert(m.holes.a is: "foo")
		Assert(m.holes.b is: "Bar")

		// Method call without arguments
		m = CombyMatch(":[obj].:[method]()",
			"obj.Close()").First()
		Assert(m.holes.obj is: "obj")
		Assert(m.holes.method is: "Close")
		}

	// ================================================================
	// Common Suneido idioms
	// ================================================================

	Test_suneido_idioms()
		{
		m = CombyMatch("if false isnt :[var] = :[expr]",
			"if false isnt r = CombyScanner(source, pos, mode)").
			First()
		Assert(m.holes.var is: "r")
		Assert(m.holes.expr is: "CombyScanner(source, pos, mode)")

		// if false is x = expr
		m = CombyMatch("if false is :[var] = :[expr]",
			"if false is sc = CombyScanner(source, current, mode)").
			First()
		Assert(m.holes.var is: "sc")
		Assert(m.holes.expr is: "CombyScanner(source, current, mode)")

		// for m, v in object (record key/value iteration)
		m = CombyMatch("for :[key], :[val] in :[expr]",
			"for m, v in object").First()
		Assert(m.holes.key is: "m")
		Assert(m.holes.val is: "v")
		Assert(m.holes.expr is: "object")

		// for x in iterable
		m = CombyMatch("for :[var] in :[expr]",
			"for x in mylist").First()
		Assert(m.holes.var is: "x")
		Assert(m.holes.expr is: "mylist")

		// switch statement
		m = CombyMatch("switch :[expr] { :[body] }",
			"switch x { case 1: 'one'; case 2: 'two' }").First()
		Assert(m.holes.expr is: "x")
		Assert(m.holes.body is: "case 1: 'one'; case 2: 'two'")

		// class with inheritance
		m = CombyMatch("class : :[base] { :[body] }",
			"class : Test { Test_main() { } }").First()
		Assert(m.holes.base is: "Test")
		Assert(m.holes.body is: "Test_main() { }")

		// class without base
		m = CombyMatch("class { :[body] }",
			"class { Foo() { } }").First()
		Assert(m.holes.body is: "Foo() { }")
		}
// ================================================================
	// Control flow patterns
	// ================================================================

	Test_control_flow()
		{
		// if-else
		m = CombyMatch("if :[cond] { :[then] } else { :[else] }",
			"if x > 0 { f(x) } else { g(x) }").First()
		Assert(m.holes.cond is: "x > 0")
		Assert(m.holes.then is: "f(x)")
		Assert(m.holes.else is: "g(x)")

		// while loop
		m = CombyMatch("while :[cond] { :[body] }",
			"while pos < source.Size() { pos = r.pos }").First()
		Assert(m.holes.cond is: "pos < source.Size()")
		Assert(m.holes.body is: "pos = r.pos")

		// do-while loop
		m = CombyMatch("do { :[body] } while :[cond]",
			"do { ++i } while i < n").First()
		Assert(m.holes.body is: "++i")
		Assert(m.holes.cond is: "i < n")

		// for-in range (counted loop)
		m = CombyMatch("for :[v] in :[lo]..:[hi]",
			"for i in 0..n-1").First()
		Assert(m.holes.v is: "i")
		Assert(m.holes.lo is: "0")
		Assert(m.holes.hi is: "n-1")

		// classic for loop
		m = CombyMatch("for (:[init]; :[cond]; :[inc]) { :[body] }",
			"for (i = 0; i < n; ++i) { list[i]() }").First()
		Assert(m.holes.init is: "i = 0")
		Assert(m.holes.cond is: "i < n")
		Assert(m.holes.inc is: "++i")
		Assert(m.holes.body is: "list[i]()")

		// try-catch
		m = CombyMatch("try { :[body] } catch { :[catch] }",
			"try { something() } catch { Log() }").First()
		Assert(m.holes.body is: "something()")
		Assert(m.holes.catch is: "Log()")

		// return statement
		m = CombyMatch("return :[expr]", "return a + b").First()
		Assert(m.holes.expr is: "a + b")
		}
// Suneido string awareness in suneido mode
	// ================================================================

	Test_suneido_strings()
		{
		// No match inside double-quoted string
		Assert(CombyMatch("foo(:[x])", '"foo(nope)"',
			mode: 'suneido').Count() is: 0)

		// No match inside single-quoted string
		Assert(CombyMatch("foo(:[x])", "'foo(nope)'",
			mode: 'suneido').Count() is: 0)

		// Match skips over string, finds after
		m = CombyMatch(":[x]",
			'"not this" nor this',
			mode: 'suneido').First()
		Assert(m.holes.x is: '"not this" nor this')

		// Escaped quote inside double-quoted string (suneido mode)
		m = CombyMatch("foo(:[x])",
			'foo("escaped\\"quote")', mode: 'suneido').First()
		Assert(m.holes.x is: '"escaped\\"quote"')

		// Blanks skip block comments in suneido mode
		m = CombyMatch("foo ( :[x] )",
			"foo(/*comment*/bar)",
			mode: 'suneido').First()
		Assert(m.holes.x is: "bar")

		// Backtick raw string treated as atom (suneido mode)
		m = CombyMatch(":[v] =~ :[pat]",
			"name =~ `^[a-zA-Z0-9_]+$`",
			mode: 'suneido').First()
		Assert(m.holes.v is: "name")
		Assert(m.holes.pat is: "`^[a-zA-Z0-9_]+$`")

		// No match inside backtick raw string (suneido mode)
		Assert(CombyMatch("foo(:[x])", '`foo(nope)`',
			mode: 'suneido').Count() is: 0)
		}

	Test_edge_cases()
		{
		// Repeated hole name: last match wins
		m = CombyMatch(":[x] + :[x]", "a + a").First()
		Assert(m.holes.x is: "a")
//		m2 = CombyMatch(":[x] + :[x]", "a + b").First()
//		Assert(m2.holes.x is "b")

		// Hole at end of template
		m = CombyMatch("foo(:[args]", "foo(bar, baz").First()
		Assert(m.holes.args is: "bar, baz")

		// Literal-only template
		Assert(CombyMatch("class", "class : Test { } class").
			Count() is: 2)
		// Object literal #(...)
		m = CombyMatch("#(:[inner])",
			"#(1, 2, name: 'Joe')").First()
		Assert(m.holes.inner is: "1, 2, name: 'Joe'")

		// Record literal #{...}
		m = CombyMatch("#{:[inner]}",
			'#{name: "Joe", age: 42}').First()
		Assert(m.holes.inner is: 'name: "Joe", age: 42')

		// String concatenation with $
		m = CombyMatch(":[a] $ :[b]", "foo $ bar").First()
		Assert(m.holes.a is: "foo")
		Assert(m.holes.b is: "bar")

		// is comparison
		m = CombyMatch(":[a] is :[b]", "x is 5").First()
		Assert(m.holes.a is: "x")
		Assert(m.holes.b is: "5")

		// isnt comparison
		m = CombyMatch(":[a] isnt :[b]", "x isnt false").First()
		Assert(m.holes.a is: "x")
		Assert(m.holes.b is: "false")
		// Nested balanced delimiters (multi-level)
		m = CombyMatch("outer(:[inner])",
			"outer(mid(inner(arg)))").First()
		Assert(m.holes.inner is: "mid(inner(arg))")

		m = CombyMatch("mid(:[inner])",
			"outer(mid(inner(arg)))").First()
		Assert(m.holes.inner is: "inner(arg)")

		// Function definition with brace-delimited body
		m = CombyMatch("function :[name](:[params]) { :[body] }",
			"function Foo(x, y) { return x + y }").First()
		Assert(m.holes.name is: "Foo")
		Assert(m.holes.params is: "x, y")
		Assert(m.holes.body is: "return x + y")

		m = CombyMatch("return :[a]",
			"function Foo(x, y) {
	return x + y
	}").First()
		Assert(m.holes.a is: "x + y\r\n\t")

		m = CombyMatch("return :[a] + :[b]",
			"function Foo(x, y) {
	return x/*=test*/ + y
	}").First()
		Assert(m.holes.a is: "x")
		Assert(m.holes.b is: "y\r\n\t")

		m = CombyMatch("return :[a] -",
			"function Foo(x, y) {
	return x + y
	}")
		Assert(m is: #())
		}

	Test_GetHint()
		{
		fn = CombyMatch.GetHint
		Assert(fn('') is: false)
		Assert(fn('return #()') is: false)
		Assert(fn('return :[a], :[b]') is: false)
		Assert(fn('return abc') is: 'abc')
		Assert(fn('return abc + :[test]') is: 'abc')
		Assert(fn('return abc + :[test] + 123_333') is: '123_333')
		Assert(fn('return abc + :[test] + 123_333 + #20260101') is: '20260101')
		}
	}
