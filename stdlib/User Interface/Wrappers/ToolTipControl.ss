// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Hwnd
	{
	Name: "ToolTips"
	New()
		{
		.CreateWindow(TOOLTIPS_CLASS, "",
			WS.POPUP | TTS.NOPREFIX | TTS.ALWAYSTIP,
			WS_EX.TRANSPARENT)

		.SendMessage(TTM.SETMAXTIPWIDTH, 0, 600) /*=tip max width, also handles newlines*/

		// Without SetWindowPos, TreeView tooltip would
		// bring the parent window to the top when inactive
		SetWindowPos(.Hwnd, HWND.TOPMOST, 0, 0, 0, 0,
			SWP.NOMOVE | SWP.NOSIZE | SWP.NOACTIVATE)

		.Map = Object()
		.Map[TTN.SHOW] = 'TTN_SHOW'
		.Map[TTN.POP] = 'TTN_POP'
		.Map[TTN.GETDISPINFO] = 'TTN_GETDISPINFO'
		}

	// Windows sends notifications to parent/owner
	// Suneido redirects these to the control
	// here we're redirecting them back to the parent/owner!
	TTN_SHOW(lParam)
		{
		if .observer.Method?(#ToolTip_Show)
			.observer.ToolTip_Show(:lParam)
		nmhdr = NMHDR(lParam)
		if .Window.HwndMap.Member?(nmhdr.idFrom)
			return .Window.HwndMap[nmhdr.idFrom].Notify(TTN.SHOW, lParam)
		return 0
		}

	TTN_GETDISPINFO(lParam)
		{
		if .observer.Method?(#ToolTip_GetDispInfo)
			.observer.ToolTip_GetDispInfo(:lParam)

		uFlags = 0
		// Using StructModify to read uFlags becuase there is no direct builtin to
		// convert NMTTDISPINFO2 struct to a Suneido Object
		StructModify(NMTTDISPINFO2, lParam, { uFlags = it.uFlags })
		// If the TTF_IDISHWND flag is set, then the idFrom field of the NMHDR
		// indicates the window handle of the control. In that case, we want to
		// reflect the notification to the tool window.
		if TTF.IDISHWND is (uFlags & TTF.IDISHWND)
			{
			nmhdr = NMHDR(lParam)
			hwndTool = nmhdr.idFrom
			if .Window.HwndMap.Member?(hwndTool)
				return .Window.HwndMap[hwndTool].Notify(TTN.GETDISPINFO, lParam)
			}
		// On the other hand, if the flag isn't set, we should reflect the
		// message to our parent window.
		else
			{
			hwndParent = .Parent.Hwnd
			if .Window.HwndMap.Member?(hwndParent)
				return .Window.HwndMap[hwndParent].Notify(TTN.GETDISPINFO, lParam)
			}
		return 0
		}
	TTN_POP(lParam)
		{
		if .observer.Method?(#ToolTip_Pop)
			.observer.ToolTip_Pop(:lParam)
		return 0
		}
	observer: class{}
	Observer(ob)
		{
		.observer = ob
		}

	AddTool(hwnd, text, id = 0, flags = 0, lParam = 0, rect = false)
		{
		if .Destroyed?()
			return
		// can't use TTF.SUBCLASS
		// because if control SubClass's as well
		// Windows leaks something
		// and after 32k tooltips, window creation will fail
		// use RelayEvent instead (handled in WndProc)
		if text isnt LPSTR_TEXTCALLBACK
			text = TranslateLanguage(text)
		if rect is false
			{
			// If no rectangle is specified, put tooltip on hwnd
			id = hwnd
			flags |= TTF.IDISHWND
			rect = Object()
			}
		ti = Object(
			cbSize:		TOOLINFO.Size()
			uFlags:		flags
			hwnd:		hwnd
			uId:		id
			rect:		rect
			lpszText:	text
			lParam: 	lParam
			)
		f = text is LPSTR_TEXTCALLBACK ? SendMessageTOOLINFO2 : SendMessageTOOLINFO
		return 0 isnt f(.Hwnd, TTM.ADDTOOL, 0, ti)
		}
	RemoveTool(hwnd, id = false)
		{
		if not .Member?(#Hwnd)
			return // already destroyed
		// if id not specified, then assumes tool is on hwnd
		// if AddTool specified rect, then must specify id of 0
		ti = Object(
			cbSize:	TOOLINFO.Size()
			hwnd:	hwnd
			uId:	id is false ? hwnd : id
			)
		SendMessageTOOLINFO(.Hwnd, TTM.DELTOOL, 0, ti)
		}
	RemoveAllTools()
		{
		if .Destroyed?()
			return // already destroyed
		// Enumerate the tools and then destroy them last to first.
		n = .GetToolCount()
		ti = Object(cbSize: TOOLINFO.Size())
		tools = Object()
		for (k = 0; k < n; ++k)
			{
			// Contrary to the MSDN documentation, TTM_ENUMTOOLS doesn't have a
			// return value. It seems to randomly return 1 or 0. See:
			//     http://connect.microsoft.com/VisualStudio/feedback/details/
			//         634946/ttm-enumtools-documentation-error
			ti.lpszText = 0
			SendMessageTOOLINFO2(.Hwnd, TTM.ENUMTOOLS, k, ti)
			tools.Add(ti.Copy())
			}
		// Reverse-order deletion is based on our too-clever-by-half, but
		// probably correct, assumption that the control's backing store is
		// an array.
		for ti in tools.Reverse!()
			SendMessageTOOLINFO2(.Hwnd, TTM.DELTOOL, 0, ti)
		Assert(.GetToolCount() is: 0)
		return k
		}
	RelayEvent(hwnd, message, wParam, lParam)
		{
		if .Destroyed?()
			return // destroyed
		ClientToScreen(hwnd, pt = Object(x: LOSWORD(lParam) y: HISWORD(lParam)))
		msg = Object(
			hwnd:		hwnd
			message:	message
			lParam:		lParam
			wParam:		wParam
			time:		GetTickCount()
			pt:			pt
			)
		SendMessageMSG(.Hwnd, TTM.RELAYEVENT, 0, msg)
		}
	GetToolCount()
		{
		return .SendMessage(TTM.GETTOOLCOUNT)
		}
	Pop()
		{
		.SendMessage(TTM.POP)
		}
	Popup()
		{
		.SendMessage(TTM.POPUP)
		}
	Activate(activate? = true)
		{
		.SendMessage(TTM.ACTIVATE, activate?)
		}
	MaxTipLength: 3000
	UpdateTipText(hwnd, text, id = false)
		{
		ti = Object(
			cbSize:	TOOLINFO.Size()
			uFlags: id is false ? TTF.IDISHWND : 0
			hwnd:	hwnd
			uId:	id is false ? hwnd : id
			hinst:	0
			lpszText: TranslateLanguage(text).Ellipsis(.MaxTipLength, true)
			)
		SendMessageTOOLINFO(.Hwnd, TTM.UPDATETIPTEXT, 0, ti)
		}
	AdjustRect(fLarger, rect)
		{
		SendMessageRect(.Hwnd, TTM.ADJUSTRECT, fLarger, rect)
		}
	TrackActivate(hwnd, id = false, activate? = true)
		{
		ti = Object(
			cbSize: TOOLINFO.Size()
			uFlags: id is false ? TTF.IDISHWND : 0
			hwnd:   hwnd
			uId:    id is false ? hwnd : id
			)
		SendMessageTOOLINFO(.Hwnd, TTM.TRACKACTIVATE, activate? ? 1 : 0, ti)
		}
	TrackPosition(x, y)
		{ SendMessage(.Hwnd, TTM.TRACKPOSITION, 0, MAKELONG(x, y)) }
	GetReadOnly()			// read-only not applicable to tooltip
		{ return true }
	}

