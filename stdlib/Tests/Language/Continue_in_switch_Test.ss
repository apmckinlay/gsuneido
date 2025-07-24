// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		block =
			{
			switch
				{
			case true:
				continue
				}
			}
		Assert({ block() } throws: "block:continue")
		}
	}