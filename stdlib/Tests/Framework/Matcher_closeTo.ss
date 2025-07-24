// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
// Matches a fractional number if it is 'close to' a benchmark fractional
// number. This is done by rounding both the number and the benchmark number to
// a specified number of decimal places before comparing them.
class
	{
	Match(value, args)
		{
		benchmark     = args[0]
		decimalPlaces = args[1]
		return value.Round(decimalPlaces) is benchmark.Round(decimalPlaces)
		}
	Expected(args)
		{
		return "a value matching " $ args[0] $ " to " $ args[1] $ " decimal places"
		}
	Actual(value, args)
		{
		benchmark     = args[0]
		decimalPlaces = args[1]
		return "was " $ value $ " (at " $ decimalPlaces $
				" decimal places, the difference is " $
			   Abs(value.Round(decimalPlaces) - benchmark.Round(decimalPlaces)) $ ")"
		}
	}