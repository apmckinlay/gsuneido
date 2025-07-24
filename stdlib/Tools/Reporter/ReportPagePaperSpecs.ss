// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	sizes: (
		Letter: 	(w: 8.5, 	h: 11),
		Legal:		(w: 8.5, 	h: 14),
		Tabloid:	(w: 11, 	h: 17))

	Default: (page: 'Letter', orientation: 'Portrait', w: 8.5, 	h: 11)

	Options()
		{
		return .sizes.Members().Sort!()
		}

	GetPageSelection(width, height, defaultVal = '')
		{
		if false is size = .sizes.FindIf(
			{ it.w is width and it.h is height or it.h is width and it.w is height})
			return defaultVal
		return size
		}

	GetWidth(size, orientation)
		{
		sizeOb = Object?(size) ? size : .sizes[size]
		return orientation.Lower() is 'portrait' ? sizeOb.w : sizeOb.h
		}

	GetHeight(size, orientation)
		{
		sizeOb = Object?(size) ? size : .sizes[size]
		return orientation.Lower() is 'portrait' ? sizeOb.h : sizeOb.w
		}
	}