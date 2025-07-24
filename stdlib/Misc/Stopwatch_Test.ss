// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		_date = Date()
		sw = new Stopwatch { Stopwatch_date() { return _date } }
		_date = _date.Plus(milliseconds: 123)
		Assert(sw() is: "123 ms")
		_date = _date.Plus(milliseconds: 456)
		Assert(sw() is: "+ 456 ms = 579 ms")
		_date = _date.Plus(seconds: 10)
		Assert(sw() is: "+ 10 sec = 10.6 sec")
		}
	}