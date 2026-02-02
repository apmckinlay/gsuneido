// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// Sends SelectTab(i)
// NOTE: contents of tabs must be siblings of TabControl, not children
// NOTE: The 'data' parameter of Insert(i, text, data) should be an object.
// If it contains a 'tooltip' member,
// it is used to display a tooltip for the tab under the mouse cursor.
// eg. TabControl('one', 'two', 'three')
WndProc
	{
	Name: 'Tab'
	Xstretch: 1

	New(@tabs)
		{
		style = WS.CLIPSIBLINGS | WS.VISIBLE | TCS.TOOLTIPS
		if tabs.GetDefault(#vertical, false)
			style |= TCS.VERTICAL
		bottom? = tabs.GetDefault(#bottom, false) is true
		if bottom?
			style |= TCS.BOTTOM
		.CreateWindow("SysTabControl32", "", style)
		.close_button = tabs.GetDefault('close_button', false)
		if Number?(.close_button)
			.SubClass()
		// NOTE: bottom doesn't work with themes
		if bottom? or tabs.GetDefault(#themed, true) is false
			.DisableTheme()
		// have to register tooltip handle to get notifications
		.AddHwnd(SendMessage(.Hwnd, TCM.GETTOOLTIPS, 0, 0))
		.SetFont(text: 'M')
		.origYmin = .Ymin
		.Ymin += ScaleWithDpiFactor(6) /*= picked by trial and error, affected by dpi */
		.data = Object()

		// add the tabs
		for (tab in tabs.Values(list:))
			.Insert(.Count(), tab)

		.Map = Object()
		.Map[TCN.SELCHANGE] = "TCN_SELCHANGE"
		.Map[TCN.SELCHANGING] = "TCN_SELCHANGING"
		.Map[TTN.GETDISPINFO] = "TTN_GETDISPINFO"
		.Map[NM.CLICK] = "NM_CLICK"
		}
	Insert(i, text, data = #(tooltip: ""), image = -1)
		{
		index = SendMessageTcitem(.Hwnd, TCM.INSERTITEM, i,
			Object(mask: TCIF.TEXT | TCIF.IMAGE,
				pszText: .handleAmpersands(text), iImage: image))
		if index isnt -1
			.data.Add(data, at: index)
		}
	handleAmpersands(text)
		{
		return text.Replace('&', '&&')
		}
	Count()
		{
		return SendMessage(.Hwnd, TCM.GETITEMCOUNT, 0, 0)
		}
	Remove(i)
		{
		if (SendMessage(.Hwnd, TCM.DELETEITEM, i, 0) isnt 0)
			{
			// Adjust data
			.data.Delete(i)
			return true
			}
		return false
		}
	Select(i)
		{
		SendMessage(.Hwnd, TCM.SETCURSEL, i, 0)
		}
	Move(i, newPos)
		{
		data = .GetData(i)
		text = .GetText(i)
		.Remove(i)
		.Insert(newPos, text, :data, image: data.image)
		if newPos is 0
			.ensureFirstTabVisible()
		}
	ensureFirstTabVisible()
		{
		EnumChildWindows(.Hwnd)
			{ |child|
			pos = SendMessage(child, UDM.GETPOS32, 0, 0)
			success = HISWORD(pos) is 0 // check if it is system built-in UpDown control
			if success and LOSWORD(pos) > 0
				{
				// UDM.SETPOS32 does not scroll the tab, has to mimic mouse clicking
				for .. .Count()
					{
					SendMessage(child, WM.LBUTTONDOWN, 0, 0)
					SendMessage(child, WM.LBUTTONUP, 0, 0)
					}
				// handle drawing issue on updown arrows when switching tab from tree
				rc = .GetClientRect().ToWindowsRect()
				.Defer(uniqueID: 'EnsureFirstTabVisible')
					{ InvalidateRect(.Hwnd, rc, false) }
				}
			}
		}
	GetSelected()
		// post:	returns the index of the selected tab
		{ return .SendMessage(TCM.GETCURSEL, 0, 0) }
	SetText(i, text)
		{
		SendMessageTcitem(.Hwnd, TCM.SETITEM, i,
			Object(mask: TCIF.TEXT, pszText: text))
		}
	GetText(i)
		{
		tci = Object(mask: TCIF.TEXT, cchTextMax: 64) /*= max length */
		SendMessageTcitem(.Hwnd, TCM.GETITEM, i, tci)
		return tci.pszText
		}
	SetImage(i, img)
		{
		SendMessageTcitem(.Hwnd, TCM.SETITEM, i,
			Object(mask: TCIF.IMAGE, iImage: img))
		if img isnt .close_button and i is .iUnder
			.origImage = img
		}
	GetImage(i)
		{
		tci = Object(mask: TCIF.IMAGE)
		SendMessageTcitem(.Hwnd, TCM.GETITEM, i, tci)
		return tci.iImage
		}
	SetImageList(images)
		{
		SendMessage(.Hwnd, TCM.SETIMAGELIST, 0, images)
		}
	SetPadding(hpadding, vpadding)
		{
		.Ymin = .origYmin + vpadding * 2
		SendMessage(.Hwnd, TCM.SETPADDING, 0, hpadding | vpadding << 16)
		}
	GetData(i)
		{
		return .data[i]
		}
	ForEachTab(block)
		{
		i = 0
		.data.Each({ block(it, idx: i++) })
		}
	SetData(i, data)
		{
		.data[i] = data
		}
	TCN_SELCHANGE()
		{
		i = SendMessage(.Hwnd, TCM.GETCURSEL, 0, 0)
		.Send("SelectTab", i)
		if true isnt .Send('Tab_AllowDrag')
			.origImage = .GetImage(i)
		return 0
		}
	TCN_SELCHANGING()
		{
		if (.Send("TabControl_SelChanging") is true)
			return true
		return false
		}
	NM_CLICK()
		{
		.Send("TabClick",
			SendMessage(.Hwnd, TCM.GETCURSEL, 0, 0))
		return 0
		}
	TTN_GETDISPINFO(lParam)
		{
		p = GetMessagePos()
		point = Object(x: LOWORD(p), y: HIWORD(p))
		ScreenToClient(.Hwnd, point)
		i = .tabIndex(point.x, point.y)
		if ((i >= 0) and ((x = .GetData(i)).Member?('tooltip')))
			StructModify(NMTTDISPINFO2, lParam, { it.szText = x.tooltip })
			// Probably better to use TTM.UPDATETIPTEXT like ToolTipControl.UpdateTipText
			// but couldn't get it to work
		return 0
		}
	ContextMenu(x, y)
		{
		return .Send("TabContextMenu", x, y)
		}
	SetReadOnly(readonly /*unused*/) // don't want to disable tabs
		{
		}
	GetReadOnly()			// read-only not applicable to tabs
		{ return true }

	iUnder: -1
	MOUSEMOVE(lParam)
		{
		x = LOSWORD(lParam)
		y = HISWORD(lParam)
		i = .tabIndex(x, y)
		if .draggedTab isnt false
			.draggingTab(i, x, y)
		if i isnt .iUnder
			{
			.MOUSELEAVE()
			if i isnt -1 and .draggedTab is false
				{
				.origImage = .GetImage(i)
				.SetImage(.iUnder = i, .close_button)
				TrackMouseEvent(Object(cbSize: TRACKMOUSEEVENT.Size(),
					dwFlags: TME.LEAVE, dwHover: 500, hwndTrack: .Hwnd))
				}
			}
		return 'callsuper'
		}
	tabIndex(x, y)
		{
		return SendMessageTabHitTest(.Hwnd, TCM.HITTEST, 0, Object(pt: Object(:x, :y)))
		}

	draggingTab(i, x, y, allowSeperate? = false)
		{
		inTabBar? = .inTabBar?(x, y)
		if inTabBar?
			{
			// If i is -1, we are either past the last tab, or a in the tab padding
			i = i is -1 ? .tabIndex(x, .Ymin / 2) : i
			if i is -1 // If i is still -1, we are past the last tab
				.moveToEnd(x)
			else if .draggedTab isnt i
				.dragTab(i, x)
			}
		else if allowSeperate?
			.Send(#Tab_Separate, .draggedTab)
		SetCursor(LoadCursor(ResourceModule(), inTabBar? ? IDC.HSPLITBAR : IDC.DRAG1COPY))
		}

	padding: 15
	inTabBar?(x, y)
		{
		return y >= -.padding and y <= .Ymin + .padding and x >= 0
		}

	moveToEnd(x)
		{
		if .draggedTab is lastTab = .Count() - 1
			return
		SendMessageRect(.Hwnd, TCM.GETITEMRECT, lastTab, rc = Object())
		if x + .draggedTabWidth >= rc.right
			.moveTab(lastTab)
		}

	dragTab(i, x)
		{
		if 1 < (.draggedTab - i).Abs()
			.moveTab(i)
		else
			{
			SendMessageRect(.Hwnd, TCM.GETITEMRECT, i, rc = Object())
			moveToIdx? = i < .draggedTab
				? rc.left > x - .draggedTabWidth
				: rc.right < x + .draggedTabWidth
			if moveToIdx?
				.moveTab(i)
			}
		}

	moveTab(idx)
		{
		.Send(#MoveTab, .draggedTab, idx)
		.draggedTab = idx
		}

	LBUTTONUP(lParam)
		{
		if .draggedTab is false
			return 0
		i = .tabIndex(x = LOSWORD(lParam), y = HISWORD(lParam))
		.draggingTab(i, x, y, allowSeperate?:)
		.dragEnd()
		return 0
		}

	LBUTTONDOWN(lParam)
		{
		.dragEnd()
		x = LOSWORD(lParam)
		y = HISWORD(lParam)
		i = .onClosingIcon(x, y)
		if i isnt -1
			{
			.Send("Tab_Close", i)
			.iUnder = -1
			return 0 // suppress NM_CLICK
			}
		.dragStart(x, y)
		return 'callsuper'
		}

	draggedTab: false
	draggedTabWidth: false
	dragEnd()
		{
		if .draggedTab is false
			return
		.Send('Tab_DragEnd', .draggedTab)
		.draggedTabWidth = .draggedTab = false
		SetCursor(LoadCursor(NULL, IDC.ARROW))
		ReleaseCapture()
		}

	dragStart(x, y)
		{
		if true isnt .Send('Tab_AllowDrag')
			return
		SetCapture(.Hwnd)
		.draggedTab = .tabIndex(x, y)
		SendMessageRect(.Hwnd, TCM.GETITEMRECT, .draggedTab, rc = Object())
		.draggedTabWidth = rc.right - rc.left
		}

	RBUTTONDOWN(lParam /*unused*/)
		{
		.dragEnd()
		return 'callsuper'
		}

	onClosingIcon(x, y)
		{
		i = SendMessageTabHitTest(.Hwnd, TCM.HITTEST, 0, ht = Object(pt: Object(:x, :y)))
		return ht.flags is TCHT.ONITEMICON ? i : -1
		}
	MOUSELEAVE()
		{
		if .iUnder isnt -1
			{
			if .draggedTab is false
				.SetImage(.iUnder, .origImage)
			.iUnder = -1
			}
		return 'callsuper'
		}
	Destroy()
		{
		.DelHwnd(SendMessage(.Hwnd, TCM.GETTOOLTIPS, 0, 0))
		super.Destroy()
		}
	}
