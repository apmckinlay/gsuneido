// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
class
	{
	x1: false
	y1: false
	lineEl: false
	New(.canvas)
		{
		}

	MouseDown(x, y)
		{
		.x1 = x
		.y1 = y
		}

	MouseMove(x, y)
		{
		if .x1 is false or .y1 is false
			return

		containerWidth = .canvas.GetWidth()
		containerHeight = .canvas.GetHeight()
		.x2 = Max(Min(x, containerWidth), 0)
		.y2 = Max(Min(y, containerHeight), 0)
		if .lineEl is false
			{
			.lineEl = .canvas.Driver.AddLine(.x1, .y1, .x2, .y2, 1)
			.lineEl.SetAttribute('stroke-dasharray', '5,5')
			}
		else
			.canvas.Driver.ResizeLine(.lineEl, .x1, .y1, .x2, .y2)
		}

	MouseUp(x/*unused*/, y/*unused*/)
		{
		rtn = false
		if .lineEl isnt false
			{
			rtn = Object('Line', .x1, .y1, .x2, .y2)
			.lineEl.Remove()
			.lineEl = false
			}
		.x1 = .y1 = .x2 = .y2 = false
		return rtn
		}
	}
