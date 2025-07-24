// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		err = "invalid index column"
		Assert({ .MakeTable("(a,b,C) key(c)") } throws: err)
		Assert({ .MakeTable("(a,b,C) key(a) index(c)") } throws: err)
		}
	}
