// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		ServerSuneido.Set('TestRunningExpectedErrors', Object(.errMsg))
		Assert({ .fn(false) } throws: 'boolean does not support get')
		Assert(.fn([]) is: false)
		// does not fall through
		Assert(.fn([rundaily: true, daily_time: 300, suspended: false]) isnt: false)
		}

	fn(task)
		{
		return TranslateScheduledTask(task)
		}
	errMsg: "ERROR: could not translate task"

	TestDaily()
		{
		ServerSuneido.Set('TestRunningExpectedErrors',
			Object(.errMsg, .errMsg, .errMsg, .errMsg))
		Assert(.fn([rundaily: true, daily_time: 2360]) is: false)
		// midnight represented as 0 not 2400
		Assert(.fn([rundaily: true, daily_time: 2400]) is: false)
		Assert(.fn([rundaily: true, daily_time: 10000]) is: false)
		Assert(.fn([rundaily: true, daily_time: -5]) is: false)

		Assert(.fn([rundaily: true, daily_time: 0]) is: "at 0:00")
		Assert(.fn([rundaily: true, daily_time: 5]) is: "at 0:05")
		Assert(.fn([rundaily: true, daily_time: 50]) is: "at 0:50")
		Assert(.fn([rundaily: true, daily_time: 500]) is: "at 5:00")
		Assert(.fn([rundaily: true, daily_time: 505]) is: "at 5:05")
		Assert(.fn([rundaily: true, daily_time: 1050]) is: "at 10:50")
		Assert(.fn([rundaily: true, daily_time: 1005]) is: "at 10:05")
		Assert(.fn([rundaily: true, daily_time: 1000]) is: "at 10:00")
		Assert(.fn([rundaily: true, daily_time: 2359]) is: "at 23:59")
		}

	TestDailySkipWeekends()
		{
		Assert(.fn([rundaily: true, daily_time: '1000 skip weekends'])
			is: 'at 10:00 skip weekends')
		Assert(.fn([rundaily: true, daily_time: '0 skip weekends'])
			is: 'at 0:00 skip weekends')
		}

	TestWeekly()
		{
		ServerSuneido.Set('TestRunningExpectedErrors',
			Object(.errMsg, .errMsg, .errMsg, .errMsg))
		Assert(.fn([runweekly: true, weekly_time: 2359]) is: false)
		Assert(.fn([runweekly: true, weekly_day: -1, weekly_time: 2359]) is: false)
		Assert(.fn([runweekly: true, weekly_day: 7, weekly_time: 2359]) is: false)
		Assert(.fn([runweekly: true, weekly_day: 1, weekly_time: 2360]) is: false)
		Assert(.fn([runweekly: true, weekly_day: 1, weekly_time: 2340])
			is: "on Mon at 23:40")
		Assert(.fn([runweekly: true, weekly_day: 6, weekly_time: 340])
			is: "on Sat at 3:40")
		}

	TestMonthly()
		{
		ServerSuneido.Set('TestRunningExpectedErrors',
			Object(.errMsg, .errMsg, .errMsg))
		Assert(.fn([runmonthly: true, monthly_time: 2359]) is: false)
		Assert(.fn([runmonthly: true, monthly_day: 1, monthly_time: 2360]) is: false)
		Assert(.fn([runmonthly: true, monthly_day: 0, monthly_time: 2340])
			is: "on StartMonth at 23:40")
		Assert(.fn([runmonthly: true, monthly_day: 1, monthly_time: 2340])
			is: "on EndMonth at 23:40")
		Assert(.fn([runmonthly: true, monthly_day: 2, monthly_time: 340])
			is: "on MidMonth at 3:40")
		Assert(.fn([runmonthly: true, monthly_day: 3, monthly_time: 340])
			is: "on 2 at 3:40")
		Assert(.fn([runmonthly: true, monthly_day: 29, monthly_time: 340])
			is: "on 28 at 3:40")
		Assert(.fn([runmonthly: true, monthly_day: 30, monthly_time: 340]) is: false)
		}

	minutesInHour: 60
	minutesInDay: 1440
	minutesInWeek: 10080
	TestInterval()
		{
		ServerSuneido.Set('TestRunningExpectedErrors',
			Object(.errMsg, .errMsg, .errMsg, .errMsg, .errMsg, .errMsg, .errMsg))
		Assert(.fn([runinterval: true, interval: 1]) is: false)
		Assert(.fn([runinterval: true, interval_units: 1]) is: false)
		Assert(.fn([runinterval: true, interval: 1, interval_units: 'bob']) is: false)
		Assert(.fn([runinterval: true, interval: 'bob', interval_units: 1]) is: false)
		Assert(.fn([runinterval: true, interval: 1, interval_units: 6]) is: false)
		Assert(.fn([runinterval: true, interval: 1.5, interval_units: 1]) is: false)
		Assert(.fn([runinterval: true, interval: 1, interval_units: 1])
			is: "every 1 minutes")
		Assert(.fn([runinterval: true, interval: 0, interval_units: 1]) is: false)
		Assert(.fn([runinterval: true, interval: 1, interval_units: 2])
			is: "every " $ .minutesInHour $ " minutes")
		Assert(.fn([runinterval: true, interval: 1, interval_units: 3])
			is: "every " $ .minutesInDay $ " minutes")
		Assert(.fn([runinterval: true, interval: 1, interval_units: 4])
			is: "every " $ .minutesInWeek $ " minutes")
		Assert(.fn([runinterval: true, interval: 5, interval_units: 1])
			is: "every 5 minutes")
		Assert(.fn([runinterval: true, interval: 5, interval_units: 2])
			is: "every " $ .minutesInHour * 5 $ " minutes")
		Assert(.fn([runinterval: true, interval: 5, interval_units: 3])
			is: "every " $ .minutesInDay * 5 $ " minutes")
		Assert(.fn([runinterval: true, interval: 5, interval_units: 4])
			is: "every " $ .minutesInWeek * 5 $ " minutes")
		}
	}
