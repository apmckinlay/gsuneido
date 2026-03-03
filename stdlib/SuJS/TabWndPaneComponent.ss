// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
WndPaneComponent
	{
	New(@args)
		{
		super(@args)
		.SetStyles(Object('position': 'absolute',
			'width': '100%',
			'height': '100%'))
		}

	SetVisible(visible)
		{
		visible = not .GetHidden() and visible
		.El.SetStyle("visibility", visible ? '' : 'hidden')
		}
	}