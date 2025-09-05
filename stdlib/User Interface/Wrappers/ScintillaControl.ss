// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
WndProc
	{
	Name:		"Editor"
	Xmin: 		200
	Ymin: 		100
	Xstretch:	1
	Ystretch:	1
	DefaultFontSize: 11
	Unsortable: true
	marker_find: 0

	New(style = 0, .readonly = false, lexer = 'NULL', wrap = false, set = false,
		margin = 12, exStyle = 0, height = false, .tabthrough = false)
		{
		if height is 1
			.tabthrough = true
		if not readonly
			style |= WS.TABSTOP
		if exStyle is 0
			exStyle = WS_EX.CLIENTEDGE
		.CreateWindow("Scintilla", NULL,
			style | WS.VISIBLE | WS.VSCROLL | WS.HSCROLL |
			WS.CLIPSIBLINGS | WS.CLIPCHILDREN /* fix disappearing calltips */,
			:exStyle)
		.SubClass()
		.SETCODEPAGE(CP.ACP)
		.SetEolMode(SC.EOL_CRLF)
		.SetTechnology(1) // SC.TECHNOLOGY_DIRECTWRITE
		.SetLexer(SCLEX[lexer])

		if readonly
			.SETREADONLY(1)

		.InitFont()

		widths = #(tab: 4, scroll: 10)
		.SetTABWIDTH(widths.tab)
		.SetYCARETPOLICY(SC.CARET_SLOP | SC.CARET_EVEN, 2)
		.SetMultiPaste(SC.MULTIPASTE_EACH)

		.Send("Data")

		.setMap()
		.SetMARGINWIDTHN(1, ScaleWithDpiFactor(margin))

		.SetWrap(wrap)

		if set isnt false
			.Defer({ .Set(set) })

		.InitMarkers()

		.SetScrollWidth(widths.scroll)
		.SetScrollWidthTracking(true)
		.findreplacedata = Record()
		}

	InitFont()
		{
		.SendMessageTextIn(SCI.STYLESETFONT, SC.STYLE_DEFAULT, 	StdFonts.Mono())
		.StyleSetSize(SC.STYLE_DEFAULT, .DefaultFontSize)
		}

	setMap()
		{
		.Map = Object()
		.Map[SCEN.KILLFOCUS] = 'SCEN_KILLFOCUS'
		.Map[SCEN.SETFOCUS] = 'SCEN_SETFOCUS'
		.Map[EN.CHANGE] = 'EN_CHANGE'
		.Map[SCN.DOUBLECLICK] = 'SCN_DOUBLECLICK'
		.Map[SCN.MODIFIED] = 'SCN_MODIFIED'
		.Map[SCN.ZOOM] = 'SCN_ZOOM'
		}

	InitMarkers()
		{
		markers = .BaseStyling()
		for mem in markers.Members()
			{
			marker = markers[mem].marker.Set_default(false)
			.DefineMarker(mem, marker.marker, marker.fore, marker.back)
			}
		}

	BaseStyling()
		{
		markers = []
		markers[.marker_find] = [level: 0, marker: [SC.MARK_ROUNDRECT, back: CLR.GREEN]]
		return markers
		}

	DefineStyle(n, fore = false, back = false,
		bold = false, underline = false, italic = false)
		{
		if bold or italic
			{
			// can't seem to set bold unless you set font & size
			.StyleSetFont(n, .StyleGetFont(SC.STYLE_DEFAULT))
			.StyleSetSize(n, .StyleGetSize(SC.STYLE_DEFAULT))
			.StyleSetBold(n, bold)
			.StyleSetItalic(n, italic)
			}
		.StyleSetUnderline(n, underline)
		if fore isnt false
			.StyleSetFore(n, fore)
		if back isnt false
			.StyleSetBack(n, back)
		}
	DefineMarker(n, type, fore = false, back = false)
		{
		.MarkerDefine(n, type)
		if fore isnt false
			.MarkerSetFore(n, fore)
		if back isnt false
			.MarkerSetBack(n, back)
		}
	DefineXPMMarker(n, xpm, fore = false, back = false)
		{
		SendMessageTextIn(.Hwnd, SCI.MARKERDEFINEPIXMAP, n, xpm)
		if fore isnt false
			.MarkerSetFore(n, fore)
		if back isnt false
			.MarkerSetBack(n, back)
		}
	DefineIndicator(n, style, fore = false)
		{
		.IndicSetStyle(n, style)
		if fore isnt false
			.IndicSetFore(n, fore)
		}
	SetWordChars(chars) // override
		{
		SendMessageTextIn(.Hwnd, SCI.SETWORDCHARS, 0, chars)
		}
	GetWordChars() // override
		{
		n = .GETWORDCHARS()
		if n <= 0
			return ""
		return SendMessageTextOut(.Hwnd, SCI.GETWORDCHARS, 0, n).text
		}

	AnnotationSetText(line, text)
		{
		.SendMessageTextIn(SCI.ANNOTATIONSETTEXT, line, text)
		}

	SCIAutocShow(length, matches)
		{
		SendMessageTextIn(.Hwnd, SCI.AUTOCSHOW, length, matches.Join(' '))
		}

	Default(@args)
		// this translates .Xyz(...) to .SendMessage(SCI.XYZ, ...)
		{
		f = args[0].Upper()
		if not SCI.Member?(f)
			throw "method not found: " $ args[0]
		args[0] = SCI[f]
		return .SendMessage(@args)
		}

	EN_CHANGE()
		{
		.astWriterManager = false
		.Send("EN_CHANGE")
		return 0
		}
	SCN_MODIFIED(lParam)
		{
		modificationType = SCNotification(lParam).modificationType
		.Send("Scintilla_Modified", type: modificationType)

		// The following condition is trying to trap the event where the user
		// drags text into the control. This event does not set the focus in the control
		// by default and therefore the control does not set its dirty state.
		// This is manually setting the focus to this control to handle dirty
		userModified = SC.STARTACTION | SC.PERFORMED_USER | SC.MOD_INSERTTEXT
		if not .setMethodModifying? and .GetReadOnly() is false and not .HasFocus?() and
			(modificationType & userModified) is userModified and .Dirty?()
			.Send("NewValue", .Get())

		return 0
		}

	SetMethodModifying?()
		{
		return .setMethodModifying?
		}

	SCN_ZOOM()
		{
		.Send("Scintilla_Zoom")
		return 0
		}
	SCN_DOUBLECLICK()
		{
		.Send('Scintilla_DoubleClick')
		return 0
		}
	SCEN_KILLFOCUS()
		{
		.Send('Scintilla_KillFocus')
		if (.Dirty?())
			.Send("NewValue", .Get())
		return 0
		}

	SCEN_SETFOCUS()
		{
		.Send('Scintilla_SetFocus')
		return 0
		}

	MOUSEWHEEL(wParam, lParam)
		{
		if KeyPressed?(VK.CONTROL)
			return .WndProc.Callsuper(.Hwnd, WM.MOUSEWHEEL, wParam, lParam)

		return EditorControl.HandleVScrollEdges(wParam, lParam, .Hwnd, .WndProc)
		}

	GotoLine(line, noFocus? = false) // overload
		{
		if noFocus? is false
			.SetFocus()
		.EnsureVisibleEnforcePolicy(line)
		.GOTOLINE(line)
		}
	GetLine(line = false)
		{
		if line is false
			line = .LineFromPosition()
		n = .LineLength(line)
		if n <= 0
			return ""
		return SendMessageTextOut(.Hwnd, SCI.GETLINE, line, n).text
		}
	LineFromPosition(pos = false) // overload
		{
		if pos is false
			pos = .GetSelectionStart()
		return .LINEFROMPOSITION(pos)
		}
	GetRange(start, end) // gets from start to end - 1
		{
		start = .intoRange(start)
		end = .intoRange(end)
		s = ""
		for (i = start; i < end; i += .chunkSize)
			s $= SendMessageTextRange(.Hwnd, SCI.GETTEXTRANGE,
				i, Min(i + .chunkSize, end))
		return s
		}
	intoRange(i)
		{
		return Max(0, Min(i, .GetTextLength()))
		}
	GetCurrentReference()
		{
		// e.g. ".foo", "Alert", "Stack.", "Stack.Push"
		// used by calltips and auto complete
		max = 1024
		x = SendMessageTextOut(.Hwnd, SCI.GETCURLINE, max, max)
		line = x.text[.. x.result].RightTrim('(') // result is caret position
		return line.Extract("[._[:alnum:]?!]+$")
		}
	SelSize()
		{
		return .GetSelectionEnd() - .GetSelectionStart()
		}
	GetSelText()
		{
		sel = .GetSelect()
		return .GetRange(sel.cpMin, sel.cpMax)
		}
	GetSelect()
		{
		return Object(
			cpMin: .GetSelectionStart(),
			cpMax: .GetSelectionEnd())
		}
	SetSelect(i, n = 0)
		{
		.SetSel(i, i + n)
		}
	SetVisibleSelect(i, n)
		{
		.EnsureRangeVisible(i, i + n)
		.SetSelect(i, n)
		}
	EnsureRangeVisible(from, to)
		{
		first = .LineFromPosition(from)
		last = .LineFromPosition(to)
		for (line = first; line <= last; ++line)
			.EnsureVisible(line)
		}
	SelectLine(line /* zero-based */)
		{
		first = .PositionFromLine(line)
		last = .PositionFromLine(line + 1)
		.SetSelect(first, last - first)
		}
	SetWrap(wrap = true)
		{
		.SetWrapMode(wrap ? SC.WRAP_WORD : SC.WRAP_NONE)
		}
	PointFromPosition(pos)
		{
		return Object(
			x: .PointXFromPosition(0, pos)
			y: .PointYFromPosition(0, pos))
		}
	StyleGetFont(style)
		{
		return SendMessageTextOut(.Hwnd, SCI.STYLEGETFONT, style).text
		}
	StyleSetFont(style, font)
		{
		.SendMessageTextIn(SCI.STYLESETFONT, style, font)
		}

	On_Find()
		{
		// replaced by FindBar in LibView
		s = .GetSelText()
		if s > "" and not s.Has?('\n')
			.findreplacedata.find = s
		x = FindDialog(.findreplacedata)
		.DoFind(x)
		}

	DoFind(x)
		{
		if x is 'prev'
			return .On_Find_Previous()
		else if x is 'next'
			return .On_Find_Next()
		}

	On_Find_Next()
		{
		return .findAndMatch()
		}
	On_Find_Previous()
		{
		return .findAndMatch(prev:)
		}
	findAndMatch(prev = false)
		{
		getSelect = prev is false ? .GetSelect().cpMax : .GetSelect().cpMin
		if false is match =
			Find.DoFind(.SearchText(), getSelect, .findreplace_options, :prev)
			{
			.Send('UpdateOccurrence', num: 0, count: 0)
			return false
			}
		findOb = .numOfMatch(match)
		.Send('UpdateOccurrence', num: findOb.num, count: findOb.count)
		.SetVisibleSelect(match[0], match[1])
		return true
		}
	numOfMatch(match)
		{
		matches = Find.FindAll(.SearchText(), .findreplace_options)
		num = matches.FindIf({ it is match }) + 1
		return Object(:num, count: matches.Size())
		}

	SearchText()
		{
		return .Get()
		}

	astWriterManager: false
	Getter_AstWriterManager()
		{
		if .astWriterManager is false
			try .astWriterManager = AstWriteManager(.SearchText())
		return .astWriterManager
		}

	getter_findreplace_options()
		{
		.findreplacedata.ast = Find.NeedAst?(.findreplacedata) ? .AstWriterManager : false
		return .findreplacedata
		}

	MarkAll()
		{
		if false is matches = Find.FindAll(.SearchText(), .findreplace_options)
			return 0
		matches.Each()
			{ |m|
			.MarkerAdd(.LineFromPosition(m[0]), .marker_find)
			}
		return matches.Size()
		}
	ClearFindMarks()
		{
		.MarkerDeleteAll(.marker_find)
		}
	On_Find_Next_Selected()
		{
		.find_selected(.On_Find_Next)
		}
	On_Find_Prev_Selected()
		{
		.find_selected(.On_Find_Previous)
		}
	find_selected(nextprev)
		{
		sel = .GetSelect()
		.SelectCurrentWord()
		if "" is s = .GetSelText().Trim()
			{
			Beep()
			.SetSel(sel.cpMin, sel.cpMax)
			return
			}
		.findreplacedata.find = s
		.findreplacedata.case = true
		.findreplacedata.regex = false
		.findreplacedata.word = sel.cpMin is sel.cpMax
		.Send("FindBar_OpenFindBar")
		nextprev()
		}

	On_Replace()
		{
		s = .GetSelText()
		.findreplacedata.replaceIn = "Entire text"
		if s > ""
			if s.Has?('\n')
				.findreplacedata.replaceIn = "Selection"
			else
				.findreplacedata.find = s
		x = ReplaceDialog(this, .findreplacedata)
		if x is 'all'
			.ReplaceAll()
		}
	ReplaceOne()
		{
		sel = .GetSelect()
		size = .GetTextLength()
		if false is s = Find.DoReplace(.SearchText(), .GetSelText(), sel.cpMin, sel.cpMax,
			.findreplace_options)
			return false
		.ReplaceSel(s)
		if .findreplacedata.replaceIn is "Selection"
			{
			size_after = .GetTextLength()
			.SetSelect(sel.cpMin, sel.cpMax - sel.cpMin + (size_after - size))
			}
		return true
		}

	Paste(s)
		{
		.ReplaceSel(s)
		}
	ReplaceSel(s) // overload
		{
		s = EnsureCRLF(s)
		i = 0
		do
			{
			chunk = s[i :: .chunkSize]
			SendMessageTextIn(.Hwnd, SCI.REPLACESEL, 0, chunk)
			}
			while ((i += .chunkSize) < s.Size())
		}
	PasteOverAll(s)
		{
		.SelectAll()
		.Paste(s)
		.SetSelect(0)
		}
	ReplaceAll()
		{
		sel = .GetSelect()
		size = .GetTextLength()
		line = .GetFirstVisibleLine()
		if .findreplacedata.replaceIn isnt "Selection" // whole file
			.SelectAll()
		curSel = .GetSelect()
		if false is	s = Find.DoReplace(.SearchText(), .GetSelText(),
			curSel.cpMin, curSel.cpMax, .findreplace_options)
			return
		.ReplaceSel(s)
		.SetFirstVisibleLine(line)
		if (.findreplacedata.replaceIn is "Selection")
			{
			size_after = .GetTextLength()
			.SetSelect(sel.cpMin, sel.cpMax - sel.cpMin + (size_after - size))
			}
		}
	GetFlags()
		{
		return .get_marker_lines(1 << 1)
		}
	GetMarkerLines()
		{
		return .get_marker_lines(0xffff)
		}
	get_marker_lines(bits) // returns list of line numbers with flags
		{
		flags = Object()
		line = -1
		while 0 <= (line = .MarkerNext(line + 1, bits)) and
			line isnt 4294967295 /*= 32 bit -1 or 0xffffffff, for 64 bit gSuneido */
			flags.Add(line)
		return flags
		}
	GetMarkers()
		{
		markers = []
		line = -1
		while (0 <= (line = .MarkerNext(line + 1, 0xffff)))
			markers[line] = .MarkerGet(line)
		return markers
		}
	SetIndicator(indic, pos, len)
		{
		.SetIndicatorCurrent(indic)
		.IndicatorFillRange(pos, len)
		}
	ClearIndicator(indic, pos = 0, len = false)
		{
		if len is false
			len = .GetTextLength()
		.SetIndicatorCurrent(indic)
		.IndicatorClearRange(pos, len)
		}
	HasIndicator?(pos, indic)
		{
		if pos >= .GetTextLength()
			return false
		indics = .IndicatorAllOnFor(pos)
		return (indics & (1 << indic)) isnt 0
		}
	GetReadOnly()
		{
		return .readonly is true or .GETREADONLY() is 1
		}
	GetWindowRect()
		{
		return GetWindowRect(.Hwnd)
		}
	SetReadOnly(readonly = true)
		{
		if .readonly
			return
		if readonly and .Dirty?()
			{
// Extra logging for suggestion 25689
calls = GetCallStack(limit: 20)
RemoveAssertsFromCallStack(calls)
// avoiding broadcast/send loop
calls = FormatCallStack(calls, levels: 20).
	Replace('stdlib:Container.Broadcast:63', '').
	Replace('stdlib:Container.SetReadOnly:31', '').Lines().
	RemoveIf({ it.Blank?() }).Join('\n')
			SuneidoLog('ERROR: ScintillaControl SetReadonly(true) when dirty', :calls
				params: Object(dirty: .dirty, modify: .GetModify(),
					alreadyReadOnly: .GetReadOnly()))
			}

		return .SETREADONLY(readonly ? 1 : 0)
		}
	SetFirstVisibleLine(idx)
		{
		.LineScroll(0, idx - .GetFirstVisibleLine())
		}

	DoWithCurrentPos(block)
		{
		curLine = .GetFirstVisibleLine()
		block()
		.SetFirstVisibleLine(curLine)
		}

	Get()
		{
		return .GetRange(0, .GetTextLength())
		}
	setMethodModifying?: false // used in SCN_MODIFIED method, see comment there
	Set(s)
		{
		.astWriterManager = false
		.setMethodModifying? = true
		.ignoring_readonly()
			{
			.ClearAll()
			.append(String(s))
			}
		.EmptyUndoBuffer()
		.DocumentStart()
		.Dirty?(false)
		.setMethodModifying? = false
		}
	Trim() // remove whitespace from beginning and end
		{
		while .GetLength() > 0 and .GetAt(0).White?()
			.DeleteRange(0, 1)
		for (i = .GetLength() - 1; i >= 0 and .GetAt(i).White?(); --i)
			.DeleteRange(i, 1)
		return this // to allow chaining
		}
	AppendText(s)
		{
		.ignoring_readonly()
			{
			.append(s)
			}
		.DocumentEnd()
		}
	chunkSize: 8192
	append(s)
		{
		s = EnsureCRLF(s)
		for (i = 0; i < s.Size(); i += .chunkSize)
			{
			chunk = s[i :: .chunkSize]
			SendMessageTextIn(.Hwnd, SCI.APPENDTEXT, chunk.Size(), chunk)
			}
		}
	ignoring_readonly(block)
		{
		ro = .GETREADONLY()
		.SETREADONLY(0)
		block()
		.SETREADONLY(ro)
		}

	dirty: false
	Dirty?(dirty = "")
		{
		Assert(dirty is true or dirty is false or dirty is "")
		if dirty is false
			{
			.dirty = false
			.SetSavePoint()
			.UPDATEUI()
			}
		else if dirty is true
			{
			.dirty = true
			.UPDATEUI()
			}
		return (.dirty or (.GetModify() is 1))
		}
	UPDATEUI()
		{
		}

	On_Cut()
		{ .CUT() }
	On_Copy()
		{
		.COPY()
		s = ClipboardReadString()
		if not String?(s) or s.Size() is 0
			return
		s = s.Tr('\r').Replace('\n', '\r\n')
		ClipboardWriteString(s)
		}
	On_Paste()
		{ .PASTE() }
	On_Undo()
		{ .UNDO() }
	On_Redo()
		{ .REDO() }
	On_Delete()
		{ .CLEAR() }
	On_Select_All()
		{ .SelectAll() }

	LBUTTONDOWN(lParam)
		{
		.Send('Scintilla_LButtonDown', lParam)
		return 'callsuper'
		}
	RBUTTONDOWN(lParam)
		{
		select = .GetSelect()
		x = LOWORD(lParam)
		y = HIWORD(lParam)
		pos = .PositionFromPoint(x, y)
		if pos < select.cpMin or select.cpMax < pos
			.SetSelect(pos)
		SetFocus(.Hwnd)
		return 'callsuper'
		}
	Context_Menu: (
		"&Undo\tCtrl+Z", "&Redo\tCtrl+Y", "",
		"Cu&t\tCtrl+X", "&Copy\tCtrl+C", "&Paste\tCtrl+V", "&Delete", "",
		"Select &All\tCtrl+A", "Find...\tCtrl+F")
	ContextMenu(x, y)
		{
		ContextMenu(.Context_Menu).ShowCall(this, x, y)
		return 0
		}
	On_Context_Cut()
		{ .On_Cut() }
	On_Context_Copy()
		{ .On_Copy() }
	On_Context_Paste()
		{ .On_Paste() }
	On_Context_Undo()
		{ .On_Undo() }
	On_Context_Redo()
		{ .On_Redo() }
	On_Context_Delete()
		{ .On_Delete() }
	On_Context_Select_All()
		{ .On_Select_All() }
	On_Context_Find()
		{ .On_Find() }
	SelectCurrentWord()
		{
		// TODO handle cursor at start of word
		if .GetSelectionStart() < .GetSelectionEnd()
			return
		.WordLeft()
		.WordRightExtend()
		n = 0
		maxWordLength = 100
		while ((s = .GetSelText()) isnt '' and (s.Suffix?(' ') or s.Suffix?('\t')))
			{
			Assert(++n < maxWordLength)
			.CharLeftExtend()
			}
		}
	GetAt(pos)
		{
		return pos < 0 or pos >= .GetTextLength() ? '\x00' : .GetRange(pos, pos + 1)
		}

	FindReplaceData()
		{
		return .findreplacedata
		}

	GetState()
		{
		return Object(
			selection: .GetSelect()
			topline: .GetFirstVisibleLine()
			markers: .GetMarkers()
			linenums: .GetMarginWidthN(0)
			foldmargin: .GetMarginWidthN(2)
			)
		}
	SetState(state)
		{
		if state.Member?('selection')
			.SetSelect(state.selection.cpMin,
				state.selection.cpMax - state.selection.cpMin)
		if state.Member?('topline')
			.SetFirstVisibleLine(state.topline)
		if state.Member?('markers')
			for (line in state.markers.Members())
				.MarkerAddSet(line, state.markers[line])
		if state.Member?('linenums')
			.SetMarginWidthN(0, state.linenums)
		if state.Member?('foldmargin')
			.SetMarginWidthN(2, state.foldmargin)
		.MarkersChanged()
		}

	MarkersChanged()
		{
		.Send(#Scintilla_MarkersChanged)
		}

	GETDLGCODE(wParam, lParam)
		{
		if wParam is VK.ESCAPE and .AutocActive() is 0 or
			(.tabthrough and wParam is VK.TAB)
			return DLGC.WANTCHARS
		// prevent all the text from being selected when you tab into the editor
		return .Callsuper(.Hwnd, WM.GETDLGCODE, wParam, lParam) & ~DLGC.HASSETSEL
		}
	Destroy()
		{
		.Send("NoData")
		super.Destroy()
		}
	}
