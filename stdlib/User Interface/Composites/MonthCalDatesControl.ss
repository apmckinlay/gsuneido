// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	New(dates = "", protectBefore = "")
		{
		super(#(Vert
			(MonthCal daystate:, xstretch: 1, ystretch: 1)
			(Field xstretch: 1, readonly:)
			))
		.dates = dates.Split(',')
		.protectBefore = protectBefore is "" ? Date.Begin() : protectBefore
		for i in .dates.Members()
			.dates[i] = Date(.dates[i])
		.monthcal = .Vert.MonthCal
		.Xmin += .monthcal.Xmin * 2 + 10
		.Ymin += .monthcal.Ymin * 2 + 25
		.field = .Vert.Field
		.startdate = Date().Plus(months: (.nMonths / -2).Round(0))
		.Set(.dates)
		}
	DateSelectChange(date)
		{
		if date < .protectBefore
			{
			Alert("You cannot change dates prior to " $ .protectBefore.ShortDate())
			return
			}
		if .dates.Has?(date)
			.dates.Remove(date)
		else
			.dates.Add(date)
		.field.Set(.dates.Sort!().Map({ it.ShortDate() }).Join(','))
		.monthcal.SetFocus()
		.monthcal.SetDayState(.GetDayState(.startdate, .nMonths))
		}
	boldday(states, startdate, date)
		{
		month = date.MinusMonths(startdate)
		if (states.Member?(month))
			states[month] |= (1 << (date.Day() - 1))
		}
	GetDayState(startdate, nMonths)
		{
		.startdate = startdate
		.nMonths = nMonths
		states = Object()
		for (i = 0; i < nMonths; ++i)
			states[i] = 0
		for date in .dates
			.boldday(states, startdate, date)
		return states
		}

	nMonths: 11 // apm - where does "11" come from ???

	Get()
		{
		return .dates
		}

	Set(dates)
		{
		.dates = dates
		formatted_dates = ''
		for date in .dates
			formatted_dates $= Date(date).ShortDate() $ ","
		.field.Set(formatted_dates[.. -1])
		}
	}
