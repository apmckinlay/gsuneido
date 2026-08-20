// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		// Basic: two holes
		m = CombyMatch.All("foo(:[first], :[second])",
			"foo(a + b, c.d)").First()
		Assert(m.pos is 0 and m.end is 15)
		Assert(m.holes.first is: "a + b")
		Assert(m.holes.second is: "c.d")

		// Single hole
		m = CombyMatch.All("foo(:[x])", "foo(hello)").First()
		Assert(m.holes.x is: "hello")

		// No match
		Assert(CombyMatch.All("foo(:[x])", "bar(baz)").Count() is: 0)

		// Balanced delimiters
		m = CombyMatch.All("foo(:[x])", "foo(bar(baz))").First()
		Assert(m.holes.x is: "bar(baz)")

		// Multiple non-overlapping matches
		Assert(CombyMatch.All("foo(:[x])", "foo(a) foo(b)").Count() is: 2)

		// Hole at start
		m = CombyMatch.All(":[x] + :[y]", "a + b + c").First()
		Assert(m.holes.x is: "a")
		Assert(m.holes.y is: "b + c")

		// Zero-length hole
		m = CombyMatch.All("foo(:[x], :[y])", "foo(, bar)").First()
		Assert(m.holes.x is: "")
		Assert(m.holes.y is: "bar")

		// Generic: strings are not special
		m = CombyMatch.All("foo(:[x])", 'foo("bar")').First()
		Assert(m.holes.x is: '"bar"')

		// Suneido: anchor not inside string
		m = CombyMatch.All("foo(:[x])",
			'"foo(nope)" foo(yes)').First()
		Assert(m.holes.x is: "yes")

		// Suneido: blanks skip comments
		m = CombyMatch.All("foo ( :[x] )", "foo(/*c*/bar)").First()
		Assert(m.holes.x is: "bar")

		// Suneido: no match inside string
		Assert(CombyMatch.All("foo(:[x])", '"foo(nope)"').Count() is: 0)

		// Escape inside suneido string
		m = CombyMatch.All("foo(:[x])",
			'foo("escaped)quote")').First()
		Assert(m.holes.x is: '"escaped)quote"')

		// Empty template
		Assert(CombyMatch.All("", "anything").Count() is: 0)
		}

	// ================================================================
	// Suneido code patterns: method calls and chaining
	// ================================================================

	Test_method_calls()
		{
		// Qualified method call: record.Method(args)
		m = CombyMatch.All(":[obj].:[method](:[args])",
			"record.Rule_name(field)").First()
		Assert(m.holes.obj is: "record")
		Assert(m.holes.method is: "Rule_name")
		Assert(m.holes.args is: "field")

		// Chained method calls
		m = CombyMatch.All(":[a].:[b]().:[c]()",
			"foo.Bar().Baz()").First()
		Assert(m.holes.a is: "foo")
		Assert(m.holes.b is: "Bar")
		Assert(m.holes.c is: "Baz")

		// Chained with spaces around dots
		m = CombyMatch.All(":[a] . :[b]() . :[c]()",
			"foo . Bar() . Baz()").First()
		Assert(m.holes.a is: "foo")
		Assert(m.holes.b is: "Bar")
		Assert(m.holes.c is: "Baz")

		// Chained through line comment
		m = CombyMatch.All(":[a] . :[b]()",
			'foo // comment\n    .Bar()').First()
		Assert(m.holes.a is: "foo")
		Assert(m.holes.b is: "Bar")

		// Chained through block comment
		m = CombyMatch.All(":[a] . :[b]()",
			'foo /* comment */ .Bar()').First()
		Assert(m.holes.a is: "foo")
		Assert(m.holes.b is: "Bar")

		// Method call without arguments
		m = CombyMatch.All(":[obj].:[method]()",
			"obj.Close()").First()
		Assert(m.holes.obj is: "obj")
		Assert(m.holes.method is: "Close")
		}

	// ================================================================
	// Common Suneido idioms
	// ================================================================

	Test_suneido_idioms()
		{
		m = CombyMatch.All("if false isnt :[var] = :[expr]",
			"if false isnt r = ABC(source, pos, mode)").
			First()
		Assert(m.holes.var is: "r")
		Assert(m.holes.expr is: "ABC(source, pos, mode)")

		// if false is x = expr
		m = CombyMatch.All("if false is :[var] = :[expr]",
			"if false is sc = ABC(source, current, mode)").
			First()
		Assert(m.holes.var is: "sc")
		Assert(m.holes.expr is: "ABC(source, current, mode)")

		// for m, v in object (record key/value iteration)
		m = CombyMatch.All("for :[key], :[val] in :[expr]",
			"for m, v in object").First()
		Assert(m.holes.key is: "m")
		Assert(m.holes.val is: "v")
		Assert(m.holes.expr is: "object")

		// for x in iterable
		m = CombyMatch.All("for :[var] in :[expr]",
			"for x in mylist").First()
		Assert(m.holes.var is: "x")
		Assert(m.holes.expr is: "mylist")

		// switch statement
		m = CombyMatch.All("switch :[expr] { :[body] }",
			"switch x { case 1: 'one'; case 2: 'two' }").First()
		Assert(m.holes.expr is: "x")
		Assert(m.holes.body is: "case 1: 'one'; case 2: 'two'")

		// class with inheritance
		m = CombyMatch.All("class : :[base] { :[body] }",
			"class : Test { Test_main() { } }").First()
		Assert(m.holes.base is: "Test")
		Assert(m.holes.body is: "Test_main() { }")

		// class without base
		m = CombyMatch.All("class { :[body] }",
			"class { Foo() { } }").First()
		Assert(m.holes.body is: "Foo() { }")
		}
