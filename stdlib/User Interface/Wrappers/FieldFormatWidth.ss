// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
function (field, aveCharWidth, hdc = false)
	{
	width = 16
	fmt = Datadict(field).Format
	if fmt[0] is 'Image'
		return 100 	/*= images won't display, keep those columns small */
	else if fmt.Member?(#width)
		{
		// width function might not be handled by normal fmt.GetSize()
		width = fmt.Val_or_func(#width)
		fmtClass = false
		if hdc isnt false
			try
				{
				fmtClass = Global(fmt[0] $ 'Format')
				if fmtClass.Base?(TextFormat)
					{
					GetTextExtentPoint32(hdc, fmtClass.WidthChar, 1, sz = Object())
					aveCharWidth = sz.x
					}
				}
		}
	else
		try width = Global(fmt[0] $ 'Format').Val_or_func(#Width)
	return aveCharWidth * width
	}
