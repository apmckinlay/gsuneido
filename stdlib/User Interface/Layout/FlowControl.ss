// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
/* e.g.
Window(#(Flow ((Static hello)(Field width: 25)(Static world)(Field width: 30)
	(Static good)(Field width: 35)(Static bye))))
*/
Container
	{
	Name: 'Flow'
	Xstretch: 1
	New(controls, skip = 8)
		{
		.ctrls = controls.Map(.Construct)
		.Xmin = .ctrls.Map({ it.Xmin }).Max()
		.Ymin = .ctrls.Map({ it.Ymin }).Max() // or top/bottom ???
		.skip = ScaleWithDpiFactor(skip)
		}

	// duplicate code with Group
	GetChildren()
		{
		return .ctrls
		}
	Tally()
		{
		return .ctrls.Size()
		}
	Append(control)
		{
		.Insert(.Tally(), control)
		}
	Insert(i, control)
		{
		.ctrls.Add(ctrl = .Construct(control), at: i)
		DoStartup(ctrl)
		// TODO recalc Xmin and Ymin
		.WindowRefresh()
		}
	Remove(i)
		{
		if not .ctrls.Member?(i)
			return
		.ctrls[i].Destroy()
		.ctrls.Delete(i)
		.WindowRefresh()
		}
	// end of duplicate code

	Resize(x, y, .w, h/*unused*/)
		{
		xorg = x
		yorg = y
		xw = x + w
		row = Object()
		top = bottom = 0
		for c in .ctrls
			{
			cw = c.Xmin
			if x + .skip + cw > xw
				{
				.outputRow(row, y, top)
				row = Object()
				y += top + bottom
				x = xorg
				top = bottom = 0
				}
			row.Add([ctrl: c, :x])
			x += .skip + cw
			top = Max(top, c.Top)
			bottom = Max(bottom, c.Ymin - c.Top)
			}
		if not row.Empty?()
			{
			.outputRow(row, y, top)
			y += top + bottom
			}

		ymin_before = .Ymin
		.Ymin = y - yorg
		if .Ymin isnt ymin_before
			.WindowRefresh()
		}
	outputRow(row, y, top)
		{
		for cc in row
			{
			c = cc.ctrl
			d = c.Top is 0 ? 0 : top - c.Top
			c.Resize(cc.x, y + d, c.Xmin, c.Ymin)
			}
		}
	}