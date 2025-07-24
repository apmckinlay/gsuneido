// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ExtendColumns()
		{
		m = QueryHelper.ExtendColumns
		sf = SelectFields()
		Assert(m('stdlib', sf, #()) is: 'stdlib')
		Assert(m('stdlib', sf, #(field1)) is: 'stdlib extend field1 ')
		Assert(m('stdlib', sf, #(field1, num)) is: 'stdlib extend field1 ')
		Assert(m('stdlib', sf, #(num)) is: 'stdlib')
		}
	}