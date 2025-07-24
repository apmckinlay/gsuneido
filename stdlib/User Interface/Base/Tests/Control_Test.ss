// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		super.Setup()
		.MakeLibraryRecord(
			[name: 'Field_test_field', text: 'class
				{
				Prompt: "Test Field"
				Control: (Field, width: 20)
				}'],
			[name: 'Field_test_checkbox', text: 'Field_boolean
				{
				Prompt: "Checkbox"
				}'])
		}

	Test_Construct()
		{
		fn = function (@x) { return Control.Control_build(x) }
		Assert(fn(Object('Field')) is: Object('Field'))
		Assert(fn(Object('Field', width: 20)) is: Object('Field', width: 20))
		Assert(fn('Field') is: Object('Field'))
		Assert(fn('Field', width: 20) is: Object('Field', width: 20))

		Assert(fn('NoPrompt', 'test_field')
			is: #('Field', width: 20, name: 'test_field'))
		Assert(fn('NoPrompt', 'test_field', set: 'abc')
			is: #('Field', width: 20, set: 'abc', name: 'test_field'))

		Assert(fn('test_field')
			is: #('Pair',
				#('Static', 'Test Field', hidden: false),
				#('Field', width: 20, name: 'test_field')))
		Assert(fn('test_field', set: 'abc')
			is: #('Pair',
				#('Static', 'Test Field', hidden: false),
				#('Field', width: 20, set: 'abc', name: 'test_field')))

		Assert(fn('test_checkbox')
			is: #('CheckBox', text: 'Checkbox', hidden: false, name: 'test_checkbox'))
		Assert(fn('test_checkbox', readonly:)
			is: #('CheckBox', text: 'Checkbox', hidden: false, readonly:,
				name: 'test_checkbox'))
		}

	Test_Defer()
		{
		if not Sys.Win32?()
			return
		c = new Control { New2(){} }
		for ..10
			c.Defer({}, uniqueID: "x")
		Assert(c.Control_timers.Size() is: 1) // should not accumulate
		Assert(c.Control_unique.Size() is: 1) // should not accumulate
		}
	}