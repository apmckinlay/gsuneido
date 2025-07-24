// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
SuCanvasItem
	{
	relativeSizeDivisor: 3
	New(.x1, .y1, .x2, .y2, .width = false, .height = false)
		{
		.build()
		}

	build()
		{
		.rectEl = .Driver.AddRoundRect(.x1, .y1, .x2 - .x1, .y2 - .y1, .width, .height,
			fillColor: 'white')
		}

	BoundingRect()
		{
		return Object(x1: .x1, y1: .y1, x2: .x2, y2: .y2)
		}

	Resize(origx, origy, x, y)
		{
		changed? = false
		if .Resizing?(.x1, origx)
			{
			.x1 = x
			changed? = true
			}
		if .Resizing?(.y1, origy)
			{
			.y1 = y
			changed? = true
			}
		if .Resizing?(.x2, origx)
			{
			.x2 = x
			changed? = true
			}
		if .Resizing?(.y2, origy)
			{
			.y2 = y
			changed? = true
			}
		if changed? is true
			{
			.sortPoints(.x1, .y1, .x2, .y2)
			.width = Min(.x2 - .x1, .y2 - .y1) / .relativeSizeDivisor
			.height = .width
			.Driver.ResizeRoundRect(.rectEl, .x1, .y1, .x2 - .x1, .y2 - .y1,
				.width, .height)
			super.Resize(origx, origy, x, y)
			}
		}

	sortPoints(x1, y1, x2, y2)
		{
		.x1 = Min(x1, x2)
		.y1 = Min(y1, y2)
		.x2 = Max(x1, x2)
		.y2 = Max(y1, y2)
		}

	Move(dx, dy)
		{
		.x1 += dx
		.x2 += dx
		.y1 += dy
		.y2 += dy
		.Driver.MoveRoundRect(.rectEl, dx, dy)
		super.Move(dx, dy)
		}

	ResetSize(.x1, .y1, .x2, .y2, .width = 0, .height = 0)
		{
		.Driver.ResizeRoundRect(.rectEl, .x1, .y1, .x2 - .x1, .y2 - .y1, .width, .height)
		super.ResetSize()
		}

	GetElements()
		{
		return Object(.rectEl)
		}
	}
