// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Format
	{
	New(data = false, .width = false, .w = false, .font = false, .justify = 'left',
		.sizeWithoutWeight = false, .lineLimit = 30, access = false)
		{
		.Data = data
		.InitAccessField(access)
		}
	WidthChar: "M"
	GetSize(data = false)
		{
		if .Data isnt false
			data = .Data
		data = .Format_data(data)
		return .getSize(data)
		}

	// used by sujswebgui:ListFormatting
	Format_data(data)
		{
		if not String?(data)
			data = Display(data)
		return data.Lines().Map(#RightTrim).Join('\n').Trim('\r\n')
		}

	getSize(data)
		{
		.DoWithFont(.font)
			{ |font|
			curFont = font
			lineSpecs = _report.GetLineSpecs(font)
			if .w isnt false
				w = .w
			else
				{
				f = font.Copy()
				if .sizeWithoutWeight is true
					f.Delete(#weight)
				w = _report.GetCharWidth(.width, font, .WidthChar)
				}
			sizeOb = Object(h: lineSpecs.height, d: lineSpecs.descent, :w)
			top = sizeOb.h - sizeOb.d
			if data isnt false
				sizeOb.h = .GetDataSize(data, sizeOb.w, sizeOb.h, .lineLimit, curFont)
			sizeOb.d = sizeOb.h - top
			}

		return sizeOb
		}

	GetDataSize(data, w, h, lineLimit /*unused*/, font)
		{
		lines = .WrapDataLines(data, w, font)
		return _report.GetTextHeight(lines, h)
		}

	GetLineLimit()
		{
		return .lineLimit
		}

	Print(x, y, w, h, data = false)
		{
		// _report.AddRect(x, y, w, h, 5)
		if .Data isnt false
			data = .Data
		.PrintData(x, y, w, h, data)
		}
	PrintData(x, y, w, h, data) // called by ScintillaRichWrapFormat
		{
		data = .Format_data(data)
		.DoWithFont(.font)
			{ |font|
			if _report.Driver.GetDefault(#NoWrapSupported?, false) is true
				_report.AddText(data, x, y, w, h, font, ellipsis?:)
			else
				{
				lines = .WrapDataLines(data, w, font)
				_report.AddMultiLineText(lines, x, y, w, h, font, .justify)
				}
			}
		}
	ExportCSV(data = false)
		{
		if .Data isnt false
			data = .Data
		data = .Format_data(data)
		return .CSVExportString(data)
		}

	WrapDataLines(data, w, font) // overriden in ScintillaRichWrapFormat
		{
		data = .Format_data(data)
		.data = data
		lines = ""
		lineCount = 0
		// TextBestFit function is shared with ParseHTMLRichText
		while (false isnt line = TextBestFit(w, .data, { .measure(it, :font) }, _report))
			{
			.data = .data[line.Size() ..]
			if ++lineCount is .lineLimit
				{
				lines $= "...\n"
				break
				}
			lines $= line.RightTrim() $ '\n'
			}
		return lines[.. -1]
		}

	measure(line, font)
		{
		return _report.GetTextWidth(font, line)
		}

	Variable?()
		{
		return true
		}
	}
