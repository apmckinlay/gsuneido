// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(container, left, top, rcExclude, rect)
		{
		if left is false or top is false
			return
		windowRect = SuRender.GetClientRect()
		menuRect = SuRender.GetClientRect(container)
		adjustStyle = .getAdjustStyle(rcExclude)
		x = left
		if x + menuRect.width > windowRect.right
			{
			if adjustStyle isnt .adjustStyle.vert
				x = Max(windowRect.right - menuRect.width, 0)
			else
				x = Max(rect.left - menuRect.width, 0)
			}
		y = top
		if y + menuRect.height > windowRect.bottom
			{
			if adjustStyle isnt .adjustStyle.horz
				y = Max(windowRect.bottom - menuRect.height, 0)
			else
				y = Max(rect.top - menuRect.height, 0)
			}
		container.SetStyle(#left, x $ 'px')
		container.SetStyle(#top, y $ 'px')
		}

	adjustStyle: (any: 0, vert: 1, horz: 2)
	getAdjustStyle(rcExclude)
		{
		inf = 9999
		switch
			{
			case rcExclude is 0:
				return .adjustStyle.any
			case rcExclude.left is -inf and rcExclude.right is inf:
				return .adjustStyle.horz
			case rcExclude.top is -inf and rcExclude.bottom is inf:
				return .adjustStyle.vert
			}
		}
	}
