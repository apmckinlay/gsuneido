// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
// TODO remove StartMonth and MidMonth - can just use 1 or 15
class
	{
	// usage examples:
	// "on StartMonth at 14:30"
	// "on MidMonth at 09:00"
	// "on EndMonth at 15:00"
	// "on 5 at 06:00"  - that would be 5th day of the month (must be between 2 and 28)
	rx: '^on (StartMonth|MidMonth|EndMonth|\d\d?) '
	CallClass(when)
		{
		if false is onStr = when.Extract(.rx)
			return false
		return (false is on = .schedOn(onStr)) or
			(false is at = SchedAt(when.AfterFirst(onStr $ ' ')))
			? false : new this(on, at)
		}
	New(.on, .at)
		{
		}
	schedOn(on)
		{
		if on.Numeric?()
			{
			n = Number(on)
			return n < 2 or n > 28 ? false : n
			}
		switch (on)
			{
		case 'StartMonth':
			return 1
		case 'MidMonth':
			return 15
		case 'EndMonth':
			return Date().EndOfMonthDay()
		default :
			return false
			}
		}
	Due?(prevcheck, curtime)
		{
		return curtime.Day() is .on and .at.Due?(prevcheck, curtime)
		}
	}
