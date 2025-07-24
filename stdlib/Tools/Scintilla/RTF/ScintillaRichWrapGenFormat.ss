// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Generator
	{
	New(data, w = false, font = false, lineLimit = 6000)
		{
		.w = w is false ? _report.GetWidth() : w
		.formats = _report.PlainText?()
			? Object(Object(Object('Wrap', ScintillaRichStripHTML(data))))
			: ParseHTMLRichText.GetFormats(data, .w, _report, :font, :lineLimit)
		.row = 0
		}

	Next()
		{
		if .row >= .formats.Size()
			return false
		rowfmt = .formats[.row]
		fmt = Object('Horz')
		for col in rowfmt.Members()
			fmt.Add(rowfmt[col])
		if rowfmt.Size() is 0
			fmt.Add(Object('Text', ''))
		++.row
		return _report.Construct(fmt)
		}
	}
