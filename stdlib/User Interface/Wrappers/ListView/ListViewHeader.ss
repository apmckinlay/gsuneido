// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
HeaderControl
	{
	Name: 'ListViewHeader'
	CallClass(listview)
		{
		wndproc = new this
		wndproc.ListView = listview
		wndproc.Window = listview.Window
		wndproc.Hwnd = listview.GetHeader()
		wndproc.WndProc = wndproc
		wndproc.SubClass()
		wndproc.Controller = listview
		wndproc.Init()
		wndproc.AddHwnd(wndproc.Hwnd)
		return wndproc
		}
	New2()
		{
		}
	Init()
		{
		.Map = Object()
		.Map[TTN.SHOW] = 'TTN_SHOW'
		.tip = .Construct(ToolTipControl)
		.tip.Activate(false)
		.tip.AddTool(.Hwnd, '???')
		.SetRelay(.tip.RelayEvent)
		.Defer(.set_tip_font)
		// tooltip offsets by trial & error - may not always be right
		.offset = [x: 6, y: 2]
		}
	set_tip_font()
		{
		if .tip isnt false
			.tip.SendMessage(WM.SETFONT, .SendMessage(WM.GETFONT), false)
		}
	tip: false

	// would be nice to move this into HeaderControl
	// so ListControl would also get it
	// BUT it's relying on ListView GetStringWidth
	// MAYBE get font and measure with DrawText instead?
	curItem: false
	MOUSEMOVE(lParam)
		{
		x = LOSWORD(lParam)
		y = HISWORD(lParam)
		hti = .HitTest(x, y)
		item = hti.flags is HHT.ONHEADER ? hti.iItem : false
		if item isnt .curItem
			{
			.tip.Activate(false)
			if item isnt false and
				false isnt text = .need_tip(item)
				{
				.tip.UpdateTipText(.Hwnd, text)
				.tip.Activate(true)
				}
			.curItem = item
			}
		return 'callsuper'
		}
	need_tip(item)
		{
		r = .GetItemRect(item)
		colwid = r.right - r.left
		text = .GetItem(item).text
		strwid = .ListView.GetStringWidth(text)
		// assuming header font is same as list
		return strwid + .padding > colwid ? text : false
		}
	padding: 12 // from trial & error
	TTN_SHOW(lParam)
		{
		r = .GetItemRect(.curItem)
		.tip.AdjustRect(false, r)
		p = [x: r.left - .offset.x, y: r.top - .offset.y]
		ClientToScreen(.Hwnd, p)
		nmhdr = NMHDR(lParam)
		SetWindowPos(nmhdr.hwndFrom, 0,
			p.x, p.y, 0, 0, // rect
			SWP.NOACTIVATE | SWP.NOSIZE | SWP.NOZORDER)
		return true
		}
	NCDESTROY(@args)
		{
		// FIXME: Why does this need to be a message handler? Can't we just
		//        override Hwnd.Destroy()?
		if .tip isnt false
			{
			.tip.Destroy()
			.tip = false
			}
		super.NCDESTROY(@args)
		}
	RBUTTONDOWN(lParam)
		{
		if .ListView.GetMenu() isnt false and .ListView.GetMenu().Has?('Reset Columns')
			{
			ClientToScreen(.Hwnd, p = Object(x: LOWORD(lParam), y: HIWORD(lParam)))
			ContextMenu(#('Reset Columns')).ShowCall(.ListView, p.x, p.y)
			}
		return "callsuper"
		}
	}