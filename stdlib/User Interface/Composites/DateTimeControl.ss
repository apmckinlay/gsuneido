// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: 'DateTime'

	New(datetime = 0)
		{
		super(.layout())
		.shortcut = .FindControl('dateshortcuts')
		.monthCal = .FindControl('MonthCal')
		.ctrls = .FindControl('controls')
		if (Date?(datetime))
			.Set(datetime)
		else
			.Set(Date().NoTime().Replace(hour: 12))
		}

	layout()
		{
		return Object('Record'
			Object('Vert'
				Object('Horz'
					Object('Vert'
						Object('MonthCal'),
						Object('dateshortcuts')),
					.hourList(),
					.minuteList(),
					name: 'controls')
			 ))
		}

	hourList()
		{
		list = Settings.Get('TimeFormat').Suffix?('tt')
			? #('12 am', '1 am', '2 am', '3 am', '4 am', '5 am', '6 am',
				'7 am', '8 am', '9 am','10 am','11 am','12 pm','1 pm','2 pm',
				'3 pm','4 pm','5 pm','6 pm','7 pm', '8 pm', '9 pm', '10 pm',
				'11 pm')
			: #("0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12",
				"13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23")

		return Object('ListBox' list name: 'Hour' xmin: 75 ymin: 160)
		}

	minuteList()
		{
		return #(ListBox #('00', '05', '10', '15', '20', '25', '30', '35', '40',
			'45', '50', '55')
			name: 'Min' xmin: 50 ymin: 160)
		}

	Record_NewValue(name, value)
		{
		if name is 'dateshortcuts'
			{
			date = MonthCalDialog.ConvertShortcut(value, showTime: true)
			if date isnt ''
				{
				.Data.SetField('monthCal', date)
				.monthCal.Set(date)
				}
			}

		if name is 'MonthCal'
			{
			.Data.SetField('dateshortcuts', '')
			.shortcut.Set('')
			}
		}

	OK()
		{
		return .Get()
		}
	ListBoxDoubleClick(sel/*unused*/)
		{
		.Send('On_OK')
		}

	Get()
		{
		date = .monthCal.Get().Format('MMM dd yyyy')
		hour = .ctrls.Hour.GetText(.ctrls.Hour.GetSelected()).Split(' ')
		if (hour is #() or date is '')
			return false

		suffix = hour.Size() > 1 ? suffix = hour[1] : ''
		return Date(Display(date) $ ' ' $ hour[0] $ ':' $
			.ctrls.Min.GetText(.ctrls.Min.GetSelected()) $ suffix)
		}
	round_to_five(minutes)
		{
		minute = (Number(minutes) / 5).Round(0) * 5
		if minute is 60
			minute = 0
		return minute.Pad(2)
		}
	Set(datetime)
		{
		date = datetime.NoTime()
		.monthCal.Set(date)
		hour = .hourSelect(datetime.Hour(), Settings.Get('TimeFormat'))
		.ctrls.Hour.SetCurSel(.ctrls.Hour.FindString(hour))
		.ctrls.Min.SetCurSel(.ctrls.Min.FindString(.round_to_five(datetime.Minute())))
		}
	noon: 12
	hourSelect(hour, timeFormat)
		{
		if not timeFormat.Suffix?('tt')
			return String(hour)
		suffix = hour < .noon ? ' am' : ' pm'
		if hour is 0
			hour = 12 // midnight
		else if hour > .noon
			hour = hour - .noon
		return Display(hour) $ suffix
		}
	}
