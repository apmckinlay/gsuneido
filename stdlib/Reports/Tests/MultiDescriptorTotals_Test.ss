// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	initMoneyBags()
		{
		return Object(total_amount: new MoneyBag,
			total_amount_converted: new MoneyBag,
			total_qty1: new MoneyBag,
			total_qty2: new MoneyBag)
		}

	Test_noData()
		{
		.appended = Object()
		mbOb = .initMoneyBags()
		MultiDescriptorTotals(.Append,
			Object(
				currency: Object(
					arivc_amount: mbOb.total_amount,
					arivc_amount_converted: mbOb.total_amount_converted),
				uomField: Object(
					arivc_qty1: mbOb.total_qty1,
					arivc_qty2: mbOb.total_qty2)))

		// no data; vskip defaults to true
		Assert(.appended is: #((Vskip)))

		// no data; no vskip
		.appended = Object()
		MultiDescriptorTotals(.Append,
			Object(
				currency: Object(
					arivc_amount: mbOb.total_amount,
					arivc_amount_converted: mbOb.total_amount_converted),
				uomField: Object(
					arivc_qty1: mbOb.total_qty1,
					arivc_qty2: mbOb.total_qty2))
			vskip: false)
		Assert(.appended.Empty?())
		}

	Test_main()
		{
		.appended = Object()
		mbOb = .initMoneyBags()
		mbOb.total_amount.Plus(100, 'CAD')
		mbOb.total_amount.Plus(212, 'USD')

		mbOb.total_amount_converted.Plus(75, "CAD")
		mbOb.total_amount_converted.Plus(115, 'USD')

		mbOb.total_qty1.Plus(45, 'lbs')
		mbOb.total_qty1.Plus(8, 'ton')
		mbOb.total_qty1.Plus(11, 'cwt')
		mbOb.total_qty1.Plus(0, '')

		mbOb.total_qty2.Plus(18, 'lbs')
		mbOb.total_qty2.Plus(7, 'ton')
		mbOb.total_qty2.Plus(18, 'cwt')
		mbOb.total_qty2.Plus(9, 'hours')
		mbOb.total_qty2.Plus(3, 'm3')

		// total_amount: #(CAD: 100, USD: 212)
		// total_amount_convertetd: #(CAD: 75, USD: 115)
		// total_qty1: #(lbs: 45, ton: 8, cwt: 11, '': 0) - empty should NOT print
		// total_qty2: #(lbs: 18, ton: 7, cwt: 18, hours: 9)

		MultiDescriptorTotals(.Append,
			Object(
				currency: Object(
					amount: mbOb.total_amount,
					amount_converted: mbOb.total_amount_converted),
				uomField: Object(
					qty1: mbOb.total_qty1,
					qty2: mbOb.total_qty2))
			additionalItemsToPrint: #('customTotal')
			data: [customTotal: 5, name: 'Fred Flinstone'])

		Assert(.appended.Size() is: 6)
		// result should be:
		// qty1		qty2	uom		amount		amount_converted	cur   customTotal
		//	11		18		cwt		100			75					CAD		5
		//  		9		hours	212			115					USD
		//	45		18		lbs
		//  				m3
		//  8		7		ton
		Assert(.appended[0] is: #("_output",
			qty1: #("Total", "qty1", skip: 0),
			qty2: #("Total", "qty2", skip: 0),
			customTotal: #("Total" "customTotal", skip: 0)
			uomField: #("Total", "uomField", skip: 0),
			amount_converted: #("Total", "amount_converted", skip: 0),
			amount: #("Total", "amount", skip: 0)
			currency: #("Total", "currency", skip: 0),
			data: [qty1: 11, uomField: "cwt", qty2: 18, customTotal: 5,
				name: 'Fred Flinstone'
				amount_converted: 75, amount: 100, currency: "CAD"]))

		Assert(.appended[1] is: #("_output",
			amount: "amount", qty2: "qty2", amount_converted: "amount_converted",
			uomField: "uomField", currency: "currency", qty1: "qty1"
			data: [uomField: "hours", qty1: 0, amount: 212, qty2: 9,
				amount_converted: 115, currency: "USD"]))

		Assert(.appended[2] is: #("_output",
			qty1: "qty1", qty2: "qty2", uomField: "uomField"
			data: [uomField: "lbs", qty1: 45, qty2: 18]))

		Assert(.appended[3] is: #("_output",
			qty1: "qty1", qty2: "qty2", uomField: "uomField"
			data: [uomField: "m3", qty1: 0, qty2: 3]))
		Assert(.appended[4] is: #("_output",
			qty1: "qty1", qty2: "qty2", uomField: "uomField"
			data: [uomField: "ton", qty1: 8, qty2: 7]))
		Assert(.appended[5] is: #("Vskip"))
		}

	Test_oneMoneyBagAllZero()
		{
		.appended = Object()
		mbOb = .initMoneyBags()
		mbOb.total_amount.Plus(100, 'CAD')
		mbOb.total_amount_converted.Plus(0, "")
		mbOb.total_qty1.Plus(45, 'lbs')
		MultiDescriptorTotals(.Append,
			Object(
				currency: Object(
					amount: mbOb.total_amount,
					amount_converted: mbOb.total_amount_converted),
				uomField: Object(qty1: mbOb.total_qty1)))

		// result should be (does not print "" desc):
		// qty1		qty2	uom		amount		amount_converted	cur
		//	45				lbs		100			0					CAD		5
		Assert(.appended.Size() is: 2)
		Assert(.appended[0] is: #("_output",
			uomField: #("Total", "uomField", skip: 0),
			qty1: #("Total", "qty1", skip: 0),
			currency: #("Total", "currency", skip: 0),
			amount_converted: #("Total", "amount_converted", skip: 0),
			amount: #("Total", "amount", skip: 0),
			data: [uomField: "lbs", qty1: 45, currency: "CAD", amount_converted: 0,
				amount: 100]))
		Assert(.appended[1] is: #("Vskip"))
		}


	Append(fmt)
		{
		.appended.Add(fmt)
		}
	}