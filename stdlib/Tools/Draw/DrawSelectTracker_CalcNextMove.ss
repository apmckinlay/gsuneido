// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(dx, dy, rects, canvas)
		{
		boundingRect = .getBoundingRect(rects)
		boundingRect.left += dx
		boundingRect.right += dx
		boundingRect.top += dy
		boundingRect.bottom += dy
		adx = boundingRect.left < 0
			? Abs(boundingRect.left)
			: boundingRect.right > canvas.GetWidth()
				? canvas.GetWidth() - boundingRect.right
				: 0
		ady = boundingRect.top < 0
			? Abs(boundingRect.top)
			: boundingRect.bottom > canvas.GetHeight()
				? canvas.GetHeight() - boundingRect.bottom
				: 0
		return [x: dx + adx, y: dy + ady]
		}

	getBoundingRect(rects)
		{
		boundingRect = rects[0].Project(#left, #right, #top, #bottom)
		for rect in rects
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
	}