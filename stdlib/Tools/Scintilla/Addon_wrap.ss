// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	Init()
		{
		.SetWrapIndentMode(SC.WRAPINDENT_INDENT)
		}
	ContextMenu()
		{
		return #("Toggle Wrap")
		}
	On_Toggle_Wrap()
		{
		.SetWrap(.GetWrapMode() is SC.WRAP_NONE)
		}
	}