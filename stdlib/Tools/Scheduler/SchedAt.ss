// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
class
	{
	rx: '^at (\d\d?:\d\d)([ \t]+skip weekends)?$'
	skipRx: 'skip weekends$'
	CallClass(when)
		{
		at = when.Extract(.rx)
		skipWeekends = (when =~ .skipRx)
		return at is false or Date(at) is false ? false : new this(at, skipWeekends)
		}
	dateWidth: 4
	New(at, .skipWeekends)
		{
		.at = at.Tr(':').LeftFill(.dateWidth, '0')
		}
	Due?(prevcheck, curtime)
		{
		if prevcheck is false
			return false

		return .between(prevcheck, .at, curtime) and
			not .skipWeekend?(prevcheck, .at, curtime)
		}
	between(prevcheck, at, curtime)
		{
		prev = prevcheck.Format('HHmm')
		cur = curtime.Format('HHmm')
		if cur[..2] < prev[..2] // passed midnight
			{ // e.g. prevcheck is 23:58, curtime is 00:02, at is 23:59 or 00:01
			cur = '25' $ cur[2..]
			at = .at.Replace('^00', '25')
			}
		return prev < at and at <= cur
		}

	skipWeekend?(prevcheck, at, curtime)
		{
		if .skipWeekends isnt true
			return false
		return .isWeekend?(prevcheck, at, curtime)
		}
	friday: 4
	isWeekend?(prevcheck, at, curtime)
		{
		// using monday so we can use "> 4" instead of "is 0 or 6"
		// (just to keep code cleaner)
		curDayWeek = curtime.WeekDay('Mon')
		prevDayWeek = prevcheck.WeekDay('Mon')
		// if we get here we know 'at' is between prev and cur
		if curDayWeek <= .friday and prevDayWeek <= .friday
			return false
		if curDayWeek > .friday and prevDayWeek > .friday
			return true

		// If we get here then prev and Cur MUST be overlapping a midnight
		// 'at' will be greater than prev if 'at' is before midnight
		// i.e. if prev is 2355 and cur is 0015
		//		if 'at' is 2359 then 'at' > prev
		//		if 'at' is 0005 then 'at' < prev
		return at > prevcheck.Format('HHmm')
			? prevDayWeek > .friday
			: curDayWeek > .friday
		}
	}
