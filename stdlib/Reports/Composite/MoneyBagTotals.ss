// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(@args)
		{
		append = args[0]
		cur_col = args[1]
		type = args.Member?(2) ? args[2] : 'Total'
		currencies = Object()
		formats = Object()

		.delete_members(args)
		additional_totals = .set_additional_totals(args)
		skip = .set_skip(args)

		amounts = .get_currencies(args, currencies, formats)

		counter = .print_moneybag_totals(
			currencies, amounts, cur_col, type, args, formats, additional_totals, append)

		.print_additional_totals(additional_totals, counter, append)

		if skip is true
			append(#(Vskip))
		}

	print_additional_totals(additional_totals, counter, append)
		{
		while additional_totals.Member?(counter)
			{
			format = Object('_output')
			for field in additional_totals[counter].Members()
				format[field] = additional_totals[counter][field]
			append(format)
			++counter
			}
		}

	print_moneybag_totals(currencies, amounts, cur_col, type, args, formats,
		additional_totals, append)
		{
		first = true
		counter = 0
		for (cur in currencies)
			{
			if cur is "" and amounts[cur] is 0
				continue

			format = Object('_output')
			format[cur_col] = .cur_format(first, type, cur)

			for (fld in args.Members())
				{
				format[fld] = formats[fld]
				format[fld][1] = args[fld][cur]
				if (first)
					format[fld] = type is ""
						? format[fld]
						: Object(type, format[fld], skip: 0)
				}

			.apply_additional_fields(additional_totals, counter, format)

			append(format)
			first = false
			++counter
			}
		return counter
		}

	cur_format(first, type, cur)
		{
		return first and type isnt ""
			? Object(type, Object('Text' cur) skip: 0)
			: Object('Text', cur)
		}

	apply_additional_fields(additional_totals, counter, format)
		{
		if (additional_totals.Member?(counter))
			for field in additional_totals[counter].Members()
				format[field] = additional_totals[counter][field]
		}

	get_currencies(args, currencies, formats)
		{
		amounts = false
		for (fld in args.Members())
			{
			amounts = args[fld] = args[fld].Amounts()
			currencies.MergeUnion(amounts.Members())
			fmt = Datadict(fld)
			formats[fld] = fmt is Field_string // default returned if no Field_ definition
				? Object('OptionalNumber' mask: "-###,###,###.##") : fmt.Format.Copy()
			}
		return amounts
		}
	set_skip(args)
		{
		skip = args.GetDefault('vskip', true) is true
		args.Delete('vskip')
		return skip
		}
	set_additional_totals(args)
		{
		additional_totals = args.GetDefault('additional_totals', Object())
		args.Delete('additional_totals')
		return additional_totals
		}

	delete_members(args)
		{
		// delete these so later for loops don't see them
		// delete backwards since delete "shifts"
		args.Delete(2, 1, 0)
		}
	}
