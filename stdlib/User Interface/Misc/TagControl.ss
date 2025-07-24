// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
// displays a single tag, used by TagsControl
WndProc
	{
	Name: 'Tag'
	x_margin: 5
	y_margin: 3
	New(text, color = 0xcccccc)
		{
		.CreateWindow("SuBtnfaceArrow", "", WS.VISIBLE)
		.SubClass()
		.SetFont(size: '-1', :text)
		.Top += .y_margin
		.xtra = 2 * .TextExtent('x').x
		.text = text
		.Xmin += 2 * .x_margin + .xtra + .x_margin
		.Ymin += 2 * .y_margin
		if Object?(color)
		   color = RGB(@color)
		.color = color
		}

	Get()
		{
		return .text
		}

	ERASEBKGND()
		{
		return 1
		}
	PAINT()
		{
		hdc = BeginPaint(.Hwnd, ps = Object())
		GetClientRect(.Hwnd, r = Object())
		WithHdcSettings(hdc, .hdcSettings())
			{
			width = height = 10
			RoundRect(hdc, r.left, r.top, r.right, r.bottom, width, height)
			x = r.left + .x_margin
			y = r.top + .y_margin
			TextOut(hdc, x, y - 3 /*= textYOffset*/, 'x', 1)
			x += .xtra
			SetTextColor(hdc, 0)
			TextOut(hdc, x, y, .text, .text.Size())
			}
		EndPaint(.Hwnd, ps)
		.close_width = x
		return 0
		}

	hdcSettings()
		{
		return Object(
			GetStockObject(SO.NULL_PEN),
			.GetFont(),
			brush: .color,
			SetBkMode: TRANSPARENT
			SetTextColor: .mouseover?
				? CLR.RED
				: CLR.EnhancedButtonFace)
		}

	mouseover?: false
	MOUSEMOVE()
		{
		if .mouseover?
			return 0
		TrackMouseEvent(Object(cbSize: TRACKMOUSEEVENT.Size(),
			dwFlags: TME.LEAVE hwndTrack: .Hwnd))
		.mouseover? = true
		.Repaint()
		return 0
		}
	MOUSELEAVE()
		{
		.mouseover? = false
		.Repaint()
		return 0
		}

	LBUTTONDOWN(lParam)
		{
		x = LOWORD(lParam)
		if x < .close_width
			.Send("Tag_Remove", .text)
		return 0
		}
	}
