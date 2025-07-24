// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
class
	{
	x1: false
	y1: false
	New(.canvas)
		{
		}

	rectEl: false
	rect: false
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
		rect = .toRect(.x1, .y1, x, y)
		left = Max(rect.left, 0)
		top = Max(rect.top, 0)
		right = Min(rect.right, containerWidth)
		bottom = Min(rect.bottom, containerHeight)
		.rect = Object(:left, :top, :right, :bottom)
		w = right - left
		h = bottom - top
		if .rectEl is false
			{
			.rectEl = .canvas.Driver.AddRect(left, top, w, h, 1)
			.rectEl.SetAttribute('stroke-dasharray', '5,5')
			}
		else
			.canvas.Driver.ResizeRect(.rectEl, left, top, w, h)
		}

	toRect(x1, y1, x2, y2)
		{
		return Object(left: Min(x1, x2), top: Min(y1, y2),
			right: Max(x1, x2), bottom: Max(y1, y2))
		}

	MouseUp(x, y)
		{
		rtn = false
		if .rectEl isnt false
			{
			if .rect.right - .rect.left > 2 or .rect.bottom - .rect.top > 2
				rtn = .Rect(.rect.left, .rect.top, .rect.right, .rect.bottom)
			else
				rtn = .Point(.rect.left, .rect.top)
			.rectEl.Remove()
			.rect = .rectEl = false
			}
		else
			rtn = .Point(x, y)
		return rtn
		}

	Point(x, y)
		{
		return Object('Point', x, y)
		}

	Rect(x1, y1, x2, y2)
		{
		return Object('Rect', x1, y1, x2, y2)
		}

	varyx: 'none'
	varyy: 'none'
	ResizeDown(item, x, y)
		{
		r = item.BoundingRect()
		xmid = (r.x1 + r.x2) / 2
		ymid = (r.y1 + r.y2) / 2
		if item.InHandleArea(r.x1, r.y1, x, y)
			{ .varyx = 'left'; .varyy = 'top' }
		else if item.InHandleArea(r.x1, r.y2, x, y)
			{ .varyx = 'left'; .varyy = 'bottom' }
		else if item.InHandleArea(r.x2, r.y1, x, y)
			{ .varyx = 'right'; .varyy = 'top' }
		else if item.InHandleArea(r.x2, r.y2, x, y)
			{ .varyx = 'right'; .varyy = 'bottom' }
		else if item.InHandleArea(r.x1, ymid, x, y)
			{ .varyx = 'left'; .varyy = 'none' }
		else if item.InHandleArea(r.x2, ymid, x, y)
			{ .varyx = 'right'; .varyy = 'none' }
		else if item.InHandleArea(xmid, r.y1, x, y)
			{ .varyx = 'none'; .varyy = 'top' }
		else if item.InHandleArea(xmid, r.y2, x, y)
			{ .varyx = 'none'; .varyy = 'bottom' }
		.resize_rect = Object(left: r.x1, right: r.x2, top: r.y1, bottom: r.y2)
		item.PaintBoundingRect(r)
		}

	ResizeMove(item, x, y)
		{
		if not item.Method?(#DoResizeMove)
			{
			if .varyx isnt 'none'
				.resize_rect[.varyx] = x
			if .varyy isnt 'none'
				.resize_rect[.varyy] = y
			}
		else
			item.DoResizeMove(x, y, .varyx, .varyy, .resize_rect)
		.uninvert(.resize_rect)
		item.PaintBoundingRect(Object(x1: .resize_rect.left, y1: .resize_rect.top,
			x2: .resize_rect.right, y2: .resize_rect.bottom))
		}

	uninvert(r)
		{
		if r.left > r.right
			{
			r.Swap('left', 'right')
			.varyx = #(left: 'right', right: 'left', none: 'none')[.varyx]
			}
		if r.top > r.bottom
			{
			r.Swap('top', 'bottom')
			.varyy = #(top: 'bottom', bottom: 'top', none: 'none')[.varyy]
			}
		}

	ResizeUp(item, x /*unused*/, y /*unused*/)
		{
		item.RemoveBoundingRect()
		}
	}
