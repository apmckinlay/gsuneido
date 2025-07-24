// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
SuCanvasItem
	{
	New(.left, .top, .right, .bottom,
		.xStartArc = 0, .yStartArc = 0, .xEndArc = 0, .yEndArc = 0)
		{
		.build()
		}

	build()
		{
		.arcEl = .Driver.AddArc(.left, .top, .right, .bottom,
			.xStartArc, .yStartArc, .xEndArc, .yEndArc)
		}

	BoundingRect()
		{
		return Object(x1: .left, y1: .top, x2: .right, y2: .bottom)
		}

	Resize(origx, origy, x, y)
		{
		changed? = false
		if .Resizing?(.left, origx)
			{
			.left = x
			changed? = true
			}
		if .Resizing?(.top, origy)
			{
			.top = y
			changed? = true
			}
		if .Resizing?(.right, origx)
			{
			.right = x
			changed? = true
			}
		if .Resizing?(.bottom, origy)
			{
			.bottom = y
			changed? = true
			}
		if changed? is true
			{
			.sortPoints(.left, .top, .right, .bottom)
			.xStartArc = .left
			.xEndArc = .right
			.yStartArc = .top
			.yEndArc = .bottom
			.Driver.ResizeArc(.arcEl, .left, .top, .right, .bottom,
				.xStartArc, .yStartArc, .xEndArc, .yEndArc)
			super.Resize(origx, origy, x, y)
			}
		}

	sortPoints(left, top, right, bottom)
		{
		if .left > .right
			{
			.left = Min(left, right)
			.top = Min(top, bottom)
			.right = Max(left, right)
			.bottom = Max(top, bottom)
			}
		}

	Move(dx, dy)
		{
		.left += dx
		.right += dx
		.xStartArc += dx
		.xEndArc += dx
		.top += dy
		.bottom += dy
		.yStartArc += dy
		.yEndArc += dy
		.Driver.MoveArc(.arcEl, dx, dy)
		super.Move(dx, dy)
		}

	ResetSize(.left, .top, .right, .bottom,
		.xStartArc = 0, .yStartArc = 0, .xEndArc = 0, .yEndArc = 0)
		{
		.sortPoints(.left, .top, .right, .bottom)
		.Driver.ResizeArc(.arcEl, .left, .top, .right, .bottom,
			.xStartArc, .yStartArc, .xEndArc, .yEndArc)
		super.ResetSize()
		}

	GetElements()
		{
		return Object(.arcEl)
		}
	}
