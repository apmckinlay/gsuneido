// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_reps()
		{
		_t = Object(time: #20000101)
		timer = Timer
			{
			Timer_date()
				{
				return _t.time
				}
			}
		block = { _t.time = _t.time.Plus(seconds: 2) }
		Assert(timer(secs: 1, :block) is: "2 sec/rep")
		Assert(timer(secs: 31, :block) is: "2 sec/rep")

		block = { _t.time = _t.time.Plus(milliseconds: 20) }
		Assert(timer(secs: 1, :block) is: "50 reps/sec = 20 ms/rep")
		}
	}