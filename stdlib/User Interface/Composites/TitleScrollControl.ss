// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	Name: 'TitleScroll'
	New(title, control)
		{
		super(Object('Vert',
			Object("CenterTitle", title)
			Object('Scroll',
				Object('Border', control, 5))
			))
		.Control = .Vert.Scroll.Border.Ctrl
		.ScrollHwnd = .Vert.Scroll.Hwnd
		}
	}