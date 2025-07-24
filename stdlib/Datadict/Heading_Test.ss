// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(Heading('uom') is: 'UOM') // has both Heading and Prompt
		Assert(Heading('date') is: 'Date') // has Heading but not Prompt
		Assert(Heading('max_date') is: 'Max Date')
		}
	}