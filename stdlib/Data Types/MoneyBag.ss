// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
class
	{
	UseDeepEquals: true
	New(@args)
		{
		.bag = Object().Set_default(0)
		for cur, amt in args
			.Plus(amt, cur)
		}
	Plus(amount, currency)
		{
		.bag[currency] += amount
		return this
		}
	PlusMB(other)
		{
		amounts = other.Amounts()
		for cur in amounts.Members()
			.Plus(amounts[cur], cur)
		return this
		}
	Minus(amount, currency)
		{
		.bag[currency] -= amount
		return this
		}
	MinusMB(other)
		{
		amounts = other.Amounts()
		for cur in amounts.Members()
			.Minus(amounts[cur], cur)
		return this
		}
	Amounts()
		// returns an object with the totals for each currency
		// the object will have a member for each currency that has had a value
		// even if it is now zero
		// e.g. Object(CAD: 0, USD: 456.78)
		{
		return .bag.Copy()
		}
	Currencies()
		{
		return .bag.Members().Instantiate()
		}

	ToString()
		{
		return 'MoneyBag(' $ .bag.Map2({|cur,amt| cur $ ': ' $ amt }).Join(', ') $ ')'
		}
	Display(format = '-###,###,###.##')
		{
		return .bag.Map2({|cur,amt| amt.Format(format) $ ' ' $ cur }).Join(', ')
		}

	Zero?()
		{
		return .bag.Every?({ it is 0 })
		}

	RemoveCurrency(currency)
		{
		.bag.Delete(currency)
		return this
		}
	}
