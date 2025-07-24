// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Format
	{
	New(@args)
		{
		.font = args.GetDefault("font", false)
		.Data = args.GetDefault("data", false)
		.access = args.GetDefault('access', false)
		.formats = Object()
		.variable = false
		// have to work backwards (right to left)
		// so you know how big to make spans
		for (i = args.Size(list:) - 1; i >= 0; --i)
			{
			arg = args[i]
			if (args.Member?("widths"))
				{
				skip = args.widths.GetDefault(#skip, HskipFormat.Size.w)
				arg = String?(arg) ? Object(arg) : arg.Copy()
				arg.w = args.widths[i]
				if arg.Member?('span')
					for (j = 0; j < arg.span; ++j)
						arg.w += skip + args.widths[i + j + 1]
				}
			// top-align images
			if Object?(arg) and arg.Member?(0) and arg[0] is "Image"
				{
				arg = arg.Copy()
				arg.top = true
				}
			f = .formats[i] = _report.Construct(arg)
			if (Object?(arg) and arg.Member?("heading"))
				f.Heading = arg.heading
			if (f.Variable?())
				.variable = true
			}
		if (args.Member?("widths"))
			.widths = args.widths
		else if _report.PlainText?()
			.widths = .fakeWidths()
		else
			.fontsize() // determines widths

		size = _report.PlainText?() ? .fakeHeights() : .getHeights()

		.Ymin = .height = size.h
		.descent = size.d
		.header = Object("ColHeads", .formats, .widths, font: .font)
		}
	GetWidths()
		{
		return .widths
		}
	GetFont()
		{
		return .font
		}
	GetSize(data = #())
		{
		size = .variable ? .getHeights(data) : Object(h: .height, d: .descent)
		size.w = .widths.total
		return size
		}
	getHeights(data = #{})
		{
		if (.Data isnt false)
			data = .Data
		ascent = descent = 0
		oldfont = _report.SelectFont(.font)
		for (fmt in .formats)
			{
			value = fmt.Member?("Field")
				? ((fmt.Variable?() or data.Member?(fmt.Field)) ? data[fmt.Field] : "")
				: data
			size = fmt.GetSize(value)
			if (size.h - size.d > ascent)
				ascent = size.h - size.d
			if (size.d > descent)
				descent = size.d
			}
		_report.SelectFont(oldfont)
		extra = _report.Params.PrintLines is true ? .thickness * 2 : 0
		return Object(h: extra + Max(.Ymin, ascent + descent), d: descent)
		}
	fakeHeights()
		{
		return Object(h: 0, d: 0)
		}
	Header()
		{ return .header }
	Print(x, orgy, w, h, data = #{})
		{
		if _report.Params.PrintLines is true
			.drawline(x, orgy + h, w)

		if (.Data isnt false)
			data = .Data
		totsize = .variable ? .getHeights(data) : Object(h: .height, d: .descent)
		skip = .widths.skip
		oldfont = _report.SelectFont(.font)
		span = 0
		for (i in .formats.Members())
			{
			if (--span >= 0)
				continue // skip over spanned columns
			f = .formats[i]
			w = .widths[i]
			if (f.Member?("Span"))
				for (span = f.Span, j = 0; j < span; ++j)
					w += skip + .widths[i + j + 1]
			value = f.Member?("Field") ? data[f.Field] : data
			size = f.GetSize(value)
			y = orgy + totsize.h - totsize.d + size.d - size.h
			f.Print(x, y, w, h, value)
			colAccess = .access isnt false and .access.Member?(i) ? .access[i] : false
			f.Hotspot(x, y, w, h, data, access: colAccess)
			x += w + skip
			}
		_report.SelectFont(oldfont)
		}
	ExportCSV(data = #{})
		{
		if (.Data isnt false)
			data = .Data
		span = 0
		str = ''
		for (i in .formats.Members())
			{
			if (--span >= 0)
				continue  // skip over spanned columns
			f = .formats[i]
			if not f.Export
				continue
			if f.Member?('PrepareData')
				f.PrepareData()
			value = f.Member?("Field") ? data[f.Field] : data
			str $= f.ExportCSV(value) $ ','
			}
		return .CSVExportLine(str)
		}
	thickness: 10
	drawline(x, y, w)
		{
		y += .thickness
		_report.AddLine(x, y, x + w, y, thick: .thickness, color: 0xbbbbbb)
		}
	fontsize()
		{
		// figure out the font size (assuming all the same font)
		width = _report.GetWidth()
		// allowable range of font sizes
		lo = 4
		hi = 12
		if false is oldfont = _report.GetFont()
			oldfont = _report.GetReportDefaultFont()
		if (.font is false)
			.font = oldfont
		else
			hi = .font.GetDefault('size', 12)
		.font = .font.Copy() // make a copy cause we're going to change size
		// binary search for right size
		widths = Object()
		forever
			{
			.font.size = med = (lo + hi) / 2
			_report.SelectFont(.font)
			skip = .font.size * 15
			w = -skip
			for (i in .formats.Members())
				{
				widths[i] = .formats[i].GetSize().w
				w += widths[i] + skip
				}
			if (hi - lo < .01)
				break
			if (w < width)
				lo = med
			else
				hi = med
			}
		widths.total = w
		widths.skip = skip
		.widths = widths
		_report.SelectFont(oldfont)
		}
	fakeWidths()
		{
		widths = Object()
		for (i in .formats.Members())
			widths[i] = 100
		widths.skip = 0
		widths.total = .formats.Size() * 100
		return widths
		}
	Output(fmt, data) // called by Query to construct "Output" formats
		{
		fmt = fmt.Copy()
		fmt[0] = "Row"
		if (not fmt.Member?("data"))
			fmt.data = data
		fmt.font = .font
		fmt.widths = .widths
		for (i in .formats.Members())
			{
			f = .formats[i]
			if (f.Member?("Field") and fmt.Member?(f.Field))
				fmt[i + 1] = fmt[f.Field]
			else if (not fmt.Member?(i + 1))
				fmt[i + 1] = ""
			}
		return _report.Construct(fmt)
		}
	}
