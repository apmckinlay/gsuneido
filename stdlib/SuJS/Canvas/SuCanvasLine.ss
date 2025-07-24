// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
SuCanvasItem
	{
	New(.x1, .y1, .x2, .y2)
		{
		.build()
		}

	build()
		{
		.line = .Driver.AddLine(.x1, .y1, .x2, .y2, 1)
		}

	BoundingRect()
		{
		return Object(x1: Min(.x1, .x2), y1: Min(.y1, .y2),
			x2: Max(.x1, .x2), y2: Max(.y1, .y2))
		}

	ForeachHandle(block)
		{
		block(.x1, .y1)
		block(.x2, .y2)
		}

	Resize(origx, origy, x, y)
		{
		changed? = false
		if .Resizing?(.x1, origx) and .Resizing?(.y1, origy)
			{
			changed? = true
			.x1 = x
			.y1 = y
			}
		if .Resizing?(.x2, origx) and .Resizing?(.y2, origy)
			{
			changed? = true
			.x2 = x
			.y2 = y
			}

		if changed? is true
			{
			.Driver.ResizeLine(.line, .x1, .y1, .x2, .y2)
			super.Resize(origx, origy, x, y)
			}
		}

	Move(dx, dy)
		{
		.x1 += dx
		.y1 += dy
		.x2 += dx
		.y2 += dy
		.Driver.MoveLine(.line, dx, dy)
		super.Move(dx, dy)
		}

	ResetSize(.x1, .y1, .x2, .y2)
		{
		.Driver.ResizeLine(.line, .x1, .y1, .x2, .y2)
		super.ResetSize()
		}

	GetElements()
		{
		return Object(.line)
		}
	}
