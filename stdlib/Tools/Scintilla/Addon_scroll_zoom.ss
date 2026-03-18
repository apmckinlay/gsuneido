// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	SuJsSupport?: false
	ContextMenu()
		{
		return #("Reset Zoom")
		}

	On_Reset_Zoom()
		{
		.SETZOOM(0)
		}

	Scroll_Zoom()
		{
		return KeyPressed?(VK.CONTROL)
		}
	}