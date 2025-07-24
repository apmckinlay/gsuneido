// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// like Wrap but generates each line as a Text format
// uses DrawTextEx to determine line breaks
Generator
	{
	New(data, w = false, .font = false)
		{
		.text = String(data).Detab()
		.w = w is false ? _report.GetWidth() : w
		}
	Next()
		{
		if .text is ""
			return false

		layout = false
		.DoWithFont(.font)
			{ |font|
			if false isnt line = TextBestFit(.w, .text, { .measure(it, :font) }, _report)
				{
				.text = .text[line.Size() ..]
				layout = Object("Text", line, :font)
				}
			}
		return layout is false ? false : _report.Construct(layout)
		}

	measure(line, font)
		{
		return _report.GetTextWidth(font, line)
		}
	}
