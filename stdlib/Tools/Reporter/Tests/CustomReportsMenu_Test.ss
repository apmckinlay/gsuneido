// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_getBookLocation()
		{
		fn = CustomReportsMenu.CustomReportsMenu_getBookLocation
		Assert(fn(Object()) is: false)
		Assert(fn(Object(bookLocation: 'somewhereLocation')) is: 'somewhereLocation')
		Assert(fn(Object(tables: #())) is: false)
		Assert(fn(Object(tables: #('nonexistent'))) is: false)

		.SpyOn(Tables.GetTable).Return(Object(BookLocation: 'onBook'))
		Assert(fn(Object(tables: #('TabOne'))) is: 'onBook')
		Assert(fn(Object(tables: #('TabOne', 'TabTwo'))) is: 'onBook')
		Assert(fn(Object(bookLocation: 'somewhere', tables: #('TabOne'))) is: 'somewhere')
		}
	}