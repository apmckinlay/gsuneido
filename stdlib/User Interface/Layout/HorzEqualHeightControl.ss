// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
HorzControl
	{
	Name: 'HorzEqualHeight'
	Resize(x, y, w, h)
		{
		.Recalc()
		contentStretch = .GetContentStretch()
		xtra = Max(0, w - .Xmin)
		for c in .GetChildren()
			{
			cd = c.Top is 0 ? 0 : .Top - c.Top
			cw0 = cw = c.Xmin
			if c.Ystretch >= 0
				{
				ch = h - cd
				if c.Ymin > 0
					cw = (cw * ch / c.Ymin).Round(0)
				}
			else
				ch = c.Ymin
			if (contentStretch > 0)
				cw += (xtra * c.Xstretch / contentStretch).Round(0)
			xtra -= cw - cw0
			contentStretch -= c.Xstretch
			c.Resize(x, y + cd, cw, ch)
			x += cw
			}
		}
	}