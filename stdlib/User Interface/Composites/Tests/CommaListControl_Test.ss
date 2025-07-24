// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_numeric_valid?()
		{
		value = ''.Split(',')
		Assert(CommaListControl.CommaListControl_all_numeric?(value))

		value = '123'.Split(',')
		Assert(CommaListControl.CommaListControl_all_numeric?(value))

		value = '123, 125'.Split(',')
		Assert(CommaListControl.CommaListControl_all_numeric?(value))

		value = '123, 125test'.Split(',')
		Assert(CommaListControl.CommaListControl_all_numeric?(value) is: false)
		}
	}