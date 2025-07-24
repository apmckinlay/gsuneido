// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(BaseClass(123) is: false)
		Assert(BaseClass(#()) is: false)
		Assert(BaseClass([]) is: false)
		Assert(BaseClass(class{}) is: false)
		Assert(BaseClass(FieldControl) is: EditControl)
		Assert(BaseClass(Stack()) is: Stack)
		}
	}