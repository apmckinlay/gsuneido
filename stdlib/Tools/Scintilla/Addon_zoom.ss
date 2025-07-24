// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
// USAGE: Addon_zoom: [zoom: <true|false>, zoom_ctrl: <IE: ScintillaZoomControl>]
ScintillaAddon
	{
	ContextMenu()
		{
		return .Zoom_ctrl isnt false and .Zoom is false
			? #('Zoom...\tF6')
			: #()
		}

	On_Zoom()
		{
		EditorZoom(this, .Zoom, .Zoom_ctrl)
		}
	}
