// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	ed: class
		{
		New(.text, from, to)
			{ .SetSelect(from, to - from) }
		Get()
			{
			return .text
			}
		Addon_go_to_definition_isSuper(ref)
			{
			return Addon_go_to_definition.Addon_go_to_definition_isSuper(ref)
			}
		Addon_go_to_definition_findBeforeDot(org, wc)
			{
			return this.Eval(
				Addon_go_to_definition.Addon_go_to_definition_findBeforeDot, org, wc)
			}
		Addon_go_to_definition_findUntilNotIn(pos, wc, dir = 1)
			{
			return this.Eval(
				Addon_go_to_definition.Addon_go_to_definition_findUntilNotIn,
					pos, wc, dir)
			}
		Addon_go_to_definition_skipFirstDot?(org, orgToEndStr, isJsOrCss)
			{
			return this.Eval(
				Addon_go_to_definition.Addon_go_to_definition_skipFirstDot?,
					org, orgToEndStr, isJsOrCss)
			}
		GetAt(i)
			{ return .text[i] }
		GetSelect()
			{ return .sel }
		SetSelect(i, n)
			{ .sel = Object(cpMin: i, cpMax: i + n) }
		GetSelText()
			{ return .text[.sel.cpMin .. .sel.cpMax] }
		GetWordChars()
			{ return Addon_suneido_style.WordChars }
		}
	Test_selectRef()
		{
		f = Addon_go_to_definition.Addon_go_to_definition_selectRef

		c = new .ed("one two three", 4, 7)
		Assert(c.Eval(f) is: 'two')

		c = new .ed("one two three", 5, 5)
		Assert(c.Eval(f) is: 'two')

		c = new .ed("one two three", 4, 4)
		Assert(c.Eval(f) is: 'two')

		c = new .ed("one two three", 7, 7)
		Assert(c.Eval(f) is: 'two')

		c = new .ed("one .two three", 6, 6)
		Assert(c.Eval(f) is: '.two')

		c = new .ed("one).two three", 6, 6)
		Assert(c.Eval(f) is: ').two')

		c = new .ed("one two.three four", 10, 10)
		Assert(c.Eval(f) is: 'three')

		c = new .ed("one .two.three four", 11, 11)
		Assert(c.Eval(f) is: 'three')

		c = new .ed("one .one.two.three four", 14, 14)
		Assert(c.Eval(f) is: 'three')

		c = new .ed("one Two.three four", 11, 11)
		Assert(c.Eval(f) is: 'Two.three')

		c = new .ed("one super.Three four", 11, 11)
		Assert(c.Eval(f) is: 'super.Three')

		c = new .ed("one two.js four", 4, 4)
		Assert(c.Eval(f) is: 'two.js')

		c = new .ed("one two.js four", 8, 8)
		Assert(c.Eval(f) is: 'two.js')

		c = new .ed("one Two.css four", 8, 8)
		Assert(c.Eval(f) is: 'Two.css')

		c = new .ed("one Two.cs four", 4, 4)
		Assert(c.Eval(f) is: 'Two')

		c = new .ed("one css.two four", 8, 8)
		Assert(c.Eval(f) is: 'two')

		c = new .ed('testing colon fakelib:fake_class test', 17, 17)
		Assert(c.Eval(f) is: 'fakelib:fake_class')

		c = new .ed('testing fakelib:fake_class:15 with line number', 17, 17)
		Assert(c.Eval(f) is: 'fakelib:fake_class:15')

		c = new .ed('testing fakelib:fake_class:15 with line number', 16, 26)
		Assert(c.Eval(f) is: 'fake_class')

		c = new .ed('testing Object(:field_definition) :name shortcut', 17, 17)
		Assert(c.Eval(f) is: 'field_definition')

		c = new .ed('testing Object(field_definition:) set true shortcut', 17, 17)
		Assert(c.Eval(f) is: 'field_definition')

		c = new .ed('testing :fake_class: with extra colons', 17, 17)
		Assert(c.Eval(f) is: 'fake_class')

		c = new .ed('testing :super with extra colon', 12, 12)
		Assert(c.Eval(f) is: 'super')
		}
	Test_validRef()
		{
		f = Addon_go_to_definition.Addon_go_to_definition_validRef
		Assert(f('target'))
		Assert(f('Target'))
		Assert(f('.target'))
		Assert(f('.Target'))
		Assert(f('Foo.Target'))
		Assert(f('super.Target'))
		Assert(f('target.js'))
		Assert(f('lib:target.js:9'))
		Assert(not f('target.jss'))
		Assert(not f('Foo.target'))
		Assert(not f('foo.Target'))
		Assert(not f('.Foo.Target'))
		}

	Test_gotoObject?()
		{
		sciMock = Mock()
		sciMock.When.Send([anyArgs:]).Return('')
		cl = Addon_go_to_definition
			{
			Defer(unused) { }
			SetFocus() { }
			Addon_go_to_definition_gotoLibView(@unused) { }
			}
		newcl = cl(sciMock, [])
		method = newcl.Addon_go_to_definition_gotoObject?
		Assert(method('testObject', `#(fake)`, false) is: false)
		Assert(method('.testObject', `#(fake)`, false), msg: 'gotoObject true')
		Assert(method('.testObject', `not an ob`, false) is: false)
		}

	Test_gotoMethod?()
		{
		if BuiltDate() < #20250218
			return
		.MakeLibraryRecord([name: 'Addon_go_to_definition_FakeTest',
			text: `class
				{
				CallClass()
					{
					Print('CallClass of parent')
					}
				ParentMethod()
					{
					Print('Method of parent')
					}
				ParentMember: #(fake)
				}`])
		.MakeLibraryRecord([name: 'Addon_go_to_definition_FakeTest_Child',
			text: testClass = `Addon_go_to_definition_FakeTest
				{
				fake: true
				New(.Fakest)
					{
					Print('This is New, line 5')
					}

				Method1()
					{
					Print('This is Method1, line 9')
					}

				method2()
					{
					Print('This is method2, line 9')
					}
				}`])

		createMethod = .SpyOn(
			Addon_go_to_definition.Addon_go_to_definition_create_method).Return('')
		gotoMethod = .SpyOn(Addon_go_to_definition.Addon_go_to_definition_gotoMethodLine)

		addon = Addon_go_to_definition(.fakeEditor(testClass, 131), [])
		Assert(addon.Addon_go_to_definition_gotoMethod?('.Method1', testClass),
			msg: 'Method1')
		Assert(gotoMethod.CallLogs()[0].method is: 'Method1')

		Assert(addon.Addon_go_to_definition_gotoMethod?('.Method2', testClass) is: false)
		Assert(createMethod.CallLogs()[0].method is: 'Method2')
		Assert(gotoMethod.CallLogs()[1].method is: 'Method2')

		addon = Addon_go_to_definition(.fakeEditor(testClass, 203), [])
		Assert(addon.Addon_go_to_definition_gotoMethod?('.method2', testClass),
			msg: 'method2')
		Assert(gotoMethod.CallLogs()[2].method is: 'method2')

		Assert(addon.Addon_go_to_definition_gotoMethod?('.method3', testClass) is: false)
		Assert(createMethod.CallLogs()[1].method is: 'method3')
		Assert(gotoMethod.CallLogs()[3].method is: 'method3')

		addon = Addon_go_to_definition(.fakeEditor(testClass, 44), [])
		Assert(addon.Addon_go_to_definition_gotoMethod?('.fake', testClass), msg: '.fake')

		addon = Addon_go_to_definition(.fakeEditor(testClass, 71), [])
		Assert(addon.Addon_go_to_definition_gotoMethod?('.Fakest', testClass),
			msg: '.Fakest')

		// Members / Methods defined in parent class
		addon = Addon_go_to_definition(.fakeEditor(testClass, 23), [])
		Assert(addon.Addon_go_to_definition_gotoMethod?('.CallClass', testClass)
			is: false)

		addon = Addon_go_to_definition(.fakeEditor(testClass, 94, send?:), [])
		Assert(addon.Addon_go_to_definition_gotoMethod?('.ParentMethod', testClass),
			msg: '.ParentMethod')

		addon = Addon_go_to_definition(.fakeEditor(testClass, 94, send?:), [])
		Assert(addon.Addon_go_to_definition_gotoMethod?('.ParentMember', testClass)
			msg: '.ParentMethod 2')
		}

	fakeEditor(text, expectedPos, send? = false)
		{
		return FakeObject(
			LineFromPosition: { |pos| Assert(pos is: expectedPos); 0 }
			Get: 			{ text }
			GetSelect: 		{ [cpMin: 0] }
			GotoLine: 		{ |unused| }
			Defer: 			{ |unused| }
			SetFocus:		{ }
			Send:			{ |@args| args[0] is #CurrentTable ? #Test_lib : send? }
			)
		}
	}