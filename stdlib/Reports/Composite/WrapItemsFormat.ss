// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
Format
	{
	New(@items)
		{
		.Data = items.Member?("data") ? items.data : false
		if items.Member?("font")
			.font = items.font
		if items.Member?("skip")
			.skip = items.skip
		.items = Object()
		for (item in items.Values(list:))
			.items.Add(_report.Construct(item))
		.dimens = _report.GetDimens()
		.paperWidth = (.dimens.width - .dimens.right).InchesInTwips()
		}
	lead: 20
	font: (name: 'Arial', size: 8, weight: 'extralight')
	skip: 400
	getSize(data)
		{
		if (.Data isnt false)
			data = .Data
		x = w = h = ascent = descent = 0
		for (item in .items)
			{
			value = item.Member?("Field") ? data[item.Field] : data
			size = item.GetSize(value)
			orgx = .dimens.left.InchesInTwips()
			if ((orgx + x + size.w) > .paperWidth)
				{
				if (x > w)
					w = x
				x = 0
				h += ascent + descent + .lead
				ascent = descent = 0
				}
			if (size.h - size.d > ascent)
				ascent = size.h - size.d
			if (size.d > descent)
				descent = size.d
			x += size.w + .skip
			}
		h += ascent + descent
		return Object(:w, :h, d: 0)
		}

	GetSize(data = #{})
		{
		.DoWithFont(.font)
			{|unused|
			ob = .getSize(data)
			}
		return ob
		}
	doPrint(orgx, orgy, data)
		{
		if (.Data isnt false)
			data = .Data

		// determine the vertical sizes of each row
		x = row = 0
		ascents = Object().Set_default(0)
		descents = Object().Set_default(0)
		for (item in .items)
			{
			value = item.Member?("Field") ? data[item.Field] : data
			size = item.GetSize(value)
			if ((orgx + x + size.w) > .paperWidth)
				{
				x = 0
				++row
				}
			if (size.h - size.d > ascents[row])
				ascents[row] = size.h - size.d
			if (size.d > descents[row])
				descents[row] = size.d
			x += size.w + .skip
			}
		.print(data, orgx, orgy, ascents, descents)
		}

	Print(orgx, orgy, w /*unused*/, h /*unused*/, data = #{})
		{
		.DoWithFont(.font)
			{|unused|
			.doPrint(orgx, orgy, data)
			}
		}

	print(data, orgx, orgy, ascents, descents)
		{
		x = row = 0
		for (item in .items)
			{
			value = item.Member?("Field") ? data[item.Field] : data
			size = item.GetSize(value)
			if ((orgx + x + size.w) > .paperWidth)
				{
				x = 0
				orgy += ascents[row] + descents[row] + .lead
				++row
				}
			y = orgy + ascents[row] + size.d - size.h
			// _report.AddRect(x, y, size.w, size.h, 5)
			item.Print(orgx + x, y, size.w, size.h, value)
			item.Hotspot(orgx + x, y, size.w, size.h, data)
			x += size.w + .skip
			}
		}
	ExportCSV(data = #{})
		{
		if (.Data isnt false)
			data = .Data

		str = ''
		for (item in .items)
			{
			value = item.Member?("Field") ? data[item.Field] : data
			csvVal = item.ExportCSV(value)
			if csvVal.Prefix?('\n') and str.Suffix?(',')
				str = str[.. -1]
			str $= csvVal $ ','
			}
		return .CSVExportLine(str)
		}
	}