// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// divide input up into "pages"
// uses InputFormat
Generator
	{
	New(@items)
		{
		.H = _report.Report_dimens.H - 4.InchesInTwips() /*= calc header height*/
		if items.Member?("ymin")
			{
			.H = items.ymin
			items.Delete("ymin")
			}
		if _report.PlainText?() and .H < 0
			.H = 0
		.input = new InputFormat(@items)
		}
	Next()
		{
		if false is item = .input.Next()
			return false
		vh = 0
		vbox = VertFormat()
		vbox.H = .H
		do
			{
			if item is "pg"
				break
			h = _report.PlainText?() ? 0 : item.GetSize().h
			if item.Member?("Y") and item.Y > vh
				vh = item.Y
			if vh + h > .H
				{
				if not item.Header?
					.input.Pushback(item)
				break
				}
			vbox.AddConstructed(item)
			vh += h
			}
			while (false isnt item = .input.Next())
		return vbox
		}
	}
