// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
StaticControl
	{
	New(@args)
		{
		super(@args)
		if Suneido.User is 'default'
			.SubClass()
		.Send('Data')
		}

	HorzAdjustInListEdit: 0
	DevMenu: ('', 'Inspect Control', 'Copy Field Name', 'Go To Field Definition')

	Destroy()
		{
		.Send('NoData')
		super.Destroy()
		}
	}
