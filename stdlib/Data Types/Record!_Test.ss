// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		r = [name: #test]
		Assert(Record!(r) is: r)
		Assert(Record!({|arg| arg }, r) is: r)

		Assert({ Record!(true) } throws: "expected a record but it was Boolean")
		Assert({ Record!("") } throws: "expected a record but it was String")
		Assert({ Record!(123) } throws: "expected a record but it was Number")
		Assert({ Record!(#()) } throws: "expected a record but it was Object")
		Assert({ Record!(#20240101) } throws: "expected a record but it was Date")
		}
	}
