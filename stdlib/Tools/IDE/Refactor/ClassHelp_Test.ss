// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Class?()
		{
		Assert(not ClassHelp.Class?(""))
		Assert(not ClassHelp.Class?("function () { }"))
		Assert(not ClassHelp.Class?("fred { }"))
		Assert(ClassHelp.Class?("class { }"))
		Assert(ClassHelp.Class?("class : Fred { }"))
		Assert(ClassHelp.Class?("Fred { }"))
		Assert(not ClassHelp.Class?("Fred xxx"))
		Assert(ClassHelp.Class?("_Fred { }"))
		Assert(ClassHelp.Class?(.text))
		}
	Test_Base()
		{
		Assert(ClassHelp.SuperClass("class { }") is: false)
		Assert(ClassHelp.SuperClass("class : Fred { }") is: "Fred")
		Assert(ClassHelp.SuperClass("Fred { }") is: "Fred")
		}
	text:
"class
	{
	Zero: 123
	One()
		{
		one stuff
		}
	Two: function ()
		{
		two stuff
		}
	three()
		{
		three stuff
		}
	}"
	method:
"Added()
	{
	added stuff
	}"

	Test_AdvanceToNewline()
		{
		Assert(ClassHelp.AdvanceToNewline(.text, .text.Find('}'))
			is: .text.Find('\tTwo'))
		}
	Test_AfterMethod()
		{
		Assert(ClassHelp.AfterMethod(.text, .text.Find('one'))
			is: .text.Find('\tTwo'))
		}
	Test_AddMethod()
		{
		Assert(ClassHelp.AddMethod(.text, .text.Find('one'), .method)
			is: .text.Replace('\tTwo', '\t' $ .method $ '\r\n\tTwo'))

		Assert(ClassHelp.AddMethod(.text, .text.Find('three stuff'), .method)
			is: .text.Replace('^\t}', '\t' $ .method $ '\r\n\t}'))
		}
	Test_AddMethodAtEnd()
		{
		Assert(ClassHelp.AddMethodAtEnd(.text, .method)
			is: .text.Replace('^\t}', '\t' $ .method $ '\r\n\t}'))
		}
	Test_MethodRange()
		{
		Assert(ClassHelp.MethodRange(.text, .text.Find("one"))
			is: Object(from: .text.Find("One"), to: .text.Find("\tTwo")))
		Assert(ClassHelp.MethodRange(.text, .text.Find("two"))
			is: Object(from: .text.Find("Two") + 5, to: .text.Find("\tthree")))
		Assert(ClassHelp.MethodRange(.text, .text.Find("three stuff"))
			is: Object(from: .text.Find("three"), to: .text.Size() - 2))
		}
	Test_MethodName()
		{
		Assert(ClassHelp.MethodName(.text, 0) is: false)
		Assert(ClassHelp.MethodName(.text, .text.Find('one')) is: 'One')
		Assert(ClassHelp.MethodName(.text, .text.Find('three stuff')) is: 'three')
		}
	Test_Locals()
		{
		Assert(ClassHelp.Locals('') is: #())

		text = 'Func(a, b: c) { ++a for one in two ++three }'
		Assert(ClassHelp.Locals(text) is: #(a, c, one, two, three))
		}
	Test_LocalsInputs()
		{
		Assert(ClassHelp.LocalsInputs('') is: #())

		text = "a; ++b; c(d); e--; f = g; h += i"
		Assert(ClassHelp.LocalsInputs(text) is: #(a, b, c, d, e, g, h, i))
		}
	Test_LocalsModified()
		{
		Assert(ClassHelp.LocalsModified('') is: #())

		text = "a; ++b; c(d); e--; f = g; h += i"
		Assert(ClassHelp.LocalsModified(text) is: #(b, e, f, h))
		}
	Test_LocalsAssigned()
		{
		Assert(ClassHelp.LocalsModified('') is: #())

		text = "method(a, b) { c = 0; fn(d); ++e; f = 1; for g in h { |i,j| } }"
		Assert(ClassHelp.LocalsAssigned(text) is: #(a, b, c, f, g, i, j))
		}
	Test_Methods()
		{
		Assert(ClassHelp.Methods('') is: #())

		Assert(ClassHelp.Methods(.text) is: #(One, Two, three))
		}
	Test_FindMethod()
		{
		text =
			'class
				{
				One()
					{ two = 1 }
				two()
					{ }
				}'
		Assert(ClassHelp.FindMethod(text, 'xx') is: false)
		Assert(ClassHelp.FindMethod(text, 'One') is: text.Find('One()'))
		Assert(ClassHelp.FindMethod(text, 'two') is: text.Find('two()'))
		}
	ch: ClassHelp
		{
		ClassHelp_libraries()
			{ return _libs.Copy() }
		ClassHelp_trials()
			{ return #(tag) }
		}
	Test_FindBaseMethod()
		{
		lib1 = .MakeLibrary(
			[name: 'ATestRecord', text: 'class { Foo() { 1 } }'],
			[name: 'ATestRecord__webgui', text: 'class { Foo() { 2 } }'],
			[name: 'ATestRecord__webgui_tag', text: 'class { Foo() { 3 }}'])
		lib2 = .MakeLibrary(
			[name: 'ATestRecord', text: '_ATestRecord {}'])
		lib3 = .MakeLibrary(
			[name: 'CTestRecord', text: 'ATestRecord { Bar() { .Foo() }'],
			[name: 'CTestRecord__webgui', text: '_CTestRecord { Baz() {} }'])

		_libs = [lib1, lib2, lib3]
		fn = .ch.FindBaseMethod

		Assert(fn(lib1, 'BTestRecord', '', 'Foo') is: false)
		Assert(fn(lib2, 'ATestRecord', '_ATestRecord {}', 'Foo')
			is: Object(lib: lib1, name: 'ATestRecord'))
		Assert(fn(lib2, 'ATestRecord', '_ATestRecord {}', 'Missing') is: false)

		Assert(fn(lib3, 'CTestRecord', 'ATestRecord {}', 'Foo')
			is: Object(lib: lib1, name: 'ATestRecord'))
		Assert(fn(lib3, 'CTestRecord', 'ATestRecord {}', 'Missing') is: false)
		Assert(fn(lib3, 'CTestRecord__webgui', '_CTestRecord {}', 'Bar')
			is: Object(lib: lib3, name: 'CTestRecord'))
		Assert(fn(lib3, 'CTestRecord__webgui', '_CTestRecord {}', 'Foo')
			is: Object(lib: lib1, name: 'ATestRecord__webgui_tag'))
		}
	memberstext:
		"class
			{
			foo() { x.far = 456; x = 2 }
			Bar() { .bar = 123; .y }
			foobar: 789
			Z: 0
			}"
	Test_PrivateMembers()
		{
		Assert(ClassHelp.PrivateMembers("") is: #())
		Assert(ClassHelp.PrivateMembers(.memberstext) is: #(bar, foo, foobar))
		}
	Test_PublicMembers()
		{
		Assert(ClassHelp.PublicMembers("") is: #())
		Assert(ClassHelp.PublicMembers(.memberstext) is: #(Bar, Z))
		text = "Point { Foo: 123 bar: 456 }"
		mems = ClassHelp.PublicMembers(text)
		Assert(mems has: "Foo")
		Assert(mems has: "GetX")
		Assert(mems hasnt: "bar")
		Assert(mems hasnt: "x")
		Assert(mems hasnt: "Point_x")
		}
	Test_PublicMembersOfName()
		{
		mems = ClassHelp.PublicMembersOfName('ClassHelp_Test')
		Assert(mems has: "Test_PublicMembersOfName")
		Assert(mems hasnt: "memberstext")
		Assert(ClassHelp_Test.Members() has: "ClassHelp_Test_memberstext")
		Assert(mems hasnt: "ClassHelp_Test_memberstext")
		}
	Test_MethodRanges()
		{
		Assert(ClassHelp.MethodRanges("class { One() { fred } Two() { joe } }")
			is: #([name: "One", to: 22, from: 11], [name: "Two", to: 36, from: 26]))
		}
	Test_MethodSizes()
		{
		text = "// comment
				function ()
					{
					a = 123
					}"
		Assert(ClassHelp.MethodSizes(text) is: #([lines: 4, from: 16]))
		text = "
			class
				{
				One()
					{
					// comment followed by blank line - not counted

					}
				two()
					{
					a = 'a line'
					}
				}"
		sizes = ClassHelp.MethodSizes(text)
		sizes.Each { it.Delete(#from, #to) }
		Assert(sizes is: #([name: One, lines: 3], [name: two, lines: 4]))
		}
	Test_nonWhiteLineCount()
		{
		f = ClassHelp.ClassHelp_nonWhiteLineCount
		Assert(f('') is: 0)
		Assert(f('x') is: 1)
		Assert(f('x
			y') is: 2)
		Assert(f('x

			y') is: 2)
		Assert(f('x
			// comment
			y') is: 2)
		Assert(f('x
			/* comment */
			y') is: 2)
		Assert(f('x
			/*
			comment
			*/
			y') is: 2)
		Assert(f('x /*
			comment
			*/ y') is: 2)
		}
	Test_RetrieveParamsList()
		{
		test = function (text, expected)
			{
			body = '{ x }'
			text $= body
			pos = ClassHelp.RetrieveParamsList(text, actual = Object())
			Assert(actual is: expected)
			Assert(text[pos..] is: body)
			}
		test('()', #())
		test('(a,b,c)', #(a,b,c))
		test('(a,b,c)', #(a,b,c))
		test('(.a, _b, ._c, .Dad, parentHwnd)', #(a, b, c, dad, parentHwnd))
		test('(a = 1, b = (2), c = (1, (2), 3))', #(a,b,c))
		test('(a = function() {return true}, b = 2)', #(a, b))
		}

	Test_AllClassMembers()
		{
		test = "class
	{
	y: false
	New ()
		{
		x: true
		}
	MyMethod1()
		{
		.myVal1 = true
		}
	myMethod2(.myVal2 = 'default')
		{
		_report = false
		.myVal3 = '15'
		}
	}"
		testTwo = `#(Params
	title: "Business Partners List",
	name: 'Biz_Partner_Report'
	printParams: (bizpartner_num)
	Params:
		#(Form
			(ParamsSelect bizpartner_num, group: 0) nl
			(ParamsSelect bizpartner_city, group: 0) nl
			(ParamsSelect bizpartner_state_prov, group: 0) nl
			nl
			(Biz_PartnerRolesControl, group: 0) nl
			(bizpartner_roles_condition, group: 0) nl
			(bizpartner_status_params, group: 0) nl
			(bizpartner_print_contacts?, group: 0) nl
		))`
		calculatedMembersList = ClassHelp.AllClassMembers(test)
		Assert(calculatedMembersList.Members()
			equalsSet: #("y:", "New", "MyMethod1", "myMethod2",
				"myVal1:", "myVal2:", "myVal3:"))
		Assert(ClassHelp.AllClassMembers(testTwo).Members()
			equalsSet: #("Params:", "name:", "printParams:", "title:"))
		}



	Test_classMemberDotDeclarations()
		{
		test = "class
	{
	shouldNotFind1: false
	New ()
		{
		shouldNotFind2: true
		}
	MyMethod1(z = ClassHelp.Qc_Main(lib, name), ob = #{}, ob2 = #())
		{
		.myVal1 = true
		}
	myMethod2(.myVal2 = 'default')
		{
		_report = false
		.myVal3 = '15'
		}
	}"
		classVariables = Object()
		calculatedResult = ClassHelp.ClassHelp_classMemberDotDeclarations(
			test, classVariables, 'M')
		Assert(calculatedResult is: #('myVal1:', 'myVal2:', 'myVal3:'))

		Assert(ClassHelp.ClassHelp_classVariablesInMethodBody(
			test, classVariables, 'M', find: 'myVal1')
			is: 159)
		}
	Test_classVariablesInMethodBody()
		{
		test = "(.shouldNotFlag, .a, .b, .d)
			{
			.x = 55
			y = .a = 77
			for (.i = 0; .i < 5; .i++)
				{
				.legit = 69
				}
			.c = { |x| x * x}; .v = 'hmm'
			Object(1,2,3,4,5,6).Map({ it += .d })
			.n = true isnt false ? 5 : 4
			test = false
			.s = '15'

			data = Object()
			data.test = ''
			data.
			testTwo = ''
			fn(data).testThree = 5
			.call(.call1 = 5, .call2 = 55)
			data
			.myVal1 = 55
			fn(.myVal6 = 35)
			.
			myVal7 = 7
			fn() { .myVal19 = 90 }
			x = .myVal11 = 42
			.yy = .zz = 41
			.myVal14 = myVal15 = 66
			}"
		classVariables = Object()
		ClassHelp.ClassHelp_classVariablesInMethodBody(test, classVariables, 'M')
		Assert(classVariables is:  #("x:", "a:", "i:", "legit:", "c:", "v:", "n:", "s:",
			"call1:", "call2:", "myVal1:", "myVal6:", "myVal7:", "myVal19:", "myVal11:",
			"yy:", "zz:", "myVal14:"))

		test1 = "()
			{
			shouldNotFind2: true
			}"
		test2 = "()
			{
			.myVal1 = true
			}"
		test3 = "(.myVal2 = 'default')
			{
			_report = false
			.myVal3 = '15'
			}"
		test4 = "()
			{
			if 0 is .rc = .Send(#GetRecordControl)
			throw 'PresetsControl must be used inside a RecordControl'
			}"
		classVariables = Object()
		ClassHelp.ClassHelp_classVariablesInMethodBody(test1, classVariables, 'M')
		ClassHelp.ClassHelp_classVariablesInMethodBody(test2, classVariables, 'M')
		ClassHelp.ClassHelp_classVariablesInMethodBody(test3, classVariables, 'M')
		ClassHelp.ClassHelp_classVariablesInMethodBody(test4, classVariables, 'M')
		Assert(classVariables is: #("myVal1:", "myVal3:", "rc:"))
		}

	Test_classMemberDeclaresInParamsMethod()
		{
		test = "(.x, .myVal4 = 'default', .y, ob = #{}, ob2 = #(), #IAMANIDENTIFIER)
			{
			.shouldNotFlag = 55
			throw 'PresetsControl must be used inside a RecordControl'
			}"
		classVariables = Object()
		ClassHelp.ClassHelp_classMemberDeclaresInParamsMethod(test, classVariables, 'M')
		Assert(classVariables is: #('x:', 'myVal4:', 'y:'))

		Assert(ClassHelp.ClassHelp_classMemberDeclaresInParamsMethod(
			test, classVariables, 'M', find: 'x')
			is: 3)
		Assert(ClassHelp.ClassHelp_classMemberDeclaresInParamsMethod(
			test, classVariables, 'M', find: 'y')
			is: 28)
		}

	}




