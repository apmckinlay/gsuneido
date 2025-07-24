// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
// SuJsWebTest
Test
	{
	Test_one()
		{
		Assert(Instance?(123) is: false)
		Assert(Instance?(#()) is: false)
		Assert(Instance?(class{}) is: false)
		Assert(Instance?(Stack) is: false)

		Assert(Instance?(class{}()))
		Assert(Instance?(Stack()))
		}
	}