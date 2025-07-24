// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
class
	{
	rx: '^on (Sun|Mon|Tue|Wed|Thu|Fri|Sat) '
	CallClass(when)
		{
		return (false is on = when.Extract(.rx)) or
			(false is at = SchedAt(when.AfterFirst(on $ ' ')))
			? false : new this(on, at)
		}
	New(.on, .at)
		{
		}
	Due?(prevcheck, curtime)
		{
		return curtime.Format('ddd') is .on and .at.Due?(prevcheck, curtime)
		}
	}
