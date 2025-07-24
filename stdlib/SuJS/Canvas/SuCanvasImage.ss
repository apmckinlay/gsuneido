// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
SuCanvasItem
	{
	New(.x1, .y1, .x2, .y2, .image, .aspectRatio = 1)
		{
		.build()
		}

	build()
		{
		.imageEl = .Driver.AddImage(.x1, .y1, .x2 - .x1, .y2 - .y1,
			'data:image/jpeg;base64,' $ .image)
		.imageEl.SetAttribute('preserveAspectRatio', 'none')
		}

	BoundingRect()
		{
		return Object(x1: .x1, y1: .y1, x2: .x2, y2: .y2)
		}

	Resize(origx, origy, x, y)
		{
		changed? = false
		varyx = varyy = 'none'
		if .Resizing?(.x1, origx)
			{
			.x1 = x
			varyx = 'left'
			changed? = true
			}
		if .Resizing?(.y1, origy)
			{
			.y1 = y
			varyy = 'top'
			changed? = true
			}
		if .Resizing?(.x2, origx)
			{
			.x2 = x
			varyx = 'right'
			changed? = true
			}
		if .Resizing?(.y2, origy)
			{
			.y2 = y
			varyy = 'bottom'
			changed? = true
			}
		if changed? is true
			{
			rect = Object(left: .x1, right: .x2, top: .y1, bottom: .y2)
			CanvasImage_UpdateWithAspectRatio(varyx, varyy, rect, .aspectRatio)
			.x1 = rect.left; .x2 = rect.right; .y1 = rect.top; .y2 = rect.bottom

			.sortPoints(.x1, .y1, .x2, .y2)
			.Driver.ResizeImage(.imageEl, .x1, .y1, .x2 - .x1, .y2 - .y1)
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
		.y1 += dy
		.x2 += dx
		.y2 += dy
		.Driver.MoveImage(.imageEl, dx, dy)
		super.Move(dx, dy)
		}

	ResetSize(.x1, .y1, .x2, .y2)
		{
		.Driver.ResizeImage(.imageEl, .x1, .y1, .x2 - .x1, .y2 - .y1)
		super.ResetSize()
		}

	DoResizeMove(x, y, varyx, varyy, rect)
		{
		if varyx isnt 'none'
			rect[varyx] = x
		if varyy isnt 'none'
			rect[varyy] = y
		CanvasImage_UpdateWithAspectRatio(varyx, varyy, rect, .aspectRatio)
		}

	GetElements()
		{
		return Object(.imageEl)
		}
	}
