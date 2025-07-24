// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
HtmlTextControl
	{
	Xstretch: false
	WndClass: 'SuNobgndHand'
	Underline: true

	New(text, href)
		{
		super(text)
		.Flags = DT.SINGLELINE
		.CalcSize()
		.href = href
		.underMouse? = false
		}
	Resize(x, y, w, h /*unused*/)
		{
		r = .CalcRect(w)
		.Ymin = r.bottom + 6 /*= extra space at bottom*/
		super.Resize(x, y, w, .Ymin)
		}
	PreDraw(dc)
		{
		if (.underMouse? isnt true)
			SetTextColor(dc, CLR.BLUE)
		else
			SetTextColor(dc, GetSysColor(COLOR.HIGHLIGHT))
		}
	LBUTTONUP()
		{
		.Send('Goto', .href)
		return 0
		}
	MOUSEMOVE()
		{
		if (not .underMouse?)
			{
			.underMouse? = true
			InvalidateRect(.Hwnd, NULL, true)
			TrackMouseEvent(Object(cbSize: TRACKMOUSEEVENT.Size(),
				dwFlags: TME.LEAVE hwndTrack: .Hwnd))
			}
		return 0
		}
	MOUSELEAVE()
		{
		if (.underMouse?)
			{
			.underMouse? = false
			InvalidateRect(.Hwnd, NULL, true)
			}
		return 0
		}
	}
