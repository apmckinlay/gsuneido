// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		super.Setup()
		text = `class
			{
			CallClass() { }
			Function_doStuff() { .function_private(); .Function_public() }
			function_private() { .helper() }
			Function_public() { }
			helper() { }
			}`
		.MakeLibraryRecord([name: "HelloWorld?", :text])
		}

	Test_verify()
		{
		mock = Mock()
		mock.Func()
		mock.Verify.Func()
		Assert({ mock.Verify.Func2() } throws: 'wanted but not invoked')
		}
	Test_stub_result()
		{
		mock = Mock()
		mock.When.Func().Return(123)
		Assert(mock.Func() is: 123)
		}
	Test_stub_throw()
		{
		mock = Mock()
		mock.When.Func().Throw('error')
		Assert({ mock.Func() } throws: 'error')
		}
	Test_verify_argument_matching()
		{
		mock = Mock()
		mock.Func(123, 'xyz', abc: #(456))
		mock.Verify.Func([anyArgs:])
		mock.Verify.Func([any:], [any:], abc: [any:])
		mock.Verify.Func([is: 123], [startsWith: 'x'], abc: [anyObject:])
		}
	Test_stub_argument_matching()
		{
		mock = Mock()
		mock.When.Func([anyArgs:]).Return(123)
		Assert(mock.Func(123) is: 123)

		mock.When.Check([is: 0]).Throw('error')
		mock.Check(1)
		Assert({ mock.Check(0) } throws: 'error')
		}
	Test_never()
		{
		mock = Mock()
		mock.Verify.Never().Func()

		mock.Func()
		Assert({ mock.Verify.Never().Func() }
			throws: "should not have been invoked")
		}
	Test_specific_number_of_calls()
		{
		mock = Mock()
		mock.Verify.Times(0).Func()

		mock.Func()
		mock.Verify.Times(1).Func()
		Assert({ mock.Verify.Times(2).Func() }
			throws: "wanted 2 calls, but got 1")
		}
	Test_atLeast()
		{
		mock = Mock()
		mock.Func()
		mock.Func()
		mock.Verify.AtLeast(2).Func()
		Assert({ mock.Verify.AtLeast(3).Func() }
			throws: "wanted at least 3 calls, but got 2")
		}
	Test_atMost()
		{
		mock = Mock()
		mock.Func()
		mock.Func()
		mock.Verify.AtMost(2).Func()
		Assert({ mock.Verify.AtMost(1).Func() }
			throws: "wanted at most 1 calls, but got 2")
		}
	Test_multiple_return_values()
		{
		mock = Mock()
		mock.When.Func().Return(1, 2, 3)
		Assert(mock.Func() is: 1)
		Assert(mock.Func() is: 2)
		Assert(mock.Func() is: 3)
		}

	Test_CallThrough()
		{
		mock = Mock(Mock_Test_Class)
		.testCallThrough(mock)

		mock = Mock(Mock_Test_Class)
		.testCallThrough(mock)
		}

	testCallThrough(mock)
		{
		mock.When.Func1(2).CallThrough()
		mock.When.func2(2).CallThrough()
		mock.When.func3().CallThrough()
		Assert(mock.Eval(Mock_Test_Class.Func, 2) is: 2*2+2*(2*2+1)+1)
		mock.Verify.Times(1).Func1(2)
		mock.Verify.Times(1).func2(2)
		mock.Verify.Times(1).func3()
		}

	Test_CallThrough2()
		{
		mock = Mock(Mock_Test_Class)
		mock.When.f1().CallThrough()
		mock.When.f3().Return(3)
		Assert(mock.Eval(Mock_Test_Class.F) is: 5)
		Assert(mock.F1CalledThrough)
		Assert(mock.F2CalledThrough)
		}

	Test_callThrough_and_members()
		{
		mock = Mock(Mock_Test_Class)
		mock.When.Foo().CallThrough()
		Assert(mock.Foo() is: 15)

		mock.When.foo().CallThrough()
		Assert(mock.foo() is: 50)
		}

	Test_child_class()
		{
		mock = Mock(Mock_Test_ChildClass)
		mock.When.Bar().CallThrough()
		mock.When.bar().CallThrough()
		mock.When.Foo().CallThrough()
		Assert(mock.Bar() is: 100 + 10 + 4)
		Assert(mock.bar() is: 14 * 4)
		}

	Test_mock_inside_mock()
		{
		mock = Mock(Mock_Test_Class)
		inner = Mock(Mock_Test_ChildClass)
		mock.When.M1([anyArgs:]).CallThrough()

		Assert(mock.M1(inner) is: 'from m2')

		mock.Verify.M1(inner)
		mock.Verify.m1(inner)
		mock.Verify.m2()
		inner.Verify.M2()
		inner.Verify.Never().m3()

		mock = Mock(Mock_Test_Class)
		inner = Mock(Mock_Test_ChildClass)
		mock.When.M1([anyArgs:]).CallThrough()
		inner.When.M2().CallThrough()
		inner.When.m3().Return(false)

		Assert(mock.M1(inner) is: 'from m2')

		mock.Verify.M1(inner)
		mock.Verify.m1(inner)
		mock.Verify.m2()
		inner.Verify.M2()
		inner.Verify.m3()
		}

	Test_mock_prioritize_specific_patterns()
		{
		mock = Mock(Mock_Test_Class)
		mock.When.Func([anyArgs:]).Return('anyArgs')
		mock.When.Func('a').Return('a')
		Assert(mock.Func('hello') is: 'anyArgs')
		Assert(mock.Func('a') is: 'a')
		}

	Test_classEndsIn?()
		{
		cl = Global('HelloWorld?')
		mock = Mock(cl)
		mock.When.Function_doStuff([anyArgs:]).CallThrough()
		mock.When.function_private([anyArgs:]).CallThrough()
		mock.When.Function_public([anyArgs:]).CallThrough()

		mock.function_private()
		mock.Verify.Times(1).function_private()
		mock.Verify.Times(1).helper()

		mock.Function_doStuff()
		mock.Verify.Times(1).Function_doStuff()
		mock.Verify.AtLeast(1).function_private()
		mock.Verify.AtLeast(1).helper()
		mock.Verify.Times(1).Function_public()
		}

	Test_method_referenced_as_variable()
		{
		mock = Mock(Mock_Test_Class)
		mock.When.Foo2([anyArgs:]).Return(1)
		mock.When.foo3([anyArgs:]).Return(1)
		mock.When.Foo4([anyArgs:]).Do({|@unused| 2 })
		mock.When.foo5([anyArgs:]).Do({|@unused| 2 })

		Assert(mock.Eval(Mock_Test_Class.Foo1)
			is: #(1, 1, (1, 1), (1, 1)
				2, 2, (2, 2), (2, 2)))
		}
	}
