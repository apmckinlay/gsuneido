// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
class
	{
	snapLayouts: #()
	ContextMenu: #()
	New()
		{
		.snapLayouts = Object(SnapLeftRight())

		.ContextMenu = Object()
		for layout in .snapLayouts
			.ContextMenu.Add(@layout.ContextMenu)
		}

	ContextCall(args, window)
		{
		for layout in .snapLayouts
			layout.ContextCall(args, window)
		}

	HorzResize(newWidth, window)
		{
		for layout in .snapLayouts
			if layout.HorzResize(newWidth, window) is true
				return true
		return false
		}

	Remove(window)
		{
		for layout in .snapLayouts
			layout.Remove(window)
		}
	}
