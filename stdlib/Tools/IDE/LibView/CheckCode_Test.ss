// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(CheckCode("function () { 123 }"))
		Assert(CheckCode("function () { $ }") is: false)
		}

	Test_warnings()
		{
		for c in .funcs
			Assert(.warnings("function(){" $ c[0] $ '}') is: c[1], msg: c[0])
		for c in .getClasses()
			Assert(.warnings("class {    " $ c[0] $ '}') is: c[1], msg: c[0])
		for c in .full
			Assert(.warnings("           " $ c[0]) is: c[1], msg: c[0])
		}

	warnings(code)
		{
		results = Object()
		CheckCode(code, lib: #stdlib, :results)
		return results.Map!({ [it.pos - 11, it.len, it.msg] })
		}

	funcs: (
		(CheckCode, ()),
		("NonExistent()",
			((0, 11, "ERROR: can't find: NonExistent"))),
		(hello,
			((0, 5, "ERROR: used but not initialized: hello"))),
		("if (true) return function () { } else return false", ()),
		("[y: 3]", ()),
		("#((Horz) Fill)", ()),
		("Alert(x: 123)", ()),
		("true ? x : 456",
			((7, 1, "ERROR: used but not initialized: x"))),
		("x = 1; x", ()),
		("for x in #() \n x()", ()),
		("try Date()\ncatch(x) return x \n return false", ()),
		("y = y = 1 + y; y",
			((12, 1, "ERROR: used but not initialized: y")))
		)
	getClasses()
		{
		classes = [
			#("New() { }", ()),
			#("F(x){} G(x){}",
				(
					(2, 1, "WARNING: initialized but not used: x"),
					(9, 1, "WARNING: initialized but not used: x")
					)),
			#("F(x){} G(x){x}",
				((2, 1, "WARNING: initialized but not used: x"))),
			#("F(x){FakeFunc(x)} G(x){x}",
				((5, 8, "ERROR: can't find: FakeFunc"))),
			#("F(x){x} G(x){}",
				((10, 1, "WARNING: initialized but not used: x")))]

		classes.Append(#(
			("Foo(x :string) :string { x }", ()),
			("Foo(x: string|number) :object|false { x }", ()),
			("Foo(x :string = 0, y: object = false) :string|false {y[x] }", ())
			))
		return classes
		}

	full: (
		("NonExistent { }",
			((0, 11, "ERROR: can't find: NonExistent"))),
		("function (x, y) { y }",
			((10, 1, "WARNING: initialized but not used: x"))),
		("function (x) { x }", ()),
		("function (unused) { }", ()),
		("function (fred/*unused*/) { }", ())
		)

	Test_ignore_web()
		{
		Assert(CheckCode('}') isnt:)
		Assert(CheckCode('}', "foo.js"))
		Assert(CheckCode('}', "foo.JS"))
		Assert(CheckCode('}', "foo.css"))
		Assert(CheckCode('}', "foo.CSS"))
		}

	Test_prev_def()
		{
		pd = CheckCode.CheckCode_prev_def
		Assert(pd(#Control, #testlib, #(stdlib, testlib)))
		Assert(pd(#Control, #testlib, #(testlib, stdlib)) is: false)
		Assert(pd(#NonExistent!, #testlib, #(stdlib, testlib)) is: false)
		}

	Test_extra_checks()
		{
		test = function(code, expected, name = false)
			{
			cc = new CheckCode(false, name, "function(){\n" $ code $ "\n}",
				results = [])
			cc.CheckCode_extra_checks()
			errors? = cc.CheckCode_errors?
			warnings? = cc.CheckCode_warnings?
			if expected is true
				{
				Assert(errors? is: false)
				Assert(warnings? is: false)
				Assert(results is: [])
				}
			else
				Assert(results[0].msg has: expected)
			}
		test("", true)
		test("Date().Begin()", "use Date.Begin/End")
		test("Date().End()", "use Date.Begin/End")
		test('.Trim() isnt ""', "use .Blank?")
		test("false is x.Find", "use .Has?")
		test("false is x.FindIf", "use .Any?")
		test("false is i = x.Find", true)
		test("false is i = x.FindIf", true)
		test("//Date().End()", true)
		test("/*Date().End()*/", true)
		test('"Date().End()"', true)
		test("catch (unused)", "omit (unused)")
		test("LocalCmds()", "instance not needed")
		test("check0 = str + 0", "use Number()")
		test("check0 = 0 + str", "use Number()")
		test("check0 = 0+str", "use Number()")
		test("check0 = str+0", "use Number()")
		test("check0 = 0\r\n\t\t+page", "use Number()")
		test("check0 = 0\r\n\t\t++page", true)
		test("check0 = (check0 + c.Asc()) % 255", true)
		test("Thread.Sleep()", true)
		test("AppSleep()", true)
		test("\tThread.Sleep()", true)
		test("Sleep()", "use Thread.Sleep")
		test("\tSleep()", "use Thread.Sleep")
		test(" Sleep()", "use Thread.Sleep")
		test("{ F(it) }", "unnecessary block")
		test("{F(it)}", "unnecessary block")
		test("{ .F(it) }", "unnecessary block")
		test("{ X.F(it) }", "unnecessary block")
		test("{ x.F(it) }", true)

		test(".Merge(Object(", "use assignments")
		test(".Merge(Record(", "use assignments")
		test(".Merge(#(", "use assignments")
		test(".Merge(#{", "use assignments")
		test(".Merge([", "use assignments")
		test("[.key].Merge(args)", true)

		test("\t.foo = 123", true)
		test("\t.foo = 123", true, name: #Rule_themall_Test)
		test("\t.foo = 123", "rules should not have side effects", name: #Rule_themall)

		test("5.Times(", " use for ..")
		test("mock.Verify.Times(", true)

		test("func(foo :bar)", true) // this used to be invalid, now is an annotation
		test("func(x :string)", true)
		test("func(x: string)", true)
		test("func(x: string|number)", true)
		test("func(x :string, y :number|false)", true)
		test("func() :bar", true)
		test("func() :string|false", true)
		test("1 + :bar", "must be preceded")
		}

	Test_extra_warning()
		{
		test = function(code, warning?)
			{
			cc = new CheckCode(false, false, "function(){\n" $ code $ "\n}",
				results = [])
			cc.CheckCode_extra_checks()
			errors? = cc.CheckCode_errors?
			warnings? = cc.CheckCode_warnings?
			if warning?
				{
				Assert(warnings?)
				Assert(results isSize: 1)
				}
			else
				{
				Assert(errors? is: false)
				Assert(warnings? is: false)
				Assert(results is: [])
				}
			}
		test('Print("testing")', true)
		test('Inspect("testing")', true)
		test("StackTrace()", true)
		test('// Print("testing")', false)
		test('// Inspect("testing")', false)
		test("// StackTrace()", false)
		}

	Test_comby_matches()
		{
		code = "function (foo) {\r\n\tfoo(1)\r\n\tfoo(\r\n\t\t2)\r\n}"
		cc = .combyChecker(code, #(("foo(:[x])", "avoid foo")), results = [])
		Assert(cc() is: false)
		Assert(results isSize: 2)
		Assert(results[0]
			is: Object(pos: 19, len: 6, line: 1,
				msg: "ERROR: use of foo(1) - avoid foo"))
		Assert(results[1]
			is: Object(pos: 28, len: 10, line: 2,
				msg: "ERROR: use of foo( 2) - avoid foo"))
		Assert(cc.CheckCode_errors?)
		Assert(cc.CheckCode_warnings? is: false)
		}

	Test_comby_warning()
		{
		cc = .combyChecker("function (foo) { foo(1) }",
			#(("foo(:[x])", "", warning:)), results = [])
		Assert(cc() is: "")
		Assert(results isSize: 1)
		Assert(results[0].msg is: "WARNING: use of foo(1)")
		Assert(cc.CheckCode_errors? is: false)
		Assert(cc.CheckCode_warnings?)
		}

	Test_comby_ignored()
		{
		for code in #('"foo(1)"', "// foo(1)", "/* foo(1) */", "bar(1)")
			{
			cc = .combyChecker(code, #(("foo(:[x])", "avoid foo")), results = [])
			cc.CheckCode_extra_checks()
			Assert(results is: [])
			Assert(cc.CheckCode_errors? is: false)
			Assert(cc.CheckCode_warnings? is: false)
			}
		}

	Test_comby_filter()
		{
		code = "foo(1)\r\nfoo(2)"
		rules = #(("foo(:[x])", "avoid foo",
				filter: function(code, pos, len) { return code[pos::len] is "foo(1)" }))
		cc = .combyChecker(code, rules, results = [])
		cc.CheckCode_comby()
		Assert(results isSize: 1)
		Assert(results[0]
			is: Object(pos: 8, len: 6, line: 1,
				msg: "ERROR: use of foo(2) - avoid foo"))
		}

	Test_comby_noHint()
		{
		cc = .combyChecker("true", #(("true", "avoid true")), results = [])
		cc.CheckCode_comby()
		Assert(results isSize: 1)
		Assert(results[0].msg is: "ERROR: use of true - avoid true")
		}

	combyChecker(code, rules, results)
		{
		cc = new CheckCode(false, false, code, results)
		cc.CheckCode_badComby = rules
		return cc
		}
	}
