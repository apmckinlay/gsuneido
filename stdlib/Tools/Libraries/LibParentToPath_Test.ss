// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		x = Query1(#stdlib, group: -1, name: #LibParentToPath)
		Assert(LibParentToPath(#stdlib, x.parent) is: "/Tools/Libraries")
		}
	}