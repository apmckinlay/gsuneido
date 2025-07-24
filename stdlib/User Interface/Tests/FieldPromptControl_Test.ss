// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.uppercase = 'Field_test_uppercase' $ Display(Timestamp()).Replace("#|\.", "")
		.lowercase = 'Field_test_lowercase' $ Display(Timestamp()).Replace("#|\.", "")
		.noPrompts = 'Field_test_noPrompts' $ Display(Timestamp()).Replace("#|\.", "")
		.MakeLibraryRecord([
			name: .uppercase,
			text: 'Field_number_custom
				{
				Prompt: "Uppercase"
				Control_mask: "-###,###,###.##"
				Format_mask: "-###,###,###.##"
				}'])
		.MakeLibraryRecord([
			name: .lowercase,
			text: 'Field_number_custom
				{
				Prompt: "lowercase"
				Control_mask: "-###,###,###.##"
				Format_mask: "-###,###,###.##"
				}'])
		.MakeLibraryRecord([
			name: .noPrompts,
			text: 'Field_number_custom
				{
				Control_mask: "-###,###,###.##"
				Format_mask: "-###,###,###.##"
				}'])
		}

	Test_allowAddField?()
		{
		fn = FieldPromptControl.FieldPromptControl_allowAddField?
		Assert(fn('num', #(), #()))
		Assert(fn('num', #('num'), #()) is: false)
		Assert(fn('num', #(), #(num)) is: false)
		Assert(fn('num_internal', #(), #()) is: false)
		}

	Test_SetFieldMap()
		{
		mock = Mock(FieldPromptControl)
		fields = mock.FieldPromptControl_map = Object()
		mock.When.SetList([anyArgs:]).Return(true)
		mock.When.FieldPromptControl_add_field([anyArgs:]).CallThrough()

		mock.Eval(FieldPromptControl.SetFieldMap, :fields)
		list = mock.FieldPromptControl_map
		Assert(list is: Object())

		upper = .uppercase.AfterFirst("_")
		lower = .lowercase.AfterFirst("_")
		noPrompts = .noPrompts.AfterFirst("_")
		noDataDict = "test_noDataDict" $ Display(Timestamp()).Replace("#|\.", "")
		fields.Add(upper, lower, noPrompts, noDataDict)

		mock.Eval(FieldPromptControl.SetFieldMap, :fields)
		list = mock.FieldPromptControl_map

		mock.Verify.Times(4).FieldPromptControl_add_field([anyArgs:])
		Assert(list isSize: 2)

		Assert(list hasMember: "Uppercase")
		Assert(list hasMember: "lowercase")
		Assert(list has: upper)
		Assert(list has: lower)

		Assert(list hasntMember: noPrompts)
		Assert(list hasntMember: noDataDict)
		Assert(list hasnt: noPrompts)
		Assert(list hasnt: noDataDict)
		}

	Test_validateParams()
		{
		fn = FieldPromptControl.FieldPromptControl_validateParams

		Assert(fn(#(), false, 'listField'), msg: 'only listfield')
		Assert(fn(#(value1, value2), false, false), msg: 'only fields')
		Assert(fn(#(), 'tableName', false), msg: 'only table')
		Assert(fn(#(value1, value2), 'tableName', false), msg: 'fields table')

		msg = "ERROR: FieldPromptControl - listField cannot be used in combination " $
			"with table or fields"
		ServerSuneido.Set('TestRunningExpectedErrors', Object(msg, msg, msg))
		Assert(fn(#(value1, value2), 'tableName', 'listField') is: false, msg: 'all')
		Assert(fn(#(value1, value2), false, 'listField') is: false, msg: 'fields list')
		Assert(fn(#(), 'tableName', 'listField') is: false, msg: 'table listfield')
		}
	}