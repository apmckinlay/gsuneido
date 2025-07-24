// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
class
	{
	rx: '^every (\d+) minutes?$'
	CallClass(when)
		{
		every = when.Extract(.rx)
		return every is false or every is '0' ? false : new this(Number(every))
		}
	New(.every)
		{
		}
	Due?(prevcheck, curtime)
		{
		for (t = prevcheck is false ? curtime : prevcheck.Plus(minutes: 1);
			t <= curtime; t = t.Plus(minutes: 1))
			{
			minute = t.MinusDays(#20000101) * 24 * 60 + t.Hour() * 60 + t.Minute()
			if minute % .every is 0
				return true
			}
		return false
		}
	}
