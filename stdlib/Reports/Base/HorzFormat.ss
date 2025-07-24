// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Format
	{
	Xstretch: ""
	Ystretch: ""
	New(@items)
		{
		.Data = items.Member?("data") ? items.data : false
		.font = items.Member?("font") ? items.font : false
		oldfont = _report.SelectFont(.font)
		.items = Object()
		for (i = 0; items.Member?(i); ++i)
			{
			item = _report.Construct(items[i])
			.items.Add(item)
			}
		_report.SelectFont(oldfont)

		ystretch = xstretch = 0
		for (item in .items)
			{
			xstretch += item.Xstretch
			if (item.Ystretch > ystretch)
				ystretch = item.Ystretch
			}
		if (.Xstretch is "")
			.Xstretch = xstretch
		if (.Ystretch is "")
			.Ystretch = ystretch
		.content_xstretch = xstretch
		}
	GetSize(data = #{}, widths = false)
		{
		if (.Data isnt false)
			data = .Data
		oldfont = _report.SelectFont(.font)
		w = ascent = descent = 0
		for (i in .items.Members())
			{
			item = .items[i]
			size = item.GetSize(item.Member?("Field") ? data[item.Field] : data)
			if (item.X > w)
				w = item.X
			if widths isnt false
				widths.Add(size.w)
			w += size.w
			if (size.h - size.d > ascent)
				ascent = size.h - size.d
			if (size.d > descent)
				descent = size.d
			}
		_report.SelectFont(oldfont)
		return Object(w: Max(.Xmin, w), h: Max(.Ymin, ascent + descent), d: descent)
		}
	Print(orgx, orgy, w, h, data = #{})
		{
		// _report.AddRect(orgx, orgy, w, h, 5)
		if (.Data isnt false)
			data = .Data
		oldfont = _report.SelectFont(.font)
		totsize = .GetSize(data)
		f = Max(0, (w - totsize.w) / .content_xstretch)
		scale = w >= totsize.w ? 1 : w / totsize.w
		for (x = i = 0; .items.Member?(i); ++i)
			{
			item = .items[i]
			value = item.Member?("Field") ? data[item.Field] : data
			size = item.GetSize(value)
			if (item.X > x)
				x = item.X
			w = scale * size.w + f * item.Xstretch
			y = orgy + totsize.h - totsize.d + size.d - size.h
			item.Print(orgx + x, y, w, h - (y - orgy), value)
			item.Hotspot(orgx + x, y, w, h - (y - orgy), data)
			x += w
			}
		_report.SelectFont(oldfont)
		}
	ExportCSV(data = '')
		{
		if (.Data isnt false)
			data = .Data
		strAll = ''
		for (i = 0; .items.Member?(i); ++i)
			{
			item = .items[i]
			if not item.Export
				continue
			value = item.Member?("Field") ? data[item.Field] : data
			str = item.ExportCSV(value)
			//if the item already have a new line at end, should not need comma.
			strAll $= str.Suffix?('\n')
				? str
				: Opt(str, ',')
			}
		return .CSVExportLine(strAll)
		}
	}