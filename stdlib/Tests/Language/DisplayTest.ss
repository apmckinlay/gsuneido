// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(Display(Abs) is: 'Abs /* stdlib function */')

		Assert(Display(Dates) is: 'Dates /* stdlib class */')

		Assert(Display(Dates.Time) is: 'Dates.Time /* stdlib method */')
		}
	}