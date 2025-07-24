// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
function (key, zoomArgs = #(), pressed = false)
	{
	if key is VK.F6
		EditorZoom(@zoomArgs)
	else if KeyPressed?(VK.CONTROL, :pressed)
		{
		if key is VK.F
			.On_Find()
		else if key is VK.P
			.On_Print()
		else if key is VK.A
			.On_Select_All()
		else
			return 'callsuper'
		}
	else if key is VK.F3
		KeyPressed?(VK.SHIFT, :pressed) ? .On_Find_Previous() : .On_Find_Next()
	else
		return 'callsuper'
	return 0
	}
