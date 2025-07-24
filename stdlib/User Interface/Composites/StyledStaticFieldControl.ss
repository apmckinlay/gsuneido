// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
HorzControl
	{
	New()
		{
		.Send('Data') // so we know about record changes
		}

	Set(value)
		{
		.RemoveAll()
		if value isnt ''
			.Insert(0, Object('StyledStatic', value))
		}

	Destroy()
		{
		.Send('NoData')
		super.Destroy()
		}
	}