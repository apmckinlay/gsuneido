// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
function (nparts, amount, round = 2)
	{
	amounts = Object().AddMany!(1, nparts)
	return Prorate_Amount(amounts, amount, round)
	}