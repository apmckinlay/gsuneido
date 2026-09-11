// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
ChooseField
	{
	Name: #ChooseDates
	New(mandatory = false, .protectBeforeField = "")
		{
		super(#(Field, readonly:), :mandatory)
		}

	Getter_DialogControl()
		{
		return [MonthCalDatesDialog, .Get(),
			protectBefore: .Send(#GetField, .protectBeforeField)]
		}

	getter_dates()
		{
		return .dates = Object()
		}

	Get()
		{
		return .dates.Join(',')
		}

	Set(val)
		{
		.dates = Object()
		formattedDates = Object()
		for date in val.Split(',')
			{
			.dates.Add(date)
			Assert(Date?(d = Date(date)))
			formattedDates.Add(d.ShortDate())
			}
		.Field.Set(formattedDates.Join(','))
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
