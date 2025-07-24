// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (amounts, prorate_amount, round = 2)
	{
	if amounts.Every?({ it is 0 })
		amounts = amounts.Copy().Map!({|unused| 1 }) // Map! to handle named members
	prorate_remaining = prorate_amount.Abs()
	tot_remaining = amounts.SumWith(#Abs)
	prorate_amounts = Object().Set_default(0)
	// need sort for repeatable/consistent results
	for i in amounts.Members().Sort!()
		{
		amount = amounts[i].Abs()
		prorate = (prorate_remaining * (amount / tot_remaining)).Round(round)
		prorate_amounts[i] = prorate_amount.Sign() * prorate
		prorate_remaining -= prorate
		tot_remaining -= amount
		}
	return prorate_amounts
	}