// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
ChooseField
	{
	Name: ChooseYearMonth
	New(mandatory = false)
		{
		super(Object('Field', width: 5,
			status: "a four digit year and two digit month e.g. 200312 for Dec. 2003"),
			:mandatory)
		}

	Getter_DialogControl()
		{
		return Object('YearMonth', .Field.Get())
		}

	Valid?()
		{
		if '' is s = .Get()
			return super.Valid?()
		// could possibly remove the size checking since Date constructor no longer
		// accepts things like "2009040101"
		validSize = 6
		return s.Size() is validSize and false isnt Date('#' $ s $ '01')
		}
	}