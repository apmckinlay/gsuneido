// Copyright (C) 2016 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(task)
		{
		for schedType in #(daily, interval, weekly, monthly)
			if task['run' $ schedType] is true and
				false isnt result =
					this['Convert_' $ schedType](task)
					return result
		SuneidoLog("ERROR: could not translate task", calls:)
		return false
		}

	Convert_daily(task)
		{
		return .convertTimeFormat(task.daily_time)
		}

	Convert_weekly(task)
		{
		if not Number?(task.weekly_day) or
			task.weekly_day < 0 or task.weekly_day > 6 or /*= 7 days a week */
			false is at = .convertTimeFormat(task.weekly_time)
			return false

		dayOb = #(Sun Mon Tue Wed Thu Fri Sat)
		return "on " $ dayOb[task.weekly_day] $ ' ' $ at
		}

	Convert_interval(task)
		{
		// minutes, hours, days, and weeks all expressed in minutes for SchedEvery class
		intervals = #(1: 1, 2: 60, 3: 1440, 4: 10080)
		if not intervals.Member?(task.interval_units) or not Number?(task.interval) or
			not task.interval.Int?() or task.interval <= 0
			return false

		return "every " $ (task.interval * intervals[task.interval_units]) $ ' minutes'
		}

	Convert_monthly(task)
		{
		if not Number?(task.monthly_day) or
			(task.monthly_day < 0 or task.monthly_day > 29) or /*= day of month */
			false is at = .convertTimeFormat(task.monthly_time)
			return false

		on = String(task.monthly_day - 1)
		if task.monthly_day is 0
			on = 'StartMonth'
		else if task.monthly_day is 1
			on = 'EndMonth'
		else if task.monthly_day is 2
			on = 'MidMonth'
		return 'on ' $ on $ ' ' $ at
		}

	convertTimeFormat(time)
		{
		skipWeekends = false
		if String?(time)
			{
			skipWeekends = (time =~ 'skip weekends$')
			time = time.BeforeFirst(' skip weekends')
			}
		else
			time = String(time)
		if '' is hour = time[..-2]
			hour = '0'
		if '' is minute = time[-2..]
			minute = '0'
		return .atTime(minute, hour, skipWeekends)
		}

	atTime(minute, hour, skipWeekends)
		{
		minute = minute.LeftFill(2, '0')
		if not hour.Numeric?() or Number(hour) > 23 or Number(hour) < 0 /*= last hour */
			return false
		if not minute.Numeric?() or Number(minute) > 59 or Number(minute) < 0/*=last min*/
			return false
		return 'at ' $ hour $ ':' $ minute $ (skipWeekends ? ' skip weekends' : '')
		}
	}