// ================================================================
	// Control flow patterns
	// ================================================================

	Test_control_flow()
		{
		// if-else
		m = CombyMatch.All("if :[cond] { :[then] } else { :[else] }",
			"if x > 0 { f(x) } else { g(x) }").First()
		Assert(m.holes.cond is: "x > 0")
		Assert(m.holes.then is: "f(x)")
		Assert(m.holes.else is: "g(x)")

		// while loop
		m = CombyMatch.All("while :[cond] { :[body] }",
			"while pos < source.Size() { pos = r.pos }").First()
		Assert(m.holes.cond is: "pos < source.Size()")
		Assert(m.holes.body is: "pos = r.pos")

		// do-while loop
		m = CombyMatch.All("do { :[body] } while :[cond]",
			"do { ++i } while i < n").First()
		Assert(m.holes.body is: "++i")
		Assert(m.holes.cond is: "i < n")

		// for-in range (counted loop)
		m = CombyMatch.All("for :[v] in :[lo]..:[hi]",
			"for i in 0..n-1").First()
		Assert(m.holes.v is: "i")
		Assert(m.holes.lo is: "0")
		Assert(m.holes.hi is: "n-1")

		// classic for loop
		m = CombyMatch.All("for (:[init]; :[cond]; :[inc]) { :[body] }",
			"for (i = 0; i < n; ++i) { list[i]() }").First()
		Assert(m.holes.init is: "i = 0")
		Assert(m.holes.cond is: "i < n")
		Assert(m.holes.inc is: "++i")
		Assert(m.holes.body is: "list[i]()")

		// try-catch
		m = CombyMatch.All("try { :[body] } catch { :[catch] }",
			"try { something() } catch { Log() }").First()
		Assert(m.holes.body is: "something()")
		Assert(m.holes.catch is: "Log()")

		// return statement
		m = CombyMatch.All("return :[expr]", "return a + b").First()
		Assert(m.holes.expr is: "a + b")
		}

	// ================================================================
	Test_suneido_strings()
		{
		// No match inside double-quoted string
		Assert(CombyMatch.All("foo(:[x])", '"foo(nope)"').Count() is: 0)

		// No match inside single-quoted string
		Assert(CombyMatch.All("foo(:[x])", "'foo(nope)'").Count() is: 0)

		// Match skips over string, finds after
		m = CombyMatch.All(":[x]", '"not this" nor this').First()
		Assert(m.holes.x is: '"not this" nor this')

		// Escaped quote inside double-quoted string
		m = CombyMatch.All("foo(:[x])",
			'foo("escaped\\"quote")').First()
		Assert(m.holes.x is: '"escaped\\"quote"')

		// Blanks skip block comments in
		m = CombyMatch.All("foo ( :[x] )", "foo(/*comment*/bar)").First()
		Assert(m.holes.x is: "bar")

		// Backtick raw string treated as atom
		m = CombyMatch.All(":[v] =~ :[pat]", "name =~ `^[a-zA-Z0-9_]+$`").First()
		Assert(m.holes.v is: "name")
		Assert(m.holes.pat is: "`^[a-zA-Z0-9_]+$`")

		// No match inside backtick raw string
		Assert(CombyMatch.All("foo(:[x])", '`foo(nope)`').Count() is: 0)
		}

	Test_edge_cases()
		{
		// Repeated hole name: last match wins
		m = CombyMatch.All(":[x] + :[x]", "a + a").First()
		Assert(m.holes.x is: "a")
//		m2 = CombyMatch.All(":[x] + :[x]", "a + b").First()
//		Assert(m2.holes.x is "b")

		// Hole at end of template
		m = CombyMatch.All("foo(:[args]", "foo(bar, baz").First()
		Assert(m.holes.args is: "bar, baz")

		// Literal-only template
		Assert(CombyMatch.All("class", "class : Test { } class").
			Count() is: 2)
		// Object literal #(...)
		m = CombyMatch.All("#(:[inner])",
			"#(1, 2, name: 'Joe')").First()
		Assert(m.holes.inner is: "1, 2, name: 'Joe'")

		// Record literal #{...}
		m = CombyMatch.All("#{:[inner]}",
			'#{name: "Joe", age: 42}').First()
		Assert(m.holes.inner is: 'name: "Joe", age: 42')

		// String concatenation with $
		m = CombyMatch.All(":[a] $ :[b]", "foo $ bar").First()
		Assert(m.holes.a is: "foo")
		Assert(m.holes.b is: "bar")

		// is comparison
		m = CombyMatch.All(":[a] is :[b]", "x is 5").First()
		Assert(m.holes.a is: "x")
		Assert(m.holes.b is: "5")

		// isnt comparison
		m = CombyMatch.All(":[a] isnt :[b]", "x isnt false").First()
		Assert(m.holes.a is: "x")
		Assert(m.holes.b is: "false")
		// Nested balanced delimiters (multi-level)
		m = CombyMatch.All("outer(:[inner])",
			"outer(mid(inner(arg)))").First()
		Assert(m.holes.inner is: "mid(inner(arg))")

		m = CombyMatch.All("mid(:[inner])",
			"outer(mid(inner(arg)))").First()
		Assert(m.holes.inner is: "inner(arg)")

		// Function definition with brace-delimited body
		m = CombyMatch.All("function :[name](:[params]) { :[body] }",
			"function Foo(x, y) { return x + y }").First()
		Assert(m.holes.name is: "Foo")
		Assert(m.holes.params is: "x, y")
		Assert(m.holes.body is: "return x + y")

		m = CombyMatch.All("return :[a]",
			"function Foo(x, y) {
	return x + y
	}").First()
		Assert(m.holes.a is: "x + y\r\n\t")

		m = CombyMatch.All("return :[a] + :[b]",
			"function Foo(x, y) {
	return x/*=test*/ + y
	}").First()
		Assert(m.holes.a is: "x")
		Assert(m.holes.b is: "y\r\n\t")

		m = CombyMatch.All("return :[a] -",
			"function Foo(x, y) {
	return x + y
	}")
		Assert(m is: #())
		}


	Test_Match()
		{
		// Forward: first match from pos 0
		m = CombyMatch("foo(:[x])", "foo(a) foo(b)")
		Assert(m.holes.x is "a")
		Assert(m.pos is 0)

		// Forward: skip first, match second from pos 7
		m = CombyMatch("foo(:[x])", "foo(a) foo(b)", pos: 7)
		Assert(m.holes.x is "b")
		Assert(m.pos is 7)

		// Forward: no match beyond all candidates
		r = CombyMatch("foo(:[x])", "foo(a) foo(b)", pos: 20)
		Assert(r is false)

		// Forward: no matching pattern
		r = CombyMatch("bar(:[x])", "foo(a) foo(b)")
		Assert(r is false)

		// Forward: empty template
		r = CombyMatch("", "anything")
		Assert(r is false)

		// Forward: single hole
		m = CombyMatch("foo(:[x])", "foo(hello)")
		Assert(m.holes.x is "hello")
		Assert(m.pos is 0 and m.end is 10)

		// Forward: balanced delimiters in hole
		m = CombyMatch("foo(:[x])", "foo(bar(baz))")
		Assert(m.holes.x is "bar(baz)")

		// Backward: last match from end
		m = CombyMatch("foo(:[x])", "foo(a) foo(b)",
			pos: "foo(a) foo(b)".Size(), prev: true)
		Assert(m.holes.x is "b")

		// Backward: from pos just after first match
		// "foo(a)" ends at 6, pos:6 should find it
		m = CombyMatch("foo(:[x])", "foo(a) foo(b)", pos: 6, prev: true)
		Assert(m.holes.x is "a")
		Assert(m.pos is 0)

		// Backward: match end must not exceed pos
		// "foo(a)" ends at 6, pos:4 should give no match
		r = CombyMatch("foo(:[x])", "foo(a) foo(b)", pos: 4, prev: true)
		Assert(r is false)

		// Backward: pos before any match
		r = CombyMatch("foo(:[x])", "foo(a) foo(b)", pos: 1, prev: true)
		Assert(r is false)

		// Backward: literal-only template, match from end
		s = "class : Test { } class"
		m = CombyMatch("class", s, pos: s.Size(), prev: true)
		Assert(m.pos is 17) // second "class"
		Assert(m.end is 22)

		// Backward: literal-only from midpoint
		m = CombyMatch("class", s, pos: 10, prev: true)
		Assert(m.pos is 0)
		}

	Test_expression_holes()
		{
		m = CombyMatch(':[x:e]', 'foo')
		Assert(m.holes.x is: 'foo')
		m = CombyMatch(':[x:e]', 'foo.bar')
		Assert(m.holes.x is: 'foo.bar')
		m = CombyMatch(':[x:e]', 'foo[
			bar]')
		Assert(m.holes.x is: 'foo[
			bar]')
		m = CombyMatch(':[x:e]', 'foo.
			bar')
		Assert(m.holes.x is: 'foo.
			bar')

		m = CombyMatch(':[x:e]', 'function(foo, bar)')
		Assert(m.holes.x is: 'function(foo, bar)')
		m = CombyMatch(':[x:e]', 'foo bar')
		Assert(m.holes.x is: 'foo')
		m = CombyMatch(':[x:e] + :[y:e]', 'a + b')
		Assert(m.holes.x is: 'a')
		Assert(m.holes.y is: 'b')
		m = CombyMatch(':[x:e]', 'f(g(h()))')
		Assert(m.holes.x is: 'f(g(h()))')
		m = CombyMatch(':[x:e]()', 'bar + foo()')
		Assert(m.holes.x is: 'foo')
		m = CombyMatch(':[x:e](:[y:e])', 'foo(bar)')
		Assert(m.holes.x is: 'foo')
		Assert(m.holes.y is: 'bar')
		m = CombyMatch('foo(:[x:e]', 'foo(bar')
		Assert(m.holes.x is: 'bar')
		matches = CombyMatch.All(':[x:e]', 'foo bar baz')
		Assert(matches.Count() is: 3)
		Assert(matches[0].holes.x is: 'foo')
		Assert(matches[1].holes.x is: 'bar')
		Assert(matches[2].holes.x is: 'baz')
		m = CombyMatch(':[x:e][:[y:e]]', 'arr[0]')
		Assert(m.holes.x is: 'arr')
		Assert(m.holes.y is: '0')
		matches = CombyMatch.All(':[x:e]', 'a + foo.bar().baz()')
		Assert(matches[0].holes.x is: 'a')
		Assert(matches[1].holes.x is: '+')
		Assert(matches[2].holes.x is: 'foo.bar().baz()')
		}

	// ================================================================
	// Search patterns in longer, realistic code snippets
	// ================================================================

	Test_longer_contexts()
		{
		// Typical function: two-line body with early return pattern
		src = "function Parse(source, mode = false)\r\n" $
			"\t{\r\n" $
			"\tif false isnt r = Scanner(source)\r\n" $
			"\t\treturn r\r\n" $
			"\treturn false\r\n" $
			"\t}"

		m = CombyMatch.All('if false isnt :[var] = :[expr:e]', src)[0]
		Assert(m.holes.var is: 'r')
		Assert(m.holes.expr is: 'Scanner(source)')

		// for-in loop inside function
		src = "function Iter(c)\r\n" $
			"\t{\r\n" $
			"\tfor x in c.List()\r\n" $
			"\t\tx()\r\n" $
			"\t}"
		m = CombyMatch.All('for :[v:e] in :[e:e]', src)[0]
		Assert(m.holes.v is: 'x')
		Assert(m.holes.e is: 'c.List()')

		// try-catch inside a function
		src = "function Try()\r\n" $
			"\t{\r\n" $
			"\ttry { risky() }\r\n" $
			"\tcatch (e, 'block:break') { ; }\r\n" $
			"\tcatch { Log('fail') }\r\n" $
			"\t}"
		m = CombyMatch.All('try { :[body] } catch', src)[0]
		Assert(m.holes.body.Has?('risky()'))

		// Multiple catch clauses with pattern string
		m = CombyMatch.All('catch (:[var], :[pat])', src)[0]
		Assert(m.holes.var is: 'e')
		Assert(m.holes.pat is: "'block:break'")

		// switch inside a function body
		src = "function Dispatch(code)\r\n" $
			"\t{\r\n" $
			"\tswitch code\r\n" $
			"\t\t{\r\n" $
			"\t\tcase 1: 'one'\r\n" $
			"\t\tcase 2: 'two'\r\n" $
			"\t\t}\r\n" $
			"\t}"
		m = CombyMatch.All('switch :[e] { :[body] }', src)[0]
		Assert(m.holes.e is: 'code')
		Assert(m.holes.body.Has?("case 1: 'one'"))

		// class with inheritance and multiple methods
		src = "class : Test\r\n" $
			"\t{\r\n" $
			"\tSetup() { }\r\n" $
			"\tTest_foo() { Assert(1 is 1) }\r\n" $
			"\tTeardown() { }\r\n" $
			"\t}"
		m = CombyMatch.All('class : :[base] { :[body] }', src)[0]
		Assert(m.holes.base is: 'Test')
		Assert(m.holes.body.Has?('Setup()'))
		Assert(m.holes.body.Has?('Teardown()'))

		// return inside a nested block
		src = "function Find(items, key)\r\n" $
			"\t{\r\n" $
			"\tfor x in items\r\n" $
			"\t\tif x.name is key\r\n" $
			"\t\t\treturn x\r\n" $
			"\treturn false\r\n" $
			"\t}"
		matches = CombyMatch.All('return :[e:e]', src)
		Assert(matches.Count() is: 2)
		Assert(matches[0].holes.e is: 'x')
		Assert(matches[1].holes.e is: 'false')

		// for m, v in object inside function
		src = "function Copy(ob)\r\n" $
			"\t{\r\n" $
			"\tfor m, v in ob\r\n" $
			"\t\tresult[m] = v\r\n" $
			"\t}"
		m = CombyMatch.All('for :[m:e], :[v:e] in :[e:e]', src)[0]
		Assert(m.holes.m is: 'm')
		Assert(m.holes.v is: 'v')
		Assert(m.holes.e is: 'ob')

		// do-while loop in context
		src = "function Pop(list)\r\n" $
			"\t{\r\n" $
			"\tdo { result = list.Pop() }\r\n" $
			"\t\twhile list.Size() > 0\r\n" $
			"\t}"
		m = CombyMatch.All('do { :[body] } while :[cond]', src)[0]
		Assert(m.holes.body is: 'result = list.Pop()')
		Assert(m.holes.cond.Has?('list.Size() > 0'))

		// while loop spanning multiple tokens
		src = "while i < tokens.Size() and tokens[i].start < pos
	i++"
		m = CombyMatch.All('while :[cond]', src)[0]
		Assert(m.holes.cond.Has?('tokens.Size()'))
		Assert(m.holes.cond.Has?('tokens[i].start'))

		// class without inheritance, multi-line body
		src = "class\r\n" $
			"\t{\r\n" $
			"\tNew() { }\r\n" $
			"\tCall(x) { return x }\r\n" $
			"\t}"
		m = CombyMatch.All('class { :[body] }', src)[0]
		Assert(m.holes.body.Has?('New()'))
		Assert(m.holes.body.Has?('Call(x)'))

		// if-else with longer blocks
		src = "if x > 0 and y < 10\r\n" $
			"\t{ f(x); g(y) }\r\n" $
			"\telse { h(); k() }"
		m = CombyMatch.All('if :[cond] { :[then] } else { :[else] }',
			src)[0]
		Assert(m.holes.cond is: 'x > 0 and y < 10')
		Assert(m.holes.then is: 'f(x); g(y)')
		Assert(m.holes.else is: 'h(); k()')
		}

	// ================================================================
	// Real-world search patterns: practical refactoring / analysis
	// ================================================================

	Test_real_world_patterns()
		{
		// Capture receiver and method name: record.Rule_name(field)
		src = "function Apply(record, field)\r\n" $
			"\t{\r\n" $
			"\tif record.Member?(field)\r\n" $
			"\t\treturn record[field]\r\n" $
			"\treturn record.Rule_name(field)\r\n" $
			"\t}"
		m = CombyMatch.All('if :[rec:e].:[method:e](:[args:e])', src)[0]
		Assert(m.holes.rec is: 'record')
		Assert(m.holes.method is: 'Member?')
		Assert(m.holes.args is: 'field')

		// .Size() calls: find all instances
		matches = CombyMatch.All(':[x].Size()', src)
		Assert(matches.Count() is: 0) // no .Size() in this source

		src2 = "if items.Size() > 0\r\n" $
			"\treturn items[items.Size() - 1]"
		matches = CombyMatch.All(':[x:e].Size()', src2)
		Assert(matches.Count() is: 2)
		Assert(matches[0].holes.x is: 'items')
		Assert(matches[1].holes.x is: 'items')

		// isnt false pattern
		src = "if sc isnt false\r\n" $
			"\treturn sc\r\n" $
			"\tthrow 'not found'"
		m = CombyMatch.All(':[x:e] isnt false', src)[0]
		Assert(m.holes.x is: 'sc')

		// assignment from method: result = .SomeFn(...)
		src = "\tresult = .buildTokens(s)\r\n" $
			"\ttokens = .Parse(source)\r\n" $
			"\titems = CombyTemplate(search)"
		matches = CombyMatch.All(':[var:e] = .:[method:e](:[args:e])', src)
		Assert(matches.Count() is: 2)
		Assert(matches[0].holes.var is: 'result')
		Assert(matches[0].holes.method is: 'buildTokens')
		Assert(matches[1].holes.var is: 'tokens')
		Assert(matches[1].holes.method is: 'Parse')

		// Find all uses of a specific function
		src = "fn = CombyMatch\r\n" $
			"\tfn2 = CombyMatch.All\r\n" $
			"\tCombyMatch.GetHint(x)\r\n" $
			"\tfn3 = CombyTemplate(y)"
		matches = CombyMatch.All(
			':[name:e] = CombyMatch.:[method:e]', src)
		Assert(matches.Count() is: 1)
		Assert(matches[0].holes.method is: 'All')

		matches = CombyMatch.All(
			'CombyMatch.:[method:e](:[args:e])', src)
		Assert(matches.Count() is: 1)
		Assert(matches[0].holes.method is: 'GetHint')
		Assert(matches[0].holes.args is: 'x')
		// GetDefault calls
		src = "item.GetDefault(#expr, false)\r\n" $
			"\tob.GetDefault(name, 'default')"
		matches = CombyMatch.All(
			':[x:e].GetDefault(:[key:e], :[def:e])', src)
		Assert(matches.Count() is: 2)
		Assert(matches[0].holes.x is: 'item')
		Assert(matches[0].holes.key is: '#expr')
		Assert(matches[0].holes.def is: 'false')
		Assert(matches[1].holes.x is: 'ob')
		Assert(matches[1].holes.key is: 'name')
		Assert(matches[1].holes.def is: "'default'")
		matches = CombyMatch.All(':[x].GetDefault(:[a], :[b])', src)
		Assert(matches.Count() is: 2)

		// Object literal construction: #(1, 2, name: 'Joe')
		src = "return #(pos: 0, end: tokens[next - 1].end,\r\n" $
			"\t\tholes: env.holes)"
		m = CombyMatch.All('#(:[inner])', src)[0]
		Assert(m.holes.inner.Has?('pos:'))
		Assert(m.holes.inner.Has?('holes:'))

		// Record literal: #{name: "Joe"}
		src = "\tresult = #{name: 'Joe', age: 42, active: true}"
		m = CombyMatch.All('#{:[inner]}', src)[0]
		Assert(m.holes.inner.Has?("name: 'Joe'"))
		Assert(m.holes.inner.Has?('active: true'))

		// Expression hole: match a single term, not greedily
		src = "if expr? is true and blockLevel is 0\r\n" $
			"\tbreak"
		m = CombyMatch.All(':[a:e] is :[b:e]', src)[0]
		Assert(m.holes.a is: 'expr?')
		Assert(m.holes.b is: 'true')

		// Chained expression with expression hole
		src = "value.One().Two().Three()"
		m = CombyMatch.All(':[root:e].:[m1:e]().:[m2:e]()', src)[0]
		Assert(m.holes.root is: 'value')
		Assert(m.holes.m1 is: 'One')
		Assert(m.holes.m2 is: 'Two')

		// function definition with parameter defaults
		src = "function Start(search, s, pos = 0, prev = false)\r\n" $
			"\t{ return false }"
		m = CombyMatch.All(
			'function :[name](:[params]) { :[body] }', src)[0]
		Assert(m.holes.name is: 'Start')
		Assert(m.holes.params is: 'search, s, pos = 0, prev = false')
		Assert(m.holes.body is: 'return false')

		// Line comment in body: holes skip comments
		src = "foo // comment\r\n\t.Bar()"
		m = CombyMatch.All(':[a] . :[b]()', src)[0]
		Assert(m.holes.a is: 'foo')
		Assert(m.holes.b is: 'Bar')

		// Block comment inside body
		src = "items /* Merge step */ .Sort()"
		m = CombyMatch.All(':[x] . :[m]()', src)[0]
		Assert(m.holes.x is: 'items')
		Assert(m.holes.m is: 'Sort')

		// while with comparison containing spaces
		src = "while i < tokens.Size() and pos <= start\r\n" $
			"\ti++"
		m = CombyMatch.All('while :[cond]', src)[0]
		Assert(m.holes.cond.Has?('tokens.Size()'))
		Assert(m.holes.cond.Has?('pos <= start'))
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
