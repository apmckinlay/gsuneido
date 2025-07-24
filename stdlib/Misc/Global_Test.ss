// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(Global('PI') is: PI)
		Assert(Global('TRACE.QUERIES') is: TRACE.QUERIES)

		err = "can't find"
		Assert({ Global('Foo.') } throws: err)
		Assert({ Global('.Foo') } throws: err)
		Assert({ Global('Foo.123') } throws: err)
		}
	}