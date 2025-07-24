// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Format
	{
	New(.Data = false, .w = false, .width = false, .justify = "left",
		.font = false, .color = false, .Export = true, access = false)
		{
		.InitAccessField(access)
		}
	WidthChar: 'M'
	GetSize(data = "")
		{
		.DoWithFont(.font)
			{|font|
			lineSpecs = _report.GetLineSpecs(font)
			if _report.GetDefault('Measuring?', false)
				w = .measure(data, font)
			else if .w isnt false
				w = .w
			else if .width isnt false
				w = _report.GetCharWidth(.width, font, .WidthChar)
			else
				w = .measure(data, font)
			}
		return Object(h: lineSpecs.height, d: lineSpecs.descent, :w)
		}

	measure(data, font)
		{
		data = .to_str(data)
		if data.Has1of?('\r\n')
			data = data.FirstLine() $ "... "
		return _report.GetTextWidth(font, data)
		}

	GetDefaultWidth()
		{
		return .width
		}

	Print(x, y, w, h, data = "")
		{
		data = .to_str(data)
		// append 200 spaces so DT.END_ELLIPSIS can be invoked
		spacesRequired = 200
		if data.Has1of?('\r\n')
			data = data.FirstLine() $ ' '.Repeat(spacesRequired)

		.DoWithFont(.font)
			{|font|
			.print(x, y, w, h, font, data)
			}
		}
	print(x, y, w, h, font, data)
		{
		ellipsis? = .w isnt false or .width isnt false
		_report.DrawWithinClip(x, y, w, h)
			{
			_report.AddText(data, x, y, w, h, font, .justify, ellipsis?, .color)
			}
		}

	ExportCSV(data = '')
		{
		data = .to_str(data)
		return .CSVExportString(data)
		}

	to_str(data)
		{
		if .Data isnt false
			data = .Data
		data = .ConvertToStr(data)
		if not String?(data)
			data = Display(data)
		return data
		}

	ConvertToStr(data)
		{
		return data
		}

	GetFont()
		{ return .font }
	GetJustify()
		{ return .justify }
	}
