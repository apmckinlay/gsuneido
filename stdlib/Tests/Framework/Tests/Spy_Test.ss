// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	src1: "function (a, b = 1)
		{
		return a + b
		}"
	src2: `class
		{
		CallClass(a, _b = 1)
			{
			c = new .c
			return "c.T: " $ c.T(a, b) $
				" - B: " $ .B(a) $
				" - SpyTest_1: " $ SpyTest_1(a, b)
			}
		B(a)
			{
			return a * 2
			}
		c: class
			{
			T(a, _b)
				{
				return a + b
				}
			}
		}`
	src3: `class
		{
		CallClass() { }
		private_func() { }
		Public_func() { }
		}`
	Setup()
		{
		.MakeLibraryRecord([name: "SpyTest_1", text: .src1])
		.MakeLibraryRecord([name: "SpyTest_2", text: .src2])
		.MakeLibraryRecord([name: "SpyTest_3?", text: .src3])
		}

	Test_SpyOnFunction()
		{
		fn = Global(#SpyTest_1)
		Assert(fn(1) is: 2)

		spy = .SpyOn("SpyTest_1")
		spy.Return("case 1", when: { |a, b| a + b > 5 })
		spy.Throw("case 2", when: function (a, b) { a + b is 5 })
		spy.Return("case 3 - 1", "case 3 - 2", when: { |a, b| a + b > 0 })

		fn = Global(#SpyTest_1)
		Assert(fn(5) is: 'case 1')
		Assert({ fn(4) } throws: 'case 2')
		Assert(fn(3) is: 'case 3 - 1')
		Assert(fn(2) is: 'case 3 - 2')
		Assert({ fn(1) } throws:)
		Assert(fn(-1, -2) is: -3)

		callLogs = spy.CallLogs()
		Assert(callLogs isSize: 6)
		Assert(callLogs is: #((a: 5, b: 1), (a: 4, b: 1), (a: 3, b: 1), (a: 2, b: 1),
			(a: 1, b: 1), (a: -1, b: -2)))

		spy.Close()
		fn = Global(#SpyTest_1)
		Assert(fn(5) is: 6)
		}

	Test_SpyOnClassMethod()
		{
		_b = 2
		spy1 = .SpyOn(Global('SpyTest_2.SpyTest_2_c.T'))
		spy1.Throw("a is 1", when: { |a| a is 1 })
		spy2 = .SpyOn(Global('SpyTest_2.B'))
		spy2.Return("override")
		spy3 = .SpyOn(Global('SpyTest_1'))
		spy3.Return("a + b > 5", when: { |a, b| a + b > 5 })
		spy4 = .SpyOn(Global('SpyTest_2'))

		fn = Global('SpyTest_2')
		Assert({ fn(1) } throws: "a is 1")
		Assert(fn(2) is: "c.T: 4 - B: override - SpyTest_1: 4")
		Assert(fn(5) is: "c.T: 7 - B: override - SpyTest_1: a + b > 5")

		spy2.Close()
		fn = Global('SpyTest_2')
		Assert({ fn(1) } throws: "a is 1")
		Assert(fn(2) is: "c.T: 4 - B: 4 - SpyTest_1: 4")
		Assert(fn(5) is: "c.T: 7 - B: 10 - SpyTest_1: a + b > 5")

		Assert(spy4.CallLogs() is: #((a: 1, b: 2), (a: 2, b: 2), (a: 5, b: 2),
			(a: 1, b: 2), (a: 2, b: 2), (a: 5, b: 2)))
		}

	Test_setupInfo()
		{
		fakeSpy = Spy
			{
			Spy_registerSpy() {}
			}
		spy1 = fakeSpy(Global("SpyTest_2"))
		Assert(Display(spy1.Target) is: "SpyTest_2.CallClass /* Test_lib method */")
		Assert(spy1.Name is: "SpyTest_2")
		Assert(spy1.Paths is: #("CallClass"))
		Assert(spy1.Lib is: "Test_lib")
		Assert(spy1.Method?)
		Assert(spy1.Params is: '(a,_b=1)')

		spy2 = fakeSpy(Global("SpyTest_3?"))
		target = "SpyTest_3?.CallClass /* Test_lib method */"
		Assert(Display(spy2.Target) is: target)
		Assert(spy2.Name is: "SpyTest_3?")
		Assert(spy2.Paths is: #("CallClass"))
		Assert(spy2.Lib is: "Test_lib")
		Assert(spy2.Method?)
		Assert(spy2.Params is: '()')
		}

	Test_CallClass()
		{
		.SpyOn(Xml).Return('<body></body>')
		Assert(Xml(#(abc: 'test')) is: '<body></body>')

		// QueryCost.CallClass and GetContributions.CallClass are defined in Memoize
		// testing spy to not override Memoize.CallClass
		.SpyOn(GetContributions).Return(#('return from GetContributions'))
		result = QueryCost('stdlib')
		Assert(result hasMember: #nrecs)
		}
	}