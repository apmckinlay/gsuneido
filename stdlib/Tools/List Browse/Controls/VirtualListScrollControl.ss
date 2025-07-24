// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	Name: 		"VirtualListScroll"
	Xstretch: 	1
	Ystretch: 	1
	Xmin: 		50
	Ymin: 		10

	hdrCornerCtrl: false
	expandExtra: false
	New(control, bar, floating, thinBorder, hdrCornerCtrl = false, expandExtra = false)
		{
		.CreateWindow("SuBtnfaceArrow", "",
			(thinBorder ? WS.BORDER : 0) | WS.VISIBLE,
			(thinBorder ? 0 : WS_EX.CLIENTEDGE) | WS_EX.CONTROLPARENT)
		.SubClass()

		.ctrl = .Construct(control)
		.bar = .Construct(bar)
		.floating = .Construct(floating)
		if hdrCornerCtrl isnt false
			{
			.hdrCornerCtrl = .Construct(hdrCornerCtrl)
			SetWindowPos(.hdrCornerCtrl.Hwnd, HWND.TOP, 0, 0, 0, 0, SWP.NOMOVE)
			}
		if expandExtra isnt false and expandExtra isnt ''
			{
			.expandExtra = .Construct(expandExtra)
			SetWindowPos(.expandExtra.WndPane.Hwnd, HWND.TOP, 0, 0, 0, 0, SWP.NOMOVE)
			.expandExtra.SetVisible(false)
			}
		}

	size: false
	Resize(x, y, w, h)
		{
		changed = .size isnt newSize = Object(:x, :y, :w, :h)
		.size = newSize

		super.Resize(x, y, w, h)

		.ResizeWindow(:changed)
		}

	ResizeWindow(changed = false)
		{
		client_rect = .GetClientRect()
		bar_width = .bar.GetWidth()
		client_w = client_rect.GetWidth()
		client_h = client_rect.GetHeight()

		if not changed
			_slowQueryLog = Object(suppressed: true, from: 'ResizeWindow')
		bar? = .Send('VirtualListScroll_Resize', client_w, client_h)

		.floating.SetVisible(false)
		if bar? is 'none'
			{
			.bar.SetVisible(false)
			.ctrl.SetVisible(false)
			}
		else if bar? is true
			{
			.bar.Resize(client_w - bar_width, 0, bar_width, client_h)
			.bar.SetVisible(true)
			.ctrl.Resize(0, 0, client_w - bar_width, client_h)
			.ctrl.SetVisible(true)
			}
		else
			{
			.bar.SetVisible(false)
			.ctrl.Resize(0, 0, client_w, client_h)
			.ctrl.SetVisible(true)

			if false isnt pos = .Send('VirtualListScroll_FloatingPosition', client_w)
				{
				w = .floating.GetWidth()
				h = .floating.GetHeight()
				.floating.SetVisible(true)
				.floating.Resize(client_w - w, client_h - h - pos, w, h)
				}
			}
		.resizeHdrCornerCtrl()
		}

	resizeHdrCornerCtrl()
		{
		if .hdrCornerCtrl isnt false
			{
			width = .Send('GetExpandBarWidth')
			.hdrCornerCtrl.Resize(0, 0, width, width)
			}
		}

	GetChildren()
		{
		children = Object(.ctrl, .bar, .floating)
		if .hdrCornerCtrl isnt false
			children.Add(.hdrCornerCtrl)
		if .expandExtra isnt false
			children.Add(.expandExtra)
		return children
		}

	GetHdrCornerCtrl()
		{
		return .hdrCornerCtrl
		}

	Destroy()
		{
		.ctrl.Destroy()
		.floating.Destroy()
		.bar.Destroy()
		if .hdrCornerCtrl isnt false
			.hdrCornerCtrl.Destroy()
		if .expandExtra isnt false
			.expandExtra.Destroy()
		super.Destroy()
		}
	}