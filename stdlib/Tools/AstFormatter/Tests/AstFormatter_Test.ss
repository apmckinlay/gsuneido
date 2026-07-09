// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_expressions()
		{
		.t("123")
		.t(#foo)
		.t("x = y = 123")
		.t("x = #foo $ f(#bar)")
		.t("a + b - c + d")
		.t("x - 1")
		.t("a * b / c")
		.t("x in (1, 2, #foo, #(9))")
		.t("x not in (1, 2)")
		.t("x ? y : z")
		.t("x.y.z")
		.t("x[2]")
		.t("x[2..-1]")
		.t("x[1..]")
		.t("x[..n]")
		.t("x[i::9]")
		.t("x[f() .. g()]")
		.t("x[i+1 ..]") // no space on the absent side
		.t("x[.. i+1]")
		.t("x[i+1 ::]")
		.t("x[i + 1 .. ]", "x[i+1 ..]")
		.t("f()")
		.t("f(1, 2, 3)")
		.t("f(1, 2, a: 3)")
		.t("f(:a, :b)")
		.t("f(a:, b:)")
		.t("f(@args)")
		.t("f(@+1args)")
		.t(".f()")
		.t("x.y()")
		.t("(.f)()")
		.t("[]")
		.t("x = [1, 2]")
		.t("x = [x: 1, y: 2]")
		.t("x = Object(x: 1, y: 2)")
		.t("x = #(1, (2), foo)")
		.t("x = #(a: 1, b:)")
		.t("return false")
		.t(#return)
		.t("return a, b")
		.t("-x")
		.t("not x")
		.t("++x")
		.t("return x--")
		.t("return .x")
		.t('s =~ "^x"')
		.t("b = { it + x }")
		.t("b = {|y| y + x }")
		.t("b = {|y, z| y + z }")
		.t("f = function() { one; two }")
		.t("f = function() { f();; }")
		.t('f = function(x = "foo", y = false) { 123 }',
			"f = function(x = #foo, y = false) { 123 }", norm:)
		.t('x = "line1\nline2"') // multiline string preserved verbatim
		.t("if a and b and c\n\t\tf()")
		}

	Test_statements()
		{
		.t("if x\n\t\tf()")
		.t("if x\n\t\tf()\n\telse\n\t\tg()")
		.t("if x\n\t\tf()\n\telse if y\n\t\tg()\n\telse\n\t\th()")
		.t("if x\n\t\t{\n\t\tf()\n\t\tg()\n\t\t}")
		.t("forever\n\t\tf()")
		.t("while x\n\t\tf()")
		.t("do\n\t\t{\n\t\tf()\n\t\t} while x")
		.t("for x in list\n\t\tf(x)")
		.t("for m, v in ob\n\t\tf(m, v)")
		.t("for (i = 0; i < n; ++i)\n\t\tf(i)")
		.t("try\n\t\tf()\n\tcatch (e)\n\t\tg(e)")
		.t("try\n\t\tf()\n\tcatch (e, 'x')\n\t\tg(e)")
		.t("switch x\n\t\t{\n\tcase 1:\n\t\tf()\n\tcase 2, 3:\n\t\tg()" $
				"\n\tdefault:\n\t\th()\n\t\t}")
		.t("b.Each()\n\t\t{\nPrint(it)\n\t\t}") // debug statements at the margin
		.t("c = class\n\t\t{\n\t\t}")
		}

	Test_unbracing()
		{
		// braces around a single statement are dropped
		.t("if x\n\t\t{\n\t\tf()\n\t\t}", "if x\n\t\tf()", norm:)
		.t("if x\n\t\t{\n\t\tf()\n\t\t}\n\telse\n\t\t{\n\t\tg()\n\t\t}",
			"if x\n\t\tf()\n\telse\n\t\tg()", norm:)
		.t("while x\n\t\t{\n\t\tf()\n\t\t}", "while x\n\t\tf()", norm:)
		.t("while x\n\t\t{\n\t\tbreak\n\t\t}", "while x\n\t\tbreak", norm:)
		.t("forever\n\t\t{\n\t\tf()\n\t\t}", "forever\n\t\tf()", norm:)
		.t("for y in list\n\t\t{\n\t\tg(y)\n\t\t}", "for y in list\n\t\tg(y)", norm:)
		.t("for (i = 0; i < n; ++i)\n\t\t{\n\t\tf(i)\n\t\t}",
			"for (i = 0; i < n; ++i)\n\t\tf(i)", norm:)
		.t("try\n\t\t{\n\t\tf()\n\t\t}\n\tcatch (e)\n\t\t{\n\t\tg(e)\n\t\t}",
			"try\n\t\tf()\n\tcatch (e)\n\t\tg(e)", norm:)
		// else { if ... } rejoins the chain as else-if
		.t("if x\n\t\tf()\n\telse\n\t\t{\n\t\tif y\n\t\t\tg()\n\t\t}",
			"if x\n\t\tf()\n\telse if y\n\t\tg()", norm:)
		// a body ending in an open if/try keeps its braces: the else/catch
		// would otherwise attach to the inner statement
		.t("if x\n\t\t{\n\t\tif y\n\t\t\tf()\n\t\t}\n\telse\n\t\tg()")
		.t("if x\n\t\t{\n\t\twhile y\n\t\t\tif z\n\t\t\t\tf()\n\t\t}\n\telse\n\t\tg()")
		// closed inner chains are safe to unbrace
		.t("if x\n\t\t{\n\t\tif y\n\t\t\tf()\n\t\telse\n\t\t\tg()\n\t\t}" $
				"\n\telse\n\t\th()",
			"if x\n\t\tif y\n\t\t\tf()\n\t\telse\n\t\t\tg()\n\telse\n\t\th()", norm:)
		.t("try\n\t\t{\n\t\tif y\n\t\t\tf()\n\t\t}\n\tcatch (e)\n\t\tg(e)",
			"try\n\t\tif y\n\t\t\tf()\n\tcatch (e)\n\t\tg(e)", norm:)
		// comments that would be displaced keep the braces
		.t("if x\n\t\t{\n\t\tf() // ok\n\t\t}")
		.t("if x\n\t\t{\n\t\tf()\n\t\t// note\n\t\t}")
		// a leading comment moves with the statement
		.t("if x\n\t\t{\n\t\t// note\n\t\tf()\n\t\t}", "if x\n\t\t// note\n\t\tf()",
			norm:)
		// blank lines inside the braces go with them
		.t("if x\n\t\t{\n\n\t\tf()\n\t\t}", "if x\n\t\tf()", norm:)
		// multiple statements and empty bodies keep their braces
		.t("if x\n\t\t{\n\t\tf()\n\t\tg()\n\t\t}")
		.t("if x\n\t\t{\n\t\t}")
		}

	Test_longStrings()
		{
		// a plain string that cannot fit splits after a word, joined with $
		.t('x = "aa bb cc dd"', 'x = "aa bb " $\n\t\t"cc dd"', width: 20, norm:)
		// the split form is stable, and pieces never become #symbols
		.t('x = "aa bb " $\n\t\t"cc dd"', width: 20)
		// no word boundary: left to overflow
		.t('x = "' $ "a-".Repeat(15) $ '"', width: 20)
		// splits anywhere in an expression, at any depth
		long = "word ".Repeat(20).Trim()
		out = AstFormatter('function ()\n\t{\n\tf(name: "' $ long $ '")\n\t}\n')
		for line in out.Lines()
			Assert(line.Detab().Size() <= 90)
		Assert(AstFormatter(out) is: out)
		}

	Test_spacing()
		{
		// a simple arithmetic chain tightens as the operand of a looser
		// operator; on its own, in assignments, and in call arguments it
		// keeps its spaces
		.t("x = a * b + c", "x = a*b + c")
		.t("x = a * b + c * d", "x = a*b + c*d")
		.t("if i < n - 1\n\t\tf()", "if i < n-1\n\t\tf()")
		.t("x = i % 2 is 0", "x = i%2 is 0")
		.t("x = n + 1 in (1, 2)", "x = n+1 in (1, 2)")
		.t("x = ob[i + 1]", "x = ob[i+1]")
		.t("x = s[i + 1 .. j - 1]", "x = s[i+1 .. j-1]")
		.t("x = s[i + 1 :: n]", "x = s[i+1 :: n]")
		.t("x = a * b")
		.t("x = a + b - c")
		.t("f(a * b)")
		.t("x = a * f() + c") // a call operand keeps the chain spaced
		.t("x = a - -b is c") // no chance of a--b
		.t("x = a $ b is c") // $ never tightens
		// $ and + are EQUAL precedence (13) in the interpreter's table:
		// the Add nested by associativity must stay spaced
		.t("x = a - b $ c")
		.t("x = a $ b + c*d") // the Mul (14) under the + (13) stays tight
		.t("x = a + b % c", "x = a + b%c") // % (14) under + (13) tightens
		}

	Test_cleanups()
		{
		.t("try\n\t\tf()\n\tcatch (unused)\n\t\tg()", "try\n\t\tf()\n\tcatch\n\t\tg()",
			norm:)
		.t("try\n\t\tf()\n\tcatch (unused, 'x')\n\t\tg()") // pattern keeps var
		// super.New(x) => super(x) is untestable here: explicit super.New
		// never compiles, and AstFormatter asserts compilable input
		.t("x = y ? true : false", "x = y", norm:)
		.t("x = y ? false : true", "x = not y", norm:)
		.t("x = y is z ? false : true", "x = y isnt z", norm:)
		.t("x = a or b ? false : true", "x = not (a or b)", norm:)
		// not (...) is always left as written
		.t("if not (a is b)\n\t\tf()")
		.t("x = not (a is b)")
		.t("s = t[0..n]", "s = t[..n]", norm:)
		.t("s = t[0::n]", "s = t[::n]", norm:)
		.t("s = t[0..]") // no upper bound: kept
		.t("c = class : Foo\n\t\t{\n\t\t}", "c = Foo\n\t\t{\n\t\t}", norm:)
		}

	Test_quotes()
		{
		// 'c' for characters, #word for one-word strings, "..." for the rest
		.t("x = 'a'")
		.t('x = "a"', "x = 'a'", norm:)
		.t("x = 'foo'", "x = #foo", norm:)
		.t('x = "hello world"')
		.t("x = 'hello world'", 'x = "hello world"', norm:)
		.t("f('one', 'two words')", 'f(#one, "two words")', norm:)
		.t("x = ''", 'x = ""', norm:)
		.t("x = 'foo?'", "x = #foo?", norm:)
		.t("x = 'Global'", "x = #Global", norm:)
		// keywords and _names make confusing symbols; keep them strings
		.t("x = 'class'", 'x = "class"', norm:)
		.t("x = '_foo'", 'x = "_foo"', norm:)
		// swap the quote kind rather than escape
		.t(`x = 'say "hi"'`)
		.t(`x = "don't"`)
		.t(`x = "say 'hi'"`)
		// backquotes, escape sequences, and multiline strings stay as written
		.t("x = `raw string`")
		.t("x = 'a\\tb'")
		// constants: identifiers are already bare, others follow the rules
		.t("x = #(foo, 'a b')", 'x = #(foo, "a b")', norm:)
		// in a constant, a $ chain may only hold plain strings, never #syms
		.t("x = #{A: 'hello' $ 'world'}", 'x = #{A: "hello" $ "world"}', norm:)
		// escaped values stay as written, without dragging the comma along
		.t("x = #('a\\tb', 'c d')", "x = #('a\\tb', \"c d\")", norm:)
		.t("c = class\n\t\t{\n\t\tTitle: 'Report'\n\t\t}",
			"c = class\n\t\t{\n\t\tTitle: #Report\n\t\t}", norm:)
		}

	Test_debugStatements()
		{
		// debug calls in statement position go to the left margin
		.t("function ()\n\t{\n\tx = 1\n\tPrint(:x)\n\t}",
			"function()\n\t{\n\tx = 1\nPrint(:x)\n\t}", wrap: false)
		.t("function ()\n\t{\n\tTracePrint(a, b)\n\t}",
			"function()\n\t{\nTracePrint(a, b)\n\t}", wrap: false)
		// a single-line body containing one is forced open
		.t("f = function() { Print(s) }", "f = function()\n\t\t{\nPrint(s)\n\t\t}")
		// control bodies: unbraced and margined
		.t("if x\n\t\tPrint(y)", "if x\nPrint(y)")
		.t("if x\n\t\t{\n\t\tPrint(y)\n\t\t}", "if x\nPrint(y)", norm:)
		// a trailing comment rides along; the braces stay for it
		.t("if x\n\t\t{\n\t\tPrint(y) // dbg\n\t\t}",
			"if x\n\t\t{\nPrint(y) // dbg\n\t\t}")
		// only bare calls to the debug functions, only in statement position
		.t("x.Print()")
		.t("x = Print(1)")
		.t("Printer(x)")
		}

	Test_classes()
		{
		.t("// Copyright (C) 2026 Suneido Software Corp.\nclass\n\t{\n\tX: 1" $
				"\n\tCallClass()\n\t\t{\n\t\treturn .X\n\t\t}\n\t}", wrap: false)
		.t("class\n\t{\n\tF()\n\t\t{\n\t\treturn 1\n\t\t}\n\n\tG()\n\t\t{" $
				"\n\t\treturn 2\n\t\t}\n\t}", wrap: false)
		// a blank line is added after every method
		.t("class\n\t{\n\tF()\n\t\t{\n\t\t}\n\tG()\n\t\t{\n\t\t}\n\t}",
			"class\n\t{\n\tF()\n\t\t{\n\t\t}\n\n\tG()\n\t\t{\n\t\t}\n\t}", wrap: false)
		// the inserted blank goes BEFORE comments leading the next method
		.t("class\n\t{\n\tF()\n\t\t{\n\t\t}\n\t// about G\n\tG()\n\t\t{\n\t\t}\n\t}",
			"class\n\t{\n\tF()\n\t\t{\n\t\t}\n\n\t// about G\n\tG()\n\t\t{\n\t\t}\n\t}",
			wrap: false)
		.t("class\n\t{\n\tF()\n\t\t{\n\t\t}\n\tX: 1\n\t}",
			"class\n\t{\n\tF()\n\t\t{\n\t\t}\n\n\tX: 1\n\t}", wrap: false)
		.t("Base\n\t{\n\tOp: (a: 1, b: 2)\n\t}", wrap: false)
		}

	Test_crlf()
		{
		Assert(AstFormatter("function ()\r\n\t{\r\n\tx = 1\r\n\t}\r\n")
			is: "function()\n\t{\n\tx = 1\n\t}\n")
		src = "function()\n\t{\n\tx = 'a\r\nb'\n\t}\n"
		Assert(AstFormatter(src) is: src)
		}

	Test_normalizations()
		{
		.t("x--\n\ty = 1", "--x\n\ty = 1", norm:)
		.t("x++\n\ty = 1", "++x\n\ty = 1", norm:)
		.t("x++")
		.t("b = {|unused| x++ }")
		.t("f(a: a, b: b)", "f(:a, :b)", norm:)
		.t("f(1\n\t\t2)", "f(1, 2)", norm:) // newline-as-comma gets the comma
		.t("Object(1, 2)", "[1, 2]", norm:)
		.t("Record(x: 1)", "[x: 1]", norm:)
		.t("x = Record()", "x = []", norm:)
		.t("x = Object()") // no unnamed args: not bracketable
		.t("x = Object(a: 1)")
		.t("this.x", ".x", norm:)
		.t('x = y["z"]', "x = y.z", norm:)
		.t("x = this['y']", "x = .y", norm:)
		.t("return throw x") // throws if the caller discards the result
		.t("x = #(a: [b: 1], c: [2], d: [])")
		.t("x = #(Filters: (['a'], []))", "x = #(Filters: ([a], []))", norm:)
		.t("x = #{a: 1}")
		.t("Assert(f(x) is: false)")
		.t("Assert(f(x), is: false)")
		.t("f(0, :a, :b)") // dropping THIS comma would reparse as f(0: a, ...)
		.t("x = #(aa: 1,\n\t\tbb: 2)", "x = #(aa: 1, bb: 2)")
		.t("class\n\t{\n\tF()\n\t\t{\n\t\treturn this.x\n\t\t}\n\t}", wrap: false)
		.t("class\n\t{\n\tF()\n\t\t{\n\t\treturn this['x']\n\t\t}\n\t}",
			"class\n\t{\n\tF()\n\t\t{\n\t\treturn this.x\n\t\t}\n\t}", wrap: false, norm:)
		.t("class\n\t{\n\tF()\n\t\t{\n\t\treturn this.X\n\t\t}\n\t}",
			"class\n\t{\n\tF()\n\t\t{\n\t\treturn .X\n\t\t}\n\t}", wrap: false, norm:)
		.t("class\n\t{\n\tF()\n\t\t{\n\t\treturn .x\n\t\t}\n\t}", wrap: false)
		.t('x = #("foo")', "x = #(foo)", norm:)
		.t("for (;;)\n\t\tf()", "forever\n\t\tf()", norm:)
		.t("while true\n\t\tf()", "forever\n\t\tf()", norm:)
		.t("while (true)\n\t\tf()", "forever\n\t\tf()", norm:)
		.t("Base\n\t{\n\tOp: #(a: 1, b: 2)\n\t}", "Base\n\t{\n\tOp: (a: 1, b: 2)\n\t}",
			wrap: false, norm:)
		.t("switch (x)\n\t\t{\n\tcase 1:\n\t\treturn 2\n\t\t}",
			"switch x\n\t\t{\n\tcase 1:\n\t\treturn 2\n\t\t}", norm:)
		.t("switch (x $ y)\n\t\t{\n\tcase #ab:\n\t\treturn 2\n\t\t}") // not plain: kept
		}

	Test_verticalTables()
		{
		// one member per line in the source keeps that shape, without commas
		.t("x = #(\n\t\t(a, 'one two'),\n\t\t(b, 'three four'))",
			'x = #(\n\t\t(a, "one two")\n\t\t(b, "three four"))', norm:)
		// key: scalar rows align values to the longest key
		.t("x = #(\n\t\tLEFT: 0x0001\n\t\tCENTERX: 0x0002)",
			"x = #(\n\t\tLEFT:    0x0001\n\t\tCENTERX: 0x0002)")
		// blank lines group rows; hex and exponent forms survive
		.t("x = #(\n\t\ta: 0xff\n\n\t\tbb: 1e3)", "x = #(\n\t\ta:  0xff\n\n\t\tbb: 1e3)")
		// wrapped source (members sharing a line) reflows with fill
		.t("x = #(1, 2,\n\t\t3)", "x = #(1, 2, 3)")
		// a key-only row keeps its comma: without it the next member
		// would be absorbed as its value
		.t("x = #(\n\t\tabc:,\n\t\tde: 1)", "x = #(\n\t\tabc:,\n\t\tde:  1)")
		}

	Test_comments()
		{
		.t("x = 1 // note")
		.t("// note\n\tx = 1")
		.t("x = 1 /* c */")
		.t("/* c */ x = 1")
		.t("x = 1\n\n\ty = 2") // blank line preserved
		.t("x = 1 // one\n\n\t// two\n\ty = 2")
		.t("b = {|x/*unused*/, y| y }") // annotations hug their token
		.t("b = {|x /*unused*/, y| y }", "b = {|x/*unused*/, y| y }")
		.t("f(2/*=default*/)")
		}

	Test_width()
		{
		.t("f(aaaa, bbbb, cccc, dddd, eeee, ffff)",
			"f(aaaa, bbbb, cccc, dddd, eeee,\n\t\tffff)", width: 40)
		.t("s = aaaaaa $ bbbbbb $ cccccc $ dddddd",
			"s = aaaaaa $ bbbbbb $ cccccc $\n\t\tdddddd", width: 40)
		.t("x = cond ? someValueOne : someValueTwo",
			"x = cond\n\t\t? someValueOne\n\t\t: someValueTwo", width: 40)
		.t("x = [aaaa: 1, bbbb: 2, cccc: 3, dddd: 4]",
			"x = [aaaa: 1, bbbb: 2,\n\t\tcccc: 3, dddd: 4]", width: 30)
		}

	t(src, expected = false, wrap = true, width = false, norm = false)
		{
		if wrap
			{
			src = "function()\n\t{\n\t" $ src $ "\n\t}\n"
			if expected isnt false
				expected = "function()\n\t{\n\t" $ expected $ "\n\t}\n"
			}
		else
			{
			src = src $ '\n'
			if expected isnt false
				expected = expected $ '\n'
			}
		if expected is false
			expected = src
		actual = width is false ? AstFormatter(src) : AstFormatter(src, :width)
		Assert(actual is: expected)
		again = width is false ? AstFormatter(actual) : AstFormatter(actual, :width)
		Assert(again is: actual, msg: "not idempotent")
		if not norm
			Assert(.toks(actual) is: .toks(src), msg: "tokens changed")
		}

	toks(s)
		{
		ob = Object()
		scan = Scanner(s)
		while scan isnt (tok = scan.Next2())
			if tok not in (#NEWLINE, #WHITESPACE, #COMMENT)
				ob.Add(scan.Text())
		return ob
		}
	}
