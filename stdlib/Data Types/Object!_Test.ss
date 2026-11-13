// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		fn = { it }
		Assert(Object!(fn, []) is: [])
		Assert(Object!(fn, #()) is: #())
		Assert(Object!([]) is: [])

		Assert({ Object!(fn, true) }
			throws: "expected an object but it was Boolean")
		Assert({ Object!(fn, "") } throws: "expected an object but it was String")
		Assert(
			{
			Object!(fn,
				class
					{
					})
			} throws: "expected an object but it was Class")
		Assert(
			{
			Object!(fn,
				class
					{
					}())
			} throws: "expected an object but it was Instance")
		}
	}
