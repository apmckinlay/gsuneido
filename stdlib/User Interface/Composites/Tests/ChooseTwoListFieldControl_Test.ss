// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_valid()
		{
		validList = #(Fred, Barney, Wilma)

		validfn = ChooseTwoListFieldControl.ChooseTwoListFieldControl_validCheck?
		Assert(validfn(#(), validList))
		Assert(validfn(#(Betty), validList) is: false)
		Assert(validfn(#(Fred), validList))
		Assert(validfn(#(Fred, Fred), validList) is: false)
		}
	}