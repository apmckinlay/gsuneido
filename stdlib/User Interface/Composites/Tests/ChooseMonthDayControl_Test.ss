// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_DateFromMonthDay()
		{
		fn = ChooseMonthDayControl.DateFromMonthDay

		Assert(fn('abc') is: false)
		Assert(fn('0101') is: #20010101)
		Assert(fn('0228') is: #20010228)
		Assert(fn('0229') is: false)
		Assert(fn('dec31') is: false)
		}

	Test_Valid?()
		{
		mock = Mock(ChooseMonthDayControl)
		mock.When.Valid?().CallThrough()
		mock.ChooseMonthDayControl_mandatory = false

		mock.When.Get().Return('0101')
		Assert(mock.Valid?())

		mock.When.Get().Return('')
		Assert(mock.Valid?())

		mock.When.Get().Return('1231')
		Assert(mock.Valid?())

		mock.When.Get().Return('1232')
		Assert(mock.Valid?() is: false)

		mock.When.Get().Return('101')
		Assert(mock.Valid?() is: false)

		mock.When.Get().Return('00101')
		Assert(mock.Valid?() is: false)

		mock.When.Get().Return('text')
		Assert(mock.Valid?() is: false)

		mock.When.Get().Return('0229')
		Assert(mock.Valid?() is: false)
		}
	}