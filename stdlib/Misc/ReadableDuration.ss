// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(secs)
		{
		if secs is 0
			return "0"
		units = #(hr, min, sec, ms, us, ns)
		div = #(3600, 60, 1, .001, .000_001, .000_000_001)
		precision = 3
		for (i = 0; i < div.Size(); ++i)
			if secs > div[i]
				return (secs / div[i]).RoundToPrecision(precision) $ ' ' $ units[i]
		}
	Between(t1, t2)
		{
		return this.CallClass(t2.MinusSeconds(t1))
		}
	}