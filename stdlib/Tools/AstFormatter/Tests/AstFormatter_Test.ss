// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Suneido.AstFormatter_Test = Object()
		.AddTeardown({ Suneido.Delete(#AstFormatter_Test) })
		wg = WaitGroup()
		AstFormatter_Test.
			Members().
			Filter({ it.Prefix?(#UnitTest_) }).
			Each(
				{|m|
				wg.Thread(
					{
					try
						AstFormatter_Test[m]()
					catch (err)
						Suneido.AstFormatter_Test.Add(m $ ": " $ err)
					})
				})
		wg.Wait()
		Assert(Suneido.AstFormatter_Test isSize: 0)
		}

	UnitTest_expressions()
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

	UnitTest_statements()
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

	UnitTest_unbracing()
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
		// a BARE loop as the then-branch: the guard must reach the loop body so
		// its trailing if keeps its braces and the else stays on the outer if
		.t("if x\n\t\tfor y in list\n\t\t\t{\n\t\t\tif z\n\t\t\t\tf()\n\t\t\t}\n\telse" $
				"\n\t\tg()")
		.t("if x\n\t\twhile y\n\t\t\t{\n\t\t\tif z\n\t\t\t\tf()\n\t\t\t}\n\telse" $
				"\n\t\tg()")
		// no trailing else: the same bare-loop body unbraces normally
		.t("if x\n\t\tfor y in list\n\t\t\t{\n\t\t\tif z\n\t\t\t\tf()\n\t\t\t}",
			"if x\n\t\tfor y in list\n\t\t\tif z\n\t\t\t\tf()", norm:)
		}

	UnitTest_longStrings()
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

	UnitTest_multilineCat()
		{
		// a $ chain containing a multiline string keeps its authored layout,
		// including hand-outdented continuations after the string
		.t('x = "l1\nl2" $\n\t\t"tail"')
		.t('x = "head" $ "l1\nl2"')
		.t("f(`a\nb` $\n\t`tail`)")
		// a trailing implicit-this member after a multiline string must keep its
		// name: the raw-slice span end has to reach past the dot (spanEnd)
		.t(".t = F('a\n\t\tb ' $ .m1)")
		.t(".t = F('line1\n\t\tline2 ' $ .aTable $\n\t\t' more ' $ .aTable2)")
		}

	UnitTest_spacing()
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
		.t("Test\n\t{\n\tx: `aaa\nbbb` $ `ccc`\n\t}", wrap: false)
		.t("(x: `aaa\nbbb` $ `ccc`)", "#(x: `aaa\nbbb` $ `ccc`)", wrap: false, norm:)
		}

	UnitTest_cleanups()
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
		// a comment hugging the 0 keeps the index, else it is orphaned onto ::
		.t("s = t[0/*=b*/::2]")
		.t("s = t[0/*=b*/..2]")
		.t("c = class : Foo\n\t\t{\n\t\t}", "c = Foo\n\t\t{\n\t\t}", norm:)
		}

	UnitTest_quotes()
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
		// backquotes stay as written; multiline strings default to "
		.t("x = `raw string`")
		.t('x = "line1\nline2"')
		.t("x = 'a +\nb'", 'x = "a +\nb"', norm:)
		.t('x = \'say "hi"\nbye\'') // a " inside pins it as written
		// escape-carrying strings swap delimiters by the same rules: escapes
		// move verbatim, but any quote character inside pins them as written
		.t("x = 'a\\tb'", 'x = "a\\tb"', norm:)
		.t("x = '{\\n'", 'x = "{\\n"', norm:)
		.t('x = "\\n"', "x = '\\n'", norm:) // a single character
		.t(`x = 'don\'t \n'`)
		// constants: identifiers are already bare, others follow the rules
		.t("x = #(foo, 'a b')", 'x = #(foo, "a b")', norm:)
		// in a constant, a $ chain may only hold plain strings, never #syms
		.t("x = #{A: 'hello' $ 'world'}", 'x = #{A: "hello" $ "world"}', norm:)
		// requoted escaped values don't drag the member comma along
		.t("x = #('a\\tb', 'c d')", 'x = #("a\\tb", "c d")', norm:)
		.t("c = class\n\t\t{\n\t\tTitle: 'Report'\n\t\t}",
			"c = class\n\t\t{\n\t\tTitle: #Report\n\t\t}", norm:)
		}

	UnitTest_debugStatements()
		{
		// a debug call written at the left margin stays there
		.t("function ()\n\t{\n\tx = 1\nPrint(:x)\n\t}",
			"function()\n\t{\n\tx = 1\nPrint(:x)\n\t}", wrap: false)
		.t("function ()\n\t{\nTracePrint(a, b)\n\t}",
			"function()\n\t{\nTracePrint(a, b)\n\t}", wrap: false)
		.t("if x\nPrint(y)")
		.t("b.Each()\n\t\t{\nPrint(it)\n\t\t}")
		// one written with the code is laid out like any other call
		.t("function ()\n\t{\n\tx = 1\n\tPrint(:x)\n\t}",
			"function()\n\t{\n\tx = 1\n\tPrint(:x)\n\t}", wrap: false)
		.t("if x\n\t\tPrint(y)")
		.t("if x\n\t\t{\n\t\tPrint(y)\n\t\t}", "if x\n\t\tPrint(y)", norm:)
		.t("f = function() { Print(s) }")
		.t("b.Each()\n\t\t{\n\t\tPrint(it)\n\t\t}")
		// a trailing comment keeps the braces
		.t("if x\n\t\t{\n\t\tPrint(y) // dbg\n\t\t}")
		// only bare calls to the debug functions, only in statement position
		.t("x.Print()")
		.t("x = Print(1)")
		.t("Printer(x)")
		}

	UnitTest_blockArgs()
		{
		// a block argument indents from the statement, not the argument list
		.t("Win.SetTimeout(\n\t\t{\n\t\t.f(a)\n\t\t}, 0)")
		.t("Win.SetTimeout(\n\t\t{\n\t\t.f(a)\n\t\t})")
		.t("ob.Each(\n\t\t{\n\t\t.f(it)\n\t\t}, 1)")
		// a block written last, inside the parens: the break before it is the
		// formatter's own, so the next pass must not read it as a vertical list
		.t("Win.SetTimeout(0,\n\t\t{\n\t\t.f(a)\n\t\t})")
		.t("lines.Fold(0,\n\t\t{|sum, line|\n\t\tsum + line\n\t\t})")
		.t("Win.SetTimeout(0, { .f(a) })")
		// the trailing block syntax is unaffected
		.t("Win.SetTimeout(0)\n\t\t{\n\t\t.f(a)\n\t\t}")
		// a /*=x*/ annotation hugs its token, prose keeps its space
		.t("x = 1/*=Annotation*/")
		.t("x = 1 /*= annotation */")
		}

	UnitTest_classes()
		{
		.t("// Copyright (C) 2026 Suneido Software Corp.\nclass\n\t{\n\tX: 1" $
				"\n\tCallClass()\n\t\t{\n\t\treturn .X\n\t\t}\n\t}", wrap: false)
		.t("class\n\t{\n\tF()\n\t\t{\n\t\treturn 1\n\t\t}\n\n\tG()\n\t\t{" $
				"\n\t\treturn 2\n\t\t}\n\t}", wrap: false)
		// empty method bodies collapse to { }; a blank line follows every method
		.t("class\n\t{\n\tF()\n\t\t{\n\t\t}\n\tG()\n\t\t{\n\t\t}\n\t}",
			"class\n\t{\n\tF() { }\n\n\tG() { }\n\t}", wrap: false)
		// the inserted blank goes BEFORE comments leading the next method
		.t("class\n\t{\n\tF()\n\t\t{\n\t\t}\n\t// about G\n\tG()\n\t\t{\n\t\t}\n\t}",
			"class\n\t{\n\tF() { }\n\n\t// about G\n\tG() { }\n\t}", wrap: false)
		.t("class\n\t{\n\tF()\n\t\t{\n\t\t}\n\tX: 1\n\t}",
			"class\n\t{\n\tF() { }\n\n\tX: 1\n\t}", wrap: false)
		.t("Base\n\t{\n\tOp: (a: 1, b: 2)\n\t}", wrap: false)
		// scalar members align into a value column; longest key sets the width
		.t("class\n\t{\n\tx: 1\n\tabcd: 2\n\t}", "class\n\t{\n\tx:    1\n\tabcd: 2\n\t}",
			wrap: false)
		// a method breaks the aligned run; each contiguous run aligns on its own
		.t("class\n\t{\n\taa: 1\n\tb: 2\n\tGo()\n\t\t{\n\t\t.x()\n\t\t}\n\t}",
			"class\n\t{\n\taa: 1\n\tb:  2\n\tGo()\n\t\t{" $ "\n\t\t.x()\n\t\t}\n\t}",
			wrap: false)
		// a complex (nested) value breaks the run
		.t("class\n\t{\n\taa: 1\n\tData: (x: 1)\n\tb: 2\n\t}", wrap: false)
		// number members keep their source form
		.t("class\n\t{\n\tA: 1\n\tBb: 0xff\n\tCcc: 100_000_000\n\t}",
			"class\n\t{\n\tA:   1\n\tBb:  0xff\n\tCcc: 100_000_000\n\t}", wrap: false)
		}

	UnitTest_emptyBody()
		{
		// an empty body has nothing to lay out, so it hugs the signature as { }
		.t("function() { }", wrap: false)
		.t("function() {}", "function() { }", wrap: false)
		.t("function()\n\t{\n\t}", "function() { }", wrap: false)
		.t("class\n\t{\n\tF() {}\n\t}", "class\n\t{\n\tF() { }\n\t}", wrap: false)
		.t("class\n\t{\n\tF(a, b) {}\n\t}", "class\n\t{\n\tF(a, b) { }\n\t}", wrap: false)
		// a body with a statement keeps the Whitesmiths block
		.t("class\n\t{\n\tF()\n\t\t{\n\t\tg()\n\t\t}\n\t}", wrap: false)
		// comments are not an empty body - they hold the block open
		.t("class\n\t{\n\tF()\n\t\t{\n\t\t// keep me\n\t\t}\n\t}", wrap: false)
		.t("class\n\t{\n\tF()\n\t\t{ // note\n\t\t}\n\t}", wrap: false)
		// a comment after the body stays out of it (Trailing must not eat the })
		.t("class\n\t{\n\tF() { }\n\n\t// about G\n\tG() { }\n\t}", wrap: false)
		// only #Function bodies collapse: class and control-flow bodies are unchanged
		.t("c = class\n\t\t{\n\t\t}")
		.t("if x\n\t\t{\n\t\t}")
		}

	UnitTest_crlf()
		{
		Assert(AstFormatter("function ()\r\n\t{\r\n\tx = 1\r\n\t}\r\n")
			is: "function()\n\t{\n\tx = 1\n\t}\n")
		src = "function()\n\t{\n\tx = 'a\r\nb'\n\t}\n"
		Assert(AstFormatter(src) is: "function()\n\t{\n\tx = \"a\r\nb\"\n\t}\n")
		}

	UnitTest_normalizations()
		{
		.t("x--\n\ty = 1", "x -= 1\n\ty = 1", norm:)
		.t("x++\n\ty = 1", "x += 1\n\ty = 1", norm:)
		.t("x++")
		.t("b = {|unused| x++ }")
		.t("f(a: a, b: b)", "f(:a, :b)", norm:)
		.t("f(1\n\t\t2)", "f(1,\n\t\t2)", norm:) // newline-as-comma gets the comma
		.t("Object(1, 2)", "[1, 2]", norm:)
		.t("Record(x: 1)", "[x: 1]", norm:)
		.t("x = Record()", "x = []", norm:)
		.t("x = Object()") // no unnamed args: not bracketable
		.t("x = Object(a: 1)")
		// [...] would flip the type, so these keep the call form
		.t("Record(1, 2)") // [1, 2] is an Object, Record(1, 2) is a Record
		.t("Record(9, a: 1)") // [9, a: 1] is an Object (has a positional member)
		.t("Object(9, a: 1)", "[9, a: 1]", norm:) // mixed Object -> Object literal
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
		.t("x = #(aa: 1,\n\t\tbb: 2)") // a source break between members is kept
		.t("class\n\t{\n\tF()\n\t\t{\n\t\treturn this.x\n\t\t}\n\t}", wrap: false)
		.t("class\n\t{\n\tF()\n\t\t{\n\t\treturn this['x']\n\t\t}\n\t}",
			"class\n\t{\n\tF()\n\t\t{\n\t\treturn this.x\n\t\t}\n\t}", wrap: false, norm:)
		.t("class\n\t{\n\tF()\n\t\t{\n\t\treturn this.X\n\t\t}\n\t}",
			"class\n\t{\n\tF()\n\t\t{\n\t\treturn .X\n\t\t}\n\t}", wrap: false, norm:)
		.t("class\n\t{\n\tF()\n\t\t{\n\t\treturn .x\n\t\t}\n\t}", wrap: false)
		// members named this* keep their privatization
		.t("class\n\t{\n\tF()\n\t\t{\n\t\treturn .thisx\n\t\t}\n\t}", wrap: false)
		.t("class\n\t{\n\tF()\n\t\t{\n\t\treturn this.thisx\n\t\t}\n\t}", wrap: false)
		.t('x = #("foo")', "x = #(foo)", norm:)
		.t("for (;;)\n\t\tf()", "forever\n\t\tf()", norm:)
		.t("while true\n\t\tf()", "forever\n\t\tf()", norm:)
		.t("while (true)\n\t\tf()", "forever\n\t\tf()", norm:)
		.t("Base\n\t{\n\tOp: #(a: 1, b: 2)\n\t}", "Base\n\t{\n\tOp: (a: 1, b: 2)\n\t}",
			wrap: false, norm:)
		.t("switch (x)\n\t\t{\n\tcase 1:\n\t\treturn 2\n\t\t}",
			"switch x\n\t\t{\n\tcase 1:\n\t\treturn 2\n\t\t}", norm:)
		.t("switch (x $ y)\n\t\t{\n\tcase #ab:\n\t\treturn 2\n\t\t}") // not plain: kept
		// digit grouping: integers of 5+ digits get _ every 3, hex every 4
		.t("x = 1234")
		.t("x = 12345", "x = 12_345", norm:)
		.t("x = 1_0000", "x = 10_000", norm:)
		.t("x = -123456", "x = -123_456", norm:)
		.t("x = 0xff")
		.t("x = 0x12349876", "x = 0x1234_9876", norm:)
		.t("x = 1e3") // not a plain integer: untouched
		.t("x = #(a: 12345)", "x = #(a: 12_345)", norm:)
		.t("class\n\t{\n\tX: 12345\n\t}", "class\n\t{\n\tX: 12_345\n\t}", wrap: false,
			norm:)
		.t("(a: 1\nFoo()\n\t{\n\treturn 5\n\t})", "#(a: 1,\n\tFoo, (),\n\t{return, 5})",
			wrap: false, norm:)
		}

	UnitTest_verticalTables()
		{
		// one member per line in the source keeps that shape, with commas
		.t("x = #(\n\t\t(a, 'one two'),\n\t\t(b, 'three four'))",
			'x = #(\n\t\t(a, "one two"),\n\t\t(b, "three four")\n\t\t)', norm:)
		// key: scalar rows align values to the longest key
		.t("x = #(\n\t\tLEFT: 0x0001\n\t\tCENTERX: 0x0002)",
			"x = #(\n\t\tLEFT:    0x0001,\n\t\tCENTERX: 0x0002\n\t\t)", norm:)
		// blank lines group rows; hex and exponent forms survive
		.t("x = #(\n\t\ta: 0xff\n\n\t\tbb: 1e3)",
			"x = #(\n\t\ta:  0xff,\n\n\t\tbb: 1e3\n\t\t)", norm:)
		// a key-only row's comma glues to the ':'
		.t("x = #(\n\t\tabc:,\n\t\tde: 1)", "x = #(\n\t\tabc:,\n\t\tde:  1\n\t\t)")
		// the close delimiter takes its own line, indented with the members
		.t("x = #(\n\t\ta,\n\t\tb,\n\t\tc)", "x = #(\n\t\ta,\n\t\tb,\n\t\tc\n\t\t)")
		// a comment before the close stays above it
		.t("x = #(\n\t\ta,\n\t\tb\n\t\t// last\n\t\t)")
		// a nested table closes at its own indent
		.t("x = #(\n\t\ta: (\n\t\t\tb,\n\t\t\tc),\n\t\td: 1)",
			"x = #(\n\t\ta: (\n\t\t\tb,\n\t\t\tc\n\t\t\t),\n\t\td: 1\n\t\t)")
		// a record literal closes the same way
		.t("x = #{\n\t\ta: 1,\n\t\tb: 2}", "x = #{\n\t\ta: 1,\n\t\tb: 2\n\t\t}")
		// a break between members in the source is kept: fill preserves the layout
		.t("x = #(1, 2,\n\t\t3)")
		// a packed header with a newline-separated body stays that way, no comments
		.t("x = #(A, b, name: 1,\n\t\tc: 2,\n\t\td: 3)")
		// a blank line between members survives the fill
		.t("x = #(a: 1,\n\t\tb: 2,\n\n\t\tc: 3)")
		}

	UnitTest_chains()
		{
		// a chain of dotted calls that cannot fit breaks after each dot
		.t("x = s.Aaaaaaaaaaaaaaaa(111_111_111).Bbbbbbbbbbbbbbbb(222_222_222)" $
				".Cccccccccccccccc(333_333_333)",
			"x = s.\n\t\tAaaaaaaaaaaaaaaa(111_111_111)." $
				"\n\t\tBbbbbbbbbbbbbbbb(222_222_222)." $
				"\n\t\tCccccccccccccccc(333_333_333)")
		// a broken chain that fits joins to one line
		.t("x = s.\n\t\tF(1).\n\t\tG(2)", "x = s.F(1).G(2)")
		// an implicit-this chain keeps its first segment glued
		.t("return .Aaaaaaaaaaaaaaaa(111_111_111).Bbbbbbbbbbbbbbbb(222_222_222)" $
				".Cccccccccccccccc(333_333_333)",
			"return .Aaaaaaaaaaaaaaaa(111_111_111).\n\t\tBbbbbbbbbbbbbbbb(222_222_222)." $
				"\n\t\tCccccccccccccccc(333_333_333)")
		// a trailing block argument keeps the chain out of the group
		.t("ob.F(aaa).Each()\n\t\t{\nPrint(it)\n\t\t}")
		}

	UnitTest_tableAlign()
		{
		// uniform rows of scalars align as a grid
		.t("x = #(\n\t\t(a: 1, b: 200)\n\t\t(a: 22, b: 3)\n\t\t(a: 333, b: 44))",
			"x = #(\n\t\t(a: 1,   b: 200),\n\t\t(a: 22,  b: 3),\n\t\t(a: 333, b: 44)" $
				"\n\t\t)", norm:)
		// authored commas between rows are kept
		.t("x = #(\n\t\t(rate: 0, min: 0),\n\t\t(rate: 2.5, min: 1348.01),\n\t\t" $
				"(rate: 3, min: 2696.01))",
			"x = #(\n\t\t(rate: 0,   min: 0),\n\t\t(rate: 2.5, min: 1348.01),\n\t\t" $
				"(rate: 3,   min: 2696.01)\n\t\t)")
		// number source forms survive in cells
		.t("x = #(\n\t\t(n: 0xff, m: 1)\n\t\t(n: 2, m: 22)\n\t\t(n: 3, m: 4))",
			"x = #(\n\t\t(n: 0xff, m: 1),\n\t\t(n: 2,    m: 22),\n\t\t(n: 3,    m: 4)" $
				"\n\t\t)", norm:)
		// mixed scalar kinds share a column
		.t("x = #(\n\t\t(v: false, n: 1)\n\t\t(v: 0.25, n: 22)\n\t\t(v: hello, n: 333))",
			"x = #(\n\t\t(v: false, n: 1),\n\t\t(v: 0.25,  n: 22)," $
				"\n\t\t(v: hello, n: 333)\n\t\t)", norm:)
		// the payroll shape: a class member table
		.t("Base\n\t{\n\tR: #(\n\t\t(a: 1, b: 200)\n\t\t(a: 22, b: 3)\n\t\t" $
				"(a: 333, b: 44))\n\t}",
			"Base\n\t{\n\tR: (\n\t\t(a: 1,   b: 200),\n\t\t(a: 22,  b: 3),\n\t\t" $
				"(a: 333, b: 44)\n\t\t)\n\t}", wrap: false, norm:)
		// different keys: no grid
		.t("x = #(\n\t\t(a: 1, b: 2)\n\t\t(alpha: 123, beta: false)\n\t\t(a: 5, b: 6))",
			"x = #(\n\t\t(a: 1, b: 2),\n\t\t(alpha: 123, beta: false)," $
				"\n\t\t(a: 5, b: 6)\n\t\t)", norm:)
		// fewer than 3 rows: no grid
		.t("x = #(\n\t\t(a: 1, b: 2)\n\t\t(a: 11, b: 22))",
			"x = #(\n\t\t(a: 1, b: 2),\n\t\t(a: 11, b: 22)\n\t\t)", norm:)
		// a comment inside a row: no grid, comment kept
		.t("x = #(\n\t\t(a: 1, b: 2)\n\t\t(a: 11 /*note*/, b: 22)\n\t\t(a: 5, b: 6))",
			"x = #(\n\t\t(a: 1, b: 2),\n\t\t(a: 11 /*note*/, b: 22),\n\t\t(a: 5, b: 6)" $
				"\n\t\t)", norm:)
		// key-only cells: no grid
		.t("x = #(\n\t\t(a:, b: 1)\n\t\t(a:, b: 22)\n\t\t(a:, b: 333))",
			"x = #(\n\t\t(a:, b: 1),\n\t\t(a:, b: 22),\n\t\t(a:, b: 333)\n\t\t)", norm:)
		// too wide to pad: no grid
		.t("x = #(\n\t\t(aaaaaaa: 1, b: 22_222_222)\n\t\t" $
				"(aaaaaaa: 11_111_111, b: 2)\n\t\t" $ "(aaaaaaa: 1, b: 2))",
			"x = #(\n\t\t(aaaaaaa: 1, b: 22_222_222),\n\t\t" $
				"(aaaaaaa: 11_111_111, b: 2),\n\t\t" $ "(aaaaaaa: 1, b: 2)\n\t\t)",
			width: 40, norm:)
		}

	UnitTest_comments()
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
		// a comment inside a scalar member slice is not re-emitted
		.t("class\n\t{\n\tx: (style: 0x0100 /*ES.NOHIDESEL*/, y: 1)\n\t}", wrap: false)
		// a margin // comment stays at the margin, even last in a body
		.t("x = 1\n//\tf()\n\ty = 2")
		.t("x = 1\n//\tf()")
		.t("c = class\n\t\t{\n\t\tF()\n\t\t\t{\n\t\t\tx = 1\n" $
				"//\t\t\t.G()\n\t\t\t}\n\t\t}")
		}

	UnitTest_width()
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

	UnitTest_verticalArgs()
		{
		// call arguments written one per line stay one per line
		.t("f(aaa,\n\t\tbbb,\n\t\tccc)")
		.t("f(\n\t\taaa,\n\t\tbbb)") // an authored lead break is kept
		.t("x = [aaa,\n\t\tbbb]")

		// a break written between two arguments is kept, the other gaps fill
		.t("Rect(x1, x2,\n\t\ty1, y2)")
		.t("f(aaa, bbb,\n\t\tccc)")
		.t("f(aaa,\n\t\tbbb, ccc, ddd)")
		.t("f(a, b,\n\t\tc: 1, d: 2)")
		.t("f(aaa, bbb,\n\t\t// why\n\t\tccc)")
		.t("Bar(aaa, bbb,\n\t\tccc)\n\t\t{\n\t\tddd\n\t\t}") // block argument
		.t("f(aaaaaaaa,\n\t\tbbbbbbbb, cccccccc)", width: 40)
		.t("f(aaaaaaaa, bbbbbbbb,\n\t\tcccccccc, dddddddd)", width: 40)

		// argument lists the author did not break still fill to the width
		.t("f(aaaaaaaa, bbbbbbbb, cccccccc, dddddddd)",
			"f(aaaaaaaa, bbbbbbbb, cccccccc,\n\t\tdddddddd)", width: 40)
		}

	UnitTest_soleArgLead()
		{
		// a sole argument breaks after ( only when that is what makes it fit
		.t("Register(a_single_unbreakable_argument)",
			"Register(\n\t\ta_single_unbreakable_argument)", width: 40)

		// it has to break either way, so it stays on the callee's line and the
		// lead break would only cost a line
		.t("Register(aaaaaaaa + bbbbbbbb + cccccccccc)",
			"Register(aaaaaaaa + bbbbbbbb +\n\t\tcccccccccc)", width: 40)

		// short enough to need neither
		.t("Register(aaaa + bbbb)", width: 40)

		// more than one argument is unaffected: they pack from the open paren
		.t("Register(aaaaaaaa, bbbbbbbb, cccccccc)",
			"Register(aaaaaaaa, bbbbbbbb,\n\t\tcccccccc)", width: 40)
		}

	UnitTest_rangeAndSubscriptEndPos()
		{
		// range and subscript nodes have no pos/end, so a chain ending in one
		// took its span end from the last child, dropping the closing bracket
		.t("x = 'l1\nl2' $ a[.. -2]")
		.t("x = 'l1\nl2' $ a[1 ..]")
		.t("x = 'l1\nl2' $ a[0 :: 2]")
		.t("x = 'l1\nl2' $ a[3]")
		.t("f('l1\nl2' $ a[1 ..], :b)")
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
