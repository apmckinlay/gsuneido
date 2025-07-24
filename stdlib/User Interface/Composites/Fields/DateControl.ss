// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
FieldControl
	{
	Name: "Date"
	ComponentName: "Date"
	convertDateCodes?: true
	New(status = "A date e.g. May 29 or May 29, 2003 or 5/29/03",
		.readonly = false, .showTime = false, .mandatory = false, style = 0,
		set = false, tabover = false, hidden = false)
		{
		super(:status, :readonly, :style, :tabover, :hidden)
		if 0 isnt result = .Send('DateControl_ConvertDateCodes', showTime: .showTime)
			.convertDateCodes? = result
		if set isnt false
			.Defer({ .Set(set) })
		.ComponentArgs = Object(readonly, showTime, tabover, hidden)
		}
	SetFontAndSize(font, size, weight, underline, width/*unused*/, height/*unused*/)
		{ // overrides EditControl
		// widen year to match what FormatValue does
		fmt = Settings.Get('ShortDateFormat').Replace('\<yy?y?\>', 'yyyy')
		if .showTime is true
			fmt $= ' ' $ Settings.Get('TimeFormat')
		// this date/time must have 2 digit year, month, day, hour, minute, second
		// can't just use format because e.g. "M" may expand to "12"
		s = #19991231.235959.Format(fmt)
		super.SetFontAndSize(font, size, weight, underline, 1, 1, text: s)
		}

	Valid?()
		{
		return .validCheck?(.Get(), .readonly, .mandatory)
		}

	validCheck?(data, readonly, mandatory)
		{
		if readonly
			return true
		if not .convertDateCodes? and String?(data) and data isnt ''
			return .validDateCode?(data, .showTime)
		return Date?(data) or (data is '' and not mandatory)
		}

	validDateCode?(str, showTime = false)
		{
		if not #(t, m, h, y, r).Has?(str[0].Lower()) and
			not #(pm, ph).Has?(str[..2].Lower())
			return false

		if str[0].Lower() is 'p'
			str = str[1..]

		time = ''
		if showTime and ('' isnt time = .extractTime(str)) and not Time?(time)
			return false
		return str[1 .. (str.Size() - time.Size())].Tr('+-') is ''
		}

	extractTime(text)
		{
		return text.Extract('\d?\d?\d?\d?$')
		}

	ValidData?(@args)
		{
		value = args[0]
		return .validCheck?(value,
			args.GetDefault('readonly', false), args.GetDefault('mandatory', false))
		}

	ValidateRange?()
		{
		return .convertDateCodes?
		}

	KillFocus()
		{
		date = .get()
		if not super.Dirty?() or date is false
			return
		if super.Get() is s = .FormatValue(date, .showTime)
			return
		SetWindowText(.Hwnd, s)
		.Dirty?(true)
		.SelectAll()
		}

	Get()
		{
		// sometimes we use milliseconds for "hidden" data
		// don't want to lose this (e.g. tabbing through Browse)
		if not .Dirty?() and Date?(.orig_value)
			return .orig_value
		date = .orig_value = .get()
		if date is false
			return super.Get()
		if .showTime isnt true and Date?(date) // may be a date code (t,m,h,y,r)
			date = .orig_value = date.NoTime()
		return date
		}
	get()
		{
		text = super.Get()
		return .ConvertToDate(text, .convertDateCodes?, .showTime)
		}

	ConvertToDate(text, convertDateCodes?, showTime = false, format = false)
		{
		try // Try is to handle invalid formats: IE: '', 'a', etc, Ref: FromulaDate
			{
			// If text is a date in string form, convert and return it accordingly
			date = Date(text, format is false ? Settings.Get('ShortDateFormat') : format)
			if date isnt false
				return showTime ? date : date.NoTime()
			}

		// else, if text was not a date, ensure it is a valid date shortcut
		// IE: t++, m-100, etc
		if not .validDateCode?(text, showTime)
			return false

		// If text is a valid shortcut, return the date OR return the date shortcut
		// Places we save the shortcut: FormulaAddFunction, ScheduledReportParamsControl
		return not convertDateCodes?
			? text
			: .convertDate(text, showTime)
		}

	convertMap: (
		t:  function (today) { today 												},
		m:  function (today) { today.Replace(day: 1) 								},
		h:  function (today) { today.Replace(day: 1).Plus(months: 1).Plus(days: -1) },
		y:  function (today) { today.Replace(month: 1, day: 1) 						},
		r:  function (today) { today.Replace(month: 12, day: 31) 					},
		pm: function (today) { today.Replace(day: 1).Plus(months: -1) 				},
		ph: function (today) { today.Replace(day: 1).Plus(days: -1)					})
	convertDate(text, showTime)
		{
		days = text.Count('+') - text.Count('-')
		shortcut = text[0].Lower() is 'p' ? text[::2].Lower() : text[0].Lower()
		date = (.convertMap[shortcut])(.today()).Plus(:days)
		return showTime ? .adjustTime(text, date) : date.NoTime()
		}

	// for tests
	// dynamic parameter is for tests where this can't easily be overridden
	// (calls are not directly from the test)
	today(_dateForTest = false)
		{
		if dateForTest isnt false
			return dateForTest
		return Date()
		}

	adjustTime(text, date)
		{
		if '' is timeStr = .extractTime(text)
			return date
		timeStr = timeStr.LeftFill(4 /*= time format length*/, '0')
		// Date treats this as Date.NoTime, but user specified timeStr so ensure we set to midnight
		return timeStr is '0000'
			? date.NoTime().Plus(milliseconds: 1)
			: Date(date.Format('yyyyMMdd') $ '.' $ timeStr, 'yMd')
		}

	orig_value: false
	Set(value)
		{
		.orig_value = value
		super.Set(Date?(value)
			? .FormatValue(value, showTime: .showTime, fromSet:)
			: value)
		}

	FormatValue(date, showTime = false, fromSet = false)
		{
		if not .convertDateCodes? and String?(date)
			return date
		s = date.ShortDate()
		fmt = Settings.Get('ShortDateFormat')
		if Date(s, fmt) isnt date.NoTime()
			{
			// convert y, yy, yyy to yyyy
			s = date.Format(fmt.Replace('\<yy?y?\>', 'yyyy'))
			if Date(s, fmt) isnt date.NoTime()
				SuneidoLog("ERROR: irreversable date", params: [:s, :fmt])
			}
		if showTime is true
			s = .AddTimeFormat(date, s, fromSet)
		return s
		}

	AddTimeFormat(date, s, fromSet)
		{
		// Only adjust the time if it is not a static method call and it is not from Set
		// because otherwise the date value should be exactly what is wanted
		if .Member?(#Hwnd) and not fromSet
			{
			// if user types a date without time,
			// default the time to 12 PM (same as choosing from calendar)
			if date is date.NoTime() and not .timeEntered?()
				date = date.Replace(hour: 12)
			}
		return s $ ' ' $ date.Time()
		}

	timeEntered?()
		{
		text = GetWindowText(.Hwnd)
		if text.Has?(':') or text.Lower().Has?('am') or text.Lower().Has?('pm')
			return true
		return false
		}
	}
