// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	oldFmt: false
	Setup()
		{
		.oldFmt = Settings.Get('ShortDateFormat')
		Settings.Set('ShortDateFormat', "yyyy-MM-dd")
		}

	Test_ConvertToStr()
		{
		fn = ChooseDatesFormat.ConvertToStr
		Assert(fn('') is: '')

		Assert(fn("2025-12-22,2025-12-23,2025-12-24")
			is: "2025-12-22,2025-12-23,2025-12-24")

		Assert(fn("2025-12-22,2025-12-23,2025-12-24,...")
			is: "2025-12-22,2025-12-23,2025-12-24,...")

		Assert(fn("2025-12-22,2025-12-23,2025-12-2...")
			is: "2025-12-22,2025-12-23,...")

		Assert(fn("2025-12-22,2025-12-23,2025-12-...")
			is: "2025-12-22,2025-12-23,...")
		}

	Teardown()
		{
		super.Teardown()
		if .oldFmt isnt false
			Settings.Set('ShortDateFormat', .oldFmt)
		}
	}