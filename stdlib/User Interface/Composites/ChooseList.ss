// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	NumLines: 10
	New(.list, select, .chooselist, .listSeparator = ' - ')
		{
		.CreateWindow("SuWhiteArrow", "", WS.VISIBLE)
		.SubClass()
		.hfont = Suneido.hfont
		.selected = select is false ? 0 : select

		.sbw = GetSystemMetrics(SM.CXVSCROLL)
		.Xmin -= .sbw

		.WithSelectObject(.hfont)
			{|hdc|
			GetTextMetrics(hdc, tm = Object())
			.lineheight = tm.Height + tm.ExternalLeading
			.Ymin = list.Size() * .lineheight
			}
		}
	Startup()
		{
		i = .selected * .lineheight
		scroll = .Parent.GetYscroll()
		.Parent.Scroll(0, scroll - i)
		}
	PAINT()
		{
		hdc = BeginPaint(.Hwnd, ps = Object())
		WithHdcSettings(hdc, [.hfont, SetBkMode: TRANSPARENT])
			{
			for i in .list.Members()
				.textOut(hdc, i)
			top = .selected * .lineheight
			FillRect(hdc,
				Object(left: 0, :top, right: .Xmin + .sbw, bottom: top + .lineheight),
				GetSysColorBrush(COLOR.HIGHLIGHT))
			WithHdcSettings(hdc, [SetTextColor: GetSysColor(COLOR.HIGHLIGHTTEXT)],
				{ .textOut(hdc, .selected) })
			}
		EndPaint(.Hwnd, ps)
		return 0
		}
	textOut(hdc, i)
		{
		item = String(.list[i])
		TextOut(hdc, 2, i * .lineheight, item, item.Size())
		}
	GETDLGCODE()
		{ return DLGC.WANTALLKEYS }
	KEYDOWN(wParam)
		{
		switch (wParam)
			{
		case VK.UP :
			.setSelect(Max(0, .selected - 1))
		case VK.DOWN :
			.setSelect(Min(.list.Size() - 1, .selected + 1))
		case VK.RETURN :
			.Result(.selected)
		case VK.ESCAPE :
			.chooselist.Result(false)
			.destroyWindow()
		default:
			}
		return 0
		}
	LBUTTONUP(lParam)
		{
		y = HIWORD(lParam)
		i = Min(.list.Size() - 1, (y / .lineheight).Int())
		.Result(i)
		return 0
		}
	ignore_mousemove: false
	MOUSEMOVE(lParam)
		{
		if .ignore_mousemove
			{
			.ignore_mousemove = false
			return 0
			}
		y = HIWORD(lParam)
		.setSelect((y / .lineheight).Int())
		return 0
		}
	setSelect(i)
		{
		if i is .selected or i >= .list.Size()
			return
		.selected = i
		i *= .lineheight
		scroll = .Parent.GetYscroll()
		height = .NumLines * .lineheight
		.ignore_mousemove = true
		// because scroll causes mousemove which changes select
		if i < scroll
			.Parent.Scroll(0, scroll - i)
		else if i >= scroll + height
			.Parent.Scroll(0, scroll + height - i - .lineheight)
		.Repaint()
		}
	Result(i)
		{
		item = String(.list[i])
		if (.listSeparator isnt '')
			item = item.BeforeFirst(.listSeparator)
		chooselist = .chooselist
		if not .Destroyed?()
			.destroyWindow()
		chooselist.Result(item)
		}
	destroyWindow()
		{
		hwnd = .Window.Hwnd
		.Window.Hwnd = 0
		DestroyWindow(hwnd)
		}
	MOUSEWHEEL(wParam)
		{
		.Send('MouseWheel', wParam)
		return 0
		}
	Resize(x, y, w, h)
		{
		super.Resize(x, y, w + .sbw, h + 1)
		}
	}
