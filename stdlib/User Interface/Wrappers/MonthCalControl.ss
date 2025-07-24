// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	Name:		"MonthCal"
	Xstretch:	1
	Ystretch:	1

	New(date = "", daystate = false)
		{
		.CreateWindow("SysMonthCal32", "", WS.VISIBLE | MCS.NOTODAY |
			(daystate is true ? MCS.DAYSTATE : 0))
		.SubClass()

		SendMessageRect(.Hwnd, MCM.GETMINREQRECT, 0, rc = Object())
		padding = 4
		.Xmin = rc.right + padding
		.Ymin = rc.bottom + padding

		.Send("Data")
		.Map = Object()
		.Map[MCN.SELCHANGE] = 'MCN_SELCHANGE'
		.Map[MCN.SELECT] = 'MCN_SELECT'
		.Map[MCN.GETDAYSTATE] = 'MCN_GETDAYSTATE'

		// set the current date selection if a date is passed in
		if Date?(date)
			.Set(date)

		.SetRange(Date.Begin(), Date.End())
		}
	GETDLGCODE()
		{ return DLGC.WANTALLKEYS }
	CHAR(wParam)
		{
		switch (wParam)
			{
		case VK.SPACE :
			.MCN_SELECT()
			return 0
		case VK.RETURN :
			.Send("On_OK")
			return 0
		case VK.ESCAPE :
			.Send("On_Cancel")
			return 0
		default:
			return 'callsuper'
			}
		}
	Set(date)
		{
		if not Date?(date)
			return
		d = Object(
			wYear: date.Year(),
			wMonth: date.Month(),
			wDay: date.Day())
		SendMessageSystemTime(.Hwnd, MCM.SETCURSEL, 0, d)
		}
	Get()
		{
		SendMessageSystemTime(.Hwnd, MCM.GETCURSEL, 0, d = Object())
		date = Date(year: d.wYear, month: d.wMonth, day: d.wDay)
		return date is false ? false : date.NoTime()
		}
	MCN_SELECT()
		{
		d = .Get()
		.Send("NewValue", d)
		.Send("DateSelectChange", d)
		return 0
		}
	MCN_SELCHANGE()
		{
		.Send("MonthCalSelChange", .Get())
		return 0
		}
	MCN_GETDAYSTATE(lParam)
		{
		StructModify(NMDAYSTATE, lParam)
			{
			d = it.stStart
			cDayState = it.cDayState
			it.cDayState = 0
			}
		// GetDayState should return a list of numbers (one number per month)
		states = .Send("GetDayState",
			Date(year: d.wYear, month: d.wMonth, day: d.wDay), cDayState)
		if Object?(states)
			.Defer({ .SetDayState(states) }, uniqueID: 'SetDayState')
		return true
		}
	SetDayState(states)
		{
		n = SendMessageText(.Hwnd, MCM.GETMONTHRANGE, GMR_DAYSTATE,
			'\x00'.Repeat(2 * SYSTEMTIME.Size()))
		s = states.Map({ LONG(Object(x: it)) }).Join().RightFill(n * 4, '\x00')
		SendMessageTextIn(.Hwnd, MCM.SETDAYSTATE, n, s)
		}
	SetRange(min, max)
		{
		range = Object(
			min: Object(wYear: min.Year(), wMonth: min.Month(), wDay: min.Day())
			max: Object(wYear: max.Year(), wMonth: max.Month(), wDay: max.Day()))
		SendMessageSTRange(.Hwnd, MCM.SETRANGE, GDTR.MIN | GDTR.MAX, range)
		}
	Destroy()
		{
		.Send("NoData")
		super.Destroy()
		}
	}
