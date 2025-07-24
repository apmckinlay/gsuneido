// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Format(format)
		{
		s = .FormatEn(format)
		if GetLanguage().name isnt 'english'
			s = .translate(s)
		return s
		}
	translate(s)
		{
		pat = "\<(Monday|Tuesday|Wednesday|Thursday|Friday|Saturday|Sunday|" $
			"Mon|Tue|Wed|Thu|Fri|Sat|Sun|" $
			"January|February|March|April|May|June|" $
				"July|August|September|October|November|December|" $
			"Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)\>"
		return s.Replace(pat, TranslateLanguage)
		}
	ShortDate()
		{
		return .Format(Settings.Get('ShortDateFormat'))
		}
	LongDate()
		{
		return .Format(Settings.Get('LongDateFormat'))
		}
	Time()
		{
		return .Format(Settings.Get('TimeFormat'))
		}
	DayOfYear()
		{
		return .MinusDays(Date(year: .Year(), month: 1, day: 1)) + 1
		}
	NoTime()
		{
		return Date(year: .Year(), month: .Month(), day: .Day(),
			hour: 0, minute: 0, second: 0, millisecond: 0)
		}
	NoTime?()
		{
		return .Hour() is 0 and .Minute() is 0 and .Second() is 0 and .Millisecond() is 0
		}
	StartOfDay() // another name for NoTime
		{
		return Date(year: .Year(), month: .Month(), day: .Day(),
			hour: 0, minute: 0, second: 0, millisecond: 0)
		}
	EndOfDay()
		{
		return this is Date.End() ? this :
			Date(year: .Year(), month: .Month(), day: .Day(),
				hour: 23, minute: 59, second: 59, millisecond: 999)
		}
	Replace(@args)
		{
		ob = Object(year: .Year(), month: .Month(), day: .Day(), hour: .Hour(),
			minute: .Minute(), second: .Second(), millisecond: .Millisecond())
		for m in args.Members()
			if (ob.Member?(m))
				ob[m] = args[m]
			else
				throw "date.Replace: invalid argument: " $ m
		return Date(@ob)
		}
	MinusMinutes(date)
		{
		return this.MinusSeconds(date) / 60
		}
	MinusHours(date)
		{
		return this.MinusSeconds(date) / (60 * 60)
		}
	MinusMonths(date)
		{
		return 12 * (.Year() - date.Year()) + (.Month() - date.Month())
		}
	ShortDateTime()
		{
		return .ShortDate() $ " " $ .Time()
		}
	ShortDateTimeSec()
		{
		return .Format(Settings.Get('ShortDateFormat') $ ' ' $
			Settings.Get('TimeFormat') $ ":ss")
		}
	StdShortDate()
		{
		return .Format('yyyy-MM-dd')
		}
	StdShortDateTime()
		{
		return .Format('yyyy-MM-dd HH:mm')
		}
	StdShortDateTimeSec()
		{
		return .Format('yyyy-MM-dd HH:mm:ss')
		}
	LongDateTime()
		{
		return .LongDate() $ " " $ .Time()
		}
	StartOfYear()
		{
		return .Replace(month: 1, day: 1).NoTime()
		}
	EndOfMonth()
		{
		return .Replace(day: 1).Plus(months: 1, days: -1)
		}
	EndOfMonthDay()
		{
		return .EndOfMonth().Day()
		}
	WeekNumber(firstday = 'sun')
		{
		yearWeekStart = .Replace(day: 1, month: 1).WeekStart(firstday)
		return 1 + .WeekStart(firstday).MinusDays(yearWeekStart) / 7 /*= days per week*/
		}
	WeekStart(firstday = 'sun')
		{
		return .Plus(days: -.WeekDay(firstday))
		}
	Quarter()
		{
		return ((.Month() - 1) / 3).Floor() + 1
		}
	IsoWeekDay()
		{
		y = .Year()
		doy = .DayOfYear()
		dow0101 = .Replace(day: 1, month: 1).WeekDay('mon') + 1
		dow = .WeekDay('mon') + 1
		weeknumber = false
		if doy <= (8 - dow0101) and dow0101 > 4
			{
			yearnumber = y - 1
			if dow0101 is 5 or (dow0101 is 6 and .isLeapYear?(y - 1))
				weeknumber = 53
			else
				weeknumber = 52
			}
		else
			yearnumber = y
		if yearnumber is y
			{
			i = .isLeapYear?(y) ? 366 : 365
			if ((i - doy) < (4 - dow))
				{
				yearnumber = y + 1
				weeknumber = 1
				}
			}
		if yearnumber is y
			{
			j = doy + (7 - dow) + (dow0101 - 1)
			weeknumber = j / 7
			if dow0101 > 4
				weeknumber -= 1
			}
		return yearnumber.Pad(4) $ "-" $ "W" $ weeknumber.Pad(2) $ "-" $ dow
		}

	isLeapYear?(year)
		{
		return (year % 4) is 0 and (year % 100) isnt 0 or year % 400 is 0
		}

	GMTime()
		{
		return .Plus(minutes: .GetLocalGMTBias())
		}
	GMTimeToLocal()
		{
		return .Plus(minutes: -.GetLocalGMTBias())
		}

	UTC(gmtBias = false)
		{
		if gmtBias is false
			gmtBias = -.GetLocalGMTBias()
		if gmtBias > 840 /*= GMT+14 */ or gmtBias < -720 /*= GMT-12*/
			return 'Invalid GMT Bias: ' $ gmtBias
		sign = gmtBias is 0
			? '='
			: gmtBias < 0
				? '-'
				: '+'
		hours = String((gmtBias / 60).RoundDown(0).Abs()) 		/*= convert to hours*/
		minutes = String((gmtBias % 60).RoundDown(0).Abs()) 	/*= convert to minutes*/
		return 'UTC' $ sign $ hours.LeftFill(2, '0') $ ':' $ minutes.RightFill(2, '0')
		}

	Minus(@args)
		{
		return .Plus(@args.Map!({ -it }))
		}
	InternetFormat()
		{
		// assumes that the time is local
		date = this
		if date isnt date.NoTime()
			date = date.GMTime()
		return date.Format('ddd, dd MMM yyyy HH:mm:ss') $ ' GMT'
		}
	UnixTime()
		{
		t = this.GMTime()
		days = t.NoTime().MinusDays(#19700101)
		seconds = t.MinusSeconds(t.NoTime()).Round(0)
		return days * 60 * 60 * 24 + seconds
		}
	FromUnix(s)
		{
		hour = 60 * 60 /* = seconds in 1 hour */
		hours = (s / hour).Int()
		milliseconds = (s - hours * hour) * 1000 /* = convert seconds to milliseconds */
		return #19700101.Plus(:hours, :milliseconds).GMTimeToLocal()
		}
	}
