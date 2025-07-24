// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		Assert(Internal?(#internal)) // Field_internal
		Assert(not Internal?(#date)) // Field_date
		Assert(not Internal?(#nonexistent))
		}
	}