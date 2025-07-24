// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
class
	{
	/*0: .Append,
		1: #(currency: #(field: value, field: value), uom: #(field: value)),
		2: 'Total' vs 'GrandTotal (totalFormat),
		additionalItemsToPrint: Object(field1, field2, field3)
			additionalItemsToPrint will ONLY print on the first line
		data: optional (for printing additional totals, only used on first line)
		vskip - optional
	*/
	CallClass(@args)
		{
		append = args[0]
		breakdownTypes = args[1]
		totalFormat = args.Member?(2) ? args[2] : 'Total'
		// #(0: #('cad', 'usd'), 1: #('lbs', 'ton'))
		breakdownDesc = Object().Set_default(Object())
		formats = Object()
		additionalItemsToPrint = args.GetDefault('additionalItemsToPrint', Object())
		data = args.GetDefault('data', [])

		.buildBreakdownDesc(breakdownTypes, breakdownDesc, formats)
		skip = .set_skip(args)

		.print(breakdownDesc, breakdownTypes, totalFormat, data, append,
			additionalItemsToPrint)

		if skip is true
			append(#(Vskip))
		}

	buildBreakdownDesc(breakdownTypes, breakdownDesc, formats)
		{
		for desc in breakdownTypes.Members()
			.initObjects(breakdownTypes[desc], breakdownDesc[desc], formats)
		}

	initObjects(moneybag, currencies, formats)
		{
		for fld in moneybag.Members()
			{
			amounts = moneybag[fld] = moneybag[fld].Amounts()
			// avoid printing 0 under empty description
			amounts.DeleteIf({ it is "" and Number(amounts[it]) is 0})
			currencies.MergeUnion(amounts.Members())
			currencies.Sort!()
			fmt = Datadict(fld)
			formats[fld] = fmt is Field_string // default returned if no Field_ definition
				? Object('OptionalNumber' mask: "-###,###,###.##") : fmt.Format.Copy()
			}
		}

	set_skip(args)
		{
		skip = args.GetDefault('vskip', true) is true
		args.Delete('vskip')
		return skip
		}

	print(breakdownDesc, breakdownTypes, totalFormat, data, append,
		additionalItemsToPrint = #())
		{
		max = 0
		breakdownDesc.Each( { if it.Size() > max  max = it.Size() })
		first = true
		for (idx = 0; idx < max; ++idx)
			{
			fmt = Object('_output', data: first ? data.Copy() : [])
			.buildFormat(fmt, idx, breakdownDesc, breakdownTypes, first, totalFormat)
			if first // print additional totals on first row
				for field in additionalItemsToPrint
					fmt[field] = Object(totalFormat, field, skip: 0)

			append(fmt)
			first = false
			}
		}

	// Sample values:
	//breakdownDesc: #(
	//			uomField: #("ton", "", "lbs", "cwt", "hours", "m3"),
	//			currency: #("USD", "CAD"))
	//breakdownTypes: #(
	//	uomField: #(qty1: #(ton: 8, "": 0, lbs: 45, cwt: 11),
	//				qty2: #(hours: 9, ton: 7, lbs: 18, cwt: 18)),
	//	currency: #(amount: #(USD: 212, CAD: 100),
	//				amount_converted: #(USD: 115, CAD: 75)))
	buildFormat(fmt, idx, breakdownDesc, breakdownTypes,
		first = false, totalFormat = 'Total')
		{
		for field in breakdownDesc.Members()
			{
			if breakdownDesc[field].Member?(idx)
				{
				desc = breakdownDesc[field][idx]

				fmt[field] = first ? Object(totalFormat, field, skip: 0) : field
				fmt.data[field] = desc
				for column in breakdownTypes[field].Members().Sort!()
					{
					value = breakdownTypes[field][column][desc]
					fmt[column] = first ? Object(totalFormat, column, skip: 0): column
					fmt.data[column] = value
					}
				}
			}
		}
	}
