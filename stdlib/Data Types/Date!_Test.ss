// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		d = #20240101
		Assert(Date!(d) is: d)
		Assert(Date!({|arg| arg }, d) is: d)

		Assert({ Date!(true) } throws: "expected a date but it was Boolean")
		Assert({ Date!("") } throws: "expected a date but it was String")
		Assert({ Date!(123) } throws: "expected a date but it was Number")
		Assert({ Date!(#()) } throws: "expected a date but it was Object")
		}
	}
