// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
EditorControl
	{
	Bgndcolor: (250, 250, 50)
	New(@args)
		{
		super(@args)
		.bgndbrush = CreateSolidBrush(.color(.Bgndcolor))
		}
	color(color)
		{
		return Object?(color)
			? RGB(color[0], color[1], color[2])
			: color
		}
	CTLCOLOREDIT(wParam)
		{
		if .hilite?()
			{
			SetBkColor(wParam, .color(.Bgndcolor))
			return .bgndbrush
			}
		return GetStockObject(SO.WHITE_BRUSH)
		}
	hilite: false
	EN_CHANGE()
		{
		hilite = .hilite?()
		if hilite isnt .hilite
			{
			.Repaint()
			.hilite = hilite
			}
		return 0
		}
	hilite?()
		{
		return not GetWindowText(.Hwnd).Blank?()
		}
	Destroy()
		{
		DeleteObject(.bgndbrush)
		super.Destroy()
		}
	}
