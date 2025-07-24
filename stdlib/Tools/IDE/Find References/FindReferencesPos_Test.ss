// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		frp = function (s, name, orig_name = false)
			{
			if orig_name is false
				orig_name = name
			return FindReferencesPos(s, 'RecName', name, orig_name)
			}
		Assert(frp("hello world", "xxx") is: false)
		Assert(frp("hello world", "hell") is: false) // whole word
		Assert(frp("hello world", "lo") is: false) // whole word
		Assert(frp("hello world", "hello") is: 0)
		Assert(frp("hello world", "hello", "Field_hello") is: 0)
		Assert(frp("hello world", "world") is: 6)
		Assert(frp("hello.world", "world", "Rule_world") is: 6)

		// constant string
		Assert(frp('Foo', 'Foo') is: 0)
		Assert(frp('"Foo"', 'Foo') is: 0)
		Assert(frp('"FooBar"', 'Foo') is: false)
		Assert(frp('#(Foo: Foo)', 'Foo') is: 2)
		Assert(frp('#(Foo: FooBar)', 'Foo') is: false)
		Assert(frp('class
			{
			Foo: Foo
			}', 'Foo') is: 16)

		Assert(frp('function () { .Foo(test.Foo + test["Foo"])}', 'Foo') is: false)

		// base class
		Assert(frp('Foo {}', 'Foo') is: 0)
		}

	Test_FindAllPos()
		{
		fn = FindReferencesPos.FindAllPos
		code = '
_Foo/*found*/
	{
	Foo(a = function () { Foo/*found*/ }
		b = #(Foo: Foo/*found*/))
		{
		c = Object(Foo/*found*/, Foo()/*found*/).Each(Foo/*found*/)
		}
	cl: class
		{
		test()
			{
			Foo/*found*/(Foo/*found*/, a: Foo/*found*/)
			}
		}
	}
	// comment Foo'
		Assert(fn(code, 'Foo', 'Foo', 'Foo')
			is: #(3, 44, 68, 107, 121, 142, 198, 211, 228))

		code = `
TestControl/*found*/
	{
	// TestControl
	Control: #(Test/*found*/)
	New()
		{
		.Test = .FindControl('Test'/*found*/)
		.width = TestControl/*found*/.Width
		.title = 'Test TestControl'
		}
	Test()
		{
		.Test(Test1['Test'] + TestControl1.TestControl)
		"TestControl is a class"
		}
	// Test.test
	Member: ('Test.Abc'/*found*/, 'TestControl.TestControl_abc'/*found*/, '_Test.Abc',
		'Test.', 'Test()'/*found*/, 'Test(flag:)'/*found*/, "call Test()"
		"Test() is a function")
	}`
		Assert(fn(code, 'Foo', 'Test', 'TestControl')
			is: #(2, 57, 109, 138, 325, 346, 411, 430))
		}
	}