// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_convertShortcut()
		{
		fn = MonthCalDialog.ConvertShortcut
		Assert(fn('Today', false) is: Date().NoTime())
		Assert(fn('Start of Current Month', false) is: Date().Replace(day: 1).NoTime())
		Assert(fn('End of Current Month', false)
			is: Date().Replace(day: 1).Plus(months: 1, days: -1).NoTime())
		Assert(fn('Start of Previous Month', false)
			is: Date().Replace(day: 1).Plus(months: -1).NoTime())
		Assert(fn('End of Previous Month', false)
			is: Date().Replace(day: 1).Plus(days: -1).NoTime())
		Assert(fn('Start of Current Year', false)
			is: Date().Replace(month: 1, day: 1). NoTime())
		Assert(fn('End of Current Year', false)
			is: Date().Replace(month: 12, day: 31).NoTime())
		}
	}