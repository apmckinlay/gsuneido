// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		.appended = Object()
		moneybag = new MoneyBag
		MoneyBagTotals(.Append, 'currency', 'GrandTotal', revenue: moneybag)
		Assert(.appended is: #((Vskip)))

		.appended = Object()
		moneybag.Plus(100, 'CAD')
		MoneyBagTotals(.Append, 'currency', 'GrandTotal', revenue: moneybag)
		Assert(.appended, is: #(
			("_output",
				revenue: ("GrandTotal",
					("OptionalNumber", 100, mask: '-###,###,###.##'), skip: 0),
				currency: ("GrandTotal", ("Text", "CAD"), skip: 0)),
			("Vskip")))

		// vskip
		.appended = Object()
		MoneyBagTotals(.Append, 'currency', 'GrandTotal', revenue: moneybag, vskip: false)
		Assert(.appended, is: #(
			("_output",
				revenue: ("GrandTotal",
					("OptionalNumber", 100, mask: '-###,###,###.##'), skip: 0),
				currency: ("GrandTotal", ("Text", "CAD"), skip: 0))))

		// additional totals
		.appended = Object()
		MoneyBagTotals(.Append, 'currency', 'GrandTotal',
			revenue: moneybag
			additional_totals: #((expenses: (Text 'Expenses')))
			)
		Assert(.appended, is: #(
			("_output",
				expenses: (Text 'Expenses')
				revenue: ("GrandTotal",
					("OptionalNumber", 100, mask: '-###,###,###.##'), skip: 0),
				currency: ("GrandTotal", ("Text", "CAD"), skip: 0))
			("Vskip")))
		}

	Append(fmt)
		{
		.appended.Add(fmt)
		}
	}