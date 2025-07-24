// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_good()
		{
		error = Object(failed: false)
		a = Object(a: false)
		mutex = Mutex()
		wg = WaitGroup()
		f = {
			for ..1000
				{
				mutex.Do()
					{
					if a.a
						error.failed = true
					a.a = true
					a.a = false
					}
				}
			}
		wg.Thread(f)
		wg.Thread(f)
		wg.Wait()
		Assert(not error.failed)
		}
	XTest_bad() // disabled because it fails on some systems
		{
		if Suneido.GoMetric("/sched/gomaxprocs:threads") < 4
			return
		for i in ..6
			if .test_bad()
				return
			else
				Thread.Sleep(8 << i)
		throw "failed"
		}
	test_bad() // returns true on success
		{
		error = Object(failed: false)
		a = Object(a: false)
		wg = WaitGroup()
		f = {
			for ..1000
				{
				if a.a
					error.failed = true
				a.a = true
				a.a = false
				}
			}
		wg.Thread(f)
		wg.Thread(f)
		wg.Wait()
		// no mutex so it should fail
		return error.failed
		}
	}