// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name: "MonthCal"
	ComponentName: "MonthCal"
	ComponentArgs: #()
	Xstretch: 1
	Ystretch: 1
	New(date = "")
		{
		.Send("Data")
		.Set(date)
		}

	Set(date)
		{
		if not Date?(date)
			date = Date()
		.selectedDate = date.NoTime()
		.Act('Set', date)
		}

	Get()
		{
		return .selectedDate
		}

	SELECT(d)
		{
		.selectedDate = d
		.Send("NewValue", d)
		.Send("DateSelectChange", d)
		return 0
		}


	SetDayState(states)
		{
		.Act('SetDayState', states, .refreshVersion)
		}
	SetRange(min, max)
		{
		.Act('SetRange', min, max)
		}

	refreshVersion: 0
	GETDAYSTATE(startdate, nMonths, .refreshVersion)
		{
		if 0 isnt states = .Send("GetDayState", startdate, nMonths)
			.SetDayState(states)
		}

	Destroy()
		{
		.Send("NoData")
		super.Destroy()
		}
	}