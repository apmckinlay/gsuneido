// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(color)
		{
		switch
			{
		case String?(color):
			return color
		case Number?(color):
			return .toColorString(color)
		case Object?(color):
			return .toColorString(RGB(@color))
			}
		}

	toColorString(color)
		{
		s = color.Hex().LeftFill(6/*=digits*/, '0')
		return  '#' $ s[4/*=r*/::2] $ s[2/*=g*/::2] $ s[0/*=b*/::2]
		}

	Reverse(cssColor)
		{
		Assert(cssColor.Prefix?('#'))
		n = 0
		n = n * 16/*=hex base*/ + .fromHex(cssColor[5/*=idx*/])
		n = n * 16/*=hex base*/ + .fromHex(cssColor[6/*=idx*/])
		n = n * 16/*=hex base*/ + .fromHex(cssColor[3/*=idx*/])
		n = n * 16/*=hex base*/ + .fromHex(cssColor[4/*=idx*/])
		n = n * 16/*=hex base*/ + .fromHex(cssColor[1/*=idx*/])
		n = n * 16/*=hex base*/ + .fromHex(cssColor[2/*=idx*/])
		return n
		}

	fromHex(c)
		{
		if c >= '0' and c <= '9'
			return c.Asc() - '0'.Asc()
		if c >= 'a' and c <= 'f'
			return c.Asc() - 'a'.Asc() + 10/*=base*/
		if c >= 'A' and c <= 'F'
			return c.Asc() - 'A'.Asc() + 10/*=base*/
		}
	}