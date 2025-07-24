// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Get()
		{
		mock = Mock(MonthDayControl)
		mock.When.Get().CallThrough()
		mock.MonthDayControl_lists = Mock()
		mock.MonthDayControl_lists.Day = Mock()
		mock.MonthDayControl_lists.Month = Mock()

		mock.MonthDayControl_lists.Day.When.Get().Return('01')
		mock.MonthDayControl_lists.Month.When.GetCurSel().Return(0)
		Assert(mock.Get() is: '0101')

		mock.MonthDayControl_lists.Day.When.Get().Return('02')
		Assert(mock.Get() is: '0102')

		mock.MonthDayControl_lists.Day.When.Get().Return('29')
		mock.MonthDayControl_lists.Month.When.GetCurSel().Return(1)
		Assert(mock.Get() is: false)

		mock.MonthDayControl_lists.Day.When.Get().Return('31')
		mock.MonthDayControl_lists.Month.When.GetCurSel().Return(3)
		Assert(mock.Get() is: false)

		mock.MonthDayControl_lists.Month.When.GetCurSel().Return(12)
		Assert(mock.Get() is: false)

		mock.MonthDayControl_lists.Month.When.GetCurSel().Return(false)
		Assert({ mock.Get() } throws:)
		}
	}