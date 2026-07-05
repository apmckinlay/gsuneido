// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
// multi-line EditControl
EditControl
	{
	Name:		"Editor"
	Unsortable: true
	Xstretch:	1
	Ystretch:	1
	DefaultHeight: 4
	Hasfocus?:	false

	New(style = 0, .readonly = false, .font = "", .size = "", .zoom = false,
		mandatory = false, set = "", height =  false, .tabthrough = false,
		hidden = false, tabover = false, width = false, weight = false,
		readOnlyBgndColor = false, status = '', textLimit = false)
		{
		super(mandatory, readonly,
			style | WS.VSCROLL | ES.MULTILINE | ES.AUTOVSCROLL | ES.WANTRETURN,
			:hidden, :tabover, :font, :size, :weight, :width, :height,
			:readOnlyBgndColor, :status)
		.SubClass()
		.editorTextLimit = .editorTextLimit(textLimit)
		.Set(set)
		.Map[EN.CHANGE] = 'EN_CHANGE'
		.findreplacedata = Record()
		.AddContextMenuItem("Find...\tCtrl+F", .On_Find)
		.AddContextMenuItem("Print...\tCtrl+P", .On_Print)
		if .zoom is false
			.AddContextMenuItem("Zoom...\tF6", .On_Zoom)
		.SendMessage(EM.SETLIMITTEXT, .editorTextLimit, 0)
		}

	editorTextLimit(textLimit)
		{
		return Number?(textLimit)
			// Ensure specified textLimit never exceeds EditorTextLimit
			? Min(EditorTextLimit, textLimit)
			: EditorTextLimit
		}

	tabthrough: false
	GETDLGCODE(wParam)
		{
		return wParam is VK.ESCAPE or (.tabthrough and wParam is VK.TAB)
			? DLGC.WANTCHARS : DLGC.WANTALLKEYS
		}

	KEYDOWN(wParam)
		{
		return .Eval(EditorKeyDownHandler, wParam, zoomArgs: .ZoomArgs())
		}

	ZoomArgs()
		{
		return [this, .zoom, font: .font, size: .size, textLimit: .editorTextLimit]
		}

	EN_KILLFOCUS()
		{
		retVal = super.EN_KILLFOCUS()
		if (.Dirty?())
			.Send("NewValue", .Get())
		return retVal
		}

	MOUSEWHEEL(wParam, lParam)
		{
		return .HandleVScrollEdges(wParam, lParam, .Hwnd, .WndProc)
		}

	HandleVScrollEdges(wParam, lParam, hwnd, wndProc)
		{
		direction = HISWORD(wParam) > 0 ? 'UP' : 'DOWN'
		GetScrollInfo(hwnd, SB.VERT,
			info = Object(cbSize: SCROLLINFO.Size(), fMask: SIF.ALL))
		pos = GetScrollPos(hwnd, SB.VERT)
		// if scrolling up at the first row or down at the last row, it scrolls the parent
		if ((direction is 'UP' and pos is info.nMin) or
			(direction is 'DOWN' and pos + info.nPage > info.nMax))
			return wndProc.Callsuper(hwnd, WM.MOUSEWHEEL, wParam, lParam)
		return 'callsuper'
		}

	On_Print()
		{
		if '' is (text = .Get().Trim())
			return
		Params.On_Print(Object('WrapGen', text),
			title: '', name: 'print_editor', previewWindow: .Window.Hwnd)
		}

	On_Find()
		{
		s = .GetSelText()
		if s > "" and not s.Has?('\n')
			.findreplacedata.find = s
		x = FindDialog(.findreplacedata)
		if x is #next
			.On_Find_Next()
		else if x is #prev
			.On_Find_Previous()
		}

	On_Find_Next()
		{
		return .findNextPrev()
		}

	On_Find_Previous()
		{
		return .findNextPrev(prev:)
		}

	findNextPrev(prev = false)
		{
		if .findreplacedata.find.Blank?()
			return false
		from = .GetSel()[prev ? 0 : 1]
		if false is match = Find.DoFind(super.Get(), from, .findreplacedata, :prev)
			return false
		.SetSel(match[0], match[0] + match[1])
		return true
		}

	On_Zoom()
		{
		EditorZoom(@.ZoomArgs())
		}

	Get()
		{
		text = GetWindowText(.Hwnd)
		text = text.Tr("\r")
		return text
		}

	Set(value)
		{
		value = String(value)
		if value.Size() > .editorTextLimit
			{
			ProgrammerError('EditorControl Set value is over limit',
				params: Object(size: value.Size(), name: .Name))
			value = value[::.editorTextLimit - 3/*=size of '...'*/] $ '...'
			}
		value = value.Replace("\n", "\r\n")
		SetWindowText(.Hwnd, value)
		.Dirty?(false)
		}

	HasFocus?()
		{
		return .Hasfocus? or super.HasFocus?()
		}

	SetReadOnly(readOnly)
		{
		super.SetReadOnly(readOnly)
		.readonly = readOnly
		}

	MakeSummary()
		{
		if .GetHidden() is true
			return ''

		summaryLimit = 60
		text = .Get().Trim()
		summary = text.BeforeFirst('\n').Trim()[.. summaryLimit]
		if text.Size() > summary.Size()
			summary $= '...'
		return summary
		}
	}
