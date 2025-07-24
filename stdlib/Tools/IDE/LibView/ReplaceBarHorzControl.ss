// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
// ugly hack to make replace field the same size as the find field
// Grid would be simpler but doesn't handle stretch
HorzControl
	{
	Resize(x, y, w, h)
		{
		fc = .Window.FindControl('find').Field
		rc = .FindControl('replace')
		xtra = Max(0, w - .Xmin)
		content_stretch = .Group_content_stretch
		for (c in .Group_ctrls)
			{
			d = c.Top is 0 ? 0 : .Top - c.Top
			w = c.Xmin
			if content_stretch > 0
				w += (xtra * c.Xstretch / content_stretch).Round(0)
			if c is rc
				w = fc.GetRect().GetWidth()
			xtra -= w - c.Xmin
			content_stretch -= c.Xstretch
			c.Resize(x, y + d, w, c.Ystretch >= 0 ? h - d : c.Ymin)
			x += w
			}
		}
	}