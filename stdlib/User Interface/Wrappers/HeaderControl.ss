// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	Name: "Header"
	Xstretch: 1

	New(@args)
		{
		.batchProcessing = true

		style = args.GetDefault('style', 0)
		.headerSelectPrompt = args.GetDefault('headerSelectPrompt', false)
		.CreateWindow("SysHeader32", "", WS.CHILD | WS.VISIBLE | style, 0, 0, 0, 0)
		.SetFont(text: "pW")
		.h = .Ymin = Max(.Ymin + 2, GetSystemMetrics(SM.CYVSCROLL))
		if args.Member?(0) and Object?(args[0])
			args = args[0]
		.tips = Object()
		for arg in args.Values(list:)
			.AddItem(arg)

		.Map = Object()
		.Map[HDN.TRACK] = 'HDN_TRACK'
		.Map[HDN.ENDTRACK] = 'HDN_ENDTRACK'
		.Map[HDN.ENDDRAG] = 'HDN_ENDDRAG'
		.Map[HDN.ITEMCLICK] = 'HDN_ITEMCLICK'
		.Map[HDN.BEGINTRACK] = 'HDN_BEGINTRACK'
		.Map[HDN.DIVIDERDBLCLICK] = 'HDN_DIVIDERDBLCLICK'
		}
	batchProcessing: false
	Startup()
		{
		.batchProcessing = false
		.resetToolTips()
		}
	AddItem(field, width = false, tip = false)
		{ .InsertItem(.GetItemCount(), field, width, tip) }
	InsertItem(idx, field, width = false, tip = false)
		{
		Assert(Integer?(idx) and 0 <= idx and idx <= .GetItemCount())
		append? = idx is .GetItemCount()
		text = .GetHeaderText(field)

		if width is false
			width = .getWidth(field, text)
		hdi = Object(
			mask:		HDI.TEXT | HDI.FORMAT | HDI.WIDTH,
			pszText:	TranslateLanguage(text),
			fmt:		HDF.STRING,
			cxy:		width
			)

		if -1 isnt SendMessageHditem(.Hwnd, HDM.INSERTITEM, idx, hdi)
			{
			.calc_Xmin(append?)
			if tip is false
				tip = text
			.tips.Add(tip, at: idx)
			append? ? .appendToolTip() : .resetToolTips()
			}
		else // Throwing to ensure the screen gets reloaded, and we get locals
			throw 'SendMessageHditem failed to add item'
		}
	headerBorder: 20 // determined by experimentation
	GetDefaultColumnWidth(field)
		{
		return .getWidth(field, .GetHeaderText(field))
		}
	getWidth(field, text)
		{
		heading_width = text is "" ? 0 : .TextExtent(text).x
		format_width = FieldFormatWidth(field, .AveCharWidth)
		return Max(heading_width, format_width) + .headerBorder
		}
	GetHeaderText(field, _capFieldPrompt = false)
		{
		if .headerSelectPrompt is 'no_prompts' or field is 'listrow_deleted'
			return capFieldPrompt isnt false ? capFieldPrompt : field

		if .headerSelectPrompt is false or field.Prefix?('custom_')
			return Datadict.PromptOrHeading(field)

		return Datadict.SelectPrompt(field, excludeTags: #(Internal))
		}
	GetItemFormat(i)
		{
		item = Object(mask: HDI.FORMAT)
		SendMessageHditem(.Hwnd, HDM.SETITEM, i, item)
		return item.fmt
		}
	SetItemFormat(i, format)
		{
		item = Object(mask: HDI.FORMAT, fmt: format)
		SendMessageHditem(.Hwnd, HDM.SETITEM, i, item)
		}
	SetItemWidth(i, width)
		{
		item = Object(mask: HDI.WIDTH, cxy: width)
		SendMessageHditem(.Hwnd, HDM.SETITEM, i, item)
		.calc_Xmin()
		SetWindowPos(.Hwnd, 0, 0, 0, Max(.w, .Xmin), .h,
			SWP.NOMOVE | SWP.NOZORDER | SWP.NOACTIVATE)
		.resetToolTips()
		}
	GetItemWidth(idx)
		{
		Assert(Integer?(idx) and 0 <= idx and idx < .GetItemCount())
		rc = .GetItemRect(idx)
		return rc.right - rc.left
		}
	GetItemRect(idx)
		{
		Assert(Integer?(idx) and 0 <= idx and idx < .GetItemCount())
		SendMessageRect(.Hwnd, HDM.GETITEMRECT, idx, rc = Object())
		return rc
		}
	itemTextBufferSize: 200
	getter_itemMask()
		{
		.itemMask = HDI.TEXT | HDI.WIDTH | HDI.FORMAT
		}
	GetItem(idx)
		{
		hdi = Object(
			mask: 			.itemMask
			cchTextMax:		.itemTextBufferSize)
		SendMessageHditem(.Hwnd, HDM.GETITEM, idx, hdi)
		return Object(text: hdi.pszText, width: hdi.cxy)
		}
	SwapItems(idx1, idx2)
		{
		Assert(Integer?(idx1) and 0 <= idx1 and idx1 < .GetItemCount())
		Assert(Integer?(idx2) and 0 <= idx2 and idx2 < .GetItemCount())
		hdi1 = Object(
			mask: 			.itemMask
			cchTextMax:		.itemTextBufferSize)
		hdi2 = Object(
			mask: 			.itemMask
			cchTextMax:		.itemTextBufferSize)
		SendMessageHditem(.Hwnd, HDM.GETITEM, idx1, hdi1)
		SendMessageHditem(.Hwnd, HDM.GETITEM, idx2, hdi2)
		SendMessageHditem(.Hwnd, HDM.SETITEM, idx1, hdi2)
		SendMessageHditem(.Hwnd, HDM.SETITEM, idx2, hdi1)
		.tips.Swap(idx1, idx2)
		.resetToolTips()
		}
	GetItemCount()
		{
		return .SendMessage(HDM.GETITEMCOUNT, 0, 0)
		}
	DeleteItem(idx)
		{
		Assert(Integer?(idx) and 0 <= idx and idx < .GetItemCount())
		.SendMessage(HDM.DELETEITEM, idx, 0)
		.tips.Delete(idx)
		.resetToolTips()
		}
	Clear()
		{
		.batchProcessing = true
		numItems = .GetItemCount()
		while (numItems-- > 0)
			.DeleteItem(0)
		.tips.Delete(all:)
		.batchProcessing = false
		.resetToolTips()
		}
	SetButtons(buttons)
		{
		if buttons
			.AddStyle(HDS.BUTTONS)
		else
			.RemStyle(HDS.BUTTONS)
		}
	Send(@args)
		{
		if .Parent.Method?(args[0])
			return .Parent[args[0]](@+1 args)
		return super.Send(@args)
		}
	calc_Xmin(append? = false)
		{
		if append? is true
			{
			.Xmin += .GetItemWidth(.GetItemCount() - 1)
			return
			}
		.Xmin = 0
		for (idx = 0; idx < .GetItemCount(); idx++)
			.Xmin += .GetItemWidth(idx)
		}
	Resize(x, y, w, h)
		{
		super.Resize(x, y, .w = w, .h = h)
		}
	GetReadOnly()			// read-only not applicable to header
		{ return true }

	trackMinWidth: 0
	HDN_BEGINTRACK(lParam)
		{
		idx = NMHEADER(lParam).iItem
		if false is .Send("Header_AllowTrack", idx)
			return true
		.trackMinWidth = .Send("Header_TrackMinWidth", idx)
		return false
		}
	HDN_TRACK(lParam)
		{
		nmh = NMHEADER(lParam)
		idx = nmh.iItem
		width = Max(nmh.pitem.cxy, .trackMinWidth)
		oldWidth = .GetItemWidth(idx)
		.Xmin += width - oldWidth
		.Send("HeaderTrack", idx, width)
		.setItem(idx, width)
		.resetToolTips()
		return 0
		}
	HDN_ENDTRACK(lParam)
		{
		nmh = NMHEADER(lParam)
		idx = nmh.iItem
		width = Max(nmh.pitem.cxy, .trackMinWidth)
		.Send("HeaderResize", idx, width)
		.Defer({ .setItem(idx, width) })
		return 0
		}
	setItem(idx, width)
		{
		SendMessageHditem(.Hwnd, HDM.SETITEM, idx, Object(mask: HDI.WIDTH, cxy: width))
		}

	HDN_ENDDRAG(lParam)
		{
		nmh = NMHEADER(lParam)
		item = nmh.iItem
		newpos = nmh.pitem.iOrder
		if newpos is -1
			return false
		if false is .Send("Header_AllowDrag", item)
			return true
		inc = newpos > item ? 1 : -1
		for (i = item; i isnt newpos; i += inc)
			.SwapItems(i, i + inc)
		.Send("HeaderReorder", item, newpos)
		return true
		}
	HDN_ITEMCLICK(lParam)
		{
		nmh = NMHEADER(lParam)
		.Send("HeaderClick", nmh.iItem, nmh.iButton)
		return 0
		}
	HitTest(x, y)
		{
		SendMessageHDHITTESTINFO(.Hwnd, HDM.HITTEST, 0, hti = Object(pt: Object(:x, :y)))
		return hti
		}
	ntooltips: 0
	resetToolTips()
		{
		if .batchProcessing is true
			return

		if .tips.Empty?()
			return

		.SubClass() // needed to relay
		tips = .Window.Tips()
		.SetRelay(tips.RelayEvent)
		for (i = 0; i < .ntooltips; ++i)
			tips.RemoveTool(.Hwnd, i)
		.ntooltips = .GetItemCount()
		for (i = 0; i < .ntooltips; ++i)
			if .tips[i] isnt false
				tips.AddTool(.Hwnd, .tips[i], i, rect: .getTipRect(i))
		}

	appendToolTip()
		{
		if .batchProcessing is true
			return

		if 0 > (last = .GetItemCount() - 1) or .tips.GetDefault(last, false) is false
			return

		.SubClass() // needed to relay
		tips = .Window.Tips()
		.SetRelay(tips.RelayEvent)
		tips.AddTool(.Hwnd, .tips[last], last, rect: .getTipRect(last))
		}

	getTipRect(i)
		{
		rect = .GetItemRect(i)
		// If this is called too early,
		// the rect returned from .GetItemRect has top = bottom = 0
		if rect.top is rect.bottom
			rect.bottom += .h
		return rect
		}

	HDN_DIVIDERDBLCLICK(lParam)
		{
		.Send("HeaderDividerDoubleClick", NMHEADER(lParam).iItem)
		return 0
		}
	}
