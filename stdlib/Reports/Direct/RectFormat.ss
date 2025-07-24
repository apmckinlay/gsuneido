// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Format
	{
	New(item, .height = false, .width = false, .thick = 10, .outside = 120, inside = 120)
		{
		.item = _report.Construct(item)
		.Xstretch = .item.Xstretch
		.Ystretch = .item.Ystretch
		.extra = inside + outside + thick
		}
	GetSize(data = false)
		{
		sz = .item.GetSize(.item.Member?("Field") ? data[.item.Field] : data)
		sz.w += 2 * .extra
		sz.h += 2 * .extra
		if .width isnt false
			sz.w = .width
		if .height isnt false
			sz.h = .height
		return sz
		}
	Print(x, y, w, h, data = false)
		{
		_report.AddRect(x + .outside, y + .outside,
			w - .outside * 2, h - .outside * 2, .thick)
		itemWidth = w - 2 * .extra
		itemHeight = h - 2 * .extra
		_report.DrawWithinClip(x + .extra, y + .extra, itemWidth, itemHeight)
			{
			.item.Print(x + .extra, y + .extra, itemWidth, itemHeight,
				.item.Member?("Field") ? data[.item.Field] : data)
			}
		}
	ExportCSV(data = false)
		{
		return .item.ExportCSV(data)
		}
	}
