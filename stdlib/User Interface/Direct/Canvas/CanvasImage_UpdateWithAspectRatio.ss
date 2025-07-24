// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(varyx, varyy, rect, aspectRatio)
		{
		if varyx is 'none'
			rect.right = rect.left +
				(rect.bottom - rect.top).Abs() * aspectRatio
		else if varyy is 'none'
			rect.bottom = rect.top +
				(rect.right - rect.left).Abs() / aspectRatio
		else
			.handleResizeBoth(varyx, varyy, rect, aspectRatio)
		}

	handleResizeBoth(varyx, varyy, rect, aspectRatio)
		{
		curWidth = (rect.right - rect.left).Abs()
		curHeight = (rect.bottom - rect.top).Abs()
		if curWidth / curHeight < aspectRatio
			{
			width = curHeight * aspectRatio
			varyx is 'left'
				? rect.left = rect.left < rect.right
					? rect.right - width
					: rect.right + width
				: rect.right = rect.left < rect.right
					? rect.left + width
					: rect.left - width
			}
		else
			{
			height = curWidth / aspectRatio
			varyy is 'top'
				? rect.top = rect.top < rect.bottom
					? rect.bottom - height
					: rect.bottom + height
				: rect.bottom = rect.top < rect.bottom
					? rect.top + height
					: rect.top - height
			}
		}
	}