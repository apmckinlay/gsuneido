// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// Contributions from Jos Schaars
// internal base for Vert, Horz, and Split, not for users
Container
	{
	Xstretch: ""
	Ystretch: ""
	New(controls)
		{
		.ctrls = Object()
		for (c in controls.Values(list:))
			.ctrls.Add(.Construct(c))

		.xmin0 = .Xmin
		.set_xstretch? = .Xstretch is ""
		.ymin0 = .Ymin
		.set_ystretch? = .Ystretch is ""
		.overlap = false isnt controls.GetDefault('overlap', false) ? 1 : 0
		.Recalc()
		}
	CalcXminByControls(@args)
		{
		.Xmin = .xmin0 = .DoCalcXminByControls(@args)
		}
	Recalc()
		{
		.Xmin = .xmin0
		.Ymin = .ymin0
		if .Dir is "vert"
			.vertRecalc()
		else
			.horzRecalc()
		}
	vertRecalc()
		{
		xmin = xstretch = ymin = ystretch = left = right = false
		maxheight = 0
		for (c in .ctrls)
			{
			xstretch = Max(xstretch, c.Xstretch)
			if (c.Left is 0)
				xmin = Max(xmin, c.Xmin)
			else
				{
				left = Max(left, c.Left)
				right = Max(right, c.Xmin - c.Left)
				}
			ymin += c.Ymin - .overlap
			ystretch += c.Ystretch
			maxheight += c.CalcMaxHeight()
			}
		xmin = Max(xmin, left + right)
		if (xmin > .xmin0)
			.Xmin = xmin
		.Left = left
		ymin += .overlap
		if (ymin >= .ymin0)
			.Ymin = ymin
		.content_stretch = ystretch
		.MaxHeight = maxheight
		.setStretch(xstretch, ystretch)
		}
	horzRecalc()
		{
		xmin = xstretch = ymin = ystretch = top = bottom = false
		for (c in .ctrls)
			{
			ystretch = Max(ystretch, c.Ystretch)
			if (c.Top is 0)
				ymin = Max(ymin, c.Ymin)
			else
				{
				top = Max(top, c.Top)
				bottom = Max(bottom, c.Ymin - c.Top)
				}
			xmin += c.Xmin - .overlap
			xstretch += c.Xstretch
			}
		ymin = Max(ymin, top + bottom)
		if (ymin > .ymin0)
			.Ymin = ymin
		.Top = top
		xmin += .overlap
		if (xmin > .xmin0)
			.Xmin = xmin
		.content_stretch = xstretch
		.setStretch(xstretch, ystretch)
		}
	setStretch(xstretch, ystretch)
		{
		if .set_xstretch?
			.Xstretch = xstretch
		if .set_ystretch?
			.Ystretch = ystretch
		}
	Append(control)
		{
		.Insert(.Tally(), control)
		}
	Insert(i, control)
		{
		c = .addCtrl(control, i)
		.WindowRefresh()
		return c
		}
	AppendAll(controls)
		{
		for ctrl in controls
			.addCtrl(ctrl, .Tally())
		.WindowRefresh()
		}
	addCtrl(ctrl, at)
		{
		.ctrls.Add(c = .Construct(ctrl), :at)
		DoStartup(c)
		return c
		}
	Remove(i)
		{
		if not .ctrls.Member?(i)
			return
		.ctrls[i].Destroy()
		.ctrls.Delete(i)
		.WindowRefresh()
		}
	RemoveAll()
		{
		for (c in .ctrls)
			c.Destroy()
		.ctrls = Object()
		.WindowRefresh()
		}
	x: false
	Resize(x, y, w, h)
		{
		.Recalc() // not sure why this is necessary, but it seems to be
		.x = x; .y = y; .w = w; .h = h; // used by BookSplit & Splitter
		if .Dir is "vert"
			.vertResize(x, y, w, h)
		else
			.horzResize(x, y, w, h)
		}
	vertResize(x, y, w, h)
		{
		xtra = Max(0, h - .Ymin)
		for c in .ctrls
			{
			d = c.Left is 0 ? 0 : .Left - c.Left
			h = c.Ymin
			if .content_stretch > 0
				{
				h += (xtra * c.Ystretch / .content_stretch).Round(0)
				if h > c.MaxHeight
					h = c.MaxHeight
				}
			xtra -= h - c.Ymin
			.content_stretch -= c.Ystretch
			c.Resize(x + d, y, c.Xstretch >= 0 ? w - d : c.Xmin, h)
			y += h - .overlap
			}
		}
	horzResize(x, y, w, h)
		{
		xtra = Max(0, w - .Xmin)
		for c in .ctrls
			{
			d = c.Top is 0 ? 0 : .Top - c.Top
			w = c.Xmin
			if .content_stretch > 0
				w += (xtra * c.Xstretch / .content_stretch).Round(0)
			xtra -= w - c.Xmin
			.content_stretch -= c.Xstretch
			c.Resize(x, y + d, w, c.Ystretch >= 0 ? h - d : c.Ymin)
			x += w - .overlap
			}
		}
	Tally()
		{ return .ctrls.Size() }
	ctrls: () // default if destroyed
	GetChildren()
		{ return .ctrls }
	GetContentStretch()
		{ return .content_stretch }
	Get()
		{
		ob = Object()
		for (ctrl in .ctrls)
			if ctrl.Method?(#Get)
				ob[ctrl.Name] = ctrl.Get()
		return ob
		}
	}