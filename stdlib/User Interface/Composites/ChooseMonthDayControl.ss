// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
ChooseField
	{
	Name: ChooseMonthDay
	New(mandatory = false)
		{
		super(Object('Field', :mandatory, width: 4,
			status: "a four digit month and day e.g. 1231 for Dec. 31"))
		.mandatory = mandatory
		}
	Getter_DialogControl()
		{ return Object('MonthDay', .Field.Get()) }
	Valid?()
		{
		val = .Get()
		if val is ''
			return not .mandatory
		if val.Size() isnt 4 /* = 2 digits for month and 2 for day */
			return false
		return .DateFromMonthDay(val) isnt false
		}
	FieldKillFocus()
		{
		.SetValid(.Valid?())
		}
	DateFromMonthDay(monthday)
		{
		return Date('#2001' $ monthday)
		}
	}
