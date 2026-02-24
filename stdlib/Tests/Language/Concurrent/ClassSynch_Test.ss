// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_good()
		{
		.a = false
		.Error = false
		wg = WaitGroup()
		wg.Thread(.good)
		wg.Thread(.good)
		wg.Wait()
		Assert(not .Error)
		}
	good()
		{
		for ..1000
			.Synchronized()
				{
				if .a
					.Error = true
				.a = true
				.a = false
				}
		}

	XTest_bad() // disabled because it fails on some systems
		{
		.a = false
		.Error = false
		wg = WaitGroup()
		wg.Thread(.bad)
		wg.Thread(.bad)
		wg.Wait()
		Assert(.Error)
		}
	bad()
		{
		for ..1000
			{
			if .a
				.Error = true
			.a = true
			.a = false
			}
		}
	}