// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_parse()
		{
		Assert(SchedEvery("garbage") is: false)
		Assert(SchedEvery("every 0 minutes") is: false)

		test = function (s, n)
			{
			sched = SchedEvery(s)
			Assert(sched.SchedEvery_every is: n)
			}
		test("every 1 minute", 1)
		test("every 10 minutes", 10)
		}
	Test_due()
		{
		test = function (prevcheck, cur, expected)
			{
			Assert(SchedEvery('every 5 minutes').Due?(prevcheck, cur), is: expected)
			}
		test(false, #20160101.0914, false)
		test(#20160101.0914, #20160101.0915, true)
		test(#20160101.0915, #20160101.0918, false) // not due again
		test(#20160101.0919, #20160101.0920, true) // due again
		test(#20160101, #20160201.0005, true)

		// 1 day interval
		test = function (prevcheck, cur, expected)
			{
			Assert(SchedEvery('every 1440 minutes').Due?(prevcheck, cur), is: expected)
			}

		// no prevcheck (first run), running at 1PM
		test(false, #20160101.1300, false)
		// no prevcheck (first run), running at midnight
		test(false, #20160102, true)
		// no prevcheck (first run), 5 minutes after midnight
		test(false, #20160101.000500000, false)
		// no prevcheck (first run), 50 seconds after midnight
		test(false, #20160102.000050000, true)
		// prev chek at noon, cur check at 1PM
		test(#20160101.1200, #20160101.1300, false)
		// prev check at noon yesterday, cur check at 5 minutes after midnight
		test(#20160101.1200, #20160102.0005, true)
		// prev check at noon yesterday, cur check at 50 seconds after midnight
		test(#20160101.1200, #20160102.000050000, true)
		}
	}
