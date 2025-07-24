// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	Name:			'VirtualListThumb'
	Ystretch:		1

	thumbRect: 		false
	thumbHeight:	40

	New(width = false)
		{
		.CreateWindow("SuBtnfaceArrowNoDblClks", windowName: "", style: WS.VISIBLE)
		.SubClass()

		.createSystemObjects()

		.width = width isnt false ? width : GetSystemMetrics(SM.CXVSCROLL)
		.Xmin = .width
		.thumbHeight = ScaleWithDpiFactor(.thumbHeight)
		}

	createSystemObjects()
		{
		.scrollPen = CreatePen(PS.SOLID, 1, RGB(r:102, g:112, b:214))
		.scrollLine = CreatePen(PS.SOLID, 1, RGB(r:102, g:112, b:214))
		.scrollBrush = CreateSolidBrush(RGB(r:235, g:235, b:235))
		.backgroundBrush = GetSysColorBrush(COLOR.BTNFACE)
		}

	PAINT()
		{
		hdc = BeginPaint(.Hwnd, ps = Object())
		if .thumbRect isnt false
			.paint(hdc)
		EndPaint(.Hwnd, ps)
		return 0
		}

	paint(hdc)
		{
		FillRect(hdc, .GetClientRect().ToWindowsRect(), .backgroundBrush)
		pen = .mouseOnThumb ? .scrollLine : .scrollPen
		WithHdcSettings(hdc, [pen, .scrollBrush, SetBkMode: TRANSPARENT])
			{
			Rectangle(hdc, .thumbRect.left, .thumbRect.top, .thumbRect.right,
				.thumbRect.bottom)
			.drawDashesOnThumb(hdc)
			}
		}

	drawDashesOnThumb(hdc)
		{
		midPoint = .thumbRect.top + .thumbHeight / 2

		// each dash is made of two lines for extra thickness
		.drawLine(hdc, midPoint, adjustFromMid: 4)
		.drawLine(hdc, midPoint, adjustFromMid: 5)

		.drawLine(hdc, midPoint, adjustFromMid: 0)
		.drawLine(hdc, midPoint, adjustFromMid: 1)

		.drawLine(hdc, midPoint, adjustFromMid: -4)
		.drawLine(hdc, midPoint, adjustFromMid: -3)
		}

	drawLine(hdc, midPoint, adjustFromMid)
		{
		horzPadding = 5
		MoveTo(hdc, .thumbRect.left + horzPadding, midPoint + adjustFromMid)
		LineTo(hdc, .thumbRect.right - horzPadding, midPoint + adjustFromMid)
		}

	pos: 'top'
	Resize(x, y, w, h)
		{
		super.Resize(x, y, w, h)
		.calcThumbRect()
		}

	calcThumbRect()
		{
		height = .GetClientRect().GetHeight()
		top = .pos is 'top'
			? 0
			: .pos is 'bottom'
				? height - .thumbHeight
				: (height - .thumbHeight) / 2

		.thumbRect = Object(left: 0, right: .Xmin, :top, bottom: top + .thumbHeight)
		}

	SetThumbPosition(pos)	// takes top, bottom, middle
		{
		if pos is .pos
			return

		.pos = pos
		.calcThumbRect()
		.Repaint()
		}

	dragging: false
	dragy: 0
	lasty: 0
	pressed: false
	LBUTTONDOWN(lParam)
		{
		x = LOSWORD(lParam)
		y = HISWORD(lParam)

		SetCapture(.Hwnd)
		TrackMouseEvent(Object(cbSize: TRACKMOUSEEVENT.Size(),
			dwFlags: TME.HOVER, hwndTrack: .Hwnd))
		point = Object(:x, :y)
		if POINTinRECT(.thumbRect, point)
			{
			.dragging = true
			.pos = false
			.dragy = .lasty = y
			InvalidateRect(.Hwnd, .GetClientRect().ToWindowsRect(), false)
			return true
			}
		else if POINTinRECT(.upThumbRect(), point)
			{
			.Send('VirtualListThumb_PageUp')
			.pressed = -1
			}
		else if POINTinRECT(.downThumbRect(), point)
			{
			.Send('VirtualListThumb_PageDown')
			.pressed = 1
			}
		return false
		}

	upThumbRect()
		{
		rect = .thumbRect.Copy()
		rect.bottom = rect.top
		rect.top = 0
		return rect
		}

	downThumbRect()
		{
		winRect = .GetClientRect()
		rect = .thumbRect.Copy()
		rect.top = rect.bottom
		rect.bottom = winRect.GetHeight()
		return rect
		}

	MOUSEHOVER()
		{
		if .pressed isnt false
			.Send('VirtualListThumb_MouseHold', .pressed)
		return 0
		}

	mouseOnThumb: false
	MOUSEMOVE(lParam)
		{
		if .dragging
			{
			if .WindowActive?()
				.draggingThumb(HISWORD(lParam))
			else
				// stop if the window is no longer active, e.g. there is a pop-up
				.LBUTTONUP()
			}

		return 0
		}

	draggingThumb(y)
		{
		diff = .lasty - y
		if diff is 0
			return

		.thumbRect.top -= diff
		.thumbRect.bottom -= diff
		InvalidateRect(.Hwnd, .GetClientRect().ToWindowsRect(), false)
		.Update()
		.lasty = y

		.Send('VirtualListThumb_Dragging', y - .dragy)
		}

	LBUTTONUP()
		{
		ReleaseCapture()
		.handleDragging()
		.handlePressed()
		.Repaint()

		return 0
		}

	handlePressed()
		{
		if .pressed isnt false
			{
			.Send('VirtualListThumb_MouseUp')
			.pressed = false
			}
		}

	handleDragging()
		{
		if .dragging
			{
			.Send('VirtualListThumb_EndDragging')
			.dragging = false
			}
		}

	destroyObject(ob)
		{
		if ob isnt false
			DeleteObject(ob)
		}

	Destroy()
		{
		.destroyObject(.scrollBrush)
		.destroyObject(.scrollPen)
		.destroyObject(.backgroundBrush)
		.destroyObject(.scrollLine)
		super.Destroy()
		}
	}
