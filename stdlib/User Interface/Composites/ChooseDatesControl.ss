// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
ChooseField
	{
	Name: ChooseDates
	dates: ""
	New(mandatory = false, .protectBeforeField = "")
		{
		super(#('Field', readonly:), :mandatory)
		}
	Getter_DialogControl()
		{
		return Object(MonthCalDatesDialog, .Get(),
			protectBefore: .Send('GetField', .protectBeforeField))
		}
	Get()
		{
		return .dates is "" or .dates.Empty?() ? '' : .dates.Join(',')
		}
	Set(val)
		{
		.dates = Object()
		format_date = ''
		for date in val.Split(',')
			{
			.dates.Add(date)
			format_date $= Date(date).ShortDate() $ ","
			}
		.Field.Set(format_date[.. -1])
		}
	Valid?()
		{
		return .Field.Get().Split(',').Every?({ Date(it) isnt false })
		}
	GetReadOnly()
		{
		return false
		}
	}