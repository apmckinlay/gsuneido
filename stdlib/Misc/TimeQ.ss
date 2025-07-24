// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(t)
		{
		if false is t = .convertToNumber(t)
			return false
		maxValidTime = 2400
		if t < 0 or t > maxValidTime or not Integer?(t)
			return false
		divisor = 100 // shift two decimal places to separate hours from minutes
		hour = (t / divisor).Int()
		min  = t % divisor
		maxHour = 24
		maxMinute = 60
		return 0 <= hour and hour < maxHour and 0 <= min and min < maxMinute
		}

	convertToNumber(t)
		{
		if t is ''
			return false
		try
			t = Number(t)
		catch
			t = false
		return t
		}
	}