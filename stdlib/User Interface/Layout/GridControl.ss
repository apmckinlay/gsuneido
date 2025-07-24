// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// controls should be list of rows, each a list of controls
// BUG: columns with xstretch don't respect xmin
Container
	{
	Xskip: 3
	Yskip: 3
	Xstretch: 0
	Ystretch: 0
	Name: 'Grid'
	New(controls, setLeft = false, overlap = false)
		{
		.nrows = controls.Size()
		if overlap is true
			.Yskip = Windows11?() ? 0 : -1
		.Ymin = -.Yskip
		.ctrls = Object()
		.ymin = Object()
		.top = Object().Set_default(0)
		.span = Object().Set_default(Object().Set_default(1))
		.calcHeights(controls)

		ncols = 0
		for row in .ctrls.Values(list:)
			if row.Size() > ncols
				ncols = row.Size()

		.calcWidth(ncols)

		// if spanned cells are too big expand the last column of span
		.expandLastColumn(ncols)

		if setLeft is true
			.Left = .xmin[0]
		}

	calcHeights(controls)
		{
		rnum = 0
		for row in controls
			{
			col = ymin = top = bottom = 0
			r = Object()
			for ctrl in row
				{
				if .span?(ctrl)
					.span[rnum][col] = ctrl.span
				c = .Construct(ctrl)
				r.Add(c)
				if (c.Top > top)
					top = c.Top
				if (c.Ymin - c.Top > bottom)
					bottom = c.Ymin - c.Top
				ym = c.Top is 0 ? c.Ymin : top + bottom
				if ym > ymin
					ymin = ym
				++col
				}
			++rnum
			.ctrls.Add(r)
			.ymin.Add(ymin)
			.top.Add(top)
			.Ymin += ymin + .Yskip
			}
		}

	span?(ctrl)
		{
		return Object?(ctrl) and ctrl.Member?('span')
		}

	calcWidth(ncols)
		{
		.Xmin = -.Xskip
		.left = Object()
		.xmin = Object()
		for (i = 0; i < ncols; ++i)
			{
			xmin = left = right = 0
			for (j = 0; j < .nrows; ++j)
				{
				if i >= .ctrls[j].Size()
					continue
				ctrl = .ctrls[j][i]
				if ctrl.Left > left
					left = ctrl.Left
				if .span[j][i] is 1
					{
					if ctrl.Xmin - ctrl.Left > right
						right = ctrl.Xmin - ctrl.Left
					xmin = .ctrlXmin(ctrl, left, right, xmin)
					}
				else
					xmin = .spanning(ctrl, left, right, xmin)
				}
			.left.Add(left)
			.xmin.Add(xmin)
			.Xmin += xmin + .Xskip
			}
		}

	ctrlXmin(ctrl, left, right, xmin)
		{
		xm = ctrl.Left is 0 ? ctrl.Xmin : left + right
		if xm > xmin
			xmin = xm
		return xmin
		}

	spanning(ctrl, left, right, xmin)
		{
		if ctrl.Left isnt 0 and left + right > xmin
			xmin = left + right
		return xmin
		}

	expandLastColumn(ncols)
		{
		for (i = 0; i < ncols; ++i)
			{
			for (j = 0; j < .nrows; ++j)
				{
				if .span[j][i] is 1
					continue
				c = .ctrls[j][i]
				if c.Left is 0
					xmin = c.Xmin
				else
					{
					right = c.Xmin - c.Left
					xmin = .left[i] + right
					}

				total = -.Xskip
				for (k = 0; k < .span[j][i]; ++k)
					total += .xmin[i + k] + .Xskip

				if xmin > total
					{
					diff = xmin - total
					.xmin[i + .span[j][i] - 1] += diff
					.Xmin += diff
					}
				}
			}
		}

	Resize(x, y, w /*unused*/, h /*unused*/)
		{
		x0 = x
		yymin = 0
		for i in .ctrls.Members()
			{
			x = x0
			row = .ctrls[i]
			colnum = 0 // not same as j if spanning
			for j in row.Members()
				{
				ctrl = row[j]
				dx = ctrl.Left is 0 ? 0 : .left[j] - ctrl.Left
				wid = -.Xskip
				for (m = 0; m < .span[i][j]; ++m, ++colnum)
					wid += .xmin[colnum] + .Xskip
				ctrl.Resize(x + dx, .ctrlY(ctrl, y, i),
					ctrl.Xstretch is false ? ctrl.Xmin : wid - dx,
					ctrl.Ystretch is false ? ctrl.Ymin : .ymin[i])
				x += wid + .Xskip
				yymin = ctrl.Ymin > .ymin[i] ? ctrl.Ymin : .ymin[i]
				}
			y += yymin + .Yskip
			}
		}

	ctrlY(ctrl, y, i)
		{
		dy = ctrl.Top is 0 ? 0 : .top[i] - ctrl.Top
		if dy < 0
			dy = 0
		return y + dy
		}

	GetChildren()
		{
		list = Object()
		for row in .ctrls
			for ctrl in row
				list.Add(ctrl)
		return list
		}
	}
