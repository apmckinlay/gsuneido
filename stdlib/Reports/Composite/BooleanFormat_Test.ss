// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ExportCSV()
		{
		fmt = BooleanFormat()
		Assert(fmt.ExportCSV() is: '"No"')
		Assert(fmt.ExportCSV(true) is: '"Yes"')
		Assert(fmt.ExportCSV(false) is: '"No"')
		Assert(fmt.ExportCSV('') is: '"No"')
		Assert(fmt.ExportCSV([]) is: '"[]"')
		Assert(fmt.ExportCSV('something invalid') is: '"something invalid"')

		fmt = BooleanFormat(true)
		Assert(fmt.ExportCSV(true) is: '"Yes"')
		Assert(fmt.ExportCSV([]) is: '"Yes"')
		Assert(fmt.ExportCSV('') is: '"Yes"')
		Assert(fmt.ExportCSV(false) is: '"Yes"')
		Assert(fmt.ExportCSV('something invalid') is: '"Yes"')


		fmt = BooleanFormat(false)
		Assert(fmt.ExportCSV(true) is: '"No"')
		Assert(fmt.ExportCSV([]) is: '"No"')
		Assert(fmt.ExportCSV('') is: '"No"')
		Assert(fmt.ExportCSV(false) is: '"No"')
		Assert(fmt.ExportCSV('something invalid') is: '"No"')
		}
	}