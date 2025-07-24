// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
WrapFormat
	{
	New(data = false, width = false, w = 2880, fontsize = 9)
		{
		super(data, :width, :w, font: [name: "@mono", size: fontsize])
		}

	WrapDataLines(data, w, font)
		{
		maxwidth = w / _report.GetCharWidth(1, font, widthChar: '9')
		data = data.Replace('\t', '    ')
		for (i = 0; data[i] is ' '; ++i)
			{}
		indent = data[.. i]
		result = ""
		oline = indent
		indent $=  '    '
		owidth = indent.Size()
		for (j = 0, n = data.Size(); i < n; i += j)
			{
			// get next word (with leading whitespace)
			for (j = 0; data[i + j] is ' '; ++j)
				{}
			j += data[i + j..].Find(' ')
			word = data[i :: j]
			wwidth = word.Size()
			if owidth + wwidth < maxwidth
				{
				oline $= word
				owidth += wwidth
				}
			else
				{
				result $= oline $ "\r\n"
				oline = indent $ word.LeftTrim()
				owidth = indent.Size() * 4 + wwidth
				}
			}
		result $= oline
		return result
		}
	}
