// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(font, oldfont, factor = 1)
		{
		font = font.Copy()
		font = .buildFontPDFSize(font, oldfont)
		return .ensure(font, factor)
		}
	buildFontPDFSize(font, oldfont)
		{
		if font.Member?(#size)
			{
			if String?(font.size) and oldfont isnt false and oldfont.Member?('size') and
				Number?(oldfont.size) // e.g. '+3' or '-1'
				font.size = Number(font.size) + oldfont.size
			}
		else if oldfont isnt false and oldfont.Member?(#size)
			font.Add(oldfont.size, at: 'size')
		return font
		}
	ensure(font, factor)
		{
		font.MergeNew(_report.GetDefaultFont())  //ensure font object is complete
		font.size *= factor
		return font
		}
	}