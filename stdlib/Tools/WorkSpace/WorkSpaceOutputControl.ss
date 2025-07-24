// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
ScintillaIDEControl
	{
	Addon_show_references:,
	Addon_suneido_log:,
	Addon_inspect:,
	Addon_suneido_style_lines:,
	Addon_zoom: [
		zoom: false,
		zoom_ctrl: function (unused, text)
			{
			// zoom_ctrl arguments depend on zoomDialog interface in stdlib:EditorZoom
			zoomArgs = [
				set: text,
				Addon_suneido_log:,
				Addon_inspect:,
				Addon_suneido_style_lines:
			]
			return DisplayCodeControl(@zoomArgs)
			}
		]
	}
