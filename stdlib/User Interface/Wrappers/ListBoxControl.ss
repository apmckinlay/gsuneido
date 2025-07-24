// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	Name:		"ListBox"
	Xmin: 		100
	Xstretch: 	1
	Ymin: 		100
	Ystretch: 	1

	New(@args)
		{
		style = WS.VISIBLE | WS.VSCROLL | WS.HSCROLL |
			LBS.NOTIFY | LBS.HASSTRINGS | LBS.NOINTEGRALHEIGHT
		if args.GetDefault("multicolumn", false) is true
			style |= LBS.MULTICOLUMN
		if args.GetDefault("sort", false) is true
			style |= LBS.SORT
		if args.GetDefault("multiSelect", false) is true
			style |= LBS.MULTIPLESEL

		.CreateWindow("listbox", "", style, WS_EX.CLIENTEDGE)
		.SubClass() // to get ContextMenu

		.themed? = args.GetDefault("themed?", false)
		.brush = CreateSolidBrush(IDE_ColorScheme.DefaultStyle.defaultBack)
		font = args.GetDefault("font", "")
		size = args.GetDefault("size", "")
		weight = args.GetDefault("weight", "")
		.SetFont(font, size, weight)

		if args.Size(list:) is 1 and Object?(args[0])
			args = args[0]
		try for s in args.Values(list:)
			.AddItem(s)

		.Map = Object()
		.Map[LBN.SELCHANGE] = 'SELCHANGE'
		}
	wid: 0
	n: 0
	AddItem(s, n = false)
		{
		s = .limitStringWidth(s)
		i = SendMessageTextIn(.Hwnd, LB.ADDSTRING, 0, s)
		.SetData(i, n is false ? .n++ : n)
		}
	InsertItem(s, i)
		{
		SendMessageTextIn(.Hwnd, LB.INSERTSTRING, i, .limitStringWidth(s))
		}
	maxCharacters: 250
	marginRight: 10
	limitStringWidth(s)
		{
		s = String(s)[.. .maxCharacters]
		ex = .TextExtent(s)
		if ex.x > .wid
			{
			.wid = ex.x
			.SendMessage(LB.SETHORIZONTALEXTENT, .wid + .marginRight, 0)
			}
		return s
		}
	LBUTTONDOWN(lParam)
		{
		.Send(#ListBoxSelect, i = .itemClicked(HISWORD(lParam)))
		return i is -1 ? 0 : 'callsuper'
		}
	itemClicked(y)
		{
		i = (y / .GetItemHeight()).Int() + .GetTopIndex()
		return i < .GetCount() ? i : -1
		}
	LBUTTONDBLCLK(lParam)
		{
		if -1 isnt i = .itemClicked(HISWORD(lParam))
			{
			.SetItemSelected(i, true)
			.Send("ListBoxDoubleClick", i)
			}
		return 0
		}
	SELCHANGE()
		{
		.Send("ListBoxSelect", .GetCurSel())
		return 0
		}
	GetCount()
		{
		return .SendMessage(LB.GETCOUNT)
		}
	Count() // deprecated (inconsistent name)
		{
		return .SendMessage(LB.GETCOUNT)
		}
	DeleteItem(i)
		{
		// WARNING: this will mess up .n
		.SendMessage(LB.DELETESTRING, i)
		// TODO: adjust wid
		}
	DeleteAll()
		{
		for (i = .Count() - 1; i >= 0; --i)
			.DeleteItem(i)
		.n = 0
		.wid = 0
		}
	SetItemSelected(i, selected? = true)
		{
		.SendMessage(LB.SETSEL, selected? ? 1 : 0, i)
		}
	GetItemSelected?(i)
		{
		return .SendMessage(LB.GETSEL, i) isnt 0
		}
	SetCurSel(i)
		{
		.SendMessage(LB.SETCURSEL, i)
		}
	GetCurSel()
		{
		return .SendMessage(LB.GETCURSEL)
		}
	GetSelected()
		{
		return .GetData(.GetCurSel())
		}
	GetAllSelected()
		{
		selected = Object()
		for i in .. .GetCount()
			if .GetItemSelected?(i)
				selected[i] = .GetText(i)
		return selected
		}
	Get()
		{
		return .GetText(.GetCurSel())
		}
	SetData(i, n)
		{
		.SendMessage(LB.SETITEMDATA, i, n)
		}
	GetData(i)
		{
		return .SendMessage(LB.GETITEMDATA, i, 0)
		}
	GetText(i)
		{
		// need to check if size is 0 because you can't create buffer with size of 0
		if LB.ERR is (size = .SendMessage(LB.GETTEXTLEN, i)) or size is 0
			return ""
		return .SendMessageTextOut(LB.GETTEXT, i, size).text
		}
	SetColumnWidth(w)
		{
		.SendMessage(LB.SETCOLUMNWIDTH, w)
		}
	FindString(text)
		{
		return SendMessageTextIn(.Hwnd, LB.FINDSTRING, -1, text)
		}
	ContextMenu(x, y)
		{
		if x is 0 and y is 0 // keyboard
			{
			i = .GetCurSel()
			if i is -1
				return 0
			pt = Object(x: 10, y: (i - .GetTopIndex() + 1) * .GetItemHeight())
			ClientToScreen(.Hwnd, pt)
			x = pt.x
			y = pt.y
			}
		else // mouse right click
			{
			pt = Object(:x, :y)
			ScreenToClient(.Hwnd, pt)
			if -1 is i = .itemClicked(pt.y)
				return 0
			.SetItemSelected(i)
			.SetCurSel(i)
			}
		.Send("ListBox_ContextMenu", x, y)
		return 0
		}
	GetItemHeight()
		{
		return .SendMessage(LB.GETITEMHEIGHT)
		}
	GetTopIndex()
		{
		return .SendMessage(LB.GETTOPINDEX)
		}
	SetTopIndex(i)
		{
		return .SendMessage(LB.SETTOPINDEX, i)
		}
	GetHorizontalExtent()
		{
		.SendMessage(LB.GETHORIZONTALEXTENT)
		}
	theme: false
	ResetTheme()
		{
		if not .themed?
			return
		.theme = IDE_ColorScheme.GetTheme()
		DeleteObject(.brush)
		.brush = CreateSolidBrush(.theme.defaultBack)
		.Repaint(true)
		}
	CTLCOLORLISTBOX(wParam)
		{
		// Cannot switch to WithBkMode as the end point is ambiguous
		SetBkMode(wParam, TRANSPARENT)
		if .themed? and .theme isnt false
			SetTextColor(wParam, .theme.defaultFore)
		return .brush
		}
	Destroy()
		{
		DeleteObject(.brush)
		super.Destroy()
		}
	}
