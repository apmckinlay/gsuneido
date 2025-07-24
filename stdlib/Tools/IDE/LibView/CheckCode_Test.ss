// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(CheckCode('function () { 123 }'))
		Assert(CheckCode('function () { $ }') is: false)
		Assert(CheckCode('function () { Curry() }') is: "")
		}
	Test_warnings()
		{
		for c in .funcs
			Assert(.warnings('function(){' $ c[0] $ '}') is: c[1], msg: c[0])
		for c in .classes
			Assert(.warnings('class {    ' $ c[0] $ '}') is: c[1], msg: c[0])
		for c in .full
			Assert(.warnings('           ' $ c[0]) is: c[1], msg: c[0])
		}
	warnings(code)
		{
		results = Object()
		CheckCode(code, lib: 'stdlib', :results)
		return results.Map!({ Object(it.pos - 11, it.len, it.msg) })
		}
	funcs: (
		('CheckCode', #())
		('NonExistent()',
			#((0, 11, "ERROR: can't find: NonExistent")))
		('hello',
			#((0, 5, "ERROR: used but not initialized: hello")))
		('if (true) return function () { }', #())
		('[y: 3]', #())
		('#((Horz) Fill)', #())
		('Alert(x: 123)', #())
		('true ? x : 456',
			#((7, 1, "ERROR: used but not initialized: x")))
		('x = 1; x', #())
		('for x in #() \n x()', #())
		('try Date()\ncatch(x) return x', #())
		('y = y = 1 + y; y',
			#((12, 1, "ERROR: used but not initialized: y")))
		)
	classes: (
		('New() { }', #())
		('F(x){} G(x){}', #(
			(2, 1, "WARNING: initialized but not used: x"),
			(9, 1, "WARNING: initialized but not used: x")))
		('F(x){} G(x){x}',
			#((2, 1, "WARNING: initialized but not used: x")))
		('F(x){FakeFunc(x)} G(x){x}',
			#((5, 8, "ERROR: can't find: FakeFunc")))
		('F(x){x} G(x){}',
			#((10, 1, "WARNING: initialized but not used: x")))
		)
	full: (
		('NonExistent { }',
			#((0, 11, "ERROR: can't find: NonExistent")))
		('function (x, y) { y }',
			#((10, 1, "WARNING: initialized but not used: x")))
		('function (x) { x }', #())
		('function (unused) { }', #())
		('function (fred/*unused*/) { }', #())
		)

	Test_ignore_web()
		{
		Assert(CheckCode('}') isnt: true)
		Assert(CheckCode('}', 'foo.js'))
		Assert(CheckCode('}', 'foo.JS'))
		Assert(CheckCode('}', 'foo.css'))
		Assert(CheckCode('}', 'foo.CSS'))
		}

	Test_prev_def()
		{
		pd = CheckCode.CheckCode_prev_def
		Assert(pd('Control', 'testlib', #(stdlib, testlib)))
		Assert(pd('Control', 'testlib', #(testlib, stdlib)) is: false)
		Assert(pd('NonExistent!', 'testlib', #(stdlib, testlib)) is: false)
		}
	Test_extra_checks()
		{
		test = function (code, expected, name = false)
			{
			cc = new CheckCode(false, name, 'function(){\n' $ code $ '\n}',
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
				{
				Assert(results[0].msg has: expected)
				}
			}
		test('', true)
		test('Date().Begin()', "use Date.Begin/End")
		test('Date().End()', "use Date.Begin/End")
		test('.Trim() isnt ""', "use .Blank?")
		test('false is x.Find', "use .Has?")
		test('false is x.FindIf', "use .Any?")
		test('func(foo :bar)', "must be preceded")
		test('false is i = x.Find', true)
		test('false is i = x.FindIf', true)
		test('//Date().End()', true)
		test('/*Date().End()*/', true)
		test('"Date().End()"', true)
		test('catch (unused)', 'omit (unused)')
		test('LocalCmds()', 'instance not needed')
		test('check0 = str + 0', 'use Number()')
		test('check0 = 0 + str', 'use Number()')
		test('check0 = 0+str', 'use Number()')
		test('check0 = str+0', 'use Number()')
		test('check0 = 0\r\n\t\t+page', 'use Number()')
		test('check0 = 0\r\n\t\t++page', true)
		test('check0 = (check0 + c.Asc()) % 255', true)
		test('Thread.Sleep()', true)
		test('AppSleep()', true)
		test('\tThread.Sleep()', true)
		test('Sleep()', 'use Thread.Sleep')
		test('\tSleep()', 'use Thread.Sleep')
		test(' Sleep()', 'use Thread.Sleep')
		test('{ F(it) }', 'unnecessary block')
		test('{F(it)}', 'unnecessary block')
		test('{ .F(it) }', 'unnecessary block')
		test('{ X.F(it) }', 'unnecessary block')
		test('{ x.F(it) }', true)

		test(".Merge(Object(", "use assignments")
		test(".Merge(Record(", "use assignments")
		test(".Merge(#(", "use assignments")
		test(".Merge(#{", "use assignments")
		test(".Merge([", "use assignments")
		test("[.key].Merge(args)", true)

		test("\t.foo = 123", true)
		test("\t.foo = 123", true, name: "Rule_themall_Test")
		test("\t.foo = 123", "rules should not have side effects", name: "Rule_themall")

		test("5.Times(", " use for ..")
		test("mock.Verify.Times(", true)
		}
	Test_extra_warning()
		{
		test = function (code, warning?)
			{
			cc = new CheckCode(false, false, 'function(){\n' $ code $ '\n}',
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
		test('Print("testing")', 		true)
		test('Inspect("testing")', 		true)
		test('StackTrace()', 			true)
		test('// Print("testing")', 	false)
		test('// Inspect("testing")', 	false)
		test('// StackTrace()', 		false)
		test('ob = Object(name: name)', true)
		test('func(name: name)',  		true)
		test('func(name: namer)',  		false)
		test('func(name: other)',  		false)
		test('func(name: name + 1)',  	false)
		test('func(name: name $ "x")',  false)
		test('func(name: name - 1)',  	false)
		test('func(name: name * 1)',  	false)
		test('func(name: name / 1)',  	false)
		test('func(name: name = 1)',  	false)
		test('func(name: name+1)',  	false)
		test('func(name: name-1)',  	false)
		test('func(name: name*1)',  	false)
		test('func(name: name/1)',  	false)
		test('func(name: name=1)',  	false)
		test('func(name: name$"x")',  	false)
		test('func(name: name or true)',	false)
		test('func(name: name and true)',	false)
		test('func(name: name | true)',	false)
		test('func(name: name & true)',	false)
		test('func(name: name.first)',	false)
		test('func(name: name.first)',	false)
		test('func(name: name[first])', false)
		test('func(name: name())',   	false)
		test('func(update: update?)',   false)
		}
	}
