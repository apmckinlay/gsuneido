// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
function (date)
	{
	widths = Object(wYear: 4, wMilliseconds: 3).Set_default(2)
	s = '#'
	for m in #(wYear, wMonth, wDay, wHour, wMinute, wSecond, wMilliseconds)
		s $= date[m].Pad(widths[m])
	datePartLength = 9
	s = s[.. datePartLength] $ '.' $ s[datePartLength ..]
	return Date(s)
	}