// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// for dialogs, used by ChooseDateControl
Controller
	{
	New(date = "")
		{
		super(.layout(date))
		.shortcut = .FindControl('dateshortcuts')
		}

	DateSelectChange(date)
		{
		.Window.Result(date)
		}

	layout(date)
		{
		return Object('Vert'
			Object('MonthCal', date),
			Object('dateshortcuts'))
		}

	ConvertShortcut(value, showTime)
		{
		prefixes = Object('Today': 't', 'Start of Current Month': 'm',
			'End of Current Month': 'h', 'Start of Previous Month': 'pm',
			'End of Previous Month': 'ph', 'Start of Current Year': 'y'
			'End of Current Year': 'r')
		if not prefixes.Member?(value)
			return ''

		prefix = prefixes[value]
		return DateControl.ConvertToDate(prefix, convertDateCodes?:, :showTime)
		}

	NewValue(value, source = false)
		{
		if source is false or source.Name isnt 'dateshortcuts'
			return
		date = .ConvertShortcut(value, showTime: false)
		if date isnt ''
			.Window.Result(date)
		}

	On_OK()
		{
		.Window.Result(.FindControl('MonthCal').Get())
		}
	}
