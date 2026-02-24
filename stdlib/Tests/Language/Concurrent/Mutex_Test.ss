// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_good()
		{
		.a = false
		.failed = false
		.mutex = Mutex()
		wg = WaitGroup()
		wg.Thread(.f)
		wg.Thread(.f)
		wg.Wait()
		Assert(not .failed)
		}
	f()
		{
		for ..1000
			{
			.mutex.Do()
				{
				if .a
					.failed = true
				.a = true
				.a = false
				}
			}
		}
	}