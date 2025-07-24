// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
StaticControl
	{
	CalcLines(text, orig_xmin)
		{
		xmin = orig_xmin isnt 0 ? orig_xmin : 700
		lines = 0
		for line in text.Lines()
			lines += .bestFit(line, xmin)
		.Xmin = xmin
		return Max(lines, 1)
		}
	bestFit(text, xmin)
		{
		lines = 0
		str = text
		while false isnt fitText = TextBestFit(xmin, str, .measure)
			{
			lines++
			if "" is str = str[fitText.Size() ..]
				break
			}
		return lines
		}
	measure(str)
		{
		return .TextExtent(str).x
		}
	}
