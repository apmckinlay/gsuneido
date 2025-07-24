// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_DisplayValues()
		{
		t1 = Timestamp()
		t2 = Timestamp()
		table = .MakeTable('(biztest_num, biztest_name, biztest_abbrev, biztest_display)
			key(biztest_num)',
			[biztest_num: t1, biztest_name: 't1 name', biztest_abbrev: 't1_abbrev',
				biztest_display: 't1 display'],
			[biztest_num: t2, biztest_name: 't2 name', biztest_abbrev: 't2_abbrev',
				biztest_display: 't2 display'])

		// test default
		Assert(IdControl.DisplayValues([#Id, table, 'biztest_num'], [t1, t2])
			is: #('t1 name', 't2 name'))
		// test named args
		Assert(IdControl.DisplayValues([#Id, table, field: 'biztest_num'], [t1, t2])
			is: #('t1 name', 't2 name'))
		Assert(IdControl.DisplayValues([#Id, query: table,
			field: 'biztest_num'], [t1, t2]) is: #('t1 name', 't2 name'))
		// test nameField
		Assert(IdControl.DisplayValues([#Id, query: table,
			field: 'biztest_num', nameField: 'biztest_display'], [t1, t2])
			is: #('t1 display', 't2 display'))
		// test allow other
		Assert(IdControl.DisplayValues([#Id, table, 'biztest_num',
			nameField: 'biztest_display', allowOther:], [t1, t2, 'TEST STRING'])
			is: #('t1 display', 't2 display', 'TEST STRING'))
		}
	}