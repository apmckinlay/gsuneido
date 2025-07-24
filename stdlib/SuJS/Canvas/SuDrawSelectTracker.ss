// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
SuDrawRectTracker
	{
	New(.canvas)
		{
		super(canvas)
		.rects = Object()
		}

	mousedown?: false
	MouseDown(x, y, _event)
		{
		.rects = Object()
		if false isnt i = .canvas.ItemAtPoint(x, y)
			{
			item = .canvas.GetAllItems()[i]
			if not .canvas.GetSelected().Has?(item)
				{
				.canvas.MaybeClearSelect()
				.canvas.Select(i)
				}
			else
				{
				if event.ctrlKey is true or event.shiftKey is true
					.canvas.UnSelect(i)
				}

			for item in .canvas.GetSelected()
				{
				r = item.BoundingRect()
				item.PaintBoundingRect(r)
				rect = Object(left: r.x1, right: r.x2, top: r.y1, bottom: r.y2, :item)
				.rects.Add(rect)
				}
			.x = .origx = x
			.y = .origy = y
			}
		else
			super.MouseDown(x, y)
		}

	dragging()
		{
		return .rects.NotEmpty?()
		}

	MouseMove(x, y)
		{
		if .dragging()
			{
			boundingRect = .getBoundingRect()
			dx = x - .x
			dy = y - .y
			boundingRect.left += dx
			boundingRect.right += dx
			boundingRect.top += dy
			boundingRect.bottom += dy

			adx = boundingRect.left < 0
				? Abs(boundingRect.left)
				: boundingRect.right > .canvas.GetWidth()
					? .canvas.GetWidth() - boundingRect.right
					: 0
			ady = boundingRect.top < 0
				? Abs(boundingRect.top)
				: boundingRect.bottom > .canvas.GetHeight()
					? .canvas.GetHeight() - boundingRect.bottom
					: 0

			for rect in .rects
				{
				rect.left += dx + adx
				rect.right += dx + adx
				rect.top += dy + ady
				rect.bottom += dy + ady
				rect.item.PaintBoundingRect(Object(
					x1: rect.left, y1: rect.top, x2: rect.right, y2: rect.bottom))
				}
			.x = x + adx
			.y = y + ady
			}
		else
			super.MouseMove(x, y)
		}

	getBoundingRect()
		{
		boundingRect = .rects[0].Project(#left, #right, #top, #bottom)
		for rect in .rects
			{
			if boundingRect.left > rect.left
				boundingRect.left = rect.left
			if boundingRect.right < rect.right
				boundingRect.right = rect.right
			if boundingRect.top > rect.top
				boundingRect.top = rect.top
			if boundingRect.bottom < rect.bottom
				boundingRect.bottom = rect.bottom
			}
		return boundingRect
		}

	MouseUp(x, y)
		{
		if .dragging()
			{
			dx = .x - .origx
			dy = .y - .origy
			move? = dx isnt 0 or dy isnt 0
			for rect in .rects
				{
				rect.item.RemoveBoundingRect()
				if move?
					rect.item.Move(dx, dy)
				}
			return false
			}
		else
			return super.MouseUp(x, y)
		}

	AddPoint(x, y)
		{
		selected = .canvas.GetSelected()
		for (item in selected)
			if (item.Contains(x, y))
				return
		.canvas.SelectPoint(x, y)
		for (item in selected)
			.canvas.Select(.canvas.GetAllItems().Find(item))
		}

	Point(x, y, _event)
		{
		if event.ctrlKey is true or event.shiftKey is true
			.AddPoint(x, y)
		else
			.canvas.SelectPoint(x, y)
		return false
		}

	Rect(x1, y1, x2, y2)
		{
		.canvas.SelectRect(Min(x1, x2), Min(y1, y2), Max(x1, x2), Max(y1, y2))
		return false
		}
	}
