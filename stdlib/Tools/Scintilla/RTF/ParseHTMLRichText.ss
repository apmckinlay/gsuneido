// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
class
	{
	GetFormats(data, w, rpt, font = false, lineLimit = false)
		{
		data = .ensureSpanned(data)
		if false is parser = .GetParsedText(data)
			return Object(Object(Object("Text", data)))

		return .parseText(parser, w, rpt, font, lineLimit)
		}

	ensureSpanned(data)
		{
		if not data.Prefix?('<span style=')
			{
			spanIdx = data.Find('<span')
			data = "<span style=''>" $ XmlEntityEncode(data[..spanIdx]) $ "</span>" $
				data[spanIdx..]
			}
		return data
		}

	GetParsedText(data)
		{
		if data.Prefix?('<span style=')
			{
			data = data.Replace('<br />', '\&crlf;')
			data = data.Replace(ScintillaRichEditorEncodeMap['\t'], '\t')
			scintillaRichParser = XmlParser
				{
				IgnorableWhitespace(s)
					{
					.Characters(XmlEntityDecode(s))
					}
				}
			try
				if false is parsed = scintillaRichParser('<html>' $ data $ '</html>')
					return false
			catch (unused, "*Invalid xml|*unmatched tag|*XmlReader:")
				return false
			return parsed
			}
		return false
		}
	reachedLineLimit(lines, lineLimit)
		{
		return false isnt lineLimit and lines >= lineLimit
		}
	parseText(parser, w, rpt, font, lineLimit)
		{
		rowFormat = Object()
		current = Object(Row: Object(), Size: 0)
		// for each format
		linefont = .parseChildren(parser, current, rowFormat, w, rpt, font, lineLimit)

		// handle the very last item
		if not current.Row.Empty?()
			rowFormat.Add(current.Row)

		if .reachedLineLimit(rowFormat.Size(), lineLimit)
			rowFormat.Add(Object(Object('Text', "...", font: linefont)))

		return rowFormat
		}
	parseChildren(parser, current, rowFormat, w, rpt, font = false, lineLimit = false)
		{
		reportFont = font isnt false ? font : rpt.GetFont()
		linefont = ''
		for child in parser.Children()
			{
			if .reachedLineLimit(rowFormat.Size(), lineLimit)
				break

			linefont = .getFont(child, reportFont.Copy())
			mainStr = child.Text().Replace('&crlf;', '\r\n')

			firstLine? = true
			lines = mainStr.Lines()
			if mainStr.Suffix?('\n')
				lines = lines.Concat([''])

			for line in lines
				{
				if firstLine? is false
					.startNewPageLine(current, rowFormat)
				else
					firstLine? = false

				if false is .parseLine(current, rowFormat,
					line, w, rpt,
					linefont, lineLimit)
					break
				}
			}
		return linefont
		}

	parseLine(current, rowFormat, str, w, rpt, linefont, lineLimit)
		{
		if .reachedLineLimit(rowFormat.Size(), lineLimit)
			return false
		// While we still have text from str left to put on the page
		while str.Size() isnt 0
			{
			// get substr that will fit in the remaining line space
			substr = .findTextFit(str, current.Size, w, linefont, rpt)

			if substr isnt ''
				{
				// add the text and measure it
				result = .measureItem(substr, linefont, rpt)
				current.Row.Add(result.item)
				current.Size += result.size
				}
			// get what text will not fit on the rest of the page
			// and try again on the next line
			if (((str = str.AfterFirst(substr)).Size() isnt 0))
				.startNewPageLine(current, rowFormat)

			if .reachedLineLimit(rowFormat.Size(), lineLimit)
				break
			}
		return true
		}

	getFont(child, font)
		{
		// need to handle text that has been appended to the end
		// might not have rtf formatting
		format = child.Attributes().Member?('style')
			? child.Attributes().style
			: ""
		if format =~ (':bold(;|$)')
			font.weight = 'bold'
		if format =~ (':italic(;|$)')
			font.italic = true
		if format.Has?(':underline')
			font.underline = true
		if format.Has?(' line-through')
			font.strikeout = true
		return font
		}

	measureItem(str, font, rpt)
		{
		item = Object('Text', str.Detab(), :font)
		fmt = rpt.Construct(item)
		return Object(size: fmt.GetSize().w, :item)
		}

	measureItemSize(str, font, rpt)
		{
		return .measureItem(str, font, rpt).size
		}

	startNewPageLine(current, rowFormat)
		{
		rowFormat.Add(current.Row)
		current.Row = Object()
		current.Size = 0
		}

	findTextFit(str, currentSize, w, font, rpt)
		{
		// measure the remaining text
		size = .measureItem(str, font, rpt).size

		// calculate how much space we have left on the line
		available = w - currentSize

		// if we've reached the end of the line, return '' to start new line
		if available <= 0
			return ''

		// if everything will fit, return full length
		if size <= available
			return str

		// if the str is one word that would fit on one line OR
		// if the first word of the str doesn't fit on the line AND
		// if it is not already on a fresh line
		// return '' to start new line
		if (((pos = str.Find1of(', ')) is str.Size() and size < w) or
			(.measureItemSize(str[.. pos], font, rpt) > available) and
			currentSize > 0)
			return ''

		// Following Function is shared with WrapFormat, doesn't handle above cases
		return TextBestFit(available, str, { .measureItemSize(it, :font, :rpt)},
			report: _report)
		}
	}
