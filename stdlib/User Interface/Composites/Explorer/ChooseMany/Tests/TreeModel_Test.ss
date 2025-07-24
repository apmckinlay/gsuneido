// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_uniqueQuery()
		{
		m = TreeModel.TreeModel_uniqueQuery
		table = 'test_table'
		prefix = 'test_table where name is "test_class" and '

		x = [name: 'test_class', group: false, num: 1]
		Assert(m(x, table) is: prefix $ 'group is -1')

		x.group = true
		Assert(m(x, table) is: prefix $ 'parent is "" and num isnt 1 and group > -1')

		x.parent = 'test_folder'
		Assert(m(x, table)
			is: prefix $ 'parent is "test_folder" and num isnt 1 and group > -1')

		x.group = x.parent = 0
		Assert(m(x, table) is: prefix $ 'parent is 0 and num isnt 1 and group > -1')

		x.group = -1
		Assert(m(x, table) is: prefix $ 'group is -1')

		x.group = 10
		x.parent = x.num = ''
		Assert(m(x, table) is: prefix $ 'parent is "" and group > -1')
		}

	Test_copyName()
		{
		m = TreeModel.TreeModel_copyName

		Assert(m('Test_Class') is: s = 'Test_Class_Copy1')
		Assert(m(s) is: s = 'Test_Class_Copy2')
		Assert(m(s) is: s = 'Test_Class_Copy3')
		Assert(m(s) is: s = 'Test_Class_Copy4')
		Assert(m(s) is: s = 'Test_Class_Copy5')

		Assert(m('EndsWith_Copy') is: s = 'EndsWith_Copy_Copy1')
		Assert(m(s) is: s = 'EndsWith_Copy_Copy2')
		Assert(m(s) is: s = 'EndsWith_Copy_Copy3')
		Assert(m(s) is: s = 'EndsWith_Copy_Copy4')
		Assert(m(s) is: s = 'EndsWith_Copy_Copy5')

		Assert(m('DoubleDigits_Copy9') is: s = 'DoubleDigits_Copy10')
		Assert(m(s) is: s = 'DoubleDigits_Copy11')
		Assert(m(s) is: s = 'DoubleDigits_Copy12')

		Assert(m('TripleDigits_Copy99') is: s = 'TripleDigits_Copy100')
		Assert(m(s) is: s = 'TripleDigits_Copy101')
		Assert(m(s) is: s = 'TripleDigits_Copy102')
		}
	}