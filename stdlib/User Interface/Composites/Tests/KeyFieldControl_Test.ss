// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_MandatoryAndAllowOther()
		{
		Assert(KeyFieldControl.ValidData?("", "field", "query"))
		Assert(KeyFieldControl.ValidData?("", "field", "query", mandatory:) is: false)
		Assert(KeyFieldControl.
			ValidData?("test", "field", "query", allowOther:, mandatory:))
		}
	}