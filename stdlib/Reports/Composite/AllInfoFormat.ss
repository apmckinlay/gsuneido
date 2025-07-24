// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
Format
	{
	New(.prefix = '', data = false, .font = false, .width = false)
		{
		.Data = data
		}

	Variable?()
		{ return true }

	GetSize(data = #{})
		{
		fmt = .make_format(data)
		return fmt.GetSize()
		}

	Print(x, y, w, h, data = #{})
		{
		fmt = .make_format(data)
		fmt.Print(x, y, w, h)
		}

	maxInfoFields: 10
	make_format(data)
		{
		if .Data isnt false
			data = .Data
		fmt = Object()
		for (i = 1; i <= .maxInfoFields; ++i)
			{
			field = .prefix $ 'info' $ i
			if data.Member?(field)
				{
				j = data[field].Find(':')
				prompt = data[field][.. j + 1]
				info = StripInfoLabel(data[field][j + 2 ..])
				font = .font is false ? _report.GetFont() : .font

				fmt.Add(Object(
					Object('Text' prompt, :font, justify: 'right')
					Object('Text' info, :font)))
				}
			}
		width = .width is false
			? 3000	/*= default width */
			: .width * 180 /*= charLength in Twips */
		format = Object('Grid', fmt, :width, top:, xstretch: .Xstretch)
		if fmt.Size() is 2
			format.Add(Object('Text', '', width: 40))
		return _report.Construct(format)
		}
	ExportCSV(data = #())
		{
		if .Data isnt false
			data = .Data
		csv = ''
		for (i = 1; i <= .maxInfoFields; ++i)
			{
			field = .prefix $ 'info' $ i
			if data.Member?(field)
				{
				j = data[field].Find(':')
				prompt = data[field][.. j + 1]
				info = data[field][j + 2 ..]
				csv $= prompt $ info $ ' '
				}
			}
		return .CSVExportString(csv)
		}
	}