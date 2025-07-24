// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Format
	{
	Xstretch: ""
	Ystretch: ""
	New(@items)
		{
		.Data = items.GetDefault("data", false)
		.font = items.GetDefault("font", false)
		.top = items.GetDefault("top", false)
		.newlines = items.GetDefault('newlines', false)
		.maxHeight = items.GetDefault('maxHeight', false)
		.items = Object()
		if items.Size() > 0
			{
			oldfont = _report.SelectFont(.font)
			for (i = 0; items.Member?(i); ++i)
				{
				item = _report.Construct(items[i])
				.items.Add(item)
				}
			_report.SelectFont(oldfont)
			}
		.set_xstretch? = .Xstretch is ""
		.set_ystretch? = .Ystretch is ""
		.recalc()
		}
	recalc()
		{
		ystretch = xstretch = 0
		for (item in .items)
			{
			ystretch += item.Ystretch
			if (item.Xstretch > xstretch)
				xstretch = item.Xstretch
			}
		if (.set_xstretch?)
			.Xstretch = xstretch
		if (.set_ystretch?)
			.Ystretch = ystretch
		.content_ystretch = ystretch
		}
	GetSize(data = #{}) // const
		{
		if .Data isnt false
			data = .Data
		oldfont = _report.SelectFont(.font)
		w = h = d = 0
		for i in .items.Members()
			{
			item = .items[i]
			if item.Y >= _report.GetDimens().H
				continue // page footer
			size = .getItemSize(item, data)
			if size.w > w
				w = size.w
			if item.Y > h
				h = item.Y
			if i is 0
				d = size.d
			else
				d += size.h
			h += size.h
			}
		_report.SelectFont(oldfont)
		return .buildSize(w, h, d)
		}
	getItemSize(item, data)
		{
		return item.GetSize(item.Member?("Field") and data.Member?(item.Field)
			? data[item.Field] : data)
		}
	buildSize(w, h, d)
		{
		w = Max(w, .Xmin)
		h = Max(h, .Ymin)
		if .maxHeight isnt false and .maxHeight < h
			{
			d -= h - .maxHeight
			h = .maxHeight
			}
		if .top isnt true
			d = 0
		return Object(:w, :h, :d)
		}
	AddConstructed(item) // used by Report for main page vbox
		{
		.items.Add(item)
		.recalc()
		}
	InsertConstructed(item)
		{
		.items.Add(item, at: 0)
		.recalc()
		}
	EraseTrailingHeader()
		{
		if 0 is n = .items.Size(list:)
			return

		while n > 0 and .items[n-1].Header?
			{
			.items.Delete(n-1)
			--n
			}

		.recalc()
		}
	Tally()
		{ return .items.Size() }

	ContentTally()
		{
		n = .items.Size(list:)
		while n > 0 and .items[n-1].Header?
			--n
		return n
		}

	Print(x, y, w, h, data = #{})
		{
		// _report.AddRect(x, y, w, h, 5)
		if .Data isnt false
			data = .Data
		size = .GetSize(data)
		f = Max(0, (h - size.h) / .content_ystretch)
		oldfont = _report.SelectFont(.font)
		for (orgy = y, y = i = 0; .items.Member?(i); ++i)
			{
			item = .items[i]
			value = item.Member?("Field") ? data[item.Field] : data
			size = item.GetSize(value)
			h = size.h + f * item.Ystretch
			if item.Y > y
				y = item.Y

			if .addEllipsis?(y, h, i)
				{
				.ellipsis.Print(x, orgy + y, w, .ellipsisHeight)
				break
				}
			item.Print(x, orgy + y, item.Xstretch > 0 ? w : size.w, h, value)
			item.Hotspot(x, orgy + y, item.Xstretch > 0 ? w : size.w, h, data)
			y += h
			}
		_report.SelectFont(oldfont)
		}
	addEllipsis?(y, h, i)
		{
		if .maxHeight is false
			return false

		if y + h > .maxHeight
			return true

		if ((i isnt .items.Size() - 1) and (.maxHeight - y - h < .ellipsisHeight))
			return true

		return false
		}
	getter_ellipsis()
		{
		return .ellipsis = _report.Construct(#('Text', '...'))
		}
	getter_ellipsisHeight()
		{
		return .ellipsisHeight = .ellipsis.GetSize().h
		}

	ExportFile: false
	ExportCSV(data = #{})
		{
		if (.Data isnt false)
			data = .Data
		str = ''
		for (i = 0; .items.Member?(i); ++i)
			{
			item = .items[i]
			if not item.Export
				continue
			value = item.Member?("Field") ? data[item.Field] : data
			s = item.ExportCSV(value)
			if .ExportFile isnt false
				{ // top level vbox
				s = s.Trim().Replace('\n\n+', '\n')
				if s isnt ''
					.ExportFile.Writeline(s)
				}
			else // other nested Vert's
				str $= s $ (.newlines ? '\n' : '')
			}
		return str
		}
	GetItems()
		{ return .items }
	}
