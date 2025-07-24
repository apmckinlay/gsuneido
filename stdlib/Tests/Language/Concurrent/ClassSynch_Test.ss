// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_good()
		{
		.a = false
		.Error = false
		wg = WaitGroup()
		f = {
			for ..1000
				.good()
			}
		wg.Thread(f)
		wg.Thread(f)
		wg.Wait()
		Assert(not .Error)
		}
	good()
		{
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
		f = {
			for ..1000
				.bad()
			}
		wg.Thread(f)
		wg.Thread(f)
		wg.Wait()
		Assert(.Error)
		}
	bad()
		{
		if .a
			.Error = true
		.a = true
		.a = false
		}
	}