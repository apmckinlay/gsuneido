// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
SuCanvasItem
	{
	New(.x1, .y1, .x2, .y2, .color = 'white', .thick = 1)
		{
		.build()
		}

	build()
		{
		.rectEl = .Driver.AddRect(.x1, .y1, .x2 - .x1, .y2 - .y1, thick: .thick,
			fillColor: .color)
		}

	BoundingRect()
		{
		return Object(x1: .x1, y1: .y1, x2: .x2, y2: .y2)
		}

	Resize(origx, origy, x, y)
		{
		x1 = .Resizing?(.x1, origx) ? x : .x1
		y1 = .Resizing?(.y1, origy) ? y : .y1
		x2 = .Resizing?(.x2, origx) ? x : .x2
		y2 = .Resizing?(.y2, origy) ? y : .y2

		.x1 = Min(x1, x2)
		.y1 = Min(y1, y2)
		.x2 = Max(x1, x2)
		.y2 = Max(y1, y2)

		.sortPoints(.x1, .y1, .x2, .y2)
		.Driver.ResizeRect(.rectEl, .x1, .y1, .x2 - .x1, .y2 - .y1)
		super.Resize(origx, origy, x, y)
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
		.Driver.MoveRect(.rectEl, dx, dy)
		super.Move(dx, dy)
		}

	ResetSize(.x1, .y1, .x2, .y2)
		{
		.Driver.ResizeRect(.rectEl, .x1, .y1, .x2 - .x1, .y2 - .y1)
		super.ResetSize()
		}

	GetElements()
		{
		return Object(.rectEl)
		}
	}
