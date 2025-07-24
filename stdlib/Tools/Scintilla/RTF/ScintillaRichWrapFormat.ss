// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
WrapFormat
	{
	formats: false

	GetDataSize(data, w, h, lineLimit, font /*unused*/)
		{
		.formats = ParseHTMLRichText.GetFormats(data, w, _report, :lineLimit)
		return h * .formats.Size()
		}

	// wrap line is handled by ParseHTMLRichText
	WrapDataLines(data, w /*unused*/, font /*unused*/)
		{
		return data
		}

	Print(x, y, w, h, data = false)
		{
		// _report.AddRect(x, y, w, h, 5)
		if .Data isnt false
			data = .Data

		if .formats is false
			.formats = ParseHTMLRichText.GetFormats(data, w, _report,
				lineLimit: .GetLineLimit())

		height = h / .formats.Size()
		top = y

		for format in .formats
			{
			left = x
			for item in format
				{
				oldfont = _report.SelectFont(item.font)
				fmt = _report.Construct(item)
				size = fmt.GetSize().w
				super.PrintData(left, top, size, height, fmt.Data)
				_report.SelectFont(oldfont)
				left += size
				}
			top += height
			}
		}

	ExportCSV(data = false)
		{
		super.ExportCSV(ScintillaRichStripHTML(data))
		}
	}
