// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
SuCanvasItem
	{
	New(.x1, .y1, .x2, .y2, .text, font = false, .justify = 'left', .rect = false)
		{
		.font = font is false ? Object() : font
		.build()
		}

	rectEl: false
	build()
		{
		.addRect()
		.lines = Object()
		w = .x2 - .x1
		y = .y1
		for line in .text.Split('\n')
			{
			textEl = .Driver.AddText(line, .x1, y, w, 0/*unused*/, .font, .justify)
			.lines.Add(textEl)
			metrics = SuRender().GetTextMetrics(textEl, line)
			.Driver.MoveText(textEl, 0, metrics.ascent) // align the top of the text to y
			textEl.style['user-select'] = 'none'
			y += metrics.height
			}
		}

	addRect()
		{
		if .rect is false
			return
		.rectEl = .Driver.AddRect(.x1, .y1, .x2 - .x1, .y2 - .y1, thick: 1,
			fillColor: .rect.fill, lineColor: .rect.line)
		}

	AfterEdit(.x1, .y1, .x2, .y2, .text, font = false, .justify = 'left', .rect = false)
		{
		.Remove()
		.font = font is false ? Object() : font
		.build()
		super.AfterEdit()
		}

	BoundingRect()
		{
		return Object(x1: .x1, y1: .y1, x2: .x2, y2: .y2)
		}

	Move(dx, dy)
		{
		.x1 += dx
		.y1 += dy
		.x2 += dx
		.y2 += dy
		for textEl in .lines
			.Driver.MoveText(textEl, dx, dy)
		if .rectEl isnt false
			.Driver.MoveRect(.rectEl, dx, dy)
		super.Move(dx, dy)
		}

	Resize(origx, origy, x, y)
		{
		if not .resize?(origx, origy)
			return

		x1 = .Resizing?(.x1, origx) ? x : .x1
		y1 = .Resizing?(.y1, origy) ? y : .y1
		x2 = .Resizing?(.x2, origx) ? x : .x2
		y2 = .Resizing?(.y2, origy) ? y : .y2

		.x1 = Min(x1, x2)
		.y1 = Min(y1, y2)
		.x2 = Max(x1, x2)
		.y2 = Max(y1, y2)

		if .rectEl isnt false
			.Driver.ResizeRect(.rectEl, .x1, .y1, .x2 - .x1, .y2 - .y1)

		.repaint()
		super.Resize(origx, origy, x, y)
		}

	resize?(origx, origy)
		{
		return .Resizing?(.x1, origx) or .Resizing?(.y1, origy) or
			.Resizing?(.x2, origx) or .Resizing?(.y2, origy)
		}

	ResetSize(.x1, .y1, .x2, .y2)
		{
		if .rectEl isnt false
			.Driver.ResizeRect(.rectEl, .x1, .y1, .x2 - .x1, .y2 - .y1)

		.repaint()
		super.ResetSize()
		}

	repaint()
		{
		w = .x2 - .x1
		y = .y1
		for textEl in .lines
			{
			.Driver.ResizeText(textEl, .x1, y, w, 0/*unused*/, .justify)
			metrics = SuRender().GetTextMetrics(textEl, textEl.textContent)
			.Driver.MoveText(textEl, 0, metrics.ascent) // align the top of the text to y
			y += metrics.height
			}
		}

	GetElements()
		{
		ob = .lines.Copy()
		if .rectEl isnt false
			ob.Add(.rectEl, at: 0)
		return ob
		}
	}
