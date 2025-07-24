// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
ChooseField
	{
	Name: 'ChooseDate'
	New(mandatory = false, field = 'Date', .checkCurrent = false, .checkValid = false,
		tabover = false, status = "", hidden = false, readonly = false)
		{
		super(.Layout(field), :mandatory, :tabover, :status, :hidden, :readonly)
		}

	Layout(field)
		{
		return Object(field)
		}

	Getter_DialogControl()
		{
		return Object(MonthCalDialog, .Field.Get(), closeButton?:)
		}

	NewValue(value)
		{
		if .checkValid isnt false and '' isnt msg = .checkValidDate(value, .checkValid)
			{
			.Set(.Send('GetField', .Name))
			.Dirty?(false)
			.AlertInfo("Invalid Date", msg)
			return
			}

		super.NewValue(value)
		.Dirty?(false) // to prevent multiple alerts on browse
		if .checkCurrent
			.check_current(value)
		}

	CheckingContributions: 'ExtraChooseDateChecking'
	checkValidDate(date, checkValid, chkReadOnly = true, record = false)
		{
		if chkReadOnly and .GetReadOnly()
			return ''

		if not Date?(date)
			return ''

		for contrib in Contributions(.CheckingContributions)
			if '' isnt msg = contrib(date, checkValid, ctrl: this, :record)
				return msg

		return ''
		}

	check_current(date)
		{
		if .GetReadOnly()
			return
		name = .Name
		field = PromptOrHeading(name)
		if field is name
			field = 'date'
		if .date_out_of_range?(date)
			.AlertInfo("Date Warning",
				"You have entered a " $ field $ " that does not fall within six " $
				"months of the current date.\nPlease check the date and make sure" $
				" it is entered properly.")
		}

	date_out_of_range?(date)
		{
		return (Date?(date) and
			(Date().Plus(months: 6) < date or Date().Plus(months: -6) > date))
		}

	ValidData?(@args)
		{
		value = args[0]

		if args.GetDefault('checkValid', false) isnt false and
			.checkValidDate(value, args.checkValid, chkReadOnly: false,
				record: args.GetDefault('record', false)) isnt ''
			return false

		return DateControl.ValidData?(@args)
		}

	ValidateRange?()
		{
		return .Field.Method?('ValidateRange?') ? .Field.ValidateRange?() : true
		}

	DisplayValues(control /*unused*/, vals)
		{
		return vals.Map(DateControl.FormatValue)
		}
	}
