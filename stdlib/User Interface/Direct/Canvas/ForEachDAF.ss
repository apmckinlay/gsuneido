// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
function (items, block)
	{
	for item in items
		{
		if item.Base?(CanvasDAF)
			block(item)
		else if item.Base?(CanvasGroup)
			ForEachDAF(item.GetItems(), block)
		}
	}