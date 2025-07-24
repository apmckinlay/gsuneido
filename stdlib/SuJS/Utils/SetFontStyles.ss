// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	map: (
		italic: 'font-style',
		weight: 'font-weight',
		fontSize: 'font-size',
		fontName: 'font-family')

	CallClass(font, styles)
		{
		if font is false
			return
		for m in .map.Members()
			{
			if font[m] isnt ''
				styles[.map[m]] = font[m]
			}
		}
	}
