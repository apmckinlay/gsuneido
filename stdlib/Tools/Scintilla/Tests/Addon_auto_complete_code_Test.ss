// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_getTextNames()
		{
		f = Addon_auto_complete_code.Addon_auto_complete_code_getTextNames
		Assert(f("") is: #())
		Assert(f("a ab abc") is: #())
		Assert(f("now is the time for all good good men") is: #(good, time))
		}

	Test_findCaller()
		{
		fn = Addon_auto_complete_code.Addon_auto_complete_code_findCaller
		text = "function ()
	{
	Foo(1,
		Bar('test', 2 + (3 * C.Fred())))
	if (1 < 2)
		Print('test')
	}"
		Assert(fn(text, 5) is: false)
		Assert(fn(text, 10) is: false)
		Assert(fn(text, 25) is: 'Foo')
		Assert(fn(text, 35) is: 'Bar')
		Assert(fn(text, 46) is: false)
		Assert(fn(text, 55) is: 'C.Fred')
		Assert(fn(text, 70) is: false)
		Assert(fn(text, 85) is: 'Print')
		Assert(fn(text, 999) is: false)

		text = "class
	{
	Foo(a, b, c)
		{
		d = a * (b + c)
		Bar(d)
			{
			Print(it)
			}
		e = ddd
		(ddd + 1)
		}
	}"
		Assert(fn(text, 3) is: false)
		Assert(fn(text, 19) is: false)
		Assert(fn(text, 45) is: false)
		Assert(fn(text, 57) is: 'Bar')
		Assert(fn(text, 76) is: 'Print')
//TODO: revisit if this does become an issue
//		Assert(fn(text, 103) is: false)

		text = "#(1, 2, 3)"
		Assert(fn(text, 3) is: false)
		}

	Test_buildParamList()
		{
		.MakeLibraryRecord(
			[name: "Foo", text: "function (aaa, abc, bbc) {}"],
			[name: "Bar", text: "class {
				New(arg1, arg2 = false) {}
				Test(@args) {} }"],
			[name: "Bar2", text: "Bar {
				test2(arg = {}) {}
				test3() {}}"],
			[name: "Bar3", text: "Bar {
				Test5: 'Test String'
				Test6: 123
				Test7: #20231019
				Test8: #(something)
			}"])
		fn = Addon_auto_complete_code.Addon_auto_complete_code_buildParamList
		Assert(fn('Foo') is: #(aaa, abc, bbc))
		Assert(fn('Bar') is: #(arg1, arg2))
		Assert(fn('Bar2.Test') is: #(args))
		Assert(fn('Bar2.Bar2_test2') is: #(arg))
		Assert(fn('Bar2.Bar2_test3') is: #())
		Assert(fn('Bar2.Test4') is: false)
		Assert(fn('Bar3.Test5') is: false)
		Assert(fn('Bar3.Test6') is: false)
		Assert(fn('Bar3.Test7') is: false)
		Assert(fn('Bar3.Test8') is: false)
		}
	}
