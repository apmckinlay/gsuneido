// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// Contributions from Jos Schaars
// TODO: fix Top: 0 - it's half handled
Container
	{
	Name:		"Form"
	Left:		''
	Dir:		'horz'
	New(@args)
		{
		if args.GetDefault('left', '') isnt ''
			.Left = args.left
		adjust = .getAdjust(args)
		.xsep = ScaleWithDpiFactor(.xsep)
		.minline = ScaleWithDpiFactor(.minline)
		.ctrls = Object()
		for (item in args.Values(list:))
			if (item is 'nl')
				.ctrls.Add(adjust)
			else if (item is 'nL')
				.ctrls.Add(0)
			else // control
				{
				.ctrls.Add(.Construct(item))
				if (Object?(item) and item.Member?("group"))
					.ctrls.Last().Group = item.group
				}
		if (not .ctrls.Empty?())
			.Recalc()
		}
	getAdjust(args)
		{
		if Windows11?()
			return 0
		return args.GetDefault('overlap', true) ? -1 : 0
		}
	xsep: 10	// horz spacing between items
	minline: 6	// min line height
	Recalc()	// determine the positions
		{
		if .ctrls.Empty?()
			return
		result = .calcLineHeightsAndGroups()
		.alignVert(result.lineheights)
		.alignHorz(result.groups)
		}
	calcLineHeightsAndGroups()
		{
		// set x positions ignoring groups, determine line heights
		// check group numbers, track group members
		g = -1
		groups = Object()
		lineheights = Object()
		x = top = bottom = 0
		for (i in .ctrls.Members())
			if (Number?(c = .ctrls[i]))
				{
				lineheights.Add(Object(h: Max(.minline, top + bottom), :top))
				x = top = bottom = 0
				g = -1
				}
			else // control
				{
				c.X = x
				x += c.Xmin + .xsep
				top = Max(top, c.Top)
				bottom = Max(bottom, c.Ymin - c.Top)
				if c.Member?("Group") and Number?(c.Group)
					{
					Assert(c.Group > g)
					g = c.Group
					if (not groups.Member?(g))
						groups[g] = Object()
					groups[g].Add(i)
					}
				if .endofline?(i) and c.Xstretch > 0 and .Xstretch > 0
					c.formStretch? = true
				}
		lineheights.Add(Object(h: Max(.minline, top + bottom), :top))
		return Object(:groups, :lineheights)
		}
	alignVert(lineheights)
		{
		// align vertically using lineheights
		// determine form height
		ymin = line = 0
		for (c in .ctrls)
			{
			if (Number?(c))
				ymin += lineheights[line++].h + c
			else if (c.Top is 0)
				c.Y = ymin
			else
				c.Y = ymin + (lineheights[line].top - c.Top)
			}
		ymin += Number?(.ctrls.Last())? 1 : lineheights[line].h
		.Ymin = Max(.Ymin, ymin)	// override .Ymin if too small
		}
	alignHorz(groups)
		{
		// align horizontally using groups
		for (gn in groups.Members().Sort!())
			{
			group = groups[gn]
			// find the alignment for the group
			xl = 0	// x + left
			for (g in group.Values())
				xl = Max(xl, .ctrls[g].X + .ctrls[g].Left)
			if (.Left is "")
				.Left = xl
			.alignGroup(group, xl)
			}
		for (c in .GetChildren())
			.Xmin = Max(.Xmin, c.X + c.Xmin)	// override .Xmin if too small
		if .Left is ""
			.Left = 0
		}
	alignGroup(group, xl)
		{
		// align group and shove over controls
		for (g in group.Values())
			{
			c = .ctrls[g]
			if (0 isnt dx = xl - (c.X + c.Left))
				for (i = g; i < .ctrls.Size(); ++i)
					{
					if (Number?(c = .ctrls[i]))
						break
					c.X += dx
					}
			}
		}
	endofline?(i)
		{
		return i + 1 >= .ctrls.Size() or Number?(.ctrls[i + 1])
		}
	Resize(x, y, w, h/*unused*/)
		{
		for (c in .GetChildren())
			{
			cw = c.Xmin
			if c.GetDefault(#formStretch?, false)
				cw = Max(w - c.X, cw)
			c.Resize(x + c.X, y + c.Y, cw, c.Ymin)
			}
		}
	GetChildren()
		{
		return .ctrls.Filter(Instance?)
		}
	}
