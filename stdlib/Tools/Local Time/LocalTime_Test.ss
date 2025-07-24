// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_calc_offset()
		{
		.testOffset('SK', 'SK', 0, 0)
		.testOffset('AB', 'SK', -1, 0)
		.testOffset('SK', 'DE', -1, 0)
		.testOffset('AK', 'NB', -5, 0)
		.testOffset('NB', 'AK', 5, 0)
		.testOffset('AK', 'NL', -5, -30)
		.testOffset('NL', 'AK', 5, 30)
		.testOffset('NL', 'NS', 0, 30)
		.testOffset('AB', 'AB', 0, 0)
		.testOffset('AB', 'BC', 1, 0)
		.testOffset('BC', 'AB', -1, 0)

		result = LocalTime.LocalTime_calc_offset('AST', 'EST')
		Assert(result.hours is: 1)
		Assert(result.minutes is: 0)

		result = LocalTime.LocalTime_calc_offset('EST', 'AST')
		Assert(result.hours is: -1)
		Assert(result.minutes is: 0)

		result = LocalTime.LocalTime_calc_offset('NST', 'AST')
		Assert(result.hours is: 0)
		Assert(result.minutes is: 30)

		result = LocalTime.LocalTime_calc_offset('AST', 'NST')
		Assert(result.hours is: 0)
		Assert(result.minutes is: -30)
		}

	testOffset(toTimeZone, fromTimeZone, offsetHours, offsetMinutes)
		{
		msgPrefix = fromTimeZone $ ' - ' $ toTimeZone
		toRec = TimeZones[toTimeZone]
		fromRec = TimeZones[fromTimeZone]
		result = LocalTime.LocalTime_calc_offset(toRec.zone, fromRec.zone)
		Assert(result.hours is: offsetHours msg: msgPrefix $ ' offset hours')
		Assert(result.minutes is: offsetMinutes msg: msgPrefix $ ' offset minutes')
		}

	Test_calc_daylightsavings()
		{
		.testDaylightSavings('AB', 'SK', #19000101.1200, #19000101.1100, #19000101.1100)
		.testDaylightSavings('AB', 'SK', #19000325.1200, #19000325.1100, #19000325.1200)
		.testDaylightSavings('SK', 'AB', #19000101.1200, #19000101.1300, #19000101.1300)
		.testDaylightSavings('SK', 'AB', #19000325.1200, #19000325.1300, #19000325.1200)
		.testDaylightSavings('MB', 'SK', #19000101.1200, #19000101.1200, #19000101.1200)
		.testDaylightSavings('MB', 'SK', #19000325.1200, #19000325.1200, #19000325.1300)
		.testDaylightSavings('SK', 'MB', #19000101.1200, #19000101.1200, #19000101.1200)
		.testDaylightSavings('SK', 'MB', #19000325.1200, #19000325.1200, #19000325.1100)
		}

	testDaylightSavings(toTimeZone, fromTimeZone, curTime, expectedCurTime, expectedRes)
		{
		msgPrefix = fromTimeZone $ ' - ' $ toTimeZone
		toRec = TimeZones[toTimeZone]
		fromRec = TimeZones[fromTimeZone]
		offset = LocalTime.LocalTime_calc_offset(toRec.zone, fromRec.zone)
		curTime = curTime.Plus(hours: offset.hours, minutes: offset.minutes)
		Assert(curTime is: expectedCurTime msg: msgPrefix $ ' offset')
		result = LocalTime.LocalTime_calc_daylightsavings(curTime, toRec, fromRec)
		Assert(result is: expectedRes
			msg: msgPrefix $ ' ' $ Display(curTime.NoTime()) $ ' after daylight')
		}
	}
