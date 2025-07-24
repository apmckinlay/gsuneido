// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	New(dates = "", protectBefore = "")
		{
		super(Object('Vert'
			Object('MonthCalDates', :dates, :protectBefore, name: 'dates')
			))
		}
	OK()
		{
		dates = .FindControl('dates').Get()
		return dates.Sort!().Join(',')
		}
	}
