// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	PercentOf(amount, percent, precision = 2)
		{
		amt = amount * Number(percent).PercentToDecimal()
		if precision isnt false
			amt = amt.Round(precision)
		return amt
		}

	ConvertToDollars(amount)
		{
		return Number(amount) * .01 /*= convert to dollar and cents */
		}
	}