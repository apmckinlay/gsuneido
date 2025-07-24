// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	dupMsg(fieldName)
		{
		return 'Field name ' $ Display(fieldName) $ ' already exists.\nDuplicate ' $
				'field names are not allowed - fields must have different names.'
		}
	Test_checkForDuplicate()
		{
		fn = CustomizeScreenControl.CustomizeScreenControl_checkForDuplicate

		mock = Mock()
		mock.When.FindString('IDontExist').Return(-1)
		Assert(fn(mock, 'IDontExist') is: "")

		mock.When.FindString('ExactMatch').Return(1)
		mock.When.GetCount().Return(5)
		mock.When.GetText([anyArgs:]).Return('ExactMatch')
		Assert(fn(mock, 'ExactMatch') is: .dupMsg('ExactMatch'))

		mock.When.FindString('PartialMatch').Return(1)
		mock.When.GetCount().Return(5)
		mock.When.GetText([anyArgs:]).Return('PartialMatchWithButNotComplete')
		Assert(fn(mock, 'PartialMatch') is: "")

		mock.When.FindString('PartialMatchWithExactAfter').Return(1)
		mock.When.GetCount().Return(5)
		mock.When.GetText([anyArgs:]).Return('PartialMatchWithExactAfter3thIndex',
			'PartialMatchWithExactAfter3thIndex', 'PartialMatchWithExactAfter',
			'PartialMatchWithExactAfter3thIndex', 'PartialMatchWithExactAfter3thIndex')
		Assert(fn(mock, 'PartialMatchWithExactAfter') is:
			.dupMsg('PartialMatchWithExactAfter'))

		//Regex
		mock.When.FindString('loa').Return(1)
		mock.When.GetCount().Return(3)
		mock.When.GetText([anyArgs:]).Return('load?')
		Assert(fn(mock, 'loa') is: "")
		}

	Test_validateCustomTabName()
		{
		fn = CustomizeScreenControl.CustomizeScreenControl_validateCustomTabName
		tabsCtrl = FakeObject(GetAllTabNames: Object('Tab1', 'Tab2'))
		tabs = Object(all_tabs: Object('Tab1', 'Tab3'))
		customizable = Mock()
		customizable.When.TabCustom?('Tab4').Return(false)
		customizable.When.TabCustom?('Tab5').Return(true)
		customizable.When.TabCustom?('Tab6').Return(true)
		customizable.When.LayoutHidden?('Tab6').Return(true)
		mock = Mock()
		Assert(mock.Eval(fn, 'Tab1', tabs, tabsCtrl, customizable) is: false)
		Assert(mock.Eval(fn, 'Tab2', tabs, tabsCtrl, customizable) is: false)
		Assert(mock.Eval(fn, 'Tab3', tabs, tabsCtrl, customizable) is: false)
		Assert(mock.Eval(fn, 'Tab4', tabs, tabsCtrl, customizable))
		Assert(mock.Eval(fn, 'Tab5', tabs, tabsCtrl, customizable) is: false)
		Assert(mock.Eval(fn, 'Tab6', tabs, tabsCtrl, customizable) is: false)
		Assert(mock.Eval(fn, 'Tab6', tabs, tabsCtrl, customizable, allowHidden:))
		mock.Verify.Times(5).AlertInfo('Add Custom Tab',
			'Reserved name or Tab name already exists. Choose another name')
		}

	Test_checkRestrictionsPerTab()
		{
		testCl = CustomizeScreenControl
			{
			CustomizeScreenControl_maxTotalFields: 10
			CustomizeScreenControl_getRecord() { return Record() }
			}
		fn = testCl.CustomizeScreenControl_checkRestrictionsPerTab
		errors = Object()
		tabsControl = FakeObject(GetAllTabCount: 0)
		sf = FakeObject()
		Assert(fn(tabsControl, sf, {|msg| errors.Add(msg) }))

		testCl = CustomizeScreenControl
			{
			CustomizeScreenControl_maxTotalFields: 10
			CustomizeScreenControl_getRecord() { return Record(Tab: 'Field1 Field2') }
			}
		fn = testCl.CustomizeScreenControl_checkRestrictionsPerTab
		tabsControl = FakeObject(GetAllTabCount: 1, TabName: 'Tab')
		sf = FakeObject(FormulaFields: #('field_1', 'field_2'))
		Assert(fn(tabsControl, sf, {|msg| errors.Add(msg) }))

		testCl = CustomizeScreenControl
			{
			CustomizeScreenControl_maxTotalFields: 10
			CustomizeScreenControl_getRecord()
				{
				return Record(HeaderTab: 'Field1\r\nField2\r\nField3\r\nField4\r\n' $
					'Field5\r\nField6\r\nField7')
				}
			}
		fn = testCl.CustomizeScreenControl_checkRestrictionsPerTab
		tabsControl = FakeObject(GetAllTabCount: 1, TabName: 'HeaderTab')
		sf = FakeObject(FormulaFields: #('field_1', 'field_2', 'field_3', 'field_4',
			'field_5', 'field_6', 'field_7'))
		Assert(fn(tabsControl, sf, {|msg| errors.Add(msg) } ) is: false)
		Assert(errors.Last() is: 'You cannot add more than 6 lines of customized fields' $
			' to the Header screen.')

		testCl = CustomizeScreenControl
			{
			CustomizeScreenControl_maxTotalFields: 10
			CustomizeScreenControl_getRecord()
				{
				return Record(HeaderTab: 'Field1 Field2 Field3 Field4 Field5 Field6 ' $
					'Field7 Field8 Field9')
				}
			}
		fn = testCl.CustomizeScreenControl_checkRestrictionsPerTab
		tabsControl = FakeObject(GetAllTabCount: 1, TabName: 'HeaderTab')
		sf = FakeObject(FormulaFields: #('field_1', 'field_2', 'field_3', 'field_4',
			'field_5', 'fild_6','field_7','field_8','field_9'))
		Assert(fn(tabsControl, sf, {|msg| errors.Add(msg) } ) is: false)
		Assert(errors.Last() is: 'You cannot add more than 8 customized fields to a ' $
			'row (HeaderTab tab.)')

		testCl = CustomizeScreenControl
			{
			CustomizeScreenControl_maxTotalFields: 10
			CustomizeScreenControl_getRecord()
				{ return Record(HeaderTab: 'Field1 Field2 Field1') }
			}
		fn = testCl.CustomizeScreenControl_checkRestrictionsPerTab
		tabsControl = FakeObject(GetAllTabCount: 1, TabName: 'HeaderTab')
		sf = FakeObject(FormulaFields: #('field_1', 'field_2', 'field_1'))
		Assert(fn(tabsControl, sf, {|msg| errors.Add(msg) } ) is: false)
		Assert(errors.Last() is: 'You cannot add the same field (' $
			Prompt('field_1') $	') to a screen more than once.')

		testCl = CustomizeScreenControl
			{
			CustomizeScreenControl_maxTotalFields: 10
			CustomizeScreenControl_getRecord()
				{
				return Record(HeaderTab: 'Field1 Field2 Field3 Field4 Field5 Field6',
					SomeOtherTab: 'Field7 Field8 Field9 Field10 Field11')
				}
			}
		fn = testCl.CustomizeScreenControl_checkRestrictionsPerTab
		tabsControl = FakeObject(GetAllTabCount: 2,
			TabName: function (i) { return i is 0 ? 'HeaderTab' : 'Some Other Tab' }
			)
		sf = FakeObject(FormulaFields: function (line)
			{
			return line.Prefix?('Field1')
				? #('field_1', 'field_2', 'field_3', 'field_4', 'field_5', 'field_6')
				: #('field_7','field_8','field_9', 'field_10', 'field_11')
			})
		Assert(fn(tabsControl, sf, {|msg| errors.Add(msg) } ) is: false)
		Assert(errors.Last()
			is: 'The total number of fields on all tab layouts cannot exceed 10 ' $
				'(Currently there are 11)')
		}
	}