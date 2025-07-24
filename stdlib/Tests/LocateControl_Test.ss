// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		.TearDownIfTablesNotExist(UserSettings.Table)
		}
	Test_layout_none()
		{
		mock = Mock()
		mock.LocateControl_query = ""
		for keys in #((), (""), ("": ""))
			{
			mock.LocateControl_keys = keys
			Assert(mock.Eval(LocateControl.LocateControl_layout,
				false, #(), #(), false, #()) is: #(Static ''))
			}

		}
	Test_layout()
		{
		// note: due to Mock, this doesn't test useLowerKeys
		mock = Mock(LocateControl)
		mock.When.getAbbrevField([anyArgs:]).CallThrough()
		mock.LocateControl_customizeQueryCols = false
		mock.LocateControl_query = "biz_partners"
		mock.LocateControl_keys = #(Abbrev: bizpartner_abbrev)

		// single key
		layout = mock.Eval(LocateControl.LocateControl_layout, false, #(), #(), false,
			#())
		Assert(layout has: #(Static Abbrev name: 'LocateBy'))

		keys =  #(Abbrev: bizpartner_abbrev, Name: bizpartner_name)
		mock.LocateControl_keys = keys

		// multiple keys, no saved setting so uses first
		layout = mock.Eval(LocateControl.LocateControl_layout, false, #(), #(), false,
			#())
		Assert(layout has: Object('ChooseList', keys.Members(),
			set: keys.Members().Sort!()[0], name: 'LocateBy'))

		// saved settings
		for field in #(bizpartner_abbrev, bizpartner_name)
			{
			UserSettings.Put('locateby:LocateControl_Test', field)
			layout = mock.Eval(LocateControl.LocateControl_layout, false, #(), #(),
				'LocateControl_Test', #())
			Assert(layout has: Object('ChooseList', keys.Members(),
				set: keys.Find(field), name: 'LocateBy'))
			}

		// invalid saved setting
		UserSettings.Put('locateby:LocateControl_Test', 'invalid!')
		layout = mock.Eval(LocateControl.LocateControl_layout, false, #(), #(),
			'LocateControl_Test', #())
		Assert(layout has: Object('ChooseList', keys.Members(),
			set: keys.Members().Sort!()[0], name: 'LocateBy'))
		}
	Test_useLowerKeys()
		{
		mock = Mock()

		keys = Object(Abbrev: 'bizpartner_abbrev', Name: 'bizpartner_name')
		mock.LocateControl_keys = keys
		cols = #(bizpartner_abbrev, bizpartner_name)
		mock.Eval(LocateControl.LocateControl_useLowerKeys, cols)
		Assert(mock.LocateControl_keys is: keys)
		cols = #(bizpartner_abbrev, bizpartner_name, bizpartner_name_lower!)
		mock.Eval(LocateControl.LocateControl_useLowerKeys, cols)
		Assert(mock.LocateControl_keys is:
			#(Abbrev: bizpartner_abbrev, 'Name*': bizpartner_name_lower!))
		}

	Test_getAbbrevField()
		{
		abbrev = LocateControl.LocateControl_getAbbrevField

		Assert(abbrev('bizpartner_name', #(bizpartner_name)) is: false)

		Assert(abbrev('bizpartner_abbrev', #(bizpartner_abbrev)) is: false)

		Assert(abbrev('bizpartner_abbrev_lower!',
			#(bizpartner_abbrev, bizpartner_abbrev_lower!))
			is: false)

		Assert(abbrev('bizpartner_name', #(bizpartner_name, bizpartner_abbrev))
			is: 'bizpartner_abbrev')

		Assert(abbrev('bizpartner_name_lower!',
			#(bizpartner_name, bizpartner_name_lower!, bizpartner_abbrev))
			is: 'bizpartner_abbrev')
		}

	Teardown()
		{
		UserSettings.Remove('locateby:LocateControl_Test')
		super.Teardown()
		}
	}
