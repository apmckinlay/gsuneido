// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_hourSelect()
		{
		fn = DateTimeControl.DateTimeControl_hourSelect

		// Testing 24 hour
		Assert(fn(0, 'HH:mm') is: '0')
		Assert(fn(6, 'HH:mm') is: '6')
		Assert(fn(11, 'HH:mm') is: '11')
		Assert(fn(12, 'HH:mm') is: '12')
		Assert(fn(23, 'HH:mm') is: '23')

		Assert(fn(0, 'H:mm') is: '0')
		Assert(fn(6, 'H:mm') is: '6')
		Assert(fn(11, 'H:mm') is: '11')
		Assert(fn(12, 'H:mm') is: '12')
		Assert(fn(23, 'H:mm') is: '23')

		// Testing AM/PM
		Assert(fn(0, 'hh:mm tt') is: '12 am')
		Assert(fn(6, 'hh:mm tt') is: '6 am')
		Assert(fn(11, 'hh:mm tt') is: '11 am')
		Assert(fn(12, 'hh:mm tt') is: '12 pm')
		Assert(fn(23, 'hh:mm tt') is: '11 pm')

		Assert(fn(0, 'h:mm tt') is: '12 am')
		Assert(fn(6, 'h:mm tt') is: '6 am')
		Assert(fn(11, 'h:mm tt') is: '11 am')
		Assert(fn(12, 'h:mm tt') is: '12 pm')
		Assert(fn(23, 'h:mm tt') is: '11 pm')
		}
	}