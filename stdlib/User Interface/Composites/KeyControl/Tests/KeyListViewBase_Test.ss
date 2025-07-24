// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_removeUnwantedKeys()
		{
		func = KeyListViewBase.KeyListViewBase_removeUnwantedKeys

		keys = Object()
		func(keys)
		Assert(keys is: #())

		keys = Object('test1_num')
		func(keys)
		Assert(keys is: #())

		keys = Object('test1_num_new')
		func(keys)
		Assert(keys is: #())

		keys = Object('test1_field1, test1_field2')
		func(keys)
		Assert(keys is: #())

		keys = Object('test1_field1, test1_field2', 'test1_num')
		func(keys)
		Assert(keys is: #())

		keys = Object('test1_name')
		func(keys)
		Assert(keys is: #('test1_name'))

		keys = Object('test1_name', 'test1_abbrev')
		func(keys)
		Assert(keys is: #('test1_name', 'test1_abbrev'))

		keys = Object('test1_num_new', 'test1_name')
		func(keys)
		Assert(keys is: #('test1_name'))

		keys = Object('test1_name', 'test1_name_lower!')
		func(keys)
		Assert(keys is: #('test1_name_lower!'))

		keys = Object('fred', 'test1_name_lower!')
		func(keys)
		Assert(keys is: #('fred', 'test1_name_lower!'))

		keys = Object('test1_name', 'test1_name_lower!', 'test1_abbrev',
			'test1_abbrev_lower!')
		func(keys)
		Assert(keys is: #('test1_name_lower!', 'test1_abbrev_lower!'))

		keys = Object('test1_name', 'test1_name_lower!', 'test1_abbrev',
			'test1_abbrev_lower!', 'test1_num', 'test1_field1, test1_field2', 'fred')
		func(keys)
		Assert(keys is: #('test1_name_lower!', 'test1_abbrev_lower!', 'fred'))
		}

	Test_optionalRestrictions()
		{
		.MakeLibraryRecord([name: 'Field_testfield_status',
			text: 'Field_string
				{
				Prompt: "Status"
				Control: (ChooseList #(active, inactive), mandatory:, width: 6)
				}'])
		.MakeLibraryRecord([name: 'Field_testfield_expiry',
			text: 'Field_date
				{
				Prompt: "Test Expiry Date"
				}'])
		.MakeLibraryRecord([name: 'Field_testfield_terminated',
			text: 'Field_date
				{
				Prompt: "Test Terminated"
				}'])

		date = Display(Date().NoTime())
		mock = Mock(KeyListViewBase)
		mock.When.KeyListViewBase_addCheckBoxAndQuery([anyArgs:]).CallThrough()
		mock.When.ChooseDateCtrl?([anyArgs:]).CallThrough()
		mock.When.BuildChooseDateQueryWhere([anyArgs:]).CallThrough()
		mock.When.ChooseListActiveInactive?([anyArgs:]).CallThrough()
		mock.When.BuildChooseListActiveInactiveWhere([anyArgs:]).CallThrough()
		mock.KeyListViewBase_optionalRestrictions = #(testfield_status, testfield_expiry,
			testfield_terminated)
		mock.KeyListViewBase_whereMap = Object()
		layout = mock.Eval(KeyListViewBase.KeyListViewBase_optionalRestrictionLayout)
		Assert(layout[1][1].name is: "show_inactive")
		whereMap = mock.KeyListViewBase_whereMap
		Assert(whereMap.show_inactive is: 'testfield_status is "active"')
		Assert(whereMap.show_terminated is: '(testfield_terminated is ""' $
			' or testfield_terminated >= ' $ date $ ')')
		Assert(whereMap.show_expired is: '(testfield_expiry is ""' $
			' or testfield_expiry >= ' $	date $ ')')

		updateFn = KeyListViewBase.KeyListViewBase_updateExtraWhere
		mock.Eval(updateFn, [])
		Assert(mock.Eval(KeyListViewBase.GetExtraWhere) is:
			' where (testfield_expiry is "" or testfield_expiry >= ' $ date $
			') and testfield_status is "active" and ' $
			'(testfield_terminated is "" or testfield_terminated >= ' $ date $ ')')

		mock.Eval(updateFn, [show_expired:])
		Assert(mock.Eval(KeyListViewBase.GetExtraWhere)
			is: ' where testfield_status is "active"' $
			' and (testfield_terminated is "" or testfield_terminated >= ' $ date $ ')')

		mock.Eval(updateFn, [show_expired:, show_terminated:])
		Assert(mock.Eval(KeyListViewBase.GetExtraWhere)
			is: ' where testfield_status is "active"')

		mock.Eval(updateFn, [show_expired:, show_terminated:, show_inactive:])
		Assert(mock.Eval(KeyListViewBase.GetExtraWhere) is: '')
		}
	}