// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(color, format = function (rgbOb) { return RGB(rgbOb.r, rgbOb.g, rgbOb.b) })
		{
		r, g, b = .rgbBase(color)
		return Object(
			bad: .rgbFade(r, g, b, fadingRate: 80)
			satisfactory: .rgbFade(r, g, b, fadingRate: 120)
			good: .rgbFade(r, g, b, fadingRate: 150)
			value: .rgbFade(r, g, b, fadingRate: 0)).
				Map!({ format(it) })
		}

	rgbBase(color)
		{
		// Extract single RGB components from color with a mask:
		r = (color & 0xFF0000) / 0x10000 	/*= RGB: red value */
		g = (color & 0x00FF00) / 0x100		/*= RGB: green value */
		b = (color & 0x0000FF)				/*= RGB: blue value */
		return r, g, b
		}

	rgbFade(r, g, b, fadingRate)
		{
		shade = Max(r, g, b)
		incr = (shade * fadingRate).PercentToDecimal().Int()
		r = Min(255 /*= max color*/, r + incr)
		g = Min(255 /*= max color*/, g + incr)
		b = Min(255 /*= max color*/, b + incr)
		return Object(:r, :g, :b)
		}
	}