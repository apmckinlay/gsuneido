// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	Xstretch: 	1
	Ystretch: 	1
	Xmin: 		50
	Ymin: 		50
	Name: 		"Scroll"
	outerH:		0

	New(control, style = 0, .vmanual = false, wndclass = "SuBtnfaceArrow",
		.vdisable = false, .dyscroll = 21, .trim = false, .noEdge = false)
		{
		exStyle = WS_EX.CONTROLPARENT
		if not noEdge
			exStyle |= WS_EX.CLIENTEDGE
		.CreateWindow(wndclass, "",
			WS.HSCROLL | WS.VSCROLL | WS.VISIBLE | style, :exStyle, w: 1, h: 1)
		.SubClass()
		SetProp(.Hwnd, "suneido_navigable", true)
		.ctrl = .Construct(control)
		.handleTrim()
		}
	handleTrim()
		{
		if .trim is false
			return
		border = 4
		w = .ctrl.Xmin + border + GetSystemMetrics(SM.CXVSCROLL)
		if w < .Xmin
			.Xmin = w
		h = .ctrl.Ymin + border + GetSystemMetrics(SM.CXVSCROLL)
		if h < .Ymin
			.Ymin = h
		if .trim is 'auto' and .ctrl.Ystretch > 0
			.MaxHeight = 99999
		else
			.MaxHeight = .ctrl.Ymin + border + GetSystemMetrics(SM.CXVSCROLL)
		}
	GetChild()
		{
		return .ctrl
		}
	GetChildren()
		{
		return Object(.ctrl)
		}
	Recalc()
		{
		.handleTrim()
		}
	outerW: false
	Adjust() // call if content size changes
		{
		if (.outerW is false)
			return
		.adjustHorz()
		if (.vmanual)
			.vertBar = .Send("AdjustVert")
		else
			.AdjustVert()

		if not .noEdge
			{
			brd3D = (GetWindowLong(.Hwnd, GWL.EXSTYLE) & WS_EX.CLIENTEDGE) isnt 0
			if (brd3D isnt (.horzBar or .vertBar))
				{
				// only show client edge border if scrollbars are visible
				SetWindowLong(.Hwnd, GWL.EXSTYLE, WS_EX.CONTROLPARENT |
					(brd3D ? 0 : WS_EX.CLIENTEDGE))
				SetWindowPos(.Hwnd, 0, 0, 0, 0, 0,
					SWP.FRAMECHANGED | SWP.NOSIZE | SWP.NOMOVE | SWP.NOZORDER)
				}
			}
		}
	adjustHorz()
		{
		inner = .getInner()
		if (.ctrl.Xmin > inner.w)
			inner.h -= GetSystemMetrics(SM.CYHSCROLL)
		if (.ctrl.Ymin > inner.h or .vmanual)
			inner.w -= GetSystemMetrics(SM.CXVSCROLL)

		SetScrollInfo(.Hwnd, SB.HORZ,
			Object(fMask: SIF.RANGE | SIF.PAGE, nMin: 0, nMax: .ctrl.Xmin,
			nPage: inner.w + 1, cbSize: SCROLLINFO.Size()), true)
		.horzBar = .ctrl.Xmin > inner.w

		.maxXscroll = Max(0, .ctrl.Xmin - inner.w)
		.xscroll = Min(.xscroll, .maxXscroll)
		}
	AdjustVert()
		{
		inner = .getInner()
		if (.ctrl.Ymin > inner.h or .vdisable is true)
			inner.w -= GetSystemMetrics(SM.CXVSCROLL)
		if (.ctrl.Xmin > inner.w)
			inner.h -= GetSystemMetrics(SM.CYHSCROLL)

		mask = SIF.RANGE | SIF.PAGE
		if (.vdisable is true)
			mask |= SIF.DISABLENOSCROLL
		// according to MSDN documentation, the range for nPage is
		// 0 to nMax - nMin + 1, this is why we add 1 to h
		SetScrollInfo(.Hwnd, SB.VERT,
			Object(fMask: mask, nMin: 0, nMax: .ctrl.Ymin,
			nPage: inner.h + 1, cbSize: SCROLLINFO.Size()), true)
		.vertBar = .ctrl.Ymin > inner.h

		.maxYscroll = Max(0, (.ctrl.Ymin - inner.h).Int())
		.yscroll = Min(.yscroll, .maxYscroll)
		}
	getInner()
		{
		brdr = .getBorder()
		return Object(w: .outerW - brdr, h: .outerH - brdr)
		}
	getBorder()
		{
		brdr = (GetWindowLong(.Hwnd, GWL.STYLE) & WS.BORDER) isnt 0 ? 1 : 0
		brdr += (GetWindowLong(.Hwnd, GWL.EXSTYLE) & WS_EX.CLIENTEDGE) isnt 0 ? 2 : 0
		return 2 * brdr
		}
	Resize(x, y, w, h)
		{
		super.Resize(x, y, w, h)
		.outerW = w
		.outerH = h

		.resize()
		// to look at some time, only needed for Tab ?
		.resize()
		}
	resize()
		{
		.Adjust()
		inner = .getInner()
		inner.w -= .vertBar ? GetSystemMetrics(SM.CXVSCROLL) : 0
		inner.h -= .horzBar ? GetSystemMetrics(SM.CXVSCROLL) : 0
		inner.w = .ctrl.Xstretch >= 0 ? Max(.ctrl.Xmin, inner.w) : .ctrl.Xmin
		inner.h = .ctrl.Ystretch >= 0 ? Max(.ctrl.Ymin, inner.h) : .ctrl.Ymin
		.ctrl.Resize(-.xscroll, -.yscroll, inner.w, inner.h)
		}
	MOUSEWHEEL(wParam)
		{
		scroll = GetWheelScrollInfo(wParam)
		if KeyPressed?(VK.CONTROL)
			{
			if true is .Send('Zoom', scroll.clicks > 0 ? 'In' : 'Out')
				return 0
			}
		dy = scroll.lines * .dyscroll
		if (scroll.page? or dy.Abs() > .getInner().h)
			dy = scroll.clicks.Sign() * .getInner().h
		if false is .Scroll(0, dy)
			.WndProc.Callsuper(.Hwnd, WM.MOUSEWHEEL, :wParam, lParam: 0)
		return 0
		}
	yscroll: 0
	VSCROLL(wParam)
		{
		if .vmanual
			.manualVScroll(wParam)
		else
			switch (LOWORD(wParam))
				{
			case SB.LINEUP		:	.Scroll(0, .dyscroll)
			case SB.LINEDOWN	:	.Scroll(0, -.dyscroll)
			case SB.PAGEUP		:	.Scroll(0, .getInner().h - .dyscroll)
			case SB.PAGEDOWN	:	.Scroll(0, .dyscroll - .getInner().h)
			case SB.TOP			:	.Scroll(0, .maxYscroll)
			case SB.BOTTOM		:	.Scroll(0, -.maxYscroll)
			case SB.THUMBTRACK	:	.Scroll(0, .yscroll - HIWORD(wParam))
			default:
				}
		return 0
		}
	manualVScroll(wParam)
		{
		switch (LOWORD(wParam))
			{
		case SB.LINEUP		:	.Send('ScrollUp')
		case SB.LINEDOWN	:	.Send('ScrollDown')
		case SB.PAGEUP		:	.Send('PageUp')
		case SB.PAGEDOWN	:	.Send('PageDown')
		case SB.THUMBTRACK	:
								pos = HIWORD(wParam)
								.Send('ScrollVertPos', pos)
								.ScrollVertPos(pos)
		default:
			}
		}
	ScrollVertPos(yscroll)
		{
		SetScrollInfo(.Hwnd, SB.VERT,
			Object(fMask: SIF.POS, nPos: yscroll, cbSize: SCROLLINFO.Size()), true)
		}
	xscroll: 0
	GetXscroll()
		{
		return .xscroll
		}
	GetYscroll()
		{
		return .yscroll
		}
	HSCROLL(wParam)
		{
		line = 30
		page = .getInner().w - line
		switch (LOWORD(wParam))
			{
		case SB.LINELEFT	:	.Scroll(line, 0)
		case SB.LINERIGHT	: 	.Scroll(-line, 0)
		case SB.PAGELEFT	:  	.Scroll(page, 0)
		case SB.PAGERIGHT	: 	.Scroll(-page, 0)
		case SB.TOP			:	.Scroll(.maxXscroll, 0)
		case SB.BOTTOM		:	.Scroll(-.maxXscroll, 0)
		case SB.THUMBTRACK	:	.Scroll(.xscroll - HIWORD(wParam), 0)
		default:
			}
		return 0
		}
	maxXscroll: 0
	maxYscroll: 0
	Scroll(dx, dy)
		// pre:	dx and dy are integers
		// post:	scrolls client area horizontally and vertically by dx and dy, respective
		{
		// horizontal movement calculations
		dx = Max(Min(dx, .xscroll), .xscroll - .maxXscroll)
		.xscroll -= dx
		SetScrollInfo(.Hwnd, SB.HORZ,
			Object(fMask: SIF.POS, nPos: .xscroll, cbSize: SCROLLINFO.Size()), true)

		// vertical movement calculations
		dy = Max(Min(dy, .yscroll), .yscroll - .maxYscroll)
		.yscroll -= dy
		SetScrollInfo(.Hwnd, SB.VERT,
			Object(fMask: SIF.POS, nPos: .yscroll, cbSize: SCROLLINFO.Size()), true)

		if (dx is 0 and dy is 0)
			return false

		children = Object()
		EnumChildWindows(.Hwnd)
			{ |childHwnd|
			if .Hwnd is GetParent(childHwnd)
				children.Add(childHwnd)
			true // continue enumerating
			}
		hwdp = BeginDeferWindowPos(children.Size())
		for (child in children)
			{
			rc = GetWindowRect(child)
			ScreenToClient(.Hwnd, pt = Object(x: rc.left, y: rc.top))
			hwdp = DeferWindowPos(hwdp, child, 0,
				pt.x + dx, pt.y + dy, 0, 0,
				SWP.NOSIZE | SWP.NOZORDER | SWP.NOACTIVATE |
					(dx.Abs() is .getInner().w - 30 ? SWP.NOCOPYBITS : 0))
			}
		EndDeferWindowPos(hwdp)
		return true
		}
	SetEnabled(enabled)
		// pre:		enabled is a Boolean value
		// post:	this is enabled iff enabled is true
		{
		Assert(Boolean?(enabled))
		.ctrl.SetEnabled(enabled)
		super.SetEnabled(enabled)
		}
	SetVisible(visible)
		// pre:		visible is a Boolean value
		// post:	this is visible iff visible is true
		{
		Assert(Boolean?(visible))
		.ctrl.SetVisible(visible)
		super.SetVisible(visible)
		}
	SetReadOnly(readonly)
		{
		.ctrl.SetReadOnly(readonly)
		}
	Update()
		{
		.ctrl.Update()
		super.Update()
		}
	Destroy()
		{
		.ctrl.Destroy()
		super.Destroy()
		}
	}
