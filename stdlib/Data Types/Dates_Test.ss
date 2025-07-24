// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// this is tests for methods defined in Dates
// DateTest is for built-in methods
// SuJsWebTest
Test
	{
	Test_DayOfYear()
		{
		Assert((Date().NoTime().Replace(month: 1, day: 1)).DayOfYear() is: 1)
		Assert(#20030220.DayOfYear() is: 51)
		}

	Test_EndOfDay()
		{
		Assert(#20030220.123456789.EndOfDay() is: #20030220.235959999)
		}

	Test_Replace()
		{
		Assert(#20021224.Replace(month: 1) is: #20020124)
		}

	Test_EndOfMonth()
		{
		Assert(#20030215.EndOfMonth() is: #20030228)
		Assert(#20030315.EndOfMonth() is: #20030331)
		Assert(#20030415.EndOfMonth() is: #20030430)
		Assert(#20030531.EndOfMonth() is: #20030531)

		Assert(#20030215.EndOfMonthDay() is: 28)
		Assert(#20030315.EndOfMonthDay() is: 31)
		Assert(#20030415.EndOfMonthDay() is: 30)
		Assert(#20030531.EndOfMonthDay() is: 31)
		}

	Test_WeekNumber()
		{
		Assert(#20030331.WeekNumber() is: 14)
		Assert(#20030331.WeekNumber(0) is: 14)
		Assert(#20030331.WeekNumber(1) is: 14)
		Assert(#20030331.WeekNumber('mon') is: 14)
		Assert(#20030331.WeekNumber('Monday') is: 14)
		Assert(#20031231.WeekNumber(1) is: 53)
		Assert(#20121231.WeekNumber(1) is: 54)

		Assert(#20180101.WeekNumber() is: 1)
		for i in ..7
			Assert(#20180101.WeekNumber(i) is: 1)
		Assert(#20181108.WeekNumber() is: 45)
		Assert(#20181108.WeekNumber(2) is: 46)
		}

	Test_IsoWeekDay()
		{
		Assert(#20030331.IsoWeekDay() is: "2003-W14-1")
		Assert(#20031231.IsoWeekDay() is: "2004-W01-3")
		Assert(#20040101.IsoWeekDay() is: "2004-W01-4")
		Assert(#20091231.IsoWeekDay() is: "2009-W53-4")
		Assert(#20100101.IsoWeekDay() is: "2009-W53-5")
		}

	Test_StdShortDate()
		{
		date = #20050501.143654111
		Assert(date.StdShortDate() is: "2005-05-01")
		}

	Test_StdShortDateTime()
		{
		date = #20050501.143654111
		Assert(date.StdShortDateTime() is: "2005-05-01 14:36")
		}

	TestStdShortDateTimeSec()
		{
		date = #20050501.143654111
		Assert(date.StdShortDateTimeSec() is: "2005-05-01 14:36:54")
		}

	Test_Quarter()
		{
		Assert(#20120101.Quarter() is: 1)
		Assert(#20120201.Quarter() is: 1)
		Assert(#20120301.Quarter() is: 1)
		Assert(#20120401.Quarter() is: 2)
		Assert(#20120501.Quarter() is: 2)
		Assert(#20120601.Quarter() is: 2)
		Assert(#20120701.Quarter() is: 3)
		Assert(#20120801.Quarter() is: 3)
		Assert(#20120901.Quarter() is: 3)
		Assert(#20121001.Quarter() is: 4)
		Assert(#20121101.Quarter() is: 4)
		Assert(#20121201.Quarter() is: 4)
		}

	Test_NoTime?()
		{
		Assert(#20200101.NoTime?())
		Assert(#20200101.000000001.NoTime?() is: false)
		Assert(#20200121.135612545.NoTime?() is: false)
		}

	Test_UTC()
		{
		Assert(Dates.UTC(-721) 	is: 'Invalid GMT Bias: -721')
		Assert(Dates.UTC(-720) 	is: 'UTC-12:00')
		Assert(Dates.UTC(0) 	is: 'UTC=00:00')
		Assert(Dates.UTC(15) 	is: 'UTC+00:15')
		Assert(Dates.UTC(-15) 	is: 'UTC-00:15')
		Assert(Dates.UTC(60) 	is: 'UTC+01:00')
		Assert(Dates.UTC(-60) 	is: 'UTC-01:00')
		Assert(Dates.UTC(360) 	is: 'UTC+06:00')
		Assert(Dates.UTC(-360) 	is: 'UTC-06:00')
		Assert(Dates.UTC(405) 	is: 'UTC+06:45')
		Assert(Dates.UTC(-405) 	is: 'UTC-06:45')
		Assert(Dates.UTC(840) 	is: 'UTC+14:00')
		Assert(Dates.UTC(841) 	is: 'Invalid GMT Bias: 841')
		}

	Test_FromUnix()
		{
		// SuJsWebTest Excluded
		Assert(Dates.FromUnix(1598369968).GMTime() is: #20200825.153928)
		// ensure ms is handled for satellites
		Assert(Dates.FromUnix(409778802.128).GMTime() is: #19821226.192642128)
		}
	}
